<template>
  <v-app app>
    <Splash :isLoading="isLoading" :hasError="hasError" :brokerUrl="brokerUrl" />
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
import { BROKER_URL } from '@/config'

export default {
  name: 'App',
  components: {
    Splash,
  },
  mounted() {
    this.$mqtt.on('connect', function () {
      console.log('connected to', this.brokerUrl)
      this.isLoading = false
      this.hasError = false
    }.bind(this))

    this.$mqtt.on('reconnect', function () {
      console.log('reconnecting to', this.brokerUrl)
      this.hasError = false
    }.bind(this))

    this.$mqtt.on('close', function () {
      console.log('connection closed')
    }.bind(this))

    this.$mqtt.on('offline', function () {
      console.log('offline')
      this.hasError = true
      this.isLoading = true
    }.bind(this))

    this.$mqtt.on('disconnect', function () {
      console.log('disconnected')
      this.isLoading = true
      this.hasError = true
    }.bind(this))

    this.$mqtt.on('error', function (error) {
      console.log('connection failed', error)
      this.isLoading = true
      this.hasError = true
    }.bind(this))
  },
  data: () => ({
    isLoading: true,
    hasError: false,
    brokerUrl: BROKER_URL,
  }),
};
</script>
