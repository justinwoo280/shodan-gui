<template>
  <div class="page-container">
    <div class="page-header">
      <div class="page-title">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#58a6ff" stroke-width="2"><rect x="2" y="2" width="20" height="8" rx="2"/><rect x="2" y="14" width="20" height="8" rx="2"/></svg>
        DNS Tools
      </div>
      <div class="page-subtitle">Domain info, forward & reverse DNS resolution</div>
    </div>

    <div style="display:flex;flex:1;overflow:hidden">
      <div class="page-content" style="flex:1;overflow-y:auto">
        <div class="tabs">
          <span class="tab" :class="{ active: tab === 'domain' }" @click="tab = 'domain'">Domain Info</span>
          <span class="tab" :class="{ active: tab === 'resolve' }" @click="tab = 'resolve'">DNS Resolve</span>
          <span class="tab" :class="{ active: tab === 'reverse' }" @click="tab = 'reverse'">Reverse DNS</span>
          <span class="tab" :class="{ active: tab === 'myip' }" @click="tab = 'myip'; getMyIP()">My IP</span>
        </div>

        <div class="alert-box error" v-if="error">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
          {{ error }}
        </div>

        <div v-if="tab === 'domain'">
          <div class="flex gap-2" style="margin-bottom:16px;flex-wrap:wrap">
            <input class="input flex-1" placeholder="Domain (e.g. google.com)" v-model="domain" @keyup.enter="lookupDomain" style="max-width:300px" />
            <select class="input select" v-model="dnsType" style="width:120px">
              <option value="">All Types</option>
              <option value="A">A</option>
              <option value="AAAA">AAAA</option>
              <option value="CNAME">CNAME</option>
              <option value="NS">NS</option>
              <option value="MX">MX</option>
              <option value="TXT">TXT</option>
              <option value="SOA">SOA</option>
            </select>
            <label class="flex items-center gap-2" style="font-size:13px;color:var(--text-secondary);cursor:pointer">
              <input type="checkbox" class="checkbox" v-model="dnsHistory" /> History
            </label>
            <button class="btn btn-primary" @click="lookupDomain" :disabled="loading || !domain.trim()">
              <div class="spinner" v-if="loading" style="width:12px;height:12px"></div>
              Lookup
            </button>
            <button class="btn btn-sm" v-if="domainResult" @click="exportDomain">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
              Export
            </button>
          </div>

          <div v-if="loading" class="loading-overlay"><div class="spinner"></div> Looking up domain...</div>

          <div v-if="domainResult && !loading">
            <div class="card" style="margin-bottom:16px">
              <div class="flex items-center gap-3">
                <div style="font-size:20px;font-weight:700;color:var(--accent)">{{ domainResult.domain }}</div>
                <span class="badge badge-blue" v-for="t in domainResult.tags" :key="t">{{ t }}</span>
                <span class="badge badge-yellow" v-if="domainResult.more">More data available</span>
              </div>
              <div style="margin-top:8px">
                <span class="text-muted" style="font-size:12px">{{ domainResult.subdomains?.length || 0 }} subdomains found</span>
              </div>
            </div>

            <div class="card" style="padding:0;overflow:hidden;margin-bottom:16px" v-if="domainResult.data?.length">
              <table class="table">
                <thead>
                  <tr>
                    <th>Subdomain</th>
                    <th>Type</th>
                    <th>Value</th>
                    <th>Last Seen</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="(rec, i) in domainResult.data" :key="i">
                    <td class="mono" style="color:var(--accent)">{{ rec.subdomain || '@' }}</td>
                    <td><span class="badge" :class="typeClass(rec.type)">{{ rec.type }}</span></td>
                    <td class="mono" style="font-size:12px">{{ rec.value }}</td>
                    <td class="text-muted" style="font-size:11px">{{ formatDate(rec.last_seen) }}</td>
                  </tr>
                </tbody>
              </table>
            </div>

            <div v-if="domainResult.subdomains?.length">
              <div class="section-title">Subdomains ({{ domainResult.subdomains.length }})</div>
              <div class="flex gap-2" style="flex-wrap:wrap">
                <span class="tag" v-for="s in domainResult.subdomains" :key="s">{{ s }}</span>
              </div>
            </div>
          </div>
        </div>

        <div v-if="tab === 'resolve'">
          <div class="flex gap-2" style="margin-bottom:16px">
            <input class="input flex-1" placeholder="Hostnames (comma-separated): google.com,bing.com" v-model="resolveInput" @keyup.enter="doResolve" />
            <button class="btn btn-primary" @click="doResolve" :disabled="loading || !resolveInput.trim()">
              <div class="spinner" v-if="loading" style="width:12px;height:12px"></div>
              Resolve
            </button>
          </div>
          <div v-if="loading" class="loading-overlay"><div class="spinner"></div></div>
          <div v-if="resolveResult && !loading" class="card" style="padding:0;overflow:hidden">
            <table class="table">
              <thead><tr><th>Hostname</th><th>IP Address</th></tr></thead>
              <tbody>
                <tr v-for="(ip, hostname) in resolveResult" :key="String(hostname)">
                  <td class="mono" style="color:var(--text-secondary)">{{ hostname }}</td>
                  <td><span class="ip-chip">{{ ip }}</span></td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div v-if="tab === 'reverse'">
          <div class="flex gap-2" style="margin-bottom:16px">
            <input class="input flex-1" placeholder="IP addresses (comma-separated): 8.8.8.8,1.1.1.1" v-model="reverseInput" @keyup.enter="doReverse" />
            <button class="btn btn-primary" @click="doReverse" :disabled="loading || !reverseInput.trim()">
              <div class="spinner" v-if="loading" style="width:12px;height:12px"></div>
              Reverse
            </button>
          </div>
          <div v-if="loading" class="loading-overlay"><div class="spinner"></div></div>
          <div v-if="reverseResult && !loading" class="card" style="padding:0;overflow:hidden">
            <table class="table">
              <thead><tr><th>IP Address</th><th>Hostnames</th></tr></thead>
              <tbody>
                <tr v-for="(hostnames, ip) in reverseResult" :key="String(ip)">
                  <td><span class="ip-chip">{{ ip }}</span></td>
                  <td>
                    <div class="flex gap-2" style="flex-wrap:wrap">
                      <span class="tag" v-for="h in hostnames" :key="h">{{ h }}</span>
                      <span class="text-muted" v-if="!hostnames?.length">No PTR records</span>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div v-if="tab === 'myip'">
          <div v-if="loading" class="loading-overlay"><div class="spinner"></div></div>
          <div v-if="myIP && !loading" class="card" style="text-align:center;padding:40px">
            <div class="text-muted" style="font-size:13px;margin-bottom:8px">Your Public IP Address</div>
            <div class="ip-chip" style="font-size:24px;padding:12px 24px;display:inline-flex">{{ myIP }}</div>
            <div class="text-muted" style="font-size:12px;margin-top:12px">As seen from the Internet via Shodan</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { DomainInfo, DNSResolve, DNSReverse, MyIP, ExportJSON } from '../../wailsjs/go/main/App'

