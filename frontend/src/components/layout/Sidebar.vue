<template>
  <nav class="sidebar">
    <div class="sidebar-logo">
      <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
        <circle cx="12" cy="12" r="10" stroke="#58a6ff" stroke-width="1.5"/>
        <circle cx="12" cy="12" r="6" stroke="#58a6ff" stroke-width="1.5" stroke-dasharray="3 2"/>
        <circle cx="12" cy="12" r="2" fill="#58a6ff"/>
      </svg>
      <span class="logo-text">Shodan GUI</span>
    </div>

    <div class="nav-section">
      <div class="nav-section-label">Search</div>
      <router-link to="/host" class="nav-item">
        <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="16"/><line x1="8" y1="12" x2="16" y2="12"/></svg>
        Host Lookup
      </router-link>
      <router-link to="/search" class="nav-item">
        <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
        Search
      </router-link>
      <router-link to="/dns" class="nav-item">
        <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="2" width="20" height="8" rx="2"/><rect x="2" y="14" width="20" height="8" rx="2"/><line x1="6" y1="6" x2="6" y2="6"/><line x1="6" y1="18" x2="6" y2="18"/></svg>
        DNS Tools
      </router-link>
    </div>

    <div class="nav-section">
      <div class="nav-section-label">Operations</div>
      <router-link to="/scans" class="nav-item">
        <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg>
        Scans
      </router-link>
      <router-link to="/alerts" class="nav-item">
        <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"/><path d="M13.73 21a2 2 0 0 1-3.46 0"/></svg>
        Alerts
      </router-link>
      <router-link to="/notifiers" class="nav-item">
        <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/><polyline points="22,6 12,13 2,6"/></svg>
        Notifiers
      </router-link>
    </div>

    <div class="sidebar-spacer"></div>

    <div class="nav-section">
      <router-link to="/settings" class="nav-item">
        <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 2.83-2.83l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>
        Settings
      </router-link>
    </div>

    <div class="api-status" v-if="store.apiInfo">
      <div class="api-status-row">
        <span class="text-muted">Query Credits</span>
        <span class="mono text-accent">{{ store.apiInfo.query_credits?.toLocaleString() }}</span>
      </div>
      <div class="api-status-row">
        <span class="text-muted">Scan Credits</span>
        <span class="mono text-success">{{ store.apiInfo.scan_credits?.toLocaleString() }}</span>
      </div>
      <div class="api-status-row">
        <span class="text-muted">Plan</span>
        <span class="badge badge-blue">{{ store.apiInfo.plan }}</span>
      </div>
    </div>
    <div class="api-status no-key" v-else-if="!store.hasKey">
      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
      No API Key
    </div>
  </nav>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useSettingsStore } from '../../stores/settings'
const store = useSettingsStore()
onMounted(() => {
  if (store.hasKey) store.loadApiInfo()
})
</script>

<style scoped>
.sidebar {
  width: var(--sidebar-width);
  background: var(--bg-secondary);
  border-right: 1px solid var(--border-subtle);
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow-y: auto;
  flex-shrink: 0;
}

.sidebar-logo {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 16px;
  border-bottom: 1px solid var(--border-subtle);
  font-weight: 700;
  font-size: 15px;
  color: var(--text-primary);
  --wails-draggable: drag;
}
.logo-text { color: var(--accent); letter-spacing: 0.02em; }

.nav-section {
  padding: 12px 8px 4px;
}
.nav-section-label {
  font-size: 10px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: var(--text-muted);
  padding: 0 8px;
  margin-bottom: 4px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 7px 10px;
  border-radius: var(--radius-sm);
  color: var(--text-secondary);
  text-decoration: none;
  font-size: 13px;
  transition: all 0.15s;
  margin-bottom: 2px;
}
.nav-item:hover {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}
.nav-item.router-link-active {
  background: var(--accent-dim);
  color: var(--accent);
}

.sidebar-spacer { flex: 1; }

.api-status {
  margin: 8px;
  padding: 10px 12px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius);
  font-size: 12px;
}
.api-status.no-key {
  display: flex;
  align-items: center;
  gap: 6px;
  color: var(--text-muted);
}
.api-status-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}
.api-status-row:last-child { margin-bottom: 0; }
</style>
