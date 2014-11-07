package config

type ConfigRedisServer struct {
}

type ConfigRedis struct {
    Breaker ConfigBreaker
    Servers map[string]*ConfigRedisServer
}
