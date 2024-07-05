<script setup>
import {useQuasar, QSpinnerFacebook, Loading} from 'quasar'
import {computed, onMounted, onUpdated, reactive, ref} from "vue";
import AudioRecord from "components/AudioRecord.vue";
import { useI18n } from "vue-i18n";
import {useConversionStore} from "stores/conversion-store.js";
import {EventsOn} from "app/wailsjs/runtime";
import {OutputFile} from "app/wailsjs/go/controllers/ConversionController";

const { uploadFile, openVideo, startConversion } = useConversionStore();

const props = defineProps({
  locale: String,
  model: Object,
  threads: Number
})

const $q = useQuasar()
let { t, locale } = useI18n({ useScope: 'global' })

let openRecordAudio = ref(false)
let inputFile = ref(null)
let consola = ref(null)

let process = reactive({
    id: null,
    metodo: null,
    input: '-',
    inputFilename: '-',
    inputFormat: null,
    output: '-',
    outputFormat: 'txt',
    outputDir: null,
    outputFolder: null,
    console: '',
    consoleLength: 1,
    processing: false,
    success: false,
    progress: 0.0,
    progressLabel: computed(() => (process.progress * 100).toFixed(0) + '%'),
    logs: 0
});

const metodos = [
    { label: 'local_file', value: 'local' },
    { label: 'Grabacion por Microfono', value: 'mic' },
    { label: 'Archivos en la Nube', value: 'nube' }
]

onUpdated(() => {
  locale = props.locale
})

onMounted(()=> {
    consoleAddText(process.consoleLength + " - " + "Bienvenidos a BTEX (Convertidor de Audio a Texto)" + "\n")

    const historyItem = JSON.parse(localStorage.getItem('historyItem'))
    if (historyItem){
        localStorage.removeItem('historyItem')
        process.id = historyItem.ID;
        process.success = historyItem.Status === 1;
        process.metodo = 'local';
        process.input = historyItem.Route + historyItem.Folder +"\\"+ historyItem.File;
        process.inputFilename = historyItem.File;
        process.inputFormat = historyItem.Extension;
        process.output = process.input.replace(process.inputFormat, process.outputFormat);
        process.outputDir = historyItem.Route + historyItem.Folder;
        process.outputFolder = historyItem.Folder;
        openDialog()
    }

})

const validaRadio = () => {
    process.input = '-';
    process.output = '-';
    process.success = false;
    process.progress = 0
    if(process.metodo === "local"){
      inputFile.value.click()
    }else if(process.metodo === "mic"){
      openRecordAudio.value = true;
    }
}

const setFile = async (e) => {
    if (e.target.files.length > 0){
      await Loading.show()
      const reader = await new FileReader();
      await reader.readAsDataURL(e.target.files[0]);
      reader.onload = async () => {
        const stringFile = await reader.result
        await consoleAddText(process.consoleLength + " - " + "Importando archivo local..." + "\n")
        await uploadFile(stringFile, e.target.files[0].name)
          .then(response => {
            if(response){
              process.id = response.id;
              process.input = `${response.folderPath}\\${response.filename}`;
              process.inputFilename = response.filename;
              process.inputFormat = response.extension;
              process.output = process.input.replace(process.inputFormat, process.outputFormat);
              process.outputDir = response.folderPath
              process.outputFolder = response.folder;
            }
          })
      };
    }
}

const start = async () => {
  await EventsOn("porciento", function (number) {
    process.progress = Number(number)/100
  })
  await startConversion(process.id, props.model, process.input, process.outputDir, props.threads)
    .then(() => {
      process.processing = !process.processing
      process.success = true;
    })

}

const stop = async () => {

}

const removeLogs = async () => {

}

const openFolder = async (execute = true) => {

}

const openOutputFile = async () => {
  await Loading.show()
  await OutputFile(process.output)
  await Loading.hide()
}

const getLogs = async () => {


}

const consoleAddText = (text) => {
    process.consoleLength++;
    process.console += text;
    consola.value.scrollTop = consola.value.scrollHeight
}

let interval = ref(null)
const initProcess = () => {
    process.progress = 0
    consoleAddText(process.consoleLength + " - Procesamiento iniciado..." + "\n")
    process.processing = !process.processing
    interval.value = setInterval(() => {
        getLogs()
    }, 3000);
    start()
}

const cancelProcess = () => {
    consoleAddText(process.consoleLength + " - Enviando señal de cancelación..." + "\n")
    clearInterval(interval.value);
    process.progress = 0
    process.processing = !process.processing
    stop()
}

const changeOutputFormat = () => {
    process.output = process.input.replace(process.inputFormat, process.outputFormat)
}

const openDialog = async () => {
  await openVideo(process.input)
}

const savedAudio = (event) => {
  process.id = event.id;
  process.input = `${event.folderPath}\\${event.filename}`;
  process.inputFilename = event.filename;
  process.inputFormat = event.extension;
  process.output = process.input.replace(process.inputFormat, process.outputFormat);
  process.outputDir = event.folderPath
  process.outputFolder = event.folder;
}
</script>

