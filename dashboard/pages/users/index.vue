<template>
  <div v-if="loading" class="loading">Loading Users...</div>
  <div v-else>
    <h1>Users</h1>
    <div class="filter-panel">
      <input type="text" placeholder="Search..." />
    </div>
    <TableView :headings="['UPN', 'Name', 'Devices']">
      <tr v-for="user in users" :key="user.upn">
        <td>
          <NuxtLink :to="'/users/' + user.upn" exact>{{ user.upn }}</NuxtLink>
        </td>
        <td>
          {{ user.name }}
        </td>
        <td>
          <NuxtLink
            v-for="device in user.devices"
            :key="device.id"
            :to="'/devices/' + device.id"
            class="table-list"
            >{{ device.name }}</NuxtLink
          >
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
      .catch((err) => {
        console.error(err)
      })
  },
})
</script>

<style></style>
