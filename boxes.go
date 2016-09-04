package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Box represent a chirurgical box
type Box struct {
	RegistrationNumber string `db:"registration_number" json:"registration_number"`
	Name               string `db:"name" json:"name"`
	LastModified       int64  `db:"last_modified" json:"last_modified"`
	Information        string `db:"information" json:"information"`
	Specialty          string `db:"specialty" json:"specialty"`
}

// GetBoxes return all the boxes we have
func (env *Env) GetBoxes(c *gin.Context) {
	type Boxes []Box
	var boxes Boxes
	_, err := env.dbmap.Select(&boxes, "SELECT * FROM box")
	if err != nil {
		log.Println(err)
		c.JSON(404, gin.H{"error": "no box(es) into the table"})
	} else {
		c.JSON(200, boxes)
	}
}

//CreateBox add a box in the listing of boxes
func (env *Env) CreateBox(c *gin.Context) {
	var box Box
	box.LastModified = time.Now().Unix()
	c.Bind(&box)
	if box.Name == "" || box.Specialty == "" || box.RegistrationNumber == "" {
		c.JSON(422, gin.H{"error": "fields are empty"})
	} else {
		err := env.dbmap.Insert(&box)
		if err != nil {
			log.Println(err)
			c.JSON(500, "Insert failed")
		} else {
			c.JSON(201, "Insert successful")
		}
	}
}

// curl -i -X POST -H "Content-Type: application/json" -d '{ "name" : "Fabien", "specialty": "CEC", "registration_number":"MAF1231HJO"}' http://localhost:5000/api/v1/boxes
