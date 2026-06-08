package main

import (
	"regexp"
	"sort"
	"strings"
)

var manufacturerAliases = map[string]string{
	"HPE":             "HPE",
	"HP":              "HP",
	"HISENSE":         "Hisense",
	"INTEL":           "Intel",
	"KINGSTON":        "Kingston",
	"MICRON":          "Micron",
	"SAMSUNG":         "Samsung",
	"SEAGATE":         "Seagate",
	"HGST":            "HGST",
	"CISCO":           "Cisco",
	"DELL":            "Dell",
	"WESTERN DIGITAL": "Western Digital",
	"WD":              "Western Digital",
}

var partNumberPattern = regexp.MustCompile(`\b[A-Z0-9][A-Z0-9._/-]{7,}\b`)
var serialNumberPattern = regexp.MustCompile(`\b(?:SN|S/N|S\.N\.|SER\.?\s*NO\.?|SERIAL|SERIAL\s*NO|SERIALNUMBER)\s*[:#\-]?\s*([A-Z0-9][A-Z0-9\-]{5,})\b`)
var labeledPartNumberPattern = regexp.MustCompile(`\b(REF|P/N|PN|PART\s*NO|PARTNUMBER)\s*[:#-]?\s*([A-Z0-9][A-Z0-9._/-]{5,})\b`)
var serialCandidatePattern = regexp.MustCompile(`\b[A-Z0-9][A-Z0-9\-]{7,}\b`)
var wwnHexPattern = regexp.MustCompile(`^[0-9A-F]{16,32}$`)

func NormalizeText(lines []string) string {
	cleaned := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		cleaned = append(cleaned, strings.ToUpper(trimmed))
	}
	return strings.Join(cleaned, " ")
}

func ExtractPartNumber(lines []string) RuleMatch {
	joined := NormalizeText(lines)
	if joined == "" {
		return RuleMatch{}
	}

	if labeled := extractLabeledPartNumbers(joined); labeled != "" {
		return RuleMatch{Value: labeled, Signals: []string{"labeled-part-number"}}
	}

	for _, line := range lines {
		if normalized := normalizePrefixedPartNumber(line); normalized != "" {
			return RuleMatch{Value: normalized, Signals: []string{"prefixed-part-number"}}
		}
	}

	candidates := partNumberPattern.FindAllString(joined, -1)
	best := ""
	bestScore := -1
	signals := []string{}

	for _, candidate := range candidates {
		if len(candidate) < 8 {
			continue
		}
		if strings.HasPrefix(candidate, "SN") && strings.Contains(candidate, "-") {
			continue
		}
		if !hasDigit(candidate) {
			continue
		}
		if isLikelySpecToken(candidate) {
			continue
		}
		score := candidateScore(candidate)
		if score > bestScore || (score == bestScore && len(candidate) > len(best)) {
			best = candidate
			bestScore = score
		}
		signals = append(signals, candidate)
	}

	return RuleMatch{Value: best, Signals: dedupeSignals(signals)}
}

func extractLabeledPartNumbers(joined string) string {
	matches := labeledPartNumberPattern.FindAllStringSubmatch(joined, -1)
	if len(matches) == 0 {
		return ""
	}

	pnCandidate := ""
	refCandidate := ""

	for _, match := range matches {
		if len(match) < 3 {
			continue
		}
		label := strings.ReplaceAll(match[1], " ", "")
		value := match[2]
		if len(value) < 6 || !hasDigit(value) {
			continue
		}

		switch label {
		case "REF":
			refCandidate = value
		case "PN", "P/N", "PARTNO", "PARTNUMBER":
			pnCandidate = value
		}
	}

	if refCandidate != "" {
		return refCandidate
	}
	return pnCandidate
}

func candidateScore(candidate string) int {
	score := 0

	if hasDigit(candidate) {
		score += 3
	}
	if strings.Contains(candidate, "-") {
		score += 3
	}
	if strings.Contains(candidate, "_") {
		score += 1
	}
	if strings.HasPrefix(candidate, "M") {
		score += 2
	}
	if strings.HasPrefix(candidate, "SN") || strings.HasPrefix(candidate, "SNR") {
		score -= 3
	}
	if strings.Contains(candidate, "REV") {
		score -= 2
	}

	return score
}

