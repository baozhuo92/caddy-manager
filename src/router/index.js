import { createRouter, createWebHistory } from 'vue-router'
import Dashboard from '../views/Dashboard.vue'
import Config from '../views/Config.vue'
import Server from '../views/Server.vue'
import Certificates from '../views/Certificates.vue'
import Proxy from '../views/Proxy.vue'

const routes = [
  {
    path: '/',
    name: 'Dashboard',
    component: Dashboard
  },
  {
    path: '/conf',
    name: 'Config',
    component: Config
  },
  {
    path: '/servers',
    name: 'Server',
    component: Server
  },
  {
    path: '/certificates',
    name: 'Certificates',
    component: Certificates
  },
  {
    path: '/proxy',
    name: 'Proxy',
    component: Proxy
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router