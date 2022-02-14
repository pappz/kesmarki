export default {
  name: 'Shutter',
  created () {
    console.log("shutter view loaded")
  },
  methods: {
    shutterUp() {
      this.$mqtt.publish('kesmarki/shutter', 'up')
    },
    shutterStop() {
      this.$mqtt.publish('kesmarki/shutter', 'stop')
    },
    shutterDown() {
      this.$mqtt.publish('kesmarki/shutter', 'down')
    }
  }
}
