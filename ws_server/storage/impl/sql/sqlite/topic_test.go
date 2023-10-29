package sqlite

import (
	"context"
	"ilserver/domain"
	"reflect"
	"testing"
)

func Test_InsertTopic(t *testing.T) {
	storage, err := Instance()
	if err != nil {
		t.Error()
		return
	}

	// ***

	topics := []domain.Topic{
		{Lang: 1, Name: "test_name_1"},
		{Lang: 2, Name: "test_name_2"},
		{Lang: 3, Name: "test_name_2"},
		//***
	}

	for i := range topics {
		_, err = storage.InsertTopic(context.Background(), topics[i])
		if err != nil {
			t.Errorf("Insert topic failed. Err: %v", err)
			return
		}
	}

	// ***

	num, err := storage.RecordCountInTable(context.Background(), "Topics")
	if err != nil {
		t.Errorf("Get record count in table failed. Err: %v", err)
		return
	}

	if num != len(topics) {
		t.Errorf("Record count in the table is not equal %v", len(topics))
		return
	}
}

func Test_AllTopics(t *testing.T) {
	storage, err := Instance()
	if err != nil {
		t.Error()
		return
	}

	topicsForInsert := []domain.Topic{
		{Lang: 1, Name: "test_name_1"},
		{Lang: 2, Name: "test_name_2"},
		{Lang: 3, Name: "test_name_3"},
		{Lang: 4, Name: "test_name_4"},
		{Lang: 5, Name: "test_name_5"},
		//***
	}

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

	// TODO:
	if !reflect.DeepEqual(topicsForInsert, readTopics) {
		t.Errorf("Topics are not equal to each other")
		return
	}
}
