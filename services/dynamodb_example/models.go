package dynamodb_example

const tableName = "things"

type Thing struct {
	ID string `json:"id"`
}
