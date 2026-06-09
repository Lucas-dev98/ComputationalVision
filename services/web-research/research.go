package main

import (
	"context"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
)

type manufacturerRule struct {
	Canonical  string
	Aliases    []string
	PNPrefixes []string
}

var manufacturerRules = []manufacturerRule{
	{Canonical: "HPE", Aliases: []string{"HPE", "HEWLETT PACKARD ENTERPRISE"}},
	{Canonical: "HP", Aliases: []string{"HP", "HEWLETT PACKARD"}},
	{Canonical: "SK hynix", Aliases: []string{"SK HYNIX", "HYNIX"}, PNPrefixes: []string{"HMA", "HMT", "HMC", "HFS", "HFB"}},
	{Canonical: "Micron", Aliases: []string{"MICRON", "CRUCIAL"}, PNPrefixes: []string{"MTA", "MTF", "MTFD"}},
	{Canonical: "Samsung", Aliases: []string{"SAMSUNG"}, PNPrefixes: []string{"M3", "M4", "MZ", "PM", "SM"}},
	{Canonical: "Intel", Aliases: []string{"INTEL"}, PNPrefixes: []string{"INTEL_", "INTEL", "SSDPE", "SSDS", "SSDP"}},
	{Canonical: "Kingston", Aliases: []string{"KINGSTON"}, PNPrefixes: []string{"KSM", "KVR", "SKC", "SUV"}},
	{Canonical: "Seagate", Aliases: []string{"SEAGATE"}, PNPrefixes: []string{"SEAG"}},
	{Canonical: "HGST", Aliases: []string{"HGST", "HITACHI GST", "HITACHI GLOBAL STORAGE"}, PNPrefixes: []string{"HUH", "HUS", "HUC"}},
	{Canonical: "Cisco", Aliases: []string{"CISCO"}, PNPrefixes: []string{"UCSC"}},
	{Canonical: "Dell", Aliases: []string{"DELL"}},
	{Canonical: "Western Digital", Aliases: []string{"WESTERN DIGITAL", "WD", "WDC"}, PNPrefixes: []string{"WDS", "WDC", "WD"}},
	{Canonical: "Hisense", Aliases: []string{"HISENSE"}},
	{Canonical: "Broadcom", Aliases: []string{"BROADCOM"}},
	{Canonical: "Mellanox", Aliases: []string{"MELLANOX"}},
	{Canonical: "NVIDIA", Aliases: []string{"NVIDIA"}},
	{Canonical: "Lenovo", Aliases: []string{"LENOVO"}},
	{Canonical: "IBM", Aliases: []string{"IBM"}},
}

var seagatePartNumberPattern = regexp.MustCompile(`^ST\d`)

type Researcher struct {
	client *http.Client
}

func NewResearcher(client *http.Client) *Researcher {
	return &Researcher{client: client}
}

func (r *Researcher) Research(ctx context.Context, req ResearchRequest) ResearchResponse {
	partNumber := strings.ToUpper(strings.TrimSpace(req.PartNumber))
	if partNumber == "" {
		return ResearchResponse{Success: false, Error: "Part Number é obrigatório"}
	}

	results, err := r.searchWeb(ctx, partNumber)
	if err != nil {
		fallback := enrichFromSignals(partNumber, req, nil)
		fallback.Success = true
		fallback.Signals = append(fallback.Signals, "web-search:unavailable")
		return fallback
	}

	response := enrichFromSignals(partNumber, req, results)
	response.Success = true
	response.Sources = truncateResults(results, 3)
	return response
}

func (r *Researcher) searchWeb(ctx context.Context, partNumber string) ([]WebResult, error) {
	query := url.QueryEscape(partNumber + " datasheet specifications manufacturer")
	endpoint := "https://duckduckgo.com/html/?q=" + query

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; inventory-web-research/1.0)")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, io.ErrUnexpectedEOF
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 2*1024*1024))
	if err != nil {
		return nil, err
	}

	return parseDuckDuckGoHTML(string(body)), nil
}

var titleURLRe = regexp.MustCompile(`<a[^>]*class="[^"]*result__a[^"]*"[^>]*href="([^"]+)"[^>]*>(.*?)</a>`)
var snippetRe = regexp.MustCompile(`<a[^>]*class="[^"]*result__snippet[^"]*"[^>]*>(.*?)</a>|<div[^>]*class="[^"]*result__snippet[^"]*"[^>]*>(.*?)</div>`)
var htmlTagRe = regexp.MustCompile(`<[^>]+>`)

