package shodan

type Location struct {
	City        interface{} `json:"city"`
	RegionCode  interface{} `json:"region_code"`
	AreaCode    interface{} `json:"area_code"`
	Longitude   float64     `json:"longitude"`
	CountryCode string      `json:"country_code"`
	CountryName string      `json:"country_name"`
	PostalCode  interface{} `json:"postal_code"`
	DMACode     interface{} `json:"dma_code"`
	Latitude    float64     `json:"latitude"`
}

type Banner struct {
	IP        interface{}            `json:"ip"`
	IPStr     string                 `json:"ip_str"`
	Port      int                    `json:"port"`
	Transport string                 `json:"transport"`
	Protocol  map[string]interface{} `json:"_shodan"`
	Hostnames []string               `json:"hostnames"`
	Domains   []string               `json:"domains"`
	Location  Location               `json:"location"`
	Org       string                 `json:"org"`
	ISP       string                 `json:"isp"`
	ASN       string                 `json:"asn"`
	OS        interface{}            `json:"os"`
	Data      string                 `json:"data"`
	Timestamp string                 `json:"timestamp"`
	Product   string                 `json:"product"`
	Version   string                 `json:"version"`
	CPE       []string               `json:"cpe"`
	Tags      []string               `json:"tags"`
	HTTP      map[string]interface{} `json:"http"`
	SSL       map[string]interface{} `json:"ssl"`
	Opts      map[string]interface{} `json:"opts"`
	Hash      int                    `json:"hash"`
}

type HostResponse struct {
	RegionCode  interface{} `json:"region_code"`
	IP          interface{} `json:"ip"`
	PostalCode  interface{} `json:"postal_code"`
	CountryCode string      `json:"country_code"`
	City        interface{} `json:"city"`
	DMACode     interface{} `json:"dma_code"`
	LastUpdate  string      `json:"last_update"`
	Latitude    float64     `json:"latitude"`
	Tags        []string    `json:"tags"`
	AreaCode    interface{} `json:"area_code"`
	CountryName string      `json:"country_name"`
	Hostnames   []string    `json:"hostnames"`
	Org         string      `json:"org"`
	Data        []Banner    `json:"data"`
	ASN         string      `json:"asn"`
	ISP         string      `json:"isp"`
	Longitude   float64     `json:"longitude"`
	Domains     []string    `json:"domains"`
	IPStr       string      `json:"ip_str"`
	OS          interface{} `json:"os"`
	Ports       []int       `json:"ports"`
}

type FacetValue struct {
	Count int    `json:"count"`
	Value string `json:"value"`
}

type SearchCountResponse struct {
	Total  int                       `json:"total"`
	Facets map[string][]FacetValue   `json:"facets"`
}

type SearchResponse struct {
	Total    int                       `json:"total"`
	Matches  []Banner                  `json:"matches"`
	Facets   map[string][]FacetValue   `json:"facets"`
	ScrollID string                    `json:"_scroll_id"`
}

type TokensResponse struct {
	Attributes map[string]interface{} `json:"attributes"`
	Errors     []string               `json:"errors"`
	String     string                 `json:"string"`
	Filters    []string               `json:"filters"`
}

type ScanResponse struct {
	Count       int    `json:"count"`
	ID          string `json:"id"`
	CreditsLeft int    `json:"credits_left"`
}

type ScanStatus struct {
	Count   int    `json:"count"`
	Status  string `json:"status"`
	ID      string `json:"id"`
	Created string `json:"created"`
}

type ScanItem struct {
	Status      string `json:"status"`
	Created     string `json:"created"`
	StatusCheck string `json:"status_check"`
	CreditsLeft int    `json:"credits_left"`
	ID          string `json:"id"`
	Size        int    `json:"size"`
}

type ScansListResponse struct {
	Matches []ScanItem `json:"matches"`
	Total   int        `json:"total"`
}

type Alert struct {
	Name        string                 `json:"name"`
	Created     string                 `json:"created"`
	Triggers    map[string]interface{} `json:"triggers"`
	HasTriggers bool                   `json:"has_triggers"`
	Expires     int                    `json:"expires"`
	Expiration  interface{}            `json:"expiration"`
	Filters     AlertFilters           `json:"filters"`
	ID          string                 `json:"id"`
	Size        int                    `json:"size"`
	Notifiers   []string               `json:"notifiers"`
}

type AlertFilters struct {
	IP []string `json:"ip"`
}

type Trigger struct {
	Name        string `json:"name"`
	Rule        string `json:"rule"`
	Description string `json:"description"`
}

type Notifier struct {
	Description interface{}            `json:"description"`
	Args        map[string]interface{} `json:"args"`
	Provider    string                 `json:"provider"`
	ID          string                 `json:"id"`
}

type NotifiersResponse struct {
	Matches []Notifier `json:"matches"`
	Total   int        `json:"total"`
}

type Query struct {
	Votes       int      `json:"votes"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Timestamp   string   `json:"timestamp"`
	Title       string   `json:"title"`
	Query       string   `json:"query"`
}

type QueriesResponse struct {
	Matches []Query `json:"matches"`
	Total   int     `json:"total"`
}

type TagResult struct {
	Count int    `json:"count"`
	Value string `json:"value"`
}

type TagsResponse struct {
	Matches []TagResult `json:"matches"`
	Total   int         `json:"total"`
}

type DNSRecord struct {
	Subdomain string `json:"subdomain"`
	Type      string `json:"type"`
	Value     string `json:"value"`
	LastSeen  string `json:"last_seen"`
}

type DomainResponse struct {
	Domain     string      `json:"domain"`
	Tags       []string    `json:"tags"`
	Data       []DNSRecord `json:"data"`
	Subdomains []string    `json:"subdomains"`
	More       bool        `json:"more"`
}

type APIInfo struct {
	ScanCredits int                    `json:"scan_credits"`
	UsageLimits map[string]interface{} `json:"usage_limits"`
	Plan        string                 `json:"plan"`
	HTTPS       bool                   `json:"https"`
	Unlocked    bool                   `json:"unlocked"`
	QueryCredits int                   `json:"query_credits"`
	MonitoredIPs int                   `json:"monitored_ips"`
	UnlockedLeft int                   `json:"unlocked_left"`
	Telnet      bool                   `json:"telnet"`
}

type APIError struct {
	Error string `json:"error"`
}
