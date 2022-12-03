export default {
  name: 'Flower',
  mounted() {
    this.$mqtt.subscribe('kesmarki/light/flower')
  },
  data() {
    return {
      ledMode: true,
    }
  },
  watch: {
    ledMode(){
      this.setMode()
    }
  },
  mqtt: {
    'kesmarki/light/flower' (data) {
      var txt = new TextDecoder().decode(data)
      var payload = JSON.parse(txt)
      switch (payload.action) {
        case 'demo':
          this.ledMode = true
          break
        case 'off':
          this.ledMode = false
          break
      } 
    }
  },
  methods: {
    setMode() {
      var msg = {
        action: ""
      } 
      switch (this.ledMode) {
        case true:
          msg.action = "demo"
          break
        case false:
          msg.action = "off"
          break
      } 

      this.$mqtt.publish('kesmarki/light/flower',  JSON.stringify(msg), {retain: true})
    }
  }
}