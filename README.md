# Industry Classification Reader/Writer for Neo4j (industry-classification-rw-neo4j)
[![Circle CI](https://circleci.com/gh/Financial-Times/industry-classifications-rw-neo4j.svg?style=shield)](https://circleci.com/gh/Financial-Times/industry-classifications-rw-neo4j)[![Go Report Card](https://goreportcard.com/badge/github.com/Financial-Times/industry-classifications-rw-neo4j)](https://goreportcard.com/report/github.com/Financial-Times/industry-classifications-rw-neo4j) [![Coverage Status](https://coveralls.io/repos/github/Financial-Times/industry-classifications-rw-neo4j/badge.svg)](https://coveralls.io/github/Financial-Times/industry-classifications-rw-neo4j) 
__An API for reading/writing industry classification into Neo4j. Expects the industry classification json supplied to be in the following format:

{"uuid":"e8f669e0-e72f-3416-ad0c-09823ec6a27f", "prefLabel":"Industrial Machinery"}

## Installation

For the first time:

`go get github.com/Financial-Times/industry-classification-rw-neo4j`

or update:

`go get -u github.com/Financial-Times/industry-classification-rw-neo4j`

## Running

`$GOPATH/bin/industry-classification-rw-neo4j --neo-url={neo4jUrl} --port={port} --batchSize=50 --graphiteTCPAddress=graphite.ft.com:2003 --graphitePrefix=content.{env}.industry-classification.rw.neo4j.{hostname} --logMetrics=false

All arguments are optional, they default to a local Neo4j install on the default port (7474), application running on port 8080, batchSize of 1024, graphiteTCPAddress of "" (meaning metrics won't be written to Graphite), graphitePrefix of "" and logMetrics false.

NB: the default batchSize is much higher than the throughput the instance data ingester currently can cope with.

To get ALL the industry classifications that are stored in Neo4J currently in this projects under "static-data" as "IndustryClassification.json" and they can be loaded using https://github.com/Financial-Times/up-restutil project.

e.g.

```
cat ../industry-classification-rw-neo4j/static-data/industryClassifications.json |  ./up-restutil put-resources uuid http://localhost:8080/industryclassifications/
```

## Updating the model
The representation of an industry classification is held in the model.go in a struct called industry classification.

## Building
Currently this is not used, so not built or deployed

The build works via git tags. To prepare a new release
- update the version in /puppet/ft-industry_classification_rw_neo4j/Modulefile, e.g. to 0.0.12
- git tag that commit using `git tag 0.0.12`
- `git push --tags`

The deploy also works via git tag and you can also select the environment to deploy to.

## Endpoints
/industryclassification/{uuid}
### PUT
The only mandatory field is the uuid, and the uuid in the body must match the one used on the path.

Every request results in an attempt to update that industry classfication

A successful PUT results in 200.

We run queries in batches. If a batch fails, all failing requests will get a 500 server error response.

Invalid json body input, or uuids that don't match between the path and the body will result in a 400 bad request response.

Example:
`curl -XPUT -H "X-Request-Id: 123" -H "Content-Type: application/json" http://localhost:8080/industryclassifications/e8f669e0-e72f-3416-ad0c-09823ec6a27f --data '{"uuid":"e8f669e0-e72f-3416-ad0c-09823ec6a27f", "prefLabel":"Industrial Machinery"}'`

### GET
Thie internal read should return what got written (i.e., there isn't a public read for industry classification and this is not intended to ever be public either)

If not found, you'll get a 404 response.

Empty fields are omitted from the response.46d1741e-8075-3a4a-872c-71d8130fcb71
`curl -H "X-Request-Id: 123" localhost:8080/industryclassifications/e8f669e0-e72f-3416-ad0c-09823ec6a27f`

### DELETE
Will return 204 if successful, 404 if not found
`curl -XDELETE -H "X-Request-Id: 123" localhost:8080/industryclassifications/e8f669e0-e72f-3416-ad0c-09823ec6a27f`

### Admin endpoints
Healthchecks: [http://localhost:8080/__health](http://localhost:8080/__health)

Ping: [http://localhost:8080/ping](http://localhost:8080/ping) or [http://localhost:8080/__ping](http://localhost:8080/__ping)

### Logging
 the application uses logrus, the logfile is initilaised in main.go.
 logging requires an env app parameter, for all enviromets  other than local logs are written to file
 when running locally logging is written to console (if you want to log locally to file you need to pass in an env parameter that is != local)
 NOTE: build-info end point is not logged as it is called every second from varnish and this information is not needed in  logs/splunk
