<template>
  <div>
    <TableView :headings="['Information']">
      <tr>
        <td>Name</td>
        <td>{{ policy.name }}</td>
      </tr>
      <tr>
        <td>Description</td>
        <td>{{ policy.description }}</td>
      </tr>
      <tr>
        <td>Managed Areas</td>
        <td>{{ policy.payloads.join(', ') }}</td>
      </tr>
    </TableView>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  props: {
    policy: {
      type: Object,
      required: true,
    },
  },
})
</script>

<style></style>
