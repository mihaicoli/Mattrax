<template>
  <div class="page-content">
    <h1 class="title">Mattrax Enrollment</h1>
    <button v-if="platform == 'Win32'" @click="triggerEnroll()">
      ENROLL THIS DEVICE
    </button>
    <p v-else>Your device is not supported!</p>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  middleware: ['auth'],
  data() {
    return {
      platform: window.navigator.platform,
    }
  },
  methods: {
    triggerEnroll() {
      window.open(
        'ms-device-enrollment:?mode=mdm&username=' +
          this.$store.state.authentication.user.upn
      )
    },
  },
})
</script>

<style>
.title {
  font-size: 29px;
}
</style>
