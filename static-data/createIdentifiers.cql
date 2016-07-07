// For every industry classification a new UPPIdentifer is created with the value of the uuid of the industry classification
MATCH (industries:IndustryClassification) 
WITH COLLECT (industries) as inds
UNWIND inds as ind
MERGE (indent:Identifier:UPPIdentifer{value:ind.uuid})-[:IDENTIFIES]->(ind)
RETURN count(indent)