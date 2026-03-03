import axios from 'axios'

// 设置axios的基础URL，所有API请求都会自动添加/api前缀
axios.defaults.baseURL = '/api'

// 添加请求拦截器，在请求头中添加用户名
axios.interceptors.request.use(
  (config) => {
    // 从本地存储中获取用户名
    const username = localStorage.getItem('username')
    if (username) {
      config.headers['X-Username'] = username
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 告警相关API
export const alertAPI = {
  list: () => axios.get('/alerts'),
  get: (id: number) => axios.get(`/alerts/${id}`),
  resolve: (id: number) => axios.put(`/alerts/${id}/resolve`),
  delete: (id: number) => axios.delete(`/alerts/${id}`)
}

// 告警介质相关API
export const mediaAPI = {
  list: () => axios.get('/media'),
  create: (data: any) => axios.post('/media', data),
  get: (id: number) => axios.get(`/media/${id}`),
  update: (id: number, data: any) => axios.put(`/media/${id}`, data),
  delete: (id: number) => axios.delete(`/media/${id}`)
}

// 告警模板相关API
export const templateAPI = {
  list: () => axios.get('/templates'),
  create: (data: any) => axios.post('/templates', data),
  get: (id: number) => axios.get(`/templates/${id}`),
  update: (id: number, data: any) => axios.put(`/templates/${id}`, data),
  delete: (id: number) => axios.delete(`/templates/${id}`)
}

// 静默规则相关API
export const silenceAPI = {
  list: () => axios.get('/silences'),
  create: (data: any) => axios.post('/silences', data),
  get: (id: number) => axios.get(`/silences/${id}`),
  delete: (id: number) => axios.delete(`/silences/${id}`)
}

// 抑制规则相关API
export const inhibitionAPI = {
  list: () => axios.get('/inhibitions'),
  create: (data: any) => axios.post('/inhibitions', data),
  get: (id: number) => axios.get(`/inhibitions/${id}`),
  update: (id: number, data: any) => axios.put(`/inhibitions/${id}`, data),
  delete: (id: number) => axios.delete(`/inhibitions/${id}`)
}

// Prometheus相关API
export const prometheusAPI = {
  alerts: () => axios.get('/prometheus/alerts'),
  rules: () => axios.get('/prometheus/rules'),
  targets: () => axios.get('/prometheus/targets')
}

// 用户相关API
export const userAPI = {
  // 获取当前用户信息
  getCurrentUser: async () => {
    const response = await axios.get('/users/current')
    return response.data
  },
  // 更新用户信息
  updateProfile: async (data: any) => {
    const response = await axios.put('/users/profile', data)
    return response.data
  },
  // 修改密码
  changePassword: async (data: any) => {
    const response = await axios.post('/users/change-password', data)
    return response.data
  },
  // 获取用户列表
  list: async () => {
    const response = await axios.get('/users')
    return response.data
  },
  // 创建用户
  create: async (data: any) => {
    const response = await axios.post('/users', data)
    return response.data
  },
  // 获取用户详情
  get: async (id: number) => {
    const response = await axios.get(`/users/${id}`)
    return response.data
  },
  // 更新用户
  update: async (id: number, data: any) => {
    console.log('调用userAPI.update:', id, data)
    const response = await axios.put(`/users/${id}`, data)
    console.log('userAPI.update响应:', response.data)
    return response.data
  },
  // 删除用户
  delete: async (id: number) => {
    const response = await axios.delete(`/users/${id}`)
    return response.data
  }
}

// 主机管理相关API
export const hostAPI = {
  // 获取主机列表
  list: () => axios.get('/hosts'),
  // 创建主机
  create: (data: any) => axios.post('/hosts', data),
  // 获取主机详情
  get: (id: number) => axios.get(`/hosts/${id}`),
  // 更新主机
  update: (id: number, data: any) => axios.put(`/hosts/${id}`, data),
  // 删除主机
  delete: (id: number) => axios.delete(`/hosts/${id}`),
  // 测试连接
  testConnection: (id: number) => axios.post(`/hosts/${id}/test-connection`),
  // 执行SSH命令
  openSSH: (id: number, data: any) => axios.post(`/hosts/${id}/ssh/command`, data),
  // 上传文件
  uploadFile: (id: number, data: FormData) => axios.post(`/hosts/${id}/ssh/upload`, data, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  }),
  // 下载文件
  downloadFile: (id: number, params: any) => axios.get(`/hosts/${id}/ssh/download`, {
    params,
    responseType: 'blob'
  }),
  // 重启主机
  restartHost: (id: number) => axios.post(`/hosts/${id}/ssh/restart`),
  // 关机主机
  shutdownHost: (id: number) => axios.post(`/hosts/${id}/ssh/shutdown`),
  // 获取服务器资源统计
  getServerStats: (id: number) => axios.get(`/hosts/${id}/stats`)
}

// 系统设置相关API
export const systemSettingAPI = {
  // 获取所有系统设置
  list: () => axios.get('/settings'),
  // 获取单个系统设置
  get: (key: string) => axios.get(`/settings/${key}`),
  // 创建系统设置
  create: (data: any) => axios.post('/settings', data),
  // 更新系统设置
  update: (key: string, data: any) => axios.put(`/settings/${key}`, data),
  // 删除系统设置
  delete: (key: string) => axios.delete(`/settings/${key}`)
}

// 邮件相关API
export const emailAPI = {
  // 测试发送邮件
  testSend: (toEmail: string) => axios.post('/email/test', { toEmail })
}

// 邮箱配置相关API
export const emailSettingAPI = {
  // 获取邮箱配置
  get: () => axios.get('/email-settings'),
  // 更新邮箱配置
  update: (data: any) => axios.put('/email-settings', data)
}

// 集群相关API
export const clusterAPI = {
  // 获取集群列表
  list: () => axios.get('/clusters'),
  // 创建集群
  create: (data: any) => axios.post('/clusters', data),
  // 更新集群
  update: (id: number, data: any) => axios.put(`/clusters/${id}`, data),
  // 删除集群
  delete: (id: number) => axios.delete(`/clusters/${id}`),
  // 测试集群连接
  testConnection: (id: number) => axios.post(`/clusters/${id}/test-connection`),
  // 获取指定集群的命名空间列表
  getNamespaces: (id: number) => axios.get(`/clusters/${id}/namespaces`)
}

// ConfigMap相关API
export const configMapAPI = {
  // 获取ConfigMap列表
  list: (clusterId: number, namespace?: string) => axios.get(`/clusters/${clusterId}/configmaps`, {
    params: namespace ? { namespace } : {}
  })
}

// Secret相关API
export const secretAPI = {
  // 获取Secret列表
  list: (clusterId: number, namespace?: string) => axios.get(`/clusters/${clusterId}/secrets`, {
    params: namespace ? { namespace } : {}
  })
}

// 节点相关API
export const nodeAPI = {
  // 获取节点列表
  list: (clusterId: number) => axios.get(`/clusters/${clusterId}/nodes`)
}

// Pod相关API
export const podAPI = {
  // 获取Pod列表
  list: (clusterId: number, namespace?: string) => axios.get(`/clusters/${clusterId}/pods`, {
    params: namespace ? { namespace } : {}
  }),
  // 获取Pod日志
  getLogs: (clusterId: number, name: string, namespace: string) => axios.get(`/clusters/${clusterId}/pods/${name}/logs`, {
    params: { namespace }
  })
}

// Deployment相关API
export const deploymentAPI = {
  // 获取Deployment列表
  list: (clusterId: number, namespace?: string) => axios.get(`/clusters/${clusterId}/deployments`, {
    params: namespace ? { namespace } : {}
  })
}

// Service相关API
export const serviceAPI = {
  // 获取Service列表
  list: (clusterId: number, namespace?: string) => axios.get(`/clusters/${clusterId}/services`, {
    params: namespace ? { namespace } : {}
  })
}

// 事件相关API
export const eventAPI = {
  // 获取事件列表
  list: (clusterId: number, params?: any) => axios.get(`/clusters/${clusterId}/events`, {
    params
  })
}

// SLS配置相关API
export const slsConfigAPI = {
  // 获取SLS配置列表
  list: () => axios.get('/sls/configs'),
  // 创建SLS配置
  create: (data: any) => axios.post('/sls/configs', data),
  // 更新SLS配置
  update: (id: number, data: any) => axios.put(`/sls/configs/${id}`, data),
  // 删除SLS配置
  delete: (id: number) => axios.delete(`/sls/configs/${id}`),
  // 激活SLS配置
  activate: (id: number) => axios.put(`/sls/configs/${id}/active`)
}

// KVM管理相关API
export const kvmAPI = {
  // 获取KVM主机列表
  list: () => axios.get('/kvm/hosts'),
  // 创建KVM主机
  create: (data: any) => axios.post('/kvm/hosts', data),
  // 获取KVM主机详情
  get: (id: number) => axios.get(`/kvm/hosts/${id}`),
  // 更新KVM主机
  update: (id: number, data: any) => axios.put(`/kvm/hosts/${id}`, data),
  // 删除KVM主机
  delete: (id: number) => axios.delete(`/kvm/hosts/${id}`),
  // 测试KVM主机连接
  testConnection: (id: number) => axios.post(`/kvm/hosts/${id}/test`),
  // 获取KVM主机下的虚拟机列表
  getVmList: (id: number) => axios.get(`/kvm/hosts/${id}/vms`),
  // 启动虚拟机
  startVm: (hostId: number, vmId: number) => axios.post(`/kvm/hosts/${hostId}/vms/${vmId}/start`),
  // 关闭虚拟机
  shutdownVm: (hostId: number, vmId: number) => axios.post(`/kvm/hosts/${hostId}/vms/${vmId}/shutdown`),
  // 强制关机虚拟机
  forceShutdownVm: (hostId: number, vmId: number) => axios.post(`/kvm/hosts/${hostId}/vms/${vmId}/destroy`),
  // 重启虚拟机
  rebootVm: (hostId: number, vmId: number) => axios.post(`/kvm/hosts/${hostId}/vms/${vmId}/reboot`)
}

// 仪表盘相关API
export const dashboardAPI = {
  // 获取仪表盘统计数据
  getStats: () => axios.get('/dashboard/stats'),
  // 获取资源概览数据
  getResources: () => axios.get('/dashboard/resources'),
  // 获取系统信息
  getSystemInfo: () => axios.get('/dashboard/system-info'),
  // 获取系统健康状态数据
  getHealthStatus: () => axios.get('/dashboard/health-status')
}

// 日志相关API
export const logAPI = {
  // 操作日志
  getOperationLogs: (params: any) => axios.get('/logs/operation', { params }),
  deleteOperationLog: (id: number) => axios.delete(`/logs/operation/${id}`),
  batchDeleteOperationLogs: (ids: number[]) => axios.delete('/logs/operation', { data: { ids } }),
  
  // 登录日志
  getLoginLogs: (params: any) => axios.get('/logs/login', { params }),
  deleteLoginLog: (id: number) => axios.delete(`/logs/login/${id}`),
  batchDeleteLoginLogs: (ids: number[]) => axios.delete('/logs/login', { data: { ids } })
}

// AI相关API
export const aiAPI = {
  // AI聊天
  chat: (data: any) => axios.post('/ai/chat', data),
  // 测试AI连接
  testConnection: (data: any) => axios.post('/ai/test-connection', data)
}

// 登录/注册相关API
export const authAPI = {
  // 用户登录
  login: (username: string, password: string) => axios.post('/users/login', {
    username,
    password
  }),
  // 用户注册
  register: (username: string, password: string, email: string) => axios.post('/users/register', {
    username,
    password,
    email
  }),
  // 忘记密码 - 验证邮箱
  verifyEmail: (email: string) => axios.post('/users/forgot-password/verify-email', {
    email
  }),
  // 忘记密码 - 重置密码
  resetPassword: (email: string, username: string) => axios.post('/users/forgot-password/reset', {
    email,
    username
  })
}

// MD5相关API
export const md5API = {
  // 创建MD5加密历史记录
  createHistory: (data: any) => axios.post('/md5-histories', data),
  // 获取MD5加密历史记录
  getHistories: () => axios.get('/md5-histories'),
  // 删除单个MD5加密历史记录
  deleteHistory: (id: number) => axios.delete(`/md5-histories/${id}`),
  // 清空MD5加密历史记录
  clearHistories: () => axios.delete('/md5-histories')
}

// 图片压缩相关API
export const imageCompressAPI = {
  // 创建图片压缩历史记录
  createHistory: (data: any) => axios.post('/image-compress-histories', data),
  // 获取图片压缩历史记录
  getHistories: () => axios.get('/image-compress-histories'),
  // 删除单个图片压缩历史记录
  deleteHistory: (id: number) => axios.delete(`/image-compress-histories/${id}`),
  // 清空图片压缩历史记录
  clearHistories: () => axios.delete('/image-compress-histories')
}