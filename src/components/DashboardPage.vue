<template>
  <div class="dashboard-container">
    <!-- Header -->
    <header class="header">
      <button 
        class="hamburger-menu"
        @click="toggleSidebar"
        :class="{ 'active': isSidebarOpen }"
      >
        <span></span>
        <span></span>
        <span></span>
      </button>
      <div class="header-left">
        <img src="../assets/logo-sidapus.png" alt="Logo" class="logo" />
        <h1>Sistem Data Perpustakaan<br>Dan Kearsipan</h1>
      </div>
      <div class="header-right">
        <div class="notification-btn" @click="navigateToNotifications">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"></path>
            <path d="M13.73 21a2 2 0 0 1-3.46 0"></path>
          </svg>
          <span v-if="hasUnreadNotifications" class="notification-dot"></span>
        </div>
        <div class="profile-btn" @click="goToSettings">
          <span>Admin Perpustakaan</span>
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
            <circle cx="12" cy="7" r="4"></circle>
          </svg>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <div class="main-content">
      <!-- Sidebar -->
      <aside class="sidebar" :class="{ 'active': isSidebarOpen }">
        <nav class="sidebar-menu">
          <button class="nav-btn active">
            <span>Dashboard</span>
          </button>
          <button class="nav-btn" @click="navigateTo('input-update')">
            <span>Input & Update Data</span>
          </button>
          <button class="nav-btn" @click="navigateTo('pengiriman')">
            <span>Pengiriman Data</span>
          </button>
          <button class="nav-btn" @click="navigateTo('validasi')">
            <span>Validasi dan Revisi dari DPK</span>
          </button>
        </nav>
        <button class="sidebar-logout-btn" @click="logout">
          <span>Keluar</span>
        </button>
      </aside>

      <!-- Sidebar Overlay for Mobile -->
      <div class="sidebar-overlay" 
           :class="{ 'active': isSidebarOpen }" 
           @click="toggleSidebar">
      </div>

      <!-- Dashboard Content -->
      <main class="dashboard-content">
        <h2>Dashboard</h2>
        
        <!-- Semester Selector -->
        <div class="semester-selector">
          <select v-model="selectedSemester" @change="fetchDashboardData">
            <option value="2025-2">Semester Genap 2025/2026</option>
            <option value="2024-1">Semester Ganjil 2024/2025</option>
            <option value="2024-2">Semester Genap 2024/2025</option>
            <option value="2023-1">Semester Ganjil 2023/2024</option>
            <option value="2023-2">Semester Genap 2023/2024</option>
          </select>
        </div>

        <!-- Statistics Chart -->
        <div class="stats-chart">
          <h3>Statistik Pengunjung</h3>
          <canvas ref="visitorChart"></canvas>
        </div>

        <!-- Distribution Charts -->
        <div class="charts-grid">
          <div class="chart-card">
            <h3>Distribusi Jenis Perpustakaan</h3>
            <canvas ref="libraryTypeChart"></canvas>
          </div>
          <div class="chart-card">
            <h3>Status Verifikasi Data</h3>
            <canvas ref="verificationChart"></canvas>
          </div>
          <div class="chart-card">
            <h3>Tren Pengunjung dan Anggota (6 Bulan Terakhir)</h3>
            <canvas ref="trendChart"></canvas>
          </div>
        </div>

        <!-- Summary Cards -->
        <div class="summary-cards">
          <div class="summary-card">
            <div class="icon-wrapper">
              <img src="../assets/total perpustakaan.png" alt="Library Icon" class="summary-icon">
            </div>
            <div class="card-content">
              <h4>Total Perpustakaan</h4>
              <p>{{ totalLibraries }}</p>
            </div>
          </div>
          <div class="summary-card">
            <div class="icon-wrapper">
              <img src="../assets/total sdm.png" alt="Staff Icon" class="summary-icon">
            </div>
            <div class="card-content">
              <h4>Total SDM</h4>
              <p>{{ totalStaff }}</p>
            </div>
          </div>
        </div>
      </main>
    </div>
  </div>
</template>

<script>
import Chart from 'chart.js/auto'

