import { defineStore } from 'pinia'
import {OpenVideo, Start, Upload} from "app/wailsjs/go/controllers/ConversionController";
import {Loading, QSpinnerFacebook, QSpinnerGears} from "quasar";

export const useConversionStore = defineStore('conversion', () =>{

  const loadingOptions = {
    spinner: QSpinnerFacebook,
    spinnerColor: 'primary',
    spinnerSize: 120,
    backgroundColor: 'black',
    message: 'Importando Archivo...',
    messageColor: 'white',
  }

  const loadingOptionsPreview = {
    spinner: QSpinnerGears,
    spinnerColor: 'primary',
    spinnerSize: 150,
    backgroundColor: 'black',
    message: '<span class="text-subtitle1 text-bold">Previsualizando Archivo de <b>Audio o Video.</b></span><br>' +
             '<span class="text-amber text-italic">Cierre el reproductor para desbloquear la ventana...</span>',
    html: true,
    messageColor: 'white',
  }

  async function uploadFile(file, filename) {
    await Loading.show(loadingOptions)
    return Upload(file, filename, filename.split('.').pop())
      .then(function (res) {
         return res
      })
      .catch(function (error) {
        return false
      })
      .finally(() => Loading.hide())
  }

  async function openVideo(file) {
    await Loading.show(loadingOptionsPreview)
    return OpenVideo(file)
      .then(function (res) {
        return true
      })
      .catch(function (error) {
        return false
      })
      .finally(() => Loading.hide())
  }

  async function startConversion(id, model, input, outputDir, threads) {
    await Loading.show(loadingOptionsPreview)
    return Start(id, model, input, outputDir, Number(threads))
      .then(function (res) {
        return res
      })
      .catch(function (error) {
        return false
      })
      .finally(() => Loading.hide())
  }

  return { uploadFile, openVideo, startConversion }

})
