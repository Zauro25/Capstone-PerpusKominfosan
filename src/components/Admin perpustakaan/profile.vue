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
            <span>{{ userProfile?.nama_lengkap || 'Admin Perpustakaan' }}</span>
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
            <button class="nav-btn active" @click="dashboard">
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
  
        <!-- Profile Content -->
       <div class="profile-content">
          <div class="profile-card">
            <div class="profile-header">
              <div class="profile-avatar">
                <div class="avatar-circle">
                  <span class="avatar-icon">👤</span>
                </div>
              </div>
              <div class="profile-info">
                <h3 class="profile-name">{{ userProfile.nama_lengkap || '-' }}</h3>
                <p class="profile-role">{{ userProfile.user_type || '-' }}</p>
              </div>
            </div>
  
            <div class="profile-form">
              <div class="form-group">
                <h3 style="color: #0E2954;">Nama Lengkap</h3>
                <div class="detail-value" style="color: black;">{{ userProfile?.nama_lengkap || '-' }}</div>
              </div>
              <div class="form-group">
                <h3 style="color: #0E2954;">Email</h3>
                <div class="detail-value" style="color: black;">{{ userProfile.email || '-' }}</div>
              </div>
              <div class="form-group">
                <h3 style="color: #0E2954;">Nomor Telepon</h3>
                <div class="detail-value" style="color: black;">{{ userProfile.no_telepon || '-' }}</div>
              </div>
              <div class="form-group">
                <h3 style="color: #0E2954;">Username</h3>
                <div class="detail-value" style="color: black;">{{ userProfile.username || '-' }}</div>
              </div>
              <div class="form-actions">
                <button @click="handleEditProfile" class="edit-btn">Edit Profil</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  import { useProfileStore } from '../../store/profilestore.js'
  import { computed, ref, onMounted } from 'vue'
  import { useRouter } from 'vue-router'
  
  export default {
    name: 'Profile',
    setup() {
      const router = useRouter()
      const profileStore = useProfileStore()
  
      const isSidebarOpen = ref(false)
      const hasUnreadNotifications = ref(false)
  
      const toggleSidebar = () => {
        isSidebarOpen.value = !isSidebarOpen.value
      }
  
      const handleEditProfile = () => {
        router.push('/profile/edit')
      }
  
      const navigateTo = (route) => {
        router.push(`/${route}`)
      }
  
      const dashboard = () => {
        router.push('/dashboard')
      }
  
      const goToSettings = () => {
        router.push('/profile')
      }
  
      const logout = () => {
        localStorage.removeItem('authToken')
        sessionStorage.removeItem('authToken')
        router.push('/login')
      }
  
      const navigateToNotifications = () => {
        router.push('/notifications')
      }
  
      const userProfile = computed(() => profileStore.userProfile)
  
      onMounted(async () => {
        const token = localStorage.getItem('authToken') || sessionStorage.getItem('authToken')
        const userType = localStorage.getItem('userType') || sessionStorage.getItem('userType')
        console.log('userProfile value:', profileStore.userProfile)
  
        if (!token || userType !== 'admin_perpustakaan') {
          return router.push('/login')
        }
  
        try {
          await profileStore.fetchUserProfile()
          console.log('User profile:', profileStore.userProfile)
        } catch (error) {
          console.error('Error loading profile:', error)
        }
      })
  
      return {
        isSidebarOpen,
        hasUnreadNotifications,
        toggleSidebar,
        handleEditProfile,
        navigateTo,
        dashboard,
        goToSettings,
        logout,
        navigateToNotifications,
        userProfile,
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
  
  /* Profile Content Styles */
  .profile-content {
    padding: 2rem;
  }
  
  .profile-card {
    background-color: white;
    border-radius: 8px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    overflow: hidden;
    max-width: 800px;
  }
  
  .profile-header {
    display: flex;
    align-items: center;
    gap: 1.5rem;
    padding: 2rem 2rem 1.5rem 2rem;
    border-bottom: 1px solid #e2e8f0;
  }
  
  .profile-avatar {
    flex-shrink: 0;
  }
  
  .avatar-circle {
    width: 80px;
    height: 80px;
    background-color: #3b82f6;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  .avatar-circle .avatar-icon {
    color: white;
    font-size: 2rem;
  }
  
  .profile-info {
    flex: 1;
  }
  
  .profile-name {
    font-size: 1.5rem;
    font-weight: 600;
    color: #1e293b;
    margin: 0 0 0.5rem 0;
  }
  
  .profile-role {
    font-size: 1rem;
    color: #64748b;
    margin: 0;
  }
  
  .profile-form {
    padding: 2rem;
  }
  
  .form-group {
    margin-bottom: 1.5rem;
  }
  
  .form-label {
    display: block;
    font-size: 0.9rem;
    font-weight: 500;
    color: #374151;
    margin-bottom: 0.5rem;
  }
  
  .form-input {
    width: 100%;
    padding: 0.75rem 1rem;
    border: 1px solid #d1d5db;
    border-radius: 6px;
    font-size: 0.9rem;
    color: #374151;
    background-color: #f9fafb;
    transition: border-color 0.2s ease;
  }
  
  .form-input:focus {
    outline: none;
    border-color: #3b82f6;
    background-color: white;
  }
  
  .form-actions {
    display: flex;
    justify-content: flex-end;
    margin-top: 2rem;
  }
  
  .edit-btn {
    background-color: #0E2954;
    color: white;
    border: none;
    padding: 0.75rem 1.5rem;
    border-radius: 6px;
    font-size: 0.9rem;
    font-weight: 500;
    cursor: pointer;
    transition: background-color 0.2s ease;
  }
  
  .edit-btn:hover {
    background-color: #1a3a6e;
  }
  
  /* Responsive Design */
  .main-content {
    display: flex;
    min-height: calc(100vh - 70px);
    padding-top: 70px; /* biar ga ketutupan header */
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
    height: calc(100vh - 70px);
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
  .profile-content {
      flex: 1;
    margin-left: 250px;
    padding: 0 2rem;
    overflow-y: auto;
    width: calc(100% - 250px);
    height: calc(100vh - 70px);
    display: flex;
    justify-content: center; /* untuk center secara horizontal */
  }
  .profile-card {
    width: 100%;
    max-width: 6500px; /* biar gak terlalu lebar */
    margin: 0 auto; /* center jika ada ruang kosong kiri-kanan */
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
  
  
  
  @media (max-width: 640px) {
    .sidebar {
      width: 100%;
      position: relative;
      height: auto;
    }
    
    .main-content {
      margin-left: 0;
    }
    
    .header-content {
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
    
    .avatar-circle {
      width: 60px;
      height: 60px;
    }
    
    .avatar-circle .avatar-icon {
      font-size: 1.5rem;
    }
    
    .profile-name {
      font-size: 1.25rem;
    }
  }
  </style>  