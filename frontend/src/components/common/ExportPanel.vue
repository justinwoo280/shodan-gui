<template>
  <div class="export-panel" v-if="visible">
    <div class="export-header">
      <div class="flex items-center gap-2">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
        <span>Export</span>
        <span class="badge badge-blue">{{ totalCount.toLocaleString() }} records</span>
      </div>
      <button class="btn btn-icon" @click="$emit('close')" title="Close">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
      </button>
    </div>

    <div class="export-body">

      <div class="export-section" v-if="query">
        <div class="export-section-title">Scope</div>
        <div class="flex gap-2" style="flex-direction:column">
          <label class="scope-option" :class="{ active: scope === 'current' }">
            <input type="radio" v-model="scope" value="current" />
            <div class="scope-info">
              <span class="scope-label">Current page</span>
              <span class="scope-desc">{{ totalCount.toLocaleString() }} records loaded</span>
            </div>
          </label>
          <label class="scope-option" :class="{ active: scope === 'all' }">
            <input type="radio" v-model="scope" value="all" />
            <div class="scope-info">
              <span class="scope-label">All results</span>
              <span class="scope-desc">
                Up to {{ maxPages * 100 }} records
                <span v-if="totalResults">({{ totalResults.toLocaleString() }} total in DB)</span>
              </span>
            </div>
          </label>
          <div v-if="scope === 'all'" class="max-pages-row">
            <span class="text-muted" style="font-size:12px;white-space:nowrap">Max pages (100/page):</span>
            <input
              class="input"
              type="number"
              v-model.number="maxPages"
              min="1"
              max="100"
              style="width:70px;padding:4px 8px;font-size:13px"
            />
            <span class="text-muted" style="font-size:12px">= {{ (maxPages * 100).toLocaleString() }} records</span>
          </div>
          <div v-if="scope === 'all'" class="alert-box warning" style="margin:0;padding:8px 10px;font-size:12px">
            Each page consumes 1 query credit. {{ maxPages }} pages = {{ maxPages }} credits.
          </div>
        </div>
      </div>

      <div class="export-section">
        <div class="export-section-title">Format</div>
        <div class="flex gap-2">
          <label class="format-option" :class="{ active: format === 'json' }">
            <input type="radio" v-model="format" value="json" />
            <span>JSON</span>
          </label>
          <label class="format-option" :class="{ active: format === 'csv' }">
            <input type="radio" v-model="format" value="csv" />
            <span>CSV</span>
          </label>
        </div>
      </div>

      <div class="export-section">
        <div class="export-section-title">
          Fields
          <div class="flex gap-2" style="margin-left:auto">
            <button class="btn btn-sm" @click="selectAll">All</button>
            <button class="btn btn-sm" @click="selectNone">None</button>
            <button class="btn btn-sm" @click="selectDefault">Default</button>
          </div>
        </div>

        <div class="field-groups">
          <div class="field-group" v-for="group in fieldGroups" :key="group.label">
            <div class="field-group-label">{{ group.label }}</div>
            <div class="field-list">
              <label
                class="field-item"
                v-for="f in group.fields"
                :key="f.key"
                :title="f.key"
              >
                <input type="checkbox" class="checkbox" :value="f.key" v-model="selectedFields" />
                <span class="field-name">{{ f.label }}</span>
                <span class="field-key">{{ f.key }}</span>
              </label>
            </div>
          </div>
        </div>
      </div>

      <div class="export-section" v-if="selectedFields.length && previewData.length">
        <div class="export-section-title">Preview (first record)</div>
        <div class="code-block" style="max-height:120px">{{ previewText }}</div>
      </div>

      <div v-if="fetchProgress" class="fetch-progress">
        <div class="spinner" style="width:14px;height:14px"></div>
        <span>{{ fetchProgress }}</span>
      </div>
    </div>

    <div class="export-footer">
      <div class="text-muted" style="font-size:12px">
        {{ selectedFields.length }} fields selected
      </div>
      <div class="flex gap-2">
        <button class="btn" @click="$emit('close')" :disabled="exporting">Cancel</button>
        <button
          class="btn btn-primary"
          @click="doExport"
          :disabled="!selectedFields.length || exporting"
        >
          <div class="spinner" v-if="exporting" style="width:12px;height:12px"></div>
          <svg v-else width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
          {{ scope === 'all' ? 'Fetch & Export' : 'Export' }} {{ format.toUpperCase() }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ExportJSON, ExportCSV, SearchAll } from '../../../wailsjs/go/main/App'

const props = defineProps<{
  visible: boolean
  data: any[]
  filename?: string
  totalCount: number
  query?: string
  facets?: string
  filterHoneypots?: boolean
  totalResults?: number
}>()

const emit = defineEmits(['close', 'error', 'success'])

const format = ref<'json' | 'csv'>('csv')
const exporting = ref(false)
const scope = ref<'current' | 'all'>('current')
const maxPages = ref(10)
const fetchProgress = ref('')