func normalizePrefixedPartNumber(line string) string {
	fields := strings.Fields(strings.ToUpper(strings.TrimSpace(line)))
	if len(fields) < 2 {
		return ""
	}

	first := strings.Trim(fields[0], ",.;:()[]{}")
	second := strings.Trim(fields[1], ",.;:()[]{}")

	if _, ok := manufacturerAliases[first]; !ok {
		return ""
	}

	if len(second) < 6 || !hasDigit(second) {
		return ""
	}

	if !partNumberPattern.MatchString(second) && !strings.Contains(second, "-") && !strings.Contains(second, "_") {
		return ""
	}

	// Apenas Intel usa prefixo de fabricante no catálogo atual (ex.: INTEL_SSDPE2MX450G7).
	if first == "INTEL" {
		return first + "_" + second
	}

	// Para os demais fabricantes, usar o PN bruto encontrado na etiqueta.
	return second
}

func ExtractSerialNumber(lines []string) RuleMatch {
	if len(lines) == 0 {
		return RuleMatch{}
	}

	partNumber := ExtractPartNumber(lines).Value

	// Regra principal: usar valor explicitamente rotulado como S/N (ou variações).
	for _, line := range lines {
		upperLine := strings.ToUpper(strings.TrimSpace(line))
		if upperLine == "" {
			continue
		}

		matches := serialNumberPattern.FindAllStringSubmatch(upperLine, -1)
		for _, match := range matches {
			if len(match) < 2 {
				continue
			}

			candidate := strings.TrimSpace(match[1])
			if candidate == "" || isLikelySpecToken(candidate) || isLikelyWWNValue(candidate, upperLine) {
				continue
			}

			return RuleMatch{
				Value:   candidate,
				Signals: []string{"serial-prefix"},
			}
		}
	}

	// Fallback: candidatos genéricos apenas quando não há S/N explícito.
	best := ""
	bestScore := -1
	for _, line := range lines {
		upperLine := strings.ToUpper(strings.TrimSpace(line))
		if upperLine == "" {
			continue
		}

		candidates := serialCandidatePattern.FindAllString(upperLine, -1)
		for _, candidate := range candidates {
			if candidate == partNumber {
				continue
			}
			if strings.Contains(upperLine, "MODEL") || strings.Contains(upperLine, "P/N") || strings.Contains(upperLine, "PN") || strings.Contains(upperLine, "REF") {
				continue
			}
			if isLikelySpecToken(candidate) || isLikelyWWNValue(candidate, upperLine) {
				continue
			}
			if !hasDigit(candidate) {
				continue
			}

			score := len(candidate)
			if strings.HasPrefix(candidate, "SN") {
				score += 5
			}
			if strings.Contains(candidate, "-") {
				score += 2
			}

			if score > bestScore {
				best = candidate
				bestScore = score
			}
		}
	}

	if best == "" {
		return RuleMatch{}
	}

	return RuleMatch{Value: best, Signals: []string{"serial-fallback"}}
}

func isLikelyWWNValue(value, line string) bool {
	normalized := strings.ToUpper(strings.TrimSpace(value))
	lineUpper := strings.ToUpper(line)

	if strings.Contains(lineUpper, "WWN") || strings.Contains(lineUpper, "WORLD WIDE NAME") || strings.Contains(lineUpper, "NAA") {
		return true
	}

	if wwnHexPattern.MatchString(normalized) {
		return true
	}

	return false
}

func DetectManufacturer(lines []string, partNumber string) RuleMatch {
	joined := NormalizeText(lines)
	for alias, canonical := range manufacturerAliases {
		if strings.Contains(joined, alias) {
			return RuleMatch{Value: canonical, Signals: []string{"manufacturer:" + canonical}}
		}
	}

	if inferred := inferManufacturerFromPartNumber(partNumber); inferred != "" {
		return RuleMatch{Value: inferred, Signals: []string{"manufacturer:inferred:" + inferred}}
	}

	return RuleMatch{}
}

