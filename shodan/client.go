package shodan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const BaseURL = "https://api.shodan.io"

type Client struct {
	APIKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		APIKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) buildURL(path string, params map[string]string) string {
	u, _ := url.Parse(BaseURL + path)
	q := u.Query()
	q.Set("key", c.APIKey)
	for k, v := range params {
		if v != "" {
			q.Set(k, v)
		}
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func (c *Client) get(path string, params map[string]string, result interface{}) error {
	resp, err := c.httpClient.Get(c.buildURL(path, params))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		var apiErr APIError
		if json.Unmarshal(body, &apiErr) == nil && apiErr.Error != "" {
			return fmt.Errorf("%s", apiErr.Error)
		}
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return json.Unmarshal(body, result)
}

func (c *Client) delete(path string) error {
	req, err := http.NewRequest("DELETE", c.buildURL(path, nil), nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		var apiErr APIError
		if json.Unmarshal(body, &apiErr) == nil && apiErr.Error != "" {
			return fmt.Errorf("%s", apiErr.Error)
		}
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) put(path string) error {
	req, err := http.NewRequest("PUT", c.buildURL(path, nil), nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		var apiErr APIError
		if json.Unmarshal(body, &apiErr) == nil && apiErr.Error != "" {
			return fmt.Errorf("%s", apiErr.Error)
		}
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) postJSON(path string, payload interface{}, result interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.buildURL(path, nil), bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		var apiErr APIError
		if json.Unmarshal(body, &apiErr) == nil && apiErr.Error != "" {
			return fmt.Errorf("%s", apiErr.Error)
		}
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	if result != nil {
		return json.Unmarshal(body, result)
	}
	return nil
}

func (c *Client) postForm(path string, fields map[string]string, result interface{}) error {
	form := url.Values{}
	for k, v := range fields {
		form.Set(k, v)
	}
	req, err := http.NewRequest("POST", c.buildURL(path, nil), strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		var apiErr APIError
		if json.Unmarshal(body, &apiErr) == nil && apiErr.Error != "" {
			return fmt.Errorf("%s", apiErr.Error)
		}
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	if result != nil {
		return json.Unmarshal(body, result)
	}
	return nil
}

func (c *Client) HostInfo(ip string, history bool, minify bool) (*HostResponse, error) {
	params := map[string]string{}
	if history {
		params["history"] = "true"
	}
	if minify {
		params["minify"] = "true"
	}
	var result HostResponse
	err := c.get("/shodan/host/"+ip, params, &result)
	return &result, err
}

func (c *Client) SearchCount(query, facets string) (*SearchCountResponse, error) {
	params := map[string]string{"query": query, "facets": facets}
	var result SearchCountResponse
	err := c.get("/shodan/host/count", params, &result)
	return &result, err
}

func (c *Client) Search(query, facets string, page int, minify bool) (*SearchResponse, error) {
	params := map[string]string{
		"query":  query,
		"facets": facets,
		"page":   fmt.Sprintf("%d", page),
	}
	if !minify {
		params["minify"] = "false"
	}
	var result SearchResponse
	err := c.get("/shodan/host/search", params, &result)
	return &result, err
}

func (c *Client) ListFacets() ([]string, error) {
	var result []string
	err := c.get("/shodan/host/search/facets", nil, &result)
	return result, err
}

func (c *Client) ListFilters() ([]string, error) {
	var result []string
	err := c.get("/shodan/host/search/filters", nil, &result)
	return result, err
}

func (c *Client) SearchTokens(query string) (*TokensResponse, error) {
	params := map[string]string{"query": query}
	var result TokensResponse
	err := c.get("/shodan/host/search/tokens", params, &result)
	return &result, err
}

func (c *Client) ListPorts() ([]int, error) {
	var result []int
	err := c.get("/shodan/ports", nil, &result)
	return result, err
}

func (c *Client) ListProtocols() (map[string]string, error) {
	var result map[string]string
	err := c.get("/shodan/protocols", nil, &result)
	return result, err
}

func (c *Client) Scan(ips string) (*ScanResponse, error) {
	var result ScanResponse
	err := c.postForm("/shodan/scan", map[string]string{"ips": ips}, &result)
	return &result, err
}

func (c *Client) ListScans() (*ScansListResponse, error) {
	var result ScansListResponse
	err := c.get("/shodan/scans", nil, &result)
	return &result, err
}

func (c *Client) ScanStatus(id string) (*ScanStatus, error) {
	var result ScanStatus
	err := c.get("/shodan/scan/"+id, nil, &result)
	return &result, err
}

func (c *Client) CreateAlert(name string, ips []string, expires int) (*Alert, error) {
	payload := map[string]interface{}{
		"name":    name,
		"filters": map[string]interface{}{"ip": ips},
		"expires": expires,
	}
	var result Alert
	err := c.postJSON("/shodan/alert", payload, &result)
	return &result, err
}

func (c *Client) AlertInfo(id string) (*Alert, error) {
	var result Alert
	err := c.get("/shodan/alert/"+id+"/info", nil, &result)
	return &result, err
}

func (c *Client) DeleteAlert(id string) error {
	return c.delete("/shodan/alert/" + id)
}

func (c *Client) EditAlert(id string, ips []string) (*Alert, error) {
	payload := map[string]interface{}{
		"filters": map[string]interface{}{"ip": ips},
	}
	var result Alert
	err := c.postJSON("/shodan/alert/"+id, payload, &result)
	return &result, err
}

func (c *Client) ListAlerts() ([]Alert, error) {
	var result []Alert
	err := c.get("/shodan/alert/info", nil, &result)
	return result, err
}

func (c *Client) ListAlertTriggers() ([]Trigger, error) {
	var result []Trigger
	err := c.get("/shodan/alert/triggers", nil, &result)
	return result, err
}

func (c *Client) EnableTrigger(alertID, trigger string) error {
	return c.put("/shodan/alert/" + alertID + "/trigger/" + trigger)
}

func (c *Client) DisableTrigger(alertID, trigger string) error {
	return c.delete("/shodan/alert/" + alertID + "/trigger/" + trigger)
}

func (c *Client) ListNotifiers() (*NotifiersResponse, error) {
	var result NotifiersResponse
	err := c.get("/notifier", nil, &result)
	return &result, err
}

func (c *Client) ListProviders() (map[string]interface{}, error) {
	var result map[string]interface{}
	err := c.get("/notifier/provider", nil, &result)
	return result, err
}

func (c *Client) DeleteNotifier(id string) error {
	return c.delete("/notifier/" + id)
}

func (c *Client) GetNotifier(id string) (*Notifier, error) {
	var result Notifier
	err := c.get("/notifier/"+id, nil, &result)
	return &result, err
}

func (c *Client) ListQueries(page int, sort, order string) (*QueriesResponse, error) {
	params := map[string]string{
		"page":  fmt.Sprintf("%d", page),
		"sort":  sort,
		"order": order,
	}
	var result QueriesResponse
	err := c.get("/shodan/query", params, &result)
	return &result, err
}

func (c *Client) SearchQueries(query string, page int) (*QueriesResponse, error) {
	params := map[string]string{
		"query": query,
		"page":  fmt.Sprintf("%d", page),
	}
	var result QueriesResponse
	err := c.get("/shodan/query/search", params, &result)
	return &result, err
}

func (c *Client) QueryTags(size int) (*TagsResponse, error) {
	params := map[string]string{"size": fmt.Sprintf("%d", size)}
	var result TagsResponse
	err := c.get("/shodan/query/tags", params, &result)
	return &result, err
}

func (c *Client) DomainInfo(domain, dnsType string, history bool, page int) (*DomainResponse, error) {
	params := map[string]string{
		"page": fmt.Sprintf("%d", page),
		"type": dnsType,
	}
	if history {
		params["history"] = "true"
	}
	var result DomainResponse
	err := c.get("/dns/domain/"+domain, params, &result)
	return &result, err
}

func (c *Client) DNSResolve(hostnames string) (map[string]string, error) {
	params := map[string]string{"hostnames": hostnames}
	var result map[string]string
	err := c.get("/dns/resolve", params, &result)
	return result, err
}

func (c *Client) DNSReverse(ips string) (map[string][]string, error) {
	params := map[string]string{"ips": ips}
	var result map[string][]string
	err := c.get("/dns/reverse", params, &result)
	return result, err
}

func (c *Client) MyIP() (string, error) {
	var result string
	err := c.get("/tools/myip", nil, &result)
	return result, err
}

func (c *Client) APIInfo() (*APIInfo, error) {
	var result APIInfo
	err := c.get("/api-info", nil, &result)
	return &result, err
}
