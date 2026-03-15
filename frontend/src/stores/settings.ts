import { defineStore } from 'pinia'
import { ref } from 'vue'
import { GetAPIKey, HasAPIKey, SetAPIKey, APIInfo } from '../../wailsjs/go/main/App'

export const useSettingsStore = defineStore('settings', () => {
  const apiKey = ref('')
  const hasKey = ref(false)
  const apiInfo = ref<any>(null)

  async function init() {
    hasKey.value = await HasAPIKey()
    if (hasKey.value) {
      apiKey.value = await GetAPIKey()
    }
  }

  async function saveKey(key: string) {
    await SetAPIKey(key)
    apiKey.value = key
    hasKey.value = true
  }

  async function loadApiInfo() {
    try {
      apiInfo.value = await APIInfo()
    } catch {}
  }

  return { apiKey, hasKey, apiInfo, init, saveKey, loadApiInfo }
})
