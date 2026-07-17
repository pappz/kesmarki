<template src="./template.html"></template>
<script src="./script.js"></script>
<style scoped>
/* Pulse the power button while a wake command is in flight (no feedback yet). */
.wol-pulse {
  animation: wol-pulse 1.2s ease-in-out infinite;
}

@keyframes wol-pulse {
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
