<template>
  <div v-if="loading" class="loading">Loading Group...</div>
  <div v-else>
    <div class="panel">
      <div class="panel-head">
        <h1>
          <GridIcon view-box="0 0 8 8" height="33" width="33" />{{ group.name }}
        </h1>
      </div>
    </div>

    <!-- <div class="w3-bar w3-black">
      <button class="w3-bar-item w3-button" @click="navigate('')">
        Overview
      </button>
    </div> -->
    <NuxtChild />
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  layout: 'dashboard',
  data() {
    return {
      loading: true,
      group: {},
    }
  },
  created() {
    this.$store
      .dispatch('groups/getByID', this.$route.params.id)
      .then((group) => {
        this.group = group
        this.loading = false
      })
      .catch((err) => this.$store.commit('dashboard/setError', err))
  },
  methods: {
    navigate(pathSuffix: string) {
      this.$router.push('/groups/' + this.$route.params.id + pathSuffix)
    },
  },
})
</script>

<style></style>
