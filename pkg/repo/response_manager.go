package repo

import (
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

const kResponses = "responses"

type MongoResponseSaver struct {
	responses *mongo.Collection
}

func NewMongoResponseSaver(conn *mongo.Client) ResponseSaver {
	return &MongoResponseSaver{
		responses: conn.Database(kDatabase).Collection(kResponses),
	}
}

func (s *MongoResponseSaver) Save(resp *http.Response) (string, error) {
	return "", nil
}

func (s *MongoResponseSaver) Get(id int64) (*http.Response, error) {
	return nil, nil
}

func (s *MongoResponseSaver) List() ([]*http.Response, error) {
	return nil, nil
}
