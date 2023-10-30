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
	err = storage.DeleteTopics(context.Background())
	if err != nil {
		t.Errorf("Failed to delete topics. Err: %v", err)
		return
	}
}

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

	readTopics, err := storage.AllTopics(context.Background())
	if len(readTopics) != len(topicsForInsert) {
		t.Errorf("Record count in the table is not equal %v", len(topicsForInsert))
		return
	}

	if !domain.IsSomeEqual(topicsForInsert, readTopics) {
		t.Errorf("Topics are not equal to each other")
		return
	}
}

func Test_Topic(t *testing.T) {

}

// private
// -----------------------------------------------------------------------

func generateFakeTopics() []domain.Topic {
	topicCount := rand.Intn(10) + 10
	topics := make([]domain.Topic, topicCount)
	for i := range topics {
		topics[i] = domain.Topic{
			Lang: rand.Intn(10),              // <--- max lang number
			Name: utility.RandomString(1000), // <--- name length
		}
	}
	return topics
}
