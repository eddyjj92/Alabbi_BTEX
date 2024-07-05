<script setup>
import { useQuasar, QSpinnerFacebook } from 'quasar'
import {computed, onMounted, reactive, ref} from "vue";
import {useProcessStore} from "stores/process-store.js";
import {useRouter} from "vue-router";
import {storeToRefs} from "pinia";

const $q = useQuasar()
const router = useRouter();
const processStore = useProcessStore()
const {processes} = storeToRefs(processStore)

const columns = [
    {
        name: 'File',
        required: true,
        label: 'Archivo',
        align: 'left',
        field: row => row.File,
        format: val => `${val}`,
        sortable: true
    },
    { name: 'Folder', align: 'left', label: 'Carpeta', field: 'Folder', sortable: true },
    { name: 'CreatedAt', align: 'left', label: 'Fecha Procesamiento', field: 'CreatedAt', sortable: true },
]

let selected = ref([]);

onMounted(async () => {
    await processStore.getProcesses()
})

const preview = () => {
    localStorage.setItem('historyItem', JSON.stringify(selected.value[0]))
    router.push({ path : '/' })
}

const removeProcess  = async () => {
  await processStore.deleteProcess(selected.value[0].ID)
  await processStore.getProcesses()
}

const removeHistory = async () => {
  await processStore.deleteHistory()
  await processStore.getProcesses()
}
</script>

<template>
  <q-page>
    <div class="q-py-md q-my-none">
      <h4 class="text-h4 q-mx-md q-pa-none q-ma-none">Historial de Conversiones</h4>
    </div>
    <div class="q-mx-md">
      <q-table
        flat bordered
        title="Historial"
        :rows="processes"
        :columns="columns"
        row-key="ID"
        selection="single"
        v-model:selected="selected"
        :pagination="{rowsPerPage: 7}"
      >
        <template v-slot:top>
          <q-btn @click="preview" :icon="selected.length > 0 ? 'play_arrow' : 'do_disturb'" :color="selected.length > 0 ? 'blue' : 'grey-4'" :disabled="selected.length === 0" :ripple="selected.length > 0" :unelevated="selected.length === 0" :text-color="selected.length === 0 ? 'black' : 'white'" class="q-mr-sm"  label="Visualizar" />
          <q-btn :icon="selected.length > 0 ? 'folder_open' : 'folder_off'" :color="selected.length > 0 ? 'blue' : 'grey-4'" :disabled="selected.length === 0" :ripple="selected.length > 0" :unelevated="selected.length === 0" :text-color="selected.length === 0 ? 'black' : 'white'" label="Salida" class="q-mr-sm" />
          <q-btn @click="removeProcess" :icon="selected.length > 0 ? 'delete_forever' : 'do_disturb'" :color="selected.length > 0 ? 'red' : 'grey-4'" :disabled="selected.length === 0" :ripple="selected.length > 0" :unelevated="selected.length === 0" :text-color="selected.length === 0 ? 'black' : 'white'" label="Eliminar Archivo" class="q-mr-sm" />
          <q-btn @click="removeHistory" :icon="processes.length > 0 ? 'delete_forever' : 'do_disturb'" :color="processes.length > 0 ? 'red' : 'grey-4'" :disabled="processes.length === 0" :ripple="processes.length > 0" :unelevated="processes.length === 0" :text-color="processes.length === 0 ? 'black' : 'white'" label="Borrar Historial" />
        </template>
      </q-table>
    </div>
  </q-page>
</template>