const tab = ref<'domain' | 'resolve' | 'reverse' | 'myip'>('domain')
const loading = ref(false)
const error = ref('')

const domain = ref('')
const dnsType = ref('')
const dnsHistory = ref(false)
const domainResult = ref<any>(null)

const resolveInput = ref('')
const resolveResult = ref<any>(null)

const reverseInput = ref('')
const reverseResult = ref<any>(null)

const myIP = ref('')

async function lookupDomain() {
  if (!domain.value.trim()) return
  loading.value = true
  error.value = ''
  domainResult.value = null
  try {
    domainResult.value = await DomainInfo(domain.value.trim(), dnsType.value, dnsHistory.value, 1)
  } catch (e: any) {
    error.value = e?.toString()
  } finally {
    loading.value = false
  }
}

async function doResolve() {
  if (!resolveInput.value.trim()) return
  loading.value = true
  error.value = ''
  resolveResult.value = null
  try {
    resolveResult.value = await DNSResolve(resolveInput.value.trim())
  } catch (e: any) {
    error.value = e?.toString()
  } finally {
    loading.value = false
  }
}

async function doReverse() {
  if (!reverseInput.value.trim()) return
  loading.value = true
  error.value = ''
  reverseResult.value = null
  try {
    reverseResult.value = await DNSReverse(reverseInput.value.trim())
  } catch (e: any) {
    error.value = e?.toString()
  } finally {
    loading.value = false
  }
}

async function getMyIP() {
  if (myIP.value) return
  loading.value = true
  error.value = ''
  try {
    myIP.value = await MyIP()
  } catch (e: any) {
    error.value = e?.toString()
  } finally {
    loading.value = false
  }
}

async function exportDomain() {
  if (!domainResult.value) return
  try {
    await ExportJSON(JSON.stringify(domainResult.value), `dns_${domain.value}`)
  } catch (e: any) {
    error.value = e?.toString()
  }
}

function typeClass(t: string) {
  const map: Record<string, string> = {
    A: 'badge-blue', AAAA: 'badge-purple', CNAME: 'badge-yellow',
    MX: 'badge-green', NS: 'badge-blue', TXT: 'badge-yellow', SOA: 'badge-red'
  }
  return map[t] || 'badge-blue'
}

function formatDate(d: string) {
  if (!d) return '—'
  return new Date(d).toLocaleString()
}
</script>
