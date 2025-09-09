package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type HTTPConf struct {
	Addr               string   `json:"addr"`
	StaticDir          string   `json:"staticDir"`
	CORSAllowedOrigins []string `json:"corsAllowedOrigins"`
}

type KafkaConf struct {
	Broker  string `json:"broker"`
	Topic   string `json:"topic"`
	GroupID string `json:"groupID"`
}

type CacheConf struct {
	Limit int `json:"limit"`
}

type PublisherConf struct {
	Broker string `json:"broker"`
	Topic  string `json:"topic"`
	Count  int    `json:"count"`
}

type DBConf struct {
	DSN string `json:"dsn"`
}

type Config struct {
	HTTP      HTTPConf      `json:"http"`
	Kafka     KafkaConf     `json:"kafka"`
	Cache     CacheConf     `json:"cache"`
	Publisher PublisherConf `json:"publisher"`
	DB        DBConf        `json:"db"`
}

var (
	once sync.Once
	cfg  Config
)

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func load() {
	path := getEnv("CONFIG_PATH", filepath.FromSlash("internal/config/config.json"))

	cfg = Config{
		HTTP:      HTTPConf{Addr: ":8081", StaticDir: "./web", CORSAllowedOrigins: []string{"*"}},
		Kafka:     KafkaConf{Broker: "kafka:29092", Topic: "orders", GroupID: "order-group"},
		Cache:     CacheConf{Limit: 100},
		Publisher: PublisherConf{Broker: "localhost:9092", Topic: "orders", Count: 4},
		DB:        DBConf{DSN: ""},
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("config: could not read %s, using defaults/env: %v", path, err)
		return
	}

	var fileCfg Config
	if err := json.Unmarshal(data, &fileCfg); err != nil {
		log.Printf("config: invalid json in %s, using defaults/env: %v", path, err)
		return
	}

	if fileCfg.HTTP.Addr != "" {
		cfg.HTTP.Addr = fileCfg.HTTP.Addr
	}
	if fileCfg.HTTP.StaticDir != "" {
		cfg.HTTP.StaticDir = fileCfg.HTTP.StaticDir
	}
	if len(fileCfg.HTTP.CORSAllowedOrigins) > 0 {
		cfg.HTTP.CORSAllowedOrigins = fileCfg.HTTP.CORSAllowedOrigins
	}

	if fileCfg.Kafka.Broker != "" {
		cfg.Kafka.Broker = fileCfg.Kafka.Broker
	}
	if fileCfg.Kafka.Topic != "" {
		cfg.Kafka.Topic = fileCfg.Kafka.Topic
	}
	if fileCfg.Kafka.GroupID != "" {
		cfg.Kafka.GroupID = fileCfg.Kafka.GroupID
	}

	if fileCfg.Cache.Limit > 0 {
		cfg.Cache.Limit = fileCfg.Cache.Limit
	}

	if fileCfg.Publisher.Broker != "" {
		cfg.Publisher.Broker = fileCfg.Publisher.Broker
	}
	if fileCfg.Publisher.Topic != "" {
		cfg.Publisher.Topic = fileCfg.Publisher.Topic
	}
	if fileCfg.Publisher.Count > 0 {
		cfg.Publisher.Count = fileCfg.Publisher.Count
	}

	if fileCfg.DB.DSN != "" {
		cfg.DB.DSN = fileCfg.DB.DSN
	}
}

func ensureLoaded() { once.Do(load) }

func HTTPAddr() string {
	ensureLoaded()
	if v := os.Getenv("HTTP_ADDR"); v != "" {
		return v
	}
	return cfg.HTTP.Addr
}

func StaticDir() string {
	ensureLoaded()
	if v := os.Getenv("STATIC_DIR"); v != "" {
		return v
	}
	return cfg.HTTP.StaticDir
}

func CORSAllowedOrigins() []string {
	ensureLoaded()
	raw := os.Getenv("CORS_ALLOWED_ORIGINS")
	if strings.TrimSpace(raw) != "" {
		parts := strings.Split(raw, ",")
		out := make([]string, 0, len(parts))
		for _, p := range parts {
			s := strings.TrimSpace(p)
			if s != "" {
				out = append(out, s)
			}
		}
		if len(out) > 0 {
			return out
		}
	}
	if len(cfg.HTTP.CORSAllowedOrigins) == 0 {
		return []string{"*"}
	}
	return cfg.HTTP.CORSAllowedOrigins
}

func KafkaBroker() string {
	ensureLoaded()
	if v := os.Getenv("KAFKA_BROKER"); v != "" {
		return v
	}
	return cfg.Kafka.Broker
}
func KafkaTopic() string {
	ensureLoaded()
	if v := os.Getenv("KAFKA_TOPIC"); v != "" {
		return v
	}
	return cfg.Kafka.Topic
}
func KafkaGroupID() string {
	ensureLoaded()
	if v := os.Getenv("KAFKA_GROUP_ID"); v != "" {
		return v
	}
	return cfg.Kafka.GroupID
}

func CacheLimit() int {
	ensureLoaded()
	if v := os.Getenv("CACHE_LIMIT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	if cfg.Cache.Limit > 0 {
		return cfg.Cache.Limit
	}
	return 100
}

func PublisherBroker() string {
	ensureLoaded()
	if v := os.Getenv("PUB_BROKER"); v != "" {
		return v
	}
	return cfg.Publisher.Broker
}
func PublisherTopic() string {
	ensureLoaded()
	if v := os.Getenv("PUB_TOPIC"); v != "" {
		return v
	}
	return cfg.Publisher.Topic
}
func PublisherCount() int {
	ensureLoaded()
	if v := os.Getenv("PUB_COUNT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	if cfg.Publisher.Count > 0 {
		return cfg.Publisher.Count
	}
	return 4
}

func DBDSN() string {
	ensureLoaded()
	if v := os.Getenv("DB_URL"); v != "" {
		return v
	}
	if cfg.DB.DSN != "" {
		return cfg.DB.DSN
	}
	return "postgres://myadmin:mypassword@pgdb:5432/mydatabase"
}
