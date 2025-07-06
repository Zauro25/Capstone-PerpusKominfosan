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

        <!-- Data Section -->
        <div class="data-section">
          <div class="content-wrapper">
          <h2>Data Perpustakaan</h2>
          <div class="table-container">
            <table class="data-table">
              <thead>
                <tr>
                  <th>Periode</th>
                  <th>Nama Perpustakaan</th>
                  <th>Nama Kepala</th>
                  <th>Tahun Berdiri</th>
                  <th>Aksi</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="library in libraries" :key="library.id">
                    <td data-label="Periode">{{ formatPeriode(library.periode) }}</td>
                    
                    <td data-label="Nama Perpustakaan">{{ library.nama_perpustakaan }}</td>
                    <td data-label="Nama Kepala">{{ library.kepala_perpustakaan }}</td>
                    <td data-label="Tahun Berdiri">{{ library.tahun_berdiri }}</td>
                  <td>
                    <div class="action-menu">
                        <button class="action-toggle" @click="toggleMenu($event, library.id)">•••</button>
                        <div v-if="activeMenu === library.id" 
                             class="dropdown-menu"
                             :style="{ top: menuPosition.top, left: menuPosition.left }">
                        <button @click="viewDetails(library)">
                          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                          </svg>
                            <span>Lihat Detail</span>
                        </button>
                        <button @click="editLibrary(library)">
                          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                          </svg>
                            <span>Edit</span>
                        </button>
                        <button class="delete" @click="deleteLibrary(library)">
                          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                          </svg>
                            <span>Hapus</span>
                        </button>
                      </div>
                    </div>
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
  name: 'DaftarData',
  setup() {
    const router = useRouter()
    const libraryStore = useLibraryStore()
    const isSidebarOpen = ref(false)
    const activeMenu = ref(null)
    const menuPosition = ref({ top: '0px', left: '0px' })
    const hasUnreadNotifications = ref(true)
    const isMobile = ref(false)
    const libraries = ref([])

    // Fetch data perpustakaan dari backend
    const fetchLibraries = async () => {
      try {
        // Ambil langsung dari store agar reactive
        await libraryStore.fetchLibraries()
        libraries.value = libraryStore.libraries
      } catch (error) {
        console.error('Gagal mengambil data perpustakaan:', error)
      }
    }

    // Format periode from "2024-1" to "Semester Ganjil 2024/2025"
    const formatPeriode = (periode) => {
      if (!periode) return ''
      
      const [tahun, semester] = periode.split('-')
      const tahunAkhir = parseInt(tahun) + 1
      const jenisSemester = semester === '1' ? 'Ganjil' : 'Genap'
      
      return `Semester ${jenisSemester} ${tahun}/${tahunAkhir}`
    }

    const isLastRow = (id) => {
      const index = libraryStore.libraries.findIndex(lib => lib.id === id);
      return index === libraryStore.libraries.length - 1;
    }

    const calculateMenuPosition = (event, id) => {
      const button = event.currentTarget;
      const rect = button.getBoundingClientRect();
      const spaceBelow = window.innerHeight - rect.bottom;
      const spaceAbove = rect.top;
      const menuHeight = 144; // Approximate height of menu

      let top;
      if (spaceBelow >= menuHeight || spaceBelow >= spaceAbove) {
        top = rect.bottom + 5;
      } else {
        top = rect.top - menuHeight - 5;
      }

      menuPosition.value = {
        top: `${top}px`,
        left: `${rect.left}px`
      };
    };

    const toggleMenu = (event, id) => {
      if (activeMenu.value === id) {
        activeMenu.value = null;
      } else {
        calculateMenuPosition(event, id);
        activeMenu.value = id;
      }
    };

    const handleClickOutside = (event) => {
      const clickedElement = event.target;
      const isToggleButton = clickedElement.closest('.action-toggle');
      const isDropdownMenu = clickedElement.closest('.dropdown-menu');
      
      if (!isToggleButton && !isDropdownMenu) {
        activeMenu.value = null;
      }
    };

    onMounted(() => {
      fetchLibraries()
      checkMobile()
      window.addEventListener('resize', checkMobile)
      document.addEventListener('click', handleClickOutside)
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
      router.push('/login')
    }

    const viewDetails = (library) => {
      libraryStore.setCurrentLibrary(library)
      router.push(`/detail/${library.id}`)
      activeMenu.value = null
    }

    const editLibrary = (library) => {
      libraryStore.setCurrentLibrary(library)
      router.push(`/input-update?edit=${library.id}`)
      activeMenu.value = null
    }

    const deleteLibrary = async (library) => {
      if (confirm('Apakah Anda yakin ingin menghapus data ini?')) {
        try {
          await libraryStore.deleteLibrary(library.id)
          await fetchLibraries() // refresh data setelah hapus
          activeMenu.value = null
        } catch (error) {
          console.error('Error deleting library:', error)
        }
      }
    }

    onUnmounted(() => {
      window.removeEventListener('resize', checkMobile)
      document.removeEventListener('click', handleClickOutside)
    })

    return {
      isSidebarOpen,
      hasUnreadNotifications,
      libraries,
      activeMenu,
      menuPosition,
      toggleSidebar,
      navigateTo,
      navigateToNotifications,
      goToSettings,
      logout,
      toggleMenu,
      viewDetails,
      editLibrary,
      deleteLibrary,
      formatPeriode,
      isLastRow,
      handleClickOutside
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

/* Existing styles for data section */
.data-section {
  padding: 2rem;
}

.content-wrapper {
  max-width: 1000px;
  margin: 0 auto;
  background: white;
  padding: 2rem;
  border-radius: 12px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

h2 {
  margin-bottom: 1.5rem;
  color: #1e293b;
  font-size: 1.5rem;
}

.table-container {
  width: 100%;
  overflow-x: auto;
  background: white;
  border-radius: 12px;
  padding: 1px;
}

.data-table {
  width: 100%;
  min-width: 800px;
  border-collapse: separate;
  border-spacing: 0;
  border: 2px solid black;
  border-radius: 12px;
  overflow: hidden;
}

.data-table th,
.data-table td {
  padding: 1rem;
  text-align: left;
  border-right: 1px solid black;
  border-bottom: 1px solid black;
}

.data-table th {
  background-color: white;
  font-weight: 600;
  color: #1e293b;
}

.data-table td {
  color: #475569;
}

.data-table tr:last-child td {
  border-bottom: none;
}

.data-table td:last-child,
.data-table th:last-child {
  border-right: none;
}

/* Column widths */
.data-table th:nth-child(1),
.data-table td:nth-child(1) {
  width: 20%; /* Periode */
}

.data-table th:nth-child(2),
.data-table td:nth-child(2) {
  width: 30%; /* Nama Perpustakaan */
}

.data-table th:nth-child(3),
.data-table td:nth-child(3) {
  width: 25%; /* Nama Kepala */
}

.data-table th:nth-child(4),
.data-table td:nth-child(4) {
  width: 15%; /* Tahun Berdiri */
}

.data-table th:nth-child(5),
.data-table td:nth-child(5) {
  width: 10%; /* Aksi */
  text-align: center;
}

/* For mobile responsiveness */
@media screen and (max-width: 768px) {
  .data-table {
    display: block;
    min-width: unset;
    border-width: 2px;
  }

  .data-table thead {
    display: none;
  }

  .data-table tbody tr {
    display: block;
    margin-bottom: 1rem;
    border: 2px solid black;
    border-radius: 12px;
    overflow: hidden;
  }

  .data-table td {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.75rem 1rem;
    border-right: none;
    border-bottom: 1px solid black;
  }

  .data-table td:last-child {
    border-bottom: none;
  }

  .data-table td::before {
    content: attr(data-label);
    font-weight: 600;
    color: #1e293b;
  }
}

.action-menu {
  position: relative;
  display: inline-block;
}

.action-toggle {
  background: none;
  border: none;
  font-size: 1.5rem;
  font-weight: 600;
  cursor: pointer;
  padding: 0.25rem 0.5rem;
  color: #64748b;
  border-radius: 4px;
}

.action-toggle:hover {
  background-color: #f1f5f9;
}

.dropdown-menu {
  position: fixed;
  background: white;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  z-index: 999;
  min-width: 160px;
  border: 1px solid #e2e8f0;
}

.dropdown-menu button {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  width: 100%;
  padding: 0.75rem 1rem;
  border: none;
  background: none;
  cursor: pointer;
  color: #475569;
  font-size: 0.875rem;
  font-weight: 600;
  transition: all 0.2s;
  white-space: nowrap;
}

.dropdown-menu button:hover {
  background-color: #f8fafc;
}

.dropdown-menu button.delete {
  color: #dc2626;
}

.dropdown-menu button.delete:hover {
  background-color: #fee2e2;
}

.dropdown-menu svg {
  width: 1.25rem;
  height: 1.25rem;
}

/* Tab Navigation Styles */
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

/* Mobile Responsive for Data Table */
@media (max-width: 768px) {
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

  .data-section {
    padding: 0;
  }

  .data-section h2 {
    margin: 1rem;
    font-size: 1.1rem;
  }

  .table-container {
    margin: 0 1rem;
    border: none;
    border-radius: 0;
  }

  .data-table {
    display: block;
  }

  .data-table thead {
    display: none;
  }

  .data-table tbody tr {
    display: block;
    padding: 1rem;
    border: 1px solid #000;
    margin-bottom: 1rem;
    border-radius: 8px;
  }

  .data-table td {
    display: flex;
    padding: 0.5rem 0;
    border: none;
    align-items: center;
  }

  .data-table td::before {
    content: attr(data-label);
    font-weight: 600;
    width: 120px;
    min-width: 120px;
    color: #1f2937;
  }

  .action-menu {
    display: flex;
    justify-content: flex-end;
    width: 100%;
    padding-top: 0.5rem;
  }

  .action-toggle {
    padding: 0.75rem;
    font-size: 1rem;
  }
}
</style>
