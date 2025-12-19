package database

import (
	"context"
	"time"

	"github.com/qiniu/qmgo"
)

var MongoDB *qmgo.Database

func Connect(uri, dbName string) (*qmgo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := qmgo.NewClient(ctx, &qmgo.Config{
		Uri: uri,
	})

	if err != nil {
		return nil, err
	}

	return client.Database(dbName), nil
}