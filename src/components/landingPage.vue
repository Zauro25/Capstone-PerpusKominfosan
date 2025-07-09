<template>
  <div class="landing-container">
    <header class="landing-header">
      <div class="header-left">
        <div class="logo-container">
          <img src="../assets/logo biru.png" alt="Logo" class="logo" />
          <span class="logo-text">SIDAPUS</span>
        </div>
        <nav class="nav-links">
          <a href="#beranda" class="nav-link">Beranda</a>
          <a href="#fitur" class="nav-link">Fitur</a>
          <a href="#tentang" class="nav-link">Tentang</a>
          <a href="#kontak" class="nav-link">Kontak</a>
        </nav>
      </div>
      <button class="login-btn" @click="navigateToLogin">Masuk</button>
    </header>

    <main class="landing-main">
      <section class="hero-section">
        <h1>Sistem Pencatatan dan Kearsipan<br>Perpustakaan Digital</h1>
        <p class="subtitle">Solusi terintegrasi untuk mengelola dan memantau seluruh data perpustakaan di Dinas Perpustakaan dengan cepat dan aman</p>
        <div class="cta-buttons">
          <button class="cta-btn secondary">Jelajahi Fitur</button>
          <button class="cta-btn primary" @click="navigateToLogin">Mulai Sekarang</button>
        </div>
      </section>

      <section class="stats-section">
        <div class="stat-item">
          <h2>150+</h2>
          <p>Perpustakaan Terdaftar</p>
        </div>
        <div class="stat-item">
          <h2>99.9%</h2>
          <p>Uptime Sistem</p>
        </div>
        <div class="stat-item">
          <h2>24/7</h2>
          <p>Akses Real Time</p>
        </div>
        <div class="stat-item">
          <h2>100%</h2>
          <p>Data Tervalidasi</p>
        </div>
      </section>

      <section class="features-section">
        <h2>Fitur Unggulan</h2>
        <div class="features-grid">
          <div class="feature-card">
            <div class="feature-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="3" width="18" height="18" rx="2" />
                <path d="M3 9h18" />
                <path d="M9 21V9" />
              </svg>
            </div>
            <h3>Dashboard Interaktif</h3>
            <p>Visualisasi data yang informatif untuk monitoring dan analisis data perpustakaan</p>
          </div>
          <div class="feature-card">
            <div class="feature-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14" />
                <polyline points="22 4 12 14.01 9 11.01" />
              </svg>
            </div>
            <h3>Verifikasi Data</h3>
            <p>Proses verifikasi data yang terstruktur untuk memastikan akurasi informasi</p>
          </div>
          <div class="feature-card">
            <div class="feature-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
                <path d="M7 11V7a5 5 0 0 1 10 0v4" />
              </svg>
            </div>
            <h3>Keamanan Tinggi</h3>
            <p>Sistem keamanan berlapis untuk melindungi data sensitif perpustakaan</p>
          </div>
          <div class="feature-card">
            <div class="feature-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z" />
                <polyline points="3.27 6.96 12 12.01 20.73 6.96" />
                <line x1="12" y1="22.08" x2="12" y2="12" />
              </svg>
            </div>
            <h3>Responsive Design</h3>
            <p>Akses dan kelola perpustakaan - responsive di semua perangkat</p>
          </div>
          <div class="feature-card">
            <div class="feature-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z" />
              </svg>
            </div>
            <h3>Sinkronisasi Otomatis</h3>
            <p>Data tersinkronisasi secara otomatis untuk memastikan informasi terkini</p>
          </div>
          <div class="feature-card">
            <div class="feature-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="20" x2="18" y2="10" />
                <line x1="12" y1="20" x2="12" y2="4" />
                <line x1="6" y1="20" x2="6" y2="14" />
              </svg>
            </div>
            <h3>Laporan Analitik</h3>
            <p>Laporan lengkap penggunaan sistem untuk pengambilan keputusan</p>
          </div>
        </div>
      </section>

      <section class="cta-section">
        <h2>Siap Bergabung?</h2>
        <p>Transformasi digital perpustakaan Yogyakarta dimulai dari sini. Mari bersama mewujudkan ekosistem literasi yang lebih baik.</p>
        <button class="cta-btn primary" @click="navigateToLogin">Mulai Sekarang</button>
      </section>
    </main>

    <footer class="landing-footer">
      <div class="footer-content">
        <div class="footer-section">
          <h4>SIDAPUS</h4>
          <p>Sistem Informasi Pencatatan Perpustakaan Digital untuk memudahkan pengelolaan digital perpustakaan di Dinas Yogyakarta</p>
        </div>
        <div class="footer-section">
          <h4>Kontak</h4>
          <p>Dinas Perpustakaan dan Kearsipan<br>
             Kota Yogyakarta<br>
             Jalan Suroto No. 9, Kotabaru<br>
             Telepon (0274) 511 314</p>
        </div>
      </div>
      <div class="footer-bottom">
        <p>© 2024 Dinas Perpustakaan dan Kearsipan Kota Yogyakarta. All rights reserved.</p>
        <p>Sistem Informasi Pencatatan Perpustakaan Digital</p>
      </div>
    </footer>
  </div>
</template>

