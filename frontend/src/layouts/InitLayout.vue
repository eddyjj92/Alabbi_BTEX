<script setup>

</script>

<template>
  <transition
    appear
    enter-active-class="animated fadeIn"
    leave-active-class="animated fadeOut"
  >
    <router-view class="absolute-full"/>
  </transition>
</template>

<style scoped>

</style>
