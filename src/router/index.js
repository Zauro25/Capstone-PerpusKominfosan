import { createRouter, createWebHistory } from 'vue-router'
import LoginPage from '../components/LoginPage.vue'
import DashboardPage from '../components/DashboardPage.vue'
import InputUpdatePage from '../components/input&updatePage.vue'
import DaftarData from '../components/DaftarData.vue'
import NotificationPage from '../components/NotificationPage.vue'
import DetailPage from '../components/DetailPage.vue'
import PengirimanDataPage from '../components/pengirimanDataPage.vue'

const routes = [
  { 
    path: '/', 
    redirect: () => {
      const lastRoute = localStorage.getItem('lastRoute');
      return lastRoute ? `/${lastRoute}` : '/login';
    }
  },
  { 
    path: '/login', 
    component: LoginPage,
    meta: { requiresAuth: false } 
  },
  { 
    path: '/dashboard', 
    component: DashboardPage,
    meta: { requiresAuth: true } 
  },
  { 
    path: '/input-update', 
    component: InputUpdatePage,
    meta: { requiresAuth: true } 
  },
  { 
    path: '/daftar-data-update', 
    component: DaftarData,
    meta: { requiresAuth: true } 
  },
  { 
    path: '/notifications', 
    component: NotificationPage,
    meta: { requiresAuth: true } 
  },
  { 
    path: '/detail/:id', 
    component: DetailPage,
    meta: { requiresAuth: true } 
  },
  { 
    path: '/pengiriman', 
    component: PengirimanDataPage,
    meta: { requiresAuth: true } 
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Remove the initial load clear since we want to persist the state
// const isInitialLoad = true;
// if (isInitialLoad) {
//   localStorage.removeItem('authToken');
//   sessionStorage.removeItem('authToken');
// }

router.beforeEach((to, from, next) => {
  const isAuthenticated = localStorage.getItem('authToken') || sessionStorage.getItem('authToken')

  if (to.meta.requiresAuth && !isAuthenticated) {
    localStorage.removeItem('lastRoute'); // Clear last route if not authenticated
    next('/login')
  } else {
    if (to.path !== '/login') {
      localStorage.setItem('lastRoute', to.path.substring(1));
    }
    next()
  }
})

export default router