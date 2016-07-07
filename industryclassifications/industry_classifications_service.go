package industryclassifications

import (
	"encoding/json"

	"github.com/Financial-Times/neo-utils-go/neoutils"
	"github.com/Sirupsen/logrus"
	"github.com/jmcvetta/neoism"
)

// CypherDriver - CypherDriver
type CypherDriver struct {
	cypherRunner neoutils.CypherRunner
	indexManager neoutils.IndexManager
}

//NewCypherDriver instantiate driver
func NewCypherDriver(cypherRunner neoutils.CypherRunner, indexManager neoutils.IndexManager) CypherDriver {
	return CypherDriver{cypherRunner, indexManager}
}

//Initialise initialisation of the indexes
func (pcd CypherDriver) Initialise() error {
	return neoutils.EnsureConstraints(pcd.indexManager, map[string]string{"IndustryClassification": "uuid"})
}

// Check - Feeds into the Healthcheck and checks whether we can connect to Neo and that the datastore isn't empty
func (pcd CypherDriver) Check() error {
	return neoutils.Check(pcd.cypherRunner)
}

// Read - reads a industry Classification given a UUID
func (pcd CypherDriver) Read(uuid string) (interface{}, bool, error) {
	results := []struct {
		UUID      string   `json:"uuid"`
		PrefLabel string   `json:"prefLabel"`
		UUIDs     []string `json:"uuids"`
	}{}

	query := &neoism.CypherQuery{
		Statement: `MATCH (n:IndustryClassification {uuid:{uuid}})
		OPTIONAL MATCH (upp:UPPIdentifier)-[:IDENTIFIES]->(n)
		RETURN n.uuid as uuid, n.prefLabel as prefLabel, collect(upp.value) as uuids`,
		Parameters: map[string]interface{}{
			"uuid": uuid,
		},
		Result: &results,
	}

	err := pcd.cypherRunner.CypherBatch([]*neoism.CypherQuery{query})

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
	r.AlternativeIdentifiers.UUIDS = result.UUIDs
	return r, true, nil
}

//Write - Writes a industry classification node
func (pcd CypherDriver) Write(thing interface{}) error {
	industryToWrite := thing.(industryClassification)

	//cleanUP all the previous IDENTIFIERS referring to that uuid
	deletePreviousIdentifiersQuery := &neoism.CypherQuery{
		Statement: `MATCH (t:Thing {uuid:{uuid}})
		OPTIONAL MATCH (t)<-[iden:IDENTIFIES]-(i)
		DELETE iden, i`,
		Parameters: map[string]interface{}{
			"uuid": industryToWrite.UUID,
		},
	}

	logrus.Infof("Query: %v", deletePreviousIdentifiersQuery)

	//create-update node for IndustryClassification
	statement := `MERGE (n:Thing {uuid: {uuid}})
				set n={allprops}
				set n :IndustryClassification:Classification:Concept`

	createIndustryClassficationQuery := &neoism.CypherQuery{
		Statement: statement,
		Parameters: map[string]interface{}{
			"uuid": industryToWrite.UUID,
			"allprops": map[string]interface{}{
				"uuid":      industryToWrite.UUID,
				"prefLabel": industryToWrite.PrefLabel,
			},
		},
	}

	queryBatch := []*neoism.CypherQuery{deletePreviousIdentifiersQuery, createIndustryClassficationQuery}

	identQuery := &neoism.CypherQuery{
		Statement: `MERGE (t:Thing {uuid:{uuid}})
					CREATE (i:Identifier {value:{uuid}})
					MERGE (t)<-[:IDENTIFIES]-(i)
					set i :Identifier:UPPIdentifier`,
		Parameters: map[string]interface{}{
			"uuid": industryToWrite.UUID,
		},
	}

	queryBatch = append(queryBatch, identQuery)
	return pcd.cypherRunner.CypherBatch([]*neoism.CypherQuery(queryBatch))

}

//Delete - Deletes a Role
func (pcd CypherDriver) Delete(uuid string) (bool, error) {
	clearNode := &neoism.CypherQuery{
		Statement: `
			MATCH (p:Thing {uuid: {uuid}})
			OPTIONAL MATCH (p)<-[ir:IDENTIFIES]-(i:Identifier)
			REMOVE p:Concept
			REMOVE p:Person
			DELETE ir, i
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

	err := pcd.cypherRunner.CypherBatch([]*neoism.CypherQuery{clearNode, removeNodeIfUnused})

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
func (pcd CypherDriver) DecodeJSON(dec *json.Decoder) (interface{}, string, error) {
	r := industryClassification{}
	err := dec.Decode(&r)
	return r, r.UUID, err

}

// Count - Returns a count of the number of roles in this Neo instance
func (pcd CypherDriver) Count() (int, error) {

	results := []struct {
		Count int `json:"c"`
	}{}

	query := &neoism.CypherQuery{
		Statement: `MATCH (n:IndustryClassification) return count(n) as c`,
		Result:    &results,
	}

	err := pcd.cypherRunner.CypherBatch([]*neoism.CypherQuery{query})

	if err != nil {
		return 0, err
	}

	return results[0].Count, nil
}
