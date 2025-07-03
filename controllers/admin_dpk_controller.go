package controllers

import (
	"net/http"
	"time"
	"strings"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"encoding/json"
	"github.com/Zauro25/Capstone-PerpusKominfosan/config"
	"github.com/Zauro25/Capstone-PerpusKominfosan/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetDashboardDPK(c *gin.Context) {
	// Get statistics
	var stats struct {
		TotalPerpustakaan int64
		TotalKoleksi      int64
		TotalPengunjung   int64
		TotalAnggota      int64
		TotalSDM          int64
		PendingVerifikasi int64
	}
	
	config.DB.Model(&models.Perpustakaan{}).Count(&stats.TotalPerpustakaan)
	config.DB.Model(&models.Koleksi{}).Count(&stats.TotalKoleksi)
	config.DB.Model(&models.Pengunjung{}).Count(&stats.TotalPengunjung)
	config.DB.Model(&models.Anggota{}).Where("status_anggota = ?", "Aktif").Count(&stats.TotalAnggota)
	config.DB.Model(&models.SDM{}).Count(&stats.TotalSDM)
	config.DB.Model(&models.Verifikasi{}).Where("status = ?", "Pending").Count(&stats.PendingVerifikasi)
	
	// Get recent activities
	var activities []models.AuditLog
	config.DB.Where("user_type = ?", "admin_perpustakaan").
		Order("timestamp desc").
		Limit(5).
		Find(&activities)
	
	// Get recent notifications
	var notifications []models.Notifikasi
	config.DB.Where("tujuan_user = ?", "admin_dpk").
		Order("tanggal_kirim desc").
		Limit(5).
		Find(&notifications)
	
	c.JSON(http.StatusOK, gin.H{
		"statistics":    stats,
		"activities":    activities,
		"notifications": notifications,
	})
}

func GetAllPerpustakaan(c *gin.Context) {
	var perpustakaan []models.Perpustakaan
	
	query := config.DB
	
	// Filter by status if provided
	if status := c.Query("status"); status != "" {
		query = query.Where("status_verifikasi = ?", status)
	}
	
	// Search by name if provided
	if search := c.Query("search"); search != "" {
		query = query.Where("nama_perpustakaan ILIKE ?", "%"+search+"%")
	}
	
	if err := query.Find(&perpustakaan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data perpustakaan"})
		return
	}
	
	c.JSON(http.StatusOK, perpustakaan)
}

