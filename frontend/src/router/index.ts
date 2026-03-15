import { createRouter, createWebHashHistory } from 'vue-router'
import HostLookup from '../views/HostLookup.vue'
import Search from '../views/Search.vue'
import DNS from '../views/DNS.vue'
import Scans from '../views/Scans.vue'
import Alerts from '../views/Alerts.vue'
import Notifiers from '../views/Notifiers.vue'
import Settings from '../views/Settings.vue'

const routes = [
  { path: '/', redirect: '/host' },
  { path: '/host', component: HostLookup },
  { path: '/search', component: Search },
  { path: '/dns', component: DNS },
  { path: '/scans', component: Scans },
  { path: '/alerts', component: Alerts },
  { path: '/notifiers', component: Notifiers },
  { path: '/settings', component: Settings },
]

export const router = createRouter({
  history: createWebHashHistory(),
  routes,
})
