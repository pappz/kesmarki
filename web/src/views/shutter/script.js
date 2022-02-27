export default {
  name: 'Shutter',
  mounted() {
    this.color = {
      rgba: {
        "a": 0,
        "r": 0,
        "g": 0,
        "b": 0
      },
    };
  },
  data() {
    return {
      color: null,
    }
  },
  watch: {
    color() {
      this.setLedColor()
    }
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
    },
    setLedColor() {
      this.$mqtt.publish('kesmarki/led',  JSON.stringify(this.color.rgba))
    }
  }
}

