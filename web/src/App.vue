<template>
  <v-app app>
    <v-main>
      <v-container fluid ma-0 pa-0>
        <router-view/>
      </v-container>
    </v-main>

    <v-bottom-navigation fixed app>
      <v-btn value="shutter" to="/shutter">
        <span>Shutter</span>
        <v-icon>mdi-window-shutter</v-icon>
      </v-btn>
    </v-bottom-navigation>

    <v-overlay ligth='false' opacity='1' :value="overlay">
      <v-progress-circular indeterminate size="130">
        Loading...
      </v-progress-circular>
    </v-overlay>
  </v-app>
</template>

<script>

export default {
  name: 'App',
  data: () => ({
    overlay: false
  }),
};
</script>
