import request from '@/utils/request'

export function createExampleDeploy() {
  return request({
    url: `/cluster/deployment/create`,
    method: 'get',
  })
}


export function getDeploymentsList(data) {
  return request({
    url: `/cluster/deployment/list`,
    method: 'post',
    data
  })
}

export function pathNamespaceDeployments(data) {
  return request({
    url: `/cluster/deployment/update`,
    method: 'post',
    data
  })
}

export function deleteNamespaceDeployment(data) {
  return request({
    url: `/cluster/deployment/delete`,
    method: 'post',
    data
  })
}

