<template>
  <div class="login-container">
    <!-- Header Section -->
    <header class="header-section">
      <div class="logo-title">
        <img src="../assets/logo-sidapus.png" alt="Logo SIDAPUS" class="logo">
        <h1>Sistem Data Perpustakaan<br>Dan Kearsipan</h1>
      </div>
    </header>

    <!-- Main Content -->
    <main class="main-content">
      <div class="welcome-section">
        <h2>Hallo Selamat Datang</h2>
        <img src="../assets/logo biru.png" alt="Logo SIDAPUS" class="welcome-logo">
        <p>Gerbang menuju administrasi<br>Perpustakaan digital</p>
      </div>

      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <input
            type="email"
            id="username"
            v-model="form.username"
            required
            placeholder="Email"
            @blur="validateEmail"
            :class="{ 'error': errors.username }"
          />
          <span class="error-message" v-if="errors.username">{{ errors.username }}</span>
        </div>

        <div class="form-group">
          <input
            type="password"
            id="password"
            v-model="form.password"
            required
            placeholder="Password"
            @blur="validatePassword"
            :class="{ 'error': errors.password }"
          />
          <span class="error-message" v-if="errors.password">{{ errors.password }}</span>
        </div>

        <div class="remember-me">
          <input
            type="checkbox"
            id="remember"
            v-model="form.rememberMe"
          />
          <label for="remember">Simpan Password</label>
        </div>

        <button type="submit" class="login-button" :disabled="isSubmitting || hasErrors">
          {{ isSubmitting ? 'Memproses...' : 'Masuk' }}
        </button>
      </form>
    </main>
  </div>
</template>

<script>
export default {
  data() {
    return {
      form: {
        username: '',
        password: '',
        rememberMe: false
      },
      errors: {
        username: '',
        password: ''
      },
      isSubmitting: false
    }
  },
  computed: {
    hasErrors() {
      return Object.values(this.errors).some(error => error !== '')
    }
  },
  methods: {
    validateEmail() {
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
      if (!this.form.username) {
        this.errors.username = 'Email wajib diisi'
      } else if (!emailRegex.test(this.form.username)) {
        this.errors.username = 'Format email tidak valid'
      } else {
        this.errors.username = ''
      }
    },
    validatePassword() {
      const hasUpperCase = /[A-Z]/.test(this.form.password)
      const hasSpecialChar = /[!@#$%^&*_\-.,/]/.test(this.form.password)
      
      if (!this.form.password) {
        this.errors.password = 'Password wajib diisi'
      } else if (this.form.password.length < 8) {
        this.errors.password = 'Password minimal 8 karakter'
      } else if (!hasUpperCase) {
        this.errors.password = 'Password harus memiliki minimal 1 huruf besar'
      } else if (!hasSpecialChar) {
        this.errors.password = 'Password harus memiliki minimal 1 karakter unik (!@#$%^&*_-.,/)'
      } else {
        this.errors.password = ''
      }
    },
    async handleLogin() {
      this.validateEmail()
      this.validatePassword()

      if (this.hasErrors) return

      this.isSubmitting = true

      try {
        // Mock login - replace with actual API call
        await new Promise(resolve => setTimeout(resolve, 1000))
        
        // Simulate user role check (replace with actual role check from API)
        const userRole = this.form.username.includes('executive') ? 'executive' : 'admin_perpustakaan'
        
        if (this.form.rememberMe) {
          localStorage.setItem('authToken', 'mock-token')
          localStorage.setItem('userEmail', this.form.username)
          localStorage.setItem('userType', userRole)
        } else {
          sessionStorage.setItem('authToken', 'mock-token')
          sessionStorage.setItem('userType', userRole)
          localStorage.removeItem('userEmail')
        }

        // Redirect based on user role
        if (userRole === 'executive') {
          this.$router.push('/dashboard-executive')
        } else {
          this.$router.push('/dashboard')
        }
      } catch (error) {
        console.error('Login failed', error)
        alert('Login gagal. Silakan cek email dan password Anda.')
      } finally {
        this.isSubmitting = false
      }
    }
  },
  mounted() {
    // Check for saved email
    const savedEmail = localStorage.getItem('userEmail')
    if (savedEmail) {
      this.form.username = savedEmail
      this.form.rememberMe = true
    }
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  width: 100%;
  display: flex;
  flex-direction: column;
  background-color: white;
  margin-top: calc(-2vh);
  margin-left: calc(-50vw + 50%);
  margin-right: calc(-50vw + 50%);
  width: calc(100vw + 50px);
}

.header-section {
  background-color: #0E2954;
  color: white;
  padding: 1rem 2rem;
}

.logo-title {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.logo {
  width: 35px;
  height: auto;
}

.logo-title h1 {
  font-size: 1.1rem;
  font-weight: 500;
  line-height: 1.3;
  margin: 0;
}

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  max-width: 400px;
  margin: 0 auto;
  width: 100%;
}

.welcome-section {
  text-align: center;
  margin-bottom: 2rem;
}

.welcome-section h2 {
  font-size: 1.5rem;
  color: #0E2954;
  margin-bottom: 1.5rem;
}

.welcome-logo {
  width: 60px;
  height: auto;
  margin: 1.5rem 0;
}

.welcome-section p {
  font-size: 1rem;
  color: #0E2954;
  line-height: 1.5;
  margin-top: 1.5rem;
}

.login-form {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.form-group {
  width: 100%;
  margin-bottom: 0.5rem;
  position: relative;
}

.form-group input {
  width: 100%;
  padding: 0.75rem 1rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
  outline: none;
  transition: border-color 0.2s;
}

.form-group input.error {
  border-color: #dc3545;
}

.form-group input:focus {
  border-color: #0E2954;
}

.error-message {
  color: #dc3545;
  font-size: 0.8rem;
  margin-top: 0.25rem;
  display: block;
  position: absolute;
  left: 0;
}

.remember-me {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin: 0.5rem 0;
}

.remember-me input[type="checkbox"] {
  width: 18px;
  height: 18px;
  border: 2px solid white;
  border-radius: 4px;
  cursor: pointer;
  appearance: none;
  -webkit-appearance: none;
  -moz-appearance: none;
  background-color: white;
  position: relative;
  margin: 0;
  transition: all 0.2s ease;
}

.remember-me input[type="checkbox"]:checked {
  background-color: white;
  border-color: black;
}

.remember-me input[type="checkbox"]:checked::before {
  content: "✓";
  position: absolute;
  color: black;
  font-size: 16px;
  font-weight: 900;
  left: 2px;
  top: -2px;
  text-shadow: 
    -0.5px -0.5px 0 white,
    0.5px -0.5px 0 white,
    -0.5px 0.5px 0 white,
    0.5px 0.5px 0 white;
}

.remember-me input[type="checkbox"]:hover {
  border-color: black;
}

.remember-me label {
  font-size: 0.875rem;
  color: #475569;
  cursor: pointer;
  user-select: none;
}

.login-button {
  width: 100%;
  padding: 0.75rem;
  background-color: #0E2954;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s;
}

.login-button:hover {
  background-color: #1a3a6e;
}

.login-button:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}

@media (max-width: 768px) {
  .main-content {
    padding: 1.5rem;
  }

  .logo-title h1 {
    font-size: 1rem;
  }

  .welcome-section h2 {
    font-size: 1.3rem;
  }

  .welcome-logo {
    width: 50px;
  }

  .welcome-section p {
    font-size: 0.9rem;
  }
}
</style>