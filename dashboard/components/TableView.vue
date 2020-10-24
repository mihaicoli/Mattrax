<template>
  <table class="table">
    <thead>
      <tr>
        <th v-for="heading in headings" :key="heading">{{ heading }}</th>
      </tr>
    </thead>
    <tbody>
      <slot />
    </tbody>
  </table>
</template>

<script lang="ts">
export default {
  props: {
    headings: {
      type: Array,
      required: true,
    },
  },
}
</script>

<style>
.table {
  width: 100%;
  border: 1px solid #dee2e6;
  border-radius: 10px;
  background-color: #fff;
  padding: 5px;
  margin: 5px;
}

.table th {
  margin: 0;
  position: relative;
  padding: 0.5em;
  padding-right: 1.75em;
  min-width: 1em;
  vertical-align: middle;
  text-align: left;
  line-height: 1;
  white-space: nowrap;
  box-shadow: 0 1px 2px 0 rgba(50, 50, 50, 0.1);
  font-weight: 700;
  text-transform: uppercase;
}

.table td {
  border-bottom: 2px solid #565458;
}

.table tr:hover {
  background-color: #aaaaaa;
}

.table td a {
  font-size: 16px;
  font-weight: 300;
  text-decoration: none;
  color: #565458;
}

.table td a:hover {
  color: black;
}

.table-list:not(:empty):not(:last-child):after {
  content: ', ';
}
</style>
