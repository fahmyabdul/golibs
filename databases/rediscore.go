package databases

import (
	"log"
	"os"

	"github.com/gomodule/redigo/redis"
)

// RedisCore :
type RedisCore struct {
	Host      string      `json:"host"`
	Auth      string      `json:"auth"`
	DB        int         `json:"db"`
	MaxIdle   int         `json:"max_idle"`
	MaxActive int         `json:"max_active"`
	Logger    *log.Logger // Optional
}

// NewRedis :
func NewRedis(host, auth string, dbIndex, maxIdle, maxActive int, nLog *log.Logger) (*redis.Pool, error) {
	if nLog == nil {
		nLog = log.New(os.Stderr, "", log.LstdFlags)
	}

	redisConf := RedisCore{
		Host:      host,
		Auth:      auth,
		DB:        dbIndex,
		MaxIdle:   maxIdle,
		MaxActive: maxActive,
		Logger:    nLog,
	}

	redisPool, err := redisConf.InitRedis()
	if err != nil {
		return nil, err
	}

	return redisPool, nil
}

// InitRedis :
func (p *RedisCore) InitRedis() (*redis.Pool, error) {
	conn, err := redis.Dial("tcp", p.Host)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	if _, err := conn.Do("AUTH", p.Auth); err != nil {
		conn.Close()
		return nil, err
	}

	if _, err := conn.Do("SELECT", p.DB); err != nil {
		conn.Close()
		return nil, err
	}

	redPool := &redis.Pool{
		MaxIdle:   p.MaxIdle,
		MaxActive: p.MaxActive, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", p.Host)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", p.Auth); err != nil {
				c.Close()
				return nil, err
			}

			if _, err := c.Do("SELECT", p.DB); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
	}

	return redPool, nil
}
