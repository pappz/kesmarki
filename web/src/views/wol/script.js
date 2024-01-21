export default {
  name: 'Wol',
  methods: {
    up() {
      this.$mqtt.publish('kesmarki/wol/budafoki', 'up')
    },
  }
}

