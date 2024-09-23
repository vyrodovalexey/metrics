package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vyrodovalexey/metrics/internal/storage"
	"net/http"
)

const (
	badrequest = "Bad Request"
)

func Update(st storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {

		switch c.Param("type") {
		case "gauge":
			err := st.AddGauge(c.Param("name"), c.Param("value"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": badrequest,
				})
			}

		case "counter":
			err := st.AddCounter(c.Param("name"), c.Param("value"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": badrequest,
				})
			}
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": badrequest,
			})
			return
		}
	}
}

func Get(st storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {

		switch c.Param("type") {
		case "gauge":
			g, e := st.GetGauge(c.Param("name"))
			if !e {
				c.JSON(http.StatusNotFound, gin.H{
					"error": badrequest,
				})
				return
			} else {
				gs := fmt.Sprintf("%v", g)
				c.String(http.StatusOK, gs)
			}
		case "counter":
			g, e := st.GetCounter(c.Param("name"))
			if !e {
				c.JSON(http.StatusNotFound, gin.H{
					"error": badrequest,
				})
				return
			} else {
				gs := fmt.Sprintf("%v", g)
				c.String(http.StatusOK, gs)
			}
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": badrequest,
			})
			return
		}
	}
}

func GetAllKeys(st storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {

		gval, cval := st.GetAllMetricNames()
		c.HTML(200, "table.tmpl", gin.H{
			"Title":         "Metric Names",
			"GaugeValues":   gval,
			"CounterValues": cval,
		})

	}
}
