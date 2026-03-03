<template>
  <div class="pvcs-view">
    <el-card shadow="never" class="page-header">
      <div class="header-content">
        <h2>存储卷声明</h2>
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
          <el-button type="primary" size="small" @click="refreshPVCs">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </div>
    </el-card>

    <el-card shadow="hover" class="pvcs-card">
      <template #header>
        <div class="card-header">
          <span>K8s PVC列表</span>
        </div>
      </template>
      
      <el-table 
        :data="pvcs" 
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
      >
        <el-table-column prop="name" label="PVC名称" width="200" />
        <el-table-column prop="namespace" label="命名空间" width="150" />
        <el-table-column prop="status" label="状态" width="120">
          <template #default="scope">
            <el-tag :type="getStatusTagType(scope.row.status)" size="small">
              {{ scope.row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="volumeName" label="卷名称" width="200" />
        <el-table-column prop="capacity" label="容量" width="120" />
        <el-table-column prop="accessModes" label="访问模式" width="200">
          <template #default="scope">
            <el-tag size="small" v-for="mode in scope.row.accessModes" :key="mode" class="mode-tag">
              {{ mode }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="storageClassName" label="存储类" width="180" />
        <el-table-column prop="createdAt" label="创建时间" width="200">
          <template #default="scope">
            {{ formatDate(scope.row.createdAt) }}
          </template>
        </el-table-column>
      </el-table>
      
      <div v-else class="no-cluster-selected">
        <el-empty description="请选择一个集群查看存储卷声明列表" />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import axios from 'axios'

// PVC数据类型定义
interface PVC {
  name: string
  namespace: string
  status: string
  volumeName: string
  capacity: string
  accessModes: string[]
  storageClassName: string
  createdAt: string
}

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
  createdAt: string
  connectionType: string
  nodeCount?: string
}

// PVC列表
const pvcs = ref<PVC[]>([])
// 列表加载状态
const loading = ref(false)
// 命名空间列表
const namespaces = ref<string[]>([])
// 命名空间加载状态
const namespacesLoading = ref(false)
// 集群列表
const clusters = ref<Cluster[]>([])
// 选中的集群
const selectedCluster = ref<number | undefined>(undefined)
// 选中的命名空间
const selectedNamespace = ref<string>('')

// 获取状态对应的标签类型
const getStatusTagType = (status: string) => {
  switch (status) {
    case 'Bound':
      return 'success'
    case 'Pending':
      return 'warning'
    case 'Lost':
      return 'danger'
    default:
      return 'info'
  }
}

// 格式化日期函数
const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  
  // 手动格式化日期，确保格式为 YYYY-MM-DD HH:MM:SS
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

// 组件挂载时加载数据
onMounted(() => {
  loadClusters()
})

// 监听集群选择变化，只加载对应的命名空间，不直接加载PVC
watch(selectedCluster, (newClusterId) => {
  if (newClusterId) {
    loadNamespaces(newClusterId)
    // 重置PVC列表，等待命名空间选择后再加载
    pvcs.value = []
    selectedNamespace.value = ''
  } else {
    // 重置命名空间和PVC列表
    namespaces.value = []
    pvcs.value = []
    selectedNamespace.value = ''
  }
})

// 监听命名空间选择变化，只有选择了有效的命名空间后才加载PVC
watch(selectedNamespace, (newNamespace) => {
  if (selectedCluster.value && newNamespace) {
    loadPVCs(selectedCluster.value)
  } else {
    // 重置PVC列表
    pvcs.value = []
  }
})

// 加载集群列表
const loadClusters = async () => {
  try {
    const response = await axios.get('/clusters')
    clusters.value = response.data
  } catch (error) {
    console.error('加载集群列表失败:', error)
    ElMessage.error('加载集群列表失败')
    clusters.value = []
  }
}

// 加载命名空间列表
const loadNamespaces = async (clusterId: number) => {
  if (!clusterId) return
  
  namespacesLoading.value = true
  try {
    // 添加分页参数，获取所有命名空间（page=1，pageSize设为较大值）
    const response = await axios.get(`/clusters/${clusterId}/namespaces`, {
      params: {
        page: 1,
        pageSize: 1000 // 设置一个足够大的值，确保获取所有命名空间
      }
    })
    // 检查后端返回的数据格式，如果是分页对象则使用data字段
    const namespacesData = response.data.data || response.data
    namespaces.value = namespacesData.map((ns: any) => ns.name)
    // 不添加空白选项，要求用户必须选择一个具体的命名空间
  } catch (error) {
    console.error('加载命名空间列表失败:', error)
    ElMessage.error('加载命名空间列表失败')
    namespaces.value = []
  } finally {
    namespacesLoading.value = false
  }
}

// 加载PVC列表
const loadPVCs = async (clusterId: number) => {
  if (!clusterId) return
  
  loading.value = true
  try {
    const params: any = {
      namespace: selectedNamespace.value
    }
    
    const response = await axios.get(`/clusters/${clusterId}/pvcs`, { params })
    pvcs.value = response.data
  } catch (error) {
    console.error('加载PVC列表失败:', error)
    ElMessage.error('加载PVC列表失败')
    pvcs.value = []
  } finally {
    loading.value = false
  }
}

// 刷新PVC列表
const refreshPVCs = () => {
  if (selectedCluster.value) {
    loadPVCs(selectedCluster.value)
  }
}
</script>

<style scoped>
.pvcs-view {
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

.pvcs-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.card-header span {
  font-size: 18px;
  font-weight: 400;
  color: #333;
}

/* 状态标签样式 */
:deep(.el-tag) {
  margin-right: 4px;
}

.mode-tag {
  background-color: #f0f9ff;
  border-color: #91d5ff;
  color: #1890ff;
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

/* 确保表格样式正确 */
:deep(.el-table) {
  width: 100%;
}

:deep(.el-table__header-wrapper th.el-table__header-cell) {
  color: #333;
  font-weight: 600;
}

:deep(.el-table__body-wrapper td.el-table__cell) {
  color: #606266;
}
</style>
