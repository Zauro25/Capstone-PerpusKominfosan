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
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M18 8C18 6.4087 17.3679 4.88258 16.2426 3.75736C15.1174 2.63214 13.5913 2 12 2C10.4087 2 8.88258 2.63214 7.75736 3.75736C6.63214 4.88258 6 6.4087 6 8C6 15 3 17 3 17H21C21 17 18 15 18 8Z" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M13.73 21C13.5542 21.3031 13.3019 21.5547 12.9982 21.7295C12.6946 21.9044 12.3504 21.9965 12 21.9965C11.6496 21.9965 11.3054 21.9044 11.0018 21.7295C10.6982 21.5547 10.4458 21.3031 10.27 21" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          <div class="notification-dot"></div>
        </div>
        <div class="profile-btn" @click="goToSettings">
          <span>Admin Perpustakaan</span>
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M20 21V19C20 17.9391 19.5786 16.9217 18.8284 16.1716C18.0783 15.4214 17.0609 15 16 15H8C6.93913 15 5.92172 15.4214 5.17157 16.1716C4.42143 16.9217 4 17.9391 4 19V21" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M12 11C14.2091 11 16 9.20914 16 7C16 4.79086 14.2091 3 12 3C9.79086 3 8 4.79086 8 7C8 9.20914 9.79086 11 12 11Z" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
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
          <button class="nav-btn active" @click="navigateTo('validasi')">
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

      <!-- Validation Content -->
      <div class="validation-content">
        <h2>Validasi dan Revisi dari DPK</h2>
        
        <div class="data-section">
          <h3>Data Perpustakaan</h3>
          
          <div class="table-container">
            <table>
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
                        'status-not-sent': library.status === 'Belum Dikirim' || !library.status,
                        'status-waiting': library.status === 'Sudah Dikirim' || library.status === 'Menunggu',
                        'status-revision': library.status === 'Revisi',
                        'status-valid': library.status === 'Valid',
                        'status-revised': library.status === 'Telah Direvisi'
                      }"
                    >
                      {{ library.status === 'Sudah Dikirim' ? 'Menunggu' : (library.status || 'Belum Dikirim') }}
                    </span>
                  </td>
                  <td>
                    <button 
                      class="edit-btn" 
                      @click="editLibrary(library)"
                    >
                      <svg width="16" height="16" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <path d="M12.7038 0.0476074L15.9524 3.29622L5.93025 13.3184L2.68164 10.0698L12.7038 0.0476074Z" fill="black"/>
                        <path d="M1.9288 11.3167L4.68332 14.0712L0 16.0001L1.9288 11.3167Z" fill="black"/>
                      </svg>
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- Revision Note -->
          <div v-if="selectedLibrary?.status === 'Revisi'" class="revision-note">
            <h4>Catatan Revisi dari Admin DPK:</h4>
            <p>{{ selectedLibrary.revisionNote }}</p>
          </div>

          <!-- Revision Form -->
          <div v-if="showRevisionForm" class="revision-form">
            <h3>Identitas Perpustakaan</h3>
            
            <div class="form-group">
              <label>Periode</label>
              <select v-model="revisionData.periode">
                <option value="Semester Genap 2025/2026">Semester Genap 2025/2026</option>
                <!-- Add more options as needed -->
              </select>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label>Nomor Induk</label>
                <input type="text" v-model="revisionData.nomorInduk" placeholder="220010814562" />
              </div>
              <div class="form-group">
                <label>Nama Perpustakaan</label>
                <input type="text" v-model="revisionData.nama" placeholder="Perpustakaan Kota Yogyakarta" />
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label>Kepala Perpustakaan</label>
                <input type="text" v-model="revisionData.kepalaPerpustakaan" placeholder="Budi Sumanto" />
              </div>
              <div class="form-group">
                <label>Tahun Berdiri</label>
                <input type="text" v-model="revisionData.tahunBerdiri" placeholder="2010" />
              </div>
            </div>

            <div class="form-group">
              <label>Alamat</label>
              <input type="text" v-model="revisionData.alamat" placeholder="Jl. Malioboro No. 123, Yogyakarta" />
            </div>

            <div class="form-row">
              <div class="form-group">
                <label>Jenis Perpustakaan</label>
                <input type="text" v-model="revisionData.jenis" placeholder="Umum" />
              </div>
              <div class="form-group">
                <label>Jumlah SDM</label>
                <input type="text" v-model="revisionData.jumlahSdm" placeholder="114" :class="{ 'error': revisionData.jumlahSdm === '114' }" />
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label>Jumlah Pengunjung</label>
                <input type="text" v-model="revisionData.jumlahPengunjung" placeholder="3000" />
              </div>
              <div class="form-group">
                <label>Anggota Aktif</label>
                <input type="text" v-model="revisionData.anggotaAktif" placeholder="1500" />
              </div>
            </div>

            <div class="form-actions">
              <button class="btn-cancel" @click="cancelRevision">Batal</button>
              <button class="btn-submit" @click="submitRevision">Revisi</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useLibraryStore } from '../../store/libraryStore.js'
