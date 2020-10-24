<template>
  <div v-if="loading" class="loading">Loading Device...</div>
  <div v-else>
    <div class="panel">
      <div class="panel-head">
        <h1>
          <PhoneIcon view-box="0 0 24 24" height="40" width="40" />{{
            device.name
          }}
        </h1>
      </div>
      <div>
        <h2 class="subtitley">{{ device.description }}</h2>
      </div>
    </div>

    <div class="w3-bar w3-black">
      <button class="w3-bar-item w3-button" @click="navigate('')">
        Overview
      </button>
      <button class="w3-bar-item w3-button" @click="navigate('/info')">
        Information
      </button>
      <button class="w3-bar-item w3-button" @click="navigate('/scope')">
        Scope
      </button>
      <button class="w3-bar-item w3-button" @click="navigate('/settings')">
        Settings
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
