package mongodb

import (
	"context"
	"time"

	"github.com/joobers/backend/pkg/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageModel struct{ Client *mongo.Client }

func (m *MessageModel) GetCollection(userID primitive.ObjectID) (*mongo.Collection, context.Context, context.CancelFunc) {
	collection := m.Client.Database("joobers").Collection("")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	return collection, ctx, cancel
}

func (m *MessageModel) Insert(authorID primitive.ObjectID, message models.Message, collection *mongo.Collection, ctx context.Context, cancel context.CancelFunc) (*mongo.InsertOneResult, error) {
	if collection == nil || ctx == nil || cancel == nil {
		collection, ctx, cancel = m.GetCollection(authorID)
		defer cancel()
	}

	result, err := collection.InsertOne(ctx, message)
	if err != nil {
		return nil, err
	}

	return result, err
}
