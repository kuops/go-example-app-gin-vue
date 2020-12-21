package request

type CreateNameSpace struct {
	Name string `json:"name"`
}

type DeleteNameSpace struct {
	Name string `json:"name"`
}

type DeleteDeployment struct {
	Name string `json:"name"`
	Namespace string `json:"namespace"`
}

type NamespacePageInfo struct {
	Page     uint64 `json:"page" form:"page"`
	PageSize uint64 `json:"pageSize" form:"pageSize"`
}

type NamespaceList struct {
	RolePageInfo
	Key  string `json:"key"`
	ParentID uint64 `json:"parent_id"`
}


type DeploymentList struct {
	RolePageInfo
	Key  string `json:"key"`
	ParentID uint64 `json:"parent_id"`
	Namespace string `json:"namespace"`
}
