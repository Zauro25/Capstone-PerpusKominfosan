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

        <!-- Form Content -->
        <div class="form-container">
          <div class="form-content">
            <h3>Periode</h3>
            <div class="form-group">
              <select v-model="formData.periode">
                <option value="" disabled selected>Pilih Periode</option>
                <option value="2025-2">Semester Genap 2025/2026</option>
                <option value="2025-1">Semester Ganjil 2025/2026</option>
                <option value="2024-2">Semester Genap 2024/2025</option>
                <option value="2024-1">Semester Ganjil 2024/2025</option>
              </select>
            </div>

            <h3>Identitas Perpustakaan</h3>
            <div class="form-row">
              <div class="form-group">
                <label>Nomor Induk</label>
                <input type="text" v-model="formData.nomorInduk" />
              </div>
              <div class="form-group">
                <label>Nama Perpustakaan</label>
                <input type="text" v-model="formData.nama" />
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label>Kepala Perpustakaan</label>
                <input type="text" v-model="formData.kepalaPerpustakaan" />
              </div>
              <div class="form-group">
                <label>Tahun Berdiri</label>
                <input type="text" v-model="formData.tahunBerdiri" />
              </div>
            </div>

            <div class="form-group">
              <label>Alamat</label>
              <input type="text" v-model="formData.alamat" />
            </div>

            <div class="form-row">
              <div class="form-group">
                <label>Jenis Perpustakaan</label>
                <select v-model="formData.jenis">
                  <option value="" disabled selected>Pilih Jenis</option>
                  <option value="Umum">Umum</option>
                  <option value="Khusus">Khusus</option>
                  <option value="Sekolah">Sekolah</option>
                  <option value="Perguruan Tinggi">Perguruan Tinggi</option>
                </select>
              </div>
              <div class="form-group">
                <label>Jumlah SDM</label>
                <input type="text" v-model="formData.jumlahSdm" />
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label>Jumlah Pengunjung</label>
                <input type="text" v-model="formData.jumlahPengunjung" />
              </div>
              <div class="form-group">
                <label>Anggota Aktif</label>
                <input type="text" v-model="formData.anggotaAktif" />
              </div>
            </div>

            <div class="form-actions">
              <button type="button" class="btn-cancel" @click="handleCancel">Batal</button>
              <button type="submit" class="btn-submit">{{ isEditMode ? 'Update' : 'Simpan' }}</button>
            </div>
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
  name: 'InputUpdatePage',
  components: {
    Header
  },
  setup() {
    const router = useRouter()
    const route = useRoute()
    const libraryStore = useLibraryStore()
    const isSidebarOpen = ref(false)
    const isEditMode = ref(false)
    
    const formData = ref({
      periode: '',
      nomorInduk: '',
      nama: '',
      kepalaPerpustakaan: '',
      alamat: '',
      tahunBerdiri: '',
      jenis: '',
      jumlahSdm: '',
      jumlahPengunjung: '',
      anggotaAktif: ''
    })

    const toggleSidebar = () => {
      isSidebarOpen.value = !isSidebarOpen.value
    }

    const navigateTo = (route) => {
      router.push(`/${route}`)
    }

    const goToSettings = () => {
      router.push('/profile')
    }

    const logout = () => {
      localStorage.removeItem('authToken')
      sessionStorage.removeItem('authToken')
      router.push('/login')
    }

    const loadLibraryData = () => {
      const id = route.query.id
      if (id && route.query.mode === 'edit') {
        isEditMode.value = true
        const library = libraryStore.getLibraryById(parseInt(id))
        if (library) {
          formData.value = {
            periode: library.periode || '',
            nomorInduk: library.nomorInduk || '',
            nama: library.nama || '',
            kepalaPerpustakaan: library.kepalaPerpustakaan || '',
            alamat: library.alamat || '',
            tahunBerdiri: library.tahunBerdiri || '',
            jenis: library.jenis || '',
            jumlahSdm: library.jumlahSdm || '',
            jumlahPengunjung: library.jumlahPengunjung || '',
            anggotaAktif: library.anggotaAktif || ''
          }
        }
      }
    }

    const handleSubmit = () => {
      if (isEditMode.value) {
        const id = parseInt(route.query.id)
        libraryStore.updateLibrary(id, formData.value)
      } else {
        libraryStore.addLibrary(formData.value)
      }
      router.push('/daftar-data-update')
    }

    const handleCancel = () => {
      router.push('/daftar-data-update')
    }

    onMounted(() => {
      loadLibraryData()
    })

    return {
      isSidebarOpen,
      isEditMode,
      formData,
      toggleSidebar,
      navigateTo,
      goToSettings,
      logout,
      handleSubmit,
      handleCancel
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

/* Main Section */
.main-section {
  flex: 1;
  margin-left: 250px;
  padding: 0;
  background-color: #f8f9fa;
  min-height: calc(100vh - 70px);
  overflow-y: auto;
  width: calc(100% - 250px);
}

/* Tab Navigation */
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

/* Form Styles */
.form-container {
  max-width: 1200px;
  margin: 2rem auto;
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  padding: 2rem;
}

.form-header {
  margin-bottom: 2rem;
}

.form-header h2 {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 600;
  color: #1f2937;
}

.input-form {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.form-section {
  margin-bottom: 2rem;
}

.form-section h3 {
  font-size: 1.25rem;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 1.5rem;
}

.form-row {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 2rem;
  margin-bottom: 1.5rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-group.full-width {
  grid-column: span 2;
}

.form-group label {
  font-size: 1rem;
  font-weight: 500;
  color: #1f2937;
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  font-size: 1rem;
  color: #1f2937;
  background-color: white;
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: #0E2954;
  box-shadow: 0 0 0 2px rgba(14, 41, 84, 0.1);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 2rem;
  padding-top: 1rem;
  border-top: 1px solid #e5e7eb;
}

.btn-cancel {
  padding: 0.75rem 2rem;
  background-color: white;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  color: #4b5563;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-submit {
  padding: 0.75rem 2rem;
  background-color: #0E2954;
  border: none;
  border-radius: 6px;
  color: white;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-cancel:hover {
  background-color: #f9fafb;
  border-color: #d1d5db;
}

.btn-submit:hover {
  background-color: #1a3a6e;
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

  .form-container {
    margin: 1rem;
    padding: 1rem;
  }

  .form-row {
    grid-template-columns: 1fr;
    gap: 1rem;
  }

  .form-group.full-width {
    grid-column: auto;
  }

  .form-actions {
    flex-direction: column;
  }

  .btn-cancel,
  .btn-submit {
    width: 100%;
  }
}

.form-container {
  padding: 2rem;
  background: white;
}

.form-content {
  max-width: 1000px;
  margin: 0 auto;
}

h3 {
  font-size: 1.25rem;
  font-weight: 600;
  color: #1f2937;
  margin: 2rem 0 1rem 0;
}

h3:first-child {
  margin-top: 0;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
  margin-bottom: 1rem;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  font-size: 1rem;
  font-weight: 500;
  color: #1f2937;
  margin-bottom: 0.5rem;
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  font-size: 1rem;
  color: #1f2937;
  background-color: white;
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: #0E2954;
  box-shadow: 0 0 0 2px rgba(14, 41, 84, 0.1);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 2rem;
  padding-top: 1rem;
  border-top: 1px solid #e5e7eb;
}

.btn-cancel {
  padding: 0.75rem 2rem;
  background-color: white;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  color: #4b5563;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-submit {
  padding: 0.75rem 2rem;
  background-color: #0E2954;
  border: none;
  border-radius: 6px;
  color: white;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-cancel:hover {
  background-color: #f9fafb;
  border-color: #d1d5db;
}

.btn-submit:hover {
  background-color: #1a3a6e;
}

@media (max-width: 768px) {
  .form-container {
    padding: 1rem;
  }

  .form-row {
    grid-template-columns: 1fr;
    gap: 1rem;
  }

  .form-actions {
    flex-direction: column;
  }

  .btn-cancel,
  .btn-submit {
    width: 100%;
  }
}
</style>