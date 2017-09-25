// +build !jenkins

package industryclassifications

import (
	"os"
	"testing"

	"encoding/json"
	"fmt"

	"github.com/Financial-Times/base-ft-rw-app-go/baseftrwapp"
	"github.com/Financial-Times/neo-utils-go/neoutils"
	"github.com/Financial-Times/organisations-rw-neo4j/organisations"
	"github.com/jmcvetta/neoism"
	"github.com/stretchr/testify/assert"
)

const (
	industryClassificationUuid = "f21a5cc0-d326-4e62-b84a-d840c2209fee"
	organisationUuid           = "f9694ba7-eab0-4ce0-8e01-ff64bccb813c"
)

var fullIndustryClassification = industryClassification{
	PrefLabel: "PrefLabel",
	UUID:      industryClassificationUuid,
}

func TestCreateAllValuesPresent(t *testing.T) {
	assert := assert.New(t)

	db := getDatabaseConnectionAndCheckClean(t, assert)
	industryClassificationDriver := getCypherDriver(db)

	defer cleanDB([]string{industryClassificationUuid}, db, t, assert)

	assert.NoError(industryClassificationDriver.Write(fullIndustryClassification, "TID_TEST"), "Failed to write industry classfication")
	actual, found, err := getCypherDriver(db).Read(fullIndustryClassification.UUID, "TID_TEST")

	assert.NoError(err)
	assert.True(found)
	assert.EqualValues(fullIndustryClassification, actual)
}

func TestDeleteWithRelationshipsMaintainsRelationships(t *testing.T) {
	assert := assert.New(t)
	db := getDatabaseConnectionAndCheckClean(t, assert)
	industryClassificationDriver := getCypherDriver(db)

	defer cleanDB([]string{industryClassificationUuid, organisationUuid}, db, t, assert)

	assert.NoError(industryClassificationDriver.Write(fullIndustryClassification, "TID_TEST"), "Failed to write industry classfication")
	writeOrganisation(assert, db)

	found, err := industryClassificationDriver.Delete(industryClassificationUuid, "TID_TEST")

	assert.True(found, "Didn't manage to delete industry classification for uuid %", industryClassificationUuid)
	assert.NoError(err, "Error deleting industry classification for uuid %s", industryClassificationUuid)

	p, found, err := industryClassificationDriver.Read(industryClassificationUuid, "TID_TEST")

	assert.Equal(industryClassification{}, p, "Found industry classification %s who should have been deleted", p)
	assert.False(found, "Found industry classification for uuid %s who should have been deleted", industryClassificationUuid)
	assert.NoError(err, "Error trying to find industry classification for uuid %s", industryClassificationUuid)
	assert.Equal(true, doesThingExistWithIdentifiers(industryClassificationUuid, db, t, assert), "Unable to find a Thing with any Identifiers, uuid: %s", industryClassificationUuid)
}

func TestDeleteWillDeleteEntireNodeIfNoRelationship(t *testing.T) {
	assert := assert.New(t)
	db := getDatabaseConnectionAndCheckClean(t, assert)
	industryClassificationDriver := getCypherDriver(db)

	defer cleanDB([]string{industryClassificationUuid}, db, t, assert)

	assert.NoError(industryClassificationDriver.Write(fullIndustryClassification, "TID_TEST"), "Failed to write industry classification")

	found, err := industryClassificationDriver.Delete(industryClassificationUuid, "TID_TEST")
	assert.True(found, "Didn't manage to delete industry classification for uuid %", industryClassificationUuid)
	assert.NoError(err, "Error deleting industry classification for uuid %s", industryClassificationUuid)

	p, found, err := industryClassificationDriver.Read(industryClassificationUuid, "TID_TEST")

	assert.Equal(industryClassification{}, p, "Found person %s who should have been deleted", p)
	assert.False(found, "Found industry classification for uuid %s who should have been deleted", industryClassificationUuid)
	assert.NoError(err, "Error trying to find industry classification for uuid %s", industryClassificationUuid)
	assert.Equal(false, doesThingExistAtAll(industryClassificationUuid, db, t, assert), "Found thing who should have been deleted uuid: %s", industryClassificationUuid)
}

