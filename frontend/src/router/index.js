import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/Home.vue'
import WalletView from '../components/Wallet.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: HomeView
  },
  {
    path: '/wallet',
    name: 'Wallet',
    component: WalletView
  },
  {
    path: '/market',
    name: 'Market',
    component: () => import('../views/Market.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
