```vue
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
            type="text"
            id="username"
            v-model="form.username"
            required
            placeholder="Email atau Username"
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

        <div class="register-link">
          <p>Belum punya akun? <router-link to="/register">Daftar sekarang</router-link></p>
        </div>
      </form>
    </main>
  </div>
</template>

<script>
import axios from 'axios'

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
  methods: {
    async handleLogin() {
    this.isSubmitting = true;

    try {
      const response = await axios.post('http://localhost:8080/login', {
        username: this.form.username,
        password: this.form.password
      }, {
        headers: {
          'Content-Type': 'application/json'
        }
      });

      // Pastikan struktur response sesuai
      const { token, user_type: userType, user } = response.data;

      if (!token || !userType) {
        throw new Error('Response tidak valid dari server');
      }

      // Simpan data login
      const storage = this.form.rememberMe ? localStorage : sessionStorage;
      storage.setItem('authToken', token);
      storage.setItem('userType', userType);
      storage.setItem('userData', JSON.stringify(user));

      // Redirect berdasarkan role
      this.redirectBasedOnRole(userType);
      
    } catch (error) {
      console.error('Login error:', error);
      
      // Error handling yang lebih baik
      let errorMsg = 'Login gagal. Silakan coba lagi.';
      
      if (error.response) {
        // Error dari server
        errorMsg = error.response.data.error || errorMsg;
      } else if (error.request) {
        // Tidak dapat terhubung ke server
        errorMsg = 'Tidak dapat terhubung ke server. Periksa koneksi Anda.';
      } else {
        // Error lain
        errorMsg = error.message || errorMsg;
      }
      
      alert(errorMsg);
    } finally {
      this.isSubmitting = false;
    }
  },

  redirectBasedOnRole(userType) {
    let route = '/';
    
    switch(userType.toLowerCase()) { // Gunakan toLowerCase untuk lebih aman
      case 'admin_perpustakaan':
        route = '/dashboard';
        break;
      case 'admin_dpk':
        route = '/admin-dpk/dashboard';
        break;
      case 'executive':
        route = '/executive/dashboard';
        break;
      default:
       console.error(`Role tidak dikenali: ${userType}`);
        return;
    }
    
    this.$router.push(route);
  },
  
  redirectBasedOnRole(userType) {
    let route = '/';
    
    switch(userType) {
      case 'admin_perpustakaan':
        route = '/dashboard';
        break;
      case 'admin_dpk':
        route = '/admin-dpk/dashboard';
        break;
      case 'executive':
        route = '/executive/dashboard';
        break;
      default:
        console.error(`Role tidak dikenali: ${userType}`);
        return;
    }
    
    this.$router.push(route);
  }
},
  mounted() {
    // Cek token yang tersimpan
    const token = localStorage.getItem('authToken') || sessionStorage.getItem('authToken')
    const userType = localStorage.getItem('userType') || sessionStorage.getItem('userType')
    
    // Jika sudah login, redirect ke dashboard sesuai role
    if (token && userType) {
      this.redirectBasedOnRole(userType)
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
}

.header-section {
  background-color: #0E2954;
  color: white;
  padding: 1rem 2rem;
  flex-shrink: 0;
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
  padding: 2rem;
  max-width: 400px;
  margin: 0 auto;
  width: 100%;
  max-height: calc(100vh - 80px);
  padding-right: 10px;
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
  border: 2px solid #333;
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
  border-color: #333;
}

.remember-me input[type="checkbox"]:checked::before {
  content: "✓";
  position: absolute;
  color: #333;
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
  border-color: #333;
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

.register-link {
  margin-top: 1rem;
  text-align: center;
}

.register-link p {
  font-size: 0.9rem;
  color: #475569;
}

.register-link a {
  color: #0E2954;
  text-decoration: none;
  font-weight: 500;
}

.register-link a:hover {
  text-decoration: underline;
}

@media (max-width: 768px) {
  .main-content {
    padding: 1.5rem;
    max-height: calc(100vh - 60px);
    padding-right: 8px;
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
```