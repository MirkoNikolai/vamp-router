package api

import (
	"net/http"

	"github.com/MirkoNikolai/vamp-router/haproxy"
	"github.com/gin-gonic/gin"
)

func GetConfig(c *gin.Context) {

	Config(c).BeginReadTrans()
	defer Config(c).EndReadTrans()

	c.JSON(http.StatusOK, Config(c))
}

func PostConfig(c *gin.Context) {

	Config(c).BeginWriteTrans()
	defer Config(c).EndWriteTrans()

	var config haproxy.Config

	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": err.Error()})
	} else {
		if err := Config(c).UpdateConfig(&config); err != nil {
			HandleError(c, err)
		} else {
			HandleReload(c, Config(c), http.StatusCreated, gin.H{"status": "updated config"})
		}
	}
}
