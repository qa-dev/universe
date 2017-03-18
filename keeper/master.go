package keeper

import "gopkg.in/mgo.v2"

type MongoMaster struct {
	mongo *mgo.Session
}

func NewMongoMaster(session *mgo.Session) *MongoMaster {
	return &MongoMaster{session}
}

func (m *MongoMaster) GetCollection(pluginName string) *mgo.Collection {
	return m.mongo.DB("subscribers").C("subscribers_" + pluginName)
}