func inferManufacturerFromPartNumber(partNumber string) string {
	pn := strings.ToUpper(strings.TrimSpace(partNumber))
	switch {
	case strings.HasPrefix(pn, "M"):
		return "Samsung"
	case strings.HasPrefix(pn, "HUH"):
		return "HGST"
	case strings.HasPrefix(pn, "INTEL_") || strings.HasPrefix(pn, "INTEL"):
		return "Intel"
	case strings.HasPrefix(pn, "SEAG"):
		return "Seagate"
	case strings.HasPrefix(pn, "UCSC"):
		return "Cisco"
	default:
		return ""
	}
}

func ClassifyCategory(lines []string) RuleMatch {
	joined := NormalizeText(lines)
	signals := []string{}

	if hasAny(joined, []string{"DDR3", "DDR4", "DDR5", "DIMM", "RDIMM", "UDIMM", "SO-DIMM", "SODIMM", "ECC", "PC3", "PC4", "PC5", "MT/S", "2RX", "1RX"}) {
		signals = append(signals, "memory")
		return RuleMatch{Value: "memory", Signals: signals}
	}

	if hasAny(joined, []string{"NVME", "NVME", "SSD", "HDD", "SAS", "SATA", "M.2", "U.2", "PCIE"}) {
		signals = append(signals, "storage")
		return RuleMatch{Value: "disk", Signals: signals}
	}

	if hasAny(joined, []string{"RJ45", "SFP", "QSFP", "NIC", "ETHERNET", "10G", "25G", "40G", "100G"}) {
		signals = append(signals, "network")
		return RuleMatch{Value: "network", Signals: signals}
	}

	return RuleMatch{Value: "unknown"}
}

func BuildNormalizedDescription(lines []string, category string) string {
	joined := NormalizeText(lines)
	if joined == "" {
		return ""
	}

	parts := []string{}

	switch category {
	case "memory":
		if token := firstMatch(joined, []string{"DDR3", "DDR4", "DDR5"}); token != "" {
			parts = append(parts, token)
		} else if inferred := inferDDRFromPC(joined); inferred != "" {
			parts = append(parts, inferred)
		}
		if token := firstCapacity(joined); token != "" {
			parts = append(parts, token)
		}
		if token := firstSpeed(joined); token != "" {
			parts = append(parts, token)
		}
		if token := firstRank(joined); token != "" {
			parts = append(parts, token)
		}
		if token := firstMatch(joined, []string{"RDIMM", "UDIMM", "SO-DIMM", "DIMM"}); token != "" {
			parts = append(parts, token)
		}
		if strings.Contains(joined, "ECC") {
			parts = append(parts, "ECC")
		}
	case "disk":
		if token := firstMatch(joined, []string{"NVME", "SSD", "HDD"}); token != "" {
			parts = append(parts, token)
		}
		if token := firstCapacity(joined); token != "" {
			parts = append(parts, token)
		}
		if token := firstMatch(joined, []string{"SAS", "SATA", "PCIE", "M.2", "U.2"}); token != "" {
			parts = append(parts, token)
		}
	case "network":
		if token := firstMatch(joined, []string{"NIC", "ADAPTER", "MODULE"}); token != "" {
			parts = append(parts, token)
		}
		if token := firstNetworkSpeed(joined); token != "" {
			parts = append(parts, token)
		}
		if token := firstMatch(joined, []string{"ETHERNET", "RJ45", "SFP", "QSFP"}); token != "" {
			parts = append(parts, token)
		}
	default:
		parts = append(parts, joined)
	}

	if len(parts) == 0 {
		return joined
	}

	return strings.Join(dedupeKeepOrder(parts), " ")
}

