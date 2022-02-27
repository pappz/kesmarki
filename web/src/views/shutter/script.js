export default {
  name: 'Shutter',
  mounted() {
    this.$mqtt.subscribe('kesmarki/led')
    this.defaultColor = {}
    this.color = {
      rgba: {
        "r": 50,
        "g": 0,
        "b": 100,
        "a": 0,
      },
    };
  },
  data() {
    return {
      color: null
    }
  },
  watch: {
    color() {
      this.setLedColor()
    }
  },
  mqtt: {
    'kesmarki/led' (data) {
      var txt = new TextDecoder().decode(data)
      this.defaultColor = JSON.parse(txt)
      console.log('default color: '+this.defaultColor)
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
      this.$mqtt.publish('kesmarki/led',  JSON.stringify(this.color.rgba), {retain: true})
    }
  }
}

