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
        <img src="../../assets/logo-sidapus.png" alt="Logo" class="logo" />
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

    <div class="main-content">
      <!-- Sidebar -->
      <aside class="sidebar" :class="{ 'active': isSidebarOpen }">
        <nav class="sidebar-menu">
          <button class="nav-btn" @click="navigateTo('dashboard')">
            <span>Dashboard</span>
          </button>
          <button class="nav-btn" @click="navigateTo('input-update')">
            <span>Input & Update Data</span>
          </button>
          <button class="nav-btn active">
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
      <div 
        class="sidebar-overlay" 
        :class="{ 'active': isSidebarOpen }" 
        @click="toggleSidebar"
      ></div>

      <!-- Main Section -->
      <main class="main-section">
        <div class="content-wrapper">
          <div class="page-header">
            <h2>Pengiriman Data</h2>
          </div>
          
          <div class="data-table">
            <h3 class="table-title">Data Perpustakaan</h3>
            <div class="table-container">
              <table>
                <thead>
                  <tr>
                    <th>Periode</th>
                    <th>Nama Perpustakaan</th>
                    <th>Nama Kepala</th>
                    <th>Tahun Berdiri</th>
                    <th>Status</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="(item, index) in libraryDataByPeriod" :key="item.id">
                    <td class="periode-cell">
                      <div class="checkbox-container">
                        <input 
                          type="checkbox" 
                          :id="'checkbox-' + index"
                          v-model="selectedItems[index]"
                          :disabled="item.status === 'Sudah Dikirim'"
                        >
                        <label :for="'checkbox-' + index">{{ formatPeriode(item.periode) }}</label>
                      </div>
                    </td>
                    <td>{{ item.nama }}</td>
                    <td>{{ item.kepalaPerpustakaan }}</td>
                    <td>{{ item.tahunBerdiri }}</td>
                    <td>
                      <span 
                        class="status-badge"
                        :class="{ 'sent': item.status === 'Sudah Dikirim' }"
                      >
                        {{ item.status || 'Belum Dikirim' }}
                      </span>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

          <div class="action-buttons">
            <button 
              class="btn btn-cancel" 
              @click="cancelSelection"
              :disabled="!hasSelectedItems"
            >
              Batal
            </button>
            <button 
              class="btn btn-send" 
              @click="sendData"
              :disabled="!hasSelectedItems"
            >
              Kirim
            </button>
          </div>
        </div>
      </main>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useLibraryStore } from '../../store/libraryStore.js'
import { useRouter } from 'vue-router'
import { useSubmissionStore } from '../../store/submissionStore.js'

