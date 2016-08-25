package main

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

// BoxComposition represent the composition of a chirurgical box
type BoxComposition struct {
	BoxID        int64 `db:"boxid" json:"boxid"`
	InstrumentID int64 `db:"instrumentid" json:"instrumentid"`
	Quantity     int   `db:"quantity" json:"quantity"`
	Missing      int   `db:"missing" json:"missing"`
}

type BoxContent struct {
	Name     string `json:"name"`
	Quantity int    ` json:"quantity"`
	Missing  int    ` json:"missing"`
}

//GetBoxComposition return the instruments inside a box
func (env *Env) GetBoxComposition(c *gin.Context) {
	//id := c.Params.ByName("id")
	type WholeContent []BoxContent
	var wholeContent WholeContent
	_, err := env.dbmap.Select(&wholeContent, "SELECT instrument.name, box_composition.quantity, box_composition.missing FROM box_composition INNER JOIN instrument ON box_composition.instrumentid=instrument.id")
	if err != nil {
		log.Println(err)
		c.JSON(404, gin.H{"error": "no instrument(s) found for the box"})
	} else {
		c.JSON(200, wholeContent)
	}
}

//AddInstrumentToBox add an instrument to the box
func (env *Env) AddInstrumentToBox(c *gin.Context) {
	id := c.Params.ByName("id")
	var boxComposition BoxComposition
	c.Bind(&boxComposition)
	if boxComposition.InstrumentID == 0 || boxComposition.Quantity == 0 {
		c.JSON(422, gin.H{"error": "fields are empty"})
	} else {
		id64, _ := strconv.ParseInt(id, 10, 64)
		boxComposition.BoxID = id64
		err := env.dbmap.Insert(&boxComposition)
		if err != nil {
			log.Println(err)
			c.JSON(500, "Insert failed")
		} else {
			c.JSON(201, "Insert successful")
		}
	}
}

// curl -i -X POST -H "Content-Type: application/json" -d '{"instrumentid":2,"quantity":5}' http://localhost:5000/api/v1/boxes/2/content
