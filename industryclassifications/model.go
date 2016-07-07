package industryclassifications

type industryClassification struct {
	UUID                   string                 `json:"uuid"`
	PrefLabel              string                 `json:"prefLabel,omitempty"`
	AlternativeIdentifiers alternativeIdentifiers `json:"alternativeIdentifiers"`
}

type alternativeIdentifiers struct {
	FactsetIdentifier string   `json:"factsetIdentifier,omitempty"`
	UUIDS             []string `json:"uuids"`
}
