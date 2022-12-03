export default {
  name: 'Shutter',
  methods: {
    shutterUp() {
      this.$mqtt.publish('kesmarki/shutter', 'up')
    },
    shutterStop() {
      this.$mqtt.publish('kesmarki/shutter', 'stop')
    },
    shutterDown() {
      this.$mqtt.publish('kesmarki/shutter', 'down')
    },
  }
}

