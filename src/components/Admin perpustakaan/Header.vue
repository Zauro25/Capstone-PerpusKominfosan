<template>
  <header class="header">
    <button 
      class="hamburger-menu"
      @click="$emit('toggle-sidebar')"
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
      <div class="notification-btn">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M18 8C18 6.4087 17.3679 4.88258 16.2426 3.75736C15.1174 2.63214 13.5913 2 12 2C10.4087 2 8.88258 2.63214 7.75736 3.75736C6.63214 4.88258 6 6.4087 6 8C6 15 3 17 3 17H21C21 17 18 15 18 8Z" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          <path d="M13.73 21C13.5542 21.3031 13.3019 21.5547 12.9982 21.7295C12.6946 21.9044 12.3504 21.9965 12 21.9965C11.6496 21.9965 11.3054 21.9044 11.0018 21.7295C10.6982 21.5547 10.4458 21.3031 10.27 21" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        <div class="notification-dot"></div>
      </div>
      <div class="admin-profile">
        <span>Admin Perpustakaan</span>
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M20 21V19C20 17.9391 19.5786 16.9217 18.8284 16.1716C18.0783 15.4214 17.0609 15 16 15H8C6.93913 15 5.92172 15.4214 5.17157 16.1716C4.42143 16.9217 4 17.9391 4 19V21" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          <path d="M12 11C14.2091 11 16 9.20914 16 7C16 4.79086 14.2091 3 12 3C9.79086 3 8 4.79086 8 7C8 9.20914 9.79086 11 12 11Z" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
      </div>
    </div>
  </header>
</template>

<script>
export default {
  name: 'Header',
  props: {
    isSidebarOpen: {
      type: Boolean,
      required: true
    }
  },
  emits: ['toggle-sidebar']
}
</script>

<style scoped>
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

.admin-profile {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  cursor: pointer;
  transition: opacity 0.2s ease;
}

.admin-profile:hover {
  opacity: 0.8;
}

.admin-profile span {
  font-size: 0.95rem;
  font-weight: 500;
  color: white;
}

@media (max-width: 768px) {
  .hamburger-menu {
    display: block;
  }

  .header {
    padding: 0.75rem 1rem;
  }

  .header-left h1 {
    font-size: 0.9rem;
  }

  .admin-profile span {
    display: none;
  }

  .header-right {
    gap: 1rem;
  }
}
</style> 