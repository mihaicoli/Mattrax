<template>
  <div class="page-content">
    <div v-if="loading" class="loading">Checking Login...</div>
    <form v-else @submit.prevent="login" class="form">
      <p v-if="errorTxt" class="error-msg">{{ errorTxt }}</p>
      <input
        v-model="user.upn"
        v-on:input="errorTxt = null"
        required
        type="email"
        placeholder="chris@otbeaumont.me"
        autocomplete="username"
      />
      <input
        v-model="user.password"
        v-on:input="errorTxt = null"
        required
        type="password"
        placeholder="password"
        autocomplete="current-password"
      />
      <button>LOGIN</button>
    </form>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  data() {
    return {
      loading: false,
      errorTxt: null,
      user: {
        upn: '',
        password: '',
      },
    }
  },
  methods: {
    login() {
      this.loading = true
      this.$store
        .dispatch('authentication/login', this.user)
        .then(() => {
          if (this.$store.state.authentication.user.aud === 'dashboard') {
            this.$router.push(
              this.$route.query?.redirect_to !== undefined
                ? <string>this.$route.query?.redirect_to
                : '/'
            )
          } else if (
            this.$store.state.authentication.user.aud === 'enrollment'
          ) {
            this.$router.push('/enroll')
          } else {
            console.error(new Error('Unknown authentication token audience'))
          }
        })
        .catch((err) => {
          this.loading = false
          this.errorTxt = err
        })
    },
  },
})
</script>

<style scoped>
.form input {
  outline: 0;
  background: #f2f2f2;
  width: 100%;
  border: 0;
  margin: 0 0 15px;
  padding: 15px;
  box-sizing: border-box;
  font-size: 14px;
}
.form .error-msg {
  margin-bottom: 5px;
  color: red;
  font-size: 13px;
}
</style>
