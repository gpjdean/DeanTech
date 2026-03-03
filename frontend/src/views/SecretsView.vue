<template>
  <div class="secrets-view">
    <el-card shadow="never" class="page-header">
      <div class="header-content">
        <h2>密钥</h2>
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
          <el-button type="primary" size="small" @click="refreshSecrets">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </div>
    </el-card>

    <el-card shadow="hover" class="secrets-card">
      <template #header>
        <div class="card-header">
          <span>K8s Secret列表</span>
          <el-tooltip content="Secret包含敏感信息，请注意保密">
            <el-icon class="warning-icon"><Warning /></el-icon>
          </el-tooltip>
        </div>
      </template>
      
      <el-table 
        :data="secrets" 
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
        <el-table-column prop="name" label="Secret名称" min-width="200" show-overflow-tooltip />
        <el-table-column prop="namespace" label="命名空间" width="150" show-overflow-tooltip />
        <el-table-column prop="cluster" label="集群" width="150" show-overflow-tooltip />
        <el-table-column prop="type" label="类型" width="150">
          <template #default="scope">
            <el-tag :type="getTypeTagType(scope.row.type)" size="medium">
              {{ scope.row.type }}
            </el-tag>
          </template>
        </el-table-column>
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
            <el-button type="primary" size="small" @click="viewSecretDetails(scope.row)">
              <el-icon><View /></el-icon>
              查看
            </el-button>
            <el-button type="danger" size="small" @click="deleteSecret(scope.row)">
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <div v-else class="no-cluster-selected">
        <el-empty description="请选择一个集群查看Secret列表" />
      </div>
    </el-card>
    
    <!-- Secret详情对话框 -->
    <el-dialog
      v-model="secretDetailsVisible"
      title="Secret详情"
      width="800px"
    >
      <div v-if="selectedSecret" class="secret-details">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="Secret名称">{{ selectedSecret.name }}</el-descriptions-item>
          <el-descriptions-item label="命名空间">{{ selectedSecret.namespace }}</el-descriptions-item>
          <el-descriptions-item label="类型">{{ selectedSecret.type }}</el-descriptions-item>
          <el-descriptions-item label="数据项数量">{{ Object.keys(selectedSecret.data).length }}</el-descriptions-item>
          <el-descriptions-item label="二进制数据项">{{ Object.keys(selectedSecret.binaryData).length }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ selectedSecret.createdAt }}</el-descriptions-item>
        </el-descriptions>
        
        <el-alert
          title="安全提示"
          type="warning"
          :closable="false"
          style="margin-top: 20px;"
        >
          此Secret包含敏感信息，请确保只有授权人员可以访问。
        </el-alert>
        
        <h4 style="margin-top: 20px;">数据内容</h4>
        <div v-if="Object.keys(selectedSecret.data).length > 0">
          <el-tabs type="border-card">
            <el-tab-pane
              v-for="(value, key) in selectedSecret.data"
              :key="key"
              :label="key"
            >
              <div class="secret-data-content">
                <div class="secret-data-header">
                  <span>Base64编码数据</span>
                  <el-button type="primary" size="small" @click="toggleSecretDataVisibility(key)">
                    {{ secretDataVisibility[key] ? '隐藏' : '显示' }}明文
                  </el-button>
                </div>
                <pre class="encoded-data">{{ value }}</pre>
                <pre v-if="secretDataVisibility[key]" class="decoded-data">{{ decodeBase64(value) }}</pre>
              </div>
            </el-tab-pane>
          </el-tabs>
        </div>
        <div v-else class="empty-text">无数据内容</div>
        
        <h4 style="margin-top: 20px;">二进制数据</h4>
        <div v-if="Object.keys(selectedSecret.binaryData).length > 0">
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
          <el-button @click="secretDetailsVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { Refresh, View, Delete, Download, Warning } from '@element-plus/icons-vue'
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

// Secret数据类型定义
interface Secret {
  name: string
  namespace: string
  cluster: string
  type: string
  data: Record<string, string>
  binaryData: Record<string, string>
  annotations: Record<string, string>
  createdAt: string
}

// Secret数据
const secrets = ref<Secret[]>([])

// 对话框相关
const secretDetailsVisible = ref(false)
const selectedSecret = ref<Secret | null>(null)
const secretDataVisibility = ref<Record<string, boolean>>({})

