// +build !jenkins

package industryclassifications

import (
	"os"
	"testing"

	"github.com/Financial-Times/base-ft-rw-app-go/baseftrwapp"
	"github.com/Financial-Times/neo-utils-go/neoutils"
	"github.com/stretchr/testify/assert"
)

var industryClassificationDriver baseftrwapp.Service

func TestDelete(t *testing.T) {
	assert := assert.New(t)
	uuid := "54321"

	industryClassificationDriver = getIndustryClassificationService(t)

	industryClassificationToDelete := industryClassification{UUID: uuid, PrefLabel: "TestIndustryClassification"}
	industryClassificationToDelete.AlternativeIdentifiers.UUIDS = []string{uuid}

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
	industryClassificationDriver = getIndustryClassificationService(t)

	industryClassificationToWrite := industryClassification{UUID: uuid, PrefLabel: "TestIndustryClassfication"}
	industryClassificationToWrite.AlternativeIdentifiers.UUIDS = []string{uuid}

	assert.NoError(industryClassificationDriver.Write(industryClassificationToWrite), "Failed to write industry classification")

	readIndustryClassificationForUUIDAndCheckFieldsMatch(t, uuid, industryClassificationToWrite)

	cleanUp(t, uuid)
}

func TestCreateHandlesSpecialCharacters(t *testing.T) {
	assert := assert.New(t)
	uuid := "12345"
	industryClassificationDriver = getIndustryClassificationService(t)

	industryClassificationToWrite := industryClassification{UUID: uuid, PrefLabel: "Honcho`s pürfèct"}
	industryClassificationToWrite.AlternativeIdentifiers.UUIDS = []string{uuid}

	assert.NoError(industryClassificationDriver.Write(industryClassificationToWrite), "Failed to write industry classification")

	readIndustryClassificationForUUIDAndCheckFieldsMatch(t, uuid, industryClassificationToWrite)

	cleanUp(t, uuid)
}

func TestCreateNotAllValuesPresent(t *testing.T) {
	assert := assert.New(t)

	uuid := "12345"
	industryClassificationDriver = getIndustryClassificationService(t)

	industryClassificationToWrite := industryClassification{UUID: uuid}
	industryClassificationToWrite.AlternativeIdentifiers.UUIDS = []string{uuid}

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

func getIndustryClassificationService(t *testing.T) service {
	assert := assert.New(t)
	url := os.Getenv("NEO4J_TEST_URL")
	if url == "" {
		url = "http://localhost:7474/db/data"
	}

	conf := neoutils.DefaultConnectionConfig()
	conf.Transactional = false
	db, err := neoutils.Connect(url, conf)
	assert.NoError(err, "Failed to connect to Neo4j")
	return NewCypherIndustryClassifcationService(db)
}

func cleanUp(t *testing.T, uuid string) {
	assert := assert.New(t)
	found, err := industryClassificationDriver.Delete(uuid)
	assert.True(found, "Didn't manage to delete industry classification for uuid %", uuid)
	assert.NoError(err, "Error deleting industry classification for uuid %s", uuid)
}
