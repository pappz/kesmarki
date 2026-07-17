const STATUS_TOPIC = 'zigbee2mqtt/valve/garden'
const SET_TOPIC = 'zigbee2mqtt/valve/garden/set'
const GET_TOPIC = 'zigbee2mqtt/valve/garden/get'

export default {
  name: 'Garden',
  data() {
    return {
      // null = unknown (no status yet); 'ON' / 'OFF' once known
      state: null,
      battery: null,
      flow: null,
      linkquality: null,
      currentDeviceStatus: null,
      // true = we have received at least one status (retained or live)
      gotStatus: false,
      // true while a set command is in flight (waiting for confirming status)
      waiting: false,
      // true until the minimum pulse time has elapsed after a tap
      minWaiting: false,
      // true while a manual refresh (/get) is in flight
      refreshing: false,
      // timestamp (ms) of the last refresh, to throttle to once per 3s
      lastRefresh: 0,
      waitTimer: null,
      minWaitTimer: null,
      refreshTimer: null,
    }
  },
  computed: {
    isOn() {
      return this.state === 'ON'
    },
    // Pulse while: no data yet (initial load), or a command was just sent.
    // A tap keeps pulsing for at least minWaiting so a fast status reply
    // doesn't make it flash. Solid once status is present.
    pending() {
      return this.waiting || this.minWaiting || !this.gotStatus
    },
    deviceStatusText() {
      switch (this.currentDeviceStatus) {
        case 'normal_state': return 'Normal'
        case 'water_shortage': return 'No water'
        case 'water_leakage': return 'Leakage'
        default: return '—'
      }
    },
    deviceStatusIcon() {
      switch (this.currentDeviceStatus) {
        case 'normal_state': return 'mdi-check-circle'
        case 'water_shortage': return 'mdi-water-alert'
        case 'water_leakage': return 'mdi-water-alert'
        default: return 'mdi-help-circle'
      }
    },
    deviceStatusColor() {
      return this.currentDeviceStatus === 'normal_state' ? 'green' : 'grey'
    }
  },
  mounted() {
    // Subscribe only: the broker replays the retained status. We never send a
    // /get — that would wake the battery-powered valve on every open. Status
    // arrives from the retained message and from the device's own reports
    // (including after our /set commands).
    this.$mqtt.subscribe(STATUS_TOPIC)
    // If no retained status exists yet, stop pulsing after a while so the
    // button becomes usable (the first /set will bring a status anyway).
    this.waitTimer = setTimeout(() => {
      this.gotStatus = true
      this.waitTimer = null
    }, 6000)
  },
  beforeDestroy() {
    if (this.waitTimer) clearTimeout(this.waitTimer)
    if (this.minWaitTimer) clearTimeout(this.minWaitTimer)
    if (this.refreshTimer) clearTimeout(this.refreshTimer)
  },
  mqtt: {
    'zigbee2mqtt/valve/garden'(data) {
      var txt = new TextDecoder().decode(data)
      try {
        var p = JSON.parse(txt)
        if (p.state !== undefined) this.state = p.state
        if (p.battery !== undefined) this.battery = p.battery
        if (p.flow !== undefined) this.flow = p.flow
        if (p.linkquality !== undefined) this.linkquality = p.linkquality
        if (p.current_device_status !== undefined) this.currentDeviceStatus = p.current_device_status
      } catch (e) {
        // ignore malformed payloads
      }
      // A status confirms the last command landed and that we know the state.
      this.gotStatus = true
      this.stopWaiting()
      // A fresh reading arrived — end the refresh spinner.
      this.refreshing = false
      if (this.refreshTimer) {
        clearTimeout(this.refreshTimer)
        this.refreshTimer = null
      }
    }
  },
  methods: {
    toggle() {
      if (this.pending) return
      this.waiting = true
      // Keep pulsing for at least ~1.5s even if the status confirms instantly.
      this.minWaiting = true
      if (this.minWaitTimer) clearTimeout(this.minWaitTimer)
      this.minWaitTimer = setTimeout(() => {
        this.minWaiting = false
        this.minWaitTimer = null
      }, 1500)
      var next = this.isOn ? 'OFF' : 'ON'
      this.$mqtt.publish(SET_TOPIC, JSON.stringify({ state: next }))
      // Clear the waiting state if no status confirms within 10s.
      this.waitTimer = setTimeout(() => {
        this.stopWaiting()
      }, 10000)
    },
    stopWaiting() {
      this.waiting = false
      if (this.waitTimer) {
        clearTimeout(this.waitTimer)
        this.waitTimer = null
      }
    },
    refresh() {
      // Throttle to once per 3s so the battery-powered valve isn't spammed.
      var now = Date.now()
      if (this.refreshing || now - this.lastRefresh < 3000) {
        return
      }
      this.lastRefresh = now
      this.refreshing = true
      // Ask the device for a fresh reading (wakes it). The status handler
      // clears the spinner when the reply arrives.
      this.$mqtt.publish(GET_TOPIC, JSON.stringify({ state: '' }))
      // Safety: clear the spinner if no reply comes within 5s.
      if (this.refreshTimer) clearTimeout(this.refreshTimer)
      this.refreshTimer = setTimeout(() => {
        this.refreshing = false
        this.refreshTimer = null
      }, 5000)
    }
  }
}
