package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	driver "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"
)

// Env contains useful data we want to share
type Env struct {
	dbmap *gorp.DbMap
}

var dbHost string
var dbUser string
var dbPassword string
var port string
var clientURL string

func init() {
	dbHost = os.Getenv("MYSQL_HOST")
	if dbHost == "" {
		dbHost = "localhost:3306"
	}
	dbUser = os.Getenv("MYSQL_USER")
	if dbUser == "" {
		dbUser = "root"
	}
	dbPassword = os.Getenv("MYSQL_PASSWORD")
	port = os.Getenv("PORT")
	if port == "" {
		port = ":5000"
	}
	clientURL = os.Getenv("CLIENT_URL")
	if clientURL == "" {
		clientURL = "http://localhost:3"
	}
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
	//Enable Cors
	r.Use(Cors())
	// Initialize routes
	v1 := r.Group("api/v1")
	// Box API Endpoints
	v1.GET("/boxes", env.GetBoxes)
	v1.POST("/boxes", env.CreateBox)
	// Instrument API Endpoints
	v1.GET("/instruments", env.GetInstruments)
	v1.POST("/instruments", env.CreateInstrument)
	v1.GET("/instruments/:id/boxes", env.GetInstrumentBoxes)
	// Specialty API Endpoints
	v1.GET("/specialties", env.GetSpecialties)
	v1.GET("/specialties/:specialty/boxes", env.GetBoxesBySpecialty)
	//BoxComposition API Endpoints
	v1.GET("/boxes/:id/content", env.GetBoxComposition)
	v1.POST("/boxes/:id/content", env.AddInstrumentToBox)
	//StatusReport API Endpoints
	v1.GET("/statusreports/reasons", env.GetReasonsStatusReport)
	v1.GET("/statusreports", env.GetStatusReports)
	v1.POST("/statusreports", env.CreateStatusReport)
	v1.PUT("/statusreports/:id", env.UpdateStatusReport)

	//IncidentReport API Endpoints
	v1.GET("/incidentreports/reasons", env.GetReasonsIncident)
	v1.GET("/incidentreports", env.GetIncidentReports)
	v1.POST("/incidentreports", env.CreateIncidentReport)

	// Run the router
	r.Run(port)
}

func initDb() *gorp.DbMap {
	config := &driver.Config{User: dbUser, Passwd: dbPassword, Net: "tcp", Addr: dbHost, DBName: "hospital"}
	db, err := sql.Open("mysql", config.FormatDSN())
	checkErr(err, "sql.Open failed")
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	dbmap.AddTableWithName(Box{}, "box").SetKeys(false, "RegistrationNumber")
	dbmap.AddTableWithName(Instrument{}, "instrument").SetKeys(false, "Ref")
	dbmap.AddTableWithName(BoxComposition{}, "box_composition").SetKeys(true, "ID").SetUniqueTogether("BoxID", "InstrumentID")
	dbmap.AddTableWithName(StatusReport{}, "status_report").SetKeys(true, "ID")
	dbmap.AddTableWithName(IncidentReport{}, "incident_report").SetKeys(true, "ID")
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create table failed")
	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

//Cors handle the CORS nonsense
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.JSON(200, "Cool cool cool")
			return
		}
		c.Next()
	}
}