func writeOrganisation(assert *assert.Assertions, db neoutils.NeoConnection) baseftrwapp.Service {
	orgRW := organisations.NewCypherOrganisationService(db)
	assert.NoError(orgRW.Initialise())
	writeJSONToService(orgRW, "./fixtures/Organisation-Parent-f9694ba7-eab0-4ce0-8e01-ff64bccb813c.json", assert)
	return orgRW
}

func writeJSONToService(service baseftrwapp.Service, pathToJSONFile string, assert *assert.Assertions) {
	f, err := os.Open(pathToJSONFile)
	assert.NoError(err)
	dec := json.NewDecoder(f)
	inst, _, errr := service.DecodeJSON(dec)
	assert.NoError(errr)
	errrr := service.Write(inst, "TID_TEST")
	assert.NoError(errrr)
}

func doesThingExistAtAll(uuid string, db neoutils.NeoConnection, t *testing.T, assert *assert.Assertions) bool {
	result := []struct {
		Uuid string `json:"thing.uuid"`
	}{}

	checkGraph := neoism.CypherQuery{
		Statement: `
			MATCH (a:Thing {uuid: "%s"}) return a.uuid
		`,
		Parameters: neoism.Props{
			"uuid": uuid,
		},
		Result: &result,
	}

	err := db.CypherBatch([]*neoism.CypherQuery{&checkGraph})
	assert.NoError(err)

	if len(result) == 0 {
		return false
	}

	return true
}

func doesThingExistWithIdentifiers(uuid string, db neoutils.NeoConnection, t *testing.T, assert *assert.Assertions) bool {

	result := []struct {
		uuid string `json:"thing.uuid"`
	}{}

	checkGraph := neoism.CypherQuery{
		Statement: `
			MATCH (a:Thing {uuid: "%s"})-[:IDENTIFIES]-(:Identifier)
			WITH collect(distinct a.uuid) as uuid
			RETURN uuid
		`,
		Parameters: neoism.Props{
			"uuid": uuid,
		},
		Result: &result,
	}

	err := db.CypherBatch([]*neoism.CypherQuery{&checkGraph})
	assert.NoError(err)

	if len(result) == 0 {
		return false
	}

	return true
}

func getDatabaseConnectionAndCheckClean(t *testing.T, assert *assert.Assertions) neoutils.NeoConnection {
	db := getDatabaseConnection(assert)
	checkDbClean([]string{industryClassificationUuid, organisationUuid}, db, t)
	return db
}

func getDatabaseConnection(assert *assert.Assertions) neoutils.NeoConnection {
	url := os.Getenv("NEO4J_TEST_URL")
	if url == "" {
		url = "http://localhost:7474/db/data"
	}

	conf := neoutils.DefaultConnectionConfig()
	conf.Transactional = false
	db, err := neoutils.Connect(url, conf)
	assert.NoError(err, "Failed to connect to Neo4j")
	return db
}

func cleanDB(uuidsToClean []string, db neoutils.NeoConnection, t *testing.T, assert *assert.Assertions) {
	qs := make([]*neoism.CypherQuery, len(uuidsToClean))
	for i, uuid := range uuidsToClean {
		qs[i] = &neoism.CypherQuery{
			Statement: fmt.Sprintf(`
			MATCH (a:Thing {uuid: "%s"})
			OPTIONAL MATCH (a)-[rel]-(i:Identifier)
			DELETE rel, i
			DETACH DELETE a`, uuid)}
	}

	err := db.CypherBatch(qs)
	assert.NoError(err)
}

func checkDbClean(uuidsCleaned []string, db neoutils.NeoConnection, t *testing.T) {
	assert := assert.New(t)

	result := []struct {
		Uuid string `json:"thing.uuid"`
	}{}

	checkGraph := neoism.CypherQuery{
		Statement: `
			MATCH (thing) WHERE thing.uuid in {uuids} RETURN thing.uuid
		`,
		Parameters: neoism.Props{
			"uuids": uuidsCleaned,
		},
		Result: &result,
	}
	err := db.CypherBatch([]*neoism.CypherQuery{&checkGraph})
	assert.NoError(err)
	assert.Empty(result)
}

func getCypherDriver(db neoutils.NeoConnection) service {
	cr := NewCypherIndustryClassifcationService(db)
	cr.Initialise()
	return cr
}
