package main

import (
	"net/http"
	"pg_service/db"
	"mapping/data"
)

func main() {
	/* INIT DATABASE CONNECTION */
	db.Open()
	defer db.Close()

	/* HANDLER FUNC */
	// Mapping
	http.HandleFunc("/point/add", data.AddNewFeatureHandler)
	

	/* START HTTP SERVER */
	http.ListenAndServe(":8080", nil)
}