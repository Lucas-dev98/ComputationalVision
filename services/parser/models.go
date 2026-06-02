package main

type ParseRequest struct {
	Text       []string `json:"text"`
	Structured any      `json:"structured,omitempty"`
}

type ParseResponse struct {
	Success               bool     `json:"success"`
	PartNumber            string   `json:"part_number,omitempty"`
	SerialNumber          string   `json:"serial_number,omitempty"`
	Manufacturer          string   `json:"manufacturer,omitempty"`
	Category              string   `json:"category,omitempty"`
	NormalizedDescription string   `json:"normalized_description,omitempty"`
	Confidence            float64  `json:"confidence"`
	Signals               []string `json:"signals,omitempty"`
	Tokens                []string `json:"tokens,omitempty"`
	Error                 string   `json:"error,omitempty"`
}

type RuleMatch struct {
	Value   string
	Signals []string
}
