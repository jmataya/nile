package main

import (
	"log"
	"net/http"

	"github.com/jmataya/nile"
)

func main() {
	r := nile.New()
	r.GET("/hello", func(c *nile.Context) nile.Response {
		return basicResp{
			Message:    "hello",
			statusCode: http.StatusOK,
		}
	})

	r.GET("/world", func(c *nile.Context) nile.Response {
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
