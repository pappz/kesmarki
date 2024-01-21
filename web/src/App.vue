<template>
  <v-app app>
    <Splash :isLoading="isLoading" />
    <div v-if="!isLoading">
      <v-main>
        <v-container fluid ma-0 pa-0>
          <router-view/>
        </v-container>
      </v-main>
    </div>
  </v-app>
</template>

<script>
import Splash from '@/views/splash'

export default {
  name: 'App',
  components: {
    Splash,
  },
  mounted() {
    this.$mqtt.on('connect', function (){
      console.log("connected")
      this.isLoading = false
    }.bind(this))
    this.$mqtt.on('disconnect', function (){
      console.log("disconnected")
      this.isLoading = true
    }.bind(this))
    this.$mqtt.on('error', function(error) {
      console.log('connection failed', error)
      this.isLoading = false
    }.bind(this))
  },
  data: () => ({
    isLoading: false,
  }),
};
</script>
