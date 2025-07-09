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
          <button class="nav-btn active">
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
      <div 
        class="sidebar-overlay" 
        :class="{ 'active': isSidebarOpen }"
        @click="toggleSidebar"
      ></div>

      <!-- Main Section -->
      <main class="main-section">
        <!-- Top Buttons -->
        <div class="input-header">
          <button 
            class="tab-button"
            @click="navigateTo('input-update')"
          >
            Input Data Perpustakaan
          </button>
          <button 
            class="tab-button active"
          >
            Daftar Data & Update
          </button>
        </div>

        <div class="detail-container">
          <div class="detail-header">
            <h2>Lihat Detail</h2>
            <button class="close-button" @click="goBack">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <line x1="18" y1="6" x2="6" y2="18"></line>
                <line x1="6" y1="6" x2="18" y2="18"></line>
              </svg>
            </button>
          </div>

          <div class="detail-form" v-if="libraryData">
            <div class="form-section">
              <h3>Periode</h3>
              <div class="form-group">
                <label>Semester</label>
                <div class="form-value">{{ libraryData.periode }}</div>
              </div>
            </div>

            <div class="form-section">
              <h3>Identitas Perpustakaan</h3>
              <div class="form-group">
                <label>Nomor Induk</label>
                <div class="form-value">{{ libraryData.nomor_induk }}</div>
              </div>
              <div class="form-group">
                <label>Nama Perpustakaan</label>
                <div class="form-value">{{ libraryData.nama_perpustakaan }}</div>
              </div>
              <div class="form-group">
                <label>Kepala Perpustakaan</label>
                <div class="form-value">{{ libraryData.kepala_perpustakaan }}</div>
              </div>
              <div class="form-group">
                <label>Alamat</label>
                <div class="form-value">{{ libraryData.alamat }}</div>
              </div>
              <div class="form-group">
                <label>Tahun Berdiri</label>
                <div class="form-value">{{ libraryData.tahun_berdiri }}</div>
              </div>
              <div class="form-group">
                <label>Jenis Perpustakaan</label>
                <div class="form-value">{{ libraryData.jenis_perpustakaan }}</div>
              </div>
            </div>

            <div class="form-section">
              <h3>SDM</h3>
              <div class="form-group">
                <label>Jumlah SDM</label>
                <div class="form-value">{{ libraryData.jumlah_sdm }}</div>
              </div>
            </div>

            <div class="form-section">
              <h3>Data Pengunjung</h3>
              <div class="form-group">
                <label>Jumlah Pengunjung</label>
                <div class="form-value">{{ libraryData.jumlah_pengunjung }}</div>
              </div>
              <div class="form-group">
                <label>Anggota Aktif</label>
                <div class="form-value">{{ libraryData.jumlah_anggota }}</div>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useLibraryStore } from '../store/libraryStore'

