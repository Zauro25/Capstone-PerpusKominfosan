import { defineStore } from 'pinia'

export const useSubmissionStore = defineStore('submission', {
  state: () => ({
    submittedData: []
  }),

  actions: {
    submitData(data) {
      // Add submission date and initial status
      const submissionData = {
        ...data,
        submissionDate: new Date().toISOString(),
        status: 'Menunggu'
      }
      this.submittedData.push(submissionData)
      // Save to localStorage
      this.saveToLocalStorage()
    },

    updateStatus(periode, newStatus) {
      const data = this.submittedData.find(item => item.periode === periode)
      if (data) {
        data.status = newStatus
        this.saveToLocalStorage()
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