import { defineStore } from 'pinia'
import {Upload} from "app/wailsjs/go/controllers/ConversionController";

export const useConversionStore = defineStore('conversion', () =>{

  async function  uploadFile(file) {
    Upload(file)
      .then(function (res) {
         return { success: true, data: res.data }
      })
      .catch(function (error) {
        return { success: false, data: error }
      })
  }

  return { uploadFile }

})
