<template>
  <div class="cluster-dashboard-view">
    <el-card shadow="never" class="page-header">
      <div class="header-content">
        <h2>仪表盘</h2>
        <div class="header-actions">
          <el-select v-model="selectedCluster" placeholder="选择集群" style="width: 200px; margin-right: 10px;">
            <el-option
              v-for="cluster in clusters"
              :key="cluster.id"
              :label="cluster.name"
              :value="cluster.id"
            />
          </el-select>
          <el-button type="primary" size="small" @click="refreshData">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- 整体加载状态 -->
    <div class="dashboard-container" v-loading="loading" element-loading-text="加载数据中...">
      <!-- 统计卡片 -->
      <div class="stats-cards">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon pod-icon">
              <el-icon size="32"><Box /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-label">Pod数量</div>
              <div class="stat-value">{{ stats.pods.total }}</div>
              <div class="stat-desc">
                <el-tag type="success" size="small">{{ stats.pods.running }} 运行中</el-tag>
                <el-tag type="warning" size="small">{{ stats.pods.pending }} 等待中</el-tag>
                <el-tag type="danger" size="small">{{ stats.pods.failed }} 失败</el-tag>
              </div>
            </div>
          </div>
        </el-card>

        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon service-icon">
              <el-icon size="32"><Link /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-label">Service数量</div>
              <div class="stat-value">{{ stats.services.total }}</div>
              <div class="stat-desc">
                <el-tag type="primary" size="small">{{ stats.services.clusterIP }} ClusterIP</el-tag>
                <el-tag type="success" size="small">{{ stats.services.nodePort }} NodePort</el-tag>
                <el-tag type="warning" size="small">{{ stats.services.loadBalancer }} LoadBalancer</el-tag>
              </div>
            </div>
          </div>
        </el-card>

        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon pvc-icon">
              <el-icon size="32"><Document /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-label">PVC数量</div>
              <div class="stat-value">{{ stats.pvc.total }}</div>
              <div class="stat-desc">
                <el-tag type="success" size="small">{{ stats.pvc.bound }} 已绑定</el-tag>
                <el-tag type="warning" size="small">{{ stats.pvc.pending }} 等待中</el-tag>
                <el-tag type="danger" size="small">{{ stats.pvc.failed }} 失败</el-tag>
              </div>
            </div>
          </div>
        </el-card>

        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon pv-icon">
              <el-icon size="32"><Cpu /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-label">PV数量</div>
              <div class="stat-value">{{ stats.pv.total }}</div>
              <div class="stat-desc">
                <el-tag type="success" size="small">{{ stats.pv.available }} 可用</el-tag>
                <el-tag type="primary" size="small">{{ stats.pv.bound }} 已绑定</el-tag>
                <el-tag type="danger" size="small">{{ stats.pv.failed }} 失败</el-tag>
              </div>
            </div>
          </div>
        </el-card>

        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon configmap-icon">
              <el-icon size="32"><Document /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-label">ConfigMap数量</div>
              <div class="stat-value">{{ stats.configMaps.total }}</div>
            </div>
          </div>
        </el-card>

        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon secret-icon">
              <el-icon size="32"><Lock /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-label">Secret数量</div>
              <div class="stat-value">{{ stats.secrets.total }}</div>
            </div>
          </div>
        </el-card>

        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon ingress-icon">
              <el-icon size="32"><Link /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-label">Ingress数量</div>
              <div class="stat-value">{{ stats.ingresses.total }}</div>
            </div>
          </div>
        </el-card>
      </div>

      <!-- 节点状态列表 -->
      <el-card shadow="hover" class="nodes-card">
        <template #header>
          <div class="card-header">
            <span>节点状态</span>
          </div>
        </template>
        <div class="nodes-content">
          <el-table :data="nodes" stripe style="width: 100%">
            <el-table-column prop="name" label="节点名称" width="200"></el-table-column>
            <el-table-column prop="status" label="状态">
              <template #default="scope">
                <el-tag :type="scope.row.status === 'Ready' ? 'success' : 'danger'" size="small">{{ scope.row.status }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="role" label="角色"></el-table-column>
            <el-table-column prop="cpuUsage" label="CPU使用率">
              <template #header>
                <div class="column-header">
                  <span>CPU使用率</span>
                  <el-tooltip content="CPU使用率通过Metrics进行取值" placement="top">
                    <el-icon class="info-icon"><InfoFilled /></el-icon>
                  </el-tooltip>
                </div>
              </template>
              <template #default="scope">
                <div v-if="scope.row.cpuUsage > 0" class="usage-progress">
                  <el-progress :percentage="Math.round(scope.row.cpuUsage)" :stroke-width="8" :color="getProgressColor(scope.row.cpuUsage)" show-text></el-progress>
                </div>
                <div v-else class="no-data">-</div>
              </template>
            </el-table-column>
            <el-table-column prop="memoryUsage" label="内存使用率">
              <template #header>
                <div class="column-header">
                  <span>内存使用率</span>
                  <el-tooltip content="内存使用率通过Metrics进行取值" placement="top">
                    <el-icon class="info-icon"><InfoFilled /></el-icon>
                  </el-tooltip>
                </div>
              </template>
              <template #default="scope">
                <div v-if="scope.row.memoryUsage > 0" class="usage-progress">
                  <el-progress :percentage="Math.round(scope.row.memoryUsage)" :stroke-width="8" :color="getProgressColor(scope.row.memoryUsage)" show-text></el-progress>
                </div>
                <div v-else class="no-data">-</div>
              </template>
            </el-table-column>
            <el-table-column prop="pods" label="限制Pod数"></el-table-column>
          </el-table>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { Refresh, Box, Link, Document, Cpu, InfoFilled, Lock } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'

// 获取路由信息
const route = useRoute()

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

// 集群统计数据
const stats = ref({
  pods: { total: 0, running: 0, pending: 0, failed: 0 },
  services: { total: 0, clusterIP: 0, nodePort: 0, loadBalancer: 0 },
  pvc: { total: 0, bound: 0, pending: 0, failed: 0 },
  pv: { total: 0, available: 0, bound: 0, failed: 0 },
  configMaps: { total: 0 },
  secrets: { total: 0 },
  ingresses: { total: 0 }
})

// 节点数据类型定义
interface Node {
  name: string
  status: string
  role: string
  cpuUsage: number
  memoryUsage: number
  pods: number
}

// 节点列表数据
const nodes = ref<Node[]>([])

// 组件挂载时加载数据
onMounted(() => {
  loadClusters()
})

// 监听集群选择变化，加载对应的仪表盘数据
watch(selectedCluster, (newClusterId) => {
  if (newClusterId) {
    loadDashboardData()
  }
})

// 加载集群列表
const loadClusters = async () => {
  try {
    // 调用真实的API获取集群列表，axios会自动添加baseURL
    const response = await axios.get('/clusters')
    clusters.value = response.data
    
    // 从URL查询参数获取clusterId
    const clusterId = parseInt(route.query.clusterId as string)
    if (clusterId && !isNaN(clusterId)) {
      // 检查集群是否存在
      const clusterExists = clusters.value.some(cluster => cluster.id === clusterId)
      if (clusterExists) {
        selectedCluster.value = clusterId
      }
    }
  } catch (error) {
    console.error('加载集群列表失败:', error)
    // 不显示错误消息给用户，只在控制台打印
    clusters.value = []
  }
}

// 加载仪表盘数据
const loadDashboardData = async () => {
  if (!selectedCluster.value) {
    ElMessage.warning('请先选择一个集群')
    return
  }
  
  loading.value = true
  try {
    // 并行获取各种数据
    const [dashboardResponse, nodesResponse] = await Promise.all([
      axios.get(`/clusters/${selectedCluster.value}/dashboard`),
      axios.get(`/clusters/${selectedCluster.value}/nodes`)
    ])
    
    // 更新统计数据，保留新添加的字段
    const dashboardStats = dashboardResponse.data.stats
    stats.value = {
      ...stats.value,
      pods: dashboardStats.pods || { total: 0, running: 0, pending: 0, failed: 0 },
      services: dashboardStats.services || { total: 0, clusterIP: 0, nodePort: 0, loadBalancer: 0 },
      pvc: dashboardStats.pvc || { total: 0, bound: 0, pending: 0, failed: 0 },
      pv: dashboardStats.pv || { total: 0, available: 0, bound: 0, failed: 0 },
      // 保留新添加的字段，确保它们不会被覆盖
      configMaps: dashboardStats.configMaps || stats.value.configMaps || { total: 0 },
      secrets: dashboardStats.secrets || stats.value.secrets || { total: 0 },
      ingresses: dashboardStats.ingresses || stats.value.ingresses || { total: 0 }
    }
    
    // 更新节点数据
    nodes.value = nodesResponse.data.map((node: any) => ({
      name: node.name,
      status: node.status,
      role: node.role || 'worker',
      cpuUsage: node.cpuUsage || 0,
      memoryUsage: node.memoryUsage || 0,
      pods: node.pods || 0
    }))
    
    loading.value = false
  } catch (error) {
    console.error('加载仪表盘数据失败:', error)
    loading.value = false
    ElMessage.error('加载仪表盘数据失败')
  }
}

// 刷新数据
const refreshData = () => {
  loadDashboardData()
}

// 获取进度条颜色
const getProgressColor = (percentage: number) => {
  if (percentage < 60) {
    return '#67C23A'
  } else if (percentage < 80) {
    return '#E6A23C'
  } else {
    return '#F56C6C'
  }
}
</script>

<style scoped>
.cluster-dashboard-view {
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

/* 统计卡片样式 */
.stats-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
}

/* 仪表盘容器样式 */
.dashboard-container {
  width: 100%;
}

.stat-card {
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1);
}

