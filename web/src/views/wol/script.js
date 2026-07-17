export default {
  name: 'Wol',
  data() {
    return {
      // null = unknown (no status received yet), true = online, false = offline
      online: null,
      // true while a wake command is in flight (button shows a spinner)
      waking: false,
      wakeTimer: null,
      wakePoll: null,
      // Keep the pulse running for at least ~2 animation cycles so a fast
      // status reply doesn't make it flash for a single frame.
      minPulseDone: false,
      pulseTimer: null,
    }
  },
  computed: {
    // Transitional state: waking up, or no status yet. Button pulses and is
    // not clickable until we know the machine's state.
    pending() {
      return this.waking || this.online === null || !this.minPulseDone
    }
  },
  mounted() {
    this.startMinPulse()
    this.$mqtt.subscribe('kesmarki/wol/budafoki/status')
    // Ask the server for a fresh probe now that we are connected/mounted.
    this.requestStatus()
  },
  beforeDestroy() {
    this.stopWaking()
    if (this.pulseTimer) {
      clearTimeout(this.pulseTimer)
    }
  },
  mqtt: {
    'kesmarki/wol/budafoki/status'(data) {
      var txt = new TextDecoder().decode(data)
      try {
        var payload = JSON.parse(txt)
        this.online = payload.online
      } catch (e) {
        this.online = null
      }
      // Once the machine reports online, the wake succeeded.
      if (this.online === true) {
        this.stopWaking()
      }
    }
  },
  methods: {
    requestStatus() {
      // Non-empty payload on purpose: an empty (0-byte) publish frame breaks
      // the WebSocket transport with this broker/proxy stack. The server
      // ignores the payload content.
      this.$mqtt.publish('kesmarki/wol/budafoki/status/get', 'get')
    },
    startMinPulse() {
      // Guarantee ~2 pulse cycles (2 x 1.2s) before the pulse can stop.
      this.minPulseDone = false
      if (this.pulseTimer) {
        clearTimeout(this.pulseTimer)
      }
      this.pulseTimer = setTimeout(() => {
        this.minPulseDone = true
        this.pulseTimer = null
      }, 2400)
    },
    stopWaking() {
      this.waking = false
      if (this.wakeTimer) {
        clearTimeout(this.wakeTimer)
        this.wakeTimer = null
      }
      if (this.wakePoll) {
        clearInterval(this.wakePoll)
        this.wakePoll = null
      }
    },
    up() {
      // Ignore taps while pulsing (waking, or status not known yet).
      if (this.pending) {
        return
      }
      this.waking = true
      this.startMinPulse()
      this.$mqtt.publish('kesmarki/wol/budafoki', 'up')
      // The machine needs time to boot, so the server's immediate ping after
      // the wake command will still report offline. Poll the status while the
      // machine boots; the status handler clears the spinner once it reports
      // online, and the timeout clears it if it never comes up.
      this.wakePoll = setInterval(() => {
        this.requestStatus()
      }, 5000)
      this.wakeTimer = setTimeout(() => {
        this.stopWaking()
      }, 60000)
    },
  }
}
