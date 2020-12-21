import request from '@/utils/request'

export function getNamespace(data) {
  return request({
    url: '/cluster/namespace/list',
    method: 'post',
    data
  })
}

export function createNamespace(data) {
  return request({
    url: '/cluster/namespace/create',
    method: 'post',
    data
  })
}

export function deleteNamespace(data) {
  return request({
    url: `/cluster/namespace/delete`,
    method: 'post',
    data
  })
}
