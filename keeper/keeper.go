package keeper

type Keeper interface {
	StoreSubscriber(pluginName string, data interface{}) error
	GetSubscribers(pluginName string, result interface{}) error
}
