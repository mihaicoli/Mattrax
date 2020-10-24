<template>
  <div v-if="loading" class="loading">Loading User...</div>
  <div v-else-if="notfound" class="notfound">User Not Found</div>
  <div v-else>
    <div class="panel">
      <div class="panel-head">
        <h1>
          <UserIcon view-box="0 0 24 24" height="40" width="40" />{{
            user.fullname
          }}
          ({{ user.upn }})
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
  middleware: ['auth'],
  data() {
    return {
      loading: true,
      notfound: false,
      user: {},
    }
  },
  created() {
    this.$store
      .dispatch('users/getByID', this.$route.params.id)
      .then((user) => {
        if (user === null) {
          this.notfound = true
        } else {
          this.user = user
        }
        this.loading = false
      })
      .catch((err) => {
        console.error(err)
      })
  },
  methods: {
    navigate(pathSuffix: string) {
      this.$router.push('/users/' + this.$route.params.id + pathSuffix)
    },
  },
})
</script>

<style></style>
