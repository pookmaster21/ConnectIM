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
	logger *Logger
}

var DB db

func Init_db(ctx context.Context, uri string, logger *Logger) {
	DB.logger = logger

	var err error
	DB.client, err = mongo.Connect(
		ctx,
		options.Client().ApplyURI(uri),
	)
	if err != nil {
		logger.Fatal(err.Error())
		return
	}

	DB.coll = DB.client.Database("ConnectIM").Collection("Users")
	_, err = DB.coll.Aggregate(ctx, bson.A{})
	if err != nil {
		logger.Fatal(err.Error())
		return
	}
	logger.Info("Inited the DB")
}

func (d *db) Insert_user(ctx context.Context, user *User) {
	_, err := d.coll.InsertOne(ctx, bson.M{
		"username": user.Username,
		"password": user.Password,
		"whatsapp": user.Whatsapp,
		"telegram": user.Telegram,
		"discord":  user.Discord,
		"prefred":  user.Prefered,
	})
	if err != nil {
		d.logger.Fatal(err.Error())
	}
}

func (d *db) Find_user(ctx context.Context, fields []string, info []any) *User {
	result := bson.M{}
	for i := range fields {
		result[fields[i]] = info[i]
	}

	err := d.coll.FindOne(ctx, result).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	if err != nil {
		d.logger.Fatal(err.Error())
	}

	jsondata, err := json.Marshal(result)
	if err != nil {
		d.logger.Fatal(err.Error())
	}

	var user User
	err = json.Unmarshal(jsondata, &user)
	if err != nil {
		d.logger.Fatal(err.Error())
	}
	return &user
}

func (d *db) Close(ctx context.Context) {
	err := d.client.Disconnect(ctx)
	if err != nil {
		d.logger.Error(err.Error())
		return
	}

	d.logger.Info("Closed the DB")
}
