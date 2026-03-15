<template>
  <div class="page-container">
    <div class="page-header">
      <div class="page-title">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#58a6ff" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/></svg>
        Query Directory
      </div>
      <div class="page-subtitle">Browse saved search queries shared by the Shodan community</div>
    </div>

    <div class="page-toolbar">
      <input class="input flex-1" placeholder="Search queries (e.g. webcam, ics, camera)" v-model="searchQ" @keyup.enter="searchQueries" style="max-width:360px" />
      <button class="btn btn-primary" @click="searchQueries" :disabled="loading">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
        Search
      </button>
      <button class="btn" @click="loadQueries" :disabled="loading">Browse All</button>
      <select class="input select" v-model="sortBy" @change="loadQueries" style="width:130px">
        <option value="timestamp">Newest</option>
        <option value="votes">Most Voted</option>
      </select>
    </div>

    <div class="page-content">
      <div class="alert-box error" v-if="error">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
        {{ error }}
      </div>

      <div class="flex gap-3" style="margin-bottom:16px;align-items:flex-start">
        <div style="flex:1">
          <div v-if="loading" class="loading-overlay"><div class="spinner"></div></div>

          <div v-if="queries.length" style="display:flex;flex-direction:column;gap:8px">
            <div class="card query-card" v-for="q in queries" :key="q.title + q.query">
              <div class="flex items-center justify-between" style="margin-bottom:6px">
                <div style="font-size:14px;font-weight:600;color:var(--text-primary)">{{ q.title }}</div>
                <div class="flex items-center gap-2">
                  <span class="badge badge-yellow">▲ {{ q.votes }}</span>
                  <span class="text-muted" style="font-size:11px">{{ formatDate(q.timestamp) }}</span>
                </div>
              </div>
              <div class="code-block" style="max-height:60px;margin-bottom:8px;cursor:pointer" @click="useQuery(q.query)" title="Click to use in search">{{ q.query }}</div>
              <div class="flex items-center justify-between">
                <div class="flex gap-2" style="flex-wrap:wrap">
                  <span class="badge badge-blue" v-for="t in q.tags?.filter((x: string) => x)" :key="t">{{ t }}</span>
                </div>
                <div class="flex gap-2">
                  <button class="btn btn-sm" @click="useQuery(q.query)" title="Open in Search">
                    <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
                    Use Query
                  </button>
                </div>
              </div>
              <div v-if="q.description" style="font-size:12px;color:var(--text-muted);margin-top:6px">{{ q.description }}</div>
            </div>
          </div>

          <div class="pagination" v-if="total > 10">
            <button class="btn btn-sm" @click="prevPage" :disabled="page <= 1">← Prev</button>
            <span class="page-info">Page {{ page }} of {{ Math.ceil(total / 10) }} ({{ total.toLocaleString() }} queries)</span>
            <button class="btn btn-sm" @click="nextPage" :disabled="page >= Math.ceil(total / 10)">Next →</button>
          </div>

          <div class="empty-state" v-if="!queries.length && !loading">
            <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
            <p>Browse the community directory or search for queries</p>
          </div>
        </div>

        <div style="width:200px;flex-shrink:0">
          <div class="card">
            <div class="section-title">Popular Tags</div>
            <div class="flex gap-2" style="flex-wrap:wrap" v-if="popularTags.length">
              <span
                class="badge badge-blue"
                v-for="t in popularTags"
                :key="t.value"
                style="cursor:pointer"
                @click="searchQ = t.value; searchQueries()"
                :title="`${t.count} queries`"
              >{{ t.value }}</span>
            </div>
            <div class="text-muted" style="font-size:12px" v-else>Loading...</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ListQueries, SearchQueries, QueryTags } from '../../wailsjs/go/main/App'

const router = useRouter()
const loading = ref(false)
const error = ref('')
const queries = ref<any[]>([])
const popularTags = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const searchQ = ref('')
const sortBy = ref('timestamp')
const isSearchMode = ref(false)

onMounted(() => { loadQueries(); loadTags() })

async function loadQueries() {
  loading.value = true
  error.value = ''
  isSearchMode.value = false
  try {
    const r = await ListQueries(page.value, sortBy.value, 'desc')
    queries.value = r?.matches || []
    total.value = r?.total || 0
  } catch (e: any) {
    error.value = e?.toString()
  } finally {
    loading.value = false
  }
}

async function searchQueries() {
  if (!searchQ.value.trim()) { loadQueries(); return }
  loading.value = true
  error.value = ''
  isSearchMode.value = true
  page.value = 1
  try {
    const r = await SearchQueries(searchQ.value.trim(), page.value)
    queries.value = r?.matches || []
    total.value = r?.total || 0
  } catch (e: any) {
    error.value = e?.toString()
  } finally {
    loading.value = false
  }
}

async function loadTags() {
  try {
    const r = await QueryTags(20)
    popularTags.value = r?.matches || []
  } catch {}
}

function useQuery(q: string) {
  router.push({ path: '/search', query: { q } })
}

async function prevPage() {
  if (page.value <= 1) return
  page.value--
  isSearchMode.value ? await searchQueries() : await loadQueries()
}

async function nextPage() {
  page.value++
  isSearchMode.value ? await searchQueries() : await loadQueries()
}

function formatDate(d: string) {
  if (!d) return '—'
  return new Date(d).toLocaleDateString()
}
</script>

<style scoped>
.query-card { transition: border-color 0.15s; }
.query-card:hover { border-color: var(--border); }
</style>
