package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jmataya/nile"
)

func main() {
	r := nile.New()
	r.GET("/hello", func(c nile.Context) nile.Response {
		return basicResp{
			Message:    "hello",
			statusCode: http.StatusOK,
		}
	})

	r.GET("/world", func(c nile.Context) nile.Response {
		return basicResp{
			Message:    "world",
			statusCode: http.StatusOK,
		}
	})

	r.GET("/products/:id", func(c nile.Context) nile.Response {
		id := c.Param("id")
		message := fmt.Sprintf("Found product %s", id)

		return basicResp{
			Message:    message,
			statusCode: http.StatusOK,
		}
	})

	r.GET("/products/:id/edit", func(c nile.Context) nile.Response {
		id := c.Param("id")
		message := fmt.Sprintf("Editing product %s", id)

		return basicResp{
			Message:    message,
			statusCode: http.StatusOK,
		}
	})

	log.Fatal(r.Start(":8000"))
}

type basicResp struct {
	Message    string
	statusCode int
}

func (b basicResp) Body() interface{} {
	return map[string]string{"message": b.Message}
}

func (b basicResp) StatusCode() int {
	return b.statusCode
}
