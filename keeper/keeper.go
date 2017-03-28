package keeper

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

const (
	SubscribersDatabase string = "subscribers"
	MaxFailCount        int    = 5
)

type Keeper struct {
	mongo        *mgo.Session
	databaseName string
}

type doFunc func() error

func NewKeeper(session *mgo.Session) *Keeper {
	return &Keeper{session, SubscribersDatabase}
}

func (k *Keeper) SetCustomDatabaseName(name string) {
	k.databaseName = name
}

func (k *Keeper) StoreSubscriber(pluginName string, data interface{}) error {
	err := k.doAction(func() error {
		return k.mongo.DB(k.databaseName).C("stored_" + pluginName).Insert(data)
	})
	return err
}

func (k *Keeper) RemoveSubscriber(pluginName string, data interface{}) error {
	err := k.doAction(func() error {
		return k.mongo.DB(k.databaseName).C("stored_" + pluginName).Remove(data)
	})
	return err
}

func (k *Keeper) GetSubscribers(pluginName string, data interface{}) error {
	err := k.doAction(func() error {
		return k.mongo.DB(k.databaseName).C("stored_" + pluginName).Find(nil).All(data)
	})
	return err
}

func (k *Keeper) doAction(fn doFunc) error {
	err := errors.New("1")
	failCount := 0
	for err != nil {
		err = fn()
		if err == nil {
			if failCount > 0 {
				log.Infof("Connected in %d tries", failCount)
			}
			break
		}
		failCount++
		log.Info(err)
		log.Info("Mongo disconnected. Waiting...")
		if failCount > MaxFailCount {
			log.Infof("Fail after %d tries", failCount)
			return err
		}
		k.mongo.Refresh()
	}
	return err
}
