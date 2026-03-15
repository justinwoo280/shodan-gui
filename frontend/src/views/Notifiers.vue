<template>
  <div class="page-container">
    <div class="page-header">
      <div class="page-title">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#58a6ff" stroke-width="2"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/><polyline points="22,6 12,13 2,6"/></svg>
        Notifiers
      </div>
      <div class="page-subtitle">Manage notification services for network alerts</div>
    </div>

    <div class="page-toolbar">
      <button class="btn btn-primary" @click="showCreate = !showCreate">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
        New Notifier
      </button>
      <button class="btn" @click="loadNotifiers" :disabled="loading">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="23 4 23 10 17 10"/><path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"/></svg>
        Refresh
      </button>
    </div>

    <div class="page-content">
      <div class="alert-box error" v-if="error">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
        {{ error }}
      </div>
      <div class="alert-box success" v-if="successMsg">✓ {{ successMsg }}</div>

      <div class="card" style="margin-bottom:16px" v-if="showCreate">
        <div class="section-title">Create Notifier</div>
        <div class="form-group">
          <label class="form-label">Provider</label>
          <select class="input select" v-model="newNotifier.provider" @change="onProviderChange">
            <option value="">Select provider...</option>
            <option v-for="(cfg, name) in providers" :key="String(name)" :value="String(name)">{{ name }}</option>
          </select>
        </div>
        <div class="form-group">
          <label class="form-label">Description</label>
          <input class="input" v-model="newNotifier.description" placeholder="My email notifier" />
        </div>
        <div class="form-group" v-for="field in providerFields" :key="field">
          <label class="form-label">{{ field }}</label>
          <input class="input" v-model="newNotifier.args[field]" :placeholder="field" />
        </div>
        <div class="flex gap-2">
          <button class="btn btn-primary" @click="createNotifier" :disabled="submitting || !newNotifier.provider">
            <div class="spinner" v-if="submitting" style="width:12px;height:12px"></div>
            Create
          </button>
          <button class="btn" @click="showCreate = false">Cancel</button>
        </div>
      </div>

      <div v-if="loading && !notifiers.length" class="loading-overlay"><div class="spinner"></div> Loading...</div>

      <div v-if="notifiers.length" style="display:flex;flex-direction:column;gap:8px">
        <div class="card" v-for="n in notifiers" :key="n.id">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <span class="badge" :class="providerClass(n.provider)">{{ n.provider }}</span>
              <div>
                <div style="font-size:13px;font-weight:600">{{ n.description || 'No description' }}</div>
                <div style="font-size:11px;margin-top:2px">
                  <span class="tag mono">{{ n.id }}</span>
                  <span class="text-muted" style="margin-left:8px" v-for="(v, k) in n.args" :key="String(k)">
                    {{ k }}: {{ v }}
                  </span>
                </div>
              </div>
            </div>
            <button class="btn btn-sm btn-danger" @click="deleteNotifier(n.id)">Delete</button>
          </div>
        </div>
      </div>

      <div class="empty-state" v-if="!notifiers.length && !loading">
        <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/><polyline points="22,6 12,13 2,6"/></svg>
        <p>No notifiers configured. Add email, Slack, Telegram or other notification services.</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ListNotifiers, ListProviders, DeleteNotifier } from '../../wailsjs/go/main/App'

const loading = ref(false)
const submitting = ref(false)
const error = ref('')
const successMsg = ref('')
const notifiers = ref<any[]>([])
const providers = ref<any>({})
const showCreate = ref(false)
const newNotifier = ref<{ provider: string; description: string; args: Record<string, string> }>({
  provider: '', description: '', args: {}
})

const providerFields = computed(() => {
  if (!newNotifier.value.provider || !providers.value[newNotifier.value.provider]) return []
  return providers.value[newNotifier.value.provider]?.required || []
})

onMounted(() => { loadNotifiers(); loadProviders() })

async function loadNotifiers() {
  loading.value = true
  error.value = ''
  try {
    const r = await ListNotifiers()
    notifiers.value = r?.matches || []
  } catch (e: any) {
    error.value = e?.toString()
  } finally {
    loading.value = false
  }
}

async function loadProviders() {
  try {
    providers.value = await ListProviders() || {}
  } catch {}
}

function onProviderChange() {
  newNotifier.value.args = {}
}

async function createNotifier() {
  submitting.value = true
  error.value = ''
  try {
    alert('Note: Notifier creation via REST API requires form params — use the Shodan website for full management.')
    showCreate.value = false
  } catch (e: any) {
    error.value = e?.toString()
  } finally {
    submitting.value = false
  }
}

async function deleteNotifier(id: string) {
  if (!confirm('Delete this notifier?')) return
  try {
    await DeleteNotifier(id)
    successMsg.value = 'Notifier deleted'
    await loadNotifiers()
  } catch (e: any) {
    error.value = e?.toString()
  }
}

function providerClass(p: string) {
  const map: Record<string, string> = {
    email: 'badge-blue', slack: 'badge-purple', telegram: 'badge-blue',
    webhook: 'badge-yellow', pagerduty: 'badge-red'
  }
  return map[p] || 'badge-blue'
}
</script>
