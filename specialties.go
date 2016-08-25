package main

import "github.com/gin-gonic/gin"

// Specialty represent a specialty in the hospital
type Specialty struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

//GetSpecialties return all the specialties in the hospital
func (env *Env) GetSpecialties(c *gin.Context) {
	type Specialties []Specialty
	var specialties = Specialties{Specialty{1, "CEC"}, Specialty{2, "COELIO"}, Specialty{3, "DIGESTIF"}, Specialty{4, "ESTHETIQUE"}, Specialty{5, "GYNECO"}, Specialty{6, "OPHTALMO"}, Specialty{7, "ORL"}, Specialty{8, "ORTHO"}, Specialty{9, "PEDIATRIE"}, Specialty{10, "THORACIQUE"}, Specialty{11, "UROLOGIE"}}

	c.JSON(200, specialties)

}
