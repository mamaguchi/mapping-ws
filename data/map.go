package data

import (
    // "os"
    "net/http"
    "encoding/json"
    // "time"
    // "strconv"
    // "strings"
    "errors"
    "fmt"
    // "log"
    "context"
    "github.com/jackc/pgx"
    "github.com/jackc/pgx/pgxpool"
    // "github.com/jackc/pgconn"
    "pg_service/db"
    "pg_service/util"
    // "pg-service//auth"
)

type NewGeomIn struct {
    PointWKT string     `json:"pointWKT"`
}

type Geom struct {
    Id int              `json:"id"`
    PointWKT string     `json:"pointWKT"`
}

type GeomsOut struct {
    Geoms []Geom        `json:"geoms"`
}

type DelGeomIn struct {
    Id int              `json:"id"`
}


func AddNewGeom(conn *pgxpool.Pool, geom string) error {
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

func AddNewGeomHandler(w http.ResponseWriter, r *http.Request) {
    util.SetDefaultHeader(w)
    if (r.Method == "OPTIONS") { return }
    fmt.Println("[AddNewGeomHandler] request received")    

    // VERIFY AUTH TOKEN
    // authToken := strings.Split(r.Header.Get("Authorization"), " ")[1]
    // if !auth.VerifyTokenHMAC(authToken) {
    //     util.SendUnauthorizedStatus(w)
    //     return
    // }

    var geom NewGeomIn
    err := json.NewDecoder(r.Body).Decode(&geom)
    if err != nil {
        util.SendInternalServerErrorStatus(w, err)
        return
    }
    
    db.CheckDbConn()
    err = AddNewGeom(db.Conn, geom.PointWKT)
    if err != nil {                        
        util.SendInternalServerErrorStatus(w, err)
        return 
    }    
}

func GetGeoms(conn *pgxpool.Pool) ([]byte, error) {
    sql := 
        `select id, ST_AsText(geom) as wkt
           from wbk.crud`

    rows, err := conn.Query(context.Background(), sql)
    if err != nil {
        return nil, err 
    }

    var geoms GeomsOut
    for rows.Next() {
        var id int 
        var pointWKT string 

        err = rows.Scan(&id, &pointWKT)
        if err != nil {
            return nil, err
        }

        geom := Geom{
            Id: id,
            PointWKT: pointWKT,
        }
        geoms.Geoms = append(geoms.Geoms, geom)
    }
    if len(geoms.Geoms) == 0 {
        return nil, pgx.ErrNoRows
    }
    outputJson, err := json.Marshal(geoms)
    return outputJson, err
}

func GetGeomsHandler(w http.ResponseWriter, r *http.Request) {
    util.SetDefaultHeader(w)
    if (r.Method == "OPTIONS") { return }
    fmt.Println("[GetGeomsHandler] request received")    

    // VERIFY AUTH TOKEN
    // authToken := strings.Split(r.Header.Get("Authorization"), " ")[1]
    // if !auth.VerifyTokenHMAC(authToken) {
    //     util.SendUnauthorizedStatus(w)
    //     return
    // }   

    // var wbkcase Wbkcase
    // err := json.NewDecoder(r.Body).Decode(&wbkcase)
    // if err != nil {
    //     util.SendInternalServerErrorStatus(w, err)
    //     return
    // }
    
    db.CheckDbConn()
    closeContactsJson, err := GetGeoms(db.Conn)
    if err != nil {      
        if err == pgx.ErrNoRows { 		
            util.SendStatusNotFound(w, err)
			return 
		}   
        util.SendInternalServerErrorStatus(w, err)
        return 
    }
    fmt.Printf("%s\n", closeContactsJson)
    fmt.Fprintf(w, "%s", closeContactsJson)
}

func UpdateGeom(conn *pgxpool.Pool, g Geom) error {
    sql := 
        `update wbk.crud
           set geom = ST_PointFromText($1, 4326)
           where id = $2
             returning id`

    // _, err := conn.Exec(context.Background(), sql, 
    //     g.PointWKT, g.Id)

    row := conn.QueryRow(context.Background(), sql, 
        g.PointWKT, g.Id)
    var id int
    err := row.Scan(&id)        

    if err != nil {
        return err
    }
    return nil
}

func UpdateGeomHandler(w http.ResponseWriter, r *http.Request) {
    util.SetDefaultHeader(w)
    if (r.Method == "OPTIONS") { return }
    fmt.Println("[UpdateGeomHandler] request received")    

    // VERIFY AUTH TOKEN
    // authToken := strings.Split(r.Header.Get("Authorization"), " ")[1]
    // if !auth.VerifyTokenHMAC(authToken) {
    //     util.SendUnauthorizedStatus(w)
    //     return
    // }

    var g Geom
    err := json.NewDecoder(r.Body).Decode(&g)
    if err != nil {
        util.SendInternalServerErrorStatus(w, err)
        return
    }
    
    db.CheckDbConn()
    err = UpdateGeom(db.Conn, g)
    if err != nil {     
        util.SendInternalServerErrorStatus(w, err)
        return 
    }    
}

func DelGeom(conn *pgxpool.Pool, id int) error {
    sql := 
        `with deleted as (delete from wbk.crud
           where id = $1 returning id)
         select count(id) from deleted`

    // _, err := conn.Exec(context.Background(), sql, 
    //     id)

    row := conn.QueryRow(context.Background(), sql, 
        id)
    var deletedCount int
    err := row.Scan(&deletedCount)        

    if err != nil {
        return err
    }
    if deletedCount == 0 {
        return errors.New("Error: unable to delete geom")
    }
    return nil
}

func DelGeomHandler(w http.ResponseWriter, r *http.Request) {
    util.SetDefaultHeader(w)
    if (r.Method == "OPTIONS") { return }
    fmt.Println("[DelGeomHandler] request received")    

    // VERIFY AUTH TOKEN
    // authToken := strings.Split(r.Header.Get("Authorization"), " ")[1]
    // if !auth.VerifyTokenHMAC(authToken) {
    //     util.SendUnauthorizedStatus(w)
    //     return
    // }

    var dgi DelGeomIn
    err := json.NewDecoder(r.Body).Decode(&dgi)
    if err != nil {
        util.SendInternalServerErrorStatus(w, err)
        return
    }
    
    db.CheckDbConn()
    err = DelGeom(db.Conn, dgi.Id)
    if err != nil {     
        util.SendInternalServerErrorStatus(w, err)
        return 
    }    
}