package main

import (
	"fmt"
	"igin/igin"
	"net/http"
)

func main() {
	initEngine().Run(":8000")
}

func initEngine() *igin.Engine {
	e := igin.NewEngine()
	registerRouter(e)

	e.Use(func(c *igin.Context) {
		fmt.Println("middleware2 run ...")
		c.Next()
		fmt.Println("middleware2 end...")
	})

	return e
}

func registerRouter(e *igin.Engine) {
	e.GET("/a", func(c *igin.Context) {
		fmt.Println("aaaaa")
		c.GetResponseWriter().WriteHeader(http.StatusOK)
		c.String("success")
	})

	e.GET("/b", func(c *igin.Context) {
		result := map[string]interface{}{
			"hello": "b",
		}
		fmt.Println("bbbbb")
		c.JSON(result)
	})

	cGroup := e.Group("/c")
	cGroup.Use(func(c *igin.Context) {
		fmt.Println("middleware1 run ...")
		c.Next()
		fmt.Println("middleware1 end...")
	})
	cGroup.GET("/c1", func(c *igin.Context) {
		fmt.Println("c1c1c1")
	})

}
