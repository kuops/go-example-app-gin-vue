package v1

import (
	"context"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/kuops/go-example-app/server/pkg/log"
	"github.com/kuops/go-example-app/server/pkg/request"
	"github.com/kuops/go-example-app/server/pkg/response"
	"github.com/kuops/go-example-app/server/pkg/store/mysql"
	"github.com/kuops/go-example-app/server/pkg/store/redis"
	corev1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

type handler struct {
	dao *dao
	cache redis.Interface
	enforcer *casbin.Enforcer
	clientset *kubernetes.Clientset
	dynamic dynamic.Interface
}

func newHandler(mysqlClient *mysql.Client,redisClient redis.Interface,enforcer *casbin.Enforcer,clientset *kubernetes.Clientset,dynamic dynamic.Interface) *handler {
	return &handler{
		dao: &dao{
			db: mysqlClient.Database().DB,
		},
		cache: redisClient,
		enforcer: enforcer,
		dynamic: dynamic,
		clientset: clientset,
	}
}

// @Tags 命名空间
// @Summary 分页查看命名空间列表
// @Produce  application/json
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.NamespaceList true "data"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /api/v1/cluster/namespace/list [get]
func (h *handler)ListNamespace(c *gin.Context) {
	var req request.NamespaceList
	_ = c.ShouldBindJSON(&req)
	var err error
	var total int64

	if req.PageSize == 0 || req.Page == 0 {
		log.Error("分页参数参数非法")
		response.FailWithMessage("参数不正确", c)
		return
	}

	namespaceRes := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "namespaces"}
	namespaces, err := h.dynamic.Resource(namespaceRes).List(context.TODO(),metav1.ListOptions{})

	if err != nil {
		log.Errorf("获取失败 %v",err)
		response.FailWithMessage("获取失败", c)
		return
	}

	var itemsList []unstructured.Unstructured

	if len(namespaces.Items) > 0 {
		if req.Key != "" {
			for _, v := range namespaces.Items {
				name := v.GetName()
				if strings.Contains(name, req.Key) {
					itemsList = append(itemsList, v)
				}
			}
			namespaces = &unstructured.UnstructuredList{
				Object: namespaces.Object,
				Items:  itemsList,
			}
		}
	}

	total = int64(len(namespaces.Items))
	offset := (req.Page -1) * req.PageSize
	limit:= offset + req.PageSize

	if offset < uint64(total) {
		if limit < uint64(total) {
			namespaces = &unstructured.UnstructuredList{
				Object: namespaces.Object,
				Items: namespaces.Items[offset:limit],
			}
		} else {
			namespaces = &unstructured.UnstructuredList{
				Object: namespaces.Object,
				Items: namespaces.Items[offset:len(namespaces.Items)],
			}
		}
	} else {
		namespaces = &unstructured.UnstructuredList{
			Object: namespaces.Object,
			Items: nil,
		}
	}

	response.OkWithDetailed(response.PageResult{
		List:     namespaces,
		Total:    total,
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	}, "获取成功", c)
}


// @Tags 命名空间
// @Summary 创建命名空间
// @Produce  application/json
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.CreateNameSpace true "data"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /api/v1/cluster/namespace/create [post]
func (h *handler)CreateNamespace(c *gin.Context) {
	var ns request.CreateNameSpace
	var Namespace = &corev1.Namespace{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       corev1.NamespaceSpec{},
		Status:     corev1.NamespaceStatus{},
	}

	err := c.ShouldBindJSON(&ns)
	if err != nil {
		response.FailWithMessage("缺少必填字段", c)
		return
	}
	Namespace.ObjectMeta.Name = ns.Name
	if _,err := h.clientset.CoreV1().Namespaces().Create(context.TODO(),Namespace,metav1.CreateOptions{});err != nil {
			log.Errorf("获取失败 %v",err)
			response.FailWithMessage("获取失败", c)
			return
	}
	response.OkWithDetailed(Namespace,"创建成功",c)
}

