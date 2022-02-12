import Vue from 'vue'
import App from './App.vue'
import VueMqtt from 'vue-mqtt';
import './registerServiceWorker'
import router from './router'
import store from './store'
import vuetify from './plugins/vuetify'
import './assets/sass/main.scss'

Vue.config.productionTip = false

Vue.use(VueMqtt, 'ws://192.168.0.87:1882', {clientId: 'WebClient-' + parseInt(Math.random() * 100000)})

new Vue({
  router,
  store,
  vuetify,
  render: h => h(App)
}).$mount('#app')
