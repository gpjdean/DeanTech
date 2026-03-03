<template>
  <div class="nodes-view">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <h2>节点管理</h2>
        <div class="header-stats">
          <div class="stat-item">
            <span class="stat-value">{{ nodes.length }}</span>
            <span class="stat-label">活跃节点</span>
          </div>
          <div class="stat-item">
            <span class="stat-value">{{ clusters.length }}</span>
            <span class="stat-label">可用集群</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 节点列表卡片 -->
    <el-card shadow="hover" class="nodes-card">
      <template #header>
        <div class="card-header">
          <span>K8s节点列表</span>
          <div class="header-actions">
            <el-select v-model="selectedClusterId" placeholder="选择集群" style="width: 200px; margin-right: 10px;">
              <el-option 
                v-for="cluster in clusters" 
                :key="cluster.id" 
                :label="cluster.name" 
                :value="cluster.id" 
              />
            </el-select>
            <el-button type="primary" size="small" @click="loadNodes">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
          </div>
        </div>
      </template>
      
      <div v-if="selectedClusterId" class="nodes-container">
        <!-- 添加表格滚动容器 -->
        <div class="table-scroll-container">
          <el-table 
            :data="nodes" 
            style="width: 100%" 
            stripe
            v-loading="loading"
            border
          >
            <el-table-column prop="name" label="节点名称" min-width="180" />
            <el-table-column prop="ip" label="IP地址" min-width="150" />
            <el-table-column prop="role" label="节点角色" min-width="120">
              <template #default="scope">
                <el-tag :type="scope.row.role === 'master' ? 'primary' : 'success'">
                  {{ scope.row.role }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" min-width="120">
              <template #default="scope">
                <el-tag :type="scope.row.status === 'Ready' ? 'success' : 'danger'">
                  {{ scope.row.status }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="os" label="操作系统" min-width="150" show-overflow-tooltip />
            <el-table-column prop="kernelVersion" label="内核版本" min-width="150" show-overflow-tooltip />
            <el-table-column prop="containerRuntime" label="容器运行时" min-width="150" show-overflow-tooltip />
            <el-table-column prop="cpu" label="CPU核心" min-width="100" />
            <el-table-column prop="memory" label="内存(MB)" min-width="120" />
            <el-table-column prop="pods" label="限制Pod数量" min-width="100" />
            <el-table-column label="操作" min-width="120" fixed="right">
              <template #default="scope">
                <el-button type="primary" size="small" @click="showNodeDetails(scope.row)">
                  详情
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
        
        <!-- 节点统计信息 -->
        <div class="nodes-stats" v-if="nodes.length > 0">
          <div class="stat-card">
            <div class="stat-title">节点资源概览</div>
            <div class="stat-grid">
              <div class="stat-item">
                <span class="stat-value">{{ totalCpu }}</span>
                <span class="stat-label">总CPU核心</span>
              </div>
              <div class="stat-item">
                <span class="stat-value">{{ totalMemory }}</span>
                <span class="stat-label">总内存(MB)</span>
              </div>
              <div class="stat-item">
                <span class="stat-value">{{ totalPods }}</span>
                <span class="stat-label">总Pod数量</span>
              </div>
              <div class="stat-item">
                <span class="stat-value">{{ readyNodes }}</span>
                <span class="stat-label">就绪节点</span>
              </div>
              <div class="stat-item">
                <span class="stat-value">{{ notReadyNodes }}</span>
                <span class="stat-label">未就绪节点</span>
              </div>
              <div class="stat-item">
                <span class="stat-value">{{ masterNodes }}</span>
                <span class="stat-label">Master节点</span>
              </div>
              <div class="stat-item">
                <span class="stat-value">{{ workerNodes }}</span>
                <span class="stat-label">Worker节点</span>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <div v-else class="no-cluster-selected">
        <el-empty description="请选择一个集群查看节点列表" />
      </div>
    </el-card>

    <!-- 节点详情对话框 -->
    <el-dialog
      v-model="nodeDetailsVisible"
      title="节点详情"
      width="90%"
      :max-width="'900px'"
      :close-on-click-modal="false"
    >
      <div v-if="selectedNode" class="node-details">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="节点名称">{{ selectedNode.name }}</el-descriptions-item>
          <el-descriptions-item label="IP地址">{{ selectedNode.ip }}</el-descriptions-item>
          <el-descriptions-item label="节点角色">{{ selectedNode.role }}</el-descriptions-item>
          <el-descriptions-item label="节点状态">{{ selectedNode.status }}</el-descriptions-item>
          <el-descriptions-item label="操作系统">{{ selectedNode.os }}</el-descriptions-item>
          <el-descriptions-item label="内核版本">{{ selectedNode.kernelVersion }}</el-descriptions-item>
          <el-descriptions-item label="容器运行时">{{ selectedNode.containerRuntime }}</el-descriptions-item>
          <el-descriptions-item label="CPU核心">{{ selectedNode.cpu }}</el-descriptions-item>
          <el-descriptions-item label="内存(MB)">{{ selectedNode.memory }}</el-descriptions-item>
          <el-descriptions-item label="限制Pod数量">{{ selectedNode.pods }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
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

// 节点数据类型定义
interface Node {
  id: number
  clusterID: number
  name: string
  ip: string
  role: string
  status: string
  os: string
  kernelVersion: string
  containerRuntime: string
  cpu: number
  memory: number
  pods: number
}

// 集群数据
const clusters = ref<Cluster[]>([])
const selectedClusterId = ref<number | undefined>(undefined)

// 节点数据
const nodes = ref<Node[]>([])

// 节点详情对话框状态
const nodeDetailsVisible = ref(false)
const selectedNode = ref<Node | null>(null)

// 组件挂载时加载数据
onMounted(() => {
  loadClusters()
})

// 监听集群选择变化，自动加载节点数据
watch(selectedClusterId, (newClusterId) => {
  if (newClusterId) {
    loadNodes()
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

// 加载节点列表
const loadNodes = async () => {
  if (!selectedClusterId.value) {
    ElMessage.warning('请先选择一个集群')
    return
  }
  
  loading.value = true
  try {
    // 先清空列表，确保数据实时性
    nodes.value = []
    
    // 调用真实的API获取节点列表，axios会自动添加baseURL
    const response = await axios.get(`/clusters/${selectedClusterId.value}/nodes`, {
      params: {
        _t: Date.now() // 添加时间戳，防止浏览器缓存
      },
      headers: {
        'Cache-Control': 'no-cache' // 禁用缓存
      }
    })
    nodes.value = response.data
    loading.value = false
  } catch (error) {
    console.error('加载节点列表失败:', error)
    loading.value = false
    ElMessage.error('加载节点列表失败')
  }
}

// 显示节点详情
const showNodeDetails = (node: Node) => {
  selectedNode.value = node
  nodeDetailsVisible.value = true
}

// 计算属性：节点统计信息
const totalCpu = computed(() => {
  return nodes.value.reduce((sum, node) => sum + node.cpu, 0)
})

const totalMemory = computed(() => {
  return nodes.value.reduce((sum, node) => sum + node.memory, 0)
})

const totalPods = computed(() => {
  return nodes.value.reduce((sum, node) => sum + node.pods, 0)
})

const readyNodes = computed(() => {
  return nodes.value.filter(node => node.status === 'Ready').length
})

const notReadyNodes = computed(() => {
  return nodes.value.filter(node => node.status !== 'Ready').length
})

const masterNodes = computed(() => {
  return nodes.value.filter(node => node.role === 'master').length
})

const workerNodes = computed(() => {
  return nodes.value.filter(node => node.role === 'worker').length
})
</script>

<style scoped>
.nodes-view {
  width: 100%;
  padding: 0;
  background-color: #f5f7fa;
  min-height: calc(100vh - 120px);
}

/* 页面头部样式 */
.page-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 40px 0;
  margin-bottom: 30px;
  border-radius: 16px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
  margin: 0 20px 30px 20px;
}

.header-content {
  display: flex;
  align-items: center;
  gap: 20px;
  margin: 0 auto;
  padding: 0 24px;
  justify-content: space-between;
}

.header-content h2 {
  margin: 0 0 8px 0;
  font-size: 32px;
  font-weight: 700;
  color: white;
  letter-spacing: -0.5px;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

/* 头部统计信息 */
.header-stats {
  display: flex;
  gap: 24px;
}

.stat-item {
  text-align: center;
  background: rgba(255, 255, 255, 0.15);
  padding: 16px 24px;
  border-radius: 12px;
  backdrop-filter: blur(10px);
  min-width: 100px;
  transition: all 0.3s ease;
}

.stat-item:hover {
  background: rgba(255, 255, 255, 0.25);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.stat-value {
  display: block;
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 4px;
}

.stat-label {
  display: block;
  font-size: 14px;
  opacity: 0.9;
}

/* 节点卡片样式 */
.nodes-card {
  margin: 0 20px 30px 20px;
  border-radius: 16px;
  background: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  border: 1px solid #e4e7ed;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
}

.nodes-card:hover {
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.15);
  border-color: #667eea;
  transform: translateY(-2px);
}

/* 卡片头部样式 */
.card-header {
  display: flex;
  justify-content: flex-start;
  align-items: center;
  gap: 12px;
  margin-bottom: 0;
  padding: 20px 24px;
  border-bottom: 1px solid #f0f2f5;
  border-radius: 16px 16px 0 0;
  background: linear-gradient(135deg, #fafafa 0%, #f0f2f5 100%);
}

.card-header span {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

/* 卡片头部右侧操作区域 */
.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-left: auto;
}

/* 节点容器 */
.nodes-container {
  padding: 24px;
  background-color: #fafafa;
}

/* 表格滚动容器 */
.table-scroll-container {
  overflow-x: auto;
  margin-bottom: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  background: white;
}

/* 节点统计信息 */
.nodes-stats {
  margin-top: 24px;
}

.stat-card {
  background: white;
  padding: 20px;
  border-radius: 12px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  border: 1px solid #e4e7ed;
}

.stat-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 16px;
}

.stat-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 16px;
}

.nodes-stats .stat-item {
  background: #f0f2f5;
  color: #303133;
  padding: 16px;
  border-radius: 12px;
  text-align: center;
  transition: all 0.3s ease;
  backdrop-filter: none;
}

.nodes-stats .stat-item:hover {
  background: linear-gradient(135deg, #ecf5ff 0%, #f0f9ff 100%);
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(64, 158, 255, 0.2);
  border: 1px solid #d9ecff;
}

.nodes-stats .stat-value {
  font-size: 20px;
  font-weight: 700;
  color: #303133;
}

.nodes-stats .stat-label {
  font-size: 13px;
  color: #606266;
  margin-top: 4px;
}

/* 表格样式 */
.nodes-container :deep(.el-table) {
  width: 100%;
  min-width: 1000px;
  border-radius: 12px;
  overflow: hidden;
  background: white;
  box-shadow: none;
}

.nodes-container :deep(.el-table__header-wrapper) {
  background: #fafafa;
  border-bottom: 1px solid #e4e7ed;
  border-radius: 12px 12px 0 0;
}

.nodes-container :deep(.el-table__header-wrapper th) {
  font-weight: 600;
  color: #303133;
  background: #fafafa;
  border-bottom: 1px solid #e4e7ed;
}

.nodes-container :deep(.el-table__body-wrapper) {
  background: white;
}

.nodes-container :deep(.el-table__body tr:hover > td) {
  background-color: #ecf5ff;
  border-radius: 6px;
}

/* 优化标签样式 */
.nodes-container :deep(.el-tag) {
  border-radius: 12px;
  padding: 0 12px;
  font-size: 12px;
  font-weight: 500;
}

/* 优化按钮样式 */
.nodes-container :deep(.el-button--primary) {
  border-radius: 8px;
  font-weight: 500;
  padding: 8px 16px;
  transition: all 0.3s ease;
}

.nodes-container :deep(.el-button--primary:hover) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
}

/* 无集群选择样式 */
.no-cluster-selected {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 300px;
  background-color: #fafafa;
  border-radius: 8px;
  margin: 20px;
}

/* 节点详情对话框样式 */
.node-details {
  padding: 20px;
}

/* 响应式设计 */
@media (max-width: 1400px) {
  .header-content {
    padding: 0 20px;
  }
  
  .nodes-card {
    margin: 0 20px 30px 20px;
  }
}

@media (max-width: 1200px) {
  .header-content {
    flex-direction: column;
    text-align: center;
    gap: 20px;
  }
  
  .header-stats {
    width: 100%;
    justify-content: center;
  }
  
  .stat-grid {
    grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
    gap: 14px;
  }
}

@media (max-width: 1024px) {
  .stat-grid {
    grid-template-columns: repeat(auto-fit, minmax(130px, 1fr));
    gap: 12px;
  }
  
  .nodes-card {
    margin: 0 16px 24px 16px;
  }
  
  .nodes-container {
    padding: 20px;
  }
}

@media (max-width: 768px) {
  .page-header {
    padding: 30px 0;
  }
  
  .header-content h2 {
    font-size: 24px;
  }
  
  .header-stats {
    flex-direction: column;
    gap: 12px;
    width: 100%;
  }
  
  .stat-item {
    width: 100%;
  }
  
  .card-header {
    flex-direction: column;
    align-items: stretch;
    gap: 16px;
  }
  
  .header-actions {
    margin-left: 0;
    justify-content: flex-start;
    flex-wrap: wrap;
  }
  
  .nodes-card {
    margin: 0 16px 24px 16px;
  }
  
  .nodes-container {
    padding: 16px;
  }
  
  .stat-grid {
    grid-template-columns: 1fr 1fr;
  }
  
  /* 调整表格列显示，在手机端优化显示效果 */
  .nodes-container :deep(.el-table) {
    min-width: 800px;
  }
}

@media (max-width: 480px) {
  .page-header {
    padding: 20px 0;
  }
  
  .header-content {
    padding: 0 12px;
  }
  
  .header-content h2 {
    font-size: 20px;
  }
  
  .stat-item {
    padding: 12px 16px;
  }
  
  .stat-value {
    font-size: 20px;
  }
  
  .stat-label {
    font-size: 13px;
  }
  
  .nodes-card {
    margin: 0 8px 16px 8px;
  }
  
  .nodes-container {
    padding: 12px;
  }
  
  .card-header {
    padding: 16px;
  }
  
  .stat-title {
    font-size: 14px;
    margin-bottom: 12px;
  }
  
  .stat-grid {
    gap: 10px;
  }
  
  .nodes-stats .stat-item {
    padding: 12px;
  }
  
  .nodes-stats .stat-value {
    font-size: 18px;
  }
  
  .nodes-stats .stat-label {
    font-size: 12px;
  }
  
  /* 进一步优化表格在极小屏幕上的显示 */
  .nodes-container :deep(.el-table) {
    min-width: 700px;
  }
  
  /* 对话框优化 */
  .node-details {
    padding: 12px;
  }
  
  :deep(.el-descriptions) {
    :deep(.el-descriptions__item-label) {
      font-size: 13px;
    }
    
    :deep(.el-descriptions__content) {
      font-size: 13px;
    }
  }
}
</style>
