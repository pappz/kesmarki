import Vue from 'vue'
import VueRouter from 'vue-router'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    redirect: 'shutter',
  },
  {
    path: '/shutter',
    name: 'Shutter',
    component: () => import('@/views/shutter'),
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
