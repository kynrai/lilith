package dynamodb

import (
	"context"
	"errors"
	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	awsdynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var _ Repo = (*dynamodb)(nil)

type Repo interface {
	Getter
	Putter
}

type (
	Getter interface {
		Get(ctx context.Context, tableName string, key, v interface{}) error
	}
	GetterFunc func(ctx context.Context, tableName string, key, v interface{}) error
)

func (f GetterFunc) Get(ctx context.Context, tableName string, key, v interface{}) error {
	return f(ctx, tableName, key, v)
}

type (
	Putter interface {
		Put(ctx context.Context, tableName string, v interface{}) error
	}
	PutterFunc func(ctx context.Context, tableName string, v interface{}) error
)

func (f PutterFunc) Put(ctx context.Context, tableName string, v interface{}) error {
	return f(ctx, tableName, v)
}

type dynamodb struct {
	client   *awsdynamodb.DynamoDB
	endpoint string
	region   string
}

type dynamodbOption func(*dynamodb)

func New(opts ...dynamodbOption) *dynamodb {
	d := &dynamodb{}
	// TODO make session for all use cases, local testing.
	// URL based and running in AWS environment
	d.client = awsdynamodb.New(nil)
	return d
}

func (d *dynamodb) Get(ctx context.Context, tableName string, key, v interface{}) error {
	akey, err := dynamodbattribute.MarshalMap(key)
	if err != nil {
		return err
	}
	result, err := d.client.GetItemWithContext(ctx, &awsdynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       akey,
	})
	if err != nil {
		return err
	}
	return dynamodbattribute.UnmarshalMap(result.Item, v)
}

func (d *dynamodb) GetMulti(ctx context.Context, tableName string) error {
	return nil
}

func (d *dynamodb) Put(ctx context.Context, tableName string, v interface{}) error {
	av, err := dynamodbattribute.MarshalMap(v)
	if err != nil {
		return err
	}
	_, err = d.client.PutItemWithContext(ctx, &awsdynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	})
	return err
}

func (d *dynamodb) PutMulti(ctx context.Context, tableName string, vals interface{}) error {
	// Check vals are a a slice of values
	v := reflect.ValueOf(vals)
	if v.Kind() != reflect.Slice {
		return errors.New("vals is not a slice")
	}

	// Number of items
	l := v.Len()
	writeBatch := make([]*awsdynamodb.WriteRequest, l)
	for i := 0; i < l; i++ {
		// Marshal items and put into dynamo write requets
		item, err := dynamodbattribute.MarshalMap(v.Index(i).Interface())
		if err != nil {
			return err
		}
		writeBatch[i] = &awsdynamodb.WriteRequest{PutRequest: &awsdynamodb.PutRequest{Item: item}}
	}

	unprocessed := []*awsdynamodb.WriteRequest{}
	batcher := NewBatcher(writeBatch, 0)

	for !batcher.Done() {
		puts := &awsdynamodb.BatchWriteItemInput{
			RequestItems: map[string][]*awsdynamodb.WriteRequest{
				tableName: batcher.Next(unprocessed),
			},
		}
		resp, err := d.client.BatchWriteItemWithContext(ctx, puts)
		if err != nil {
			return err
		}
		unprocessed = resp.UnprocessedItems[tableName]
	}

	return nil
}

func InterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

func (d *dynamodb) Del(ctx context.Context) error {
	return nil
}

func (d *dynamodb) DelMulti(ctx context.Context) error {
	return nil
}

func (d *dynamodb) Client() *awsdynamodb.DynamoDB {
	return d.client
}
