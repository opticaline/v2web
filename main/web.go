package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/opticaline/v2ray-web"
	"log"
	"net/http"
	"time"
)

func Database(db *core.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// Set example variable
		c.Set("db", db)

		// before request

		c.Next()

		// after request
		latency := time.Since(t)
		log.Print(latency)

		// access the status we are sending
		status := c.Writer.Status()
		log.Println(status)
	}
}

type QueryAll struct {
	Start time.Time `form:"start" time_format:"2006-01-02T15:04:05" time_utc:"1"`
	End   time.Time `form:"end" time_format:"2006-01-02T15:04:05" time_utc:"1"`
}

func getQuery(c *gin.Context) (QueryAll, error) {
	var query QueryAll
	if query.Start.IsZero() {
		query.Start = time.Now().Add(time.Hour * time.Duration(-24))
	}
	if query.End.IsZero() {
		query.End = time.Now()
	}
	if c.ShouldBind(&query) == nil {
		return query, nil
	} else {
		return query, errors.New("wrong query")
	}
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	db := core.Create()

	config, err := (&core.Config{}).Load()
	if err != nil {
		fmt.Println(err)
	}

	sg := core.StatsGenerator{ApiHost: config.ApiHost, ApiPort: config.ApiPort, DB: db}
	r.GET("/api/traffic", func(c *gin.Context) {
		var traffics []core.Traffic
		query, err := getQuery(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, nil)
		}
		db.Table("traffics").
			Select("`date`, data_type, name, 'all' as type, sum(value) as value").
			Where("date BETWEEN ? AND ?", query.Start, query.End).
			Group("`name`, `date`, `data_type`").
			Order("date").
			Find(&traffics)
		c.JSON(http.StatusOK, traffics)
	})

	r.GET("/view/traffic", func(c *gin.Context) {
		m := make(map[string]interface{})
		m["data_api"] = "/api/traffic"
		c.HTML(http.StatusOK, "traffic.html", m)
	})

	r.GET("/api/traffic/user/:name", func(c *gin.Context) {
		var traffics []core.Traffic
		query, err := getQuery(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, nil)
		}
		name := c.Param("name")
		if query.End.IsZero() {
			query.End = time.Now()
		}
		db.Table("traffics").
			Select("`date`, data_type, name, 'all' as type, sum(value) as value").
			Where("date BETWEEN ? AND ? AND name = ? AND data_type = 'user'", query.Start, query.End, name).
			Group("`name`, `date`, `data_type`").
			Order("date").
			Find(&traffics)
		c.JSON(http.StatusOK, traffics)
	})

	r.GET("/view/traffic/user/:name", func(c *gin.Context) {
		query, err := getQuery(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, nil)
		}
		name := c.Param("name")
		m := make(map[string]interface{})
		m["name"] = name
		m["query"] = query
		m["data_api"] = "/api/traffic/user/" + name
		c.HTML(http.StatusOK, "traffic.html", m)
	})

	r.GET("/api/user", func(c *gin.Context) {
		var traffics []core.Traffic
		db.Table("traffics").
			Select("name").
			Where("data_type='user'").
			Group("name").Scan(&traffics)
		names := make([]string, len(traffics))
		for i, v := range traffics {
			names[i] = v.Name
		}
		m := make(map[string]interface{})
		m["users"] = names
		c.JSON(http.StatusOK, m)
	})

	r.GET("/view/user", func(c *gin.Context) {
		var traffics []core.Traffic
		db.Table("traffics").
			Select("name").
			Where("data_type='user'").
			Group("name").Scan(&traffics)
		names := make([]string, len(traffics))
		for i, v := range traffics {
			names[i] = v.Name
		}
		m := make(map[string]interface{})
		m["users"] = names
		c.HTML(http.StatusOK, "user.html", m)
	})

	sg.Start()
	err = r.Run(fmt.Sprintf("%s:%d", "0.0.0.0", config.ServerPort))
	if err != nil {
		log.Fatalln(err)
	}
}
