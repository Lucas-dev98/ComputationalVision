package main

type ResearchRequest struct {
	PartNumber            string   `json:"part_number"`
	Manufacturer          string   `json:"manufacturer,omitempty"`
	Category              string   `json:"category,omitempty"`
	NormalizedDescription string   `json:"normalized_description,omitempty"`
	Tokens                []string `json:"tokens,omitempty"`
}

type WebResult struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet"`
}

type ResearchResponse struct {
	Success               bool        `json:"success"`
	PartNumber            string      `json:"part_number"`
	Found                 bool        `json:"found"`
	Manufacturer          string      `json:"manufacturer,omitempty"`
	Category              string      `json:"category,omitempty"`
	NormalizedDescription string      `json:"normalized_description,omitempty"`
	Confidence            float64     `json:"confidence"`
	Sources               []WebResult `json:"sources,omitempty"`
	Signals               []string    `json:"signals,omitempty"`
	Error                 string      `json:"error,omitempty"`
}
