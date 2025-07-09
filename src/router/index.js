import { createRouter, createWebHistory } from 'vue-router'
import LoginPage from '../components/LoginPage.vue'
import DashboardPage from '../components/DashboardPage.vue'
import InputUpdatePage from '../components/input&updatePage.vue'
import DaftarData from '../components/DaftarData.vue'
import PengirimanDataPage from '../components/pengirimanDataPage.vue'
import ValidasiPage from '../components/validasi.vue'
import NotificationPage from '../components/NotificationPage.vue'
import LandingPage from '../components/LandingPage.vue'

const routes = [
  {
    path: '/',
    name: 'landing',
    component: LandingPage
  },
  {
    path: '/login',
    name: 'login',
    component: LoginPage
  },
  {
    path: '/dashboard',
    name: 'dashboard',
    component: DashboardPage,
    meta: { requiresAuth: true }
  },
  {
    path: '/input-update',
    name: 'input-update',
    component: InputUpdatePage,
    meta: { requiresAuth: true }
  },
  {
    path: '/daftar-data-update',
    name: 'daftar-data-update',
    component: DaftarData,
    meta: { requiresAuth: true }
  },
  {
    path: '/pengiriman',
    name: 'pengiriman',
    component: PengirimanDataPage,
    meta: { requiresAuth: true }
  },
  {
    path: '/validasi',
    name: 'validasi',
    component: ValidasiPage,
    meta: { requiresAuth: true }
  },
  {
    path: '/notifications',
    name: 'notifications',
    component: NotificationPage,
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const isAuthenticated = localStorage.getItem('authToken') || sessionStorage.getItem('authToken')
  
  if (to.matched.some(record => record.meta.requiresAuth)) {
    if (!isAuthenticated) {
      next('/login')
    } else {
      next()
    }
  } else {
    if (isAuthenticated && (to.path === '/login' || to.path === '/')) {
      next('/dashboard')
    } else {
      next()
    }
  }
})

export default router