func ParseLines(lines []string) ParseResponse {
	partNumber := ExtractPartNumber(lines)
	serialNumber := ExtractSerialNumber(lines)
	manufacturer := DetectManufacturer(lines, partNumber.Value)
	category := ClassifyCategory(lines)

	confidence := 0.45
	signals := []string{}
	if partNumber.Value != "" {
		confidence += 0.2
		signals = append(signals, partNumber.Signals...)
	}
	if serialNumber.Value != "" {
		confidence += 0.1
		signals = append(signals, serialNumber.Signals...)
	}
	if manufacturer.Value != "" {
		confidence += 0.1
		signals = append(signals, manufacturer.Signals...)
	}
	if category.Value != "unknown" {
		confidence += 0.15
		signals = append(signals, category.Signals...)
	}

	if confidence > 0.99 {
		confidence = 0.99
	}

	return ParseResponse{
		Success:               true,
		PartNumber:            partNumber.Value,
		SerialNumber:          serialNumber.Value,
		Manufacturer:          manufacturer.Value,
		Category:              category.Value,
		NormalizedDescription: BuildNormalizedDescription(lines, category.Value),
		Confidence:            confidence,
		Signals:               dedupeSignals(signals),
		Tokens:                tokenizeLines(lines),
	}
}

func hasAny(joined string, tokens []string) bool {
	for _, token := range tokens {
		if strings.Contains(joined, token) {
			return true
		}
	}
	return false
}

func hasDigit(value string) bool {
	for _, char := range value {
		if char >= '0' && char <= '9' {
			return true
		}
	}
	return false
}

func isLikelySpecToken(value string) bool {
	if strings.HasSuffix(value, "GB") || strings.HasSuffix(value, "MHZ") || strings.HasSuffix(value, "RPM") {
		return true
	}
	if strings.Contains(value, "PC4-") || strings.Contains(value, "PC3-") {
		return true
	}
	return false
}

func firstMatch(joined string, tokens []string) string {
	for _, token := range tokens {
		if strings.Contains(joined, token) {
			return token
		}
	}
	return ""
}

func firstCapacity(joined string) string {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`\d+\s?(?:GB|TB)`),
		regexp.MustCompile(`\d+\s?G`),
	}

	for _, pattern := range patterns {
		match := pattern.FindString(joined)
		if match == "" {
			continue
		}
		clean := strings.ReplaceAll(match, " ", "")
		if strings.HasSuffix(clean, "G") {
			clean += "B"
		}
		return clean
	}

	return ""
}

func firstSpeed(joined string) string {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`\d{4}\s?MHZ`),
		regexp.MustCompile(`\d{4,5}\s?MT/S`),
		regexp.MustCompile(`PC[345]-\d{4,5}[A-Z]*`),
	}
	for _, pattern := range patterns {
		if match := pattern.FindString(joined); match != "" {
			return strings.ReplaceAll(match, " ", "")
		}
	}
	return ""
}

func firstRank(joined string) string {
	pattern := regexp.MustCompile(`[12]RX[48]`)
	return pattern.FindString(joined)
}

func inferDDRFromPC(joined string) string {
	switch {
	case strings.Contains(joined, "PC3"):
		return "DDR3"
	case strings.Contains(joined, "PC4"):
		return "DDR4"
	case strings.Contains(joined, "PC5"):
		return "DDR5"
	default:
		return ""
	}
}

func firstNetworkSpeed(joined string) string {
	pattern := regexp.MustCompile(`\b\d+G\b`)
	return pattern.FindString(joined)
}

func tokenizeLines(lines []string) []string {
	tokens := make([]string, 0)
	for _, line := range lines {
		for _, token := range strings.Fields(strings.ToUpper(line)) {
			token = strings.TrimSpace(strings.Trim(token, ",.;:()[]{}"))
			if token != "" {
				tokens = append(tokens, token)
			}
		}
	}
	return dedupeStrings(tokens)
}

func dedupeSignals(values []string) []string {
	return dedupeStrings(values)
}

func dedupeStrings(values []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(values))
	for _, value := range values {
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
	seen := make(map[string]struct{})
	result := make([]string, 0, len(values))
	for _, value := range values {
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
