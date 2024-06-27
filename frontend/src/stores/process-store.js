import { defineStore } from 'pinia'
import {ref} from "vue";
import {GetAll} from "app/wailsjs/go/controllers/ConversionController";

export const useProcessStore = defineStore('process', () =>{

  const processes = ref([])

  async function  getProcesses() {
    return await GetAll()
      .then(response => {
        processes.value = response
      })
      .catch(error => {
        return false
      })
  }

  return { processes,  getProcesses }

})
