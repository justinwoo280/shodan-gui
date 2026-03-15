<template>
  <div class="page-container">
    <div class="page-header">
      <div class="page-title">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#58a6ff" stroke-width="2"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 2.83-2.83l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>
        Settings
      </div>
      <div class="page-subtitle">Configure API key and view account information</div>
    </div>

    <div class="page-content">
      <div class="grid-2" style="gap:20px;max-width:900px">
        <div>
          <div class="card" style="margin-bottom:16px">
            <div class="section-title">API Configuration</div>

            <div class="alert-box error" v-if="error">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
              {{ error }}
            </div>
            <div class="alert-box success" v-if="successMsg">✓ {{ successMsg }}</div>

            <div class="form-group">
              <label class="form-label">Shodan API Key</label>
              <div class="flex gap-2">
                <input
                  class="input input-mono flex-1"
                  :type="showKey ? 'text' : 'password'"
                  v-model="keyInput"
                  placeholder="Enter your Shodan API key"
                  @keyup.enter="saveKey"
                />
                <button class="btn btn-icon" @click="showKey = !showKey" :title="showKey ? 'Hide' : 'Show'">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <template v-if="showKey"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/><line x1="1" y1="1" x2="23" y2="23"/></template>
                    <template v-else><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/></template>
                  </svg>
                </button>
              </div>
              <div class="form-hint">
                Get your API key at
                <a href="https://account.shodan.io" target="_blank" style="color:var(--accent)">account.shodan.io</a>
              </div>
            </div>

            <div class="flex gap-2">
              <button class="btn btn-primary" @click="saveKey" :disabled="saving || !keyInput.trim()">
                <div class="spinner" v-if="saving" style="width:12px;height:12px"></div>
                <svg v-else width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"/><polyline points="17 21 17 13 7 13 7 21"/><polyline points="7 3 7 8 15 8"/></svg>
                Save API Key
              </button>
              <button class="btn" @click="testKey" :disabled="testing || !keyInput.trim()">
                <div class="spinner" v-if="testing" style="width:12px;height:12px"></div>
                Test Connection
              </button>
            </div>

            <div class="alert-box info" v-if="store.hasKey" style="margin-top:12px;margin-bottom:0">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
              API key is configured and saved locally
            </div>
          </div>

          <div class="card">
            <div class="section-title">About</div>
            <div style="font-size:13px;color:var(--text-secondary);line-height:1.8">
              <div><strong>Shodan GUI</strong> — Full-featured desktop client for the Shodan API</div>
              <div class="text-muted" style="margin-top:4px">Built with Go + Wails v2 + Vue 3</div>
              <div style="margin-top:12px">
                <a href="https://developer.shodan.io/api" target="_blank" class="btn btn-sm">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/><polyline points="15 3 21 3 21 9"/><line x1="10" y1="14" x2="21" y2="3"/></svg>
                  API Documentation
                </a>
              </div>
            </div>
          </div>
        </div>

        <div>
          <div class="card" v-if="apiInfo">
            <div class="section-title">Account & Plan</div>
            <div class="grid-2" style="gap:10px;margin-bottom:16px">
              <div class="stat-card">
                <div class="stat-icon">🔍</div>
                <div class="stat-value">{{ apiInfo.query_credits?.toLocaleString() }}</div>
                <div class="stat-label">Query Credits</div>
              </div>
              <div class="stat-card">
                <div class="stat-icon">📡</div>
                <div class="stat-value">{{ apiInfo.scan_credits?.toLocaleString() }}</div>
                <div class="stat-label">Scan Credits</div>
              </div>
              <div class="stat-card">
                <div class="stat-icon">🌐</div>
                <div class="stat-value">{{ apiInfo.monitored_ips?.toLocaleString() }}</div>
                <div class="stat-label">Monitored IPs</div>
              </div>
              <div class="stat-card">
                <div class="stat-icon">🔑</div>
                <div class="stat-value" style="font-size:14px">{{ apiInfo.unlocked_left?.toLocaleString() }}</div>
                <div class="stat-label">Unlocked Left</div>
              </div>
            </div>

            <div class="form-group" style="margin-bottom:10px">
              <div class="form-label">Plan</div>
              <span class="badge badge-blue" style="font-size:13px;padding:4px 12px">{{ apiInfo.plan }}</span>
            </div>

            <div class="flex gap-2" style="flex-wrap:wrap">
              <div class="form-group" style="margin-bottom:6px;margin-right:16px">
                <div class="form-label">HTTPS API</div>
                <span class="badge" :class="apiInfo.https ? 'badge-green' : 'badge-red'">{{ apiInfo.https ? 'Enabled' : 'Disabled' }}</span>
              </div>
              <div class="form-group" style="margin-bottom:6px;margin-right:16px">
                <div class="form-label">Telnet</div>
                <span class="badge" :class="apiInfo.telnet ? 'badge-green' : 'badge-red'">{{ apiInfo.telnet ? 'Enabled' : 'Disabled' }}</span>
              </div>
              <div class="form-group" style="margin-bottom:6px">
                <div class="form-label">Unlocked</div>
                <span class="badge" :class="apiInfo.unlocked ? 'badge-green' : 'badge-red'">{{ apiInfo.unlocked ? 'Yes' : 'No' }}</span>
              </div>
            </div>

            <div v-if="apiInfo.usage_limits" style="margin-top:8px">
              <div class="form-label">Usage Limits</div>
              <div style="font-size:12px;color:var(--text-muted);font-family:var(--font-mono)">
                <div v-for="(v, k) in apiInfo.usage_limits" :key="String(k)">
                  {{ k }}: {{ v === -1 ? 'Unlimited' : v }}
                </div>
              </div>
            </div>
          </div>

          <div class="card" v-if="store.hasKey && !apiInfo" style="margin-top:0">
            <button class="btn btn-primary w-full" @click="loadApiInfo" :disabled="loadingInfo">
              <div class="spinner" v-if="loadingInfo" style="width:12px;height:12px"></div>
              Load Account Info
            </button>
          </div>

          <div class="empty-state" v-if="!store.hasKey">
            <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 2.83-2.83l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>
            <p>Configure your API key to view account details</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { SetAPIKey, APIInfo } from '../../wailsjs/go/main/App'
