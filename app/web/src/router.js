import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    component: () => import('./views/HomeView.vue')
  },
  {
    path: '/news',
    component: () => import('./views/NewsView.vue')
  },
  {
    path: '/go',
    meta: { technology: 'go' },
    component: () => import('./views/ContributionsView.vue')
  },
  {
    path: '/node',
    meta: { technology: 'node' },
    component: () => import('./views/ContributionsView.vue')
  },
  {
    path: '/impressum',
    component: () => import('./views/ImpressumView.vue')
  },
  {
    path: '/privacy',
    component: () => import('./views/PrivacyView.vue')
  },
  {
    path: '/:pathMatch(.*)*',
    component: () => import('./views/404View.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
