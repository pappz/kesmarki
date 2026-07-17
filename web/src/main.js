import Vue from 'vue'
import App from './App.vue'
import VueMqtt from 'vue-mqtt';
import './registerServiceWorker'
import router from './router'
import store from './store'
import vuetify from './plugins/vuetify'
import { BROKER_URL } from './config'
import './assets/sass/main.scss'

Vue.config.productionTip = false

function isWebView () {
    let regex = /\[KesmarkiApp \d+\]$/i
    if (regex.test(window.navigator.userAgent)) {
        return true
    }
}

function getPassword () {
    if(isWebView()) {
        return window.KesmarkiApp.getPassword()
    } else {
        return "unknown"
    }
}

Vue.use(VueMqtt, BROKER_URL, {
      clientId: 'WebClient-' + parseInt(Math.random() * 100000),
      username: 'webapp',
      keepalive: 10,
      password: getPassword(),
      queueQoSZero: false
})

var app = new Vue({
  router,
  store,
  vuetify,
  render: h => h(App),
}).$mount('#app')

export default app
