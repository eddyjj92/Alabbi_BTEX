import { defineStore } from 'pinia'
import {Upload} from "app/wailsjs/go/controllers/ConversionController";
import {Loading, QSpinnerFacebook} from "quasar";

export const useConversionStore = defineStore('conversion', () =>{

  const loadingOptions = {
    spinner: QSpinnerFacebook,
    spinnerColor: 'primary',
    spinnerSize: 120,
    backgroundColor: 'black',
    message: 'Importando Archivo...',
    messageColor: 'white',
  }

  async function uploadFile(file, filename) {
    await Loading.show(loadingOptions)
    Upload(file, filename)
      .then(function (res) {
         return { success: true, data: res.data }
      })
      .catch(function (error) {
        return { success: false, data: error }
      })
      .finally(() => Loading.hide())
  }

  return { uploadFile }

})
