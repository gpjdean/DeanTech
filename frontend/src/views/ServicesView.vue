<template>
  <div class="services-view">
    <el-card shadow="never" class="page-header">
      <div class="header-content">
        <h2>服务</h2>
        <div class="header-actions">
          <el-select v-model="selectedCluster" placeholder="选择集群" style="width: 200px; margin-right: 10px;">
            <el-option
              v-for="cluster in clusters"
              :key="cluster.id"
              :label="cluster.name"
              :value="cluster.id"
            />
          </el-select>
          <el-select v-model="selectedNamespace" placeholder="选择命名空间" style="width: 200px; margin-right: 10px;" :loading="namespacesLoading" filterable>
            <el-option
              v-for="namespace in namespaces"
              :key="namespace"
              :label="namespace"
              :value="namespace"
            />
          </el-select>
          <el-button type="primary" size="small" @click="refreshServices">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </div>
    </el-card>

    <el-card shadow="hover" class="services-card">
      <template #header>
        <div class="card-header">
          <span>K8s Service列表</span>
        </div>
      </template>
      
      <el-table 
        :data="services" 
        style="width: 100%; color: #303133; font-size: 14px; line-height: 1.5;"
        stripe
        v-loading="loading"
        v-if="selectedCluster"
        :row-style="{ height: 'auto', backgroundColor: '#ffffff' }"
        :cell-style="{ 
          'white-space': 'normal', 
          'word-break': 'break-all', 
          'min-height': '40px', 
          'padding': '8px 12px', 
          'color': '#303133', 
          'font-size': '14px',
          'line-height': '1.5'
        }"
        :header-cell-style="{ 
          'color': '#303133', 
          'font-weight': 'bold', 
          'background-color': '#f5f7fa',
          'font-size': '14px'
        }"
        border
        :row-key="(row: any) => `${row.namespace}-${row.name}`"
      >
        <el-table-column prop="name" label="Service名称" min-width="200" show-overflow-tooltip />
        <el-table-column prop="namespace" label="命名空间" width="150" show-overflow-tooltip />
        <el-table-column prop="cluster" label="集群" width="150" show-overflow-tooltip />
        <el-table-column prop="type" label="类型" width="120">
          <template #default="scope">
            <el-tag 
              :type="getTypeTagType(scope.row.type)" 
              size="small" 
              effect="dark" 
              style="border-radius: 6px; padding: 3px 10px; font-size: 12px; font-weight: 500;"
            >
              {{ scope.row.type || 'Unknown' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="clusterIP" label="Cluster IP" width="150" show-overflow-tooltip />
        <el-table-column prop="externalIP" label="外部IP" width="150" show-overflow-tooltip />
        <el-table-column prop="ports" label="端口" min-width="200">
          <template #default="scope">
            <div v-for="port in scope.row.ports" :key="port.port" class="port-item" style="color: #303133; font-size: 14px; margin: 4px 0;">
              {{ port.port }}/{{ port.protocol }} → {{ port.targetPort }}
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="selector" label="标签选择器" min-width="250" show-overflow-tooltip>
          <template #default="scope">
            <div style="display: flex; flex-wrap: wrap; gap: 6px; min-height: 36px; align-content: flex-start;">
              <el-tag v-if="Object.keys(scope.row.selector).length === 0" type="info" size="small" effect="light">
                无
              </el-tag>
              <el-tag
                v-for="(value, key) in scope.row.selector"
                :key="key"
                type="info"
                size="small"
                effect="light"
                style="margin: 0; max-width: calc(100% - 12px); overflow: hidden; text-overflow: ellipsis; white-space: nowrap;"
              >
                {{ key }}={{ value }}
              </el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="120" align="center">
          <template #default="scope">
            <span style="font-weight: 500; color: #606266;">
              {{ formatAge(scope.row.createdAt) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" @click="viewServiceDetails(scope.row)">
              <el-icon><View /></el-icon>
              查看
            </el-button>
            <el-button type="danger" size="small" @click="deleteService(scope.row)">
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <div v-else class="no-cluster-selected">
        <el-empty description="请选择一个集群查看Service列表" />
      </div>
    </el-card>
    
    <!-- Service详情对话框 -->
    <el-dialog
      v-model="serviceDetailsVisible"
      title="Service详情"
      width="800px"
    >
      <div v-if="selectedService" class="service-details">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="Service名称">{{ selectedService.name }}</el-descriptions-item>
          <el-descriptions-item label="命名空间">{{ selectedService.namespace }}</el-descriptions-item>
          <el-descriptions-item label="类型">{{ selectedService.type }}</el-descriptions-item>
          <el-descriptions-item label="Cluster IP">{{ selectedService.clusterIP }}</el-descriptions-item>
          <el-descriptions-item label="外部IP">{{ selectedService.externalIP || 'N/A' }}</el-descriptions-item>
          <el-descriptions-item label="外部名称">{{ selectedService.externalName || 'N/A' }}</el-descriptions-item>
          <el-descriptions-item label="会话亲和性">{{ selectedService.sessionAffinity || 'N/A' }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ selectedService.createdAt }}</el-descriptions-item>
        </el-descriptions>
        
        <h4 style="margin-top: 20px;">端口配置</h4>
        <el-table :data="selectedService.ports" style="width: 100%; margin-top: 10px;">
          <el-table-column prop="name" label="名称" width="120" />
          <el-table-column prop="port" label="端口" width="100" />
          <el-table-column prop="targetPort" label="目标端口" width="120" />
          <el-table-column prop="protocol" label="协议" width="100" />
          <el-table-column prop="nodePort" label="NodePort" width="120" />
        </el-table>
        
        <h4 style="margin-top: 20px;">标签选择器</h4>
        <div class="selector-detail" v-if="Object.keys(selectedService.selector).length > 0">
          <el-tag v-for="(value, key) in selectedService.selector" :key="key" type="info" size="small" style="margin-right: 8px; margin-bottom: 8px;">
            {{ key }}={{ value }}
          </el-tag>
        </div>
        <div v-else class="empty-text">无标签选择器</div>
        
        <h4 style="margin-top: 20px;">注解</h4>
        <div class="annotations-detail" v-if="Object.keys(selectedService.annotations).length > 0">
          <el-tag v-for="(value, key) in selectedService.annotations" :key="key" type="success" size="small" style="margin-right: 8px; margin-bottom: 8px;">
            {{ key }}={{ value }}
          </el-tag>
        </div>
        <div v-else class="empty-text">无注解</div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="serviceDetailsVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { Refresh, View, Delete } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'

// 加载状态
const loading = ref(false)

// 命名空间加载状态
const namespacesLoading = ref(false)

// 集群数据类型定义
interface Cluster {
  id: number
  name: string
  apiURL: string
  token: string
  kubeconfig: string
  description: string
  status: string
  version: string
  isActive: boolean
}

// 集群列表
const clusters = ref<Cluster[]>([])

// 命名空间列表
const namespaces = ref<string[]>([])

// 选择的集群和命名空间
const selectedCluster = ref<number | undefined>(undefined)
const selectedNamespace = ref<string>('')

// Service数据类型定义
interface ServicePort {
  name?: string
  port: number
  targetPort: number
  protocol: string
  nodePort?: number
}

interface Service {
  name: string
  namespace: string
  cluster: string
  type: string
  clusterIP: string
  externalIP?: string
  externalName?: string
  ports: ServicePort[]
  selector: Record<string, string>
  sessionAffinity?: string
  annotations: Record<string, string>
  createdAt: string
}

// Service数据
const services = ref<Service[]>([])

// 对话框相关
const serviceDetailsVisible = ref(false)
const selectedService = ref<Service | null>(null)

// 组件挂载时加载数据
onMounted(() => {
  loadClusters()
})

// 监听集群选择变化，获取对应的命名空间列表
watch(selectedCluster, (newClusterId) => {
  if (newClusterId) {
    loadNamespaces(newClusterId)
  } else {
    // 重置命名空间列表
    namespaces.value = []
    selectedNamespace.value = ''
  }
})

// 监听命名空间选择变化，自动刷新Service列表
watch(selectedNamespace, () => {
  if (selectedCluster.value) {
    loadServices()
  }
})

// 加载集群列表
const loadClusters = async () => {
  try {
    // 调用真实的API获取集群列表，axios会自动添加baseURL
    const response = await axios.get('/clusters')
    clusters.value = response.data
    // 不自动选择集群，让用户手动选择
  } catch (error) {
    console.error('加载集群列表失败:', error)
    // 不显示错误消息给用户，只在控制台打印
    clusters.value = []
  }
}

// 加载命名空间列表
const loadNamespaces = async (clusterId: number) => {
  try {
    namespacesLoading.value = true
    // 调用API获取指定集群的命名空间列表，axios会自动添加baseURL
    // 添加分页参数，获取所有命名空间
    const response = await axios.get(`/clusters/${clusterId}/namespaces`, {
      params: {
        page: 1,
        pageSize: 1000 // 设置一个足够大的值，确保获取所有命名空间
      }
    })
    // 检查后端返回的数据格式，如果是分页对象则使用data字段
    const namespacesData = response.data.data || response.data
    // 只提取命名空间名称，确保是字符串数组
    namespaces.value = namespacesData.map((ns: any) => ns.name || '')
    // 过滤掉空字符串
    namespaces.value = namespaces.value.filter((ns: string) => ns !== '')
    // 不自动选择命名空间，让用户手动选择
    selectedNamespace.value = ''
    namespacesLoading.value = false
  } catch (error) {
    console.error('加载命名空间列表失败:', error)
    // 不显示错误消息给用户，只在控制台打印
    namespaces.value = []
    selectedNamespace.value = ''
    namespacesLoading.value = false
  }
}

// 取消令牌，用于取消之前的请求
let cancelTokenSource: any = null

// 加载Service列表
const loadServices = async () => {
  if (!selectedCluster.value) {
    ElMessage.warning('请先选择一个集群')
    return
  }
  
  // 取消之前的请求
  if (cancelTokenSource) {
    cancelTokenSource.cancel('Operation canceled by the user.')
  }
  
  // 创建新的取消令牌
  cancelTokenSource = axios.CancelToken.source()
  
  loading.value = true
  try {
    // 调用API获取Service列表，包含命名空间参数，axios会自动添加baseURL
    const response = await axios.get(`/clusters/${selectedCluster.value}/services`, {
      params: {
        namespace: selectedNamespace.value || '',
        _t: Date.now() // 添加时间戳，防止浏览器缓存
      },
      headers: {
        'Cache-Control': 'no-cache' // 禁用缓存
      },
      cancelToken: cancelTokenSource.token
    })
    // 为每个Service添加集群名称
    const clusterName = clusters.value.find(c => c.id === selectedCluster.value)?.name || ''
    const newServices = response.data.map((service: any) => ({
      ...service,
      cluster: clusterName
    }))
    
    // 先更新数据，然后设置loading为false，确保数据正确显示
    services.value = newServices
    loading.value = false
  } catch (error) {
    // 如果是取消请求的错误，不显示错误信息
    if (axios.isCancel(error)) {
      console.log('Request canceled:', error.message)
      return
    }
    console.error('加载Service列表失败:', error)
    loading.value = false
    ElMessage.error('加载Service列表失败')
    // 确保在错误情况下，列表为空，显示正确的空状态
    services.value = []
  }
}

// 刷新Service列表
const refreshServices = () => {
  loadServices()
}

// 获取类型标签类型
const getTypeTagType = (type: string) => {
  switch (type) {
    case 'ClusterIP':
      return 'primary'
    case 'NodePort':
      return 'success'
    case 'LoadBalancer':
      return 'warning'
    case 'ExternalName':
      return 'info'
    default:
      return 'info'
  }
}

// 格式化相对时间（AGE）
const formatAge = (createdAt: string) => {
  if (!createdAt) return ''
  const now = new Date()
  const created = new Date(createdAt)
  const diffInSeconds = Math.floor((now.getTime() - created.getTime()) / 1000)
  
  if (diffInSeconds < 60) {
    return `${diffInSeconds}s`
  } else if (diffInSeconds < 3600) {
    return `${Math.floor(diffInSeconds / 60)}m`
  } else if (diffInSeconds < 86400) {
    return `${Math.floor(diffInSeconds / 3600)}h`
  } else if (diffInSeconds < 2592000) {
    return `${Math.floor(diffInSeconds / 86400)}d`
  } else if (diffInSeconds < 31536000) {
    return `${Math.floor(diffInSeconds / 2592000)}M`
  } else {
    return `${Math.floor(diffInSeconds / 31536000)}y`
  }
}

// 查看Service详情
const viewServiceDetails = (service: Service) => {
  selectedService.value = service
  serviceDetailsVisible.value = true
}

// 删除Service
const deleteService = async (service: Service) => {
  try {
    // 这里应该调用真实的API删除Service
    // 暂时使用模拟数据
    const index = services.value.findIndex((s) => s.namespace === service.namespace && s.name === service.name)
    if (index !== -1) {
      services.value.splice(index, 1)
    }
    ElMessage.success('删除Service成功')
  } catch (error) {
    console.error('删除Service失败:', error)
    ElMessage.error('删除Service失败')
  }
}
</script>

<style scoped>
.services-view {
  width: 100%;
  padding: 0;
}

.page-header {
  margin-bottom: 20px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-content h2 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: #333;
}

.header-actions {
  display: flex;
  align-items: center;
}

.services-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.port-item {
  margin-bottom: 4px;
  font-size: 13px;
}

.selector-item {
  display: flex;
  gap: 4px;
  margin-bottom: 4px;
  font-size: 13px;
}

.selector-key {
  font-weight: 500;
  color: #606266;
}

.selector-value {
  color: #303133;
}

.service-details {
  padding: 10px 0;
}

.selector-detail {
  margin-top: 10px;
}

.annotations-detail {
  margin-top: 10px;
}

.empty-text {
  margin-top: 10px;
  color: #909399;
  font-style: italic;
}

.no-cluster-selected {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 300px;
  background-color: #fafafa;
  border-radius: 8px;
  margin-top: 20px;
}
</style>