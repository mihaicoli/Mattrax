<template>
  <div v-if="loading" class="loading">Loading Device Information...</div>
  <div v-else>
    <TableView
      v-for="(data, tableKey) in info"
      :key="tableKey"
      :headings="[tableKey]"
    >
      <tr v-for="(item, rowKey) in data" :key="rowKey">
        <td>{{ rowKey }}</td>
        <td>{{ item }}</td>
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
      info: {},
    }
  },
  created() {
    this.$store
      .dispatch('devices/getInformationByID', this.$route.params.id)
      .then((deviceInfo) => {
        this.info = deviceInfo
        this.loading = false
      })
      .catch((err) => {
        console.error(err)
      })
  },
})
</script>

<style></style>