func GetDetailPerpustakaan(c *gin.Context) {
	id := c.Param("id")
	
	var perpustakaan models.Perpustakaan
	if err := config.DB.Preload("Koleksi").
		Preload("SDM").
		Preload("Pengunjung").
		Preload("Anggota").
		Preload("AdminPerpustakaan").
		First(&perpustakaan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Perpustakaan tidak ditemukan"})
		return
	}
	
	c.JSON(http.StatusOK, perpustakaan)
}
func VerifikasiAdminPerpustakaan(c *gin.Context) {
    adminDPKID := c.GetUint("user_id") // ID admin DPK yang melakukan verifikasi
    
    var req models.VerifyAdminRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Start transaction
    tx := config.DB.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // Get admin perpustakaan data with preloaded perpustakaan
    var adminPerpus models.AdminPerpustakaan
    if err := tx.Preload("Perpustakaan").First(&adminPerpus, req.AdminPerpustakaanID).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusNotFound, gin.H{"error": "Admin perpustakaan tidak ditemukan"})
        return
    }

    // Check if already verified
    if adminPerpus.IsActive && req.Status == "approved" {
        tx.Rollback()
        c.JSON(http.StatusBadRequest, gin.H{"error": "Admin perpustakaan sudah diverifikasi sebelumnya"})
        return
    }

    // Update admin status
    isApproved := req.Status == "approved"
    adminPerpus.IsActive = isApproved
    if err := tx.Save(&adminPerpus).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui status admin"})
        return
    }

    // Update perpustakaan verification status if approved
    if isApproved {
        adminPerpus.Perpustakaan.StatusVerifikasi = "Verified"
        if err := tx.Save(&adminPerpus.Perpustakaan).Error; err != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui status perpustakaan"})
            return
        }
    }

    // Create verification record
    verifikasi := models.Verifikasi{
        PerpustakaanID:    adminPerpus.PerpustakaanID,
        JenisData:        "AdminPerpustakaan",
        Status:           req.Status,
        CatatanRevisi:    req.Catatan,
		TanggalVerifikasi: func(t time.Time) *time.Time { return &t }(time.Now()),
        AdminDPKID:      adminDPKID,
    }
    if err := tx.Create(&verifikasi).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat record verifikasi"})
        return
    }

    // Create notification for admin perpustakaan
    judul := "Akun Anda Telah Diverifikasi"
    isi := "Akun admin perpustakaan Anda telah diverifikasi dan diaktifkan."
    if !isApproved {
        judul = "Verifikasi Akun Ditolak"
        isi = "Verifikasi akun admin perpustakaan Anda ditolak. Catatan: " + req.Catatan
    }

    notifikasi := models.Notifikasi{
        JudulNotifikasi: judul,
        IsiNotifikasi:   isi,
        JenisNotifikasi: "info",
        TujuanUser:      "admin_perpustakaan",
        UserID:         &adminPerpus.ID,
        RelatedID:      &adminPerpus.ID,
        RelatedType:    "admin_perpustakaan",
        IsRead:         false,
        TanggalKirim:   time.Now(),
        ExpiredAt:      time.Now().Add(7 * 24 * time.Hour),
        ActionLink:     "/admin-perpustakaan/profile",
    }
    if err := tx.Create(&notifikasi).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat notifikasi"})
        return
    }

    // Create audit log
    auditLog := models.AuditLog{
        UserType:  "admin_dpk",
        UserID:    adminDPKID,
        Action:    "VERIFY_ADMIN_PERUSTAKAAN",
        TableName: "admin_perpustakaan",
        RecordID:  adminPerpus.ID,
        NewValues: stringifyMap(map[string]interface{}{
            "status":        req.Status,
            "is_active":     isApproved,
            "admin_dpk_id":  adminDPKID,
            "catatan":       req.Catatan,
        }),
        IPAddress: c.ClientIP(),
        UserAgent: c.GetHeader("User-Agent"),
        Timestamp: time.Now(),
    }
    if err := tx.Create(&auditLog).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat audit log"})
        return
    }

    // Commit transaction
    if err := tx.Commit().Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyelesaikan transaksi"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Verifikasi admin perpustakaan berhasil",
        "status":  req.Status,
        "admin": gin.H{
            "id": adminPerpus.ID,
            "username": adminPerpus.Username,
            "is_active": adminPerpus.IsActive,
        },
        "perpustakaan": gin.H{
            "id": adminPerpus.Perpustakaan.ID,
            "nama_perpustakaan": adminPerpus.Perpustakaan.NamaPerpustakaan,
            "status_verifikasi": adminPerpus.Perpustakaan.StatusVerifikasi,
        },
    })
}
func VerifikasiData(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var req models.VerifikasiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	var perpustakaan models.Perpustakaan
	if err := config.DB.First(&perpustakaan, req.PerpustakaanID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Perpustakaan tidak ditemukan"})
		return
	}
	
	if perpustakaan.StatusVerifikasi != "Terkirim" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data belum dikirim untuk verifikasi"})
		return
	}
	
	now := time.Now()
	
	// Update perpustakaan status
	if err := config.DB.Model(&perpustakaan).Updates(map[string]interface{}{
		"status_verifikasi": req.Status,
		"catatan_revisi":    req.CatatanRevisi,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui status perpustakaan"})
		return
	}
	
	// Update verification record
	var verifikasi models.Verifikasi
	if err := config.DB.Where("perpustakaan_id = ? AND status = ?", req.PerpustakaanID, "Pending").
		First(&verifikasi).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record verifikasi tidak ditemukan"})
		return
	}
	
	verifikasi.Status = req.Status
	verifikasi.CatatanRevisi = req.CatatanRevisi
	verifikasi.TanggalVerifikasi = &now
	verifikasi.AdminDPKID = userID
	
	if err := config.DB.Save(&verifikasi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui record verifikasi"})
		return
	}
	
	// Create log revisi if status is Direvisi
	if req.Status == "Direvisi" {
		logRevisi := models.LogRevisi{
			PerpustakaanID: req.PerpustakaanID,
			JenisRevisi:    "Perpustakaan",
			CatatanRevisi:  req.CatatanRevisi,
			StatusSebelum:  "Terkirim",
			StatusSesudah:  "Direvisi",
			AdminDPKID:     userID,
			TanggalRevisi:  now,
		}
		
		if err := config.DB.Create(&logRevisi).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat log revisi"})
			return
		}
	}
	
	// Create notification for perpustakaan admin
	judul := "Data Disetujui"
	if req.Status == "Direvisi" {
		judul = "Data Perlu Revisi"
	}
	
	notifikasi := models.Notifikasi{
		JudulNotifikasi: judul,
		IsiNotifikasi:   "Data perpustakaan " + perpustakaan.NamaPerpustakaan + " telah diverifikasi dengan status: " + req.Status,
		JenisNotifikasi: "Info",
		TujuanUser:      "admin_perpustakaan",
		UserID:         &perpustakaan.ID,
		TanggalKirim:    now,
	}
	if err := config.DB.Create(&notifikasi).Error; err != nil {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan notifikasi"})
	return
}
	if err := notificationService.SendVerificationNotification(req.PerpustakaanID, req.Status, req.CatatanRevisi); err != nil {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengirim notifikasi"})
	return
}
	
	// Create audit log
	auditLog := models.AuditLog{
		UserType:  "admin_dpk",
		UserID:    userID,
		Action:    "VERIFIKASI_DATA",
		TableName: "perpustakaan",
		RecordID:  req.PerpustakaanID,
		NewValues: stringifyMap(map[string]interface{}{
			"status_verifikasi": req.Status,
			"catatan_revisi":    req.CatatanRevisi,
		}),
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
		Timestamp: now,
	}
	config.DB.Create(&auditLog)
	
	c.JSON(http.StatusOK, gin.H{"message": "Verifikasi data berhasil disimpan"})
}

