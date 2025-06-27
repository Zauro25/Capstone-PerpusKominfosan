package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/Zauro25/Capstone-PerpusKominfosan/controllers"
	"github.com/Zauro25/Capstone-PerpusKominfosan/middleware"
)

func SetupRoutes(r *gin.Engine) {
	
	api := r.Group("")

	// Public routes
	api.POST("/login", controllers.Login)
	api.POST("/register", controllers.Register)
	api.POST("/register-admin-dpk", controllers.RegisterAdminDPK)
	api.POST("/register-executive", controllers.RegisterExecutive)
	
	// Protected routes
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	protected.Use(middleware.AuditLogMiddleware())

	// Auth routes
	protected.POST("/logout", controllers.Logout)
	protected.POST("/change-password", controllers.ChangePassword)
	protected.GET("/profile", controllers.GetProfile)
	protected.PUT("/profile", controllers.UpdateProfile)

	// Admin Perpustakaan routes
	adminPerpus := protected.Group("/admin-perpustakaan")
	adminPerpus.Use(middleware.RequireRole("admin_perpustakaan"))
	{
		adminPerpus.GET("/dashboard", controllers.GetDashboardPerpustakaan)
		adminPerpus.GET("/data", controllers.GetDataPerpustakaan)
		adminPerpus.PUT("/data", controllers.UpdateDataPerpustakaan)
		adminPerpus.POST("/input-data", controllers.InputDataPerpustakaan)
		adminPerpus.POST("/send-data", controllers.SendDataToDPK)
		adminPerpus.GET("/history", controllers.GetHistoryPengiriman)
		adminPerpus.GET("/notifications", controllers.GetNotifications)
		adminPerpus.PUT("/notifications/:id/read", controllers.MarkNotificationAsRead)
	}

	// Admin DPK routes
	adminDPK := protected.Group("/admin-dpk")
	adminDPK.Use(middleware.RequireRole("admin_dpk"))
	{
		adminDPK.GET("/dashboard", controllers.GetDashboardDPK)
		adminDPK.GET("/perpustakaan", controllers.GetAllPerpustakaan)
		adminDPK.GET("/perpustakaan/:id", controllers.GetDetailPerpustakaan)
		adminDPK.POST("/verifikasi", controllers.VerifikasiData)
		adminDPK.GET("/verifikasiakun", controllers.GetPendingVerification)
		adminDPK.POST("/verifikasiakun", controllers.VerifyAdminPerpustakaan)
		adminDPK.POST("/laporan", controllers.GenerateLaporan)
		adminDPK.GET("/laporan/:id/download", controllers.DownloadLaporan)
		adminDPK.POST("/notifications/broadcast", controllers.SendBroadcastNotification)
		adminDPK.GET("/audit-logs", controllers.GetAuditLogs)
		adminDPK.POST("/send-reminder", controllers.SendReminder)
		
		// Admin perpustakaan management
		adminDPK.GET("/pending-admin-verifications", controllers.GetPendingAdminVerifications)
		adminDPK.POST("/verify-admin-perpustakaan", controllers.VerifyAdminPerpustakaan)
		adminDPK.GET("/manage-accounts", controllers.GetAllAdminPerpustakaan)
		adminDPK.POST("/manage-accounts", controllers.CreateAdminPerpustakaan)
		adminDPK.PUT("/manage-accounts/:id", controllers.UpdateAdminPerpustakaan)
		adminDPK.DELETE("/manage-accounts/:id", controllers.DeleteAdminPerpustakaan)
		adminDPK.PUT("/manage-accounts/:id/activate", controllers.ActivateAdminPerpustakaan)
		adminDPK.PUT("/manage-accounts/:id/deactivate", controllers.DeactivateAdminPerpustakaan)
		adminDPK.POST("/manage-accounts/:id/reset-password", controllers.ResetPasswordAdminPerpustakaan)
	}

	// Executive routes
	executive := protected.Group("/executive")
	executive.Use(middleware.RequireRole("executive"))
	{
		executive.GET("/dashboard", controllers.GetDashboardExecutive)
		executive.GET("/statistics", controllers.GetStatistiksPerpustakaan)
		executive.GET("/laporan", controllers.GetLaporanExecutive)
	}

	// Common routes (accessible by all authenticated users)
	protected.GET("/notifications", controllers.GetNotifications)
	protected.PUT("/notifications/:id/read", controllers.MarkNotificationAsRead)
}