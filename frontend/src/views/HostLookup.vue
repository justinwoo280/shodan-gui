<template>
  <div class="page-container">
    <div class="page-header">
      <div class="page-title">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#58a6ff" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="16"/><line x1="8" y1="12" x2="16" y2="12"/></svg>
        Host Lookup
      </div>
      <div class="page-subtitle">Retrieve all services & information found on a given host IP</div>
    </div>

    <div class="page-toolbar">
      <input
        class="input flex-1"
        placeholder="Enter IP address (e.g. 8.8.8.8)"
        v-model="ip"
        @keyup.enter="lookup"
        style="max-width: 360px"
      />
      <label class="flex items-center gap-2" style="font-size:13px;color:var(--text-secondary);cursor:pointer">
        <input type="checkbox" class="checkbox" v-model="history" /> History
      </label>
      <label class="flex items-center gap-2" style="font-size:13px;color:var(--text-secondary);cursor:pointer">
        <input type="checkbox" class="checkbox" v-model="minify" /> Minify
      </label>
      <button class="btn btn-primary" @click="lookup" :disabled="loading || !ip.trim()">
        <div class="spinner" v-if="loading" style="width:12px;height:12px"></div>
        <svg v-else width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
        Lookup
      </button>
      <template v-if="result">
        <button class="btn btn-sm" @click="exportJSON">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
          Export JSON
        </button>
      </template>
    </div>

    <div class="page-content">
      <div class="alert-box error" v-if="error">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
        {{ error }}
      </div>

      <div class="loading-overlay" v-if="loading">
        <div class="spinner"></div>
        <span>Looking up {{ ip }}...</span>
      </div>

      <div v-if="result && !loading">
        <div class="host-overview card" style="margin-bottom:16px">
          <div class="host-ip">
            <span class="ip-chip" style="font-size:18px;padding:6px 14px">{{ result.ip_str }}</span>
            <span class="badge badge-blue" style="font-size:12px">{{ result.country_name }}</span>
            <span class="badge badge-purple" v-if="result.org">{{ result.org }}</span>
            <span class="badge badge-yellow" v-if="result.asn">{{ result.asn }}</span>
          </div>
          <div class="host-meta" style="margin-top:12px">
            <div class="grid-4" style="gap:12px">
              <div class="stat-card">
                <div class="stat-icon">🔌</div>
                <div class="stat-value">{{ result.ports?.length || 0 }}</div>
                <div class="stat-label">Open Ports</div>
              </div>
              <div class="stat-card">
                <div class="stat-icon">🏢</div>
                <div class="stat-value" style="font-size:14px">{{ result.isp || '—' }}</div>
                <div class="stat-label">ISP</div>
              </div>
              <div class="stat-card">
                <div class="stat-icon">🌐</div>
                <div class="stat-value" style="font-size:14px">{{ result.country_code || '—' }}</div>
                <div class="stat-label">Country</div>
              </div>
              <div class="stat-card">
                <div class="stat-icon">🕒</div>
                <div class="stat-value" style="font-size:12px">{{ formatDate(result.last_update) }}</div>
                <div class="stat-label">Last Updated</div>
              </div>
            </div>
          </div>
        </div>

        <div style="margin-bottom:12px" v-if="result.hostnames?.length">
          <div class="section-title">Hostnames</div>
          <div class="flex gap-2" style="flex-wrap:wrap">
            <span class="tag" v-for="h in result.hostnames" :key="h">{{ h }}</span>
          </div>
        </div>

        <div style="margin-bottom:12px" v-if="result.ports?.length">
          <div class="section-title">Open Ports</div>
          <div class="flex gap-2" style="flex-wrap:wrap">
            <span class="port-chip" v-for="p in result.ports" :key="p">{{ p }}</span>
          </div>
        </div>

        <div style="margin-bottom:12px" v-if="result.tags?.length">
          <div class="section-title">Tags</div>
          <div class="flex gap-2" style="flex-wrap:wrap">
            <span class="badge badge-blue" v-for="t in result.tags" :key="t">{{ t }}</span>
          </div>
        </div>

        <div style="margin-bottom:16px" v-if="result.data?.length">
          <div class="section-title">Services ({{ result.data.length }})</div>
          <div class="service-list">
            <div
              class="service-item card"
              v-for="(svc, i) in result.data"
              :key="i"
              @click="selectedService = selectedService === i ? null : i"
              :class="{ expanded: selectedService === i }"
            >
              <div class="service-header">
                <div class="flex items-center gap-2">
                  <span class="port-chip">{{ svc.port }}/{{ svc.transport }}</span>
                  <span class="text-accent" style="font-size:13px;font-weight:600">{{ svc.product || 'Unknown' }}</span>
                  <span class="text-muted" style="font-size:12px" v-if="svc.version">v{{ svc.version }}</span>
                </div>
                <div class="flex items-center gap-2">
                  <span class="text-muted" style="font-size:11px">{{ formatDate(svc.timestamp) }}</span>
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" :style="{transform: selectedService === i ? 'rotate(180deg)' : 'none', transition: '0.2s'}"><polyline points="6 9 12 15 18 9"/></svg>
                </div>
              </div>
              <div class="service-body" v-if="selectedService === i">
                <pre class="code-block" style="max-height:200px;margin-top:10px">{{ svc.data }}</pre>
                <JsonViewer :data="svc" title="Banner Data" style="margin-top:10px" />
              </div>
            </div>
          </div>
        </div>

        <JsonViewer :data="result" title="Raw Response" :show-export="true" @export="exportJSON" />
      </div>

      <div class="empty-state" v-if="!result && !loading && !error">
        <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="16"/><line x1="8" y1="12" x2="16" y2="12"/></svg>
        <p>Enter an IP address above to retrieve host information from Shodan</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { HostInfo, ExportJSON } from '../../wailsjs/go/main/App'
import JsonViewer from '../components/common/JsonViewer.vue'

const ip = ref('')
const history = ref(false)
const minify = ref(false)
const loading = ref(false)
const error = ref('')
const result = ref<any>(null)
const selectedService = ref<number | null>(null)

async function lookup() {
  if (!ip.value.trim()) return
  loading.value = true
  error.value = ''
  result.value = null
  selectedService.value = null
  try {
    result.value = await HostInfo(ip.value.trim(), history.value, minify.value)
  } catch (e: any) {
    error.value = e?.toString() || 'Unknown error'
  } finally {
    loading.value = false
  }
}

async function exportJSON() {
  if (!result.value) return
  try {
    await ExportJSON(JSON.stringify(result.value), `host_${ip.value}`)
  } catch (e: any) {
    error.value = e?.toString()
  }
}

function formatDate(d: string) {
  if (!d) return '—'
  return new Date(d).toLocaleString()
}
</script>

<style scoped>
.host-ip { display: flex; align-items: center; gap: 10px; flex-wrap: wrap; }
.service-list { display: flex; flex-direction: column; gap: 8px; }
.service-item {
  cursor: pointer;
  transition: border-color 0.15s;
  padding: 12px 16px;
}
.service-item:hover { border-color: var(--border); }
.service-item.expanded { border-color: var(--accent); }
.service-header { display: flex; align-items: center; justify-content: space-between; }
.service-body { border-top: 1px solid var(--border-subtle); padding-top: 10px; margin-top: 10px; }
</style>
