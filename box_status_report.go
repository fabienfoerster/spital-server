package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

//StatusReport represent a status report
type StatusReport struct {
	ID            int64  `db:"id" json:"id"`
	BoxID         int64  `db:"boxid" json:"boxid"`
	InstrumentID  int64  `db:"instrumentid" json:"instrumentid"`
	Specialty     string `db:"specialty" json:"specialty"`
	Interlocutor  string `db:"interlocutor" json:"interlocutor"`
	DateGoing     int64  `db:"dategoing" json:"dategoing"`
	CommentGoing  string `db:"commentgoind" json:"commentgoing"`
	DateComing    int64  `db:"datecoming" json:"datecoming"`
	CommentComing string `db:"commentcoming" json:"commentcoming"`
	Reason        string `db:"reason" json:"reason"`
}

var reasons = [...]string{"Instru manquant", "Rajout", "Réparation", "Instru en plus", "Prob stérilité/propreté", "Modification boite", "Création boite", "Boite ouverte non utilisée pour réappro"}

//GetReasons return the possible reasons for a status report
func (env *Env) GetReasons(c *gin.Context) {
	c.JSON(200, reasons)
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
	if statusReport.BoxID == 0 || statusReport.InstrumentID == 0 || statusReport.Specialty == "" || statusReport.Interlocutor == "" || statusReport.DateGoing == 0 || statusReport.Reason == "" {
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

func (env *Env) handleReason(boxID int64, instruID int64, reason string) bool {
	if reason == "Rajout" {
		return env.modifyInstrumentCount(boxID, instruID, 1)
	}
	if reason == "Instru manquant" {
		return env.modifyInstrumentCount(boxID, instruID, -1)
	}
	return true
}
