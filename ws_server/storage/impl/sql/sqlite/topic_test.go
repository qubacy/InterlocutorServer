package sqlite

import (
	"context"
	"ilserver/domain"
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
		storage.InsertTopic(context.Background(), topics[i])
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