export default {
  name: 'DetailPage',
  setup() {
    const router = useRouter()
    const libraryStore = useLibraryStore()
    const isSidebarOpen = ref(false)
    const hasUnreadNotifications = ref(true)
    const isMobile = ref(false)
    const libraryData = ref(null)

    onMounted(() => {
      checkMobile()
      window.addEventListener('resize', checkMobile)
      document.addEventListener('click', handleClickOutside)
      loadLibraryData()
    })

    const checkMobile = () => {
      isMobile.value = window.innerWidth <= 768
    }

    const toggleSidebar = () => {
      isSidebarOpen.value = !isSidebarOpen.value
      if (window.innerWidth <= 768) {
        document.body.style.overflow = isSidebarOpen.value ? 'hidden' : ''
      }
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

    const navigateTo = (route) => {
      router.push(`/${route}`)
      localStorage.setItem('lastRoute', route)
    }

    const navigateToNotifications = () => {
      router.push('/notifications')
    }

    const goToSettings = () => {
      router.push('/settings')
    }

    const logout = () => {
      localStorage.removeItem('lastRoute')
      router.push('/login')
    }

    const goBack = () => {
      router.push('/daftar-data-update')
    }

    const loadLibraryData = async () => {
      try {
        const id = parseInt(router.currentRoute.value.params.id)
        const data = await libraryStore.getLibraryById(id)
        if (data) {
          libraryData.value = data
        } else {
          router.push('/daftar-data-update')
        }
      } catch (error) {
        console.error('Error loading library data:', error)
        router.push('/daftar-data-update')
      }
    }

    const formatDate = (dateString) => {
      return new Date(dateString).toLocaleDateString('id-ID', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
      })
    }

    return {
      isSidebarOpen,
      hasUnreadNotifications,
      toggleSidebar,
      navigateTo,
      navigateToNotifications,
      goToSettings,
      logout,
      goBack,
      libraryData,
      formatDate
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

.profile-btn svg {
  width: 24px;
  height: 24px;
  stroke: white;
}

/* Main Content and Sidebar */
.main-content {
  display: flex;
  height: calc(100vh - 70px);
  margin-top: 70px;
}

/* Sidebar Styles */
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

/* Mobile Styles */
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

/* Detail Content Styles */
.main-section {
  flex: 1;
  margin-left: 250px;
  padding: 0;
  background-color: #f8f9fa;
  min-height: calc(100vh - 70px);
  overflow-y: auto;
}

.detail-container {
  max-width: 1200px;
  margin: 0 auto;
}

.input-header {
  display: flex;
  justify-content: space-between;
  padding: 0;
  margin: 0;
  border-bottom: 1px solid #e5e7eb;
  background-color: white;
  position: sticky;
  top: 0;
  z-index: 10;
}

.tab-button {
  background: none;
  border: none;
  color: #1f2937;
  font-size: 1.25rem;
  font-weight: 600;
  padding: 1rem;
  cursor: pointer;
  transition: all 0.2s;
  border-bottom: 2px solid transparent;
  margin: 0;
  flex: 1;
  text-align: center;
}

.tab-button:hover {
  color: #0E2954;
}

.tab-button.active {
  color: #0E2954;
  border-bottom: 2px solid #0E2954;
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem 2rem;
  border-bottom: 1px solid #e5e7eb;
  background-color: white;
}

.detail-header h2 {
  font-size: 1.25rem;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.close-button {
  background: none;
  border: none;
  color: #ef4444;
  cursor: pointer;
  padding: 0.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: opacity 0.2s;
}

.close-button:hover {
  opacity: 0.8;
}

.detail-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  padding: 0 2rem 2rem 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

.form-section {
  margin-bottom: 1.5rem;
}

.form-section h3 {
  font-size: 1rem;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
  margin-bottom: -0.5rem;
}

.form-group {
  display: flex;
  margin-bottom: 0;
  padding: 0.75rem 0;
  border-bottom: 1px solid #e2e2e2;
}

.form-group label {
  flex: 0 0 200px;
  color: #1f2937;
  font-weight: 400;
  font-size: 0.9rem;
}

.form-group .form-value {
  flex: 1;
  color: #1f2937;
  font-size: 0.9rem;
}

/* Overlay */
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

@media (max-width: 768px) {
  .hamburger-menu {
    display: block;
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

  .input-header {
    flex-direction: row;
    border-bottom: none;
  }

  .tab-button {
    font-size: 0.9rem;
    padding: 1rem 0.5rem;
    border-bottom: 1px solid #e5e7eb;
    white-space: normal;
    line-height: 1.2;
    height: auto;
  }

  .tab-button.active {
    background-color: rgba(14, 41, 84, 0.05);
  }

  .detail-form {
    padding: 1rem;
    padding-top: 0.5rem;
  }

  .header-left h1 {
    font-size: 0.9rem;
  }

  .profile-btn span {
    display: none;
  }
}

@media (max-width: 360px) {
  .tab-button {
    font-size: 0.75rem;
    padding: 0.75rem 0.25rem;
  }

  .form-group {
    padding: 0.5rem;
  }

  .form-group label {
    font-size: 0.85rem;
  }

  .form-group .form-value {
    font-size: 0.85rem;
  }
}
</style>