func parseDuckDuckGoHTML(content string) []WebResult {
	titleMatches := titleURLRe.FindAllStringSubmatch(content, -1)
	snippetMatches := snippetRe.FindAllStringSubmatch(content, -1)

	results := make([]WebResult, 0, len(titleMatches))
	for i, match := range titleMatches {
		if len(match) < 3 {
			continue
		}

		title := sanitizeHTMLText(match[2])
		link := strings.TrimSpace(html.UnescapeString(match[1]))
		snippet := ""
		if i < len(snippetMatches) {
			snippet = sanitizeHTMLText(firstNonEmpty(snippetMatches[i][1], snippetMatches[i][2]))
		}

		if title == "" || link == "" {
			continue
		}
		results = append(results, WebResult{Title: title, URL: link, Snippet: snippet})
	}

	return truncateResults(results, 8)
}

func enrichFromSignals(partNumber string, req ResearchRequest, results []WebResult) ResearchResponse {
	joinedText := strings.ToUpper(strings.Join(req.Tokens, " ")) + " " + strings.ToUpper(req.NormalizedDescription)
	signals := []string{}

	if req.Manufacturer != "" {
		signals = append(signals, "manufacturer:parser")
	}
	if req.Category != "" && req.Category != "unknown" {
		signals = append(signals, "category:parser")
	}

	webBlob := strings.ToUpper(partNumber)
	for _, item := range results {
		webBlob += " " + strings.ToUpper(item.Title) + " " + strings.ToUpper(item.Snippet) + " " + strings.ToUpper(item.URL)
	}

	manufacturer := strings.TrimSpace(req.Manufacturer)
	if manufacturer == "" {
		if detected := detectManufacturerFromText(webBlob + " " + joinedText); detected != "" {
			manufacturer = detected
			signals = append(signals, "manufacturer:web:"+detected)
		}
	}
	if manufacturer == "" {
		manufacturer = inferManufacturerFromPartNumber(partNumber)
		if manufacturer != "" {
			signals = append(signals, "manufacturer:inferred")
		}
	}
	if manufacturer == "" {
		if inferred := inferManufacturerFromTokens(webBlob + " " + joinedText); inferred != "" {
			manufacturer = inferred
			signals = append(signals, "manufacturer:token-inferred")
		}
	}

	category := strings.TrimSpace(strings.ToLower(req.Category))
	if category == "" || category == "unknown" {
		switch {
		case hasAny(webBlob, []string{"DDR3", "DDR4", "DDR5", "RDIMM", "UDIMM", "DIMM", "ECC", "PC4", "PC5"}):
			category = "memory"
			signals = append(signals, "category:web:memory")
		case hasAny(webBlob, []string{"SSD", "NVME", "HDD", "SAS", "SATA", "U.2", "M.2", "PCIE"}):
			category = "disk"
			signals = append(signals, "category:web:disk")
		case hasAny(webBlob, []string{"RJ45", "NIC", "ETHERNET", "SFP", "QSFP", "10G", "25G", "40G", "100G"}):
			category = "network"
			signals = append(signals, "category:web:network")
		default:
			category = "unknown"
		}
	}

	normalizedDescription := strings.TrimSpace(req.NormalizedDescription)
	if normalizedDescription == "" {
		normalizedDescription = buildAutoDescription(partNumber, webBlob, category, manufacturer)
		if normalizedDescription != "" {
			signals = append(signals, "description:web")
		}
	}

	confidence := 0.35
	if len(results) > 0 {
		confidence += 0.25
	}
	if manufacturer != "" {
		confidence += 0.15
	}
	if category != "" && category != "unknown" {
		confidence += 0.15
	}
	if normalizedDescription != "" {
		confidence += 0.1
	}
	if confidence > 0.99 {
		confidence = 0.99
	}

	return ResearchResponse{
		PartNumber:            partNumber,
		Found:                 len(results) > 0,
		Manufacturer:          manufacturer,
		Category:              category,
		NormalizedDescription: normalizedDescription,
		Confidence:            confidence,
		Signals:               dedupeSort(signals),
	}
}

