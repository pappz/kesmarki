<template src="./template.html"></template>
<script src="./script.js"></script>
<style scoped>
.view-image {
  border-radius: 12px;
  overflow: hidden;
}

/* One row below the image: valve centered, refresh pinned to the right. */
.garden-controls {
  position: relative;
  max-width: 500px;
}

.garden-controls-center {
  /* keeps the valve button horizontally centered in the row */
}

/* Refresh sits at the right edge of the row, top-aligned with the valve
   button (both start at the top of the row). */
.garden-refresh {
  position: absolute;
  right: 8px;
  top: 0;
  width: 56px !important;
  height: 56px !important;
}

.garden-refresh .v-icon {
  font-size: 30px;
}

.garden-state {
  font-size: 0.9rem;
  opacity: 0.85;
}

.stat {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  padding: 8px 4px;
}

.stat-value {
  margin-top: 4px;
  font-size: 1.05rem;
  font-weight: 600;
}

.stat-label {
  font-size: 0.7rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  opacity: 0.6;
}

/* Pulse the valve button while a command is in flight or status is unknown. */
.garden-pulse {
  animation: garden-pulse 1.2s ease-in-out infinite;
}

@keyframes garden-pulse {
  0% {
    box-shadow: 0 0 0 0 rgba(255, 158, 44, 0.6);
    opacity: 1;
  }
  70% {
    box-shadow: 0 0 0 18px rgba(255, 158, 44, 0);
    opacity: 0.75;
  }
  100% {
    box-shadow: 0 0 0 0 rgba(255, 158, 44, 0);
    opacity: 1;
  }
}
</style>
