package main

import (
"encoding/json"
"fmt"
"log"
"net/http"
"strings"
"time"
)
type Article struct {
	ID                string    `json:"id"`
	Title             string    `json:"title"`
	Subtitle          string    `json:"subtitle"`
	Content           string    `json:"content"`
  CreationTime      time.Time `json:"time"`
}

var Articles []Article


//Home page which says hello
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")

}


//Return All Articles in json format
func return_all_articles(w http.ResponseWriter, r *http.Request) {

	
	json.NewEncoder(w).Encode(Articles)
}


//Returns 1 single Article for the provided ID
func return_one_article(w http.ResponseWriter, r *http.Request) {

  var Urls =  strings.Split(r.URL.Path, "/")
  var Temp_Article = Articles[0]
  Temp_Article.ID = "NOTFOUND"
  for i := 0; i < len(Articles); i++{
      if Urls[2]==Articles[i].ID {
        Temp_Article = Articles[i]
      }

    }
  if Temp_Article.ID == "NOTFOUND" {

    w.WriteHeader(http.StatusNotFound)
    return
  }

  if len(Urls) != 3 {

    w.WriteHeader(http.StatusNotFound)
    return
  }
  json.NewEncoder(w).Encode(Temp_Article)

}

//Creates Article Through json provided in POST request

func create(w http.ResponseWriter, r *http.Request) {
  var temp_article Article

    // Try to decode the request body into the struct. If there is an error,
    // respond to the client with the error message and a 400 status code.
    err := json.NewDecoder(r.Body).Decode(&temp_article)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    temp_article.CreationTime = time.Now()
    // Do something with the Person struct...
    Articles = append(Articles,temp_article)
    json.NewEncoder(w).Encode(temp_article)
}

//Searches for articles with keywords provided
func search(w http.ResponseWriter, r *http.Request) {

    var Found_Articles  []Article
    var q1 = r.URL.Query().Get("q")

    for i := 0; i < len(Articles); i++{
        contains_title := strings.Contains(Articles[i].Title,q1)
        contains_subtitle := strings.Contains(Articles[i].Subtitle,q1)
        contains_content := strings.Contains(Articles[i].Content,q1)
        present := contains_title || contains_subtitle || contains_content
        if present {
          Found_Articles = append(Found_Articles,Articles[i])
        }

      }
  json.NewEncoder(w).Encode(Found_Articles)
}

//decides whether to create an article or return all articles based on request type
func create_or_return(w http.ResponseWriter, r *http.Request) {
  if r.Method == "POST" {
    create(w,r)
    return
  }

  if r.Method == "GET" {
    return_all_articles(w,r)
    return
  }


}

//handles all requests
func handleRequests() {
    http.HandleFunc("/", homePage)
    http.HandleFunc("/articles",create_or_return)
    http.HandleFunc("/articles/",return_one_article)
    http.HandleFunc("/articles/search",search)
    log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
  Articles = []Article{
  Article{ID: "1", Title: "Hi", Subtitle: "Article", Content: "Content"},
  Article{ID: "2", Title: "Hi Again", Subtitle: "Article 2", Content: "Content2"},
}

    handleRequests()

}
