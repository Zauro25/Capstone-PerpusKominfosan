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
          <button class="nav-btn" @click="navigateTo('dashboard')">
            <span>Dashboard</span>
          </button>
          <button class="nav-btn" @click="navigateTo('input-update')">
            <span>Input & Update Data</span>
          </button>
          <button class="nav-btn" @click="navigateTo('pengiriman')">
            <span>Pengiriman Data</span>
          </button>
          <button class="nav-btn active">
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
        <div class="content-container">
          <h1 class="page-title">Validasi dan Revisi dari DPK</h1>
          
          <div class="data-section">
            <h2>Data Perpustakaan</h2>
            <div class="table-container">
              <table class="data-table">
                <thead>
                  <tr>
                    <th>Priode</th>
                    <th>Nama Perpustakaan</th>
                    <th>Nama Kepala</th>
                    <th>Status</th>
                    <th>Aksi</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="library in libraries" :key="library.id">
                    <td>{{ formatPeriode(library.periode) }}</td>
                    <td>{{ library.nama }}</td>
                    <td>{{ library.kepalaPerpustakaan }}</td>
                    <td>
                      <span 
                        class="status-badge"
                        :class="{
                          'menunggu': library.status === 'Menunggu',
                          'revisi': library.status === 'Revisi',
                          'valid': library.status === 'Valid'
                        }"
                      >
                        {{ library.status }}
                      </span>
                    </td>
                    <td class="action-buttons">
                      <button 
                        class="action-btn validate-btn" 
                        @click="updateValidationStatus(library.id, 'Valid')"
                        v-if="library.status !== 'Valid'"
                      >
                        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                          <path d="M20 6L9 17l-5-5"/>
                        </svg>
                        <span>Validasi</span>
                      </button>
                      <button 
                        class="action-btn revise-btn" 
                        @click="updateValidationStatus(library.id, 'Revisi')"
                        v-if="library.status !== 'Revisi'"
                      >
                        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                          <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                          <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                        </svg>
                        <span>Revisi</span>
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </main>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useLibraryStore } from '../store/libraryStore'

