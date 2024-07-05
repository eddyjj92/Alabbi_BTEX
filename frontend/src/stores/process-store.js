import { defineStore } from 'pinia'
import {ref} from "vue";
import {Delete, DeleteHistory, GetAll} from "app/wailsjs/go/controllers/ConversionController";
import {Loading} from "quasar";

export const useProcessStore = defineStore('process', () =>{
  const processes = ref([])

  async function  getProcesses() {
    await Loading.show()
    return await GetAll()
      .then(response => {
        processes.value = response
      })
      .catch(error => {
        return false
      })
      .finally(() => Loading.hide())
  }

  async function deleteProcess(id) {
    await Loading.show()
    return await Delete(Number(id))
      .then(response => {
        return true
      })
      .catch(error => {
        return false
      })
      .finally(() => Loading.hide())
  }

  async function deleteHistory() {
    await Loading.show()
    return await DeleteHistory()
      .then(response => {
        return true
      })
      .catch(error => {
        return false
      })
      .finally(() => Loading.hide())
  }

  return { processes,  getProcesses, deleteProcess, deleteHistory }

})
