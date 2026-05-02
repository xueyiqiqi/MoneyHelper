import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { public: true },
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
    meta: { public: true },
  },
  {
    path: '/',
    component: () => import('@/components/AppLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        redirect: '/bills',
      },
      {
        path: 'spaces',
        name: 'SpaceList',
        component: () => import('@/views/SpaceList.vue'),
      },
      {
        path: 'spaces/:id',
        name: 'SpaceDetail',
        component: () => import('@/views/SpaceDetail.vue'),
      },
      {
        path: 'bills',
        name: 'BillList',
        component: () => import('@/views/BillList.vue'),
      },
      {
        path: 'bills/new',
        name: 'BillCreate',
        component: () => import('@/views/BillEdit.vue'),
      },
      {
        path: 'bills/:id',
        name: 'BillEdit',
        component: () => import('@/views/BillEdit.vue'),
      },
      {
        path: 'analysis',
        name: 'Analysis',
        component: () => import('@/views/Analysis.vue'),
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/Profile.vue'),
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    next('/login')
  } else if (to.meta.public && authStore.isLoggedIn) {
    next('/spaces')
  } else {
    next()
  }
})

export default router

