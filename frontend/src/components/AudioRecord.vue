<script setup>

import { onMounted, onUpdated, reactive, ref} from "vue";
import {useConversionStore} from "stores/conversion-store";

const { uploadFile } = useConversionStore();

const props = defineProps({
    open: Boolean,
    process: Object
})

const emits = defineEmits(['save'])

let openAudio = ref(false)

onMounted(() => {
    openAudio.value = props.open
})

onUpdated(() => {
    openAudio.value = props.open
})

let audio = reactive({
  minutes: 0,
  seconds: 0,
  recording: false,
  audio: null,
  file: null,
})


let mediaRecorder;
let audioChunks = [];
let interval = ref(null)

const startRecord = async () => {
  audio.minutes = 0;
  audio.seconds = 0;
  audio.recording = !audio.recording;
  interval.value = setInterval(() =>{
    audio.seconds += 1
    if (audio.seconds === 60){
      audio.seconds = 0
      audio.minutes += 1
    }
  }, 1000);

  let stream = await navigator.mediaDevices.getUserMedia({ audio: true });
  mediaRecorder = new MediaRecorder(stream);

  mediaRecorder.ondataavailable = event => {
    audioChunks.push(event.data);
  };

  mediaRecorder.onstop = async () => {
    let audioBlob = new Blob(audioChunks, { type: 'audio/wav' });
    audioChunks = [];
    let audioUrl = URL.createObjectURL(audioBlob);
    audio.audio = new Audio(audioUrl);
    audio.file = audioBlob;
  };

  mediaRecorder.start();
}

const timeFormat = (time) => {
  if (time < 10){
    return `0${time}`;
  }
  return time;
}

const stopRecord = () => {
  clearInterval(interval.value);
  audio.recording = !audio.recording;
  mediaRecorder.stop()
}

const playAudio = () => {
  audio.audio.play();
}

const saveAudio = async () => {
  const reader = await new FileReader();
  await reader.readAsDataURL(audio.file);
  reader.onload = async () => {
    const stringFile = await reader.result
    await uploadFile(stringFile, `grabacion_${Math.random()}.wav`)
      .then(response => {
        emits("save", response)
        openAudio.value = false
      })
  };
}

const dialogHide = () => {
  if (mediaRecorder){
    mediaRecorder.stop()
  }
  audio.minutes = 0;
  audio.seconds = 0;
  audio.recording = false;
  audio.audio = null;
  audio.file = null;
}
</script>

<template>
    <q-dialog
            v-model="openAudio"
            @beforeHide="emits('close')"
            persistent
            @hide="dialogHide"
    >
        <q-card style="width: 580px; max-width: 80vw;">
            <q-card-section class="row items-center">
                <div class="text-h6">Grabar Audio</div>
                <q-space />
                <q-btn icon="close" flat round dense v-close-popup />
            </q-card-section>

            <q-card-section>
              <span>{{timeFormat(audio.minutes)}}:{{timeFormat(audio.seconds)}}</span>
              <q-linear-progress v-if="!audio.recording" :value="0" />
              <q-linear-progress v-else  indeterminate />
              <q-btn v-if="!audio.recording" @click="startRecord" class="q-mr-sm q-mt-sm" color="primary">
                <q-icon size="sm" class="float-left" name="mic" style="margin-left: -5px; margin-right: 5px" />
                <span>Grabar</span>
              </q-btn>
              <q-btn v-else @click="stopRecord" class="q-mr-sm q-mt-sm" color="red">
                <q-icon size="sm" class="float-left" name="stop_circle" style="margin-left: -5px; margin-right: 5px" />
                <span>Parar</span>
              </q-btn>
              <q-btn :disable="audio.recording" @click="playAudio" class="q-mr-sm q-mt-sm" color="secondary">
                <q-icon size="sm" class="float-left" name="play_circle_outline" style="margin-left: -5px; margin-right: 5px" />
                <span>Reproducir</span>
              </q-btn>
              <q-btn :disable="audio.recording" @click="saveAudio" class="q-mr-sm q-mt-sm" color="green">
                <q-icon size="sm" class="save" name="save" style="margin-left: -5px; margin-right: 5px" />
                <span>Guardar</span>
              </q-btn>
              <q-btn v-close-popup class="q-mr-sm q-mt-sm" color="grey-6">
                <q-icon size="sm" class="save" name="cancel" style="margin-left: -5px; margin-right: 5px" />
                <span>Cancelar</span>
              </q-btn>
            </q-card-section>
        </q-card>
    </q-dialog>
</template>

<style scoped>

</style>
