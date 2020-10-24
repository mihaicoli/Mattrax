<template>
  <div v-if="loading" class="loading">Loading Device Scope...</div>
  <div v-else>
    <TableView :headings="['Groups']">
      <tr v-for="group in groups" :key="group.id">
        <td>
          <NuxtLink :to="'/groups/' + group.id" exact>{{
            group.name
          }}</NuxtLink>
        </td>
      </tr>
    </TableView>
    <TableView :headings="['Policies']">
      <tr v-for="policy in policies" :key="policy.name">
        <td>
          <NuxtLink :to="'/policies/' + policy.id" exact>{{
            policy.name
          }}</NuxtLink>
        </td>
      </tr>
    </TableView>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  data() {
    return {
      loading: true,
      groups: [],
      policies: [],
    }
  },
  created() {
    this.$store
      .dispatch('devices/getScopeByID', this.$route.params.id)
      .then((scope) => {
        this.groups = scope.groups
        this.policies = scope.policies
        this.loading = false
      })
      .catch((err) => {
        console.error(err)
      })
  },
})
</script>

<style></style>