export default {
  data() {
    return {
      selectedSemester: '2025-2',
      showNotifications: false,
      hasUnreadNotifications: false,
      notifications: [],
      totalLibraries: 0,
      totalStaff: 0,
      isSidebarOpen: false,
      isMobile: false,
      charts: {
        visitorChart: null,
        libraryTypeChart: null,
        verificationChart: null,
        trendChart: null
      }
    }
  },
  async created() {
    await this.fetchDashboardData()
    await this.fetchNotifications()
    this.checkMobile()
    window.addEventListener('resize', this.checkMobile)
    document.addEventListener('click', this.handleClickOutside)
  },
  mounted() {
    this.initializeCharts()
  },
  beforeDestroy() {
    window.removeEventListener('resize', this.checkMobile)
    document.removeEventListener('click', this.handleClickOutside)
  },
  methods: {
    checkMobile() {
      this.isMobile = window.innerWidth <= 768
    },
    toggleSidebar() {
      this.isSidebarOpen = !this.isSidebarOpen
      // Prevent body scroll when sidebar is open on mobile
      if (window.innerWidth <= 768) {
        document.body.style.overflow = this.isSidebarOpen ? 'hidden' : ''
      }
    },
    handleClickOutside(event) {
      // Close sidebar when clicking outside on mobile
      if (this.isSidebarOpen && window.innerWidth <= 768) {
        const sidebar = document.querySelector('.sidebar')
        const menuToggle = document.querySelector('.hamburger-menu')
        if (!sidebar?.contains(event.target) && !menuToggle?.contains(event.target)) {
          this.toggleSidebar()
        }
      }
    },
    async fetchDashboardData() {
      try {
        const response = await fetch(`/api/dashboard-data?semester=${this.selectedSemester}`)
        const data = await response.json()
        
        this.totalLibraries = data.totalLibraries || 0
        this.totalStaff = data.totalStaff || 0
        
        this.updateCharts(data)
      } catch (error) {
        console.error('Error fetching dashboard data:', error)
      }
    },
    async fetchNotifications() {
      try {
        const response = await fetch('/api/notifications')
        const data = await response.json()
        this.notifications = data.notifications || [
          { id: 1, message: 'Data perpustakaan baru telah ditambahkan', time: '5 menit yang lalu', read: false },
          { id: 2, message: 'Validasi data menunggu persetujuan', time: '1 jam yang lalu', read: false }
        ]
        this.hasUnreadNotifications = this.notifications.some(n => !n.read)
      } catch (error) {
        console.error('Error fetching notifications:', error)
      }
    },
    initializeCharts() {
      const commonDoughnutOptions = {
        responsive: true,
        maintainAspectRatio: false,
        cutout: '60%',
        plugins: {
          legend: {
            position: 'bottom',
            labels: {
              boxWidth: 12,
              padding: 15,
              font: {
                size: 11
              }
            }
          }
        },
        layout: {
          padding: {
            top: 10,
            bottom: 10,
            left: 10,
            right: 10
          }
        }
      };

      // Initialize visitor chart
      this.charts.visitorChart = new Chart(this.$refs.visitorChart, {
        type: 'line',
        data: {
          labels: ['Jan', 'Feb', 'Mar', 'Apr', 'Mei', 'Jun', 'Jul', 'Agu', 'Sep', 'Okt', 'Nov', 'Des'],
          datasets: [
            {
              label: 'Jumlah Pengunjung',
              borderColor: '#4318FF',
              backgroundColor: '#4318FF',
              data: [],
              tension: 0.4
            },
            {
              label: 'Jumlah Anggota',
              borderColor: '#E31A1A',
              backgroundColor: '#E31A1A',
              data: [],
              tension: 0.4
            }
          ]
        },
        options: {
          responsive: true,
          maintainAspectRatio: false,
          plugins: {
            legend: {
              position: 'top',
              labels: {
                boxWidth: 12,
                padding: 15,
                font: {
                  size: 11
                }
              }
            }
          },
          scales: {
            y: {
              beginAtZero: true,
              max: 1000,
              ticks: {
                stepSize: 200,
                font: {
                  size: 10
                }
              }
            },
            x: {
              ticks: {
                font: {
                  size: 10
                }
              }
            }
          }
        }
      });

      // Initialize library type chart (donut)
      this.charts.libraryTypeChart = new Chart(this.$refs.libraryTypeChart, {
        type: 'doughnut',
        data: {
          labels: ['Perpustakaan Umum', 'Perpustakaan Sekolah', 'Perpustakaan Khusus'],
          datasets: [{
            data: [],
            backgroundColor: ['#0E2954', '#4318FF', '#E31A1A']
          }]
        },
        options: commonDoughnutOptions
      });

      // Initialize verification status chart (donut)
      this.charts.verificationChart = new Chart(this.$refs.verificationChart, {
        type: 'doughnut',
        data: {
          labels: ['Data Valid', 'Data Revisi', 'Data Menunggu'],
          datasets: [{
            data: [],
            backgroundColor: ['#0E2954', '#4318FF', '#E31A1A']
          }]
        },
        options: commonDoughnutOptions
      });

      // Initialize trend chart (donut)
      this.charts.trendChart = new Chart(this.$refs.trendChart, {
        type: 'doughnut',
        data: {
          labels: ['Data Pengunjung', 'Data Anggota Aktif'],
          datasets: [{
            data: [],
            backgroundColor: ['#0E2954', '#4318FF']
          }]
        },
        options: commonDoughnutOptions
      });
    },
    updateCharts(data) {
      if (this.charts.visitorChart && data.visitorData) {
        this.charts.visitorChart.data.datasets[0].data = data.visitorData
        this.charts.visitorChart.data.datasets[1].data = data.memberData
        this.charts.visitorChart.update()
      }
      if (this.charts.libraryTypeChart && data.libraryTypeData) {
        this.charts.libraryTypeChart.data.datasets[0].data = data.libraryTypeData
        this.charts.libraryTypeChart.update()
      }
      if (this.charts.verificationChart && data.verificationData) {
        this.charts.verificationChart.data.datasets[0].data = data.verificationData
        this.charts.verificationChart.update()
      }
      if (this.charts.trendChart && data.trendData) {
        this.charts.trendChart.data.datasets[0].data = data.trendData
        this.charts.trendChart.update()
      }
    },
    toggleNotifications() {
      this.showNotifications = !this.showNotifications
    },
    async readNotification(id) {
      try {
        await fetch(`/api/notifications/${id}/read`, { method: 'POST' })
        this.notifications = this.notifications.map(n => 
          n.id === id ? { ...n, read: true } : n
        )
        this.hasUnreadNotifications = this.notifications.some(n => !n.read)
      } catch (error) {
        console.error('Error marking notification as read:', error)
      }
    },
    async markAllAsRead() {
      try {
        await fetch('/api/notifications/mark-all-read', { method: 'POST' })
        this.notifications = this.notifications.map(n => ({ ...n, read: true }))
        this.hasUnreadNotifications = false
      } catch (error) {
        console.error('Error marking all notifications as read:', error)
      }
    },
    navigateTo(route) {
      this.$router.push(`/${route}`)
    },
    goToSettings() {
      this.$router.push('/profile')
    },
    logout() {
      localStorage.removeItem('authToken')
      sessionStorage.removeItem('authToken')
      this.$router.push('/login')
    },
    navigateToNotifications() {
      this.$router.push('/notifications');
    }
  }
}
</script>

