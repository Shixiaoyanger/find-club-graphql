package main

import (
	"find-club-graphql/controllers"
	"net/http"
)

func main() {

	http.HandleFunc("/graphql", controllers.Graphql)
	http.ListenAndServe(":8080", nil)

}
