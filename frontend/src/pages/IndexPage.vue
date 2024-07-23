<template>
  <video ref="video" style="width: 100vw; height: 100vh" src="~assets/intro.mp4" />
</template>

<script setup>
import {onMounted, ref} from "vue";
import {useRouter} from "vue-router";

const router = useRouter()
const video = ref(null)

onMounted(() => {
  video.value.play()
  video.value.addEventListener('pause', () => {
    if (video.value.currentTime === video.value.duration) {
      router.push({ path: "/conversion" }) ;
    }
  })
})
</script>
