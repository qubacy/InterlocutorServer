package sqlite

import (
	"context"
	"ilserver/domain"
	"ilserver/utility"
	"math/rand"
	"testing"
)

func Test_InsertTopic(t *testing.T) {
	storage, err := Instance()
	if err != nil {
		t.Error()
		return
	}

	// ***

	topicsForInsert := generateFakeTopics()
	for i := range topicsForInsert {
		_, err = storage.InsertTopic(context.Background(), topicsForInsert[i])
		if err != nil {
			t.Errorf("Insert topic failed. Err: %v", err)
			return
		}
	}

	// ***

	num, err := storage.RecordCountInTable(context.Background(), "Topics")
	if err != nil {
		t.Errorf("Could not get the number of records in the table. Err: %v", err)
		return
	}

	if num != len(topicsForInsert) {
		t.Errorf("Record count in the table is not equal %v", len(topicsForInsert))
		return
	}

	// ***

	// testing?
	deleteTopicsWithChecks(t, storage)
}

func Test_InsertTopics(t *testing.T) {
	storage, err := Instance()
	if err != nil {
		t.Error()
		return
	}
	topicsForInsert := generateFakeTopics()
	insertTopicsWithChecks(t, storage, topicsForInsert)

	// ***

	num, err := storage.RecordCountInTable(context.Background(), "Topics")
	if err != nil {
		t.Errorf("Could not get the number of records in the table. Err: %v", err)
		return
	}

	if num != len(topicsForInsert) {
		t.Errorf("Record count in the table is not equal %v", len(topicsForInsert))
		return
	}

	// ***

	deleteTopicsWithChecks(t, storage)
}

// -----------------------------------------------------------------------

func Test_AllTopics(t *testing.T) {
	storage, err := Instance()
	if err != nil {
		t.Error()
		return
	}

	topicsForInsert := generateFakeTopics()
	for i := range topicsForInsert {
		_, err = storage.InsertTopic(context.Background(), topicsForInsert[i])
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
	storage, err := Instance()
	if err != nil {
		t.Error()
		return
	}

	topicsForInsert := generateFakeTopics()
	err = storage.InsertTopics(context.Background(), topicsForInsert)
	if err != nil {
		t.Errorf("Insert topics failed. Err: %v", err)
		return
	}

	// ***

	allTopicsWithChecks(t, storage, topicsForInsert)

	// ***

	deleteTopicsWithChecks(t, storage)
}

func Test_Topic(t *testing.T) {
	storage, err := Instance()
	if err != nil {
		t.Error()
		return
	}

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

	if !domain.Contains(topicsForInsert, randomTopic) {
		t.Errorf("Topic not found in saved topics")
		return
	}

	// ***

	deleteTopicsWithChecks(t, storage)
}

func Test_RandomTopic(t *testing.T) {
	storage, err := Instance()
	if err != nil {
		t.Error()
		return
	}

	topicsForInsert := generateFakeTopics()
	insertTopicsWithChecks(t, storage, topicsForInsert)

	// ***

	randomTopic, err := storage.RandomTopic(context.Background(), 1) // unstable test!
	if err != nil {
		t.Errorf("Select topic failed. Err: %v", err)
		return
	}

	if !domain.Contains(topicsForInsert, randomTopic) {
		t.Errorf("Topic not found in saved topics")
		return
	}

	// ***

	deleteTopicsWithChecks(t, storage)
}

// delete
// -----------------------------------------------------------------------

func Test_DeleteTopic(t *testing.T) {
	// TODO:
}

// private help functions.
// -----------------------------------------------------------------------

func insertTopicsWithChecks(t *testing.T, storage *Storage, topics []domain.Topic) {
	err := storage.InsertTopics(context.Background(), topics)
	if err != nil {
		t.Errorf("Insert topics failed. Err: %v", err)
		return
	}
}

func allTopicsWithChecks(t *testing.T, storage *Storage, storedTopics []domain.Topic) []domain.Topic {
	readTopics, err := storage.AllTopics(context.Background())
	if err != nil {
		t.Errorf("Select all topics failed. Err: %v", err)
		return nil
	}
	if len(readTopics) != len(storedTopics) {
		t.Errorf("Record count in the table is not equal %v", len(storedTopics))
		return nil
	}
	if !domain.IsSomeEqual(storedTopics, readTopics) {
		t.Errorf("Topics are not equal to each other")
		return nil
	}

	return readTopics
}

func deleteTopicsWithChecks(t *testing.T, storage *Storage) {
	err := storage.DeleteTopics(context.Background())
	if err != nil {
		t.Errorf("Failed to delete topics. Err: %v", err)
		return
	}
}

// ***

func generateFakeTopics() []domain.Topic {
	topicCount := rand.Intn(10) + 10
	topics := make([]domain.Topic, topicCount)
	for i := range topics {
		topics[i] = domain.Topic{
			Lang: rand.Intn(2),               // <--- max lang number
			Name: utility.RandomString(1000), // <--- name length
		}
	}
	return topics
}
