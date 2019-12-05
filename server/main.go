package main

import (
	. "github.com/fishedee/app/log"
	. "github.com/fishedee/app/middleware"
	. "github.com/fishedee/app/router"
	. "github.com/fishedee/language"
	"html/template"
	"io/ioutil"
	"net/http"
)

var (
	log              Log
	listen           string = ":8299"
	categoryTemplate *template.Template
)

func run() {
	defer CatchCrash(func(e Exception) {
		log.Critical("server crash! %v", e.Error())
	})

	//预先读取模板
	var err error
	categoryTemplate, err = template.ParseFiles("../static/list.html")
	if err != nil {
		panic(err)
	}

	//启动服务器
	routerFactory := NewRouterFactory()
	routerFactory.Use(NewLogMiddleware(log, nil))
	routerFactory.Static("/", "../static")
	routerFactory.GET("/:category", getCategory)

	log.Debug("server is running... listen %v", listen)
	server := &http.Server{
		Addr:    listen,
		Handler: routerFactory.Create(),
	}
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

type Image struct {
	Src string
}

type CategoryTemplateData struct {
	Title  string
	Images []Image
}

type CategoryInfo struct {
	Title string
	Src   string
}

func getCategory(w http.ResponseWriter, r *http.Request, param RouterParam) {
	category := param[0].Value

	categroyMap := map[string]CategoryInfo{
		"nature.html":   CategoryInfo{"自然", "/gallery/nature"},
		"food.html":     CategoryInfo{"美食", "/gallery/food"},
		"face.html":     CategoryInfo{"人物", "/gallery/face"},
		"travel.html":   CategoryInfo{"旅游", "/gallery/travel"},
		"wildlife.html": CategoryInfo{"动物", "/gallery/wildlife"},
		"city.html":     CategoryInfo{"城市", "/gallery/city"},
	}
	categoryInfo, isExist := categroyMap[category]
	if isExist == false {
		panic("not exist category " + category)
	}

	fileInfo, err := ioutil.ReadDir("../static" + categoryInfo.Src)
	if err != nil {
		panic(err)
	}

	result := CategoryTemplateData{}
	result.Title = categoryInfo.Title
	result.Images = []Image{}

	for _, file := range fileInfo {
		result.Images = append(result.Images, Image{
			Src: categoryInfo.Src + "/" + file.Name(),
		})
	}
	categoryTemplate.Execute(w, result)
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
