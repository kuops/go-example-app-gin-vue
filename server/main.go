package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

type Info struct {
	Roles        []string `json:"roles"`
	Introduction string   `json:"introduction"`
	Avatar       string   `json:"avatar"`
	Name         string   `json:"name"`
}

type Table struct {
	Total int    `json:"total"`
	Items []Item `json:"items"`
}

type Item struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Time      string `json:"display_time"`
	Pageviews uint32 `json:"pageviews"`
	Status    string `json:"status"`
}

func main() {

	var info = &Info{
		Avatar:       "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		Roles:        []string{"admin"},
		Introduction: "I am a super administrato",
		Name:         "Admin",
	}

	var token = map[string]string{
		"token": "admin-token",
	}

	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = []string{"*"}
	r.Use(cors.New(config))

	r.GET("api/v1/vue-admin-template/user/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": 20000,
			"data": info,
		})
	})

	r.POST("api/v1/vue-admin-template/user/login", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": 20000,
			"data": token,
		})
	})

	r.POST("api/v1/vue-admin-template/user/logout", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": 20000,
			"data": "success",
		})
	})

	var items = []Item{
		{
			ID:        0,
			Title:     "First",
			Author:    "anonymous",
			Time:      time.Now().Format("2006-01-02 15:01:05"),
			Pageviews: 6666,
			Status:    "draft",
		},
		{
			ID:        1,
			Title:     "Second",
			Author:    "anonymous",
			Time:      time.Now().Format("2006-01-02 15:01:05"),
			Pageviews: 7777,
			Status:    "published",
		},
		{
			ID:        1,
			Title:     "third",
			Author:    "anonymous",
			Time:      time.Now().Format("2006-01-02 15:01:05"),
			Pageviews: 5555,
			Status:    "deleted",
		},
	}
	var table = &Table{
		Total: len(items),
		Items: items,
	}

	r.GET("api/v1/vue-admin-template/table/list", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": 20000,
			"data": table,
		})
	})

	r.Run(":8080")
}