<script>
import { useRouter } from 'vue-router'

export default {
  name: 'LandingPage',
  setup() {
    const router = useRouter()

    const navigateToLogin = () => {
      router.push('/login')
    }

    return {
      navigateToLogin
    }
  }
}
</script>

<style scoped>
.landing-container {
  min-height: 100vh;
  background-color: #0E2954;
  color: white;
}

.landing-header {
  padding: 1rem 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  background-color: white;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  z-index: 1000;
  height: 70px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 3rem;
}

.logo-container {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.logo {
  height: 35px;
  width: auto;
}

.logo-text {
  font-size: 1.25rem;
  font-weight: 600;
  color: #0E2954;
}

.nav-links {
  display: flex;
  gap: 2rem;
  align-items: center;
}

.nav-link {
  color: #1f2937;
  text-decoration: none;
  font-weight: 500;
  transition: color 0.2s ease;
}

.nav-link:hover {
  color: #0E2954;
}

.login-btn {
  padding: 0.5rem 2rem;
  border: none;
  border-radius: 8px;
  background-color: #0E2954;
  color: white;
  font-weight: 500;
  transition: all 0.2s ease;
}

.login-btn:hover {
  background-color: #0a1f3d;
  transform: translateY(-1px);
}

.landing-main {
  padding-top: 80px;
}

.hero-section {
  text-align: center;
  padding: 4rem 2rem;
  max-width: 800px;
  margin: 0 auto;
}

.hero-section h1 {
  font-size: 2.5rem;
  font-weight: 700;
  margin-bottom: 1.5rem;
  line-height: 1.2;
}

.subtitle {
  font-size: 1.1rem;
  color: rgba(255, 255, 255, 0.9);
  margin-bottom: 2rem;
}

.cta-buttons {
  display: flex;
  gap: 1rem;
  justify-content: center;
}

.cta-btn {
  padding: 0.75rem 2rem;
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.2s ease;
}

.cta-btn.primary {
  background-color: white;
  color: #0E2954;
}

.cta-btn.primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(255, 255, 255, 0.2);
}

.cta-btn.secondary {
  border: 2px solid white;
  color: white;
}

.cta-btn.secondary:hover {
  background-color: rgba(255, 255, 255, 0.1);
}

.stats-section {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 2rem;
  padding: 4rem 2rem;
  background-color: rgba(255, 255, 255, 0.05);
  text-align: center;
}

.stat-item h2 {
  font-size: 2.5rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
  background: linear-gradient(45deg, #4318FF, #9747FF);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.stat-item p {
  color: rgba(255, 255, 255, 0.9);
}

.features-section {
  padding: 4rem 2rem;
  text-align: center;
}

.features-section h2 {
  font-size: 2rem;
  margin-bottom: 3rem;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

.feature-card {
  background-color: rgba(255, 255, 255, 0.05);
  padding: 2rem;
  border-radius: 12px;
  transition: all 0.3s ease;
}

.feature-card:hover {
  transform: translateY(-5px);
  background-color: rgba(255, 255, 255, 0.1);
}

.feature-icon {
  width: 48px;
  height: 48px;
  margin: 0 auto 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  color: white;
}

.feature-card h3 {
  margin-bottom: 1rem;
  font-size: 1.25rem;
}

.feature-card p {
  color: rgba(255, 255, 255, 0.8);
  font-size: 0.95rem;
  line-height: 1.5;
}

.cta-section {
  text-align: center;
  padding: 4rem 2rem;
  background-color: rgba(255, 255, 255, 0.05);
}

.cta-section h2 {
  font-size: 2rem;
  margin-bottom: 1rem;
}

.cta-section p {
  color: rgba(255, 255, 255, 0.9);
  max-width: 600px;
  margin: 0 auto 2rem;
}

.landing-footer {
  background-color: rgba(0, 0, 0, 0.2);
  padding: 4rem 2rem 2rem;
}

.footer-content {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 4rem;
  max-width: 1200px;
  margin: 0 auto;
  margin-bottom: 3rem;
}

.footer-section h4 {
  font-size: 1.25rem;
  margin-bottom: 1rem;
}

.footer-section p {
  color: rgba(255, 255, 255, 0.8);
  line-height: 1.6;
}

.footer-bottom {
  text-align: center;
  padding-top: 2rem;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.footer-bottom p {
  color: rgba(255, 255, 255, 0.6);
  font-size: 0.9rem;
  margin-bottom: 0.5rem;
}

@media (max-width: 1024px) {
  .features-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .hero-section h1 {
    font-size: 2rem;
  }

  .stats-section {
    grid-template-columns: repeat(2, 1fr);
  }

  .features-grid {
    grid-template-columns: 1fr;
  }

  .footer-content {
    grid-template-columns: 1fr;
    gap: 2rem;
  }

  .nav-links {
    display: none;
  }

  .header-left {
    gap: 1rem;
  }
}

@media (max-width: 480px) {
  .cta-buttons {
    flex-direction: column;
  }

  .stats-section {
    grid-template-columns: 1fr;
  }
}
</style>
    grid-template-columns: 1fr;    flex-direction: column;
    gap: 1rem;}

/* Add smooth scrolling */
html {
  scroll-behavior: smooth;
