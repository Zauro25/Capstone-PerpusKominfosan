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
        <!-- Tab Navigation -->
        <div class="tab-navigation">
          <button 
            class="tab-button active"
          >
            Input Data Perpustakaan
          </button>
          <button 
            class="tab-button"
            @click="navigateTo('daftar-data-update')"
          >
            Daftar Data & Update
          </button>
        </div>

        <!-- Content Container -->
        <div class="content-container">
          <!-- Input Form -->
          <div class="form-container">
            <form @submit.prevent="handleSubmit" class="input-form">
              <h2>Periode</h2>
              <div class="form-group">
                <select 
                  id="periode"
                  v-model="form.periode" 
                  required
                >
                  <option value="" disabled selected>Pilih Periode</option>
                  <option value="2025-2">Semester Genap 2025/2026</option>
                  <option value="2024-1">Semester Ganjil 2024/2025</option>
                  <option value="2024-2">Semester Genap 2024/2025</option>
                  <option value="2023-1">Semester Ganjil 2023/2024</option>
                  <option value="2023-2">Semester Genap 2023/2024</option>
                </select>
              </div>

              <h2>Identitas Perpustakaan</h2>
              <div class="form-grid">
                <div class="form-group">
                  <label for="nomorInduk">Nomor Induk</label>
                  <input 
                    type="number" 
                    id="nomorInduk" 
                    v-model.number="form.nomorInduk"
                    min="0"
                    @input="validateNomorInduk"
                    placeholder="Masukkan nomor induk"
                    required 
                  />
                  <span v-if="nomorIndukError" class="error-message">{{ nomorIndukError }}</span>
                </div>
                <div class="form-group">
                  <label for="namaPerpustakaan">Nama Perpustakaan</label>
                  <input 
                    type="text" 
                    id="namaPerpustakaan" 
                    v-model="form.namaPerpustakaan" 
                    required 
                  />
                </div>
                <div class="form-group">
                  <label for="kepalaPerpustakaan">Kepala Perpustakaan</label>
                  <input 
                    type="text" 
                    id="kepalaPerpustakaan" 
                    v-model="form.kepalaPerpustakaan" 
                    required 
                  />
                </div>
                <div class="form-group">
                  <label for="tahunBerdiri">Tahun Berdiri</label>
                  <input 
                    type="number" 
                    id="tahunBerdiri" 
                    v-model.number="form.tahunBerdiri" 
                    required 
                    min="1900"
                    max="2100"
                  />
                </div>
              </div>

              <div class="form-group full-width">
                <label for="alamat">Alamat</label>
                <input 
                  type="text" 
                  id="alamat" 
                  v-model="form.alamat" 
                  required 
                />
              </div>

              <div class="form-grid">
                <div class="form-group">
                  <label for="jenisPerpustakaan">Jenis Perpustakaan</label>
                  <select 
                    id="jenisPerpustakaan" 
                    v-model="form.jenisPerpustakaan" 
                    required
                  >
                    <option value="" disabled selected>Pilih Jenis</option>
                    <option value="Umum">Umum</option>
                    <option value="Khusus">Khusus</option>
                    <option value="Sekolah">Sekolah</option>
                    <option value="Perguruan Tinggi">Perguruan Tinggi</option>
                  </select>
                </div>
                <div class="form-group">
                  <label for="jumlahSDM">Jumlah SDM</label>
                  <input 
                    type="number" 
                    id="jumlahSDM" 
                    v-model.number="form.jumlahSDM" 
                    required 
                    min="0"
                  />
                </div>
                <div class="form-group">
                  <label for="jumlahPengunjung">Jumlah Pengunjung</label>
                  <input 
                    type="number" 
                    id="jumlahPengunjung" 
                    v-model.number="form.jumlahPengunjung" 
                    required 
                    min="0"
                  />
                </div>
                <div class="form-group">
                  <label for="jumlahAnggota">Anggota Aktif</label>
                  <input 
                    type="number" 
                    id="jumlahAnggota" 
                    v-model.number="form.jumlahAnggota" 
                    required 
                    min="0"
                  />
                </div>
              </div>

              <div class="button-group">
                <button type="button" class="cancel-btn" @click="handleCancel">
                  Batal
                </button>
                <button type="submit" class="submit-btn">
                  {{ isEditing ? 'Update Data' : 'Input Data' }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </main>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useRoute } from 'vue-router'
import { useLibraryStore } from '../store/libraryStore'
import './styles/form.css'