func GetPendingVerification(c *gin.Context) {
	var verifikasi []models.Verifikasi
	
	if err := config.DB.Preload("Perpustakaan").
		Preload("AdminDPK").
		Where("status = ?", "Pending").
		Order("tanggal_verifikasi asc").
		Find(&verifikasi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data verifikasi"})
		return
	}
	
	c.JSON(http.StatusOK, verifikasi)
}

func GenerateLaporan(c *gin.Context) {
	adminDPKID := c.GetUint("user_id") // ID admin DPK yang membuat laporan
    var req models.LaporanRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        fmt.Println("Binding error:", err) // Add logging
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    fmt.Println("Request body:", req)

    // 1. Ambil data statistik
    var statistik struct {
        TotalPerpustakaan int64
        TotalKoleksi      int64
        TotalPengunjung   int64
        TotalAnggota      int64
        TotalSDM          int64
        PerpustakaanByJenis []struct {
            Jenis  string
            Jumlah int64
        }
        KunjunganBulanan []struct {
            Bulan string
            Jumlah int64
        }
    }

    // Hitung total data
    config.DB.Model(&models.Perpustakaan{}).Count(&statistik.TotalPerpustakaan)
    config.DB.Model(&models.Koleksi{}).Count(&statistik.TotalKoleksi)
    config.DB.Model(&models.Pengunjung{}).Count(&statistik.TotalPengunjung)
    config.DB.Model(&models.Anggota{}).Count(&statistik.TotalAnggota)
    config.DB.Model(&models.SDM{}).Count(&statistik.TotalSDM)

    // Data per jenis perpustakaan
    config.DB.Model(&models.Perpustakaan{}).
        Select("jenis_perpustakaan as jenis, count(*) as jumlah").
        Group("jenis_perpustakaan").
        Scan(&statistik.PerpustakaanByJenis)

    // Data kunjungan bulanan (contoh)
    config.DB.Raw(`
        SELECT 
            to_char(tanggal_kunjungan, 'YYYY-MM') as bulan,
            sum(jumlah_pengunjung) as jumlah
        FROM pengunjungs
        GROUP BY bulan
        ORDER BY bulan
    `).Scan(&statistik.KunjunganBulanan)

    // 2. Siapkan data visualisasi
    chartData := gin.H{
        "total_data": gin.H{
            "perpustakaan": statistik.TotalPerpustakaan,
            "koleksi":      statistik.TotalKoleksi,
            "pengunjung":   statistik.TotalPengunjung,
            "anggota":      statistik.TotalAnggota,
            "sdm":          statistik.TotalSDM,
        },
        "jenis_perpustakaan": gin.H{
            "labels":   []string{},
            "datasets": []gin.H{{
                "data":            []int{},
                "backgroundColor": []string{"#FF6384", "#36A2EB", "#FFCE56"},
            }},
        },
        "kunjungan_bulanan": gin.H{
            "labels":   []string{},
            "datasets": []gin.H{{
                "label":           "Jumlah Pengunjung",
                "data":            []int{},
                "backgroundColor": "rgba(75, 192, 192, 0.2)",
                "borderColor":     "rgba(75, 192, 192, 1)",
            }},
        },
    }

    // Isi data chart
    for _, j := range statistik.PerpustakaanByJenis {
        chartData["jenis_perpustakaan"].(gin.H)["labels"] = append(
            chartData["jenis_perpustakaan"].(gin.H)["labels"].([]string),
            j.Jenis,
        )
        chartData["jenis_perpustakaan"].(gin.H)["datasets"].([]gin.H)[0]["data"] = append(
            chartData["jenis_perpustakaan"].(gin.H)["datasets"].([]gin.H)[0]["data"].([]int),
            int(j.Jumlah),
        )
    }

    for _, k := range statistik.KunjunganBulanan {
        chartData["kunjungan_bulanan"].(gin.H)["labels"] = append(
            chartData["kunjungan_bulanan"].(gin.H)["labels"].([]string),
            k.Bulan,
        )
        chartData["kunjungan_bulanan"].(gin.H)["datasets"].([]gin.H)[0]["data"] = append(
            chartData["kunjungan_bulanan"].(gin.H)["datasets"].([]gin.H)[0]["data"].([]int),
            int(k.Jumlah),
        )
    }

    // 3. Generate file laporan
    fileName := fmt.Sprintf("laporan_%s_%s.%s", req.JenisLaporan, time.Now().Format("20060102"), req.FormatFile)
    filePath := filepath.Join("storage", "reports", fileName)

    if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat direktori penyimpanan"})
        return
    }

    var fileContent string
	switch req.FormatFile {
	case "csv":
		fileContent = generateCSVReport(statistik, chartData)
	case "json":
		fileContent = generateJSONReport(statistik, chartData)
	case "pdf":
		fileContent = generatePDFReport(statistik, chartData) // Implementasi khusus
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format laporan tidak didukung"})
		return
	}

    if err := os.WriteFile(filePath, []byte(fileContent), 0644); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan laporan"})
        return
    }

    // 4. Simpan record laporan
    chartDataJSON, _ := json.Marshal(chartData)
    laporan := models.Laporan{
        Judul:          req.Judul,
        Periode:        req.Periode,
        JenisLaporan:   req.JenisLaporan,
        FilePath:       filePath,
        FormatFile:     req.FormatFile,
        ChartData:      string(chartDataJSON),
        TanggalGenerate: time.Now(),
        AdminDPKID:     adminDPKID,
        Status:         "Generated",
    }

    if err := config.DB.Create(&laporan).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data laporan"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Laporan berhasil dibuat",
        "data": gin.H{
            "laporan": laporan,
            "chart_data": chartData,
        },
        "download_url": fmt.Sprintf("/api/v1/laporan/%d/download", laporan.ID),
    })
}

