import request from '@/utils/request'

export function requestEditPwd(data) {
  return request({
    url: '/user/editpwd',
    method: 'post',
    data
  })
}

export function requestList(data) {
  return request({
    url: '/user/list',
    method: 'post',
    data
  })
}

export function requestDetail(id) {
  return request({
    url: '/user/detail',
    method: 'get',
    params: { id }
  })
}

export function requestUpdate(data) {
  return request({
    url: '/user/update',
    method: 'post',
    data
  })
}

export function requestCreate(data) {
  return request({
    url: '/user/create',
    method: 'post',
    data
  })
}

export function requestDelete(data) {
  return request({
    url: '/user/delete',
    method: 'post',
    data
  })
}

export function requestUsersRoleIDList(id) {
  return request({
    url: '/user/usersroleidlist',
    method: 'get',
    params: { id }
  })
}

export function requestSetRole(id, data) {
  return request({
    url: '/user/setrole',
    method: 'post',
    params: { id },
    data
  })
}

