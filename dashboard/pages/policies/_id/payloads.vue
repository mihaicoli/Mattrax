<template>
  <div>
    <h1>Payloads</h1>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  props: {
    policy: {
      type: Object,
      required: true,
    },
  },
})
</script>

<style></style>
