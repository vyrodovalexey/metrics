package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vyrodovalexey/metrics/internal/storage"
	"net/http"
)

func Update(st *storage.MemStorage) gin.HandlerFunc {
	return func(c *gin.Context) {

		switch c.Param("type") {
		case "gauge":
			err := st.AddGauge(c.Param("name"), c.Param("value"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Bad Request",
				})
			}
		case "counter":
			err := st.AddCounter(c.Param("name"), c.Param("value"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Bad Request",
				})
			}
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Bad Request",
			})
			return
		}
	}
}

func Get(st *storage.MemStorage) gin.HandlerFunc {
	return func(c *gin.Context) {

		switch c.Param("type") {
		case "gauge":
			g, e := st.GetGauge(c.Param("name"))
			if !e {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Bad Request",
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
					"error": "Bad Request",
				})
				return
			} else {
				gs := fmt.Sprintf("%v", g)
				c.String(http.StatusOK, gs)
			}
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Bad Request",
			})
			return
		}
	}
}

func GetAllKeys(st *storage.MemStorage) gin.HandlerFunc {
	return func(c *gin.Context) {

		gval, cval := st.GetAllMetricNames()
		c.HTML(200, "table.tmpl", gin.H{
			"Title":         "Metric Names",
			"GaugeValues":   gval,
			"CounterValues": cval,
		})

	}
}