import { useSettingsStore } from '../stores/settings'

const store = useSettingsStore()
const keyInput = ref('')
const showKey = ref(false)
const saving = ref(false)
const testing = ref(false)
const loadingInfo = ref(false)
const error = ref('')
const successMsg = ref('')
const apiInfo = ref<any>(null)

onMounted(async () => {
  keyInput.value = store.apiKey
  if (store.hasKey) await loadApiInfo()
})

async function saveKey() {
  saving.value = true
  error.value = ''
  successMsg.value = ''
  try {
    await store.saveKey(keyInput.value.trim())
    successMsg.value = 'API key saved successfully'
    await loadApiInfo()
  } catch (e: any) {
    error.value = e?.toString()
  } finally {
    saving.value = false
  }
}

async function testKey() {
  testing.value = true
  error.value = ''
  successMsg.value = ''
  try {
    await SetAPIKey(keyInput.value.trim())
    const info = await APIInfo()
    apiInfo.value = info
    successMsg.value = `Connected! Plan: ${info.plan} | Query credits: ${info.query_credits?.toLocaleString()}`
  } catch (e: any) {
    error.value = `Connection failed: ${e?.toString()}`
  } finally {
    testing.value = false
  }
}

async function loadApiInfo() {
  loadingInfo.value = true
  try {
    apiInfo.value = await APIInfo()
    store.apiInfo = apiInfo.value
  } catch (e: any) {
    error.value = e?.toString()
  } finally {
    loadingInfo.value = false
  }
}
</script>
