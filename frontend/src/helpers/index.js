import {useAuthStore} from "stores/auth-store.js";
import {useRouter} from "vue-router";

export const helpers = () => {

  const goTimeFormat = (time) => {
    return time.substring(0, 16).replace("T", " ")
  }

  function base64ToBlob(base64) {
    // Decodificar la cadena base64
    const binaryString = atob(base64);
    const len = binaryString.length;
    const bytes = new Uint8Array(len);
    for (let i = 0; i < len; i++) {
      bytes[i] = binaryString.charCodeAt(i);
    }
    // Crear el blob
    const blob = new Blob([bytes], { type: 'image/png' });
    return blob;
  }

  return { goTimeFormat, base64ToBlob }
}
