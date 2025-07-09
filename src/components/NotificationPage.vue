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
      <main class="notification-content">
        <h2 class="page-title">Notifikasi</h2>

        <div class="notifications-list">
          <div class="notification-item needs-revision">
            <div class="notification-time">18 Juni 2025, 14:30</div>
            <div class="notification-status revision">Status: Perlu Revisi</div>
            <div class="notification-note">
              <strong>Catatan:</strong> Jumlah Pengunjung bulan Oktober tampaknya tidak realistis (tulis 50 orang, normalnya 200-300 orang). Mohon dicek kembali dan perbaiki. Juga perlu melengkapi data SDM pustakawan yang masih kosong.
            </div>
          </div>

          <div class="notification-item valid">
            <div class="notification-time">17 Juni 2025, 9:30</div>
            <div class="notification-status valid">Status: Valid</div>
          </div>
        </div>
      </main>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

export default {
  name: 'NotificationPage',
  setup() {
    const router = useRouter()
    const isSidebarOpen = ref(false)
    const hasUnreadNotifications = ref(true)

    const toggleSidebar = () => {
      isSidebarOpen.value = !isSidebarOpen.value
      if (window.innerWidth <= 768) {
        document.body.style.overflow = isSidebarOpen.value ? 'hidden' : ''
      }
    }

    const navigateTo = (route) => {
      if (window.innerWidth <= 768) {
        toggleSidebar()
      }
      router.push(`/${route}`)
    }

    const navigateToNotifications = () => {
      router.push('/notifications')
    }

    const goToSettings = () => {
      router.push('/settings')
    }

    const logout = () => {
      localStorage.removeItem('authToken')
      sessionStorage.removeItem('authToken')
      router.push('/login')
    }

    return {
      isSidebarOpen,
      hasUnreadNotifications,
      toggleSidebar,
      navigateTo,
      navigateToNotifications,
      goToSettings,
      logout
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

.notification-btn.active {
  color: #4318FF;
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
.notification-content {
  flex: 1;
  margin-left: 250px;
  padding: 2rem;
  background-color: #F3F4F6;
  min-height: calc(100vh - 70px);
  overflow-y: auto;
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

  .notification-content {
    margin-left: 0;
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

/* Notification styles */
.page-title {
  font-family: 'Poppins', system-ui, -apple-system, sans-serif;
  font-size: 1.5rem;
  font-weight: 600;
  color: #0E2954;
  margin-bottom: 1.5rem;
}

.notifications-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.notification-item {
  background-color: #EFF6FF;
  border-radius: 0.5rem;
  padding: 1.5rem;
  border-left: 4px solid;
}

.notification-item.needs-revision {
  border-left-color: #F97316;
}

.notification-item.valid {
  border-left-color: #22C55E;
  background-color: #F0FDF4;
}

.notification-time {
  font-family: 'Poppins', system-ui, -apple-system, sans-serif;
  font-size: 0.875rem;
  color: #4B5563;
  margin-bottom: 0.5rem;
}

.notification-status {
  font-family: 'Poppins', system-ui, -apple-system, sans-serif;
  font-weight: 500;
  margin-bottom: 0.75rem;
}

.notification-status.revision {
  color: #F97316;
}

.notification-status.valid {
  color: #22C55E;
}

.notification-note {
  font-family: 'Poppins', system-ui, -apple-system, sans-serif;
  color: #1F2937;
  line-height: 1.5;
}

.notification-note strong {
  font-weight: 600;
}
</style> 