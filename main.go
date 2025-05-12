package main

import (
	"encoding/json"
	"html/template"
	"log"
	"os"
	
	"github.com/russross/blackfriday/v2"
)

type Page struct {
	Title string
	Body  template.HTML
}

type Config struct {
	Title string `json:"title"`
}

type Site struct {
	Config
	Posts []Page
}

func main() {
	var err error
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	configData, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	err = json.Unmarshal(configData, &config)
	if err != nil {
		log.Fatal(err)
	}

	files, err := os.ReadDir("posts")
	if err != nil {
		log.Fatal(err)
	}

	posts := make([]Page, len(files))
	for i, file := range files {
		raw, err := os.ReadFile("posts/" + file.Name())
		if err != nil {
			log.Fatal(err)
		}
		body := template.HTML(blackfriday.Run(raw))
		posts[i] = Page{file.Name(), body}
	}

	site := Site{config, posts}

	err = os.RemoveAll("public")
	if err != nil {
		log.Fatal(err)
	}
	err = os.Mkdir("public", 0777)
	if err != nil {
		log.Fatal(err)
	}
	err = os.Mkdir("public/posts", 0777)
	if err != nil {
		log.Fatal(err)
	}

	index, err := os.Create("public/index.html") 
	if err != nil {
		log.Fatal(err)
	}
	defer index.Close()

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(index, site)

	tmpl, err = template.ParseFiles("templates/post.html")
	if err != nil {
		log.Fatal(err)
	}
	for _, post := range site.Posts {
		f, err := os.Create("public/posts/" + post.Title + ".html")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		tmpl.Execute(f, post)
	}
}