import { useSubmissionStore } from '../../store/submissionStore.js' // Added import for submissionStore
import Header from './Header.vue'

export default {
  name: 'ValidasiPage',
  components: {
    Header
  },
  setup() {
    const router = useRouter()
    const libraryStore = useLibraryStore()
    const submissionStore = useSubmissionStore() // Added submissionStore
    const isSidebarOpen = ref(false)
    const showRevisionForm = ref(false)
    const selectedLibrary = ref(null)
    const revisionData = ref({
      periode: '',
      nomorInduk: '',
      nama: '',
      kepalaPerpustakaan: '',
      tahunBerdiri: '',
      alamat: '',
      jenis: '',
      jumlahSdm: '',
      jumlahPengunjung: '',
      anggotaAktif: ''
    })

    // Format periode from "2024-1" to "Semester Ganjil 2024/2025"
    const formatPeriode = (periode) => {
      if (!periode) return ''
      
      const [tahun, semester] = periode.split('-')
      const tahunAkhir = parseInt(tahun) + 1
      const jenisSemester = semester === '1' ? 'Ganjil' : 'Genap'
      
      return `Semester ${jenisSemester} ${tahun}/${tahunAkhir}`
    }

    const toggleSidebar = () => {
      isSidebarOpen.value = !isSidebarOpen.value
    }

    const navigateTo = (route) => {
      router.push(`/${route}`)
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
      localStorage.removeItem('authToken')
      sessionStorage.removeItem('authToken')
      router.push('/login')
    }

    const editLibrary = (library) => {
      selectedLibrary.value = library
      if (library.status === 'Revisi') {
        showRevisionForm.value = true
        // Populate revision form with existing data
        revisionData.value = {
          periode: library.periode,
          nomorInduk: library.nomorInduk || '',
          nama: library.nama,
          kepalaPerpustakaan: library.kepalaPerpustakaan,
          tahunBerdiri: library.tahunBerdiri || '',
          alamat: library.alamat || '',
          jenis: library.jenis || '',
          jumlahSdm: library.jumlahSdm || '',
          jumlahPengunjung: library.jumlahPengunjung || '',
          anggotaAktif: library.anggotaAktif || ''
        }
      }
    }

    const cancelRevision = () => {
      showRevisionForm.value = false
      selectedLibrary.value = null
      // Reset form data
      revisionData.value = {
        periode: '',
        nomorInduk: '',
        nama: '',
        kepalaPerpustakaan: '',
        tahunBerdiri: '',
        alamat: '',
        jenis: '',
        jumlahSdm: '',
        jumlahPengunjung: '',
        anggotaAktif: ''
      }
    }

    const submitRevision = () => {
      if (selectedLibrary.value) {
        // Create updated data object
        const updatedData = {
          ...selectedLibrary.value,
          ...revisionData.value,
          status: 'Telah Direvisi'
        }
        
        // Update in library store
        libraryStore.updateLibrary(selectedLibrary.value.id, updatedData)
        
        // Also update in submission store if needed
        submissionStore.updateStatus(selectedLibrary.value.periode, 'Telah Direvisi')
        
        // Close form and reset
        showRevisionForm.value = false
        selectedLibrary.value = null
        
        // Reset form data
        revisionData.value = {
          periode: '',
          nomorInduk: '',
          nama: '',
          kepalaPerpustakaan: '',
          tahunBerdiri: '',
          alamat: '',
          jenis: '',
          jumlahSdm: '',
          jumlahPengunjung: '',
          anggotaAktif: ''
        }
      }
    }

    // Load data on component mount
    onMounted(() => {
      // Load data from store
      libraryStore.loadFromLocalStorage()

      // Set initial status for new data if not set
      const libraries = libraryStore.libraries
      libraries.forEach(library => {
        if (!library.status) {
          libraryStore.updateLibrary(library.id, {
            ...library,
            status: 'Belum Dikirim'
          })
        }
      })
    })

    return {
      isSidebarOpen,
      libraries: libraryStore.libraries,
      showRevisionForm,
      selectedLibrary,
      revisionData,
      toggleSidebar,
      navigateTo,
      navigateToNotifications,
      goToSettings,
      logout,
      editLibrary,
      cancelRevision,
      submitRevision,
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
  gap: 2rem;
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
  top: 0;
  right: 0;
  width: 8px;
  height: 8px;
  background-color: #ef4444;
  border-radius: 50%;
}

.profile-btn {
  display: flex;
  align-items: center;
  gap: 0.75rem;
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
  overflow: hidden;
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

.status-badge.revised {
  background-color: #dbeafe;
  color: #1e40af;
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

.validation-content {
  flex: 1;
  margin-left: 250px;
  padding: 2rem;
  background-color: white;
  overflow-y: auto;
  height: 100%;
}

h2 {
  font-size: 1.5rem;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 2.5rem;
}

.data-section h3 {
  font-size: 1.25rem;
  font-weight: 500;
  color: #1e293b;
  margin-bottom: 1.5rem;
}

.table-container {
  background-color: white;
  border-radius: 0.5rem;
  overflow: hidden;
  border: 1px solid #e2e8f0;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th {
  background-color: white;
  color: #1e293b;
  font-weight: 600;
  text-align: left;
  padding: 1rem;
  border-bottom: 1px solid #e2e8f0;
}

td {
  padding: 1rem;
  color: #1e293b;
  border-bottom: 1px solid #e2e8f0;
  background-color: white;
}

tr:last-child td {
  border-bottom: none;
}

.status-badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 0.25rem;
  font-size: 0.875rem;
  font-weight: 500;
}

.status-not-sent {
  background-color: #f3f4f6;
  color: #4b5563;
}

.status-waiting {
  background-color: #fef9c3;
  color: #854d0e;
}

.status-revision {
  background-color: #fecaca;
  color: #991b1b;
}

.status-valid {
  background-color: #dcfce7;
  color: #166534;
}

.status-revised {
  background-color: #dbeafe;
  color: #1e40af;
}

.edit-btn {
  background: none;
  border: none;
  padding: 0.25rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto;
  transition: all 0.2s ease;
}

.edit-btn:hover {
  opacity: 0.7;
}

.edit-btn svg {
  width: 16px;
  height: 16px;
}

.edit-btn i {
  font-size: 1rem;
}

/* Mobile responsive adjustments */
@media (max-width: 768px) {
  .validation-content {
    margin-left: 0;
    padding: 1rem;
  }

  table {
    font-size: 0.875rem;
  }

  th, td {
    padding: 0.75rem;
  }

  .status-badge {
    padding: 0.25rem 0.5rem;
  }
}

.revision-note {
  margin-top: 2rem;
  padding: 1rem;
  background-color: #fef2f2;
  border-radius: 0.5rem;
  border: 1px solid #fecaca;
}

.revision-note h4 {
  color: #991b1b;
  font-size: 1rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
}

.revision-note p {
  color: #991b1b;
  font-size: 0.875rem;
  line-height: 1.5;
}

.revision-form {
  margin-top: 2rem;
  padding: 2rem;
  background-color: white;
  border-radius: 0.5rem;
  border: 1px solid #e2e8f0;
  margin-bottom: 2rem;
}

.revision-form h3 {
  font-size: 1.25rem;
  font-weight: 500;
  color: #1e293b;
  margin-bottom: 1.5rem;
}

.form-row {
  display: flex;
  gap: 1rem;
  margin-bottom: 1rem;
}

.form-group {
  flex: 1;
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  font-size: 0.875rem;
  font-weight: 500;
  color: #1e293b;
  margin-bottom: 0.5rem;
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #e2e8f0;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  color: #1e293b;
}

.form-group input.error {
  border-color: #ef4444;
  background-color: #fef2f2;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 2rem;
}

.btn-cancel {
  padding: 0.5rem 1rem;
  background-color: white;
  border: 1px solid #e2e8f0;
  border-radius: 0.375rem;
  color: #64748b;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-submit {
  padding: 0.5rem 1rem;
  background-color: #1e40af;
  border: none;
  border-radius: 0.375rem;
  color: white;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-cancel:hover {
  background-color: #f8fafc;
}

.btn-submit:hover {
  background-color: #1e3a8a;
}
</style>