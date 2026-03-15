package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"shodan-gui/shodan"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type Config struct {
	APIKey string `json:"api_key"`
}

type App struct {
	ctx        context.Context
	client     *shodan.Client
	configPath string
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}
	dir := filepath.Join(configDir, "shodan-gui")
	os.MkdirAll(dir, 0755)
	a.configPath = filepath.Join(dir, "config.json")

	cfg := a.loadConfig()
	if cfg.APIKey != "" {
		a.client = shodan.NewClient(cfg.APIKey)
	}
}

func (a *App) loadConfig() Config {
	data, err := os.ReadFile(a.configPath)
	if err != nil {
		return Config{}
	}
	var cfg Config
	json.Unmarshal(data, &cfg)
	return cfg
}

func (a *App) saveConfig(cfg Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(a.configPath, data, 0600)
}

func (a *App) checkClient() error {
	if a.client == nil {
		return fmt.Errorf("API key not set. Please configure your Shodan API key in Settings.")
	}
	return nil
}

func (a *App) SetAPIKey(key string) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return fmt.Errorf("API key cannot be empty")
	}
	a.client = shodan.NewClient(key)
	cfg := Config{APIKey: key}
	return a.saveConfig(cfg)
}

func (a *App) GetAPIKey() string {
	cfg := a.loadConfig()
	return cfg.APIKey
}

func (a *App) HasAPIKey() bool {
	return a.client != nil
}

func (a *App) HostInfo(ip string, history bool, minify bool) (*shodan.HostResponse, error) {
	if err := a.checkClient(); err != nil {
		return nil, err
	}
	return a.client.HostInfo(ip, history, minify)
}

func (a *App) SearchCount(query, facets string) (*shodan.SearchCountResponse, error) {
	if err := a.checkClient(); err != nil {
		return nil, err
	}
	return a.client.SearchCount(query, facets)
}

func (a *App) Search(query, facets string, page int, minify bool) (*shodan.SearchResponse, error) {
	if err := a.checkClient(); err != nil {
		return nil, err
	}
	return a.client.Search(query, facets, page, minify)
}

func (a *App) SearchAll(query, facets string, maxPages int, filterHoneypots bool) ([]shodan.Banner, error) {
	if err := a.checkClient(); err != nil {
		return nil, err
	}
	if maxPages <= 0 {
		maxPages = 10
	}
	if maxPages > 100 {
		maxPages = 100
	}
	var all []shodan.Banner
	for page := 1; page <= maxPages; page++ {
		res, err := a.client.Search(query, facets, page, true)
		if err != nil {
			return nil, err
		}
		for _, m := range res.Matches {
			if filterHoneypots {
				isHoneypot := false
				for _, t := range m.Tags {
					if t == "honeypot" {
						isHoneypot = true
						break
					}
				}
				if isHoneypot {
					continue
				}
			}
			all = append(all, m)
		}
		if len(res.Matches) < 100 || (res.Total > 0 && len(all) >= res.Total) {
			break
		}
	}
	return all, nil
}

func (a *App) ListFacets() ([]string, error) {
	if err := a.checkClient(); err != nil {
		return nil, err
	}
	return a.client.ListFacets()
}

func (a *App) ListFilters() ([]string, error) {
	if err := a.checkClient(); err != nil {
		return nil, err
	}
	return a.client.ListFilters()
}

func (a *App) SearchTokens(query string) (*shodan.TokensResponse, error) {
	if err := a.checkClient(); err != nil {
		return nil, err
	}
	return a.client.SearchTokens(query)
}

func (a *App) ListPorts() ([]int, error) {
	if err := a.checkClient(); err != nil {
		return nil, err
	}
	return a.client.ListPorts()
}

func (a *App) ListProtocols() (map[string]string, error) {
	if err := a.checkClient(); err != nil {
		return nil, err
	}
	return a.client.ListProtocols()
}

func (a *App) Scan(ips string) (*shodan.ScanResponse, error) {
	if err := a.checkClient(); err != nil {
		return nil, err
	}
	return a.client.Scan(ips)
}

func (a *App) ListScans() (*shodan.ScansListResponse, error) {
	if err := a.checkClient(); err != nil {
		return nil, err
	}
	return a.client.ListScans()
}

func (a *App) ScanStatus(id string) (*shodan.ScanStatus, error) {
	if err := a.checkClient(); err != nil {
		return nil, err
	}
	return a.client.ScanStatus(id)
}

func (a *App) CreateAlert(name string, ips []string, expires int) (*shodan.Alert, error) {
	if err := a.checkClient(); err != nil {
		return nil, err
	}
	return a.client.CreateAlert(name, ips, expires)
}

func (a *App) AlertInfo(id string) (*shodan.Alert, error) {
	if err := a.checkClient(); err != nil {
		return nil, err
	}
	return a.client.AlertInfo(id)
}

func (a *App) DeleteAlert(id string) error {
	if err := a.checkClient(); err != nil {
		return err
	}
	return a.client.DeleteAlert(id)
}

func (a *App) EditAlert(id string, ips []string) (*shodan.Alert, error) {
	if err := a.checkClient(); err != nil {
		return nil, err
	}
	return a.client.EditAlert(id, ips)
}