export default {
  name: 'ValidasiPage',
  setup() {
    const router = useRouter()
    const libraryStore = useLibraryStore()
    const isSidebarOpen = ref(false)
    const hasUnreadNotifications = ref(true)
    const libraries = ref([])
    const isMobile = ref(false)

    // Load library data with validation status
    const loadLibraryData = () => {
      try {
        // Get all libraries from store
        const allLibraries = libraryStore.libraries
        // Add validation status if not present
        libraries.value = allLibraries.map(lib => ({
          ...lib,
          status: lib.status || 'Menunggu' // Default status if not set
        }))
      } catch (error) {
        console.error('Error loading library data:', error)
      }
    }

    // Update library validation status
    const updateValidationStatus = async (libraryId, newStatus) => {
      try {
        await libraryStore.updateLibrary(libraryId, {
          status: newStatus
        })
        // Reload data to reflect changes
        loadLibraryData()
      } catch (error) {
        console.error('Error updating validation status:', error)
      }
    }

    // Edit library data
    const editLibrary = (id) => {
      router.push(`/input-update?edit=${id}`)
    }

    // Format period string
    const formatPeriode = (periode) => {
      if (!periode) return ''
      const [year, semester] = periode.split('-')
      const semesterText = semester === '1' ? 'Ganjil' : 'Genap'
      const nextYear = parseInt(year) + 1
      return `Semester ${semesterText} ${year}/${nextYear}`
    }

    onMounted(() => {
      checkMobile()
      window.addEventListener('resize', checkMobile)
      document.addEventListener('click', handleClickOutside)
      // Load initial data
      loadLibraryData()
    })

    const checkMobile = () => {
      isMobile.value = window.innerWidth <= 768
    }

    const toggleSidebar = () => {
      isSidebarOpen.value = !isSidebarOpen.value
      if (isMobile.value) {
        document.body.style.overflow = isSidebarOpen.value ? 'hidden' : ''
      }
    }

    const handleClickOutside = (event) => {
      if (isSidebarOpen.value && isMobile.value) {
        const sidebar = document.querySelector('.sidebar')
        const menuToggle = document.querySelector('.hamburger-menu')
        if (!sidebar?.contains(event.target) && !menuToggle?.contains(event.target)) {
          toggleSidebar()
        }
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
        case 'daftar-data-update':
          router.push('/daftar-data-update')
          break
        default:
          router.push(`/${route}`)
      }
    }

    const navigateToNotifications = () => {
      if (isMobile.value) {
        toggleSidebar()
      }
      router.push('/notifications')
    }

    const goToSettings = () => {
      if (isMobile.value) {
        toggleSidebar()
      }
      router.push('/settings')
    }

    const logout = () => {
      // Clear auth tokens
      localStorage.removeItem('authToken')
      sessionStorage.removeItem('authToken')
      localStorage.removeItem('lastRoute')
      router.push('/login')
    }

    // Cleanup on component unmount
    onUnmounted(() => {
      window.removeEventListener('resize', checkMobile)
      document.removeEventListener('click', handleClickOutside)
    })

    return {
      isSidebarOpen,
      hasUnreadNotifications,
      libraries,
      toggleSidebar,
      navigateTo,
      navigateToNotifications,
      goToSettings,
      logout,
      formatPeriode,
      editLibrary,
      updateValidationStatus
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

/* Main Content */
.main-content {
  display: flex;
  height: calc(100vh - 70px);
  margin-top: 70px;
}

/* Main Section */
.main-section {
  flex: 1;
  margin-left: 250px;
  padding: 0;
  background-color: #f8f9fa;
  min-height: calc(100vh - 70px);
  overflow-y: auto;
}

/* Additional styles specific to validation page */
.page-title {
  font-size: 1.5rem;
  color: #1f2937;
  margin: 1.5rem 2rem;
  font-weight: 600;
}

.data-section {
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  margin: 1.5rem 2rem;
  padding: 1.5rem;
}

.data-section h2 {
  color: #1f2937;
  font-size: 1.2rem;
  margin-bottom: 1.5rem;
  font-weight: 600;
}

.table-container {
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

.data-table th {
  background-color: #f9fafb;
  padding: 1rem;
  font-weight: 600;
  color: #4b5563;
  border-bottom: 1px solid #e5e7eb;
}

.data-table td {
  padding: 1rem;
  border-bottom: 1px solid #e5e7eb;
  color: #1f2937;
}

.status-badge {
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.875rem;
  font-weight: 500;
}

.status-badge.menunggu {
  background-color: #FEF3C7;
  color: #92400E;
}

.status-badge.revisi {
  background-color: #FEE2E2;
  color: #991B1B;
}

.status-badge.valid {
  background-color: #D1FAE5;
  color: #065F46;
}

.edit-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 0.5rem;
  color: #4B5563;
  border-radius: 4px;
  transition: all 0.2s ease;
}

.edit-btn:hover {
  background-color: #F3F4F6;
  color: #1F2937;
}

.action-buttons {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 0.75rem;
  border: none;
  border-radius: 6px;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.validate-btn {
  background-color: #D1FAE5;
  color: #065F46;
}

.validate-btn:hover {
  background-color: #A7F3D0;
}

.revise-btn {
  background-color: #FEE2E2;
  color: #991B1B;
}

.revise-btn:hover {
  background-color: #FCA5A5;
}

/* Mobile responsive adjustments */
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

  .page-title {
    margin: 1rem;
    font-size: 1.25rem;
  }

  .data-section {
    margin: 1rem;
    padding: 1rem;
  }

  .data-table th,
  .data-table td {
    padding: 0.75rem 0.5rem;
    font-size: 0.875rem;
  }

  .status-badge {
    padding: 0.25rem 0.5rem;
    font-size: 0.75rem;
  }

  .action-buttons {
    flex-direction: column;
    gap: 0.25rem;
  }

  .action-btn {
    width: 100%;
    justify-content: center;
    padding: 0.375rem 0.5rem;
    font-size: 0.75rem;
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
</style>
