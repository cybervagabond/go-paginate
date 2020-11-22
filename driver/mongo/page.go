package mongo

import (
	"context"
	"github.com/raphaelvigee/go-paginate/driver"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoDriverPage struct {
	coll       *mongo.Collection
	pageInfo   driver.PageInfo
	cursorFunc func(i int64) (interface{}, error)
}

var _ driver.Page = (*mongoDriverPage)(nil)

func (m mongoDriverPage) Cursor(i int64) (interface{}, error) {
	return m.cursorFunc(i)
}

func (m mongoDriverPage) Query(dst interface{}) error {
	if m.coll == nil {
		return nil
	}

	// TODO: get filter from user
	filter := bson.M{}

	// run mongodb query
	res, err := m.coll.Find(context.TODO(), filter)
	if err != nil {
		return err
	}

	// decode find result into destination
	err = res.All(context.TODO(), dst)
	if err != nil {
		return err
	}

	return nil
}

func (m mongoDriverPage) Count() (int64, error) {
	if m.coll == nil {
		return 0, nil
	}

	var filter bson.D

	c, err := m.coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, err
	}

	return c, err
}

func (m mongoDriverPage) Info() driver.PageInfo {
	return m.pageInfo
}
