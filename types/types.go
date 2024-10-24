package types

import "time"

type Config struct {
	ChatID         int64         `json:"chat_id"`
	Token          string        `json:"token"`
	Interval       time.Duration `json:"interval"`
	NotifyInterval time.Duration `json:"notify_interval"`
	RpcList        []string      `json:"rpc_list"`
}

type Chan struct {
	Err error
	Rpc string
}

type Notifys map[string]time.Time
