<script setup>
import {onMounted, ref} from 'vue'
import {useI18n} from "vue-i18n";
import {Loading, useQuasar} from "quasar";
import es from 'quasar/lang/es';


const { t, locale } = useI18n({ useScope: 'global' })

const leftDrawerOpen = ref(false)
const $q = useQuasar();


onMounted(() => {
  $q.lang.set(es)
})

const toggleLeftDrawer =  () => {
  leftDrawerOpen.value = !leftDrawerOpen.value
}


let modelo = ref('small');

const modelos = [
  { label: 'Pequeño (Prioridad: Velocidad)', value: 'small' },
  { label: 'Medio (Prioridad: Equilibrado)', value: 'medium' },
  { label: 'Grande (Prioridad: Eficiencia)', value: 'large-v2' }
]

let threads = ref(2)



const cancelLoading = () => {
  alert(Loading.isActive)
  if (Loading.isActive){
    Loading.hide()
  }else{
    Loading.show()
  }
}
</script>

<template>
  <q-layout view="hHh lpR fFf">

    <q-header elevated class="bg-primary text-white">
      <q-toolbar>
        <q-btn dense flat round icon="menu" @click="toggleLeftDrawer" />
        <q-toolbar-title>
          <img class="q-mt-sm" style="width: 300px;" src="~assets/img/Marca1.svg">
        </q-toolbar-title>
        <q-select
          v-model="locale"
          :options="[{ label: 'EN', value: 'en-US'},
                 { label: 'ES', value: 'es'}]"
          :label="$t('language')"
          label-color="white"
          color="black"
          dense
          borderless
          emit-value
          map-options
          options-dense
          style="min-width: 80px"
        />
      </q-toolbar>

    </q-header>

    <q-drawer show-if-above v-model="leftDrawerOpen" side="left" bordered>
      <q-list bordered padding class="rounded-borders" style="max-width: 350px">
        <q-item-label header class="text-bold q-pb-xs">{{$t('main_menu')}}</q-item-label>

        <q-item to="/" clickable v-ripple>
          <q-item-section avatar top>
            <q-avatar rounded>
              <img src="~assets/img/Conversión.svg" alt="">
            </q-avatar>
          </q-item-section>

          <q-item-section>
            <q-item-label lines="1">{{$t('conversion')}}</q-item-label>
            <q-item-label caption>{{$t('convert_audio_to_text')}}</q-item-label>
          </q-item-section>

          <q-item-section side>
            <q-icon name="info" color="green" />
          </q-item-section>
        </q-item>

        <q-item to="/history" clickable v-ripple>
          <q-item-section avatar top>
            <q-avatar rounded>
              <img src="~assets/img/Historial.svg" alt="">
            </q-avatar>
          </q-item-section>

          <q-item-section>
            <q-item-label lines="1">{{$t('history')}}</q-item-label>
            <q-item-label caption>{{$t('my_conversions')}}</q-item-label>
          </q-item-section>

          <q-item-section side>
            <q-icon name="info" />
          </q-item-section>
        </q-item>

        <q-separator spaced />
        <q-item-label header class="text-bold q-pb-xs">{{$t('about')}}</q-item-label>

        <q-item clickable v-ripple>
          <q-item-section avatar top>
            <q-avatar rounded>
              <img src="~assets/img/about.svg" alt="">
            </q-avatar>
          </q-item-section>

          <q-item-section>
            <q-item-label lines="1">BTEX</q-item-label>
            <q-item-label caption>{{$t('system_description')}}</q-item-label>
          </q-item-section>

          <q-item-section side>
            <q-icon name="info" />
          </q-item-section>
        </q-item>

        <q-item clickable v-ripple>
          <q-item-section avatar top>
            <q-avatar rounded>
              <img src="~assets/img/Desarr.svg" alt="">
            </q-avatar>
          </q-item-section>

          <q-item-section>
            <q-item-label lines="1">{{$t('developers')}}</q-item-label>
            <q-item-label caption>{{$t('contact_information')}}</q-item-label>
          </q-item-section>

          <q-item-section side>
            <q-icon name="info" color="amber" />
          </q-item-section>
        </q-item>
        <q-separator spaced />

        <q-item-label header class="text-bold q-pb-xs">{{$t('Modelos de Conversion')}}</q-item-label>

        <q-item class="q-px-sm q-pb-none">
          <q-item-section>
            <q-option-group
              :options="modelos"
              type="radio"
              v-model="modelo"
              class="float-left"
            />
          </q-item-section>
        </q-item>
        <q-separator spaced />

        <q-item-label header class="text-bold q-pb-xs">{{$t('Subprocesos')}}</q-item-label>

        <q-item class="q-px-sm q-pa-none">
          <q-item-section>
            <q-slider
              v-model="threads"
              marker-labels
              :min="1"
              :max="8"
              class="q-mx-sm"
              style="max-width: 95%"
            />
          </q-item-section>
        </q-item>
      </q-list>
    </q-drawer>

    <q-page-container>
      <router-view :locale="locale" :model="modelo" :threads="threads" />
    </q-page-container>

  </q-layout>
</template>

<style scoped>

</style>
