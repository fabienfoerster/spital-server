package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"
)

// Env contains useful data we want to share
type Env struct {
	dbmap *gorp.DbMap
}

// Instrument present a chirurgical instrument
type Instrument struct {
	ID   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

// Specialty describe the specialty in the hospital
type Specialty struct {
	ID   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

// BoxComposition represent the composition of a chirurgical box
type BoxComposition struct {
	BoxID        int64 `db:"boxid" json:"boxid"`
	InstrumentID int64 `db:"instrumentid" json:"instrumentid"`
}

var dbHost string
var dbUser string
var dbPassword string

func init() {
	dbHost = os.Getenv("MYSQL_HOST")
	dbUser = os.Getenv("MYSQL_USER")
	if dbUser == "" {
		dbUser = "root"
	}
	dbPassword = os.Getenv("MYSQL_PASSWORD")
}

func main() {
	// Initialize DB
	dbmap := initDb()
	defer dbmap.Db.Close()

	//Add the DB to the Env
	env := &Env{dbmap}

	// Routing
	// Initialize routing
	r := gin.Default()
	// Initialize routes
	v1 := r.Group("api/v1")
	v1.GET("/boxes", env.GetBoxes)
	v1.POST("/boxes", env.CreateBox)
	// Run the router
	r.Run(":5000")
}

func initDb() *gorp.DbMap {
	datasource := fmt.Sprintf("%s:%s@%s/hospital", dbUser, dbPassword, dbHost)
	db, err := sql.Open("mysql", datasource)
	checkErr(err, "sql.Open failed")
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	dbmap.AddTableWithName(Box{}, "box").SetKeys(true, "ID")
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create table failed")
	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
