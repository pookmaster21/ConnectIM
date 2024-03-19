package db

import (
	"context"
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	. "github.com/pookmaster21/ConnectIM/types"
)

type db struct {
	coll   *mongo.Collection
	client *mongo.Client
}

var (
	DB     db
	logger *Logger
)

func InitDB(ctx context.Context, uri string) {
	logger = NewLogger()

	var err error
	DB.client, err = mongo.Connect(
		ctx,
		options.Client().ApplyURI(uri),
	)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	DB.coll = DB.client.Database("ConnectIM").Collection("Users")
	_, err = DB.coll.Aggregate(ctx, bson.A{})
	if err != nil {
		logger.Error(err.Error())
		return
	}
	logger.Info("Inited the DB")
}

func (d *db) InsertUser(ctx context.Context, user *User) {
	_, err := d.coll.InsertOne(ctx, bson.M{
		"username": user.Username,
		"password": user.Password,
		"whatsapp": user.Whatsapp,
		"telegram": user.Telegram,
		"discord":  user.Discord,
		"prefered": user.Prefered,
	})
	if err != nil {
		logger.Error(err.Error())
		return
	}
}

func (d *db) FindUser(ctx context.Context, fields []string, info []any) *User {
	result := bson.M{}
	for i := range fields {
		result[fields[i]] = info[i]
	}

	err := d.coll.FindOne(ctx, result).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	if err != nil {
		logger.Error(err.Error())
		return nil
	}

	jsondata, err := json.Marshal(result)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}

	var user User
	err = json.Unmarshal(jsondata, &user)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	return &user
}

func (d *db) Close(ctx context.Context) {
	err := d.client.Disconnect(ctx)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Info("Closed the DB")
}

func (d *db) DeleteUser(ctx context.Context, fields []string, info []any) {
	result := bson.M{}
	for i := range fields {
		result[fields[i]] = info[i]
	}

	_, err := d.coll.DeleteOne(ctx, result)
	if err != nil {
		logger.Error(err.Error())
		return
	}
}