func (a *App) ListAlerts() ([]shodan.Alert, error) {
	if err := a.checkClient(); err != nil {
		return nil, err
	}
	return a.client.ListAlerts()
}

func (a *App) ListAlertTriggers() ([]shodan.Trigger, error) {
	if err := a.checkClient(); err != nil {
		return nil, err
	}
	return a.client.ListAlertTriggers()
}

func (a *App) EnableTrigger(alertID, trigger string) error {
	if err := a.checkClient(); err != nil {
		return err
	}
	return a.client.EnableTrigger(alertID, trigger)
}

func (a *App) DisableTrigger(alertID, trigger string) error {
	if err := a.checkClient(); err != nil {
		return err
	}
	return a.client.DisableTrigger(alertID, trigger)
}

func (a *App) ListNotifiers() (*shodan.NotifiersResponse, error) {
	if err := a.checkClient(); err != nil {
		return nil, err
	}
	return a.client.ListNotifiers()
}

func (a *App) ListProviders() (map[string]interface{}, error) {
	if err := a.checkClient(); err != nil {
		return nil, err
	}
	return a.client.ListProviders()
}

func (a *App) DeleteNotifier(id string) error {
	if err := a.checkClient(); err != nil {
		return err
	}
	return a.client.DeleteNotifier(id)
}

func (a *App) GetNotifier(id string) (*shodan.Notifier, error) {
	if err := a.checkClient(); err != nil {
		return nil, err
	}
	return a.client.GetNotifier(id)
}

func (a *App) ListQueries(page int, sort, order string) (*shodan.QueriesResponse, error) {
	if err := a.checkClient(); err != nil {
		return nil, err
	}
	return a.client.ListQueries(page, sort, order)
}

func (a *App) SearchQueries(query string, page int) (*shodan.QueriesResponse, error) {
	if err := a.checkClient(); err != nil {
		return nil, err
	}
	return a.client.SearchQueries(query, page)
}

func (a *App) QueryTags(size int) (*shodan.TagsResponse, error) {
	if err := a.checkClient(); err != nil {
		return nil, err
	}
	return a.client.QueryTags(size)
}

func (a *App) DomainInfo(domain, dnsType string, history bool, page int) (*shodan.DomainResponse, error) {
	if err := a.checkClient(); err != nil {
		return nil, err
	}
	return a.client.DomainInfo(domain, dnsType, history, page)
}

func (a *App) DNSResolve(hostnames string) (map[string]string, error) {
	if err := a.checkClient(); err != nil {
		return nil, err
	}
	return a.client.DNSResolve(hostnames)
}

func (a *App) DNSReverse(ips string) (map[string][]string, error) {
	if err := a.checkClient(); err != nil {
		return nil, err
	}
	return a.client.DNSReverse(ips)
}

func (a *App) MyIP() (string, error) {
	if err := a.checkClient(); err != nil {
		return "", err
	}
	return a.client.MyIP()
}

func (a *App) APIInfo() (*shodan.APIInfo, error) {
	if err := a.checkClient(); err != nil {
		return nil, err
	}
	return a.client.APIInfo()
}

func (a *App) ExportJSON(data string, defaultName string) (string, error) {
	path, err := wailsRuntime.SaveFileDialog(a.ctx, wailsRuntime.SaveDialogOptions{
		DefaultFilename: defaultName + "_" + time.Now().Format("20060102_150405") + ".json",
		Filters: []wailsRuntime.FileFilter{
			{DisplayName: "JSON Files", Pattern: "*.json"},
		},
	})
	if err != nil || path == "" {
		return "", err
	}
	var prettyJSON map[string]interface{}
	var prettyJSONArr []interface{}
	var out []byte
	if json.Unmarshal([]byte(data), &prettyJSON) == nil {
		out, _ = json.MarshalIndent(prettyJSON, "", "  ")
	} else if json.Unmarshal([]byte(data), &prettyJSONArr) == nil {
		out, _ = json.MarshalIndent(prettyJSONArr, "", "  ")
	} else {
		out = []byte(data)
	}
	err = os.WriteFile(path, out, 0644)
	if err != nil {
		return "", err
	}
	return path, nil
}

func (a *App) ExportCSV(jsonData string, defaultName string) (string, error) {
	path, err := wailsRuntime.SaveFileDialog(a.ctx, wailsRuntime.SaveDialogOptions{
		DefaultFilename: defaultName + "_" + time.Now().Format("20060102_150405") + ".csv",
		Filters: []wailsRuntime.FileFilter{
			{DisplayName: "CSV Files", Pattern: "*.csv"},
		},
	})
	if err != nil || path == "" {
		return "", err
	}

	var rows []map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &rows); err != nil {
		return "", fmt.Errorf("data must be a JSON array of objects")
	}

	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	if len(rows) == 0 {
		return path, nil
	}

	headers := make([]string, 0)
	for k := range rows[0] {
		headers = append(headers, k)
	}
	w.Write(headers)

	for _, row := range rows {
		record := make([]string, len(headers))
		for i, h := range headers {
			val := row[h]
			if val == nil {
				record[i] = ""
			} else {
				b, _ := json.Marshal(val)
				s := string(b)
				s = strings.Trim(s, `"`)
				record[i] = s
			}
		}
		w.Write(record)
	}

	return path, nil
}
