<template>
  <div class="page-container" style="position:relative">
    <div class="page-header">
      <div class="page-title">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#58a6ff" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
        Search
      </div>
      <div class="page-subtitle">Search the Shodan database using query syntax and facets</div>
    </div>

    <div class="page-toolbar">
      <input
        class="input flex-1"
        placeholder='Query (e.g. apache country:DE, product:nginx port:80)'
        v-model="query"
        @keyup.enter="search()"
        style="max-width: 500px"
      />
      <input
        class="input"
        placeholder="Facets (e.g. country,org)"
        v-model="facets"
        style="width:180px"
      />
      <button class="btn btn-primary" @click="search()" :disabled="loading || !query.trim()">
        <div class="spinner" v-if="loading" style="width:12px;height:12px"></div>
        <svg v-else width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
        Search
      </button>
      <button class="btn" @click="countOnly()" :disabled="loading || !query.trim()" title="Count only (no credits used)">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg>
        Count
      </button>
      <button class="btn" @click="analyzeTokens()" :disabled="!query.trim()" title="Analyze query tokens">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="16 18 22 12 16 6"/><polyline points="8 6 2 12 8 18"/></svg>
        Tokens
      </button>

      <div class="toolbar-divider"></div>

      <label class="filter-toggle" :class="{ active: filterHoneypots }" title="Filter out honeypot-tagged results">
        <input type="checkbox" class="checkbox" v-model="filterHoneypots" />
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 15v-4H7l5-8v4h4l-5 8z"/></svg>
        Filter Honeypots
        <span class="badge badge-red" v-if="filterHoneypots && honeypotCount > 0">-{{ honeypotCount }}</span>
      </label>

      <button
        class="btn"
        :class="{ 'btn-primary': showExport }"
        @click="showExport = !showExport"
        :disabled="!displayMatches.length"
        title="Export results"
      >
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
        Export
      </button>
    </div>

    <div class="page-content">
      <div class="alert-box error" v-if="error">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
        {{ error }}
      </div>
      <div class="alert-box success" v-if="exportSuccess">✓ File exported successfully</div>

      <div class="loading-overlay" v-if="loading">
        <div class="spinner"></div>
        <span>Searching...</span>
      </div>

      <div v-if="countResult !== null && !loading" class="card" style="margin-bottom:16px">
        <div class="flex items-center gap-3">
          <div class="stat-value" style="font-size:28px;font-family:var(--font-mono);color:var(--accent)">
            {{ countResult.total?.toLocaleString() }}
          </div>
          <div>
            <div style="font-size:14px;font-weight:600">Total Results</div>
            <div class="text-muted" style="font-size:12px">No query credits consumed</div>
          </div>
        </div>
        <div v-if="countResult.facets && Object.keys(countResult.facets).length" style="margin-top:16px">
          <div class="section-title">Facets</div>
          <div class="grid-2" style="gap:12px">
            <div v-for="(values, key) in countResult.facets" :key="String(key)" class="card" style="padding:12px">
              <div style="font-size:11px;font-weight:700;text-transform:uppercase;color:var(--text-muted);margin-bottom:8px">{{ key }}</div>
              <div v-for="v in (values as any[]).slice(0,10)" :key="v.value" class="facet-row">
                <span style="font-size:12px;color:var(--text-secondary)">{{ v.value }}</span>
                <span class="badge badge-blue" style="font-size:11px">{{ v.count?.toLocaleString() }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-if="tokensResult && !loading" class="card" style="margin-bottom:16px">
        <div class="section-title">Query Analysis</div>
        <div class="grid-2">
          <div>
            <div class="form-label">String</div>
            <span class="tag">{{ tokensResult.string || '(none)' }}</span>
          </div>
          <div>
            <div class="form-label">Active Filters</div>
            <div class="flex gap-2" style="flex-wrap:wrap">
              <span class="badge badge-purple" v-for="f in tokensResult.filters" :key="f">{{ f }}</span>
              <span class="text-muted" v-if="!tokensResult.filters?.length">None</span>
            </div>
          </div>
        </div>
        <div v-if="tokensResult.errors?.length" class="alert-box error" style="margin-top:10px;margin-bottom:0">
          {{ tokensResult.errors.join(', ') }}
        </div>
      </div>

      <div v-if="result && !loading">
        <div class="flex items-center justify-between" style="margin-bottom:12px">
          <div class="result-count">
            <strong style="color:var(--accent)">{{ result.total?.toLocaleString() }}</strong> total
            — <strong style="color:var(--text-primary)">{{ displayMatches.length }}</strong> showing
            <template v-if="filterHoneypots && honeypotCount > 0">
              (<span class="text-error">{{ honeypotCount }} honeypots hidden</span>)
            </template>
            — page {{ currentPage }}
          </div>
          <div class="flex gap-2 items-center">
            <div class="tabs" style="margin-bottom:0;border-bottom:none">
              <span class="tab" :class="{ active: viewMode === 'table' }" @click="viewMode = 'table'">Table</span>
              <span class="tab" :class="{ active: viewMode === 'cards' }" @click="viewMode = 'cards'">Cards</span>
              <span class="tab" :class="{ active: viewMode === 'raw' }" @click="viewMode = 'raw'">Raw JSON</span>
            </div>
          </div>
        </div>

        <div v-if="viewMode === 'raw'">
          <JsonViewer :data="{ total: result.total, matches: displayMatches, facets: result.facets }" title="Search Results" />
        </div>

        <div v-else-if="viewMode === 'table'" class="card" style="padding:0;overflow:hidden">
          <table class="table">
            <thead>
              <tr>
                <th>IP</th>
                <th>Port</th>
                <th>Product</th>
                <th>Country</th>
                <th>Org</th>
                <th>Tags</th>
                <th>Timestamp</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="(m, i) in displayMatches"
                :key="i"
                @click="selectedBanner = selectedBanner === i ? null : i"
                style="cursor:pointer"
              >
                <td><span class="ip-chip">{{ m.ip_str }}</span></td>
                <td><span class="port-chip">{{ m.port }}/{{ m.transport }}</span></td>
                <td>
                  <span style="color:var(--text-primary)">{{ m.product || '—' }}</span>
                  <span class="text-muted" style="font-size:11px;margin-left:4px" v-if="m.version">v{{ m.version }}</span>
                </td>
                <td>{{ m.location?.country_code || '—' }}</td>
                <td>{{ m.org || '—' }}</td>
                <td>
                  <div class="flex gap-1" style="flex-wrap:wrap">
                    <span
                      v-for="t in (m.tags || [])"
                      :key="t"
                      class="badge"
                      :class="t === 'honeypot' ? 'badge-red' : 'badge-blue'"
                    >{{ t }}</span>
                  </div>
                </td>
                <td class="text-muted" style="font-size:11px">{{ formatDate(m.timestamp) }}</td>
              </tr>
            </tbody>
          </table>
          <div v-if="selectedBanner !== null" style="padding:16px;border-top:1px solid var(--border-subtle)">
            <JsonViewer :data="displayMatches[selectedBanner]" title="Banner Details" />
          </div>
        </div>

        <div v-else-if="viewMode === 'cards'" style="display:flex;flex-direction:column;gap:8px">
          <div class="card result-card" v-for="(m, i) in displayMatches" :key="i">
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <span class="ip-chip">{{ m.ip_str }}</span>
                <span class="port-chip">{{ m.port }}/{{ m.transport }}</span>
                <span class="text-accent" style="font-size:13px;font-weight:600">{{ m.product || '' }}</span>
                <span class="text-muted" style="font-size:12px" v-if="m.version">v{{ m.version }}</span>
              </div>
              <div class="flex items-center gap-2">
                <span
                  v-for="t in (m.tags || [])"
                  :key="t"
                  class="badge"
                  :class="t === 'honeypot' ? 'badge-red' : 'badge-blue'"
                >{{ t }}</span>
                <span class="badge badge-blue" v-if="m.location?.country_code">{{ m.location.country_code }}</span>
                <span class="text-muted" style="font-size:11px">{{ m.org }}</span>
              </div>
            </div>
            <pre class="code-block" style="max-height:80px;margin-top:8px" v-if="m.data">{{ m.data.slice(0, 300) }}</pre>
          </div>
        </div>

        <div class="pagination" v-if="result.total > 100">
          <button class="btn btn-sm" @click="goPage(currentPage - 1)" :disabled="currentPage <= 1">← Prev</button>
          <span class="page-info">Page {{ currentPage }} of {{ Math.ceil(result.total / 100) }}</span>
          <button class="btn btn-sm" @click="goPage(currentPage + 1)" :disabled="currentPage >= Math.ceil(result.total / 100)">Next →</button>
        </div>

        <div v-if="result.facets && Object.keys(result.facets).length" style="margin-top:16px">
          <div class="section-title">Facets</div>
          <div class="grid-2" style="gap:12px">
            <div v-for="(values, key) in result.facets" :key="String(key)" class="card" style="padding:12px">
              <div style="font-size:11px;font-weight:700;text-transform:uppercase;color:var(--text-muted);margin-bottom:8px">{{ key }}</div>
              <div v-for="v in (values as any[]).slice(0,10)" :key="v.value" class="facet-row">
                <span style="font-size:12px;color:var(--text-secondary)">{{ v.value }}</span>
                <span class="badge badge-blue" style="font-size:11px">{{ v.count?.toLocaleString() }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="empty-state" v-if="!result && !loading && !error && countResult === null && !tokensResult">
        <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
        <p>Enter a search query to find devices and services in the Shodan database</p>
      </div>
    </div>

    <ExportPanel
      :visible="showExport"
      :data="displayMatches"
      :total-count="displayMatches.length"
      :filename="`search_${query.slice(0,30)}`"
      :query="query"
      :facets="facets"
      :filter-honeypots="filterHoneypots"
      :total-results="result?.total"
      @close="showExport = false"
      @error="(e: string) => error = e"
      @success="onExportSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Search as ShodanSearch, SearchCount, SearchTokens } from '../../wailsjs/go/main/App'
import JsonViewer from '../components/common/JsonViewer.vue'
import ExportPanel from '../components/common/ExportPanel.vue'

const route = useRoute()
const query = ref('')

onMounted(() => {
  if (route.query.q) {
    query.value = String(route.query.q)
    search()
  }
})

const facets = ref('')
const loading = ref(false)
const error = ref('')
const result = ref<any>(null)
const countResult = ref<any>(null)
const tokensResult = ref<any>(null)
const currentPage = ref(1)
const viewMode = ref<'table' | 'cards' | 'raw'>('table')
const selectedBanner = ref<number | null>(null)
const filterHoneypots = ref(true)
const showExport = ref(false)
const exportSuccess = ref(false)

const honeypotCount = computed(() => {
  if (!result.value?.matches) return 0
  return result.value.matches.filter((m: any) =>
    Array.isArray(m.tags) && m.tags.includes('honeypot')
  ).length
})

const displayMatches = computed(() => {
  if (!result.value?.matches) return []
  if (!filterHoneypots.value) return result.value.matches
  return result.value.matches.filter((m: any) =>
    !(Array.isArray(m.tags) && m.tags.includes('honeypot'))
  )
})

async function search(page = 1) {
  if (!query.value.trim()) return
  loading.value = true
  error.value = ''
  result.value = null
  countResult.value = null
  tokensResult.value = null
  currentPage.value = page
  selectedBanner.value = null
  showExport.value = false
  try {
    result.value = await ShodanSearch(query.value.trim(), facets.value, page, true)
  } catch (e: any) {
    error.value = e?.toString() || 'Unknown error'
  } finally {
    loading.value = false
  }
}

async function countOnly() {
  if (!query.value.trim()) return
  loading.value = true
  error.value = ''
  result.value = null
  countResult.value = null
  tokensResult.value = null
  try {
    countResult.value = await SearchCount(query.value.trim(), facets.value)
  } catch (e: any) {
    error.value = e?.toString() || 'Unknown error'
  } finally {
    loading.value = false
  }
}

async function analyzeTokens() {
  if (!query.value.trim()) return
  error.value = ''
  result.value = null
  countResult.value = null
  try {
    tokensResult.value = await SearchTokens(query.value.trim())
  } catch (e: any) {
    error.value = e?.toString() || 'Unknown error'
  }
}

async function goPage(page: number) {
  if (page < 1) return
  await search(page)
}

function onExportSuccess() {
  exportSuccess.value = true
  setTimeout(() => exportSuccess.value = false, 3000)
}

function formatDate(d: string) {
  if (!d) return '—'
  return new Date(d).toLocaleString()
}
</script>

<style scoped>
.toolbar-divider {
  width: 1px;
  height: 20px;
  background: var(--border);
  flex-shrink: 0;
}

.filter-toggle {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 5px 12px;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: 13px;
  color: var(--text-secondary);
  transition: all 0.15s;
  white-space: nowrap;
}
.filter-toggle:hover { border-color: var(--text-muted); color: var(--text-primary); }
.filter-toggle .checkbox { display: none; }
.filter-toggle.active {
  border-color: var(--error);
  color: var(--error);
  background: var(--error-dim);
}

.facet-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 4px 0;
  border-bottom: 1px solid var(--border-subtle);
}
.facet-row:last-child { border-bottom: none; }
.result-card { transition: border-color 0.15s; }
.result-card:hover { border-color: var(--border); }
</style>