export default {
  name: 'PengirimanData',
  setup() {
    const router = useRouter()
    const libraryStore = useLibraryStore()
    const selectedItems = ref({})
    const isSidebarOpen = ref(false)
    const hasUnreadNotifications = ref(true)
    const isMobile = ref(false)
    const submissionStore = useSubmissionStore()

    // Computed property to get library data grouped by period
    const libraryDataByPeriod = computed(() => {
      const data = libraryStore.libraries
      // Sort by period in descending order (newest first)
      return data.sort((a, b) => {
        const [yearA, semA] = a.periode.split('-')
        const [yearB, semB] = b.periode.split('-')
        return yearB - yearA || semB - semA
      })
    })

    // Computed property to check if any items are selected
    const hasSelectedItems = computed(() => {
      return Object.values(selectedItems.value).some(value => value)
    })

    // Format period string
    const formatPeriode = (periode) => {
      const [year, semester] = periode.split('-')
      return semester === '1' 
        ? `Semester Ganjil ${year}/${parseInt(year) + 1}`
        : `Semester Genap ${year}/${parseInt(year) + 1}`
    }

    // Send selected data
    const sendData = async () => {
      try {
        const selectedData = libraryDataByPeriod.value.filter((_, index) => selectedItems.value[index])
        
        // Update both library and submission stores
        selectedData.forEach(data => {
          // Update library store
          libraryStore.updateLibrary(data.id, {
            ...data,
            status: 'Sudah Dikirim'
          })

          // Update submission store
          submissionStore.submitData({
            id: data.id,
            periode: data.periode,
            nama: data.nama,
            kepalaPerpustakaan: data.kepalaPerpustakaan,
            tahunBerdiri: data.tahunBerdiri,
            status: 'Sudah Dikirim'
          })
        })

        // Clear selections
        selectedItems.value = {}
        alert('Data berhasil dikirim!')
      } catch (error) {
        console.error('Error sending data:', error)
        alert('Gagal mengirim data. Silakan coba lagi.')
      }
    }

    // Cancel selection
    const cancelSelection = () => {
      selectedItems.value = {}
    }

    // Methods for navigation and UI
    const toggleSidebar = () => {
      isSidebarOpen.value = !isSidebarOpen.value
      if (window.innerWidth <= 768) {
        document.body.style.overflow = isSidebarOpen.value ? 'hidden' : ''
      }
    }

    const navigateTo = (route) => {
      // Close sidebar on navigation if mobile
      if (isMobile.value) {
        toggleSidebar()
      }

      switch (route) {
        case 'dashboard':
          router.push('/dashboard')
          break
        case 'input-update':
          router.push('/input-update')
          break
        case 'pengiriman':
          router.push('/pengiriman')
          break
        case 'validasi':
          router.push('/validasi')
          break
        default:
          router.push(`/${route}`)
      }
    }

    const navigateToNotifications = () => {
      router.push('/notifications')
    }

    const goToSettings = () => {
      const userType = localStorage.getItem('userType') || sessionStorage.getItem('userType')
      switch(userType) {
        case 'admin_perpustakaan':
          router.push('/profile')
          break
        case 'executive':
          router.push('/profile-executive')
          break
        case 'admin_dpk':
          router.push('/profile-dpk')
          break
        default:
          router.push('/profile')
      }
    }

    const logout = () => {
      // Clear auth tokens
      localStorage.removeItem('authToken')
      sessionStorage.removeItem('authToken')
      localStorage.removeItem('lastRoute')
      router.push('/login')
    }

    onMounted(() => {
      checkMobile()
      window.addEventListener('resize', checkMobile)
      document.addEventListener('click', handleClickOutside)
    })

    const checkMobile = () => {
      isMobile.value = window.innerWidth <= 768
    }

    const handleClickOutside = (event) => {
      if (isSidebarOpen.value && window.innerWidth <= 768) {
        const sidebar = document.querySelector('.sidebar')
        const menuToggle = document.querySelector('.hamburger-menu')
        if (!sidebar?.contains(event.target) && !menuToggle?.contains(event.target)) {
          toggleSidebar()
        }
      }
    }

    const submitData = (data) => {
      submissionStore.submitData({
        id: Date.now(),
        periode: data.periode,
        namaPerpustakaan: data.namaPerpustakaan,
        kepala: data.kepala,
        // ... other data fields ...
      })
      // Show success message or handle UI feedback
    }

    return {
      selectedItems,
      isSidebarOpen,
      hasUnreadNotifications,
      libraryDataByPeriod,
      hasSelectedItems,
      formatPeriode,
      sendData,
      cancelSelection,
      toggleSidebar,
      navigateTo,
      navigateToNotifications,
      goToSettings,
      logout,
      submitData
    }
  }
}
</script>

<style scoped>
/* Reset default margins and padding */
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
}

/* Header Styles */
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

/* Hamburger Menu */
.hamburger-menu {
  display: none;
  background: none;
  border: none;
  cursor: pointer;
  padding: 0.5rem;
  z-index: 1001;
}

.hamburger-menu span {
  display: block;
  width: 25px;
  height: 3px;
  background-color: white;
  margin: 5px 0;
  transition: all 0.3s ease;
}

.hamburger-menu.active span:nth-child(1) {
  transform: rotate(45deg) translate(5px, 5px);
}

.hamburger-menu.active span:nth-child(2) {
  opacity: 0;
}

.hamburger-menu.active span:nth-child(3) {
  transform: rotate(-45deg) translate(5px, -5px);
}

/* Sidebar */
.sidebar {
  width: 250px;
  background-color: #0E2954;
  position: fixed;
  top: 70px;
  bottom: 0;
  left: 0;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  z-index: 998;
  padding: 0;
  height: calc(100vh - 70px);
  margin-top: 0;
}

.sidebar-menu {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  padding: 1rem;
  padding-top: 1rem;
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
  font-family: inter, sans-serif;
  font-size: 1rem;
}

.nav-btn:hover {
  background-color: rgba(255, 255, 255, 0.1);
  transform: translateX(5px);
}

.nav-btn.active {
  background-color: #4318FF;
}

.sidebar-logout-btn {
  padding: 0.75rem 1rem;
  margin: 1rem;
  margin-top: auto;
  border: none;
  border-radius: 8px;
  background: transparent;
  color: white;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  cursor: pointer;
  transition: all 0.2s ease;
  font-family: inter, sans-serif;
  font-size: 1rem;
}

.sidebar-logout-btn:hover {
  background-color: rgba(255, 255, 255, 0.1);
  transform: translateX(5px);
}

/* Sidebar Overlay */
.sidebar-overlay {
  display: none;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 997;
}

.sidebar-overlay.active {
  display: block;
}

/* Main Content Styles */
.main-content {
  display: flex;
  height: calc(100vh - 70px);
  margin-top: 70px;
}

