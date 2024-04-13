//go:build k8s

package config

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(vchat-mysql:13316)/v_chat",
	},
	Redis: RedisConfig{
		Addr: "vchat-redis:6379",
	},
}
