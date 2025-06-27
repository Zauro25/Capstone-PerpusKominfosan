package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Zauro25/Capstone-PerpusKominfosan/config"
	"github.com/Zauro25/Capstone-PerpusKominfosan/models"
)

func GetDashboardExecutive(c *gin.Context) {
	// Get statistics
	var stats struct {
		TotalPerpustakaan int64
		TotalKoleksi      int64
		TotalPengunjung   int64
		TotalAnggota      int64
		TotalSDM          int64
	}
	
	config.DB.Model(&models.Perpustakaan{}).Count(&stats.TotalPerpustakaan)
	config.DB.Model(&models.Koleksi{}).Count(&stats.TotalKoleksi)
	config.DB.Model(&models.Pengunjung{}).Count(&stats.TotalPengunjung)
	config.DB.Model(&models.Anggota{}).Where("status_anggota = ?", "Aktif").Count(&stats.TotalAnggota)
	config.DB.Model(&models.SDM{}).Count(&stats.TotalSDM)
	
	// Get recent reports
	var laporan []models.Laporan
	config.DB.Order("tanggal_generate desc").Limit(5).Find(&laporan)
	
	// Get recent activities
	var activities []models.AuditLog
	config.DB.Where("user_type = ?", "admin_dpk").
		Order("timestamp desc").
		Limit(5).
		Find(&activities)
	
	c.JSON(http.StatusOK, gin.H{
		"statistics": stats,
		"reports":    laporan,
		"activities": activities,
	})
}

func GetStatistiksPerpustakaan(c *gin.Context) {
	// Get statistics by perpustakaan type
	var stats []struct {
		JenisPerpustakaan string
		Total             int64
		TotalKoleksi      int64
		TotalPengunjung   int64
	}
	
	config.DB.Model(&models.Perpustakaan{}).
		Select("jenis_perpustakaan, COUNT(*) as total").
		Group("jenis_perpustakaan").
		Scan(&stats)
	
	// Add koleksi and pengunjung counts
	for i := range stats {
		config.DB.Model(&models.Koleksi{}).
			Joins("JOIN perpustakaan ON koleksi.perpustakaan_id = perpustakaan.id").
			Where("perpustakaan.jenis_perpustakaan = ?", stats[i].JenisPerpustakaan).
			Count(&stats[i].TotalKoleksi)
			
		config.DB.Model(&models.Pengunjung{}).
			Joins("JOIN perpustakaan ON pengunjung.perpustakaan_id = perpustakaan.id").
			Where("perpustakaan.jenis_perpustakaan = ?", stats[i].JenisPerpustakaan).
			Count(&stats[i].TotalPengunjung)
	}
	
	c.JSON(http.StatusOK, stats)
}

func GetLaporanExecutive(c *gin.Context) {
	var laporan []models.Laporan
	
	query := config.DB
	
	// Filter by periode if provided
	if periode := c.Query("periode"); periode != "" {
		query = query.Where("periode = ?", periode)
	}
	
	// Filter by jenis if provided
	if jenis := c.Query("jenis"); jenis != "" {
		query = query.Where("jenis_laporan = ?", jenis)
	}
	
	if err := query.Order("tanggal_generate desc").Find(&laporan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data laporan"})
		return
	}
	
	c.JSON(http.StatusOK, laporan)
}