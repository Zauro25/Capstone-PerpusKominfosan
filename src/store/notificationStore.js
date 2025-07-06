import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useNotificationStore = defineStore('notification', () => {
  const notifications = ref([])
  const hasUnreadNotifications = ref(false)
  const lastReadTimestamp = ref(Date.now())

  // Fungsi untuk menambah notifikasi baru
  const addNotification = (notification) => {
    const newNotification = {
      id: Date.now(),
      timestamp: Date.now(),
      read: false,
      ...notification
    }
    notifications.value.unshift(newNotification) // Menambah di awal array
    updateUnreadStatus()
  }

  // Fungsi untuk menandai notifikasi sebagai sudah dibaca
  const markAsRead = (notificationId) => {
    const notification = notifications.value.find(n => n.id === notificationId)
    if (notification) {
      notification.read = true
      updateUnreadStatus()
    }
  }

  // Fungsi untuk menandai semua notifikasi sebagai sudah dibaca
  const markAllAsRead = () => {
    notifications.value.forEach(notification => {
      notification.read = true
    })
    lastReadTimestamp.value = Date.now()
    updateUnreadStatus()
  }

  // Fungsi untuk mengupdate status unread
  const updateUnreadStatus = () => {
    hasUnreadNotifications.value = notifications.value.some(notification => !notification.read)
  }

  // Fungsi untuk mendapatkan jumlah notifikasi yang belum dibaca
  const getUnreadCount = () => {
    return notifications.value.filter(notification => !notification.read).length
  }

  return {
    notifications,
    hasUnreadNotifications,
    lastReadTimestamp,
    addNotification,
    markAsRead,
    markAllAsRead,
    getUnreadCount
  }
}) 