// @Tags 命名空间
// @Summary 删除命名空间
// @Produce  application/json
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.DeleteNameSpace true "data"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /api/v1/cluster/namespace/delete [delete]
func (h *handler)DeleteNamespace(c *gin.Context) {
	var ns request.DeleteNameSpace
	err := c.ShouldBindJSON(&ns)
	if err != nil {
		response.FailWithMessage("缺少必填字段", c)
		return
	}
	if err := h.clientset.CoreV1().Namespaces().Delete(context.TODO(),ns.Name,metav1.DeleteOptions{});err != nil {
		log.Errorf("删除失败 %v",err)
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功",c)
}


// @Tags 无状态部署
// @Summary 查看无状态部署
// @Produce  application/json
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.DeploymentList true "data"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /api/v1/cluster/deployment/list [post]
func (h *handler)ListDeployment(c *gin.Context) {
	var req request.DeploymentList
	_ = c.ShouldBindJSON(&req)
	var err error
	var total int64
	var deployments *unstructured.UnstructuredList

	if req.PageSize == 0 || req.Page == 0 {
		log.Error("分页参数参数非法")
		response.FailWithMessage("参数不正确", c)
		return
	}

	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	deployments, err = h.dynamic.Resource(deploymentRes).List(context.TODO(),metav1.ListOptions{})

	if req.Namespace != "" {
		deployments, err = h.dynamic.Resource(deploymentRes).Namespace(req.Namespace).List(context.TODO(),metav1.ListOptions{})
	}

	if err != nil {
		log.Errorf("获取失败 %v",err)
		response.FailWithMessage("获取失败", c)
		return
	}


	var itemsList []unstructured.Unstructured

	if len(deployments.Items) > 0 {
		if req.Key != "" {
			for _, v := range deployments.Items {
				name := v.GetName()
				if strings.Contains(name, req.Key) {
					itemsList = append(itemsList, v)
				}
			}
			deployments = &unstructured.UnstructuredList{
				Object: deployments.Object,
				Items:  itemsList,
			}
		}
	}

	total = int64(len(deployments.Items))
	offset := (req.Page -1) * req.PageSize
	limit:= offset + req.PageSize

	if offset < uint64(total) {
		if limit < uint64(total) {
			deployments = &unstructured.UnstructuredList{
				Object: deployments.Object,
				Items: deployments.Items[offset:limit],
			}
		} else {
			deployments = &unstructured.UnstructuredList{
				Object: deployments.Object,
				Items: deployments.Items[offset:len(deployments.Items)],
			}
		}
	} else {
		deployments = &unstructured.UnstructuredList{
			Object: deployments.Object,
			Items: nil,
		}
	}

	response.OkWithDetailed(response.PageResult{
		List:     deployments,
		Total:    total,
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	}, "获取成功", c)
}

// @Tags 无状态部署
// @Summary 更新无状态部署
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body req true "ID, 父级ID, URL"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新菜单信息成功"}"
// @Router /api/v1/cluster/deployment/update [post]
func (h *handler) UpdateDeployment(c *gin.Context) {
	var req appsv1.Deployment
	var namespace string
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Errorf("缺少必填字段",err)
		response.FailWithMessage("缺少必填字段", c)
		return
	}
	if req.ObjectMeta.Namespace == "" {
		namespace = "default"
	} else {
		namespace = req.ObjectMeta.Namespace
	}

	deployments, err := h.clientset.AppsV1().Deployments(namespace).Update(context.TODO(),&req,metav1.UpdateOptions{})

	if err != nil {
		log.Errorf("更新失败 %v",err)
		response.FailWithMessage("更新失败", c)
		return
	}

	response.OkWithDetailed(deployments, "更新成功", c)
}


// @Tags 无状态部署
// @Summary 更新无状态部署
// @Produce  application/json
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.DeleteDeployment true "data"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /api/v1/cluster/deployment/delete [post]

func (h *handler)DeleteDeployment(c *gin.Context) {
	var req request.DeleteDeployment
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage("缺少必填字段", c)
		return
	}
	if err := h.clientset.AppsV1().Deployments(req.Namespace).Delete(context.TODO(),req.Name,metav1.DeleteOptions{}); err != nil {
		log.Errorf("删除失败 %v",err)
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功",c)
}


// @Tags 无状态部署
// @Summary 创建无状态部署
// @Produce  application/json
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.DeleteDeployment true "data"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /api/v1/cluster/deployment/create [post]

func (h *handler)CreateDeployment(c *gin.Context) {
	deploymentSpec := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: func() *int32 { i := int32(2); return &i }(),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "web",
							Image: "nginx:latest",
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
	if deploy,err := h.clientset.AppsV1().Deployments("default").Create(context.TODO(),deploymentSpec,metav1.CreateOptions{}); err != nil {
		log.Errorf("创建失败 %v",err)
		response.FailWithMessage("创建失败", c)
		return
	} else {
		response.OkWithDetailed(deploy,"创建成功",c)
	}
}
