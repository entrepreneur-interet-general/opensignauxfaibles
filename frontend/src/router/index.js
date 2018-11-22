import Vue from 'vue'
import Router from 'vue-router'
import Dashboard from '@/components/Dashboard'
import Data from '@/components/Data'
import Admin from '@/components/Admin'
import Browse from '@/components/Browse'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Dashboard',
      component: Dashboard
    },
    {
      path: '/data',
      name: 'Data',
      component: Data
    },
    {
      path: '/admin',
      name: 'Admin',
      component: Admin
    },
    {
      path: '/browse',
      name: 'Browse',
      component: Browse
    }
  ]
})
