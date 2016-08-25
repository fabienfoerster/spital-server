package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Box represent a chirurgical box
type Box struct {
	ID           int64  `db:"id" json:"id"`
	Name         string `db:"name" json:"name"`
	LastModified int64  `db:"last_modified" json:"last_modified"`
	Information  string `db:"information" json:"information"`
	Specialty    string `db:"specialty" json:"specialty"`
}

// GetBoxes return all the boxes we have
func (env *Env) GetBoxes(c *gin.Context) {
	type Boxes []Box
	var boxes Boxes
	_, err := env.dbmap.Select(&boxes, "SELECT * FROM box")
	if err != nil {
		c.JSON(404, gin.H{"error": "no box(es) into the table"})
	} else {
		c.JSON(200, boxes)
	}

}

//CreateBox add a box in the listing of boxes
func (env *Env) CreateBox(c *gin.Context) {
	var box Box
	box.LastModified = time.Now().UnixNano()
	c.Bind(&box)
	if box.Name == "" || box.Specialty == "" {
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
