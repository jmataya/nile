package main

import (
	"log"
	"net/http"

	"github.com/jmataya/nile"
	"github.com/jmataya/nile/routing"
)

func main() {
	r := nile.New()
	err := r.GET("/hello", func(c *routing.Context) routing.Response {
		return basicResp{
			Message:    "hello",
			statusCode: http.StatusOK,
		}
	})

	if err != nil {
		log.Fatal(err)
	}

	r.GET("/world", func(c *routing.Context) routing.Response {
		return basicResp{
			Message:    "world",
			statusCode: http.StatusOK,
		}
	})

	log.Fatal(r.Start(":8000"))
}

type basicResp struct {
	Message    string `json:"message"`
	statusCode int
}

func (b basicResp) StatusCode() int {
	return b.statusCode
}

func (b basicResp) Status() string {
	return http.StatusText(b.statusCode)
}
