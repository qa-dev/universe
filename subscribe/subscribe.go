package subscribe

type Subscribe struct {
	EventName   string `json:"event_name"`
	WebHookPath string `json:"webhook"`
}

type SubscribeData struct {
	EventName  string      `json:"event_name"`
	PluginName string      `json:"plugin_name"`
	Data       interface{} `json:"data"`
}

type UnsubscribeData struct {
	EventName  string      `json:"event_name"`
	PluginName string      `json:"plugin_name"`
	Data       interface{} `json:"data"`
}
