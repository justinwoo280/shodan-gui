export namespace shodan {
	
	export class APIInfo {
	    scan_credits: number;
	    usage_limits: Record<string, any>;
	    plan: string;
	    https: boolean;
	    unlocked: boolean;
	    query_credits: number;
	    monitored_ips: number;
	    unlocked_left: number;
	    telnet: boolean;
	
	    static createFrom(source: any = {}) {
	        return new APIInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.scan_credits = source["scan_credits"];
	        this.usage_limits = source["usage_limits"];
	        this.plan = source["plan"];
	        this.https = source["https"];
	        this.unlocked = source["unlocked"];
	        this.query_credits = source["query_credits"];
	        this.monitored_ips = source["monitored_ips"];
	        this.unlocked_left = source["unlocked_left"];
	        this.telnet = source["telnet"];
	    }
	}
	export class AlertFilters {
	    ip: string[];
	
	    static createFrom(source: any = {}) {
	        return new AlertFilters(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ip = source["ip"];
	    }
	}
	export class Alert {
	    name: string;
	    created: string;
	    triggers: Record<string, any>;
	    has_triggers: boolean;
	    expires: number;
	    expiration: any;
	    filters: AlertFilters;
	    id: string;
	    size: number;
	    notifiers: string[];
	
	    static createFrom(source: any = {}) {
	        return new Alert(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.created = source["created"];
	        this.triggers = source["triggers"];
	        this.has_triggers = source["has_triggers"];
	        this.expires = source["expires"];
	        this.expiration = source["expiration"];
	        this.filters = this.convertValues(source["filters"], AlertFilters);
	        this.id = source["id"];
	        this.size = source["size"];
	        this.notifiers = source["notifiers"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class Location {
	    city: any;
	    region_code: any;
	    area_code: any;
	    longitude: number;
	    country_code: string;
	    country_name: string;
	    postal_code: any;
	    dma_code: any;
	    latitude: number;
	
	    static createFrom(source: any = {}) {
	        return new Location(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.city = source["city"];
	        this.region_code = source["region_code"];
	        this.area_code = source["area_code"];
	        this.longitude = source["longitude"];
	        this.country_code = source["country_code"];
	        this.country_name = source["country_name"];
	        this.postal_code = source["postal_code"];
	        this.dma_code = source["dma_code"];
	        this.latitude = source["latitude"];
	    }
	}
	export class Banner {
	    ip: any;
	    ip_str: string;
	    port: number;
	    transport: string;
	    _shodan: Record<string, any>;
	    hostnames: string[];
	    domains: string[];
	    location: Location;
	    org: string;
	    isp: string;
	    asn: string;
	    os: any;
	    data: string;
	    timestamp: string;
	    product: string;
	    version: string;
	    cpe: string[];
	    tags: string[];
	    http: Record<string, any>;
	    ssl: Record<string, any>;
	    opts: Record<string, any>;
	    hash: number;
	
	    static createFrom(source: any = {}) {
	        return new Banner(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ip = source["ip"];
	        this.ip_str = source["ip_str"];
	        this.port = source["port"];
	        this.transport = source["transport"];
	        this._shodan = source["_shodan"];
	        this.hostnames = source["hostnames"];
	        this.domains = source["domains"];
	        this.location = this.convertValues(source["location"], Location);
	        this.org = source["org"];
	        this.isp = source["isp"];
	        this.asn = source["asn"];
	        this.os = source["os"];
	        this.data = source["data"];
	        this.timestamp = source["timestamp"];
	        this.product = source["product"];
	        this.version = source["version"];
	        this.cpe = source["cpe"];
	        this.tags = source["tags"];
	        this.http = source["http"];
	        this.ssl = source["ssl"];
	        this.opts = source["opts"];
	        this.hash = source["hash"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DNSRecord {
	    subdomain: string;
	    type: string;
	    value: string;
	    last_seen: string;
	
	    static createFrom(source: any = {}) {
	        return new DNSRecord(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.subdomain = source["subdomain"];
	        this.type = source["type"];
	        this.value = source["value"];
	        this.last_seen = source["last_seen"];
	    }
	}
	export class DomainResponse {
	    domain: string;
	    tags: string[];
	    data: DNSRecord[];
	    subdomains: string[];
	    more: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DomainResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.domain = source["domain"];
	        this.tags = source["tags"];
	        this.data = this.convertValues(source["data"], DNSRecord);
	        this.subdomains = source["subdomains"];
	        this.more = source["more"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class FacetValue {
	    count: number;
	    value: string;
	
	    static createFrom(source: any = {}) {
	        return new FacetValue(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.count = source["count"];
	        this.value = source["value"];
	    }
	}
	export class HostResponse {
	    region_code: any;
	    ip: any;
	    postal_code: any;
	    country_code: string;
	    city: any;
	    dma_code: any;
	    last_update: string;
	    latitude: number;
	    tags: string[];
	    area_code: any;
	    country_name: string;
	    hostnames: string[];
	    org: string;
	    data: Banner[];
	    asn: string;
	    isp: string;
	    longitude: number;
	    domains: string[];
	    ip_str: string;
	    os: any;
	    ports: number[];
	
	    static createFrom(source: any = {}) {
	        return new HostResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.region_code = source["region_code"];
	        this.ip = source["ip"];
	        this.postal_code = source["postal_code"];
	        this.country_code = source["country_code"];
	        this.city = source["city"];
	        this.dma_code = source["dma_code"];
	        this.last_update = source["last_update"];
	        this.latitude = source["latitude"];
	        this.tags = source["tags"];
	        this.area_code = source["area_code"];
	        this.country_name = source["country_name"];
	        this.hostnames = source["hostnames"];
	        this.org = source["org"];
	        this.data = this.convertValues(source["data"], Banner);
	        this.asn = source["asn"];
	        this.isp = source["isp"];
	        this.longitude = source["longitude"];
	        this.domains = source["domains"];
	        this.ip_str = source["ip_str"];
	        this.os = source["os"];
	        this.ports = source["ports"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class Notifier {
	    description: any;
	    args: Record<string, any>;
	    provider: string;
	    id: string;
	
	    static createFrom(source: any = {}) {
	        return new Notifier(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.description = source["description"];
	        this.args = source["args"];
	        this.provider = source["provider"];
	        this.id = source["id"];
	    }
	}
	export class NotifiersResponse {
	    matches: Notifier[];
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new NotifiersResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.matches = this.convertValues(source["matches"], Notifier);
	        this.total = source["total"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Query {
	    votes: number;
	    description: string;
	    tags: string[];
	    timestamp: string;
	    title: string;
	    query: string;
	
	    static createFrom(source: any = {}) {
	        return new Query(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.votes = source["votes"];
	        this.description = source["description"];
	        this.tags = source["tags"];
	        this.timestamp = source["timestamp"];
	        this.title = source["title"];
	        this.query = source["query"];
	    }
	}
	export class QueriesResponse {
	    matches: Query[];
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new QueriesResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.matches = this.convertValues(source["matches"], Query);
	        this.total = source["total"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class ScanItem {
	    status: string;
	    created: string;
	    status_check: string;
	    credits_left: number;
	    id: string;
	    size: number;
	
	    static createFrom(source: any = {}) {
	        return new ScanItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.created = source["created"];
	        this.status_check = source["status_check"];
	        this.credits_left = source["credits_left"];
	        this.id = source["id"];
	        this.size = source["size"];
	    }
	}
	export class ScanResponse {
	    count: number;
	    id: string;
	    credits_left: number;
	
	    static createFrom(source: any = {}) {
	        return new ScanResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.count = source["count"];
	        this.id = source["id"];
	        this.credits_left = source["credits_left"];
	    }
	}
	export class ScanStatus {
	    count: number;
	    status: string;
	    id: string;
	    created: string;
	
	    static createFrom(source: any = {}) {
	        return new ScanStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.count = source["count"];
	        this.status = source["status"];
	        this.id = source["id"];
	        this.created = source["created"];
	    }
	}
	export class ScansListResponse {
	    matches: ScanItem[];
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new ScansListResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.matches = this.convertValues(source["matches"], ScanItem);
	        this.total = source["total"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SearchCountResponse {
	    total: number;
	    facets: Record<string, Array<FacetValue>>;
	
	    static createFrom(source: any = {}) {
	        return new SearchCountResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total = source["total"];
	        this.facets = this.convertValues(source["facets"], Array<FacetValue>, true);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SearchResponse {
	    total: number;
	    matches: Banner[];
	    facets: Record<string, Array<FacetValue>>;
	    _scroll_id: string;
	
	    static createFrom(source: any = {}) {
	        return new SearchResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total = source["total"];
	        this.matches = this.convertValues(source["matches"], Banner);
	        this.facets = this.convertValues(source["facets"], Array<FacetValue>, true);
	        this._scroll_id = source["_scroll_id"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TagResult {
	    count: number;
	    value: string;
	
	    static createFrom(source: any = {}) {
	        return new TagResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.count = source["count"];
	        this.value = source["value"];
	    }
	}
	export class TagsResponse {
	    matches: TagResult[];
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new TagsResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.matches = this.convertValues(source["matches"], TagResult);
	        this.total = source["total"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TokensResponse {
	    attributes: Record<string, any>;
	    errors: string[];
	    string: string;
	    filters: string[];
	
	    static createFrom(source: any = {}) {
	        return new TokensResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.attributes = source["attributes"];
	        this.errors = source["errors"];
	        this.string = source["string"];
	        this.filters = source["filters"];
	    }
	}
	export class Trigger {
	    name: string;
	    rule: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new Trigger(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.rule = source["rule"];
	        this.description = source["description"];
	    }
	}

}

