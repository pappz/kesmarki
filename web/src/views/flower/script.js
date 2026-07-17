export default {
  name: 'Flower',
  mounted() {
    this.$mqtt.subscribe('kesmarki/light/flower')
  },
  data() {
    return {
      ledMode: true,
      // True while applying an incoming broker state, so the watcher
      // doesn't publish it straight back (which caused a feedback loop).
      applyingRemote: false,
    }
  },
  watch: {
    ledMode() {
      if (this.applyingRemote) {
        return
      }
      this.setMode()
    }
  },
  mqtt: {
    'kesmarki/light/flower' (data) {
      var txt = new TextDecoder().decode(data)
      var payload = JSON.parse(txt)
      var mode
      switch (payload.action) {
        case 'demo':
          mode = true
          break
        case 'off':
          mode = false
          break
        default:
          return
      }
      if (mode === this.ledMode) {
        return
      }
      this.applyingRemote = true
      this.ledMode = mode
      this.$nextTick(() => {
        this.applyingRemote = false
      })
    }
  },
  methods: {
    setMode() {
      var msg = {
        action: this.ledMode ? "demo" : "off"
      }
      this.$mqtt.publish('kesmarki/light/flower', JSON.stringify(msg), { retain: true })
    }
  }
}