// Fungsi pembantu untuk generate report
func generateCSVReport(data interface{}, chartData gin.H) string {
    var content strings.Builder
    content.WriteString("Total Perpustakaan,Total Koleksi,Total Pengunjung,Total Anggota,Total SDM,Chart Data\n")
    stats := data.(struct {
        TotalPerpustakaan, TotalKoleksi, TotalPengunjung, TotalAnggota, TotalSDM int64
        PerpustakaanByJenis []struct{ Jenis string; Jumlah int64 }
        KunjunganBulanan []struct{ Bulan string; Jumlah int64 }
    })
    chartJSON, _ := json.Marshal(chartData)
    content.WriteString(fmt.Sprintf("%d,%d,%d,%d,%d,%q\n",
        stats.TotalPerpustakaan, stats.TotalKoleksi, stats.TotalPengunjung,
        stats.TotalAnggota, stats.TotalSDM, string(chartJSON)))
    return content.String()
}

func generateJSONReport(data interface{}, chartData gin.H) string {
    report := gin.H{
        "statistik": data,
        "visualisasi": chartData,
    }
    jsonData, _ := json.MarshalIndent(report, "", "  ")
    return string(jsonData)
}

func generatePDFReport(data interface{}, chartData gin.H) string {
    // Convert chartData to JSON
    chartJSON, _ := json.Marshal(chartData)

    // Create a temporary HTML file for chart rendering
    htmlContent := `
    <html>
    <body>
        <h1>Laporan Statistik</h1>
        <div id="chart"></div>
        <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
        <script>
            const data = ` + string(chartJSON) + `;
            new Chart(document.getElementById('chart'), {
                type: 'bar',
                data: {
                    labels: data.jenis_perpustakaan.labels,
                    datasets: data.jenis_perpustakaan.datasets
                }
            });
        </script>
    </body>
    </html>`

    // Save HTML temporarily
    htmlPath := "temp_chart.html"
    if err := os.WriteFile(htmlPath, []byte(htmlContent), 0644); err != nil {
        return ""
    }

    // Use wkhtmltopdf to convert HTML to PDF
    pdfPath := filepath.Join("storage", "reports", fmt.Sprintf("laporan_%s.pdf", time.Now().Format("20060102")))
    cmd := exec.Command("wkhtmltopdf", htmlPath, pdfPath)
    if err := cmd.Run(); err != nil {
        return ""
    }

    // Clean up temporary HTML
    os.Remove(htmlPath)
    return pdfPath
}
func DownloadLaporan(c *gin.Context) {
	id := c.Param("id")
	
	var laporan models.Laporan
	if err := config.DB.First(&laporan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Laporan tidak ditemukan"})
		return
	}
	
	// Check if file exists
	if _, err := os.Stat(laporan.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File laporan tidak ditemukan"})
		return
	}
	
	// Set headers for download
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(laporan.FilePath)))
	c.Header("Content-Type", "application/octet-stream")
	
	// Serve the file
	c.File(laporan.FilePath)
}
func generateReportContent(perpustakaans []models.Perpustakaan) string {
    // Example report generation
    var content strings.Builder
    content.WriteString("Laporan Data Perpustakaan\n")
    content.WriteString("=========================\n\n")
    
    for _, p := range perpustakaans {
        content.WriteString(fmt.Sprintf("Nama: %s\n", p.NamaPerpustakaan))
        content.WriteString(fmt.Sprintf("Alamat: %s\n", p.Alamat))
        content.WriteString(fmt.Sprintf("Jumlah SDM: %d\n", p.JumlahSDM))
        content.WriteString(fmt.Sprintf("Jumlah Pengunjung: %d\n", p.JumlahPengunjung))
        content.WriteString(fmt.Sprintf("Jumlah Anggota: %d\n", p.JumlahAnggota))
        content.WriteString("-------------------------\n")
    }
    
    return content.String()
}
func calculateStatistics(perpustakaans []models.Perpustakaan) map[string]interface{} {
    stats := make(map[string]interface{})
    
    // Calculate totals
    var totalPerpustakaan, totalSDM, totalPengunjung, totalAnggota int
    
    for _, p := range perpustakaans {
        totalPerpustakaan++
        totalSDM += p.JumlahSDM
        totalPengunjung += p.JumlahPengunjung
        totalAnggota += p.JumlahAnggota
    }
    
    // Calculate averages
    avgSDM := float64(totalSDM) / float64(totalPerpustakaan)
    avgPengunjung := float64(totalPengunjung) / float64(totalPerpustakaan)
    avgAnggota := float64(totalAnggota) / float64(totalPerpustakaan)
    
    // Group by jenis perpustakaan
    jenisStats := make(map[string]map[string]interface{})
    for _, p := range perpustakaans {
        if _, ok := jenisStats[p.JenisPerpustakaan]; !ok {
            jenisStats[p.JenisPerpustakaan] = make(map[string]interface{})
            jenisStats[p.JenisPerpustakaan]["count"] = 0
            jenisStats[p.JenisPerpustakaan]["sdm"] = 0
            jenisStats[p.JenisPerpustakaan]["pengunjung"] = 0
            jenisStats[p.JenisPerpustakaan]["anggota"] = 0
        }
        
        jenisStats[p.JenisPerpustakaan]["count"] = jenisStats[p.JenisPerpustakaan]["count"].(int) + 1
        jenisStats[p.JenisPerpustakaan]["sdm"] = jenisStats[p.JenisPerpustakaan]["sdm"].(int) + p.JumlahSDM
        jenisStats[p.JenisPerpustakaan]["pengunjung"] = jenisStats[p.JenisPerpustakaan]["pengunjung"].(int) + p.JumlahPengunjung
        jenisStats[p.JenisPerpustakaan]["anggota"] = jenisStats[p.JenisPerpustakaan]["anggota"].(int) + p.JumlahAnggota
    }
    
    stats["total_perpustakaan"] = totalPerpustakaan
    stats["total_sdm"] = totalSDM
    stats["total_pengunjung"] = totalPengunjung
    stats["total_anggota"] = totalAnggota
    stats["avg_sdm"] = avgSDM
    stats["avg_pengunjung"] = avgPengunjung
    stats["avg_anggota"] = avgAnggota
    stats["by_jenis"] = jenisStats
    
    return stats
}

