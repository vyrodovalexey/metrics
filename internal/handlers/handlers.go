package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vyrodovalexey/metrics/internal/model"
	"github.com/vyrodovalexey/metrics/internal/storage"
	"net/http"
	"os"
)

const (
	badrequest = "Bad Request"
)

func UpdateFromBodyJSON(st storage.Storage, f *os.File, p bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("Content-Type") != "application/json" {
			c.JSON(http.StatusUnsupportedMediaType, gin.H{
				"error": badrequest,
			})
			return
		} else {
			m := &model.Metrics{}
			body := c.Request.Body
			err := m.BodyToMetric(&body)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": badrequest,
				})
				return
			}
			err = st.UpdateMetric(m, f, p)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
			st.GetMetric(m)
			c.JSON(http.StatusOK, m)
			return
		}
	}
}

func UpdateFromURLPath(st storage.Storage, f *os.File, p bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		m := &model.Metrics{}
		err := m.URLPathToMetric(c.Param("type"), c.Param("name"), c.Param("value"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		err = st.UpdateMetric(m, f, p)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		st.GetMetric(m)
		if m.MType == "gauge" {
			c.String(http.StatusOK, fmt.Sprintf("%v", *m.Value))
		} else {
			c.String(http.StatusOK, fmt.Sprintf("%d", *m.Delta))
		}
		return

	}
}

func Get(st storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		m := &model.Metrics{}
		err := m.URLPathToMetric(c.Param("type"), c.Param("name"), c.Param("value"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		b := st.GetMetric(m)
		if !b {
			c.JSON(http.StatusNotFound, gin.H{
				"error": badrequest,
			})
			return
		}
		if m.MType == "gauge" {
			c.String(http.StatusOK, fmt.Sprintf("%v", *m.Value))
		} else {
			c.String(http.StatusOK, fmt.Sprintf("%d", *m.Delta))
		}
		return
	}
}

func GetBodyJSON(st storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Request.Header.Get("Content-Type") != "application/json" {
			c.JSON(http.StatusUnsupportedMediaType, gin.H{
				"error": badrequest,
			})
			return
		} else {
			m := &model.Metrics{}
			body := c.Request.Body
			err := m.BodyToMetric(&body)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": badrequest,
				})
				return
			}
			b := st.GetMetric(m)
			if !b {
				c.JSON(http.StatusNotFound, gin.H{
					"error": badrequest,
				})
				return
			}
			c.JSON(http.StatusOK, m)
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
