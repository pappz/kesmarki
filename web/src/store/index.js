import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    user: {
      inited: true,
    }
  },
  getters: {
    user(state) {
      return state.user
    }
  },
  mutations: {
  },
  actions: {
  },
  modules: {
  }
});
