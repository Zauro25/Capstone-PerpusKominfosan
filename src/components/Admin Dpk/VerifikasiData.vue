<template>
  <div class="dashboard-container">
    <!-- Header -->
    <header class="header">
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
          <span>Admin DPK</span>
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
          <router-link to="/dashboard-dpk" class="nav-btn">
            <span>Dashboard</span>
          </router-link>
          <router-link to="/verifikasi-dpk" class="nav-btn active">
            <span>Verifikasi Data</span>
          </router-link>
          <router-link to="/laporan" class="nav-btn">
            <span>Laporan</span>
          </router-link>
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

      <!-- Verification Content -->
      <main class="verification-content">
        <h2 class="page-title">Verifikasi Data</h2>
        
        <!-- Filter Section -->
        <div class="filter-section">
          <div class="custom-dropdown">
            <button class="dropdown-toggle" @click="toggleSemesterDropdown">
              {{ formatSelectedSemester }}
              <svg 
                xmlns="http://www.w3.org/2000/svg" 
                width="16" 
                height="16" 
                viewBox="0 0 24 24" 
                fill="none" 
                stroke="currentColor" 
                stroke-width="2" 
                stroke-linecap="round" 
                stroke-linejoin="round"
                :class="{ 'rotate': showSemesterDropdown }"
              >
                <path d="M6 9l6 6 6-6"/>
              </svg>
            </button>
            <div class="dropdown-menu" v-if="showSemesterDropdown">
              <div 
                class="dropdown-item"
                :class="{ 'active': selectedSemester === 'all' }"
                @click="selectSemester('all')"
              >
                Semua Periode
              </div>
              <div 
                class="dropdown-item"
                :class="{ 'active': selectedSemester === '2025-2' }"
                @click="selectSemester('2025-2')"
              >
                Semester Genap 2025/2026
              </div>
              <div 
                class="dropdown-item"
                :class="{ 'active': selectedSemester === '2025-1' }"
                @click="selectSemester('2025-1')"
              >
                Semester Ganjil 2025/2026
              </div>
              <div 
                class="dropdown-item"
                :class="{ 'active': selectedSemester === '2024-2' }"
                @click="selectSemester('2024-2')"
              >
                Semester Genap 2024/2025
              </div>
              <div 
                class="dropdown-item"
                :class="{ 'active': selectedSemester === '2024-1' }"
                @click="selectSemester('2024-1')"
              >
                Semester Ganjil 2024/2025
              </div>
              <div 
                class="dropdown-item"
                :class="{ 'active': selectedSemester === '2023-2' }"
                @click="selectSemester('2023-2')"
              >
                Semester Genap 2023/2024
              </div>
              <div 
                class="dropdown-item"
                :class="{ 'active': selectedSemester === '2023-1' }"
                @click="selectSemester('2023-1')"
              >
                Semester Ganjil 2023/2024
              </div>
            </div>
          </div>
          
          <div class="custom-dropdown">
            <button class="dropdown-toggle" @click="toggleStatusDropdown">
              Status Data
              <svg 
                xmlns="http://www.w3.org/2000/svg" 
                width="16" 
                height="16" 
                viewBox="0 0 24 24" 
                fill="none" 
                stroke="currentColor" 
                stroke-width="2" 
                stroke-linecap="round" 
                stroke-linejoin="round"
                :class="{ 'rotate': showStatusDropdown }"
              >
                <path d="M6 9l6 6 6-6"/>
              </svg>
            </button>
            <div class="dropdown-menu" v-if="showStatusDropdown">
              <div 
                class="dropdown-item"
                :class="{ 'active': selectedStatus === 'all' }"
                @click="selectStatus('all')"
              >
                Semua Data
              </div>
              <div 
                class="dropdown-item"
                :class="{ 'active': selectedStatus === 'Menunggu' }"
                @click="selectStatus('Menunggu')"
              >
                Menunggu
              </div>
              <div 
                class="dropdown-item"
                :class="{ 'active': selectedStatus === 'Revisi' }"
                @click="selectStatus('Revisi')"
              >
                Revisi
              </div>
              <div 
                class="dropdown-item"
                :class="{ 'active': selectedStatus === 'Valid' }"
                @click="selectStatus('Valid')"
              >
                Valid
              </div>
              <div 
                class="dropdown-item"
                :class="{ 'active': selectedStatus === 'Telah Direvisi' }"
                @click="selectStatus('Telah Direvisi')"
              >
                Telah Direvisi
              </div>
            </div>
          </div>
        </div>

        <!-- Data Table -->
        <div class="data-section">
          <div class="table-container">
            <table>
              <thead>
                <tr>
                  <th>Nomor Induk</th>
                  <th>Nama Perpustakaan</th>
                  <th>Nama Kepala</th>
                  <th>Status</th>
                  <th>Aksi</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="library in filteredLibraries" :key="library.id">
                  <td>{{ library.nomorInduk }}</td>
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
                    <button class="review-btn" @click="reviewData(library)">
                      Review
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Review Modal -->
        <div v-if="showReviewModal" class="modal">
          <div class="modal-content">
            <div class="modal-header">
              <h3>Review Data Perpustakaan</h3>
              <button class="close-btn" @click="closeReviewModal">&times;</button>
            </div>

            <div class="modal-body">
              <div class="library-details">
                <div class="detail-section">
                  <h4>Identitas Perpustakaan</h4>
                  <div class="detail-row">
                    <span class="label">Nomor Induk:</span>
                    <span class="value">{{ selectedLibrary?.nomorInduk }}</span>
                  </div>
                  <div class="detail-row">
                    <span class="label">Nama Perpustakaan:</span>
                    <span class="value">{{ selectedLibrary?.nama }}</span>
                  </div>
                  <div class="detail-row">
                    <span class="label">Kepala Perpustakaan:</span>
                    <span class="value">{{ selectedLibrary?.kepalaPerpustakaan }}</span>
                  </div>
                  <div class="detail-row">
                    <span class="label">Alamat:</span>
                    <span class="value">{{ selectedLibrary?.alamat }}</span>
                  </div>
                  <div class="detail-row">
                    <span class="label">Tahun Berdiri:</span>
                    <span class="value">{{ selectedLibrary?.tahunBerdiri }}</span>
                  </div>
                  <div class="detail-row">
                    <span class="label">Jenis Perpustakaan:</span>
                    <span class="value">{{ selectedLibrary?.jenis }}</span>
                  </div>
                </div>

                <div class="detail-section">
                  <h4>Data SDM dan Pengunjung</h4>
                  <div class="detail-row">
                    <span class="label">Jumlah SDM:</span>
                    <span class="value">{{ selectedLibrary?.jumlahSdm }}</span>
                  </div>
                  <div class="detail-row">
                    <span class="label">Jumlah Pengunjung:</span>
                    <span class="value">{{ selectedLibrary?.jumlahPengunjung }}</span>
                  </div>
                  <div class="detail-row">
                    <span class="label">Anggota Aktif:</span>
                    <span class="value">{{ selectedLibrary?.anggotaAktif }}</span>
                  </div>
                </div>

                <div class="review-form">
                  <h4>Catatan Review</h4>
                  <textarea 
                    v-model="reviewNote"
                    placeholder="Masukkan catatan review jika diperlukan..."
                    rows="4"
                  ></textarea>
                </div>
              </div>
            </div>

            <div class="modal-footer">
              <button 
                class="btn-reject" 
                @click="updateStatus('Revisi')"
                :disabled="!reviewNote && selectedStatus === 'Revisi'"
              >
                Revisi
              </button>
              <button 
                class="btn-approve"
                @click="updateStatus('Valid')"
              >
                Valid
              </button>
            </div>
          </div>
        </div>
      </main>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useLibraryStore } from '../../store/libraryStore.js'

