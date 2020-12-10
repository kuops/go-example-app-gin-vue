package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
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

	kubeconfig, err := clientcmd.BuildConfigFromFlags("", "/Users/houshiying/.kube/config")
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(kubeconfig)


	r.GET("api/v1/vue-admin-template/cluster/namespaces", func(c *gin.Context) {
		namespaces, _ := clientset.CoreV1().Namespaces().List(context.TODO(),metav1.ListOptions{})
		c.JSON(200, gin.H{
			"code": 20000,
			"data": namespaces,
		})
	})

	var Namespace = &v1.Namespace{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       v1.NamespaceSpec{},
		Status:     v1.NamespaceStatus{},
	}

	type CreateNameSpace struct {
		Name string `json:"name"`
	}
	var createNS = CreateNameSpace{}
	var status string
	var code int

	r.POST("api/v1/vue-admin-template/cluster/namespaces", func(c *gin.Context) {
		err := c.ShouldBindJSON(&createNS)
		if err == nil {
			Namespace.ObjectMeta.Name = createNS.Name
			if _,err := clientset.CoreV1().Namespaces().Create(context.TODO(),Namespace,metav1.CreateOptions{});err != nil {
				status  = "failed"
				code = 20001
			}
				status = "success"
				code = 20000
		} else {
			status = "error"
			code = 20002
		}
		c.JSON(200, gin.H{
			"code": code,
			"data": map[string]string{
				"state": status,
			},
		})
	})

	type DleteNameSpace struct {
		Name string `json:"name"`
	}
	var deleteNS = DleteNameSpace{}
	r.DELETE("api/v1/vue-admin-template/cluster/namespaces", func(c *gin.Context) {
		err := c.ShouldBindJSON(&deleteNS)
		if err == nil {
			if err := clientset.CoreV1().Namespaces().Delete(context.TODO(),deleteNS.Name,metav1.DeleteOptions{});err != nil {
				status  = "failed"
				code = 20001
			}
			status = "success"
			code = 20000
		} else {
			status = "error"
			code = 20002
		}
		c.JSON(200, gin.H{
			"code": code,
			"data": map[string]string{
				"state": status,
			},
		})
	})

	r.Run(":8080")
}
