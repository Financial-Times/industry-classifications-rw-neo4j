// +build !jenkins

package industryclassifications

import (
	"os"
	"testing"

	"github.com/Financial-Times/base-ft-rw-app-go/baseftrwapp"
	"github.com/Financial-Times/neo-utils-go/neoutils"
	"github.com/jmcvetta/neoism"
	"github.com/stretchr/testify/assert"
)

var industryClassificationDriver baseftrwapp.Service

func TestDelete(t *testing.T) {
	assert := assert.New(t)
	uuid := "12345"

	industryClassificationDriver = getIndustryClassificationCypherDriver(t)

	industryClassificationToDelete := industryClassification{UUID: uuid, PrefLabel: "TestIndustryClassification"}

	assert.NoError(industryClassificationDriver.Write(industryClassificationToDelete), "Failed to write industry classification")

	found, err := industryClassificationDriver.Delete(uuid)
	assert.True(found, "Didn't manage to delete industry classification for uuid %", uuid)
	assert.NoError(err, "Error deleting industry classification for uuid %s", uuid)

	p, found, err := industryClassificationDriver.Read(uuid)

	assert.Equal(industryClassification{}, p, "Found industry classification %s who should have been deleted", p)
	assert.False(found, "Found industry classification for uuid %s who should have been deleted", uuid)
	assert.NoError(err, "Error trying to find industry classification for uuid %s", uuid)
}

func TestCreateAllValuesPresent(t *testing.T) {
	assert := assert.New(t)
	uuid := "12345"
	industryClassificationDriver = getIndustryClassificationCypherDriver(t)

	industryClassificationToWrite := industryClassification{UUID: uuid, PrefLabel: "TestIndustryClassfication"}

	assert.NoError(industryClassificationDriver.Write(industryClassificationToWrite), "Failed to write industry classification")

	readIndustryClassificationForUUIDAndCheckFieldsMatch(t, uuid, industryClassificationToWrite)

	cleanUp(t, uuid)
}

func TestCreateHandlesSpecialCharacters(t *testing.T) {
	assert := assert.New(t)
	uuid := "12345"
	industryClassificationDriver = getIndustryClassificationCypherDriver(t)

	roleToWrite := industryClassification{UUID: uuid, PrefLabel: "Honcho`s pürfèct"}

	assert.NoError(industryClassificationDriver.Write(roleToWrite), "Failed to write industry classification")

	readIndustryClassificationForUUIDAndCheckFieldsMatch(t, uuid, roleToWrite)

	cleanUp(t, uuid)
}

func TestCreateNotAllValuesPresent(t *testing.T) {
	assert := assert.New(t)
	uuid := "12345"
	industryClassificationDriver = getIndustryClassificationCypherDriver(t)

	industryClassificationToWrite := industryClassification{UUID: uuid}

	assert.NoError(industryClassificationDriver.Write(industryClassificationToWrite), "Failed to write industry classification")

	readIndustryClassificationForUUIDAndCheckFieldsMatch(t, uuid, industryClassificationToWrite)

	cleanUp(t, uuid)
}

func readIndustryClassificationForUUIDAndCheckFieldsMatch(t *testing.T, uuid string, expectedIndustryClassification industryClassification) {
	assert := assert.New(t)
	storedIndustryClassification, found, err := industryClassificationDriver.Read(uuid)

	assert.NoError(err, "Error finding industry classification for uuid %s", uuid)
	assert.True(found, "Didn't find industry classification for uuid %s", uuid)
	assert.Equal(expectedIndustryClassification, storedIndustryClassification, "industry classification should be the same")
}

func getIndustryClassificationCypherDriver(t *testing.T) CypherDriver {
	assert := assert.New(t)
	url := os.Getenv("NEO4J_TEST_URL")
	if url == "" {
		url = "http://localhost:7474/db/data"
	}

	db, err := neoism.Connect(url)
	assert.NoError(err, "Failed to connect to Neo4j")
	return NewCypherDriver(neoutils.StringerDb{db}, db)
}

func cleanUp(t *testing.T, uuid string) {
	assert := assert.New(t)
	found, err := industryClassificationDriver.Delete(uuid)
	assert.True(found, "Didn't manage to delete industry classification for uuid %", uuid)
	assert.NoError(err, "Error deleting industry classification for uuid %s", uuid)
}
