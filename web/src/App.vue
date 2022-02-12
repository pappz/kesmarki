<template>
  <v-app>
    <v-app-bar app
               color="indigo"
               dark>
      <v-app-bar-nav-icon
          v-if="$route.meta.mainPage"
          @click.stop="drawer = !drawer"
      />


    </v-app-bar>

    <v-main>
      <v-container fluid ma-0 pa-0>
        <router-view/>
      </v-container>
    </v-main>

    <v-bottom-navigation fixed app v-model="value">
      <v-btn value="shutter" to="/shutter">
        <span>Shutter</span>
        <v-icon>mdi-window-shutter</v-icon>
      </v-btn>
    </v-bottom-navigation>

  </v-app>
</template>

<script>

export default {
  name: 'App',
  data: () => ({
    drawer: false,
    group: null,
    value: 'recent'
  }),
  watch: {
    group () {
      this.drawer = false
    },
  },
  methods: {
    goBack() {
      return this.$router.go(-1)
    }
  }
};
</script>
