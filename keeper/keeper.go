package keeper

import (
	"gopkg.in/mgo.v2"
)

const SubscribersDatabase string = "subscribers"

type Keeper struct {
	mongo        *mgo.Session
	databaseName string
}

func NewKeeper(session *mgo.Session) *Keeper {
	return &Keeper{session, SubscribersDatabase}
}

func (k *Keeper) SetCustomDatabaseName(name string) {
	k.databaseName = name
}

func (k *Keeper) StoreSubscriber(pluginName string, data interface{}) error {
	return k.mongo.DB(k.databaseName).C("stored_" + pluginName).Insert(data)
}

func (k *Keeper) RemoveSubscriber(pluginName string, data interface{}) error {
	return k.mongo.DB(k.databaseName).C("stored_" + pluginName).Remove(data)
}

func (k *Keeper) GetSubscribers(pluginName string, data interface{}) error {
	return k.mongo.DB(k.databaseName).C("stored_" + pluginName).Find(nil).All(data)
}
