import Vue from 'vue'
import VueRouter from 'vue-router'
import Shutter from '@/views/shutter'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    redirect: 'shutter',
  },
  {
    path: '/shutter',
    name: 'Shutter',
    component: Shutter,
    meta: {
      title: "Shutter"
    }
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