export default {
  name: 'InputUpdatePage',
  setup() {
    const router = useRouter()
    const route = useRoute()
    const libraryStore = useLibraryStore()
    const isSidebarOpen = ref(false)
    const hasUnreadNotifications = ref(true)
    const isMobile = ref(false)
    const isEditing = ref(false)
    const editingId = ref(null)
    const nomorIndukError = ref('')

    const form = ref({
      periode: '',
      namaPerpustakaan: '',
      alamat: '',
      kepalaPerpustakaan: '',
      jenisPerpustakaan: '',
      tahunBerdiri: '',
      nomorInduk: '',
      jumlahSDM: 0,
      jumlahPengunjung: 0,
      jumlahAnggota: 0
    })

    const validateNomorInduk = (event) => {
      const value = event.target.value
      if (value === '') {
        nomorIndukError.value = 'Nomor induk harus diisi'
        return
      }
      if (isNaN(value)) {
        nomorIndukError.value = 'Nomor induk harus berupa angka'
        return
      }
      if (value < 0) {
        nomorIndukError.value = 'Nomor induk tidak boleh negatif'
        return
      }
      nomorIndukError.value = ''
    }

   onMounted(async () => {
    checkMobile()
    window.addEventListener('resize', checkMobile)
    document.addEventListener('click', handleClickOutside)

    // Cek parameter edit dari URL
      const editId = route.query.edit
      if (editId) {
        try {
          const data = await libraryStore.getLibraryById(parseInt(editId))
          if (data) {
          // Mapping data dari store ke form
            form.value = {
              periode: data.periode,
              namaPerpustakaan: data.nama_perpustakaan,
              alamat: data.alamat,
              kepalaPerpustakaan: data.kepala_perpustakaan,
              jenisPerpustakaan: data.jenis_perpustakaan,
              tahunBerdiri: data.tahun_berdiri,
              nomorInduk: data.nomor_induk,
              jumlahSDM: data.jumlah_sdm,
              jumlahPengunjung: data.jumlah_pengunjung,
              jumlahAnggota: data.jumlah_anggota
            }
            isEditing.value = true
            editingId.value = parseInt(editId)
          }
        } catch (error) {
            console.error('Error loading library data:', error)
            router.push('/daftar-data-update')
        }
      }
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

    const handleSubmit = async () => {
  if (nomorIndukError.value) return

  try {
    if (isEditing.value) {

      await libraryStore.updateLibrary(editingId.value, {
        periode: form.value.periode,
        nama_perpustakaan: form.value.namaPerpustakaan,
        alamat: form.value.alamat,
        kepala_perpustakaan: form.value.kepalaPerpustakaan,
        jenis_perpustakaan: form.value.jenisPerpustakaan,
        tahun_berdiri: form.value.tahunBerdiri,
        nomor_induk: form.value.nomorInduk,
        jumlah_sdm: form.value.jumlahSDM,
        jumlah_pengunjung: form.value.jumlahPengunjung,
        jumlah_anggota: form.value.jumlahAnggota
      })
    } else {
      await libraryStore.addLibrary({
        periode: form.value.periode,
        nama_perpustakaan: form.value.namaPerpustakaan,
        alamat: form.value.alamat,
        kepala_perpustakaan: form.value.kepalaPerpustakaan,
        jenis_perpustakaan: form.value.jenisPerpustakaan,
        tahun_berdiri: form.value.tahunBerdiri,
        nomor_induk: form.value.nomorInduk,
        jumlah_sdm: form.value.jumlahSDM,
        jumlah_pengunjung: form.value.jumlahPengunjung,
        jumlah_anggota: form.value.jumlahAnggota
      })
    }
    router.push('/daftar-data-update')
  } catch (error) {
    console.error('Error saving data:', error)
    alert(error.response?.data?.error || 'Gagal menyimpan data')
  }
} 

    const handleCancel = () => {
      // Reset form to initial state
      form.value = {
        periode: '',
        nomorInduk: '',
        namaPerpustakaan: '',
        kepalaPerpustakaan: '',
        tahunBerdiri: '',
        alamat: '',
        jenisPerpustakaan: '',
        jumlahSDM: 0,
        jumlahPengunjung: 0,
        anggotaAktif: 0
      }
      // Navigate back to daftar data page
      router.push('/daftar-data-update')
    }

    // Check if we're editing an existing entry
    return {
      isSidebarOpen,
      hasUnreadNotifications,
      toggleSidebar,
      navigateTo,
      navigateToNotifications,
      goToSettings,
      logout,
      handleSubmit,
      handleCancel,
      form,
      isEditing,
      editingId,
      nomorIndukError,
      validateNomorInduk
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

/* Tab Navigation */
.tab-navigation {
  display: flex;
  background: white;
  border-bottom: 1px solid #e5e7eb;
  margin-bottom: 1rem;
  justify-content: space-between;
  padding: 0 4rem;
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

  .tab-navigation {
    margin: 0;
    padding: 0 1rem;
  }

  .tab-button {
    flex: 1;
    padding: 1rem 0.5rem;
    font-size: 0.9rem;
    text-align: center;
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