package dynamodb_example

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/kynrai/lilith/pkg/dynamodb_example/models"
)

const tableName = "things"

var _ Repo = (*repo)(nil)

type Repo interface {
	Getter
	MultiGetter
	Putter
}

type (
	Getter interface {
		Get(ctx context.Context, id string) (*models.Thing, error)
	}
	GetterFunc func(ctx context.Context, id string) (*models.Thing, error)
)

func (f GetterFunc) Get(ctx context.Context, id string) (*models.Thing, error) {
	return f(ctx, id)
}

type (
	MultiGetter interface {
		GetMulti(ctx context.Context, ids ...string) ([]*models.Thing, error)
	}
	MultiGetterFunc func(ctx context.Context, ids ...string) ([]*models.Thing, error)
)

func (f MultiGetterFunc) GetMulti(ctx context.Context, ids ...string) ([]*models.Thing, error) {
	return f(ctx, ids...)
}

type (
	Putter interface {
		Put(ctx context.Context, t *models.Thing) error
	}
	PutterFunc func(ctx context.Context, t *models.Thing) error
)

func (f PutterFunc) Put(ctx context.Context, t *models.Thing) error {
	return f(ctx, t)
}

type repo struct {
	db dynamodb.DynamoDB
}

func New(db dynamodb.DynamoDB) *repo {
	return &repo{db}
}

func (r *repo) Get(ctx context.Context, id string) (*models.Thing, error) {
	t := models.Thing{}
	getResp, err := r.db.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": &dynamodb.AttributeValue{
				S: aws.String(id),
			},
		},
		TableName: aws.String("tableName"),
	})
	if err != nil {
		return nil, err
	}
	return &t, dynamodbattribute.UnmarshalMap(getResp.Item, &t)
}

func (r *repo) GetMulti(ctx context.Context, ids ...string) ([]*models.Thing, error) {
	keys := make([]map[string]*dynamodb.AttributeValue, 0, len(ids))
	for _, v := range ids {
		keys = append(keys, map[string]*dynamodb.AttributeValue{
			"id": &dynamodb.AttributeValue{
				S: aws.String(v),
			},
		})
	}
	getResps, err := r.db.BatchGetItemWithContext(ctx, &dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			tableName: &dynamodb.KeysAndAttributes{
				Keys: keys,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	var things []*models.Thing
	return things, dynamodbattribute.UnmarshalListOfMaps(getResps.Responses[tableName], &things)
}

func (r *repo) Put(ctx context.Context, t *models.Thing) error {
	return nil
}
