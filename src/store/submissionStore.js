import { defineStore } from 'pinia'

export const useSubmissionStore = defineStore('submission', {
  state: () => ({
    submittedData: []
  }),

  actions: {
    submitData(data) {
      // Add submission date and set status to "Sudah Dikirim"
      const submissionData = {
        ...data,
        submissionDate: new Date().toISOString(),
        status: 'Sudah Dikirim' // Always set to "Sudah Dikirim" when submitting
      }
      
      // Check if data already exists for this period
      const existingIndex = this.submittedData.findIndex(item => item.periode === data.periode)
      
      if (existingIndex !== -1) {
        // Update existing submission
        this.submittedData[existingIndex] = {
          ...this.submittedData[existingIndex],
          ...submissionData
        }
      } else {
        // Add new submission
        this.submittedData.push(submissionData)
      }
      
      // Save to localStorage
      this.saveToLocalStorage()
    },

    updateStatus(periode, newStatus) {
      // Only allow "Belum Dikirim" or "Sudah Dikirim" as valid statuses
      if (newStatus !== 'Belum Dikirim' && newStatus !== 'Sudah Dikirim') {
        return
      }

      const data = this.submittedData.find(item => item.periode === periode)
      if (data) {
        data.status = newStatus
        this.saveToLocalStorage()
      } else {
        // If no submission exists, create one with the new status
        this.submitData({
          periode,
          status: newStatus
        })
      }
    },

    loadFromLocalStorage() {
      const saved = localStorage.getItem('submittedData')
      if (saved) {
        this.submittedData = JSON.parse(saved)
      }
    },

    saveToLocalStorage() {
      localStorage.setItem('submittedData', JSON.stringify(this.submittedData))
    }
  }
}) 