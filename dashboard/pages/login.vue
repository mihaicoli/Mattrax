<template>
  <div>
    <form class="login" @submit.prevent="login">
      <h1>Sign in</h1>
      <label>Email</label>
      <input
        v-model="upn"
        required
        type="email"
        placeholder="user@example.com"
      />
      <label>Password</label>
      <input v-model="password" type="password" placeholder="password" />
      <hr />
      <button type="submit">Login</button>
    </form>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

// TODO: Password manager help, warning on login error

export default Vue.extend({
  data() {
    return {
      upn: '',
      password: '',
      errorTxt: '',
    }
  },
  methods: {
    login() {
      this.$store
        .dispatch('authentication/login', {
          upn: this.upn,
          password: this.password,
        })
        .then(() => {
          this.$router.push('/')
        })
        .catch((err) => {
          console.error(err)
        })
    },
  },
})
</script>

<style></style>
