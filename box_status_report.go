package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

//StatusReport represent a status report
type StatusReport struct {
	ID             int64  `db:"id" json:"id"`
	BoxID          string `db:"boxid" json:"boxid"`
	BoxName        string `db:"boxname" json:"boxname"`
	InstrumentID   string `db:"instrumentid" json:"instrumentid"`
	InstrumentName string `db:"instrumentname" json:"instrumentname"`
	Specialty      string `db:"specialty" json:"specialty"`
	Interlocutor   string `db:"interlocutor" json:"interlocutor"`
	DateGoing      string `db:"dategoing" json:"dategoing"`
	CommentGoing   string `db:"commentgoind" json:"commentgoing"`
	DateComing     string `db:"datecoming" json:"datecoming"`
	CommentComing  string `db:"commentcoming" json:"commentcoming"`
	Reason         string `db:"reason" json:"reason"`
}

var reasonsStatusReport = [...]string{"Instru manquant", "Rajout", "Réparation", "Instru en plus", "Prob stérilité/propreté", "Modification boite", "Création boite", "Boite ouverte non utilisée pour réappro"}

//GetReasonsStatusReport return the possible reasons for a status report
func (env *Env) GetReasonsStatusReport(c *gin.Context) {
	c.JSON(200, reasonsStatusReport)
}

//GetStatusReports return all the status report
func (env *Env) GetStatusReports(c *gin.Context) {
	type StatusReports []StatusReport
	var statusReports StatusReports
	_, err := env.dbmap.Select(&statusReports, "SELECT * FROM status_report")
	if err != nil {
		log.Println(err)
		c.JSON(404, gin.H{"error": "no status report(s) into the table"})
	} else {
		c.JSON(200, statusReports)
	}
}

//CreateStatusReport create a new status report
func (env *Env) CreateStatusReport(c *gin.Context) {
	var statusReport StatusReport
	c.Bind(&statusReport)
	fmt.Println(statusReport)
	if statusReport.BoxID == "" || statusReport.InstrumentID == "" || statusReport.Specialty == "" || statusReport.Interlocutor == "" || statusReport.DateGoing == "" || statusReport.Reason == "" {
		c.JSON(422, gin.H{"error": "fields are empty"})
	} else {
		err := env.dbmap.Insert(&statusReport)
		if err != nil {
			log.Println(err)
			c.JSON(500, "Insert failed")
		} else {
			ok := env.handleReason(statusReport.BoxID, statusReport.InstrumentID, statusReport.Reason)
			if !ok {
				log.Println("Modify count of box composition related to status report failed ...")
			}
			c.JSON(201, "Insert successful")
		}
	}
}

// curl -i -X POST -H "Content-Type: application/json" -d '{"boxid":1,"instrumentid":2,"specialty":"CEC","interlocutor":"Pauline","dategoing":145223364,"reason":"Rajout"}' http://localhost:5000/api/v1/statusreports

//UpdateStatusReport update a status report
func (env *Env) UpdateStatusReport(c *gin.Context) {
	id := c.Params.ByName("id")
	var statusReport StatusReport
	err := env.dbmap.SelectOne(&statusReport, "SELECT * FROM status_report WHERE id=?", id)
	if err != nil {
		log.Println(err)
		c.JSON(404, gin.H{"error": "status report not found"})
	}
	var json StatusReport
	c.Bind(&json)
	if json.BoxID == "" || json.InstrumentID == "" || json.Specialty == "" || json.Interlocutor == "" || json.DateGoing == "" || json.Reason == "" {
		c.JSON(422, gin.H{"error": "fields are empty"})
	}
	json.ID = statusReport.ID
	_, err = env.dbmap.Update(&json)
	if err != nil {
		log.Println(err)
		c.JSON(500, "Update failed")
	}
	c.JSON(200, json)
}

func (env *Env) handleReason(boxID string, instruID string, reason string) bool {
	if reason == "Rajout" {
		return env.modifyInstrumentCount(boxID, instruID, 1)
	}
	if reason == "Instru manquant" {
		return env.modifyInstrumentCount(boxID, instruID, -1)
	}
	return true
}
