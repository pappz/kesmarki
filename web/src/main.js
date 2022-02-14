import Vue from 'vue'
import App from './App.vue'
import VueMqtt from 'vue-mqtt';
import './registerServiceWorker'
import router from './router'
import store from './store'
import vuetify from './plugins/vuetify'
import './assets/sass/main.scss'

Vue.config.productionTip = false

Vue.use(VueMqtt, 'wss://kesm.webkeyapp.com', {clientId: 'WebClient-' + parseInt(Math.random() * 100000)})

var app = new Vue({
  router,
  store,
  vuetify,
  render: h => h(App),
}).$mount('#app')

export default app
