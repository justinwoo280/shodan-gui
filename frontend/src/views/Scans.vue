<template>
  <div class="page-container">
    <div class="page-header">
      <div class="page-title">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#58a6ff" stroke-width="2"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg>
        On-Demand Scans
      </div>
      <div class="page-subtitle">Request Shodan to crawl specific IPs or netblocks (requires scan credits)</div>
    </div>

    <div class="page-toolbar">
      <button class="btn btn-primary" @click="showCreate = !showCreate">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
        New Scan
      </button>
      <button class="btn" @click="loadScans" :disabled="loading">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" :class="{spinning: loading}"><polyline points="23 4 23 10 17 10"/><path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"/></svg>
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
        <div class="section-title">Submit New Scan</div>
        <div class="form-group">
          <label class="form-label">IPs / Netblocks</label>
          <textarea class="input" v-model="newScanIPs" placeholder="8.8.8.8, 1.1.1.1, 192.168.0.0/24&#10;One IP per line or comma-separated" rows="3"></textarea>
          <div class="form-hint">Consumes 1 scan credit per IP address</div>
        </div>
        <div class="flex gap-2">
          <button class="btn btn-primary" @click="submitScan" :disabled="submitting || !newScanIPs.trim()">
            <div class="spinner" v-if="submitting" style="width:12px;height:12px"></div>
            Submit Scan
          </button>
          <button class="btn" @click="showCreate = false">Cancel</button>
        </div>
        <div v-if="scanResult" class="card" style="margin-top:12px;background:var(--bg-tertiary)">
          <div class="flex gap-3 items-center">
            <div><div class="form-label">Scan ID</div><span class="tag mono">{{ scanResult.id }}</span></div>
            <div><div class="form-label">IPs</div><span class="stat-value" style="font-size:18px;font-family:var(--font-mono)">{{ scanResult.count }}</span></div>
            <div><div class="form-label">Credits Left</div><span class="text-success mono">{{ scanResult.credits_left?.toLocaleString() }}</span></div>
          </div>
        </div>
      </div>

      <div v-if="loading && !scans.length" class="loading-overlay"><div class="spinner"></div> Loading scans...</div>

      <div v-if="scans.length" class="card" style="padding:0;overflow:hidden">
        <table class="table">
          <thead>
            <tr>
              <th>Scan ID</th>
              <th>Status</th>
              <th>Size</th>
              <th>Created</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="scan in scans" :key="scan.id">
              <td><span class="tag mono">{{ scan.id }}</span></td>
              <td>
                <span class="badge" :class="statusClass(scan.status)">{{ scan.status }}</span>
              </td>
              <td class="mono">{{ scan.size }}</td>
              <td class="text-muted" style="font-size:12px">{{ formatDate(scan.created) }}</td>
              <td>
                <button class="btn btn-sm" @click="checkStatus(scan.id)">Check Status</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="card" v-if="statusResult" style="margin-top:16px">
        <div class="section-title">Scan Status</div>
        <div class="grid-4">
          <div><div class="form-label">ID</div><span class="tag mono">{{ statusResult.id }}</span></div>
          <div><div class="form-label">Status</div><span class="badge" :class="statusClass(statusResult.status)">{{ statusResult.status }}</span></div>
          <div><div class="form-label">Count</div><span class="mono">{{ statusResult.count }}</span></div>
          <div><div class="form-label">Created</div><span class="text-muted" style="font-size:12px">{{ formatDate(statusResult.created) }}</span></div>
        </div>
      </div>

      <div class="empty-state" v-if="!scans.length && !loading">
        <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg>
        <p>No active scans. Create a new scan to crawl specific IPs with Shodan.</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ListScans, Scan, ScanStatus } from '../../wailsjs/go/main/App'

const loading = ref(false)
const submitting = ref(false)
const error = ref('')
const successMsg = ref('')
const scans = ref<any[]>([])
const showCreate = ref(false)
const newScanIPs = ref('')
const scanResult = ref<any>(null)
const statusResult = ref<any>(null)

onMounted(loadScans)

async function loadScans() {
  loading.value = true
  error.value = ''
  try {
    const r = await ListScans()
    scans.value = r?.matches || []
  } catch (e: any) {
    error.value = e?.toString()
  } finally {
    loading.value = false
  }
}

async function submitScan() {
  submitting.value = true
  error.value = ''
  successMsg.value = ''
  scanResult.value = null
  try {
    const ips = newScanIPs.value.split(/[\n,]/).map(s => s.trim()).filter(Boolean).join(',')
    scanResult.value = await Scan(ips)
    successMsg.value = `Scan submitted: ${scanResult.value.id}`
    await loadScans()
  } catch (e: any) {
    error.value = e?.toString()
  } finally {
    submitting.value = false
  }
}

async function checkStatus(id: string) {
  error.value = ''
  try {
    statusResult.value = await ScanStatus(id)
  } catch (e: any) {
    error.value = e?.toString()
  }
}

function statusClass(s: string) {
  if (s === 'DONE') return 'badge-green'
  if (s === 'PROCESSING') return 'badge-yellow'
  if (s === 'QUEUE') return 'badge-blue'
  return 'badge-purple'
}

function formatDate(d: string) {
  if (!d) return '—'
  return new Date(d).toLocaleString()
}
</script>

<style scoped>
.spinning { animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }
</style>
