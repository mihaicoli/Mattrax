<template>
  <div v-if="loading" class="loading">Loading Groups...</div>
  <div v-else>
    <h1>Groups</h1>
    <div class="filter-panel">
      <input type="text" placeholder="Search..." disabled />
    </div>
    <TableView :headings="['ID', 'Name']">
      <tr v-for="group in groups" :key="group.id">
        <td>
          <NuxtLink :to="'/groups/' + group.id" exact>{{ group.id }}</NuxtLink>
        </td>
        <td>
          {{ group.name }}
        </td>
      </tr>
    </TableView>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  layout: 'dashboard',
  data() {
    return {
      loading: true,
      groups: [],
    }
  },
  created() {
    this.$store
      .dispatch('groups/getAll')
      .then((groups) => {
        this.groups = groups
        this.loading = false
      })
      .catch((err) => this.$store.commit('dashboard/setError', err))
  },
})
</script>

<style></style>
