<template>
  <div class="configmaps-view">
    <el-card shadow="never" class="page-header">
      <div class="header-content">
        <h2>配置</h2>
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
          <el-button type="primary" size="small" @click="refreshConfigMaps">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </div>
    </el-card>

    <el-card shadow="hover" class="configmaps-card">
      <template #header>
        <div class="card-header">
          <span>K8s ConfigMap列表</span>
        </div>
      </template>
      
      <el-table 
        :data="configmaps" 
        style="width: 100%; color: #303133; font-size: 14px; line-height: 1.5;"
        stripe
        v-loading="loading"
        v-if="selectedCluster"
        :row-style="{ height: 'auto', backgroundColor: '#ffffff' }"
        :cell-style="{ 
          'white-space': 'normal', 
          'word-break': 'break-all', 
          'min-height': '40px', 
          'padding': '12px', 
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
        empty-text="No Data"
        :row-key="(row: any) => `${row.namespace}-${row.name}`"
      >
        <el-table-column prop="name" label="ConfigMap名称" min-width="200" show-overflow-tooltip />
        <el-table-column prop="namespace" label="命名空间" width="150" show-overflow-tooltip />
        <el-table-column prop="cluster" label="集群" width="150" show-overflow-tooltip />
        <el-table-column prop="data" label="数据项数量" width="120">
          <template #default="scope">
            <el-tag type="info" size="medium">{{ Object.keys(scope.row.data).length }} 项</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="binaryData" label="二进制数据项" width="120">
          <template #default="scope">
            <el-tag type="warning" size="medium">{{ Object.keys(scope.row.binaryData).length }} 项</el-tag>
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
            <el-button type="primary" size="small" @click="viewConfigMapDetails(scope.row)">
              <el-icon><View /></el-icon>
              查看
            </el-button>
            <el-button type="danger" size="small" @click="deleteConfigMap(scope.row)">
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <div v-else class="no-cluster-selected">
        <el-empty description="请选择一个集群查看ConfigMap列表" />
      </div>
    </el-card>
    
    <!-- ConfigMap详情对话框 -->
    <el-dialog
      v-model="configMapDetailsVisible"
      title="ConfigMap详情"
      width="800px"
    >
      <div v-if="selectedConfigMap" class="configmap-details">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="ConfigMap名称">{{ selectedConfigMap.name }}</el-descriptions-item>
          <el-descriptions-item label="命名空间">{{ selectedConfigMap.namespace }}</el-descriptions-item>
          <el-descriptions-item label="数据项数量">{{ Object.keys(selectedConfigMap.data).length }}</el-descriptions-item>
          <el-descriptions-item label="二进制数据项">{{ Object.keys(selectedConfigMap.binaryData).length }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ selectedConfigMap.createdAt }}</el-descriptions-item>
        </el-descriptions>
        
        <h4 style="margin-top: 20px;">数据内容</h4>
        <div v-if="Object.keys(selectedConfigMap.data).length > 0">
          <el-tabs type="border-card">
            <el-tab-pane
              v-for="(value, key) in selectedConfigMap.data"
              :key="key"
              :label="key"
            >
              <div class="configmap-data-content">
                <pre>{{ value }}</pre>
              </div>
            </el-tab-pane>
          </el-tabs>
        </div>
        <div v-else class="empty-text">无数据内容</div>
        
        <h4 style="margin-top: 20px;">二进制数据</h4>
        <div v-if="Object.keys(selectedConfigMap.binaryData).length > 0">
          <el-table :data="binaryDataList" style="width: 100%;">
            <el-table-column prop="key" label="键" width="200" />
            <el-table-column prop="size" label="大小" width="100" />
            <el-table-column prop="actions" label="操作" width="150" fixed="right">
              <template #default>
                <el-button type="primary" size="small">
                  <el-icon><Download /></el-icon>
                  下载
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
        <div v-else class="empty-text">无二进制数据</div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="configMapDetailsVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import { Refresh, View, Delete, Download } from '@element-plus/icons-vue'
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

// ConfigMap数据类型定义
interface ConfigMap {
  name: string
  namespace: string
  cluster: string
  data: Record<string, string>
  binaryData: Record<string, string>
  annotations: Record<string, string>
  createdAt: string
}

// ConfigMap数据
const configmaps = ref<ConfigMap[]>([])

// 对话框相关
const configMapDetailsVisible = ref(false)
const selectedConfigMap = ref<ConfigMap | null>(null)

// 二进制数据列表
const binaryDataList = computed(() => {
  if (!selectedConfigMap.value) return []
  return Object.entries(selectedConfigMap.value.binaryData).map(([key, value]) => ({
    key,
    size: `${Math.ceil(value.length / 1024)} KB`,
    value
  }))
})

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

// 监听命名空间选择变化，自动刷新ConfigMap列表
watch(selectedNamespace, () => {
  if (selectedCluster.value) {
    loadConfigMaps()
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

// 加载ConfigMap列表
const loadConfigMaps = async () => {
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
    // 调用API获取ConfigMap列表，axios会自动添加baseURL
    // 只有当选择了命名空间时才传递namespace参数
    const params: any = {
      _t: Date.now() // 添加时间戳，防止浏览器缓存
    }
    if (selectedNamespace.value) {
      params.namespace = selectedNamespace.value
    }
    const response = await axios.get(`/clusters/${selectedCluster.value}/configmaps`, {
      params,
      headers: {
        'Cache-Control': 'no-cache' // 禁用缓存
      },
      cancelToken: cancelTokenSource.token
    })
    // 为每个ConfigMap添加集群名称
    const clusterName = clusters.value.find(c => c.id === selectedCluster.value)?.name || ''
    const newConfigmaps = response.data.map((configmap: any) => ({
      ...configmap,
      cluster: clusterName
    }))
    
    // 先更新数据，然后设置loading为false，确保数据正确显示
    configmaps.value = newConfigmaps
    loading.value = false
  } catch (error) {
    // 如果是取消请求的错误，不显示错误信息
    if (axios.isCancel(error)) {
      console.log('Request canceled:', error.message)
      return
    }
    console.error('加载ConfigMap列表失败:', error)
    loading.value = false
    ElMessage.error('加载ConfigMap列表失败')
    // 确保在错误情况下，列表为空，显示正确的空状态
    configmaps.value = []
  }
}

// 刷新ConfigMap列表
const refreshConfigMaps = () => {
  loadConfigMaps()
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

// 查看ConfigMap详情
const viewConfigMapDetails = (configMap: ConfigMap) => {
  selectedConfigMap.value = configMap
  configMapDetailsVisible.value = true
}

// 删除ConfigMap
const deleteConfigMap = async (configMap: ConfigMap) => {
  try {
    // 这里应该调用真实的API删除ConfigMap
    // 暂时使用模拟数据
    const index = configmaps.value.findIndex((cm) => cm.namespace === configMap.namespace && cm.name === configMap.name)
    if (index !== -1) {
      configmaps.value.splice(index, 1)
    }
    ElMessage.success('删除ConfigMap成功')
  } catch (error) {
    console.error('删除ConfigMap失败:', error)
    ElMessage.error('删除ConfigMap失败')
  }
}
</script>

<style scoped>
.configmaps-view {
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

.configmaps-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.configmap-details {
  padding: 10px 0;
}

.configmap-data-content {
  background-color: #f5f7fa;
  padding: 10px;
  border-radius: 4px;
  overflow-x: auto;
  max-height: 400px;
  overflow-y: auto;
}

.configmap-data-content pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
  font-family: 'Courier New', Courier, monospace;
  font-size: 13px;
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