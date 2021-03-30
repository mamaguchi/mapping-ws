package data

import (
    // "os"
    "net/http"
    "encoding/json"
    // "time"
    // "strconv"
    // "strings"
    // "errors"
    "fmt"
    // "log"
    "context"
    // "github.com/jackc/pgx"
    "github.com/jackc/pgx/pgxpool"
    // "github.com/jackc/pgconn"
    "pg_service/db"
    "pg_service/util"
    // "pg-service//auth"
)

type Geom struct {
    Point string    `json:"point"`
}

func AddNewFeature(conn *pgxpool.Pool, geom string) error {
	sql :=
		`insert into wbk.crud
		( geom )
		values
		( ST_PointFromText($1, 4326) )`
		
	_, err := conn.Exec(context.Background(), sql, geom)
	if err != nil {
		return err 
	}
	return nil 
}

func AddNewFeatureHandler(w http.ResponseWriter, r *http.Request) {
    util.SetDefaultHeader(w)
    if (r.Method == "OPTIONS") { return }
    fmt.Println("[AddNewFeatureHandler] request received")    

    // VERIFY AUTH TOKEN
    // authToken := strings.Split(r.Header.Get("Authorization"), " ")[1]
    // if !auth.VerifyTokenHMAC(authToken) {
    //     util.SendUnauthorizedStatus(w)
    //     return
    // }

    var geom Geom
    err := json.NewDecoder(r.Body).Decode(&geom)
    if err != nil {
        util.SendInternalServerErrorStatus(w, err)
        return
    }
    
    db.CheckDbConn()
    err = AddNewFeature(db.Conn, geom.Point)
    if err != nil {                        
        util.SendInternalServerErrorStatus(w, err)
        return 
    }    
}