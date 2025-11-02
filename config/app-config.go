package config

type AppConfig struct {
	ChatAddress          string `env:"CHAT_ADDRESS"           env-default:"localhost:7072" yaml:"chatAddress"`
	ChatHandshakeTimeout int    `env:"CHAT_HANDSHAKE_TIMEOUT" env-default:"5"              yaml:"chatHandshakeTimeout"`
	ChatReadBufferSize   int    `env:"CHAT_READ_BUFFER_SIZE"  env-default:"1024"           yaml:"chatReadBufferSize"`
	ChatWriteBufferSize  int    `env:"CHAT_WRITE_BUFFER_SIZE" env-default:"1024"           yaml:"chatWriteBufferSize"`
}