func buildAutoDescription(partNumber, blob, category, manufacturer string) string {
	parts := []string{}
	if manufacturer != "" {
		parts = append(parts, strings.ToUpper(manufacturer))
	}
	if category != "" && category != "unknown" {
		parts = append(parts, strings.ToUpper(category))
	}
	if capacity := firstRegex(blob, `\b\d+\s?(?:GB|TB)\b`); capacity != "" {
		parts = append(parts, strings.ReplaceAll(capacity, " ", ""))
	}
	if speed := firstRegex(blob, `\b(?:\d{4,5}\s?MT/S|\d{4}\s?MHZ|PC[345]-\d{4,5}[A-Z]*)\b`); speed != "" {
		parts = append(parts, strings.ReplaceAll(speed, " ", ""))
	}
	if bus := firstRegex(blob, `\b(?:SAS|SATA|NVME|PCIE|RJ45|SFP|QSFP)\b`); bus != "" {
		parts = append(parts, bus)
	}
	if len(parts) == 0 {
		return partNumber
	}
	return strings.Join(dedupeKeepOrder(parts), " ")
}

func firstRegex(content, pattern string) string {
	re := regexp.MustCompile(pattern)
	return re.FindString(content)
}

func sanitizeHTMLText(raw string) string {
	clean := html.UnescapeString(raw)
	clean = htmlTagRe.ReplaceAllString(clean, " ")
	clean = strings.TrimSpace(strings.Join(strings.Fields(clean), " "))
	return clean
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func truncateResults(values []WebResult, max int) []WebResult {
	if len(values) <= max {
		return values
	}
	return values[:max]
}

func inferManufacturerFromPartNumber(partNumber string) string {
	pn := strings.ToUpper(strings.TrimSpace(partNumber))
	switch {
	case hasAnyPrefix(pn, []string{"HMA", "HMT", "HMC", "HFS", "HFB"}):
		return "SK hynix"
	case hasAnyPrefix(pn, []string{"MTA", "MTF", "MTFD"}):
		return "Micron"
	case hasAnyPrefix(pn, []string{"M3", "M4", "MZ", "PM", "SM"}):
		return "Samsung"
	case hasAnyPrefix(pn, []string{"HUH", "HUS", "HUC"}):
		return "HGST"
	case hasAnyPrefix(pn, []string{"INTEL_", "INTEL", "SSDPE", "SSDS", "SSDP"}) || strings.Contains(pn, "SSDPE"):
		return "Intel"
	case seagatePartNumberPattern.MatchString(pn) || hasAnyPrefix(pn, []string{"SEAG"}):
		return "Seagate"
	case hasAnyPrefix(pn, []string{"UCSC"}):
		return "Cisco"
	case hasAnyPrefix(pn, []string{"KSM", "KVR", "SKC", "SUV"}):
		return "Kingston"
	case hasAnyPrefix(pn, []string{"WDS", "WDC", "WD"}):
		return "Western Digital"
	default:
		return ""
	}
}

func detectManufacturerFromText(value string) string {
	normalizedValue := " " + normalizeManufacturerTokens(value) + " "
	for _, rule := range manufacturerRules {
		for _, alias := range rule.Aliases {
			normalizedAlias := normalizeManufacturerTokens(alias)
			if normalizedAlias != "" && strings.Contains(normalizedValue, " "+normalizedAlias+" ") {
				return rule.Canonical
			}
		}
	}
	return ""
}

func normalizeManufacturerTokens(value string) string {
	upper := strings.ToUpper(strings.TrimSpace(value))
	if upper == "" {
		return ""
	}

	replacer := strings.NewReplacer(
		"-", " ",
		"_", " ",
		"/", " ",
		".", " ",
		",", " ",
		":", " ",
		";", " ",
		"(", " ",
		")", " ",
	)
	return strings.Join(strings.Fields(replacer.Replace(upper)), " ")
}

func hasAnyPrefix(value string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(value, prefix) {
			return true
		}
	}
	return false
}

func inferManufacturerFromTokens(value string) string {
	normalized := normalizeManufacturerTokens(value)
	if normalized == "" {
		return ""
	}

	for _, token := range strings.Fields(normalized) {
		if len(token) < 3 {
			continue
		}
		if inferred := inferManufacturerFromPartNumber(token); inferred != "" {
			return inferred
		}
	}

	return ""
}

func hasAny(joined string, tokens []string) bool {
	for _, token := range tokens {
		if strings.Contains(joined, token) {
			return true
		}
	}
	return false
}

func dedupeSort(values []string) []string {
	seen := map[string]struct{}{}
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}

func dedupeKeepOrder(values []string) []string {
	seen := map[string]struct{}{}
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}
