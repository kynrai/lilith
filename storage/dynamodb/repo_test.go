package dynamodb

import (
	"testing"

	awsdynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func TestGetWriteBatch(t *testing.T) {
	t.Parallel()
	type writeable struct{ ID int }

	writeBatch := make([]*awsdynamodb.WriteRequest, 60)
	for i := 0; i < 60; i++ {
		item, err := dynamodbattribute.MarshalMap(&writeable{i})
		if err != nil {
			t.Fatal(err)
		}
		writeBatch[i] = &awsdynamodb.WriteRequest{PutRequest: &awsdynamodb.PutRequest{Item: item}}
	}

	unprocessed := make([]*awsdynamodb.WriteRequest, 10)
	for i := 0; i < 10; i++ {
		item, err := dynamodbattribute.MarshalMap(&writeable{i})
		if err != nil {
			t.Fatal(err)
		}
		unprocessed[i] = &awsdynamodb.WriteRequest{PutRequest: &awsdynamodb.PutRequest{Item: item}}
	}

	for _, tc := range []struct {
		name        string
		batch       []*awsdynamodb.WriteRequest
		offset      int
		wantLen     int
		unprocessed []*awsdynamodb.WriteRequest
	}{
		{
			name:    "0-25",
			batch:   writeBatch,
			offset:  0,
			wantLen: 25,
		},
		{
			name:    "25-50",
			batch:   writeBatch,
			offset:  25,
			wantLen: 25,
		},
		{
			name:    "50-60 out of bounds",
			batch:   writeBatch,
			offset:  50,
			wantLen: 10,
		},
		{
			name:    "60-60 out of bounds",
			batch:   writeBatch,
			offset:  60,
			wantLen: 0,
		},
		{
			name:        "0-25 2 unprocessed",
			batch:       writeBatch,
			offset:      0,
			wantLen:     25,
			unprocessed: unprocessed,
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			batcher := NewBatcher(tc.batch, tc.offset)
			writeBatch := batcher.Next(tc.unprocessed)
			if got := len(writeBatch); got != tc.wantLen {
				t.Fatalf("got: %d, want: %d", got, tc.wantLen)
			}
		})
	}
}

func TestGetWriteBatchNext(t *testing.T) {
	t.Parallel()
	type writeable struct {
		ID int
	}
	writeBatch := make([]*awsdynamodb.WriteRequest, 60)
	for i := 0; i < 60; i++ {
		item, err := dynamodbattribute.MarshalMap(&writeable{i})
		if err != nil {
			t.Fatal(err)
		}
		writeBatch[i] = &awsdynamodb.WriteRequest{PutRequest: &awsdynamodb.PutRequest{Item: item}}
	}
	for _, tc := range []struct {
		name        string
		batch       []*awsdynamodb.WriteRequest
		offset      int
		wantLen     int
		unprocessed []*awsdynamodb.WriteRequest
	}{
		{
			name:    "0-25",
			batch:   writeBatch,
			offset:  25,
			wantLen: 10,
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			batcher := NewBatcher(tc.batch, tc.offset)
			writeBatch := batcher.Next(nil)
			writeBatch = batcher.Next(nil)
			if got := len(writeBatch); got != tc.wantLen {
				t.Fatalf("got: %d, want: %d", got, tc.wantLen)
			}
		})
	}
}
