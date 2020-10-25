<template>
  <div v-if="loading" class="loading">Creating User...</div>
  <div v-else>
    <div class="panel">
      <div class="panel-head">
        <h1>
          <UserIcon view-box="0 0 24 24" height="40" width="40" />Create New
          User
        </h1>
      </div>
      <div>
        <form class="create-form" @submit.prevent="createUser">
          <input
            v-model="user.fullname"
            type="text"
            placeholder="Chris Green"
            required
          />
          <input
            v-model="user.upn"
            type="email"
            placeholder="chris@otbeaumont.me"
            required
          />
          <div class="password-field">
            <input
              v-model="user.password"
              :type="passwordVisible ? 'text' : 'password'"
              placeholder="password"
              minlength="8"
              maxlength="100"
            />
            <span
              @mousedown="passwordVisible = true"
              @mouseup="passwordVisible = false"
            >
              <EyeIcon view-box="0 0 8 8" height="24" width="24" />
            </span>
          </div>

          <button type="submit" class="submit">Create User</button>
        </form>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  layout: 'dashboard',
  data() {
    return {
      loading: false,
      passwordVisible: false,
      user: {
        fullname: '',
        upn: '',
        password: '',
      },
    }
  },
  methods: {
    createUser() {
      this.loading = true
      this.$store
        .dispatch('users/createUser', this.user)
        .then(() => this.$router.push('/users/' + this.user.upn))
        .catch((err) => this.$store.commit('dashboard/setError', err))
    },
  },
})
</script>

<style>
.create-form {
  padding: 10px;
}

.create-form input {
  display: block;
  margin: 10px;
  padding: 5px;
  width: 100%;
  max-width: 300px;
}

.password-field {
  width: 300px;
  position: relative;
}

.password-field span {
  position: absolute;
  top: 0;
  right: 0px;
  z-index: 2;
  top: 2.5px;
  cursor: pointer;
}

.submit {
  background-color: var(--primary-color-accent);
  border: none;
  color: white;
  padding: 10px 32px;
  text-align: center;
  text-decoration: none;
  display: inline-block;
  font-size: 16px;
  margin: 10px;
  width: 100%;
  max-width: 300px;
}
</style>
