package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vyrodovalexey/metrics/internal/model"
	"github.com/vyrodovalexey/metrics/internal/storage"
	"net/http"
)

const (
	badrequest = "Bad Request"
)

func UpdateJSON(st storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("Content-Type") != "application/json" {
			c.JSON(http.StatusUnsupportedMediaType, gin.H{
				"error": badrequest,
			})
			return
		} else {
			var metrics model.Metrics
			err := json.NewDecoder(c.Request.Body).Decode(&metrics)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": badrequest,
				})
				return
			}
			
			switch metrics.MType {
			case "gauge":
				st.AddGauge(metrics.ID, *metrics.Value)
				g, e := st.GetGauge(metrics.ID)
				if !e {
					c.JSON(http.StatusNotFound, gin.H{
						"error": badrequest,
					})
					return
				} else {
					metrics.Value = &g
					c.JSON(http.StatusOK, metrics)
				}

			case "counter":
				st.AddCounter(metrics.ID, *metrics.Delta)
				g, e := st.GetCounter(metrics.ID)
				if !e {
					c.JSON(http.StatusNotFound, gin.H{
						"error": badrequest,
					})
					return
				} else {
					metrics.Delta = &g
					c.JSON(http.StatusOK, metrics)
				}
			default:
				c.JSON(http.StatusBadRequest, gin.H{
					"error": badrequest,
				})
				return
			}
		}
	}
}

func Update(st storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {

		switch c.Param("type") {
		case "gauge":
			err := st.AddGaugeAsString(c.Param("name"), c.Param("value"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": badrequest,
				})
			}

		case "counter":
			err := st.AddCounterAsString(c.Param("name"), c.Param("value"))
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

func GetJSON(st storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Request.Header.Get("Content-Type") != "application/json" {
			c.JSON(http.StatusUnsupportedMediaType, gin.H{
				"error": badrequest,
			})
			return
		} else {
			var metrics model.Metrics
			err := json.NewDecoder(c.Request.Body).Decode(&metrics)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": badrequest,
				})
				return
			}
			switch metrics.MType {
			case "gauge":
				g, e := st.GetGauge(metrics.ID)
				if !e {
					c.JSON(http.StatusNotFound, gin.H{
						"error": badrequest,
					})
					return
				} else {
					metrics.Value = &g
					c.JSON(http.StatusOK, metrics)
				}
			case "counter":
				g, e := st.GetCounter(metrics.ID)
				if !e {
					c.JSON(http.StatusNotFound, gin.H{
						"error": badrequest,
					})
					return
				} else {
					metrics.Delta = &g
					c.JSON(http.StatusOK, metrics)
				}
			default:
				c.JSON(http.StatusBadRequest, gin.H{
					"error": badrequest,
				})
				return
			}
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
