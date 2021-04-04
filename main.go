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
	http.HandleFunc("/geom/add", data.AddNewGeomHandler)
	http.HandleFunc("/geoms/get", data.GetGeomsHandler)
	http.HandleFunc("/geom/update", data.UpdateGeomHandler)
	http.HandleFunc("/geom/del", data.DelGeomHandler)
	http.HandleFunc("/address/point/get", data.GetAddressPointHandler)
	

	/* START HTTP SERVER */
	http.ListenAndServe(":8080", nil)
}