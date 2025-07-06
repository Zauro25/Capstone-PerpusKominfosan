package controllers

import (
	"net/http"
	"time"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/Zauro25/Capstone-PerpusKominfosan/config"
	"github.com/Zauro25/Capstone-PerpusKominfosan/models"
)

func GetDashboardPerpustakaan(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var adminPerpus models.AdminPerpustakaan
	if err := config.DB.Preload("Perpustakaan").First(&adminPerpus, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin perpustakaan tidak ditemukan"})
		return
	}

	perpustakaanID := adminPerpus.PerpustakaanID
	
	// Get statistics
	var stats struct {
		TotalKoleksi    int64
		TotalPengunjung int64
		TotalAnggota    int64
		TotalSDM        int64
	}
	
	config.DB.Model(&models.Koleksi{}).Where("perpustakaan_id = ?", perpustakaanID).Count(&stats.TotalKoleksi)
	config.DB.Model(&models.Pengunjung{}).Where("perpustakaan_id = ?", perpustakaanID).Count(&stats.TotalPengunjung)
	config.DB.Model(&models.Anggota{}).Where("perpustakaan_id = ? AND status_anggota = ?", perpustakaanID, "Aktif").Count(&stats.TotalAnggota)
	config.DB.Model(&models.SDM{}).Where("perpustakaan_id = ?", perpustakaanID).Count(&stats.TotalSDM)
	
	// Get recent activities
	var activities []models.AuditLog
	config.DB.Where("user_type = ? AND user_id = ?", "admin_perpustakaan", userID).
		Order("timestamp desc").
		Limit(5).
		Find(&activities)
	
	// Get notifications
	var notifications []models.Notifikasi
	config.DB.Where("tujuan_user = ? AND (user_id = ? OR user_id IS NULL)", "admin_perpustakaan", userID).
		Order("tanggal_kirim desc").
		Limit(5).
		Find(&notifications)
	
	c.JSON(http.StatusOK, gin.H{
		"perpustakaan": adminPerpus.Perpustakaan,
		"statistics":   stats,
		"activities":   activities,
		"notifications": notifications,
	})
}

func GetDataPerpustakaan(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var adminPerpus models.AdminPerpustakaan
	if err := config.DB.Preload("Perpustakaan").First(&adminPerpus, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin perpustakaan tidak ditemukan"})
		return
	}
	
	perpustakaanID := adminPerpus.PerpustakaanID
	
	var perpustakaan models.Perpustakaan
	if err := config.DB.Preload("Koleksi").
		Preload("SDM").
		Preload("Pengunjung").
		Preload("Anggota").
		First(&perpustakaan, perpustakaanID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data perpustakaan tidak ditemukan"})
		return
	}
	
	c.JSON(http.StatusOK, perpustakaan)
}

func UpdateDataPerpustakaan(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var adminPerpus models.AdminPerpustakaan
	if err := config.DB.Preload("Perpustakaan").First(&adminPerpus, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin perpustakaan tidak ditemukan"})
		return
	}
	
	perpustakaanID := adminPerpus.PerpustakaanID
	
	var req models.Perpustakaan
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Update only allowed fields
	updateData := map[string]interface{}{
		"nama_perpustakaan": req.NamaPerpustakaan,
		"alamat":            req.Alamat,
		"jenis_perpustakaan": req.JenisPerpustakaan,
		"nomor_induk":       req.NomorInduk,
		"JumlahSDM":       req.JumlahSDM,
		"jumlahPengunjung": req.JumlahPengunjung,
		"JumlahAnggota":    req.JumlahAnggota,
	}
	
	if err := config.DB.Model(&models.Perpustakaan{}).Where("id = ?", perpustakaanID).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data perpustakaan"})
		return
	}
	
	// Create audit log
	auditLog := models.AuditLog{
		UserType:  "admin_perpustakaan",
		UserID:    userID,
		Action:    "UPDATE_PERUSTAKAAN",
		TableName: "perpustakaan",
		RecordID:  perpustakaanID,
		NewValues: stringifyMap(updateData),
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
		Timestamp: time.Now(),
	}
	config.DB.Create(&auditLog)
	
	c.JSON(http.StatusOK, gin.H{"message": "Data perpustakaan berhasil diperbarui"})
}

