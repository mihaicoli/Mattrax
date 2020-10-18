<template>
  <div v-if="loading" class="loading">Loading Device...</div>
  <div v-else>
    <h1>Device: {{ device.name }}</h1>
    <h2>{{ device.description }}</h2>
    <div class="w3-bar w3-black">
      <button class="w3-bar-item w3-button" @click="navigate('')">
        Overview
      </button>
      <button class="w3-bar-item w3-button" @click="navigate('/metadata')">
        Metadata
      </button>
      <button class="w3-bar-item w3-button" @click="navigate('/policies')">
        Policies
      </button>
    </div>
    <NuxtChild />
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  layout: 'dashboard',
  middleware: ['auth'],
  data() {
    return {
      loading: true,
      device: {},
    }
  },
  created() {
    this.$store
      .dispatch('devices/getByID', this.$route.params.id)
      .then((device) => {
        this.device = device
        this.loading = false
      })
      .catch((err) => {
        console.error(err)
      })
  },
  methods: {
    navigate(pathSuffix: string) {
      this.$router.push('/devices/' + this.$route.params.id + pathSuffix)
    },
  },
})
</script>

<style></style>
