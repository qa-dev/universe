package keeper

import (
	"gopkg.in/mgo.v2"
)

const SUBSCRIBERS_DB string = "subscribers"
const COLLECTION_PREFIX string = "stored_"

type MongoKeeper struct {
	mongo *mgo.Session
}

func NewMongoKeeper(session *mgo.Session) Keeper {
	session.SetMode(mgo.Monotonic, true)
	return &MongoKeeper{session}
}

func (k *MongoKeeper) StoreSubscriber(pluginName string, data interface{}) error {
	err := k.mongo.DB(SUBSCRIBERS_DB).C(COLLECTION_PREFIX + pluginName).Insert(data)
	return err
}

func (k *MongoKeeper) GetSubscribers(pluginName string, result interface{}) error {
	c := k.mongo.DB(SUBSCRIBERS_DB).C(COLLECTION_PREFIX + pluginName)

	err := c.Find(nil).All(result)
	return err
}