func generateReportFile(format string, stats map[string]interface{}) string {
    // In a real implementation, this would generate an actual file
    // For simulation, just return a path
    timestamp := time.Now().Format("20060102_150405")
    return "/reports/statistik_" + timestamp + "." + strings.ToLower(format)
}
func SendBroadcastNotification(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var req models.NotifikasiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Create notification for all perpustakaan admins
	notifikasi := models.Notifikasi{
		JudulNotifikasi: req.JudulNotifikasi,
		IsiNotifikasi:   req.IsiNotifikasi,
		JenisNotifikasi: req.JenisNotifikasi,
		TujuanUser:      "admin_perpustakaan",
		TanggalKirim:    time.Now(),
	}
	
	if err := config.DB.Create(&notifikasi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengirim notifikasi"})
		return
	}
	
	// Create audit log
	auditLog := models.AuditLog{
		UserType:  "admin_dpk",
		UserID:    userID,
		Action:    "BROADCAST_NOTIFICATION",
		TableName: "notifikasi",
		RecordID:  notifikasi.ID,
		NewValues: stringifyMap(map[string]interface{}{
			"judul_notifikasi": req.JudulNotifikasi,
			"isi_notifikasi":   req.IsiNotifikasi,
			"jenis_notifikasi": req.JenisNotifikasi,
		}),
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
		Timestamp: time.Now(),
	}
	config.DB.Create(&auditLog)
	
	c.JSON(http.StatusOK, gin.H{"message": "Notifikasi berhasil dikirim ke semua admin perpustakaan"})
}

