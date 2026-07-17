export default {
  name: 'Wol',
  data() {
    return {
      // null = unknown (no status received yet), true = online, false = offline
      online: null,
      // true while a wake command is in flight (button shows a spinner)
      waking: false,
      // true until the minimum wake pulse time has elapsed after a tap
      minWaking: false,
      wakeTimer: null,
      wakePoll: null,
      minWakeTimer: null,
    }
  },
  computed: {
    // Pulse while: no status yet (initial load), or a wake was just triggered.
    // A tap keeps pulsing for at least minWaking so a fast status reply doesn't
    // make it flash. Not clickable while pulsing.
    pending() {
      return this.waking || this.minWaking || this.online === null
    }
  },
  mounted() {
    this.$mqtt.subscribe('kesmarki/wol/budafoki/status')
    // Ask the server for a fresh probe now that we are connected/mounted.
    this.requestStatus()
  },
  beforeDestroy() {
    this.stopWaking()
    if (this.minWakeTimer) {
      clearTimeout(this.minWakeTimer)
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
      // Keep pulsing for at least ~1.5s even if the status confirms instantly.
      this.minWaking = true
      if (this.minWakeTimer) clearTimeout(this.minWakeTimer)
      this.minWakeTimer = setTimeout(() => {
        this.minWaking = false
        this.minWakeTimer = null
      }, 1500)
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
      }, 10000)
    },
  }
}
