package sqlite

import (
	"context"
	"ilserver/domain"
	"ilserver/pkg/utility"
	storage "ilserver/storage/control"
	"math/rand"
	"testing"
)

func Test_InsertTopic(t *testing.T) {
	storage, err := Instance()
	if err != nil {
		t.Fatal()
		return
	}

	// ***

	topicsForInsert := generateFakeTopics()
	for i := range topicsForInsert {
		_, err := storage.InsertTopic(context.Background(), topicsForInsert[i])
		if err != nil {
			t.Errorf("Insert topic failed. Err: %v", err)
			return
		}
	}

	// ***

	num := recordCountInTableWithChecks(t, storage, "Topics")
	if num != len(topicsForInsert) {
		t.Errorf("Record count in the table is not equal %v", len(topicsForInsert))
		return
	}

	// ***

	// testing?
	deleteTopicsWithChecks(t, storage)
}

func Test_InsertTopics(t *testing.T) {
	storage := instanceWithChecks(t)

	// ***

	topicsForInsert := generateFakeTopics()
	insertTopicsWithChecks(t, storage, topicsForInsert)

	// ***

	num := recordCountInTableWithChecks(t, storage, "Topics")
	if num != len(topicsForInsert) {
		t.Errorf("Record count in the table is not equal %v", len(topicsForInsert))
		return
	}

	// ***

	deleteTopicsWithChecks(t, storage)
}

// -----------------------------------------------------------------------

func Test_AllTopics(t *testing.T) {
	storage := instanceWithChecks(t)

	// ***

	topicsForInsert := generateFakeTopics()
	for i := range topicsForInsert {
		_, err := storage.InsertTopic(context.Background(), topicsForInsert[i])
		if err != nil {
			t.Errorf("Insert topic failed. Err: %v", err)
			return
		}
	}

	// ***

	allTopicsWithChecks(t, storage, topicsForInsert)

	// ***

	deleteTopicsWithChecks(t, storage)
}

func Test_AllTopics_V1(t *testing.T) {
	storage := instanceWithChecks(t)

	topicsForInsert := generateFakeTopics()
	insertTopicsWithChecks(t, storage, topicsForInsert)

	allTopicsWithChecks(t, storage, topicsForInsert)

	deleteTopicsWithChecks(t, storage)
}

func Test_Topic(t *testing.T) {
	storage := instanceWithChecks(t)

	// ***

	topicsForInsert := generateFakeTopics()
	insertTopicsWithChecks(t, storage, topicsForInsert)

	// ***

	readTopics := allTopicsWithChecks(t, storage, topicsForInsert)

	// ***

	randomTopicIndex := rand.Intn(len(readTopics))
	randomTopicIdr := readTopics[randomTopicIndex].Idr
	randomTopic, err := storage.Topic(context.Background(), randomTopicIdr)
	if err != nil {
		t.Errorf("Select topic failed. Err: %v", err)
		return
	}

	if !topicsForInsert.Contains(randomTopic) {
		t.Errorf("Topic not found in saved topics")
		return
	}

	// ***

	deleteTopicsWithChecks(t, storage)
}

func Test_RandomTopic(t *testing.T) {
	storage := instanceWithChecks(t)

	// ***

	topicsForInsert := generateFakeTopics()
	insertTopicsWithChecks(t, storage, topicsForInsert)

	// ***

	randomTopic, err := storage.RandomTopic(context.Background(), 1) // unstable test!
	if err != nil {
		t.Errorf("Select topic failed. Err: %v", err)
		return
	}

	if !topicsForInsert.Contains(randomTopic) {
		t.Errorf("Topic not found in saved topics")
		return
	}

	// ***

	deleteTopicsWithChecks(t, storage)
}

// delete
// -----------------------------------------------------------------------

func Test_DeleteTopic(t *testing.T) {
	// TODO: or enough?
}

// private help functions.
// -----------------------------------------------------------------------

func instanceWithChecks(t *testing.T) storage.Storage {
	storage, err := Instance()
	if err != nil {
		t.Fatalf("Failed to get storage object. Err: %v", err)
		return nil
	}
	return storage
}

func recordCountInTableWithChecks(t *testing.T, storage storage.Storage, tableName string) int {
	num, err := storage.RecordCountInTable(context.Background(), tableName)
	if err != nil {
		t.Errorf("Could not get the number of records in the table. Err: %v", err)
		return 0
	}
	return num
}

func insertTopicsWithChecks(t *testing.T, storage storage.Storage, topics domain.TopicList) {
	err := storage.InsertTopics(context.Background(), topics)
	if err != nil {
		t.Errorf("Insert topics failed. Err: %v", err)
		return
	}
}

func allTopicsWithChecks(t *testing.T, storage storage.Storage, storedTopics domain.TopicList) domain.TopicList {
	readTopics, err := storage.AllTopics(context.Background())
	if err != nil {
		t.Errorf("Select all topics failed. Err: %v", err)
		return nil
	}
	if len(readTopics) != len(storedTopics) {
		t.Errorf("Record count in the table is not equal %v", len(storedTopics))
		return nil
	}
	if !storedTopics.Eq(readTopics) {
		t.Errorf("Topics are not equal to each other")
		return nil
	}

	return readTopics
}

func deleteTopicsWithChecks(t *testing.T, storage storage.Storage) {
	err := storage.DeleteTopics(context.Background())
	if err != nil {
		t.Errorf("Failed to delete topics. Err: %v", err)
		return
	}
}

// ***

func generateFakeTopics() domain.TopicList {
	count := rand.Intn(10) + 10
	entities := make(domain.TopicList, count)
	for i := range entities {
		entities[i] = domain.Topic{
			Lang: rand.Intn(2),              // <--- max lang number
			Name: utility.RandomString(100), // <--- name length
		}
	}
	return entities
}