func GetAuditLogs(c *gin.Context) {
	var logs []models.AuditLog
	
	query := config.DB
	
	// Filter by user type if provided
	if userType := c.Query("user_type"); userType != "" {
		query = query.Where("user_type = ?", userType)
	}
	
	// Filter by action if provided
	if action := c.Query("action"); action != "" {
		query = query.Where("action LIKE ?", "%"+action+"%")
	}
	
	// Filter by date range if provided
	if startDate := c.Query("start_date"); startDate != "" {
		if endDate := c.Query("end_date"); endDate != "" {
			query = query.Where("timestamp BETWEEN ? AND ?", startDate, endDate)
		} else {
			query = query.Where("timestamp >= ?", startDate)
		}
	}
	
	if err := query.Order("timestamp desc").Limit(100).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil log audit"})
		return
	}
	
	c.JSON(http.StatusOK, logs)
}

func GetAllAdminPerpustakaan(c *gin.Context) {
	var admins []models.AdminPerpustakaan
	
	query := config.DB.Preload("Perpustakaan")
	
	// Filter by active status if provided
	if active := c.Query("active"); active != "" {
		isActive := active == "true"
		query = query.Where("is_active = ?", isActive)
	}
	
	// Search by name if provided
	if search := c.Query("search"); search != "" {
		query = query.Where("nama_lengkap ILIKE ? OR username ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	
	if err := query.Find(&admins).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data admin perpustakaan"})
		return
	}
	
	c.JSON(http.StatusOK, admins)
}

func CreateAdminPerpustakaan(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var req struct {
		Username       string `json:"username" binding:"required"`
		Password       string `json:"password" binding:"required,min=6"`
		NamaLengkap    string `json:"nama_lengkap" binding:"required"`
		Email          string `json:"email" binding:"required,email"`
		NoTelepon      string `json:"no_telepon"`
		PerpustakaanID uint   `json:"perpustakaan_id" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Check if username already exists
	var existingAdmin models.AdminPerpustakaan
	if err := config.DB.Where("username = ?", req.Username).First(&existingAdmin).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username sudah digunakan"})
		return
	}
	
	// Check if email already exists
	if err := config.DB.Where("email = ?", req.Email).First(&existingAdmin).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah digunakan"})
		return
	}
	
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi password"})
		return
	}
	
	// Create admin
	admin := models.AdminPerpustakaan{
		Username:       req.Username,
		Password:       string(hashedPassword),
		NamaLengkap:    req.NamaLengkap,
		Email:          req.Email,
		NoTelepon:      req.NoTelepon,
		PerpustakaanID: req.PerpustakaanID,
		IsActive:       true,
	}
	
	if err := config.DB.Create(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat admin perpustakaan"})
		return
	}
	
	// Create audit log
	auditLog := models.AuditLog{
		UserType:  "admin_dpk",
		UserID:    userID,
		Action:    "CREATE_ADMIN_PERUSTAKAAN",
		TableName: "admin_perpustakaan",
		RecordID:  admin.ID,
		NewValues: stringifyMap(map[string]interface{}{
			"username":        req.Username,
			"nama_lengkap":    req.NamaLengkap,
			"email":           req.Email,
			"perpustakaan_id": req.PerpustakaanID,
		}),
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
		Timestamp: time.Now(),
	}
	config.DB.Create(&auditLog)
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Admin perpustakaan berhasil dibuat",
		"admin":   admin,
	})
}

func UpdateAdminPerpustakaan(c *gin.Context) {
	userID := c.GetUint("user_id")
	adminID := c.Param("id")
	
	var req struct {
		NamaLengkap string `json:"nama_lengkap"`
		Email       string `json:"email" binding:"email"`
		NoTelepon   string `json:"no_telepon"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	var admin models.AdminPerpustakaan
	if err := config.DB.Preload("Perpustakaan").First(&admin, adminID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin perpustakaan tidak ditemukan"})
		return
	}
	
	// Check if email already exists (if changed)
	if req.Email != "" && req.Email != admin.Email {
		var existingAdmin models.AdminPerpustakaan
		if err := config.DB.Where("email = ?", req.Email).First(&existingAdmin).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah digunakan"})
			return
		}
	}
	
	// Update fields
	updateData := map[string]interface{}{}
	if req.NamaLengkap != "" {
		updateData["nama_lengkap"] = req.NamaLengkap
	}
	if req.Email != "" {
		updateData["email"] = req.Email
	}
	if req.NoTelepon != "" {
		updateData["no_telepon"] = req.NoTelepon
	}
	
	if len(updateData) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tidak ada data yang diperbarui"})
		return
	}
	
	if err := config.DB.Model(&admin).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui admin perpustakaan"})
		return
	}
	
	// Create audit log
	auditLog := models.AuditLog{
		UserType:  "admin_dpk",
		UserID:    userID,
		Action:    "UPDATE_ADMIN_PERUSTAKAAN",
		TableName: "admin_perpustakaan",
		RecordID:  admin.ID,
		NewValues: stringifyMap(updateData),
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
		Timestamp: time.Now(),
	}
	config.DB.Create(&auditLog)
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Admin perpustakaan berhasil diperbarui",
		"admin":   admin,
	})
}

