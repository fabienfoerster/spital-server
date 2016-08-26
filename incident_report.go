package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

//IncidentReport represent an incident report
type IncidentReport struct {
	ID       int64  `db:"id" json:"id"`
	Reason   string `db:"reason" json:"reason"`
	Quantity int    `db:"quantity" json:"quantity"`
	Comment  string `db:"comment" json:"comment"`
	Date     int64  `db:"date" json:"date"`
}

var reasonsIncident = [...]string{"Pliage pasteur non conforme", "Température inadéquate", "Stérilité", "Sous-sachet troué", "Trace de brûlure", "Sachet pasteur troué"}

//GetReasonsIncident return the possible reasons for an incident report
func (env *Env) GetReasonsIncident(c *gin.Context) {
	c.JSON(200, reasonsIncident)
}

//GetIncidentReports return all the incident reports
func (env *Env) GetIncidentReports(c *gin.Context) {
	type IncidentReports []IncidentReport
	var incidentReports IncidentReports
	_, err := env.dbmap.Select(&incidentReports, "SELECT * FROM incident_report")
	if err != nil {
		log.Println(err)
		c.JSON(404, gin.H{"error": "no incident report(s) into the table"})
	} else {
		c.JSON(200, incidentReports)
	}
}

//CreateIncidentReport create a new incident report
func (env *Env) CreateIncidentReport(c *gin.Context) {
	var incidentReport IncidentReport
	c.Bind(&incidentReport)
	if incidentReport.Reason == "" || incidentReport.Quantity == 0 || incidentReport.Date == 0 {
		c.JSON(422, gin.H{"error": "fields are empty"})
	} else {
		err := env.dbmap.Insert(&incidentReport)
		if err != nil {
			log.Println(err)
			c.JSON(500, "Insert failed")
		} else {
			c.JSON(201, "Insert successful")
		}

	}
}

// curl -i -X POST -H "Content-Type: application/json" -d '{ "reason":"Stérilité","quantity":2,"date":14552225}' http://localhost:5000/api/v1/incidentreports
