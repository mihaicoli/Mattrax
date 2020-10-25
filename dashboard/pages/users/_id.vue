<template>
  <div v-if="loading" class="loading">Loading User...</div>
  <div v-else>
    <div class="panel">
      <div class="panel-head">
        <h1>
          <UserIcon view-box="0 0 24 24" height="40" width="40" />{{
            user.fullname
          }}
        </h1>
      </div>
      <div class="panel-body">
        <div class="datapoint">
          <h2>User Principal Name:</h2>
          <p>{{ user.upn }}</p>
        </div>
        <div v-if="user.azuread_oid" class="datapoint">
          <h2>Azure AD OID:</h2>
          <p>{{ user.azuread_oid }}</p>
        </div>
        <div v-if="user.azuread_oid" class="datapoint">
          <h2>Permission Level:</h2>
          <p>
            {{
              user.permission_level.charAt(0).toUpperCase() +
              user.permission_level.slice(1)
            }}
          </p>
        </div>
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
      user: {},
    }
  },
  created() {
    this.$store
      .dispatch('users/getByID', this.$route.params.id)
      .then((user) => {
        this.user = user
        this.loading = false
      })
      .catch((err) => this.$store.commit('dashboard/setError', err))
  },
  methods: {
    navigate(pathSuffix: string) {
      this.$router.push('/users/' + this.$route.params.id + pathSuffix)
    },
  },
})
</script>

<style>
.datapoint {
  margin: 4px 10px;
}

.datapoint h2 {
  display: inline-block;
  font-size: 1.1em;
  font-weight: 700;
  padding: 0px;
}

.datapoint p {
  display: inline-block;
}
</style>
