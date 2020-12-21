import request from '@/utils/request'

export function requestList(data) {
  return request({
    url: '/role/list',
    method: 'post',
    data
  })
}

export function requestDetail(id) {
  return request({
    url: '/role/detail',
    method: 'get',
    params: { id }
  })
}

export function requestUpdate(data) {
  return request({
    url: '/role/update',
    method: 'post',
    data
  })
}

export function requestCreate(data) {
  return request({
    url: '/role/create',
    method: 'post',
    data
  })
}

export function requestDelete(data) {
  return request({
    url: '/role/delete',
    method: 'post',
    data
  })
}

export function requestRoleMenuIDList(roleid) {
  return request({
    url: '/role/rolemenuidlist',
    method: 'get',
    params: { roleid }
  })
}

export function requestSetRole(data) {
  return request({
    url: '/role/setrole',
    method: 'post',
    data
  })
}

export function requestAll() {
  return request({
    url: '/role/allrole',
    method: 'get'
  })
}
