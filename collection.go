package main

import (
	"context"
	"fmt"
	"log"
	"time"

	chroma "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/types"
)

func (d *chromaDB) getOrCreateCollection(name string) (*chroma.Collection, bool) {
	existingCollection, err := d.client.GetCollection(context.TODO(), name, d.openaiEf)
	if err == nil {
		log.Println("Collection already exists: " + name)
		return existingCollection, true
	}
	log.Printf("Creating collection \"%s\"", name)
	newCollection, err := d.client.CreateCollection(
		context.TODO(),
		name,
		nil,
		true,
		d.openaiEf,
		types.L2,
	)
	if err != nil {
		log.Fatalf("Error creating collection: %s \n", err)
	}
	return newCollection, false
}

func (d *chromaDB) makeCollectionWithRecords(name string, records []string) (*chroma.Collection, error) {
	log.Println("Inserting records to collection \"" + name + "\" with " + fmt.Sprint(len(records)) + " records")
	startTime := time.Now()
	defer func() {
		elapsedTime := time.Since(startTime)
		log.Printf("Finished inserting %d records into \"%s\" in %.2fs", len(records), name, elapsedTime.Seconds())
	}()
	batchSize := 1000

	rs, err := types.NewRecordSet(
		types.WithEmbeddingFunction(d.openaiEf),
		types.WithIDGenerator(types.NewULIDGenerator()),
	)
	if err != nil {
		log.Fatalf("Error creating record set: %s \n", err)
	}
	// remove duplicates
	uniqueRecords := make(map[string]struct{})
	var deduplicatedRecords []string

	for _, record := range records {
		if _, exists := uniqueRecords[record]; !exists {
			uniqueRecords[record] = struct{}{}
			deduplicatedRecords = append(deduplicatedRecords, record)
		}
	}

	records = deduplicatedRecords
	recordLength := len(records)
	log.Printf("%d unique records for \"%s\"", recordLength, name)

	collection, existed := d.getOrCreateCollection(name)
	if existed {
		return collection, nil
	}

	// Insert records in batches of `batchSize` records, and the last batch will have the remaining records
	insertedCount := 0
	for start := 0; start < recordLength; start += batchSize {
		end := start + batchSize
		if end > recordLength {
			end = recordLength
		}
		for index, record := range records[start:end] {
			if record != "" {
				rs.WithRecord(types.WithDocument(record))
				insertedCount++
				if start == 0 && (index == 0 || index == end-start-1) {
					log.Printf("Processing record %d of batch starting at %d", index, start)
				}
			}
		}
		_, err = rs.BuildAndValidate(context.TODO())
		if err != nil {
			log.Fatalf("Error validating record set: %s \n", err)
		}

		_, err = collection.AddRecords(context.Background(), rs)
		if err != nil {
			log.Fatalf("Error adding documents: %s \n", err)
		}
	}
	log.Printf("Inserted %d records into collection \"%s\"", insertedCount, name)
	return collection, nil
}

func (d *chromaDB) query(collectionName string, query string, numResults int) ([]string, error) {
	log.Println("Querying collection \"" + collectionName + "\" with query: " + query)
	startTime := time.Now()
	defer func() {
		elapsedTime := time.Since(startTime)
		log.Printf("Queried collection \"%s\" in %.2fs\n", collectionName, elapsedTime.Seconds())
	}()
	collection, err := d.client.GetCollection(context.TODO(), collectionName, d.openaiEf)
	if err != nil {
		log.Fatalf("Error getting collection: %s \n", err)
	}
	qr, err := collection.Query(context.TODO(), []string{query}, int32(numResults), nil, nil, nil)
	if err != nil {
		return nil, err
	}
	var results []string
	for _, doc := range qr.Documents {
		results = append(results, doc...)
	}
	return results, nil
}
