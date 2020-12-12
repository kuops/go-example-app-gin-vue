import request from '@/utils/request'

export function getNamespaces() {
  return request({
    url: '/vue-admin-template/cluster/namespaces',
    method: 'get'
  })
}

export function createNamespace(data) {
  return request({
    url: '/vue-admin-template/cluster/namespaces',
    method: 'post',
    data
  })
}

export function deleteNamespace(name) {
  return request({
    url: `/vue-admin-template/cluster/namespaces/${name}`,
    method: 'delete',
  })
}

export function getDeployments() {
  return request({
    url: `/vue-admin-template/cluster/deployments`,
    method: 'get',
  })
}

export function getNamespaceDeployments(name) {
  return request({
    url: `/vue-admin-template/cluster/namespaces/${name}/deployments`,
    method: 'get',
  })
}

export function pathNamespaceDeployments(namespace,deployment,data) {
  return request({
    url: `/vue-admin-template/cluster/namespaces/${namespace}/deployments/${deployment}`,
    method: 'patch',
    data
  })
}

export function deleteNamespaceDeployment(namespace,deployment) {
  return request({
    url: `/vue-admin-template/cluster/namespaces/${namespace}/deployments/${deployment}`,
    method: 'delete',
  })
}
