<template>
  <div class="json-viewer" v-if="data !== null && data !== undefined">
    <div class="json-toolbar">
      <span class="json-title">{{ title }}</span>
      <div class="flex gap-2">
        <button class="btn btn-sm" @click="copyJSON">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg>
          {{ copied ? 'Copied!' : 'Copy' }}
        </button>
        <button class="btn btn-sm" @click="$emit('export')" v-if="showExport">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
          Export
        </button>
        <button class="btn btn-sm" @click="collapsed = !collapsed">
          {{ collapsed ? 'Expand' : 'Collapse' }}
        </button>
      </div>
    </div>
    <pre class="json-content" v-show="!collapsed">{{ formatted }}</pre>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

const props = defineProps<{
  data: any
  title?: string
  showExport?: boolean
}>()
defineEmits(['export'])

const collapsed = ref(false)
const copied = ref(false)

const formatted = computed(() => {
  try {
    return JSON.stringify(props.data, null, 2)
  } catch {
    return String(props.data)
  }
})

async function copyJSON() {
  await navigator.clipboard.writeText(formatted.value)
  copied.value = true
  setTimeout(() => copied.value = false, 2000)
}
</script>

<style scoped>
.json-viewer {
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius);
  overflow: hidden;
  background: var(--bg-primary);
}
.json-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  background: var(--bg-tertiary);
  border-bottom: 1px solid var(--border-subtle);
}
.json-title {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}
.json-content {
  padding: 12px;
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-secondary);
  overflow: auto;
  max-height: 500px;
  white-space: pre;
  line-height: 1.6;
}
</style>
