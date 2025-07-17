import { createRouter, createWebHistory } from 'vue-router'
import LoginPage from '../components/LoginPage.vue'
import DashboardPage from '../components/Admin perpustakaan/DashboardPage.vue'
import InputUpdatePage from '../components/Admin perpustakaan/input&updatePage.vue'
import DaftarData from '../components/Admin perpustakaan/DaftarData.vue'
import PengirimanDataPage from '../components/Admin perpustakaan/pengirimanDataPage.vue'
import ValidasiPage from '../components/Admin perpustakaan/validasi.vue'
import NotificationPage from '../components/Admin perpustakaan/NotificationPage.vue'
import LandingPage from '../components/LandingPage.vue'
import Profile from '../components/Admin perpustakaan/profile.vue'
import ProfileEdit from '../components/Admin perpustakaan/profileedit.vue'
import DashboardExecutive from '../components/Executive/DashboardExecutive.vue'
import ProfileExecutive from '../components/Executive/profileExecutive.vue'

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
    path: '/dashboard-executive',
    name: 'dashboard-executive',
    component: DashboardExecutive,
    meta: { requiresAuth: true }
  },
  {
    path: '/profile-executive',
    name: 'profile-executive',
    component: ProfileExecutive,
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
  },
  {
    path: '/profile',
    name: 'profile',
    component: Profile
  },
  {
    path: '/profile/edit',
    name: 'profile-edit',
    component: ProfileEdit
  },
  {
    path: '/detail/:id',
    name: 'detail',
    component: () => import('../components/Admin perpustakaan/DetailPage.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/input-update/:id',
    name: 'edit-library',
    component: () => import('../components/Admin perpustakaan/input&updatePage.vue'),
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