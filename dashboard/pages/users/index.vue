<template>
  <div v-if="loading" class="loading">Loading Users...</div>
  <div v-else>
    <h1>Users</h1>
    <div class="filter-panel">
      <button @click="$router.push('/users/new')">Create New User</button>
      <input type="text" placeholder="Search..." disabled />
    </div>
    <TableView :headings="['UPN', 'Name']">
      <tr v-for="user in users" :key="user.upn">
        <td>
          <NuxtLink :to="'/users/' + user.upn" exact>{{ user.upn }}</NuxtLink>
        </td>
        <td>
          {{ user.fullname }}
        </td>
      </tr>
    </TableView>
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
      users: [],
    }
  },
  created() {
    this.$store
      .dispatch('users/getAll')
      .then((users) => {
        this.users = users
        this.loading = false
      })
      .catch((err) => this.$store.commit('dashboard/setError', err))
  },
})
</script>

<style></style>
