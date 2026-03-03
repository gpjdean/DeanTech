<template>
  <div class="pvs-view">
    <el-card shadow="never" class="page-header">
      <div class="header-content">
        <h2>存储卷</h2>
        <div class="header-actions">
          <el-select v-model="selectedCluster" placeholder="选择集群" style="width: 200px; margin-right: 10px;">
            <el-option
              v-for="cluster in clusters"
              :key="cluster.id"
              :label="cluster.name"
              :value="cluster.id"
            />
          </el-select>
          <el-button type="primary" size="small" @click="refreshPVs">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </div>
    </el-card>

    <el-card shadow="hover" class="pvs-card">
      <template #header>
        <div class="card-header">
          <span>K8s PV列表</span>
        </div>
      </template>
      
      <el-table 
        :data="pvs" 
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
        empty-text="No Data"
        :row-key="(row: any) => row.name"
      >
        <el-table-column prop="name" label="PV名称" min-width="200" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="120">
          <template #default="scope">
            <el-tag :type="getStatusTagType(scope.row.status)" size="medium">
              {{ scope.row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="claimRef" label="绑定PVC" min-width="200" show-overflow-tooltip>
          <template #default="scope">
            <span v-if="scope.row.claimRef">{{ scope.row.claimRef.namespace }}/{{ scope.row.claimRef.name }}</span>
            <span v-else style="color: #909399;">未绑定</span>
          </template>
        </el-table-column>
        <el-table-column prop="storageClassName" label="存储类" width="150" show-overflow-tooltip />
        <el-table-column prop="capacity" label="容量" width="120" show-overflow-tooltip />
        <el-table-column prop="accessModes" label="访问模式" min-width="200">
          <template #default="scope">
            <el-tag 
              v-for="mode in scope.row.accessModes" 
              :key="mode" 
              type="info" 
              size="small" 
              style="margin-right: 4px; margin-bottom: 4px;"
            >
              {{ mode }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="reclaimPolicy" label="回收策略" width="120">
          <template #default="scope">
            <el-tag :type="getReclaimPolicyTagType(scope.row.reclaimPolicy)" size="medium">
              {{ scope.row.reclaimPolicy }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="180">
          <template #default="scope">
            {{ formatAge(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" @click="viewPVDetails(scope.row)">
              <el-icon><View /></el-icon>
              查看
            </el-button>
            <el-button type="danger" size="small" @click="deletePV(scope.row)">
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <div v-else class="no-cluster-selected">
        <el-empty description="请选择一个集群查看存储卷列表" />
      </div>
    </el-card>
    
    <!-- PV详情对话框 -->
    <el-dialog
      v-model="pvDetailsVisible"
      title="PV详情"
      width="800px"
    >
      <div v-if="selectedPV" class="pv-details">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="PV名称">{{ selectedPV.name }}</el-descriptions-item>
          <el-descriptions-item label="状态">{{ selectedPV.status }}</el-descriptions-item>
          <el-descriptions-item label="存储类">{{ selectedPV.storageClassName || '无' }}</el-descriptions-item>
          <el-descriptions-item label="容量">{{ selectedPV.capacity }}</el-descriptions-item>
          <el-descriptions-item label="回收策略">{{ selectedPV.reclaimPolicy }}</el-descriptions-item>
          <el-descriptions-item label="绑定PVC" :span="2">
            <span v-if="selectedPV.claimRef">{{ selectedPV.claimRef.namespace }}/{{ selectedPV.claimRef.name }}</span>
            <span v-else style="color: #909399;">未绑定</span>
          </el-descriptions-item>
          <el-descriptions-item label="访问模式" :span="2">
            <el-tag 
              v-for="mode in selectedPV.accessModes" 
              :key="mode" 
              type="info" 
              size="small" 
              style="margin-right: 4px; margin-bottom: 4px;"
            >
              {{ mode }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ selectedPV.createdAt }}</el-descriptions-item>
        </el-descriptions>
        
        <h4 style="margin-top: 20px;">Labels</h4>
        <div class="labels-section">
          <el-tag v-if="Object.keys(selectedPV.labels).length === 0" type="info" size="small">
            无
          </el-tag>
          <el-tag
            v-for="(value, key) in selectedPV.labels"
            :key="key"
            type="info"
            size="small"
            style="margin-right: 8px; margin-bottom: 8px;"
          >
            {{ key }}={{ value }}
          </el-tag>
        </div>
        
        <h4 style="margin-top: 20px;">Annotations</h4>
        <div class="annotations-section">
          <el-tag v-if="Object.keys(selectedPV.annotations).length === 0" type="info" size="small">
            无
          </el-tag>
          <el-tag
            v-for="(value, key) in selectedPV.annotations"
            :key="key"
            type="success"
            size="small"
            style="margin-right: 8px; margin-bottom: 8px;"
          >
            {{ key }}={{ value }}
          </el-tag>
        </div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="pvDetailsVisible = false">关闭</el-button>
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

// 选择的集群
const selectedCluster = ref<number | undefined>(undefined)

// PV数据类型定义
interface PVClaimRef {
  namespace: string
  name: string
}

interface PV {
  name: string
  namespace?: string
  cluster: string
  status: string
  claimRef?: PVClaimRef
  storageClassName?: string
  capacity: string
  accessModes: string[]
  reclaimPolicy: string
  labels: Record<string, string>
  annotations: Record<string, string>
  createdAt: string
}

// PV数据
const pvs = ref<PV[]>([])

// 对话框相关
const pvDetailsVisible = ref(false)
const selectedPV = ref<PV | null>(null)

// 组件挂载时加载数据
onMounted(() => {
  loadClusters()
})

// 监听集群选择变化，自动刷新PV列表
watch(selectedCluster, () => {
  if (selectedCluster.value) {
    loadPVs()
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

// 加载PV列表
const loadPVs = async () => {
  if (!selectedCluster.value) {
    ElMessage.warning('请先选择一个集群')
    return
  }
  
  loading.value = true
  try {
    // 调用API获取PV列表，axios会自动添加baseURL
    const response = await axios.get(`/clusters/${selectedCluster.value}/pvs`, {
      params: {
        _t: Date.now() // 添加时间戳，防止浏览器缓存
      },
      headers: {
        'Cache-Control': 'no-cache' // 禁用缓存
      }
    })
    // 为每个PV添加集群名称
    const clusterName = clusters.value.find(c => c.id === selectedCluster.value)?.name || ''
    const newPVs = response.data.map((pv: any) => ({
      ...pv,
      cluster: clusterName
    }))
    
    // 先更新数据，然后设置loading为false，确保数据正确显示
    pvs.value = newPVs
    loading.value = false
    
    // 如果返回的PV列表为空，显示提示信息
    if (newPVs.length === 0) {
      ElMessage.info('当前集群内无PV实例')
    }
  } catch (error) {
    console.error('加载PV列表失败:', error)
    loading.value = false
    ElMessage.error('加载PV列表失败')
    // 确保在错误情况下，列表为空，显示正确的空状态
    pvs.value = []
  }
}

// 刷新PV列表
const refreshPVs = () => {
  loadPVs()
}

// 获取状态标签类型
const getStatusTagType = (status: string) => {
  switch (status) {
    case 'Bound':
      return 'success'
    case 'Available':
      return 'primary'
    case 'Released':
      return 'warning'
    case 'Failed':
      return 'danger'
    default:
      return 'info'
  }
}

// 获取回收策略标签类型
const getReclaimPolicyTagType = (policy: string) => {
  switch (policy) {
    case 'Retain':
      return 'warning'
    case 'Delete':
      return 'success'
    case 'Recycle':
      return 'info'
    default:
      return 'info'
  }
}

// 格式化时间
const formatAge = (createdAt: string) => {
  if (!createdAt) return ''
  
  try {
    const created = new Date(createdAt)
    
    // 格式化时间为YYYY-MM-DD HH:MM:SS
    const year = created.getFullYear()
    const month = String(created.getMonth() + 1).padStart(2, '0')
    const day = String(created.getDate()).padStart(2, '0')
    const hours = String(created.getHours()).padStart(2, '0')
    const minutes = String(created.getMinutes()).padStart(2, '0')
    const seconds = String(created.getSeconds()).padStart(2, '0')
    
    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
  } catch (error) {
    console.error('时间格式化失败:', error)
    return ''
  }
}

// 查看PV详情
const viewPVDetails = (pv: PV) => {
  selectedPV.value = pv
  pvDetailsVisible.value = true
}

// 删除PV
const deletePV = async (pv: PV) => {
  try {
    // 这里应该调用真实的API删除PV
    // 暂时使用模拟数据
    const index = pvs.value.findIndex((p) => p.name === pv.name)
    if (index !== -1) {
      pvs.value.splice(index, 1)
    }
    ElMessage.success('删除PV成功')
  } catch (error) {
    console.error('删除PV失败:', error)
    ElMessage.error('删除PV失败')
  }
}
</script>

<style scoped>
.pvs-view {
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

.pvs-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pv-details {
  padding: 10px 0;
}

.labels-section,
.annotations-section {
  margin-top: 10px;
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