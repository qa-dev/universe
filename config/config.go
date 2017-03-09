package config

type App struct {
	Host string `json:"host"`
	Port uint   `json:"port"`
}
type Rmq struct {
	Uri        string `json:"uri"`
	EventQueue string `json:"event_queue"`
}

type Config struct {
	App App `json:"app"`
	Rmq Rmq `json:"rmq"`
}
