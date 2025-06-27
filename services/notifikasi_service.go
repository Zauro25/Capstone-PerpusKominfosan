package services

import (
	"gorm.io/gorm"
	"github.com/Zauro25/Capstone-PerpusKominfosan/models"
	"time"
)

type NotificationService struct {
	db *gorm.DB
}

func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{db: db}
}

// CreateNotification membuat notifikasi baru
func (s *NotificationService) CreateNotification(notification *models.Notifikasi) error {
	if notification.TanggalKirim.IsZero() {
		notification.TanggalKirim = time.Now()
	}
	return s.db.Create(notification).Error
}

// GetUserNotifications mengambil notifikasi untuk user tertentu
func (s *NotificationService) GetUserNotifications(userID uint, userType string, limit int, unreadOnly bool) ([]models.Notifikasi, error) {
	var notifications []models.Notifikasi

	query := s.db.Where("(tujuan_user = ? OR tujuan_user = 'all') AND (user_id = ? OR user_id IS NULL)", 
		userType, userID)

	if unreadOnly {
		query = query.Where("is_read = ?", false)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Order("tanggal_kirim DESC").Find(&notifications).Error
	return notifications, err
}

// MarkAsRead menandai notifikasi sebagai telah dibaca
func (s *NotificationService) MarkAsRead(notificationID uint, userID uint) error {
	return s.db.Model(&models.Notifikasi{}).
		Where("id = ? AND (user_id = ? OR user_id IS NULL)", notificationID, userID).
		Update("is_read", true).Error
}

// MarkAllAsRead menandai semua notifikasi user sebagai telah dibaca
func (s *NotificationService) MarkAllAsRead(userID uint, userType string) error {
	return s.db.Model(&models.Notifikasi{}).
		Where("(tujuan_user = ? OR tujuan_user = 'all') AND (user_id = ? OR user_id IS NULL) AND is_read = ?", 
			userType, userID, false).
		Update("is_read", true).Error
}

// DeleteExpiredNotifications menghapus notifikasi yang sudah kadaluarsa
func (s *NotificationService) DeleteExpiredNotifications() error {
	return s.db.Where("expired_at < ?", time.Now()).Delete(&models.Notifikasi{}).Error
}

// SendVerificationNotification mengirim notifikasi verifikasi
func (s *NotificationService) SendVerificationNotification(perpustakaanID uint, status string, catatan string) error {
	var perpustakaan models.Perpustakaan
	if err := s.db.First(&perpustakaan, perpustakaanID).Error; err != nil {
		return err
	}

	judul := "Verifikasi Data Perpustakaan"
	isi := "Data perpustakaan " + perpustakaan.NamaPerpustakaan + " telah diverifikasi dengan status: " + status
	jenis := "info"
	if status == "Direvisi" {
		jenis = "warning"
		isi += "\nCatatan: " + catatan
	}

	notification := models.Notifikasi{
		JudulNotifikasi: judul,
		IsiNotifikasi:   isi,
		JenisNotifikasi: jenis,
		TujuanUser:      "admin_perpustakaan",
		UserID:         &perpustakaan.ID,
		RelatedID:      &perpustakaanID,
		RelatedType:    "perpustakaan",
		TanggalKirim:    time.Now(),
		ExpiredAt:       time.Now().Add(30 * 24 * time.Hour), // Expire dalam 30 hari
		ActionLink:      "/admin-perpustakaan/data?status=revision",
	}

	return s.CreateNotification(&notification)
}

// SendDataSubmittedNotification mengirim notifikasi data terkirim
func (s *NotificationService) SendDataSubmittedNotification(perpustakaanID uint) error {
	var perpustakaan models.Perpustakaan
	if err := s.db.First(&perpustakaan, perpustakaanID).Error; err != nil {
		return err
	}

	notification := models.Notifikasi{
		JudulNotifikasi: "Data Terkirim",
		IsiNotifikasi:   "Data perpustakaan " + perpustakaan.NamaPerpustakaan + " telah berhasil dikirim ke DPK untuk verifikasi",
		JenisNotifikasi: "success",
		TujuanUser:      "admin_perpustakaan",
		UserID:         &perpustakaan.ID,
		RelatedID:      &perpustakaanID,
		RelatedType:    "perpustakaan",
		TanggalKirim:    time.Now(),
		ExpiredAt:       time.Now().Add(7 * 24 * time.Hour),
	}

	return s.CreateNotification(&notification)
}

// SendReminderNotification mengirim notifikasi pengingat
func (s *NotificationService) SendReminderNotification(perpustakaanID uint, message string) error {
	var perpustakaan models.Perpustakaan
	if err := s.db.First(&perpustakaan, perpustakaanID).Error; err != nil {
		return err
	}

	notification := models.Notifikasi{
		JudulNotifikasi: "Pengingat",
		IsiNotifikasi:   message,
		JenisNotifikasi: "warning",
		TujuanUser:      "admin_perpustakaan",
		UserID:         &perpustakaan.ID,
		RelatedID:      &perpustakaanID,
		RelatedType:    "perpustakaan",
		TanggalKirim:    time.Now(),
		ExpiredAt:       time.Now().Add(7 * 24 * time.Hour),
		ActionLink:      "/admin-perpustakaan/data",
	}

	return s.CreateNotification(&notification)
}