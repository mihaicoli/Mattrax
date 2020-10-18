<template>
  <aside class="sidebar">
    <NuxtLink to="/" exact>
      <HomeIcon />
      <span>Dashboard</span>
    </NuxtLink>
    <NuxtLink to="/devices" exact>
      <PhoneIcon />
      <span>Devices</span>
    </NuxtLink>
    <NuxtLink to="/users" exact>
      <UserIcon />
      <span>Users</span>
    </NuxtLink>
    <NuxtLink to="/policies" exact>
      <BookIcon />
      <span>Policies</span>
    </NuxtLink>
    <!-- <NuxtLink to="/reporting" exact>
      <DatabaseIcon />
      <span>Reporting</span>
    </NuxtLink> -->
    <NuxtLink to="/settings" exact>
      <CogIcon />
      <span>Settings</span>
    </NuxtLink>
  </aside>
</template>

<script>
import HomeIcon from '@/assets/icon/home.svg?inline'
import PhoneIcon from '@/assets/icon/phone.svg?inline'
import UserIcon from '@/assets/icon/user.svg?inline'
import BookIcon from '@/assets/icon/book.svg?inline'
// import DatabaseIcon from '@/assets/icon/database.svg?inline'
import CogIcon from '@/assets/icon/cog.svg?inline'

export default {
  components: {
    HomeIcon,
    PhoneIcon,
    UserIcon,
    BookIcon,
    // DatabaseIcon,
    CogIcon,
  },
}
</script>

<style>
.sidebar {
  height: 100%;
  width: 250px;
  position: fixed;
  left: 0;
  top: 50px;
  background: var(--secondary-color);
  color: var(--light-text-color);
}

.sidebar a {
  display: block;
  padding: 10px 15px;
  text-decoration: none;
  color: inherit;
}

.sidebar a:hover {
  background: var(--secondary-color-accent);
}

.sidebar a svg {
  display: inline-block;
  vertical-align: middle;
}

.sidebar a span {
  padding-left: 8px;
  vertical-align: middle;
}
</style>
