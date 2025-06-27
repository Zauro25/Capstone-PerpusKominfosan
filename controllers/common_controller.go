package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/Zauro25/Capstone-PerpusKominfosan/config"
	"github.com/Zauro25/Capstone-PerpusKominfosan/models"
	"github.com/Zauro25/Capstone-PerpusKominfosan/services"
)

var notificationService = services.NewNotificationService(config.DB)

func GetNotifications(c *gin.Context) {
	userID := c.GetUint("user_id")
	userType := c.GetString("user_type")
	
	// Get query parameters
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}
	
	unreadOnly := c.Query("unread_only") == "true"
	
	notifications, err := notificationService.GetUserNotifications(userID, userType, limit, unreadOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil notifikasi"})
		return
	}
	
	c.JSON(http.StatusOK, notifications)
}

func MarkNotificationAsRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	notificationID := c.Param("id")
	
	id, err := strconv.ParseUint(notificationID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID notifikasi tidak valid"})
		return
	}
	
	if err := notificationService.MarkAsRead(uint(id), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menandai notifikasi"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Notifikasi berhasil ditandai sebagai dibaca"})
}

func MarkAllNotificationsAsRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	userType := c.GetString("user_type")
	
	if err := notificationService.MarkAllAsRead(userID, userType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menandai notifikasi"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Semua notifikasi berhasil ditandai sebagai dibaca"})
}

func GetNotificationCount(c *gin.Context) {
	userID := c.GetUint("user_id")
	userType := c.GetString("user_type")
	
	var count int64
	err := config.DB.Model(&models.Notifikasi{}).
		Where("(tujuan_user = ? OR tujuan_user = 'all') AND (user_id = ? OR user_id IS NULL) AND is_read = ?", 
			userType, userID, false).
		Count(&count).Error
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung notifikasi"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"count": count})
}