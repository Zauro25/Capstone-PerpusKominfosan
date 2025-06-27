package tasks

import (
	"log"
	"github.com/Zauro25/Capstone-PerpusKominfosan/config"
	"github.com/Zauro25/Capstone-PerpusKominfosan/models"
	"github.com/Zauro25/Capstone-PerpusKominfosan/services"
	"time"
)

func StartNotificationTasks() {
	ticker := time.NewTicker(24 * time.Hour) // Jalankan setiap 24 jam
	go func() {
		for range ticker.C {
			cleanupExpiredNotifications()
			sendPeriodicReminders()
		}
	}()
}

func cleanupExpiredNotifications() {
	notificationService := services.NewNotificationService(config.DB)
	if err := notificationService.DeleteExpiredNotifications(); err != nil {
		log.Printf("Gagal membersihkan notifikasi kadaluarsa: %v", err)
	} else {
		log.Println("Berhasil membersihkan notifikasi kadaluarsa")
	}
}

func sendPeriodicReminders() {
	// Dapatkan daftar perpustakaan yang belum mengirim data
	var perpustakaan []models.Perpustakaan
	if err := config.DB.Where("status_verifikasi != ?", "Terkirim").Find(&perpustakaan).Error; err != nil {
		log.Printf("Gagal mendapatkan daftar perpustakaan: %v", err)
		return
	}

	notificationService := services.NewNotificationService(config.DB)
	currentSemester := getCurrentSemester()

	for _, p := range perpustakaan {
		message := "Pengingat: Silakan lengkapi dan kirim data perpustakaan untuk semester " + currentSemester
		if err := notificationService.SendReminderNotification(p.ID, message); err != nil {
			log.Printf("Gagal mengirim pengingat ke perpustakaan %d: %v", p.ID, err)
		}
	}
}

func getCurrentSemester() string {
	now := time.Now()
	year := now.Year()
	semester := "1"
	if now.Month() >= 7 {
		semester = "2"
	}
	return string(year) + "-" + semester
}