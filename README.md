# Сервис обработки и отображения заказов

Этот проект представляет собой демонстрационный микросервис на Go, который получает данные о заказах из очереди сообщений (Kafka), сохраняет их в базу данных (PostgreSQL), кэширует в памяти и предоставляет через HTTP API с простым веб-интерфейсом.

## 🚀 Как запустить проект

Для запуска проекта вам понадобится установленный **Docker**

1.  **Клонируйте репозиторий:**
    ```bash
    git clone <URL вашего репозитория>
    cd <имя папки проекта>
    ```

2.  **Запустите все сервисы с помощью Docker Compose:**
    Эта команда соберет Go-приложение, запустит контейнеры с PostgreSQL, Zookeeper и Kafka, а также сам сервис.

    ```bash
    docker-compose up --build
    ```

    После выполнения команды будут доступны:
  -   **Веб-интерфейс сервиса**: [http://localhost:8081/](http://localhost:8081/) (порт можно изменить переменной `HTTP_ADDR`)
    -   **База данных PostgreSQL**: `localhost:5432`
    -   **Kafka**: `localhost:9092`

---

## 🏗️ Архитектура и структура проекта

Проект состоит из нескольких ключевых компонентов, работающих вместе.

```
.
├── cmd/
│   └── main.go             # Главная точка входа в приложение
├── internal/
│   ├── cache/
│   │   └── cache.go        # Реализация LRU-кэша
│   ├── db/
│   │   └── postgres.go     # Подключение к PostgreSQL
│   ├── handlers/
│   │   └── handlers.go     # HTTP-обработчики (API)
│   ├── kafka/
│   │   └── kafka.go        # Подписчик на Kafka
│   ├── model/
│   │   └── model.go        # Структуры данных (модели)
│   ├── repository/
│   │   └── requests.go     # Логика работы с БД (SQL-запросы)
│   └── service/
│       └── service.go      # Бизнес-логика сервиса
├── web/
│   ├── index.html          # Веб-интерфейс
│   └── script.js           # JS для веб-интерфейса
├── publisher/
│   └── main.go             # Скрипт для отправки тестовых сообщений
├── docker-compose.yml      # Файл для оркестрации контейнеров
├── Dockerfile              # Файл для сборки Go-приложения
├── go.mod
├── go.sum
└── init.sql                # SQL-скрипт для инициализации БД
```

### Как это работает

1.  **Запуск**: `docker-compose` запускает четыре сервиса: `pgdb` (PostgreSQL), `zookeeper`, `kafka` и `go` (наш сервис).
2.  **Инициализация БД**: Контейнер `pgdb` при первом запуске выполняет скрипт `init.sql`, создавая необходимую структуру таблиц.
3.  **Старт Go-сервиса**:
    -   Приложение подключается к PostgreSQL.
    -   Загружает последние заказы из БД в LRU-кэш для быстрого доступа.
    -   Подключается к Kafka и подписывается на топик `orders`.
    -   Запускает HTTP-сервер на порту `8081`.
4.  **Получение нового заказа**:
    -   Сообщение с данными заказа публикуется в топик `orders` в Kafka.
    -   Go-сервис получает сообщение, парсит JSON.
    -   Сохраняет заказ в PostgreSQL через транзакцию.
    -   Если сохранение успешно, заказ добавляется в кэш.
5.  **Запрос данных через API**:
    -   Пользователь вводит ID заказа в веб-интерфейсе и нажимает "Найти".
  -   Frontend отправляет `GET` запрос на относительный путь `/order/{order_uid}` (базовый хост/порт берется из текущего окна).
    -   Сервис сначала ищет заказ в кэше. Если находит — мгновенно возвращает.
    -   Если в кэше нет — делает запрос в БД, сохраняет найденный результат в кэш и возвращает его.

---

## Важные фрагмены кода

### `docker-compose.yml`

```yaml
version: '3.8'

services:
  pgdb:
    # ...
  
  zookeeper:
    image: confluentinc/cp-zookeeper:7.0.1
    container_name: zookeeper
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000

  kafka:
    image: confluentinc/cp-kafka:7.0.1
    container_name: kafka
    depends_on:
      - zookeeper
    ports:
      - "9092:9092"
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: 'zookeeper:2181'
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:29092,PLAINTEXT_HOST://localhost:9092
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
      KAFKA_CREATE_TOPICS: "orders:1:1"

  go:
    build: .
    container_name: service-go
    ports:
      - "8081:8081"
    depends_on:
      - pgdb
      - kafka
    restart: on-failure
```

### `cmd/main.go`

Точка входа, где инициализируются и связываются все компоненты: база данных, кэш, сервис, Kafka и HTTP-роутер.

```go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"test-task/internal/cache"
	"test-task/internal/db"
	"test-task/internal/handlers"
	"test-task/internal/kafka"
	"test-task/internal/repository"
	"test-task/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	
	ctx, cancel := context.WithCancel(context.Background())
	kafkaSubscriber := kafka.NewKafkaSubscriber(orderService)
	go kafkaSubscriber.Subscribe(ctx)
	defer kafkaSubscriber.Close()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	cancel()
}
```

### `internal/kafka/kafka.go`
Kafka.

```go
package kafka

import (
	"context"
	"encoding/json"
	"log"
	"test-task/internal/model"
	"test-task/internal/service"

	"github.com/segmentio/kafka-go"
)

func (ks *KafkaSubscriber) Subscribe(ctx context.Context) {
	log.Println("Subscribed to Kafka topic:", topic)
	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping Kafka subscriber...")
			return
		default:
			m, err := ks.reader.FetchMessage(ctx)
			if err != nil {
				log.Printf("could not fetch message: %v", err)
				continue
			}

			var order model.Order
			if err := json.Unmarshal(m.Value, &order); err != nil {
				log.Printf("Failed to unmarshal order: %v", err)
				continue
			}

			ks.service.ProcessNewOrder(order)

			if err := ks.reader.CommitMessages(ctx, m); err != nil {
				log.Printf("failed to commit messages: %v", err)
			}
		}
	}
}
```

---

## Тесты

Чтобы проверить систему, нужно опубликовать тестовое сообщение в Kafka.

1.  **У вас есть файл `publisher/main.go`**:

    ```go
    package main

    import (
    	"context"
    	"encoding/json"
    	"io/ioutil"
    	"log"

    	"github.com/segmentio/kafka-go"
    )
    
    const (
        topic  = "orders"
        broker = "localhost:9092"
    )

    func main() {
    	data, err := ioutil.ReadFile("model.json")
    	if err != nil {
    		log.Fatalf("Failed to read model file: %v", err)
    	}

    	w := &kafka.Writer{
    		Addr:     kafka.TCP(broker),
    		Topic:    topic,
    		Balancer: &kafka.LeastBytes{},
    	}

    	err = w.WriteMessages(context.Background(), kafka.Message{Value: data})
    	if err != nil {
    		log.Fatalf("failed to write messages: %v", err)
    	}
        log.Println("Message published to Kafka successfully!")
    }
    ```

2.  **Добавьте `model.json`** в корень проекта с тестовыми данными.

3.  **Запустите паблишер** (настраивается через `PUB_BROKER`, `PUB_TOPIC`, `PUB_COUNT`):
    Выполните команду в терминале в корне проекта:
    ```bash
    go run publisher/main.go
    ```

4.  **Проверьте результат**:
    -   В логах `docker-compose` вы должны увидеть, что сервис получил и обработал заказ.
    -   Откройте [http://localhost:8081/](http://localhost:8081/), введите `order_uid` из вашего `model.json` (например, `b563feb7b2b84b6test`) и нажмите "Найти". Вы должны увидеть полную информацию о заказе.
