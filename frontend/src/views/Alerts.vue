<template>
  <div class="page-container">
    <div class="page-header">
      <div class="page-title">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#58a6ff" stroke-width="2"><path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"/><path d="M13.73 21a2 2 0 0 1-3.46 0"/></svg>
        Network Alerts
      </div>
      <div class="page-subtitle">Monitor IP ranges and get notified when Shodan discovers new services</div>
    </div>

    <div class="page-toolbar">
      <button class="btn btn-primary" @click="showCreate = !showCreate">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
        New Alert
      </button>
      <button class="btn" @click="loadAlerts" :disabled="loading">
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
        <div class="section-title">Create Network Alert</div>
        <div class="form-group">
          <label class="form-label">Alert Name</label>
          <input class="input" v-model="newAlert.name" placeholder="My Network Monitor" />
        </div>
        <div class="form-group">
          <label class="form-label">IPs / Netblocks</label>
          <textarea class="input" v-model="newAlert.ips" placeholder="8.8.8.8&#10;1.1.1.1&#10;192.168.0.0/24&#10;(one per line)" rows="4"></textarea>
        </div>
        <div class="form-group">
          <label class="form-label">Expires (seconds, 0 = never)</label>
          <input class="input" type="number" v-model.number="newAlert.expires" min="0" />
        </div>
        <div class="flex gap-2">
          <button class="btn btn-primary" @click="createAlert" :disabled="submitting || !newAlert.name || !newAlert.ips">
            <div class="spinner" v-if="submitting" style="width:12px;height:12px"></div>
            Create Alert
          </button>
          <button class="btn" @click="showCreate = false">Cancel</button>
        </div>
      </div>

      <div v-if="loading && !alerts.length" class="loading-overlay"><div class="spinner"></div> Loading alerts...</div>

      <div v-if="alerts.length" style="display:flex;flex-direction:column;gap:8px">
        <div class="card" v-for="alert in alerts" :key="alert.id" :class="{ expanded: selectedAlert === alert.id }">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div>
                <div style="font-size:14px;font-weight:600">{{ alert.name }}</div>
                <div style="font-size:11px;margin-top:2px">
                  <span class="tag mono" style="margin-right:6px">{{ alert.id }}</span>
                  <span class="text-muted">{{ alert.size }} IP{{ alert.size !== 1 ? 's' : '' }}</span>
                  <span class="text-muted" style="margin-left:8px">Created {{ formatDate(alert.created) }}</span>
                </div>
              </div>
            </div>
            <div class="flex items-center gap-2">
              <span class="badge badge-green" v-if="alert.has_triggers">Triggers Active</span>
              <button class="btn btn-sm" @click="selectedAlert = selectedAlert === alert.id ? null : alert.id">
                {{ selectedAlert === alert.id ? 'Collapse' : 'Details' }}
              </button>
              <button class="btn btn-sm btn-danger" @click="deleteAlert(alert.id)">Delete</button>
            </div>
          </div>

          <div v-if="selectedAlert === alert.id" style="border-top:1px solid var(--border-subtle);margin-top:12px;padding-top:12px">
            <div class="grid-2">
              <div>
                <div class="form-label">Monitored IPs</div>
                <div class="flex gap-2" style="flex-wrap:wrap;margin-top:4px">
                  <span class="ip-chip" v-for="ip in alert.filters?.ip" :key="ip">{{ ip }}</span>
                </div>
              </div>
              <div>
                <div class="form-label">Triggers</div>
                <div v-if="Object.keys(alert.triggers || {}).length" class="flex gap-2" style="flex-wrap:wrap;margin-top:4px">
                  <span class="badge badge-green" v-for="(_, t) in alert.triggers" :key="String(t)">{{ t }}</span>
                </div>
                <span class="text-muted" style="font-size:12px" v-else>No triggers enabled</span>
              </div>
            </div>

            <div style="margin-top:12px">
              <div class="form-label">Enable Trigger</div>
              <div class="flex gap-2" style="margin-top:4px;flex-wrap:wrap">
                <select class="input select" v-model="selectedTrigger" style="width:200px">
                  <option value="">Select trigger...</option>
                  <option v-for="t in triggers" :key="t.name" :value="t.name">{{ t.name }} — {{ t.description }}</option>
                </select>
                <button class="btn btn-success btn-sm" @click="enableTrigger(alert.id)" :disabled="!selectedTrigger">Enable</button>
                <button class="btn btn-danger btn-sm" @click="disableTrigger(alert.id)" :disabled="!selectedTrigger">Disable</button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="empty-state" v-if="!alerts.length && !loading">
        <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1"><path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"/><path d="M13.73 21a2 2 0 0 1-3.46 0"/></svg>
        <p>No network alerts configured. Create one to monitor IP ranges for new discoveries.</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ListAlerts, CreateAlert, DeleteAlert, ListAlertTriggers, EnableTrigger, DisableTrigger } from '../../wailsjs/go/main/App'

const loading = ref(false)
const submitting = ref(false)
const error = ref('')
const successMsg = ref('')
const alerts = ref<any[]>([])
const triggers = ref<any[]>([])
const showCreate = ref(false)
const selectedAlert = ref<string | null>(null)
const selectedTrigger = ref('')

const newAlert = ref({ name: '', ips: '', expires: 0 })

onMounted(() => { loadAlerts(); loadTriggers() })

async function loadAlerts() {
  loading.value = true
  error.value = ''
  try {
    alerts.value = await ListAlerts() || []
  } catch (e: any) {
    error.value = e?.toString()
  } finally {
    loading.value = false
  }
}

async function loadTriggers() {
  try {
    triggers.value = await ListAlertTriggers() || []
  } catch {}
}

async function createAlert() {
  submitting.value = true
  error.value = ''
  successMsg.value = ''
  try {
    const ips = newAlert.value.ips.split(/[\n,]/).map(s => s.trim()).filter(Boolean)
    await CreateAlert(newAlert.value.name, ips, newAlert.value.expires)
    successMsg.value = `Alert "${newAlert.value.name}" created`
    newAlert.value = { name: '', ips: '', expires: 0 }
    showCreate.value = false
    await loadAlerts()
  } catch (e: any) {
    error.value = e?.toString()
  } finally {
    submitting.value = false
  }
}

async function deleteAlert(id: string) {
  if (!confirm('Delete this alert?')) return
  error.value = ''
  try {
    await DeleteAlert(id)
    successMsg.value = 'Alert deleted'
    await loadAlerts()
  } catch (e: any) {
    error.value = e?.toString()
  }
}

async function enableTrigger(alertID: string) {
  if (!selectedTrigger.value) return
  try {
    await EnableTrigger(alertID, selectedTrigger.value)
    successMsg.value = `Trigger "${selectedTrigger.value}" enabled`
    await loadAlerts()
  } catch (e: any) {
    error.value = e?.toString()
  }
}

async function disableTrigger(alertID: string) {
  if (!selectedTrigger.value) return
  try {
    await DisableTrigger(alertID, selectedTrigger.value)
    successMsg.value = `Trigger "${selectedTrigger.value}" disabled`
    await loadAlerts()
  } catch (e: any) {
    error.value = e?.toString()
  }
}

function formatDate(d: string) {
  if (!d) return '—'
  return new Date(d).toLocaleString()
}
</script>
