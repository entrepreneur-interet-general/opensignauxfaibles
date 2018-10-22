package main

import (
	"github.com/gin-gonic/gin"
)

func debug(c *gin.Context) {
	naf, err := loadNAF()
	if err != nil {
		c.JSON(500, err)
	}
	c.JSON(200, naf)
}
