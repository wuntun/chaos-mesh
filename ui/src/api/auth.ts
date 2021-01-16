import http from 'api/http'

let tokenInterceptorId: number

export const token = (token: string) => {
  if (tokenInterceptorId !== undefined) {
    http.interceptors.request.eject(tokenInterceptorId)
  }

  tokenInterceptorId = http.interceptors.request.use((config) => {
    config.headers = {
      ...config.headers,
      Authorization: `Bearer ${token}`,
    }

    return config
  })
}

export const resetToken = () => http.interceptors.request.eject(tokenInterceptorId)

let nsInterceptorId: number

export const namespace = (ns: string) => {
  if (nsInterceptorId !== undefined) {
    http.interceptors.request.eject(nsInterceptorId)
  }

  nsInterceptorId = http.interceptors.request.use((config) => {
    if (/state|experiments|events(\/dry)?|archives$/g.test(config.url!)) {
      config.params = {
        ...config.params,
        namespace: ns === 'All' ? null : ns,
      }
    }

    return config
  })
}

let nodeInterceptorId: number

export const node = (name: string) => {
  if (nodeInterceptorId !== undefined) {
    http.interceptors.request.eject(nodeInterceptorId)
  }

  nodeInterceptorId = http.interceptors.request.use((config) => {
    config.params = {
      ...config.params,
      name,
    }

    return config
  })
}

export const resetNode = () => http.interceptors.request.eject(nodeInterceptorId)
