import { createRouter, createWebHashHistory } from 'vue-router'
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
  },
  {
    path: '/geniebot',
    name: 'GenieBot',
    component: () => import('../views/GenieBot.vue')
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
