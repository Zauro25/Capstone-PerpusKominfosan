<template>
  <div class="dashboard-container">
    <Header 
      :isSidebarOpen="isSidebarOpen"
      @toggle-sidebar="toggleSidebar"
    />

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

      <!-- Sidebar Overlay -->
      <div 
        class="sidebar-overlay" 
        :class="{ 'active': isSidebarOpen }" 
        @click="toggleSidebar"
      ></div>

      <!-- Main Section -->
      <main class="main-section">
        <!-- Tab Navigation -->
        <div class="tab-navigation">
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

        <!-- Detail Content -->
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

          <div class="detail-form" v-if="library">
            <div class="form-section">
              <h3>Priode</h3>
              <div class="form-group">
                <label>Semester</label>
                <div class="form-value">{{ formatPeriode(library.periode) }}</div>
              </div>
            </div>

            <div class="form-section">
              <h3>Identitas Perpustakaan</h3>
              <div class="form-group">
                <label>Nomor Induk</label>
                <div class="form-value">{{ library.nomorInduk || '226789567' }}</div>
              </div>
              <div class="form-group">
                <label>Nama Perpustakaan</label>
                <div class="form-value">{{ library.nama }}</div>
              </div>
              <div class="form-group">
                <label>Kepala Perpustakaan</label>
                <div class="form-value">{{ library.kepalaPerpustakaan }}</div>
              </div>
              <div class="form-group">
                <label>Alamat</label>
                <div class="form-value">{{ library.alamat || 'Jl.Yang Lurus No.123' }}</div>
              </div>
              <div class="form-group">
                <label>Tahun Berdiri</label>
                <div class="form-value">{{ library.tahunBerdiri }}</div>
              </div>
              <div class="form-group">
                <label>Jenis Perpustakaan</label>
                <div class="form-value">{{ library.jenis || 'Umum' }}</div>
              </div>
            </div>

            <div class="form-section">
              <h3>SDM</h3>
              <div class="form-group">
                <label>Jumlah SDM</label>
                <div class="form-value">{{ library.jumlahSdm || '114' }}</div>
              </div>
            </div>

            <div class="form-section">
              <h3>Data Pengunjung</h3>
              <div class="form-group">
                <label>Jumlah Pengunjung</label>
                <div class="form-value">{{ library.jumlahPengunjung || '3000' }}</div>
              </div>
              <div class="form-group">
                <label>Anggota Aktif</label>
                <div class="form-value">{{ library.anggotaAktif || '1500' }}</div>
              </div>
            </div>
          </div>

          <div v-else class="not-found">
            <p>Data perpustakaan tidak ditemukan</p>
          </div>
        </div>
      </main>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useLibraryStore } from '../store/libraryStore'
import Header from './Header.vue'

export default {
  name: 'DetailPage',
  components: {
    Header
  },
  setup() {
    const router = useRouter()
    const route = useRoute()
    const libraryStore = useLibraryStore()
    const isSidebarOpen = ref(false)
    const library = ref(null)

    const formatPeriode = (periode) => {
      if (!periode) return 'Genap 2025/2026'
      
      const [tahun, semester] = periode.split('-')
      const tahunAkhir = parseInt(tahun) + 1
      const jenisSemester = semester === '1' ? 'Ganjil' : 'Genap'
      
      return `${jenisSemester} ${tahun}/${tahunAkhir}`
    }

    const toggleSidebar = () => {
      isSidebarOpen.value = !isSidebarOpen.value
    }

    const navigateTo = (route) => {
      router.push(`/${route}`)
    }

    const goBack = () => {
      router.push('/daftar-data-update')
    }

    const logout = () => {
      localStorage.removeItem('authToken')
      sessionStorage.removeItem('authToken')
      router.push('/login')
    }

    onMounted(() => {
      const id = parseInt(route.params.id)
      library.value = libraryStore.getLibraryById(id)
      
      if (!library.value) {
        router.push('/daftar-data-update')
      }
    })

    return {
      isSidebarOpen,
      library,
      toggleSidebar,
      navigateTo,
      goBack,
      logout,
      formatPeriode
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

/* Main Content */
.main-content {
  display: flex;
  height: calc(100vh - 70px);
  margin-top: 70px;
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

.main-section {
  flex: 1;
  margin-left: 250px;
  padding: 0;
  background-color: #f8f9fa;
  min-height: calc(100vh - 70px);
  overflow-y: auto;
  width: calc(100% - 250px);
}

.tab-navigation {
  display: flex;
  background: white;
  border-bottom: 1px solid #e5e7eb;
  margin-bottom: 1rem;
  justify-content: space-between;
  padding: 0 4rem;
  position: sticky;
  top: 0;
  z-index: 10;
}

.tab-button {
  padding: 1rem 2rem;
  border: none;
  background: transparent;
  color: #1f2937;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
  font-size: 1.1rem;
}

.tab-button:hover {
  color: #0E2954;
}

.tab-button.active {
  color: #0E2954;
  font-weight: 600;
}

.tab-button.active::after {
  content: '';
  position: absolute;
  bottom: -1px;
  left: 0;
  width: 100%;
  height: 2px;
  background-color: #0E2954;
}

.detail-container {
  max-width: 1200px;
  margin: 2rem auto;
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem 2rem;
  border-bottom: 1px solid #e5e7eb;
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
  padding: 2rem;
}

.form-section {
  margin-bottom: 2rem;
}

.form-section h3 {
  font-size: 1.1rem;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 1rem;
}

.form-group {
  display: flex;
  padding: 0.75rem 0;
  border-bottom: 1px solid #e5e7eb;
}

.form-group:last-child {
  border-bottom: none;
}

.form-group label {
  flex: 0 0 200px;
  color: #1f2937;
  font-weight: 500;
}

.form-group .form-value {
  flex: 1;
  color: #4b5563;
}

.not-found {
  text-align: center;
  padding: 2rem;
  color: #64748b;
}

@media (max-width: 768px) {
  .sidebar {
    transform: translateX(-100%);
    transition: transform 0.3s ease;
  }

  .sidebar.active {
    transform: translateX(0);
  }

  .main-section {
    margin-left: 0;
    width: 100%;
  }

  .tab-navigation {
    padding: 0 1rem;
  }

  .tab-button {
    padding: 1rem;
    font-size: 0.9rem;
  }

  .detail-container {
    margin: 1rem;
  }

  .detail-header {
    padding: 1rem;
  }

  .detail-form {
    padding: 1rem;
  }

  .form-group {
    flex-direction: column;
    gap: 0.5rem;
  }

  .form-group label {
    flex: none;
  }
}
</style> 