package api

import (
	"github.com/gin-gonic/gin"
	"github.com/MirkoNikolai/vamp-router/haproxy"
	"net/http"
)

func GetRoutes(c *gin.Context) {

	Config(c).BeginReadTrans()
	defer Config(c).EndReadTrans()

	result := Config(c).GetRoutes()
	if Config(c).GetRoutes() != nil {
		c.JSON(200, result)
	} else {
		c.String(http.StatusNotFound, "no routes found")
	}

}

func GetRoute(c *gin.Context) {

	Config(c).BeginReadTrans()
	defer Config(c).EndReadTrans()

	routeName := c.Params.ByName("route")

	if result, err := Config(c).GetRoute(routeName); err != nil {
		HandleError(c, err)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func PutRoute(c *gin.Context) {

	Config(c).BeginWriteTrans()
	defer Config(c).EndWriteTrans()

	var route haproxy.Route
	routeName := c.Params.ByName("route")

	if err := c.ShouldBindJSON(&route); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if err := Config(c).UpdateRoute(routeName, &route); err != nil {
			HandleError(c, err)
		} else {
			HandleReload(c, Config(c), http.StatusOK, gin.H{"status": "updated route"})
		}
	}
}

func PostRoute(c *gin.Context) {

	Config(c).BeginWriteTrans()
	defer Config(c).EndWriteTrans()

	var route haproxy.Route

	// This will infer what binder to use depending on the content-type header.
	if err := c.ShouldBindJSON(&route); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if err := Config(c).AddRoute(route); err != nil {
			HandleError(c, err)
		} else {
			HandleReload(c, Config(c), http.StatusCreated, gin.H{"status": "created route"})
		}
	}
}

func DeleteRoute(c *gin.Context) {

	Config(c).BeginReadTrans()
	defer Config(c).EndReadTrans()

	routeName := c.Params.ByName("route")

	if err := Config(c).DeleteRoute(routeName); err != nil {
		HandleError(c, err)
	} else {
		HandleReload(c, Config(c), http.StatusNoContent, gin.H{})
	}
}

func GetRouteServices(c *gin.Context) {

	Config(c).BeginReadTrans()
	defer Config(c).EndReadTrans()

	routeName := c.Params.ByName("route")

	result, err := Config(c).GetRouteServices(routeName)
	if err != nil {
		HandleError(c, err)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func GetRouteService(c *gin.Context) {

	Config(c).BeginReadTrans()
	defer Config(c).EndReadTrans()

	routeName := c.Params.ByName("route")
	serviceName := c.Params.ByName("service")

	result, err := Config(c).GetRouteService(routeName, serviceName)
	if err != nil {
		HandleError(c, err)
	} else {
		c.JSON(http.StatusOK, result)
	}

}

func PutRouteService(c *gin.Context) {

	Config(c).BeginWriteTrans()
	defer Config(c).EndWriteTrans()

	var service haproxy.Service
	routeName := c.Params.ByName("route")
	serviceName := c.Params.ByName("service")

	if err := c.ShouldBindJSON(&service); err != nil {
		c.String(500, "Invalid JSON")
	} else {
		if err := Config(c).UpdateRouteService(routeName, serviceName, &service); err != nil {
			HandleError(c, err)
		} else {
			HandleReload(c, Config(c), 200, gin.H{"status": "updated service"})
		}		
	}
}

func PutRouteServices(c *gin.Context) {

	Config(c).BeginWriteTrans()
	defer Config(c).EndWriteTrans()

	var services []*haproxy.Service
	routeName := c.Params.ByName("route")

	if err := c.ShouldBindJSON(&services) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if err := Config(c).UpdateRouteServices(routeName, services); err != nil {
			HandleError(c, err)
		} else {
			HandleReload(c, Config(c), http.StatusOK, gin.H{"status": "updated services"})
		}
	}
}

func PostRouteService(c *gin.Context) {

	Config(c).BeginWriteTrans()
	defer Config(c).EndWriteTrans()

	var services []*haproxy.Service
	routeName := c.Params.ByName("route")

	if err:= c.ShouldBindJSON(&services); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": err.Error()}) 
	} else {
		if err := Config(c).AddRouteServices(routeName, services); err != nil {
			HandleError(c, err)
		} else {
			HandleReload(c, Config(c), http.StatusCreated, gin.H{"status": "created service(s)"})
		}

	}
}

func DeleteRouteService(c *gin.Context) {

	Config(c).BeginWriteTrans()
	defer Config(c).EndWriteTrans()

	routeName := c.Params.ByName("route")
	serviceName := c.Params.ByName("service")

	if err := Config(c).DeleteRouteService(routeName, serviceName); err != nil {
		HandleError(c, err)
	} else {
		HandleReload(c, Config(c), http.StatusNoContent, gin.H{})
	}
}

func GetServiceServers(c *gin.Context) {

	Config(c).BeginReadTrans()
	defer Config(c).EndReadTrans()

	routeName := c.Params.ByName("route")
	serviceName := c.Params.ByName("service")

	result, err := Config(c).GetServiceServers(routeName, serviceName)
	if err != nil {
		HandleError(c, err)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func GetServiceServer(c *gin.Context) {

	Config(c).BeginReadTrans()
	defer Config(c).EndReadTrans()

	routeName := c.Params.ByName("route")
	serviceName := c.Params.ByName("service")
	serverName := c.Params.ByName("server")

	result, err := Config(c).GetServiceServer(routeName, serviceName, serverName)
	if err != nil {
		HandleError(c, err)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func DeleteServiceServer(c *gin.Context) {

	Config(c).BeginWriteTrans()
	defer Config(c).EndWriteTrans()

	routeName := c.Params.ByName("route")
	serviceName := c.Params.ByName("service")
	serverName := c.Params.ByName("server")

	if err := Config(c).DeleteServiceServer(routeName, serviceName, serverName); err != nil {
		HandleError(c, err)
	} else {
		HandleReload(c, Config(c), http.StatusNoContent, gin.H{})
	}
}

func PostServiceServer(c *gin.Context) {

	Config(c).BeginWriteTrans()
	defer Config(c).EndWriteTrans()

	var server haproxy.Server
	routeName := c.Params.ByName("route")
	serviceName := c.Params.ByName("service")

	if c.Bind(&server) != nil {
		if err := Config(c).AddServiceServer(routeName, serviceName, &server); err != nil {
			HandleError(c, err)
		} else {
			HandleReload(c, Config(c), http.StatusCreated, gin.H{"status": "created server"})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "bad request"})
	}
}

func PutServiceServer(c *gin.Context) {

	Config(c).BeginWriteTrans()
	defer Config(c).EndWriteTrans()

	var server haproxy.Server
	routeName := c.Params.ByName("route")
	serviceName := c.Params.ByName("service")
	serverName := c.Params.ByName("server")

	if c.Bind(&server) != nil {
		if err := Config(c).UpdateServiceServer(routeName, serviceName, serverName, &server); err != nil {
			HandleError(c, err)
		} else {
			HandleReload(c, Config(c), http.StatusOK, gin.H{"status": "updated server"})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "bad request"})
	}
}
