export default {
  name: 'Flower',
  mounted() {
    this.$mqtt.subscribe('kesmarki/light/flower')
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
    'kesmarki/light/flower' (data) {
      var txt = new TextDecoder().decode(data)
      this.defaultColor = JSON.parse(txt)
      console.log('default color: '+this.defaultColor)
    }
  },
  methods: {
    setLedColor() {
      this.$mqtt.publish('kesmarki/led',  JSON.stringify(this.color.rgba), {retain: true})
    }
  }
}