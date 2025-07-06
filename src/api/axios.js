import axios from 'axios'

const api = axios.create({
  baseURL: 'http://localhost:8080',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Tambahkan interceptor untuk request
api.interceptors.request.use(
  config => {
    const token = localStorage.getItem('authToken') || sessionStorage.getItem('authToken')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
      console.log('[AXIOS] JWT token attached:', token)
    } else {
      console.log('[AXIOS] No JWT token found')
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// Tambahkan interceptor untuk response
api.interceptors.response.use(
  response => response,
  error => {
    if (error.response && error.response.status === 401) {
      // Handle unauthorized
      localStorage.removeItem('authToken')
      sessionStorage.removeItem('authToken')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export default api