func DeleteAdminPerpustakaan(c *gin.Context) {
	userID := c.GetUint("user_id")
	adminID := c.Param("id")
	
	var admin models.AdminPerpustakaan
	if err := config.DB.First(&admin, adminID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin perpustakaan tidak ditemukan"})
		return
	}
	
	if err := config.DB.Delete(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus admin perpustakaan"})
		return
	}
	
	// Create audit log
	auditLog := models.AuditLog{
		UserType:  "admin_dpk",
		UserID:    userID,
		Action:    "DELETE_ADMIN_PERUSTAKAAN",
		TableName: "admin_perpustakaan",
		RecordID:  admin.ID,
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
		Timestamp: time.Now(),
	}
	config.DB.Create(&auditLog)
	
	c.JSON(http.StatusOK, gin.H{"message": "Admin perpustakaan berhasil dihapus"})
}

func ActivateAdminPerpustakaan(c *gin.Context) {
	userID := c.GetUint("user_id")
	adminID := c.Param("id")
	
	var admin models.AdminPerpustakaan
	if err := config.DB.First(&admin, adminID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin perpustakaan tidak ditemukan"})
		return
	}
	
	if admin.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Admin sudah aktif"})
		return
	}
	
	admin.IsActive = true
	if err := config.DB.Save(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengaktifkan admin perpustakaan"})
		return
	}
	
	// Create audit log
	auditLog := models.AuditLog{
		UserType:  "admin_dpk",
		UserID:    userID,
		Action:    "ACTIVATE_ADMIN_PERUSTAKAAN",
		TableName: "admin_perpustakaan",
		RecordID:  admin.ID,
		NewValues: "is_active: true",
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
		Timestamp: time.Now(),
	}
	config.DB.Create(&auditLog)
	
	c.JSON(http.StatusOK, gin.H{"message": "Admin perpustakaan berhasil diaktifkan"})
}

