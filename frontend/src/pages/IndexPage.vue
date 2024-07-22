<template>
  <video ref="video" style="width: 100vw; height: 100vh" src="~assets/intro.mp4">

  </video>
</template>

<script setup>
import {onMounted, ref} from "vue";

const video = ref(null)

onMounted(() => {
  video.value.play()
  video.value.addEventListener('pause', () => {
    if (video.value.currentTime === video.value.duration) {
      window.location.href = "conversion";
    }
  })
})
</script>
