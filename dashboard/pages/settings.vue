<template>
  <div v-if="loading" class="loading">Loading Settings...</div>
  <div v-else>
    <h1>Settings</h1>
    <div class="w3-bar w3-black">
      <button class="w3-bar-item w3-button" @click="navigate('')">
        Overview
      </button>
      <button class="w3-bar-item w3-button" @click="navigate('users')">
        Users
      </button>
    </div>
    <NuxtChild />
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  layout: 'dashboard',
  data() {
    return {
      loading: false,
    }
  },
  methods: {
    navigate(pathSuffix: string) {
      this.$router.push('/settings/' + pathSuffix)
    },
  },
})
</script>

<style></style>
