<template>
  <div v-if="loading" class="loading">Loading Policy...</div>
  <div v-else>
    <div class="panel">
      <div class="panel-head">
        <h1>
          <BookIcon view-box="0 0 24 24" height="40" width="40" />{{
            policy.name
          }}
        </h1>
      </div>
      <div>
        <h2 class="subtitley">{{ policy.description }}</h2>
      </div>
    </div>

    <div class="w3-bar w3-black">
      <button class="w3-bar-item w3-button" @click="navigate('')">
        Overview
      </button>
      <button class="w3-bar-item w3-button" @click="navigate('/payloads')">
        Payloads
      </button>
    </div>
    <NuxtChild :policy="policy" />
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
      .catch((err) => this.$store.commit('dashboard/setError', err))
  },
  methods: {
    navigate(pathSuffix: string) {
      this.$router.push('/policies/' + this.$route.params.id + pathSuffix)
    },
  },
})
</script>

<style></style>
