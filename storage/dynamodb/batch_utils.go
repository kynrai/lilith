package dynamodb

import awsdynamodb "github.com/aws/aws-sdk-go/service/dynamodb"

const BatchWriteLimit = 25

type batcher struct {
	Batch     []*awsdynamodb.WriteRequest
	Offset    int
	itemsDone int
}

func NewBatcher(batch []*awsdynamodb.WriteRequest, offset int) *batcher {
	return &batcher{Batch: batch[offset:]}
}

func (b *batcher) Next(unprocessed []*awsdynamodb.WriteRequest) []*awsdynamodb.WriteRequest {
	limit := BatchWriteLimit - len(unprocessed)
	if len(b.Batch)+len(unprocessed) <= 25 {
		return append(b.Batch, unprocessed...)
	}

	l := len(b.Batch)
	if l <= 25 {
		return b.Batch
	}
	end := b.Offset + limit
	if end > l {
		end = l
	}
	defer func() { b.Offset += limit }()
	return append(b.Batch[b.Offset:end], unprocessed...)
}

func (b *batcher) Done() bool {
	return false
}
