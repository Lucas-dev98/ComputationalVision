package main

import "testing"

func TestParseDuckDuckGoHTML(t *testing.T) {
	htmlFixture := `
	<div class="result">
	  <a class="result__a" href="https://example.com/datasheet">M393A4K40DB3-CWE Datasheet</a>
	  <a class="result__snippet">Samsung DDR4 RDIMM ECC 32GB module specification.</a>
	</div>`

	results := parseDuckDuckGoHTML(htmlFixture)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].URL != "https://example.com/datasheet" {
		t.Fatalf("unexpected URL: %s", results[0].URL)
	}
}

func TestEnrichFromSignals(t *testing.T) {
	req := ResearchRequest{PartNumber: "M393A4K40DB3-CWE"}
	results := []WebResult{{
		Title:   "Samsung M393A4K40DB3-CWE DDR4 RDIMM",
		URL:     "https://example.com",
		Snippet: "32GB ECC module PC4-3200",
	}}

	response := enrichFromSignals(req.PartNumber, req, results)
	if response.Manufacturer != "Samsung" {
		t.Fatalf("expected Samsung manufacturer, got %q", response.Manufacturer)
	}
	if response.Category != "memory" {
		t.Fatalf("expected memory category, got %q", response.Category)
	}
	if response.Confidence <= 0.5 {
		t.Fatalf("expected confidence > 0.5, got %.2f", response.Confidence)
	}
}

func TestFallbackInferenceWithoutWebResults(t *testing.T) {
	req := ResearchRequest{PartNumber: "INTEL_SSDPE2MX450G7"}
	response := enrichFromSignals(req.PartNumber, req, nil)
	if response.Manufacturer != "Intel" {
		t.Fatalf("expected Intel manufacturer, got %q", response.Manufacturer)
	}
}
