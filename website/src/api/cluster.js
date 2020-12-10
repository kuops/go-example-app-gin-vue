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

export function deleteNamespace(data) {
  return request({
    url: '/vue-admin-template/cluster/namespaces',
    method: 'delete',
    data
  })
}