/* 使用率进度条样式 */
.usage-progress {
  display: flex;
  align-items: center;
  gap: 10px;
}

.usage-value {
  font-size: 12px;
  font-weight: 500;
  min-width: 45px;
  text-align: right;
}

.no-data {
  color: #909399;
  font-size: 12px;
  height: 20px;
  display: flex;
  align-items: center;
}

/* 列标题样式 */
.column-header {
  display: flex;
  align-items: center;
  gap: 5px;
}

/* 信息图标样式 */
.info-icon {
  font-size: 14px;
  color: #909399;
  cursor: help;
  transition: color 0.3s ease;
}

.info-icon:hover {
  color: #409eff;
}

.stat-icon {
  width: 64px;
  height: 64px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32px;
}

.pod-icon {
  background-color: #e6f7ff;
  color: #1890ff;
}

.service-icon {
  background-color: #f6ffed;
  color: #52c41a;
}

.pvc-icon {
  background-color: #fff7e6;
  color: #faad14;
}

.pv-icon {
  background-color: #f9f0ff;
  color: #722ed1;
}

/* ConfigMap图标样式 */
.configmap-icon {
  background-color: #f6ffed;
  color: #52c41a;
}

/* Secret图标样式 */
.secret-icon {
  background-color: #fff1f0;
  color: #f5222d;
}

/* Ingress图标样式 */
.ingress-icon {
  background-color: #e6f7ff;
  color: #1890ff;
}

.stat-info {
  flex: 1;
}

.stat-label {
  font-size: 14px;
  color: #606266;
  margin-bottom: 4px;
}

.stat-value {
  font-size: 32px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 8px;
}

.stat-desc {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>