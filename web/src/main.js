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
        return process.env.VUE_APP_MQTT_PASSWORD || "unknown"
    }
}

function makeClientId () {
    // Must be unique per running client: two clients sharing a clientId make
    // the broker kick each other in an endless reconnect loop. Combine time
    // and randomness so separate tabs/devices never collide.
    var rnd = Math.random().toString(36).slice(2, 10)
    var t = Date.now().toString(36)
    return 'WebClient-' + t + '-' + rnd
}

Vue.use(VueMqtt, BROKER_URL, {
      clientId: makeClientId(),
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
