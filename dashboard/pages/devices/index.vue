<template>
  <div v-if="loading" class="loading">Loading Devices...</div>
  <div v-else>
    <h1>Devices</h1>
    <div class="filter-panel">
      <input type="text" placeholder="Search..." disabled />
    </div>
    <TableView :headings="['Name', 'Owner', 'Model', 'Groups']">
      <tr v-for="device in devices" :key="device.id">
        <td>
          <NuxtLink :to="'/devices/' + device.id" exact>{{
            device.name
          }}</NuxtLink>
        </td>
        <td>
          <NuxtLink :to="'/users/' + device.owner" exact>{{
            device.owner
          }}</NuxtLink>
        </td>
        <td>{{ device.model }}</td>
        <td>
          <NuxtLink
            v-for="group in device.groups"
            :key="group.id"
            :to="'/groups/' + group.id"
            class="group-list"
            >{{ group.name }}</NuxtLink
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
      devices: [],
    }
  },
  created() {
    this.$store
      .dispatch('devices/getAll')
      .then((devices) => {
        this.devices = devices
        this.loading = false
      })
      .catch((err) => {
        console.error(err)
      })
  },
})
</script>

<style></style>