func SendDataToDPK(c *gin.Context) {
    userID := c.GetUint("user_id")
    
    var adminPerpus models.AdminPerpustakaan
    if err := config.DB.Preload("Perpustakaan").First(&adminPerpus, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": "Admin perpustakaan tidak ditemukan",
        })
        return
    }
    
    perpustakaanID := adminPerpus.PerpustakaanID
    
    var req models.SendDataRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }
    
    if req.PerpustakaanID != perpustakaanID {
        c.JSON(http.StatusForbidden, gin.H{
            "error": "Anda tidak memiliki akses ke perpustakaan ini",
        })
        return
    }
    
    // Check if data is complete
    var perpustakaan models.Perpustakaan
    if err := config.DB.Preload("AdminPerpustakaan").
        Preload("Koleksi").
        Preload("SDM").
        Preload("Pengunjung").
        Preload("Anggota").
        First(&perpustakaan, perpustakaanID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": "Perpustakaan tidak ditemukan",
        })
        return
    }
    
    if perpustakaan.NamaPerpustakaan == "" || perpustakaan.Alamat == "" || 
       perpustakaan.JenisPerpustakaan == "" || perpustakaan.NomorInduk == 0 {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Data perpustakaan belum lengkap",
            "required_fields": []string{
                "nama_perpustakaan", 
                "alamat", 
                "jenis_perpustakaan", 
                "nomor_induk",
            },
        })
        return
    }
    
    // Start transaction
    tx := config.DB.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    
    // Update perpustakaan status
    now := time.Now()
    updates := map[string]interface{}{
        "status_verifikasi": "Terkirim",
        "tanggal_kirim":     now,
        "catatan_revisi":    req.CatatanKirim,
    }
    
    if err := tx.Model(&models.Perpustakaan{}).
        Where("id = ?", perpustakaanID).
        Updates(updates).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Gagal mengupdate status perpustakaan",
        })
        return
    }
    
    // Create verification record
    verifikasi := models.Verifikasi{
        PerpustakaanID:    perpustakaanID,
        JenisData:        "Perpustakaan",
        Status:           "Pending",
        CatatanRevisi:    req.CatatanKirim,
        TanggalVerifikasi: &now,
    }
    
    if err := tx.Create(&verifikasi).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Gagal membuat record verifikasi",
        })
        return
    }
    
    // Create notification
    notifikasi := models.Notifikasi{
        JudulNotifikasi: "Data Perpustakaan Baru",
        IsiNotifikasi:   "Perpustakaan " + perpustakaan.NamaPerpustakaan + " telah mengirim data untuk diverifikasi",
        JenisNotifikasi: "info",
        TujuanUser:      "admin_dpk",
        RelatedID:      &perpustakaanID,
        RelatedType:    "perpustakaan",
        TanggalKirim:    now,
        ExpiredAt:       now.Add(7 * 24 * time.Hour),
        ActionLink:     "/admin-dpk/verifikasi/" + strconv.FormatUint(uint64(perpustakaanID), 10),
    }
    
    if err := tx.Create(&notifikasi).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Gagal membuat notifikasi",
        })
        return
    }
    
    // Create audit log
    auditLog := models.AuditLog{
        UserType:  "admin_perpustakaan",
        UserID:    userID,
        Action:    "SEND_DATA_TO_DPK",
        TableName: "perpustakaan",
        RecordID:  perpustakaanID,
        NewValues: stringifyMap(updates),
        IPAddress: c.ClientIP(),
        UserAgent: c.GetHeader("User-Agent"),
        Timestamp: now,
    }
    
    if err := tx.Create(&auditLog).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Gagal membuat audit log",
        })
        return
    }
    
    // Commit transaction
    if err := tx.Commit().Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Gagal menyelesaikan transaksi",
        })
        return
    }
    
    // Get updated perpustakaan data with relations
    var updatedPerpustakaan models.Perpustakaan
    if err := config.DB.Preload("AdminPerpustakaan").
        Preload("Koleksi").
        Preload("SDM").
        Preload("Pengunjung").
        Preload("Anggota").
        First(&updatedPerpustakaan, perpustakaanID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Gagal mengambil data perpustakaan terbaru",
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "data": updatedPerpustakaan,
        "message": "Data berhasil dikirim ke DPK untuk verifikasi",
    })
}
func GetHistoryPengiriman(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var adminPerpus models.AdminPerpustakaan
	if err := config.DB.Preload("Perpustakaan").First(&adminPerpus, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin perpustakaan tidak ditemukan"})
		return
	}
	
	perpustakaanID := adminPerpus.PerpustakaanID
	
	var verifikasi []models.Verifikasi
	if err := config.DB.Where("perpustakaan_id = ?", perpustakaanID).
		Order("tanggal_verifikasi desc").
		Find(&verifikasi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil riwayat verifikasi"})
		return
	}
	
	c.JSON(http.StatusOK, verifikasi)
}
func InputDataPerpustakaan(c *gin.Context) {
    userID := c.GetUint("user_id")
    // Get admin data with preloaded perpustakaan
    var adminPerpus models.AdminPerpustakaan
    if err := config.DB.Preload("Perpustakaan").First(&adminPerpus, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Admin perpustakaan tidak ditemukan"})
        return
    }

    var req models.InputDataPerpustakaanRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Start transaction
    tx := config.DB.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan internal"})
        }
    }()

    // Update perpustakaan data sesuai models.Perpustakaan
    perpustakaan := adminPerpus.Perpustakaan
    perpustakaan.Periode = req.Periode
    perpustakaan.NamaPerpustakaan = req.NamaPerpustakaan
    perpustakaan.Alamat = req.Alamat
    perpustakaan.KepalaPerpustakaan = req.KepalaPerpustakaan
    perpustakaan.JenisPerpustakaan = req.JenisPerpustakaan
    perpustakaan.TahunBerdiri = req.TahunBerdiri
    perpustakaan.NomorInduk = req.NomorInduk
    perpustakaan.JumlahSDM = req.JumlahSDM
    perpustakaan.JumlahPengunjung = req.JumlahPengunjung
    perpustakaan.JumlahAnggota = req.JumlahAnggota
    perpustakaan.StatusVerifikasi = "Draft"
    perpustakaan.TanggalKirim = nil
    perpustakaan.CatatanRevisi = ""

    if err := tx.Save(&perpustakaan).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data perpustakaan"})
        return
    }

    // Create audit log with admin info
    auditLog := models.AuditLog{
        UserType:  "admin_perpustakaan",
        UserID:    userID,
        Action:    "UPDATE_PERUSTAKAAN_DATA",
        TableName: "perpustakaan",
        RecordID:  perpustakaan.ID,
        NewValues: stringifyMap(map[string]interface{}{
            "periode":            req.Periode,
            "nama_perpustakaan":  req.NamaPerpustakaan,
            "alamat":             req.Alamat,
            "kepala_perpustakaan": req.KepalaPerpustakaan,
            "jenis_perpustakaan": req.JenisPerpustakaan,
            "tahun_berdiri":      req.TahunBerdiri,
            "nomor_induk":        req.NomorInduk,
            "jumlah_sdm":         req.JumlahSDM,
            "jumlah_pengunjung":  req.JumlahPengunjung,
            "jumlah_anggota":     req.JumlahAnggota,
            "tanggal_kirim":      time.Now(), // Tanggal kirim masih null karena belum dikirim
            "updated_by":         adminPerpus.NamaLengkap + " (" + adminPerpus.Username + ")",
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

    // Preload admin data for response
    if err := config.DB.Preload("AdminPerpustakaan").First(&perpustakaan, perpustakaan.ID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat data perpustakaan"})
        return
    }

    // Format response
    response := gin.H{
        "data": gin.H{
            "id":                 perpustakaan.ID,
            "created_at":         perpustakaan.CreatedAt,
            "updated_at":         perpustakaan.UpdatedAt,
            "periode":            perpustakaan.Periode,
            "nama_perpustakaan":  perpustakaan.NamaPerpustakaan,
            "alamat":             perpustakaan.Alamat,
            "kepala_perpustakaan": perpustakaan.KepalaPerpustakaan,
            "jenis_perpustakaan": perpustakaan.JenisPerpustakaan,
            "tahun_berdiri":      perpustakaan.TahunBerdiri,
            "nomor_induk":        perpustakaan.NomorInduk,
            "jumlah_sdm":         perpustakaan.JumlahSDM,
            "jumlah_pengunjung":  perpustakaan.JumlahPengunjung,
            "jumlah_anggota":     perpustakaan.JumlahAnggota,
            "status_verifikasi":  perpustakaan.StatusVerifikasi,
            "tanggal_kirim":      perpustakaan.TanggalKirim,
            "catatan_revisi":     perpustakaan.CatatanRevisi,
            "admin_perpustakaan": gin.H{
                "id":           adminPerpus.ID,
                "nama_lengkap": adminPerpus.NamaLengkap,
                "username":     adminPerpus.Username,
            },
            "koleksi":    perpustakaan.Koleksi,
            "sdm":        perpustakaan.SDM,
            "pengunjung": perpustakaan.Pengunjung,
            "anggota":    perpustakaan.Anggota,
        },
        "message": "Data perpustakaan berhasil diperbarui",
    }

    c.JSON(http.StatusOK, response)
}
