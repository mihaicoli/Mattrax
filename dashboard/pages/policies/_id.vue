<template>
  <div v-if="loading" class="loading">Loading Policy...</div>
  <div v-else>
    <h1>Device: {{ policy.name }}</h1>
    <h2>{{ policy.description }}</h2>
    <div class="w3-bar w3-black">
      <button class="w3-bar-item w3-button" @click="navigate('')">
        Overview
      </button>
      <button class="w3-bar-item w3-button" @click="navigate('/payloads')">
        Payloads
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
      policy: {},
    }
  },
  created() {
    this.$store
      .dispatch('policies/getByID', this.$route.params.id)
      .then((policy) => {
        this.policy = policy
        this.loading = false
      })
      .catch((err) => {
        console.error(err)
      })
  },
  methods: {
    navigate(pathSuffix: string) {
      this.$router.push('/policies/' + this.$route.params.id + pathSuffix)
    },
  },
})
</script>

<style></style>