export default {
  name: 'VerifikasiData',
  setup() {
    const router = useRouter()
    const libraryStore = useLibraryStore()
    const isSidebarOpen = ref(false)
    const hasUnreadNotifications = ref(false)
    const selectedSemester = ref('all') // Default to show all periods
    const selectedStatus = ref('all')
    const showReviewModal = ref(false)
    const selectedLibrary = ref(null)
    const reviewNote = ref('')
    const showSemesterDropdown = ref(false)
    const showStatusDropdown = ref(false)

    // Format the selected semester for display
    const formatSelectedSemester = computed(() => {
      if (selectedSemester.value === 'all') {
        return 'Semua Periode'
      }
      const [year, semester] = selectedSemester.value.split('-')
      const nextYear = parseInt(year) + 1
      return semester === '1' 
        ? `Semester Ganjil ${year}/${nextYear}`
        : `Semester Genap ${year}/${nextYear}`
    })

    const toggleSemesterDropdown = () => {
      showSemesterDropdown.value = !showSemesterDropdown.value
      if (showSemesterDropdown.value) {
        showStatusDropdown.value = false
      }
    }

    const selectSemester = (semester) => {
      selectedSemester.value = semester
      showSemesterDropdown.value = false
    }

    // Filter libraries based on selected semester and status
    const filteredLibraries = computed(() => {
      let libraries = libraryStore.libraries.filter(lib => 
        lib.status === 'Sudah Dikirim' || 
        lib.status === 'Menunggu' || 
        lib.status === 'Revisi' || 
        lib.status === 'Valid' ||
        lib.status === 'Telah Direvisi'
      )
      
      if (selectedSemester.value !== 'all') {
        libraries = libraries.filter(lib => lib.periode === selectedSemester.value)
      }
      
      if (selectedStatus.value !== 'all') {
        libraries = libraries.filter(lib => lib.status === selectedStatus.value)
      }
      
      return libraries
    })

    const toggleSidebar = () => {
      isSidebarOpen.value = !isSidebarOpen.value
    }

    const navigateToNotifications = () => {
      router.push('/notifications')
    }

    const goToSettings = () => {
      router.push('/profile-dpk')
    }

    const logout = () => {
      localStorage.removeItem('authToken')
      sessionStorage.removeItem('authToken')
      router.push('/login')
    }

    const reviewData = (library) => {
      selectedLibrary.value = library
      reviewNote.value = library.revisionNote || ''
      showReviewModal.value = true
    }

    const closeReviewModal = () => {
      showReviewModal.value = false
      selectedLibrary.value = null
      reviewNote.value = ''
    }

    const updateStatus = (newStatus) => {
      if (selectedLibrary.value) {
        libraryStore.updateLibrary(selectedLibrary.value.id, {
          ...selectedLibrary.value,
          status: newStatus,
          revisionNote: newStatus === 'Revisi' ? reviewNote.value : ''
        })
        closeReviewModal()
      }
    }

    const toggleStatusDropdown = () => {
      showStatusDropdown.value = !showStatusDropdown.value
    }

    const selectStatus = (status) => {
      selectedStatus.value = status
      showStatusDropdown.value = false
    }

    // Close dropdown when clicking outside
    const handleClickOutside = (event) => {
      const statusDropdown = document.querySelector('.custom-dropdown:nth-child(2)')
      const semesterDropdown = document.querySelector('.custom-dropdown:nth-child(1)')
      
      if (statusDropdown && !statusDropdown.contains(event.target)) {
        showStatusDropdown.value = false
      }
      if (semesterDropdown && !semesterDropdown.contains(event.target)) {
        showSemesterDropdown.value = false
      }
    }

    onMounted(() => {
      const token = localStorage.getItem('authToken') || sessionStorage.getItem('authToken')
      const userType = localStorage.getItem('userType') || sessionStorage.getItem('userType')

      if (!token || userType !== 'admin_dpk') {
        router.push('/login')
      }

      // Load library data
      libraryStore.loadFromLocalStorage()
      document.addEventListener('click', handleClickOutside)
    })

    onUnmounted(() => {
      document.removeEventListener('click', handleClickOutside)
    })

    return {
      isSidebarOpen,
      hasUnreadNotifications,
      selectedSemester,
      selectedStatus,
      showReviewModal,
      selectedLibrary,
      reviewNote,
      filteredLibraries,
      toggleSidebar,
      navigateToNotifications,
      goToSettings,
      logout,
      reviewData,
      closeReviewModal,
      updateStatus,
      showStatusDropdown,
      toggleStatusDropdown,
      selectStatus,
      showSemesterDropdown,
      toggleSemesterDropdown,
      selectSemester,
      formatSelectedSemester,
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

.main-content {
  display: flex;
  min-height: calc(100vh - 70px);
  margin-top: 70px;
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
}

.sidebar-menu {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  padding: 1rem;
  padding-top: 1rem;
  flex: 1 0 auto;
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
  text-decoration: none;
}

.nav-btn:hover {
      background-color: #0E2954;
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

.verification-content {
  flex: 1;
  margin-left: 250px;
  padding: 2rem;
  background-color: #F7F7F7;
  overflow-y: auto;
}

.page-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 2rem;
}

.filter-section {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
}

.semester-select,
.status-select {
  padding: 0.75rem;
  border: 1px solid #e2e8f0;
  border-radius: 0.5rem;
  font-size: 0.875rem;
  min-width: 200px;
  background-color: white;
}

.data-section {
  background-color: white;
  border-radius: 0.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.table-container {
  width: 100%;
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th {
  background-color: #f8fafc;
  padding: 1rem;
  text-align: left;
  font-weight: 600;
  color: #1e293b;
  border-bottom: 1px solid #e2e8f0;
}

td {
  padding: 1rem;
  border-bottom: 1px solid #e2e8f0;
  color: #1e293b;
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

.review-btn {
  padding: 0.5rem 1rem;
  background-color: #5932EA;
  color: white;
  border: none;
  border-radius: 0.25rem;
  font-size: 0.875rem;
  cursor: pointer;
  transition: background-color 0.2s;
}

.review-btn:hover {
  background-color: #4318FF;
}

/* Modal Styles */
.modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background-color: white;
  border-radius: 0.5rem;
  width: 90%;
  max-width: 800px;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-header {
  padding: 1.5rem;
  border-bottom: 1px solid #e2e8f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h3 {
  margin: 0;
  font-size: 1.25rem;
  color: #1e293b;
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #64748b;
}

.modal-body {
  padding: 1.5rem;
}

.library-details {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.detail-section {
  border: 1px solid #e2e8f0;
  border-radius: 0.5rem;
  padding: 1.5rem;
}

.detail-section h4 {
  margin: 0 0 1rem 0;
  color: #1e293b;
  font-size: 1rem;
}

.detail-row {
  display: flex;
  margin-bottom: 0.75rem;
}

.detail-row .label {
  width: 200px;
  color: #64748b;
  font-size: 0.875rem;
}

.detail-row .value {
  color: #1e293b;
  font-size: 0.875rem;
}

.review-form {
  margin-top: 1.5rem;
}

.review-form h4 {
  margin: 0 0 1rem 0;
  color: #1e293b;
  font-size: 1rem;
}

.review-form textarea {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #e2e8f0;
  border-radius: 0.5rem;
  resize: vertical;
  font-size: 0.875rem;
}

.modal-footer {
  padding: 1.5rem;
  border-top: 1px solid #e2e8f0;
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
}

.btn-reject,
.btn-approve {
  padding: 0.75rem 1.5rem;
  border-radius: 0.5rem;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-reject {
  background-color: white;
  border: 1px solid #e2e8f0;
  color: #64748b;
}

.btn-reject:hover {
  background-color: #f8fafc;
}

.btn-approve {
  background-color: #5932EA;
  border: none;
  color: white;
}

.btn-approve:hover {
  background-color: #4318FF;
}

.btn-reject:disabled,
.btn-approve:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Mobile Responsive */
@media (max-width: 768px) {
  .verification-content {
    margin-left: 0;
    padding: 1rem;
  }

  .filter-section {
    flex-direction: column;
  }

  .modal-content {
    width: 95%;
    margin: 1rem;
  }

  .detail-row {
    flex-direction: column;
  }

  .detail-row .label {
    width: 100%;
    margin-bottom: 0.25rem;
  }
}

.custom-dropdown {
  position: relative;
  width: 250px; /* Made wider to accommodate semester text */
}

.dropdown-toggle {
  width: 100%;
  padding: 0.75rem;
  background: white;
  border: 1px solid #e2e8f0;
  border-radius: 0.5rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  cursor: pointer;
  color: #1e293b;
  font-size: 0.875rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.dropdown-toggle svg {
  transition: transform 0.2s ease;
}

.dropdown-toggle svg.rotate {
  transform: rotate(180deg);
}

.dropdown-menu {
  position: absolute;
  top: calc(100% + 4px);
  left: 0;
  width: 100%;
  background: white;
  border: 1px solid #e2e8f0;
  border-radius: 0.5rem;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  z-index: 1000;
}

.dropdown-item {
  padding: 0.75rem;
  cursor: pointer;
  transition: background-color 0.2s;
  color: #1e293b;
  font-size: 0.875rem;
}

.dropdown-item:hover {
  background-color: #f8fafc;
}

.dropdown-item.active {
  background-color: #f1f5f9;
  font-weight: 500;
}
</style>