<style scoped>
/* Reset default margins and padding */
html, body {
  margin: 0;
  padding: 0;
  height: 100%;
  width: 100%;
  overflow: hidden;
}

* {
  box-sizing: border-box;
  font-family: inter, sans-serif;
}

.dashboard-container {
  height: 100vh;
  width: 100%;
  background-color: #ffffff;
  display: flex;
  flex-direction: column;
  position: relative;
}

.header {
  background-color: #0E2954;
  color: white;
  padding: 0.75rem 1.5rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
  height: 70px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.logo {
  height: 35px;
  width: auto;
}

.header-left h1 {
  color: white;
  font-size: 1.1rem;
  line-height: 1.3;
  margin: 0;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  margin-left: auto;
}

.notification-btn {
  position: relative;
  cursor: pointer;
  display: flex;
  align-items: center;
  color: white;
  transition: opacity 0.2s ease;
}

.notification-btn:hover {
  opacity: 0.8;
}

.notification-dot {
  position: absolute;
  top: -2px;
  right: -2px;
  width: 8px;
  height: 8px;
  background-color: #FF4B4B;
  border-radius: 50%;
}

.profile-btn {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.5rem;
  color: white;
  cursor: pointer;
  transition: opacity 0.2s ease;
}

.profile-btn:hover {
  opacity: 0.8;
}

.profile-btn span {
  font-size: 0.95rem;
  font-weight: 500;
}

.profile-btn svg {
  width: 24px;
  height: 24px;
  stroke: white;
}

.main-content {
  display: flex;
  min-height: calc(100vh - 70px);
}

.sidebar {
  width: 250px;
  background-color: #0E2954;
  position: fixed;
  top: 70px;
  bottom: 0;
  left: 0;
  display: flex;
  flex-direction: column;
  z-index: 998;
  padding: 0;
  height: calc(100vh - 70px);
  margin-top: 0;
}

.sidebar-logo-group {
  padding: 1.5rem 1rem 1rem 1rem;
  text-align: left;
}

.sidebar-menu {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  padding: 1rem;
  padding-top: 1rem;
  flex: 1 0 auto;
}

.sidebar-logout-btn {
  margin-top: auto;
  margin-bottom: 1.5rem;
  margin-left: 1rem;
  margin-right: 1rem;
  background: #0E2954;
  color: white;
  border: none;
  border-radius: 8px;
  padding: 0.75rem 1rem;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
}
.sidebar-logout-btn:hover {
  background: #1a3a6e;
}

.nav-btn {
  width: 100%;
  padding: 0.75rem 1rem;
  margin-bottom: 0.5rem;
  border: none;
  border-radius: 8px;
  background: transparent;
  color: white;
  text-align: left;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  cursor: pointer;
  transition: all 0.2s ease;
}

.nav-btn:hover {
  background-color: rgba(255, 255, 255, 0.1);
  transform: translateX(5px);
}

.nav-btn.active {
  background-color: #4318FF;
}

.nav-btn i {
  width: 20px;
  transition: transform 0.2s ease;
}

.nav-btn:hover i {
  transform: scale(1.1);
}

.dashboard-content {
  flex: 1;
  padding: 0 2rem;
  overflow-y: auto;
  margin-left: 250px;
  width: calc(100% - 250px);
  height: calc(100vh - 70px);
  transition: margin-left 0.3s ease, width 0.3s ease;
  background-color: white;
  margin-top: 70px;
}

.sub-header {
  background-color: #0E2954;
  padding: 1rem 2rem;
  margin: -2rem -2rem 2rem -2rem;
}

.page-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: #0E2954;
  margin-bottom: 1.5rem;
}

