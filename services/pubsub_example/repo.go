package pubsub_example

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/pubsub"
)

var _ Repo = (*repo)(nil)

type Repo interface {
	// Getter
	// Setter
}

// type (
// 	Getter interface {
// 		Get(ctx context.Context, get_id) (get_resp, error)
// 	}
// 	GetterFunc func(ctx context.Context, get_id) (get_resp, error)
// )

// func (f GetterFunc) Get(ctx context.Context, get_id) (get_resp, error) {
// 	return f(ctx)
// }

// type (
// 	Setter interface {
// 		Set(ctx context.Context, set_id) error
// 	}
// 	SetterFunc func(ctx context.Context, set_id) error
// )

// func (f SetterFunc) Set(ctx context.Context, set_id) error {
// 	return f(ctx)
// }

type repo struct {
	pb *pubsub.Client
}

func New() *repo {
	ctx := context.Background()
	pb, err := pubsub.NewClient(ctx, projectID())
	if err != nil {
		log.Fatal(err)
	}
	// Create the topic if it doesn't exist.
	if exists, err := pb.Topic(topic).Exists(ctx); err != nil {
		log.Fatal(err)
	} else if !exists {
		if _, err := pb.CreateTopic(ctx, topic); err != nil {
			log.Fatal(err)
		}
	}
	return &repo{pb}
}

func (r *repo) Pub(ctx context.Context, t *Thing) error {
	b := &bytes.Buffer{}
	if err := json.NewEncoder(b).Encode(t); err != nil {
		return err
	}
	result := r.pb.Topic(topic).Publish(ctx, &pubsub.Message{
		Data: b.Bytes(),
	})
	_, err := result.Get(ctx)
	return err
}

func (r *repo) Sub(ctx context.Context) error {
	sub, err := r.pb.CreateSubscription(ctx, "test", pubsub.SubscriptionConfig{Topic: r.pb.Topic(topic)})
	if err != nil {
		return err
	}
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
	})
	if err != nil {
		return err
	}
	return nil
}

// ProjectID will attempt to get the Google Cloud Project ID with the following rules:
// 1) Use the metadata API to get ID, this will only work in Google Cloud
// 2) Any failure or timeout (3s) will presume that the code is running outside the cloud
// in which case a default project ID is returned.
func projectID() string {
	const defaultID = "project-id"
	req, err := http.NewRequest(http.MethodGet, "http://metadata.google.internal/computeMetadata/v1/project/project-id", nil)
	if err != nil {
		return defaultID
	}
	req.Header.Add("Metadata-Flavor", "Google")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return defaultID
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return defaultID
		}
		return string(b)
	}
	return defaultID
}
