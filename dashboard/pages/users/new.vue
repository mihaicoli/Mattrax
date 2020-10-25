<template>
  <div>
    <h1>Create New User</h1>
    <input v-model="user.fullname" type="text" placeholder="Full Name" />
    <input v-model="user.upn" type="email" placeholder="UPN" />
    <input v-model="user.password" type="password" placeholder="Password" />
    <button @click.prevent="createUser">Create User</button>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  layout: 'dashboard',
  middleware: ['auth'],
  data() {
    return {
      user: {
        fullname: '',
        upn: '',
        password: '',
      },
    }
  },
  methods: {
    createUser() {
      this.$store
        .dispatch('users/createUser', this.user)
        .then(() => this.$router.push('/users'))
        .catch((err) => this.$store.commit('dashboard/setError', err))
    },
  },
})
</script>

<style></style>