.semester-selector {
  margin-top: 2rem;
}

.semester-selector select {
  padding: 0.5rem 1rem;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 1rem;
  width: 100%;
  max-width: 300px;
}

.stats-chart {
  background: white;
  padding: 1.5rem;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  margin-bottom: 2rem;
  height: 300px;
}

.charts-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.chart-card {
  background: white;
  padding: 1.5rem;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  height: 250px;
  position: relative;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.chart-card h3 {
  margin: 0 0 1rem 0;
  font-size: 0.9rem;
  color: #333;
  flex-shrink: 0;
  text-align: center;
}

.chart-card canvas {
  flex: 1;
  max-height: calc(100% - 2rem);
  width: 100% !important;
  height: 100% !important;
  object-fit: contain;
}

.summary-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1.5rem;
  margin-top: 2rem;
}

.summary-card {
  background: white;
  padding: 1.5rem;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  gap: 1.5rem;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.summary-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.icon-wrapper {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.summary-icon {
  width: 32px;
  height: 32px;
  object-fit: contain;
}

.card-content {
  flex: 1;
}

.card-content h4 {
  margin: 0;
  color: #4b5563;
  font-size: 0.9rem;
  font-weight: 500;
}

.card-content p {
  margin: 0.5rem 0 0;
  font-size: 1.75rem;
  font-weight: 600;
  color: #0E2954;
}

/* Responsive Design */
@media (max-width: 768px) {
  .header {
    padding: 0.5rem 1rem;
  }

  .header-left h1 {
    font-size: 0.9rem;
    line-height: 1.2;
    margin-left: 0.5rem;
  }

  .logo {
    height: 30px;
    width: auto;
  }

  .sidebar {
    transform: translateX(-100%);
    width: 250px;
    position: fixed;
    top: 60px;
    left: 0;
    height: calc(100vh - 60px);
    z-index: 999;
    background-color: #0E2954;
    overflow-y: auto;
  }

  .sidebar.active {
    transform: translateX(0);
  }

  .sidebar-menu {
    padding: 1rem;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    margin-top: 0;
  }

  .nav-btn {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem 1rem;
    width: 100%;
    color: white;
    background: transparent;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .nav-btn.active {
    background-color: #4318FF;
  }

  .nav-btn:hover {
    background-color: rgba(255, 255, 255, 0.1);
    transform: translateX(5px);
  }

  .dashboard-content {
    margin-left: 0;
    width: 100%;
    padding: 1rem;
  }

  .header-right {
    gap: 1rem;
  }

  .profile-btn {
    padding: 0.4rem;
  }

  .profile-btn span {
    display: none;
  }

  .summary-cards {
    grid-template-columns: 1fr;
  }

  .summary-card {
    padding: 1.25rem;
  }

  .icon-wrapper {
    width: 40px;
    height: 40px;
  }

  .summary-icon {
    width: 28px;
    height: 28px;
  }

  .card-content p {
    font-size: 1.5rem;
  }
}

/* Add overlay for mobile sidebar */
.sidebar-overlay {
  display: none;
  position: fixed;
  top: 60px;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 997;
}

@media (max-width: 768px) {
  .sidebar-overlay.active {
    display: block;
  }
}

/* Hamburger Menu Styles */
.hamburger-menu {
  display: none;
  flex-direction: column;
  justify-content: space-between;
  width: 24px;
  height: 20px;
  background: none;
  border: none;
  cursor: pointer;
  padding: 0;
  margin-right: 1rem;
}

.hamburger-menu span {
  display: block;
  width: 100%;
  height: 2px;
  background-color: white;
  transition: all 0.3s ease;
}

.hamburger-menu.active span:first-child {
  transform: translateY(9px) rotate(45deg);
}

.hamburger-menu.active span:nth-child(2) {
  opacity: 0;
}

.hamburger-menu.active span:last-child {
  transform: translateY(-9px) rotate(-45deg);
}

@media screen and (max-width: 768px) {
  .hamburger-menu {
    display: flex;
  }
  
  .header {
    padding: 0 1rem;
  }
  
  .header-left {
    margin-left: 8px;
  }
}

@media (max-width: 1024px) {
  .charts-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .charts-grid {
    grid-template-columns: 1fr;
  }
  
  .chart-card {
    height: 300px;
  }
}

/* Rest of your existing styles */
.dashboard-content {
  padding: 2rem;
  background-color: #F3F4F6;
  min-height: calc(100vh - 70px);
  margin-top: 70px;
  margin-left: 250px;
}

.filter-section {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
}

.semester-select,
.library-select {
  padding: 0.75rem;
  border: 1px solid #D1D5DB;
  border-radius: 0.375rem;
  background-color: white;
  min-width: 250px;
  font-size: 0.875rem;
  color: #374151;
  cursor: pointer;
}

.semester-select:focus,
.library-select:focus {
  outline: none;
  border-color: #0E2954;
  box-shadow: 0 0 0 2px rgba(14, 41, 84, 0.1);
}

.stats-section {
  background-color: white;
  padding: 1.5rem;
  border-radius: 0.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  margin-bottom: 2rem;
}

.chart-container {
  height: 300px;
  margin-top: 1rem;
}

.charts-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.chart-card {
  background-color: white;
  padding: 1.5rem;
  border-radius: 0.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  height: 300px;
  display: flex;
  flex-direction: column;
}

.chart-card h3 {
  margin-bottom: 1rem;
  font-size: 1rem;
  color: #374151;
  flex-shrink: 0;
}

.chart-card canvas {
  flex: 1;
  min-height: 0;
}

.summary-cards {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1.5rem;
}

.summary-card {
  background-color: white;
  padding: 1.5rem;
  border-radius: 0.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  gap: 1rem;
}

.icon-wrapper img {
  width: 48px;
  height: 48px;
}

.summary-info h4 {
  color: #374151;
  margin-bottom: 0.5rem;
}

.summary-info p {
  font-size: 1.5rem;
  font-weight: 600;
  color: #1E40AF;
}

@media (max-width: 1024px) {
  .charts-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .dashboard-content {
    margin-left: 0;
  }
  
  .charts-grid {
    grid-template-columns: 1fr;
  }
  
  .summary-cards {
    grid-template-columns: 1fr;
  }
  
  .filter-section {
    flex-direction: column;
  }
}
</style>