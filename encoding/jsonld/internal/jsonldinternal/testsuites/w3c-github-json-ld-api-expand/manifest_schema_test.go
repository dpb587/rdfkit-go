package testsuite

type manifestSchema struct {
	Sequence []struct {
		ID     string                       `json:"@id"`
		Type   []string                     `json:"@type"`
		Input  string                       `json:"input"`
		Expect string                       `json:"expect"`
		Option manifestSchemaSequenceOption `json:"option"`
	} `json:"sequence"`
}

type manifestSchemaSequenceOption struct {
	Base                  string `json:"base"`
	ExpandContext         string `json:"expandContext"`
	ProcessingMode        string `json:"processingMode"`
	ProduceGeneralizedRdf bool   `json:"produceGeneralizedRdf"`
	RDFDirection          string `json:"rdfDirection"`
	SpecVersion           string `json:"specVersion"`
}
