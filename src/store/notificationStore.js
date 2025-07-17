import { defineStore } from 'pinia'

export const useNotificationStore = defineStore('notification', {
  state: () => ({
    hasUnreadNotifications: true
  }),

  actions: {
    markNotificationsAsRead() {
      this.hasUnreadNotifications = false
    },

    setUnreadNotifications() {
      this.hasUnreadNotifications = true
    }
  }
}) 