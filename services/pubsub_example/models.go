package pubsub_example

const topic = "pubsub-example"

type Thing struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}