<template>
  <q-page>
    <div class="row items-center justify-center">
      <div class="col-xs-12 col-sm-4 col-md-4 col-lg-4 q-pa-sm">
        <q-card dark bordered class="my-card">
          <q-card-section>
            <div class="text-h6">{{$t('entry_methods')}}</div>
            <div class="text-subtitle2">{{$t('select_an_audio_input_method')}}</div>
          </q-card-section>

          <q-separator dark inset />

          <q-card-section>
            <div class="row">
              <q-option-group
                class="col-10"
                :options="metodos"
                type="radio"
                v-model="process.metodo"
                @update:model-value="validaRadio"
              >
                <template v-slot:label="opt">
                  <div class="row items-center">
                    <span>{{ $t(`${opt.label}`) }}</span>
                    <q-icon :name="opt.icon" color="teal" size="1.5em" class="q-ml-sm" />
                  </div>
                </template>
              </q-option-group>
              <input type="file" ref="inputFile" hidden @change="setFile($event)">
              <div class="col-2">
                <div style="width: 100%; cursor: pointer" @click="process.metodo = 'local'; inputFile.click()">
                  <q-avatar class="q-mt-sm" size="33px" rounded>
                    <img src="~assets/img/Local.svg" alt="">
                  </q-avatar>
                </div>
                <div style="width: 100%; cursor: pointer" @click="process.metodo = 'mic'">
                  <q-avatar class="q-mt-sm" size="33px" rounded>
                    <img src="~assets/img/Microf.svg" alt="">
                  </q-avatar>
                </div>
                <div style="width: 100%; cursor: pointer" @click="process.metodo = 'nube'">
                  <q-avatar class="q-mt-sm" size="33px" rounded>
                    <img src="~assets/img/nube.svg" alt="">
                  </q-avatar>
                </div>
              </div>
            </div>


          </q-card-section>
        </q-card>
      </div>
      <div class="col-xs-12 col-sm-8 col-md-8 col-lg-8 q-pa-sm">
        <q-card dark bordered class="my-card">
          <q-card-section>
            <div class="text-h6">Archivos de Salida</div>
            <div class="text-subtitle2">Archivos de Salida Resultado Conversion</div>
          </q-card-section>

          <q-separator dark inset />

          <q-card-section>
            <q-list bordered>
              <q-item :disable="process.input === '-'" :clickable="process.input !== '-'" :v-ripple="process.input !== '-'" @click="openDialog">
                <q-item-section avatar top >
                  <q-avatar class="q-mt-sm" size="33px" rounded>
                    <img src="~assets/img/audio.svg" alt="">
                  </q-avatar>
                </q-item-section>

                <q-item-section>
                  <q-item-label lines="1">Audio</q-item-label>
                  <q-item-label style="font-size: 12px">{{ process.input }}</q-item-label>
                </q-item-section>

                <q-item-section side>
                  <q-icon name="play_circle" color="green" />
                </q-item-section>
              </q-item>
              <q-separator/>
              <q-item :disable="!process.success" :clickable="process.success" :v-ripple="process.success" @click="openOutputFile">
                <q-item-section avatar top>
                  <q-avatar class="q-mt-sm" size="33px" rounded>
                    <img src="~assets/img/Texto.svg" alt="">
                  </q-avatar>
                </q-item-section>

                <q-item-section>
                  <q-item-label lines="1">Texto</q-item-label>
                  <q-item-label style="font-size: 12px">{{ process.output }}</q-item-label>
                </q-item-section>

                <q-item-section side>
                  <q-icon name="file_open" color="green" />
                </q-item-section>
              </q-item>
            </q-list>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-xs-12 col-sm-4 col-md-4 col-lg-4 q-pa-sm q-pa-sm">
        <q-card dark bordered class="my-card">
          <q-card-section>
            <div class="text-h6">{{$t('process')}}</div>
            <div class="text-subtitle2">Salida de Consola</div>
          </q-card-section>

          <q-separator dark inset />
          <q-card-section>
            <textarea ref="consola" style="width: 100%;font-size: 11px; font-weight: bold; border-radius: 4px" rows="8" readonly v-model="process.console" class="bg-grey-5" />
          </q-card-section>
        </q-card>
      </div>
      <div class="col-xs-12 col-sm-8 col-md-8 col-lg-8 q-pa-sm q-pa-sm">
        <q-card dark bordered class="my-card">
          <q-card-section>
            <div class="text-h6">Opciones</div>
            <div class="text-subtitle2">Proceso de Conversion</div>
          </q-card-section>

          <q-separator dark inset />

          <q-card-section>
            <div class="q-pb-md q-gutter-sm">
              <q-btn :disable="process.processing || process.input === '-'" @click="initProcess" color="primary" icon="play_circle" label="Iniciar" />
              <q-btn :disable="!process.processing" @click="cancelProcess" color="red" icon="stop_circle" label="Parar" />
            </div>
            <q-linear-progress stripe rounded size="25px" :value="process.progress" color="primary">
              <div class="absolute-full flex flex-center">
                <q-badge color="white" text-color="accent" :label="process.progressLabel" />
              </div>
            </q-linear-progress>

            <div class="q-py-sm">
              <div class="text-subtitle2">Formatos de Salida de Texto</div>
              <div class="q-gutter-md">
                <q-radio size="md" :disable="process.input === '-'" dense v-model="process.outputFormat" @update:model-value="changeOutputFormat" label="txt" color="primary"  val="txt" />
                <q-radio size="md" :disable="process.input === '-'" dense v-model="process.outputFormat" @update:model-value="changeOutputFormat" label="text" color="secondary" val="text" />
                <q-radio size="md" :disable="process.input === '-'" dense v-model="process.outputFormat" @update:model-value="changeOutputFormat" label="srt" color="red"  val="srt"/>
                <q-radio size="md" :disable="process.input === '-'" dense v-model="process.outputFormat" @update:model-value="changeOutputFormat" label="vtt" color="green"  val="vtt"/>
                <q-radio size="md" :disable="process.input === '-'" dense v-model="process.outputFormat" @update:model-value="changeOutputFormat" label="json" color="cyan" val="json" />
              </div>
            </div>

          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
  <AudioRecord :open="process.metodo === 'mic'" :process="process" @close="process.metodo = null" @save="savedAudio($event)" />
</template>
<style lang="sass" scoped>
.my-card
    width: 100%
    max-height: 275px
    min-height: 275px
    background-color: #022457
</style>
