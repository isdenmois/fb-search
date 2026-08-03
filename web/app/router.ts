import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/', component: () => import('@/pages/home-page/HomePage.vue') },
  {
    path: '/admin',
    component: () => import('@/pages/admin-page/AdminPage.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
