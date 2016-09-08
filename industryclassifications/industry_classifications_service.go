package industryclassifications

import (
	"encoding/json"

	"github.com/Financial-Times/neo-utils-go/neoutils"
	"github.com/jmcvetta/neoism"
)

// CypherDriver - CypherDriver
type service struct {
	conn neoutils.NeoConnection
}

//NewCypherDriver instantiate driver
func NewCypherIndustryClassifcationService(cypherRunner neoutils.NeoConnection) service {
	return service{cypherRunner}
}

//Initialise initialisation of the indexes
func (s service) Initialise() error {
	return s.conn.EnsureConstraints(map[string]string{
		"IndustryClassification": "uuid"})
}

// Check - Feeds into the Healthcheck and checks whether we can connect to Neo and that the datastore isn't empty
func (s service) Check() error {
	return neoutils.Check(s.conn)
}

// Read - reads a industry Classification given a UUID
func (s service) Read(uuid string) (interface{}, bool, error) {
	results := []struct {
		UUID      string `json:"uuid"`
		PrefLabel string `json:"prefLabel"`
	}{}

	query := &neoism.CypherQuery{
		Statement: `MATCH (n:IndustryClassification {uuid:{uuid}}) return n.uuid
		as uuid, n.prefLabel as prefLabel`,
		Parameters: map[string]interface{}{
			"uuid": uuid,
		},
		Result: &results,
	}

	err := s.conn.CypherBatch([]*neoism.CypherQuery{query})

	if err != nil {
		return industryClassification{}, false, err
	}

	if len(results) == 0 {
		return industryClassification{}, false, nil
	}

	result := results[0]

	r := industryClassification{
		UUID:      result.UUID,
		PrefLabel: result.PrefLabel,
	}
	return r, true, nil
}

//Write - Writes a industry classification node
func (s service) Write(thing interface{}) error {
	r := thing.(industryClassification)

	params := map[string]interface{}{
		"uuid": r.UUID,
	}

	if r.PrefLabel != "" {
		params["prefLabel"] = r.PrefLabel
	}

	statement := `MERGE (n:Thing {uuid: {uuid}})
				set n={allprops}
				set n :IndustryClassification:Classification:Concept`

	query := &neoism.CypherQuery{
		Statement: statement,
		Parameters: map[string]interface{}{
			"uuid":     r.UUID,
			"allprops": params,
		},
	}

	return s.conn.CypherBatch([]*neoism.CypherQuery{query})

}

//Delete - Deletes a Role
func (s service) Delete(uuid string) (bool, error) {
	clearNode := &neoism.CypherQuery{
		Statement: `
			MATCH (p:Thing {uuid: {uuid}})
			REMOVE p:IndustryClassification
			SET p={props}
		`,
		Parameters: map[string]interface{}{
			"uuid": uuid,
			"props": map[string]interface{}{
				"uuid": uuid,
			},
		},
		IncludeStats: true,
	}

	removeNodeIfUnused := &neoism.CypherQuery{
		Statement: `
			MATCH (p:Thing {uuid: {uuid}})
			OPTIONAL MATCH (p)-[a]-(x)
			WITH p, count(a) AS relCount
			WHERE relCount = 0
			DELETE p
		`,
		Parameters: map[string]interface{}{
			"uuid": uuid,
		},
	}

	err := s.conn.CypherBatch([]*neoism.CypherQuery{clearNode, removeNodeIfUnused})

	s1, err := clearNode.Stats()
	if err != nil {
		return false, err
	}

	var deleted bool
	if s1.ContainsUpdates && s1.LabelsRemoved > 0 {
		deleted = true
	}

	return deleted, err
}

// DecodeJSON - Decodes JSON into role
func (s service) DecodeJSON(dec *json.Decoder) (interface{}, string, error) {
	r := industryClassification{}
	err := dec.Decode(&r)
	return r, r.UUID, err

}

// Count - Returns a count of the number of roles in this Neo instance
func (s service) Count() (int, error) {

	results := []struct {
		Count int `json:"c"`
	}{}

	query := &neoism.CypherQuery{
		Statement: `MATCH (n:IndustryClassification) return count(n) as c`,
		Result:    &results,
	}

	err := s.conn.CypherBatch([]*neoism.CypherQuery{query})

	if err != nil {
		return 0, err
	}

	return results[0].Count, nil
}
