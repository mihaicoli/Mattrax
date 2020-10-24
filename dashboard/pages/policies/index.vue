<template>
  <div v-if="loading" class="loading">Loading Policies...</div>
  <div v-else>
    <h1>Policies</h1>
    <div class="filter-panel">
      <input type="text" placeholder="Search..." disabled />
    </div>
    <TableView :headings="['Name', 'Description', 'Payloads']">
      <tr v-for="policy in policies" :key="policy.id">
        <td>
          <NuxtLink :to="'/policies/' + policy.id" exact>{{
            policy.name
          }}</NuxtLink>
        </td>
        <td>
          {{ policy.description }}
        </td>
        <!-- <td>
          {{ policy.payloads.join(', ') }}
        </td> -->
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
      policies: [],
    }
  },
  created() {
    this.$store
      .dispatch('policies/getAll')
      .then((policies) => {
        this.policies = policies
        this.loading = false
      })
      .catch((err) => {
        console.error(err)
      })
  },
})
</script>

<style></style>