/* Main Section */
.main-section {
  flex: 1;
  margin-left: 250px;
  background-color: #f8f9fa;
  min-height: calc(100vh - 70px);
  overflow-y: auto;
  scrollbar-width: none; /* Firefox */
  -ms-overflow-style: none; /* Internet Explorer 10+ */
}

/* WebKit browsers (Chrome, Safari) */
.main-section::-webkit-scrollbar {
  display: none;
}

.content-wrapper {
  padding: 2rem;
  height: 100%;
}

.page-header {
  margin-bottom: 2rem;
}

.page-header h2 {
  font-size: 1.5rem;
  color: #1F2937;
  font-weight: 600;
  margin: 0;
}

.data-table {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.table-title {
  font-size: 1.25rem;
  color: #1e293b;
  margin-bottom: 1.5rem;
  font-weight: 600;
}

.table-container {
  overflow-x: auto;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
  scrollbar-width: none; /* Firefox */
  -ms-overflow-style: none; /* Internet Explorer 10+ */
}

table {
  width: 100%;
  border-collapse: collapse;
  min-width: 800px;
}

thead {
  background-color: #f8fafc;
}

th {
  text-align: left;
  padding: 1rem;
  font-weight: 600;
  color: #475569;
  border-bottom: 2px solid #e2e8f0;
}

td {
  padding: 1rem;
  border-bottom: 1px solid #e2e8f0;
  color: #1e293b;
}

tr:hover {
  background-color: #f8fafc;
}

.periode-cell {
  min-width: 250px;
}

.checkbox-container {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.checkbox-container input[type="checkbox"] {
  width: 18px;
  height: 18px;
  border: 2px solid #333;
  border-radius: 4px;
  cursor: pointer;
  appearance: none;
  -webkit-appearance: none;
  -moz-appearance: none;
  background-color: white;
  position: relative;
  margin: 0;
  transition: all 0.2s ease;
}

.checkbox-container input[type="checkbox"]:checked {
  background-color: white;
  border-color: #333;
}

.checkbox-container input[type="checkbox"]:checked::before {
  content: "✓";
  position: absolute;
  color: #333;
  font-size: 16px;
  font-weight: 900;
  left: 2px;
  top: -2px;
  text-shadow: none;
  z-index: 2;
}

.checkbox-container input[type="checkbox"]:hover {
  border-color: #333;
}

/* Update status badge styles */
.status-badge {
  display: inline-block;
  padding: 0.5rem 1rem;
  border-radius: 9999px;
  font-size: 0.875rem;
  font-weight: 500;
  background-color: #fee2e2;
  color: #dc2626;
}

.status-badge.sent {
  background-color: #dcfce7;
  color: #16a34a;
}

/* Update checkbox styles for disabled state */
.checkbox-container input[type="checkbox"]:disabled {
  background-color: #f1f5f9;
  border-color: #cbd5e1;
  cursor: not-allowed;
  opacity: 0.5;
}

.checkbox-container input[type="checkbox"]:disabled:checked::before {
  opacity: 0.5;
  color: #94a3b8;
}

.checkbox-container input[type="checkbox"]:disabled + label {
  color: #94a3b8;
  cursor: not-allowed;
}

.checkbox-container label {
  font-size: 0.875rem;
  color: #475569;
  cursor: pointer;
  user-select: none;
}

.action-buttons {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 1.5rem;
  padding: 1rem;
}

.btn {
  padding: 0.75rem 2rem;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-cancel {
  background-color: #f1f5f9;
  color: #64748b;
  border: 1px solid #e2e8f0;
}

.btn-cancel:hover:not(:disabled) {
  background-color: #e2e8f0;
}

.btn-send {
  background-color: #2563eb;
  color: white;
  border: none;
}

.btn-send:hover:not(:disabled) {
  background-color: #1d4ed8;
}

/* Mobile Responsive */
@media (max-width: 768px) {
  .hamburger-menu {
    display: block;
  }

  .header-left h1 {
    font-size: 0.9rem;
  }

  .profile-btn span {
    display: none;
  }

  .sidebar {
    transform: translateX(-100%);
    transition: transform 0.3s ease;
  }

  .sidebar.active {
    transform: translateX(0);
  }

  .main-section {
    margin-left: 0;
  }

  .content-wrapper {
    padding: 1rem;
  }

  .data-table {
    margin: 0;
  }

  .action-buttons {
    padding: 1rem 0;
  }
}

/* For very small screens */
@media (max-width: 360px) {
  .header-left h1 {
    font-size: 0.8rem;
  }

  .logo {
    height: 25px;
  }
}

/* Ensure the body doesn't show scrollbar */
body {
  overflow: hidden;
  margin: 0;
}
</style>