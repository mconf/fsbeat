// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period    time.Duration `config:"period"`
	MaxBuffer int           `config:"max.buffer"`
	FSServer  string        `config:"freeswitch.server"`
	FSPort    string        `config:"freeswitch.port"`
	FSAuth    string        `config:"freeswitch.auth"`
	FSEvents  string        `config:"freeswitch.events"`
}

var DefaultConfig = Config{
	Period:    1 * time.Second,
	MaxBuffer: 20,
	FSServer:  "localhost",
	FSPort:    "8021",
	FSAuth:    "ClueCon",
	FSEvents:  "all",
}