// 二进制数据列表
const binaryDataList = computed(() => {
  if (!selectedSecret.value) return []
  return Object.entries(selectedSecret.value.binaryData).map(([key, value]) => ({
    key,
    size: `${Math.ceil(value.length / 1024)} KB`,
    value
  }))
})

// 监听Secret详情对话框关闭事件，重置数据可见性
watch(secretDetailsVisible, (newVal) => {
  if (!newVal) {
    secretDataVisibility.value = {}
  }
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

// 监听命名空间选择变化，自动刷新Secret列表
watch(selectedNamespace, () => {
  if (selectedCluster.value) {
    loadSecrets()
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

// 加载Secret列表
const loadSecrets = async () => {
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
    // 调用API获取Secret列表，axios会自动添加baseURL
    // 只有当选择了命名空间时才传递namespace参数
    const params: any = {
      _t: Date.now() // 添加时间戳，防止浏览器缓存
    }
    if (selectedNamespace.value) {
      params.namespace = selectedNamespace.value
    }
    const response = await axios.get(`/clusters/${selectedCluster.value}/secrets`, {
      params,
      headers: {
        'Cache-Control': 'no-cache' // 禁用缓存
      },
      cancelToken: cancelTokenSource.token
    })
    // 为每个Secret添加集群名称
    const clusterName = clusters.value.find(c => c.id === selectedCluster.value)?.name || ''
    const newSecrets = response.data.map((secret: any) => ({
      ...secret,
      cluster: clusterName
    }))
    
    // 先更新数据，然后设置loading为false，确保数据正确显示
    secrets.value = newSecrets
    loading.value = false
  } catch (error) {
    // 如果是取消请求的错误，不显示错误信息
    if (axios.isCancel(error)) {
      console.log('Request canceled:', error.message)
      return
    }
    console.error('加载Secret列表失败:', error)
    loading.value = false
    ElMessage.error('加载Secret列表失败')
    // 确保在错误情况下，列表为空，显示正确的空状态
    secrets.value = []
  }
}

// 刷新Secret列表
const refreshSecrets = () => {
  loadSecrets()
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

// 获取类型标签类型
const getTypeTagType = (type: string) => {
  switch (type) {
    case 'Opaque':
      return 'primary'
    case 'kubernetes.io/service-account-token':
      return 'success'
    case 'kubernetes.io/tls':
      return 'warning'
    case 'kubernetes.io/dockerconfigjson':
      return 'info'
    default:
      return 'info'
  }
}

// 查看Secret详情
const viewSecretDetails = (secret: Secret) => {
  selectedSecret.value = secret
  secretDetailsVisible.value = true
}

// 删除Secret
const deleteSecret = async (secret: Secret) => {
  try {
    // 这里应该调用真实的API删除Secret
    // 暂时使用模拟数据
    const index = secrets.value.findIndex((s) => s.namespace === secret.namespace && s.name === secret.name)
    if (index !== -1) {
      secrets.value.splice(index, 1)
    }
    ElMessage.success('删除Secret成功')
  } catch (error) {
    console.error('删除Secret失败:', error)
    ElMessage.error('删除Secret失败')
  }
}

// Base64解码函数
const decodeBase64 = (str: string) => {
  try {
    // 使用浏览器内置的 atob 函数进行解码
    return decodeURIComponent(atob(str).split('').map(function(c) {
      return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)
    }).join(''))
  } catch (error) {
    console.error('Base64解码失败:', error)
    return '解码失败'
  }
}

// 切换Secret数据可见性
const toggleSecretDataVisibility = (key: string) => {
  secretDataVisibility.value[key] = !secretDataVisibility.value[key]
}
</script>

<style scoped>
.secrets-view {
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

.secrets-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.warning-icon {
  color: #e6a23c;
  font-size: 16px;
}

.secret-details {
  padding: 10px 0;
}

.secret-data-content {
  background-color: #f5f7fa;
  padding: 10px;
  border-radius: 4px;
  overflow-x: auto;
}

.secret-data-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.encoded-data {
  background-color: #fff;
  padding: 10px;
  border-radius: 4px;
  margin-bottom: 10px;
  font-family: 'Courier New', Courier, monospace;
  font-size: 13px;
  overflow-x: auto;
}

.decoded-data {
  background-color: #f0f9ff;
  padding: 10px;
  border-radius: 4px;
  border-left: 4px solid #3b82f6;
  font-family: 'Courier New', Courier, monospace;
  font-size: 13px;
  overflow-x: auto;
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