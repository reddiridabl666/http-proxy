package repo

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (s *MongoResponseSaver) Save(requestId string, resp *http.Response) (string, error) {
	requestObjectId, err := primitive.ObjectIDFromHex(requestId)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	resp.Body = io.NopCloser(bytes.NewReader(body))

	res, err := s.responses.InsertOne(context.Background(), bson.M{
		"code":       resp.StatusCode,
		"message":    resp.Status[strings.Index(resp.Status, " ")+1:],
		"headers":    toBson(resp.Header),
		"request_id": requestObjectId,
		"body":       string(body),
	})
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (s *MongoResponseSaver) Get(id string) (*ResponseData, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	res := s.responses.FindOne(context.Background(), bson.D{{Key: "_id", Value: objectId}})
	value := &ResponseData{}

	err = res.Decode(value)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (s *MongoResponseSaver) GetByRequest(requestId string) (*ResponseData, error) {
	objectId, err := primitive.ObjectIDFromHex(requestId)
	if err != nil {
		return nil, err
	}

	res := s.responses.FindOne(context.Background(), bson.D{{Key: "request_id", Value: objectId}})
	value := &ResponseData{}

	err = res.Decode(value)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (s *MongoResponseSaver) List(limit int64) ([]*ResponseData, error) {
	ctx := context.Background()

	opts := options.Find().SetLimit(limit).SetSort(bson.D{{Key: "_id", Value: -1}})
	cursor, err := s.responses.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}

	res := make([]*ResponseData, 0, limit/2)

	err = cursor.All(ctx, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
