package main

import (
	. "github.com/fishedee/app/log"
	. "github.com/fishedee/app/middleware"
	. "github.com/fishedee/app/router"
	. "github.com/fishedee/language"
	"net/http"
)

var (
	log    Log
	listen string = ":8299"
)

func run() {
	defer CatchCrash(func(e Exception) {
		log.Critical("server crash! %v", e.Error())
	})

	routerFactory := NewRouterFactory()
	routerFactory.Use(NewLogMiddleware(log, nil))
	routerFactory.Static("/", "../static")

	log.Debug("server is running... listen %v", listen)
	server := &http.Server{
		Addr:    listen,
		Handler: routerFactory.Create(),
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}
func main() {
	var err error
	log, err = NewLog(LogConfig{
		Driver: "console",
	})
	if err != nil {
		panic(err)
	}
	run()
}
