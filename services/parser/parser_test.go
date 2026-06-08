package main

import "testing"

func TestParseLinesMemory(t *testing.T) {
	result := ParseLines([]string{
		"M393A4K40DB3-CWE",
		"Samsung DDR4 32GB 3200MHz RDIMM ECC",
		"SN: SN-1234567890",
	})

	if result.PartNumber != "M393A4K40DB3-CWE" {
		t.Fatalf("expected part number M393A4K40DB3-CWE, got %q", result.PartNumber)
	}
	if result.SerialNumber != "SN-1234567890" {
		t.Fatalf("expected serial number SN-1234567890, got %q", result.SerialNumber)
	}
	if result.Category != "memory" {
		t.Fatalf("expected category memory, got %q", result.Category)
	}
	if result.NormalizedDescription != "DDR4 32GB 3200MHZ RDIMM ECC" {
		t.Fatalf("expected normalized description DDR4 32GB 3200MHZ RDIMM ECC, got %q", result.NormalizedDescription)
	}
}

func TestParseLinesMemoryPc4Label(t *testing.T) {
	result := ParseLines([]string{
		"SAMSUNG M393A4K40DB3-CWE",
		"32G 2Rx4 PC4-3200AA-RB2-12",
		"ECC REG",
	})

	if result.Category != "memory" {
		t.Fatalf("expected category memory, got %q", result.Category)
	}
	if result.PartNumber != "M393A4K40DB3-CWE" {
		t.Fatalf("expected part number M393A4K40DB3-CWE, got %q", result.PartNumber)
	}
	if result.NormalizedDescription != "DDR4 32GB PC4-3200AA 2RX4 ECC" {
		t.Fatalf("expected normalized description DDR4 32GB PC4-3200AA 2RX4 ECC, got %q", result.NormalizedDescription)
	}
}

func TestParseLinesIntelKeepsPrefixedPartNumber(t *testing.T) {
	result := ParseLines([]string{
		"INTEL SSDPE2MX450G7",
		"NVMe SSD 450GB 3D TLC PCI-E 3.0",
	})

	if result.PartNumber != "INTEL_SSDPE2MX450G7" {
		t.Fatalf("expected part number INTEL_SSDPE2MX450G7, got %q", result.PartNumber)
	}
}

func TestParseLinesMemoryLabelFromPhoto(t *testing.T) {
	result := ParseLines([]string{
		"PN: SF472264CKHH60FSDS",
		"REF: M393A2G40B0B-CPB",
		"PC4-2133P-RA0-10 16GB 2Rx4",
	})

	if result.PartNumber != "M393A2G40B0B-CPB" {
		t.Fatalf("expected part number M393A2G40B0B-CPB, got %q", result.PartNumber)
	}
	if result.Category != "memory" {
		t.Fatalf("expected category memory, got %q", result.Category)
	}
	if result.Manufacturer != "Samsung" {
		t.Fatalf("expected manufacturer Samsung, got %q", result.Manufacturer)
	}
	if result.NormalizedDescription != "DDR4 16GB PC4-2133P 2RX4" {
		t.Fatalf("expected normalized description DDR4 16GB PC4-2133P 2RX4, got %q", result.NormalizedDescription)
	}
}

func TestParseLinesMemoryLabelWithMergedTokens(t *testing.T) {
	result := ParseLines([]string{
		"SAMSUNG",
		"16GB2Rx4PC4-2133P",
		"M393A2G40B0B-CPB",
		"S/N: 1234ABC56789",
	})

	if result.PartNumber != "M393A2G40B0B-CPB" {
		t.Fatalf("expected part number M393A2G40B0B-CPB, got %q", result.PartNumber)
	}
	if result.SerialNumber != "1234ABC56789" {
		t.Fatalf("expected serial number 1234ABC56789, got %q", result.SerialNumber)
	}
	if result.Manufacturer != "Samsung" {
		t.Fatalf("expected manufacturer Samsung, got %q", result.Manufacturer)
	}
	if result.Category != "memory" {
		t.Fatalf("expected category memory, got %q", result.Category)
	}
	if result.NormalizedDescription != "DDR4 16GB PC4-2133P 2RX4" {
		t.Fatalf("expected normalized description DDR4 16GB PC4-2133P 2RX4, got %q", result.NormalizedDescription)
	}
}

func TestParseLinesDisk(t *testing.T) {
	result := ParseLines([]string{
		"INTEL SSDPE2MX450G7",
		"NVMe SSD 450GB 3D TLC PCI-E 3.0",
	})

	if result.Category != "disk" {
		t.Fatalf("expected category disk, got %q", result.Category)
	}
	if result.PartNumber != "INTEL_SSDPE2MX450G7" {
		t.Fatalf("expected part number INTEL_SSDPE2MX450G7, got %q", result.PartNumber)
	}
}

func TestParseLinesPrefersSNOverWWN(t *testing.T) {
	result := ParseLines([]string{
		"MODEL: WDS240G2G0A",
		"WWN: 5001B448B6351C20",
		"S/N: 21493L800954",
	})

	if result.SerialNumber != "21493L800954" {
		t.Fatalf("expected serial number 21493L800954, got %q", result.SerialNumber)
	}
}

func TestParseLinesIgnoresWWNWhenNoSNLabel(t *testing.T) {
	result := ParseLines([]string{
		"MODEL: WDS240G2G0A",
		"WORLD WIDE NAME: 5001B448B6351C20",
	})

	if result.SerialNumber != "" {
		t.Fatalf("expected empty serial number, got %q", result.SerialNumber)
	}
}
