import { defineStore } from 'pinia'

export const useProfileStore = defineStore('profile', {
  state: () => ({
    userProfile: {
      nama_lengkap: '',
      email: '',
      no_telepon: '',
      username: '',
      user_type: ''
    }
  }),

  actions: {
    async fetchUserProfile() {
      try {
        const token = localStorage.getItem('authToken') || sessionStorage.getItem('authToken')
        if (!token) {
          throw new Error('No auth token found')
        }

        // Temporary: Set default profile data until API is ready
        this.userProfile = {
          nama_lengkap: 'Admin Perpustakaan',
          email: 'admin@perpustakaan.com',
          no_telepon: '08123456789',
          username: 'admin',
          user_type: 'admin_perpustakaan'
        }

        // Comment out API call until endpoint is ready
        /*
        const response = await fetch('http://localhost:3000/api/profile', {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        })

        if (!response.ok) {
          throw new Error('Failed to fetch profile')
        }

        const data = await response.json()
        this.userProfile = data
        */
      } catch (error) {
        console.error('Error fetching profile:', error)
        // Don't throw error, just use default data
        this.userProfile = {
          nama_lengkap: 'Admin Perpustakaan',
          email: 'admin@perpustakaan.com',
          no_telepon: '08123456789',
          username: 'admin',
          user_type: 'admin_perpustakaan'
        }
      }
    },

    async updateUserProfile(profileData) {
      try {
        const token = localStorage.getItem('authToken') || sessionStorage.getItem('authToken')
        if (!token) {
          throw new Error('No auth token found')
        }

        // Temporary: Just update local state until API is ready
        this.userProfile = { ...this.userProfile, ...profileData }

        // Comment out API call until endpoint is ready
        /*
        const response = await fetch('http://localhost:3000/api/profile', {
          method: 'PUT',
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(profileData)
        })

        if (!response.ok) {
          throw new Error('Failed to update profile')
        }

        const data = await response.json()
        this.userProfile = { ...this.userProfile, ...profileData }
        return data
        */
      } catch (error) {
        console.error('Error updating profile:', error)
        throw error
      }
    }
  }
}) 