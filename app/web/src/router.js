import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    component: () => import('./views/HomeView.vue')
  },
  {
    path: '/repositories',
    component: () => import('./views/RepositoriesView.vue')
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
    path: '/python',
    meta: { technology: 'python' },
    component: () => import('./views/ContributionsView.vue')
  },
  {
    path: '/impressum',
    component: () => import('./views/ImpressumView.vue')
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