func DeactivateAdminPerpustakaan(c *gin.Context) {
	userID := c.GetUint("user_id")
	adminID := c.Param("id")
	
	var admin models.AdminPerpustakaan
	if err := config.DB.First(&admin, adminID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin perpustakaan tidak ditemukan"})
		return
	}
	
	if !admin.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Admin sudah tidak aktif"})
		return
	}
	
	admin.IsActive = false
	if err := config.DB.Save(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menonaktifkan admin perpustakaan"})
		return
	}
	
	// Create audit log
	auditLog := models.AuditLog{
		UserType:  "admin_dpk",
		UserID:    userID,
		Action:    "DEACTIVATE_ADMIN_PERUSTAKAAN",
		TableName: "admin_perpustakaan",
		RecordID:  admin.ID,
		NewValues: "is_active: false",
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
		Timestamp: time.Now(),
	}
	config.DB.Create(&auditLog)
	
	c.JSON(http.StatusOK, gin.H{"message": "Admin perpustakaan berhasil dinonaktifkan"})
}

func ResetPasswordAdminPerpustakaan(c *gin.Context) {
	userID := c.GetUint("user_id")
	adminID := c.Param("id")
	
	var req struct {
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	var admin models.AdminPerpustakaan
	if err := config.DB.First(&admin, adminID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin perpustakaan tidak ditemukan"})
		return
	}
	
	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi password"})
		return
	}
	
	admin.Password = string(hashedPassword)
	if err := config.DB.Save(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan password baru"})
		return
	}
	
	// Create audit log
	auditLog := models.AuditLog{
		UserType:  "admin_dpk",
		UserID:    userID,
		Action:    "RESET_PASSWORD_ADMIN_PERUSTAKAAN",
		TableName: "admin_perpustakaan",
		RecordID:  admin.ID,
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
		Timestamp: time.Now(),
	}
	config.DB.Create(&auditLog)
	
	c.JSON(http.StatusOK, gin.H{"message": "Password admin perpustakaan berhasil direset"})
}
func SendReminder(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var req struct {
		PerpustakaanID uint   `json:"perpustakaan_id" binding:"required"`
		Message        string `json:"message" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := notificationService.SendReminderNotification(req.PerpustakaanID, req.Message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengirim pengingat"})
		return
	}
	
	// Create audit log
	auditLog := models.AuditLog{
		UserType:  "admin_dpk",
		UserID:    userID,
		Action:    "SEND_REMINDER",
		TableName: "notifikasi",
		RecordID:  req.PerpustakaanID,
		NewValues: req.Message,
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
		Timestamp: time.Now(),
	}
	config.DB.Create(&auditLog)
	
	c.JSON(http.StatusOK, gin.H{"message": "Pengingat berhasil dikirim"})
}
func GetStatistics(c *gin.Context) {
    // Statistik per jenis perpustakaan
    var statsByType []struct {
        JenisPerpustakaan string
        Total             int
        TotalKoleksi      int
        TotalPengunjung   int
        TotalAnggota      int
    }

    if err := config.DB.Model(&models.Perpustakaan{}).
        Select("jenis_perpustakaan, COUNT(*) as total, SUM(jumlah_koleksi) as total_koleksi, SUM(jumlah_pengunjung) as total_pengunjung, SUM(jumlah_anggota) as total_anggota").
        Group("jenis_perpustakaan").
        Scan(&statsByType).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil statistik"})
        return
    }

    // Statistik perkembangan bulanan
    var monthlyStats []struct {
        Bulan    string
        Kunjungan int
        Anggota   int
    }

    currentYear := time.Now().Year()
    if err := config.DB.Raw(`
        SELECT 
            TO_CHAR(tanggal_kunjungan, 'YYYY-MM') as bulan,
            SUM(jumlah_pengunjung) as kunjungan,
            COUNT(DISTINCT anggota.id) as anggota
        FROM pengunjung
        LEFT JOIN anggota ON anggota.perpustakaan_id = pengunjung.perpustakaan_id
        WHERE EXTRACT(YEAR FROM tanggal_kunjungan) = ?
        GROUP BY TO_CHAR(tanggal_kunjungan, 'YYYY-MM')
        ORDER BY bulan
    `, currentYear).Scan(&monthlyStats).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil statistik bulanan"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "by_type": statsByType,
        "monthly": monthlyStats,
    })
}