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
import DashboardDpk from '../components/Admin Dpk/DashboardDpk.vue'
import ProfileDpk from '../components/Admin Dpk/profileDpk.vue'

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
  // Admin Perpustakaan Routes
  {
    path: '/dashboard',
    name: 'dashboard',
    component: DashboardPage,
    meta: { 
      requiresAuth: true,
      allowedRoles: ['admin_perpustakaan']
    }
  },
  {
    path: '/input-update',
    name: 'input-update',
    component: InputUpdatePage,
    meta: { 
      requiresAuth: true,
      allowedRoles: ['admin_perpustakaan']
    }
  },
  {
    path: '/daftar-data-update',
    name: 'daftar-data-update',
    component: DaftarData,
    meta: { 
      requiresAuth: true,
      allowedRoles: ['admin_perpustakaan']
    }
  },
  {
    path: '/pengiriman',
    name: 'pengiriman',
    component: PengirimanDataPage,
    meta: { 
      requiresAuth: true,
      allowedRoles: ['admin_perpustakaan']
    }
  },
  {
    path: '/validasi',
    name: 'validasi',
    component: ValidasiPage,
    meta: { 
      requiresAuth: true,
      allowedRoles: ['admin_perpustakaan']
    }
  },
  {
    path: '/notifications',
    name: 'notifications',
    component: NotificationPage,
    meta: { 
      requiresAuth: true,
      allowedRoles: ['admin_perpustakaan']
    }
  },
  {
    path: '/profile',
    name: 'profile',
    component: Profile,
    meta: { 
      requiresAuth: true,
      allowedRoles: ['admin_perpustakaan']
    }
  },
  {
    path: '/profile/edit',
    name: 'profile-edit',
    component: ProfileEdit,
    meta: { 
      requiresAuth: true,
      allowedRoles: ['admin_perpustakaan']
    }
  },
  {
    path: '/detail/:id',
    name: 'detail',
    component: () => import('../components/Admin perpustakaan/DetailPage.vue'),
    meta: { 
      requiresAuth: true,
      allowedRoles: ['admin_perpustakaan']
    }
  },
  {
    path: '/input-update/:id',
    name: 'edit-library',
    component: () => import('../components/Admin perpustakaan/input&updatePage.vue'),
    meta: { 
      requiresAuth: true,
      allowedRoles: ['admin_perpustakaan']
    }
  },
  // Executive Routes
  {
    path: '/dashboard-executive',
    name: 'dashboard-executive',
    component: DashboardExecutive,
    meta: { 
      requiresAuth: true,
      allowedRoles: ['executive']
    }
  },
  {
    path: '/profile-executive',
    name: 'profile-executive',
    component: ProfileExecutive,
    meta: { 
      requiresAuth: true,
      allowedRoles: ['executive']
    }
  },
  // Admin DPK Routes
  {
    path: '/dashboard-dpk',
    name: 'dashboard-dpk',
    component: DashboardDpk,
    meta: { 
      requiresAuth: true,
      allowedRoles: ['admin_dpk']
    }
  },
  {
    path: '/profile-dpk',
    name: 'profile-dpk',
    component: ProfileDpk,
    meta: { 
      requiresAuth: true,
      allowedRoles: ['admin_dpk']
    }
  },
  {
    path: '/verifikasi-data',
    name: 'verifikasi-data',
    component: () => import('../components/Admin Dpk/VerifikasiData.vue'),
    meta: { 
      requiresAuth: true,
      allowedRoles: ['admin_dpk']
    }
  },
  {
    path: '/laporan',
    name: 'laporan',
    component: () => import('../components/Admin Dpk/Laporan.vue'),
    meta: { 
      requiresAuth: true,
      allowedRoles: ['admin_dpk']
    }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const isAuthenticated = localStorage.getItem('authToken') || sessionStorage.getItem('authToken')
  const userRole = localStorage.getItem('userType') || sessionStorage.getItem('userType')
  
  // Handle public routes
  if (!to.meta.requiresAuth) {
    if (isAuthenticated) {
      // Redirect authenticated users from login/landing to their respective dashboards
      if (to.path === '/login' || to.path === '/') {
        switch (userRole) {
          case 'admin_perpustakaan':
            next('/dashboard')
            break
          case 'executive':
            next('/dashboard-executive')
            break
          case 'admin_dpk':
            next('/dashboard-dpk')
            break
          default:
            next('/login')
        }
      } else {
        next()
      }
    } else {
      next()
    }
    return
  }

  // Handle protected routes
  if (!isAuthenticated) {
    next('/login')
    return
  }

  // Check role-based access
  const allowedRoles = to.meta.allowedRoles || []
  if (allowedRoles.length > 0 && !allowedRoles.includes(userRole)) {
    // Redirect to appropriate dashboard based on role
    switch (userRole) {
      case 'admin_perpustakaan':
        next('/dashboard')
        break
      case 'executive':
        next('/dashboard-executive')
        break
      case 'admin_dpk':
        next('/dashboard-dpk')
        break
      default:
        next('/login')
    }
    return
  }

  next()
})

export default router