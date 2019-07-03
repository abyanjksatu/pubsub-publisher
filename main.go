package main

import (
	"context"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/labstack/echo"
)

const topicName = "kecci"

func main() {
	e := echo.New()
	e.GET("/", Send)
	e.GET("/many", SendMany)
	e.Logger.Fatal(e.Start(":1234"))
}

// Send generates random integer and sends it to Cloud Pub/Sub
func Send(c echo.Context) error {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, "docker-220612")
	if err != nil {
		c.Response().Write([]byte(err.Error()))
		return c.JSON(http.StatusInternalServerError, c.Response())
	}

	topic := client.Topic(topicName)

	rand.Seed(time.Now().UnixNano())

	result := topic.Publish(ctx, &pubsub.Message{
		Data: []byte(strconv.Itoa(rand.Intn(1000))),
	})
	id, err := result.Get(ctx)
	if err != nil {
		c.Response().Write([]byte(err.Error()))
		return c.JSON(http.StatusInternalServerError, c.Response())
	}

	c.Response().WriteHeader(http.StatusCreated)
	c.Response().Write([]byte(id))
	return c.JSON(http.StatusCreated, c.Response())
}

func SendMany(c echo.Context) error {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, "docker-220612")
	if err != nil {
		c.Response().Write([]byte(err.Error()))
		return c.JSON(http.StatusInternalServerError, c.Response())
	}

	topic := client.Topic(topicName)
	var a [100]string

	for i := 0; i < 100; i++ {
		rand.Seed(time.Now().UnixNano())

		result := topic.Publish(ctx, &pubsub.Message{
			Data: []byte(strconv.Itoa(rand.Intn(1000))),
		})

		id, err := result.Get(ctx)
		if err != nil {
			id = ""
		}
		a[i] = id
	}

	if err != nil {
		c.Response().Write([]byte(err.Error()))
		return c.JSON(http.StatusInternalServerError, c.Response())
	}

	// c.Response().WriteHeader(http.StatusCreated)
	// c.Response().Write([]byte(a))
	return c.JSON(http.StatusCreated, a)
}
