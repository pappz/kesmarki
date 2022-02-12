import Vue from 'vue'
import '@mdi/font/css/materialdesignicons.css'
import Vuetify from 'vuetify'
import theme from './theme'

Vue.use(Vuetify)

const opts = {
    theme,
    icons: {
        iconfont: 'mdi', // default - only for display purposes
    },
}

export default new Vuetify(opts)