const fieldGroups = [
  {
    label: 'Network',
    fields: [
      { key: 'ip_str',    label: 'IP Address' },
      { key: 'port',      label: 'Port' },
      { key: 'transport', label: 'Transport' },
      { key: 'asn',       label: 'ASN' },
    ]
  },
  {
    label: 'Service',
    fields: [
      { key: 'product',   label: 'Product' },
      { key: 'version',   label: 'Version' },
      { key: 'data',      label: 'Banner Data' },
      { key: 'os',        label: 'OS' },
      { key: 'cpe',       label: 'CPE' },
    ]
  },
  {
    label: 'Organization',
    fields: [
      { key: 'org',       label: 'Organization' },
      { key: 'isp',       label: 'ISP' },
    ]
  },
  {
    label: 'Location',
    fields: [
      { key: 'location.country_code',  label: 'Country Code' },
      { key: 'location.country_name',  label: 'Country Name' },
      { key: 'location.city',          label: 'City' },
      { key: 'location.region_code',   label: 'Region' },
      { key: 'location.postal_code',   label: 'Postal Code' },
      { key: 'location.latitude',      label: 'Latitude' },
      { key: 'location.longitude',     label: 'Longitude' },
    ]
  },
  {
    label: 'DNS',
    fields: [
      { key: 'hostnames', label: 'Hostnames' },
      { key: 'domains',   label: 'Domains' },
    ]
  },
  {
    label: 'Metadata',
    fields: [
      { key: 'timestamp', label: 'Timestamp' },
      { key: 'tags',      label: 'Tags' },
      { key: 'hash',      label: 'Hash' },
    ]
  },
]

const DEFAULT_FIELDS = ['ip_str', 'port', 'transport', 'product', 'version', 'org', 'location.country_code', 'location.city', 'hostnames', 'timestamp']
const ALL_FIELDS = fieldGroups.flatMap(g => g.fields.map(f => f.key))

const selectedFields = ref<string[]>([...DEFAULT_FIELDS])

watch(() => props.visible, (v) => {
  if (v) {
    selectedFields.value = [...DEFAULT_FIELDS]
    scope.value = 'current'
    fetchProgress.value = ''
  }
})

function selectAll()     { selectedFields.value = [...ALL_FIELDS] }
function selectNone()    { selectedFields.value = [] }
function selectDefault() { selectedFields.value = [...DEFAULT_FIELDS] }

function extractField(obj: any, path: string): any {
  const parts = path.split('.')
  let val = obj
  for (const p of parts) {
    if (val == null) return null
    val = val[p]
  }
  if (Array.isArray(val)) return val.join(', ')
  if (val == null) return ''
  return val
}

function pickFields(obj: any, fields: string[]): Record<string, any> {
  const result: Record<string, any> = {}
  for (const f of fields) {
    const label = f.replace('.', '_')
    result[label] = extractField(obj, f)
  }
  return result
}

const previewData = computed(() =>
  props.data.map(row => pickFields(row, selectedFields.value))
)

const previewText = computed(() => {
  if (!previewData.value.length) return ''
  return JSON.stringify(previewData.value[0], null, 2)
})

async function doExport() {
  exporting.value = true
  fetchProgress.value = ''
  try {
    let sourceData: any[]

    if (scope.value === 'all' && props.query) {
      fetchProgress.value = `Fetching up to ${maxPages.value} pages...`
      sourceData = await SearchAll(
        props.query,
        props.facets || '',
        maxPages.value,
        props.filterHoneypots ?? true
      )
      fetchProgress.value = `Fetched ${sourceData.length} records, preparing export...`
    } else {
      sourceData = props.data
    }

    const transformed = sourceData.map(row => pickFields(row, selectedFields.value))
    const json = JSON.stringify(transformed)
    const name = props.filename || 'shodan_export'

    if (format.value === 'json') {
      await ExportJSON(json, name)
    } else {
      await ExportCSV(json, name)
    }
    emit('success')
    emit('close')
  } catch (e: any) {
    emit('error', e?.toString())
  } finally {
    exporting.value = false
    fetchProgress.value = ''
  }
}
</script>

<style scoped>
.export-panel {
  position: fixed;
  right: 0;
  top: 0;
  bottom: 0;
  width: 400px;
  background: var(--bg-secondary);
  border-left: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  z-index: 100;
  box-shadow: -4px 0 20px rgba(0,0,0,0.4);
}

.export-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 16px;
  border-bottom: 1px solid var(--border-subtle);
  font-size: 14px;
  font-weight: 600;
  flex-shrink: 0;
}

.export-body {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.export-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-top: 1px solid var(--border-subtle);
  flex-shrink: 0;
}

.export-section { display: flex; flex-direction: column; gap: 8px; }
.export-section-title {
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--text-muted);
  display: flex;
  align-items: center;
  gap: 8px;
}

.scope-option {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 10px 12px;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.15s;
}
.scope-option input { margin-top: 2px; flex-shrink: 0; }
.scope-option.active { border-color: var(--accent); background: var(--accent-dim); }
.scope-info { display: flex; flex-direction: column; gap: 2px; }
.scope-label { font-size: 13px; color: var(--text-primary); font-weight: 500; }
.scope-desc { font-size: 11px; color: var(--text-muted); }

.max-pages-row {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 0 0 0;
}

.format-option {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: 13px;
  color: var(--text-secondary);
  transition: all 0.15s;
}
.format-option input { display: none; }
.format-option.active {
  border-color: var(--accent);
  color: var(--accent);
  background: var(--accent-dim);
  font-weight: 600;
}

.field-groups { display: flex; flex-direction: column; gap: 12px; }
.field-group { }
.field-group-label {
  font-size: 11px;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: 4px;
  padding-bottom: 4px;
  border-bottom: 1px solid var(--border-subtle);
}
.field-list { display: flex; flex-direction: column; gap: 2px; }
.field-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 5px 6px;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: background 0.1s;
}
.field-item:hover { background: var(--bg-tertiary); }
.field-name { font-size: 13px; color: var(--text-primary); flex: 1; }
.field-key {
  font-size: 10px;
  font-family: var(--font-mono);
  color: var(--text-muted);
  white-space: nowrap;
}

.fetch-progress {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  background: var(--accent-dim);
  border: 1px solid var(--accent);
  border-radius: var(--radius-sm);
  font-size: 13px;
  color: var(--accent);
}
</style>
