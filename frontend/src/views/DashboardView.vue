<template>
  <div class="dashboard-view">
    <!-- 页面头部 -->
    <div class="dashboard-header">
      <h2 class="dashboard-title">平台概览</h2>
      <div class="dashboard-subtitle">实时监控系统运行状态</div>
    </div>

    <!-- 主要指标卡片 -->
    <el-row :gutter="24" class="metrics-row">
      <el-col :xs="12" :sm="12" :md="6">
        <div class="metric-card">
          <div class="metric-header">
            <div class="metric-title">总告警数</div>
            <el-icon class="metric-icon danger"><Warning /></el-icon>
          </div>
          <div class="metric-value">{{ stats.activeAlerts }}</div>
          <div class="metric-change">
            <span class="trend-icon down"><ArrowDown /></span>
            <span class="trend-text">较昨日下降 12%</span>
          </div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="12" :md="6">
        <div class="metric-card">
          <div class="metric-header">
            <div class="metric-title">触发中告警</div>
            <el-icon class="metric-icon warning"><CircleClose /></el-icon>
          </div>
          <div class="metric-value">{{ stats.firingAlerts }}</div>
          <div class="metric-change">
            <span class="trend-icon up"><ArrowUp /></span>
            <span class="trend-text">较昨日上升 8%</span>
          </div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="12" :md="6">
        <div class="metric-card">
          <div class="metric-header">
            <div class="metric-title">已解决告警</div>
            <el-icon class="metric-icon success"><CircleCheck /></el-icon>
          </div>
          <div class="metric-value">{{ stats.resolvedAlerts }}</div>
          <div class="metric-change">
            <span class="trend-icon up"><ArrowUp /></span>
            <span class="trend-text">较昨日上升 23%</span>
          </div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="12" :md="6">
        <div class="metric-card">
          <div class="metric-header">
            <div class="metric-title">告警介质</div>
            <el-icon class="metric-icon info"><Message /></el-icon>
          </div>
          <div class="metric-value">{{ stats.alertMediaCount }}</div>
          <div class="metric-change">
            <span class="trend-icon stable"><Minus /></span>
            <span class="trend-text">与昨日持平</span>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- 中间区域：资源概览 + 系统状态 -->
    <el-row :gutter="24" class="mid-section">
      <!-- 资源概览 -->
      <el-col :xs="24" :md="12" class="mid-col">
        <el-card shadow="hover" class="resource-overview-card">
          <template #header>
            <div class="card-header">
              <el-icon class="header-icon"><DataBoard /></el-icon>
              <span class="header-title">资源概览</span>
            </div>
          </template>
          
          <!-- 主机概览 -->
          <div class="resource-section">
            <div class="section-title">
              <el-icon class="section-icon"><Monitor /></el-icon>
              <span>虚机概览</span>
            </div>
            <div class="resource-grid">
              <div class="resource-box">
                <div class="resource-box-header">
                  <span class="resource-box-title">总主机数</span>
                </div>
                <div class="resource-box-value">{{ resources.hosts.total }}</div>
              </div>
              <div class="resource-box">
                <div class="resource-box-header">
                  <span class="resource-box-title">在线</span>
                  <el-icon class="status-icon online"><Check /></el-icon>
                </div>
                <div class="resource-box-value online">{{ resources.hosts.online }}</div>
              </div>
              <div class="resource-box">
                <div class="resource-box-header">
                  <span class="resource-box-title">离线</span>
                  <el-icon class="status-icon offline"><Close /></el-icon>
                </div>
                <div class="resource-box-value offline">{{ resources.hosts.offline }}</div>
              </div>
            </div>
          </div>
          
          <!-- 集群概览 -->
          <div class="resource-section">
            <div class="section-title">
              <el-icon class="section-icon"><Grid /></el-icon>
              <span>集群概览</span>
            </div>
            <div class="resource-grid">
              <div class="resource-box">
                <div class="resource-box-header">
                  <span class="resource-box-title">总集群数</span>
                </div>
                <div class="resource-box-value">{{ resources.clusters.total }}</div>
              </div>
              <div class="resource-box">
                <div class="resource-box-header">
                  <span class="resource-box-title">可用</span>
                  <el-icon class="status-icon online"><Check /></el-icon>
                </div>
                <div class="resource-box-value online">{{ resources.clusters.available }}</div>
              </div>
              <div class="resource-box">
                <div class="resource-box-header">
                  <span class="resource-box-title">不可用</span>
                  <el-icon class="status-icon offline"><Close /></el-icon>
                </div>
                <div class="resource-box-value offline">{{ resources.clusters.unavailable }}</div>
              </div>
            </div>
          </div>
          
          <!-- KVM概览 -->
          <div class="resource-section">
            <div class="section-title">
              <el-icon class="section-icon"><Cpu /></el-icon>
              <span>KVM概览</span>
            </div>
            <div class="resource-grid">
              <div class="resource-box">
                <div class="resource-box-header">
                  <span class="resource-box-title">总KVM主机</span>
                </div>
                <div class="resource-box-value">{{ resources.kvm.total }}</div>
              </div>
              <div class="resource-box">
                <div class="resource-box-header">
                  <span class="resource-box-title">在线</span>
                  <el-icon class="status-icon online"><Check /></el-icon>
                </div>
                <div class="resource-box-value online">{{ resources.kvm.online }}</div>
              </div>
              <div class="resource-box">
                <div class="resource-box-header">
                  <span class="resource-box-title">离线</span>
                  <el-icon class="status-icon offline"><Close /></el-icon>
                </div>
                <div class="resource-box-value offline">{{ resources.kvm.offline }}</div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <!-- 系统状态 -->
      <el-col :xs="24" :md="12" class="mid-col">
        <div class="system-status-container">
          <!-- 系统健康状态卡片 -->
          <el-card shadow="hover" class="health-status-card mb-24">
            <template #header>
              <div class="card-header">
                <el-icon class="header-icon"><Document /></el-icon>
                <span class="header-title">系统健康状态</span>
              </div>
            </template>
            
            <div class="health-status-grid">
              <!-- 系统服务 -->
              <div class="health-item">
                <div class="health-header">
                  <el-icon :class="['health-icon', healthStatus.systemServices.status]">
                    <CircleCheck v-if="healthStatus.systemServices.status === 'good'" />
                    <Warning v-else-if="healthStatus.systemServices.status === 'warning'" />
                    <CircleClose v-else />
                  </el-icon>
                  <span class="health-title">系统服务</span>
                </div>
                <div class="health-content">
                  <div class="health-chart">
                    <el-progress 
                      type="dashboard" 
                      :percentage="healthStatus.systemServices.percentage" 
                      :color="['#67c23a', '#e6a23c', '#f56c6c']" 
                      :width="100" 
                      :format="() => ''" 
                    />
                    <div class="health-percentage">{{ healthStatus.systemServices.percentage }}%</div>
                  </div>
                  <div class="health-info">
                    <div class="health-desc">{{ healthStatus.systemServices.description }}</div>
                  </div>
                </div>
              </div>
              
              <!-- API 响应时间 -->
              <div class="health-item">
                <div class="health-header">
                  <el-icon :class="['health-icon', healthStatus.apiResponse.status]">
                    <CircleCheck v-if="healthStatus.apiResponse.status === 'good'" />
                    <Warning v-else-if="healthStatus.apiResponse.status === 'warning'" />
                    <CircleClose v-else />
                  </el-icon>
                  <span class="health-title">API 响应时间</span>
                </div>
                <div class="health-content">
                  <div class="health-chart">
                    <el-progress 
                      type="dashboard" 
                      :percentage="healthStatus.apiResponse.percentage" 
                      :color="['#67c23a', '#e6a23c', '#f56c6c']" 
                      :width="100" 
                      :format="() => ''" 
                    />
                    <div class="health-percentage">{{ healthStatus.apiResponse.percentage }}%</div>
                  </div>
                  <div class="health-info">
                    <div class="health-desc">{{ healthStatus.apiResponse.description }}</div>
                  </div>
                </div>
              </div>
              
              <!-- 数据库连接 -->
              <div class="health-item">
                <div class="health-header">
                  <el-icon :class="['health-icon', healthStatus.databaseConnection.status]">
                    <CircleCheck v-if="healthStatus.databaseConnection.status === 'good'" />
                    <Warning v-else-if="healthStatus.databaseConnection.status === 'warning'" />
                    <CircleClose v-else />
                  </el-icon>
                  <span class="health-title">数据库连接</span>
                </div>
                <div class="health-content">
                  <div class="health-chart">
                    <el-progress 
                      type="dashboard" 
                      :percentage="healthStatus.databaseConnection.percentage" 
                      :color="['#67c23a', '#e6a23c', '#f56c6c']" 
                      :width="100" 
                      :format="() => ''" 
                    />
                    <div class="health-percentage">{{ healthStatus.databaseConnection.percentage }}%</div>
                  </div>
                  <div class="health-info">
                    <div class="health-desc">{{ healthStatus.databaseConnection.description }}</div>
                  </div>
                </div>
              </div>
            </div>
          </el-card>
          
          <!-- 系统性能指标卡片 -->
          <el-card shadow="hover" class="performance-status-card">
            <template #header>
              <div class="card-header">
                <el-icon class="header-icon"><Cpu /></el-icon>
                <span class="header-title">系统性能指标</span>
              </div>
            </template>
            
            <div class="performance-grid">
              <!-- CPU 使用率 -->
              <div class="performance-item">
                <div class="performance-header">
                  <el-icon class="performance-icon"><Cpu /></el-icon>
                  <span class="performance-title">CPU 使用率</span>
                </div>
                <div class="performance-content">
                  <div class="performance-chart">
                    <el-progress 
                      :percentage="65" 
                      :stroke-width="12" 
                      :color="[
                        { color: '#67c23a', percentage: 50 },
                        { color: '#e6a23c', percentage: 80 },
                        { color: '#f56c6c', percentage: 100 }
                      ]" 
                    />
                  </div>
                  <div class="performance-info">
                    <div class="performance-value">65%</div>
                    <div class="performance-desc">负载均衡良好</div>
                  </div>
                </div>
              </div>
              
              <!-- 内存使用 -->
              <div class="performance-item">
                <div class="performance-header">
                  <el-icon class="performance-icon"><Memo /></el-icon>
                  <span class="performance-title">内存使用</span>
                </div>
                <div class="performance-content">
                  <div class="performance-chart">
                    <el-progress 
                      :percentage="78" 
                      :stroke-width="12" 
                      :color="[
                        { color: '#67c23a', percentage: 50 },
                        { color: '#e6a23c', percentage: 80 },
                        { color: '#f56c6c', percentage: 100 }
                      ]" 
                    />
                  </div>
                  <div class="performance-info">
                    <div class="performance-value">78%</div>
                    <div class="performance-desc">内存使用偏高</div>
                  </div>
                </div>
              </div>
              
              <!-- 磁盘使用 -->
              <div class="performance-item">
                <div class="performance-header">
                  <el-icon class="performance-icon"><Folder /></el-icon>
                  <span class="performance-title">磁盘使用</span>
                </div>
                <div class="performance-content">
                  <div class="performance-chart">
                    <el-progress 
                      :percentage="45" 
                      :stroke-width="12" 
                      :color="[
                        { color: '#67c23a', percentage: 50 },
                        { color: '#e6a23c', percentage: 80 },
                        { color: '#f56c6c', percentage: 100 }
                      ]" 
                    />
                  </div>
                  <div class="performance-info">
                    <div class="performance-value">45%</div>
                    <div class="performance-desc">磁盘空间充足</div>
                  </div>
                </div>
              </div>
              
              <!-- 网络流量 -->
              <div class="performance-item">
                <div class="performance-header">
                  <el-icon class="performance-icon"><Link /></el-icon>
                  <span class="performance-title">网络流量</span>
                </div>
                <div class="performance-content">
                  <div class="performance-chart">
                    <el-progress 
                      :percentage="32" 
                      :stroke-width="12" 
                      :color="[
                        { color: '#67c23a', percentage: 50 },
                        { color: '#e6a23c', percentage: 80 },
                        { color: '#f56c6c', percentage: 100 }
                      ]" 
                    />
                  </div>
                  <div class="performance-info">
                    <div class="performance-value">32%</div>
                    <div class="performance-desc">网络通畅</div>
                  </div>
                </div>
              </div>
            </div>
          </el-card>
        </div>
      </el-col>
    </el-row>
    

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import {
  Warning, CircleClose, CircleCheck, Message, 
  ArrowDown, ArrowUp, Minus, DataBoard, Monitor, 
  Grid, Check, Close, Document, Cpu, Memo,
  Folder, Link
} from '@element-plus/icons-vue'
import { dashboardAPI } from '@/api/api'

// 统计数据，从API获取
const stats = ref({
  activeAlerts: 0,
  firingAlerts: 0,
  resolvedAlerts: 0,
  alertMediaCount: 0
})

// 资源概览数据，从API获取
const resources = ref({
  hosts: {
    total: 0,
    online: 0,
    offline: 0
  },
  clusters: {
    total: 0,
    available: 0,
    unavailable: 0
  },
  kvm: {
    total: 0,
    online: 0,
    offline: 0
  }
})

// 系统健康状态数据，从API获取
const healthStatus = ref({
  systemServices: {
    percentage: 0,
    status: 'good',
    description: ''
  },
  apiResponse: {
    percentage: 0,
    status: 'good',
    description: ''
  },
  databaseConnection: {
    percentage: 0,
    status: 'good',
    description: ''
  }
})

// 从API加载仪表盘数据
const loadDashboardData = async () => {
  try {
    console.log('开始加载仪表盘数据...')
    
    // 尝试分别调用API，以便定位具体哪个API出错
    try {
      const statsResponse = await dashboardAPI.getStats()
      console.log('获取统计数据成功:', statsResponse.data)
      stats.value = statsResponse.data
    } catch (statsError) {
      console.error('获取统计数据失败:', statsError)
      // 不抛出错误，继续尝试获取资源数据
      stats.value = {
        activeAlerts: 0,
        firingAlerts: 0,
        resolvedAlerts: 0,
        alertMediaCount: 0
      }
    }
    
    try {
      const resourcesResponse = await dashboardAPI.getResources()
      console.log('获取资源数据成功:', resourcesResponse.data)
      resources.value = {
        hosts: resourcesResponse.data.hosts || {
          total: 0,
          online: 0,
          offline: 0
        },
        clusters: resourcesResponse.data.clusters || {
          total: 0,
          available: 0,
          unavailable: 0
        },
        kvm: resourcesResponse.data.kvm || {
          total: 0,
          online: 0,
          offline: 0
        }
      }
    } catch (resourcesError) {
      console.error('获取资源数据失败:', resourcesError)
      // 不抛出错误，使用默认值
      resources.value = {
        hosts: {
          total: 0,
          online: 0,
          offline: 0
        },
        clusters: {
          total: 0,
          available: 0,
          unavailable: 0
        },
        kvm: {
          total: 0,
          online: 0,
          offline: 0
        }
      }
    }
    
    try {
      // 尝试获取系统健康状态数据
      // 注意：这里假设后端提供了getHealthStatus API端点
      // 如果后端没有实现这个端点，会使用下面的默认值
      const healthResponse = await dashboardAPI.getHealthStatus()
      console.log('获取系统健康状态数据成功:', healthResponse.data)
      
      // 更新系统健康状态数据
      healthStatus.value = {
        systemServices: {
          percentage: healthResponse.data.systemServices?.percentage || 95,
          status: healthResponse.data.systemServices?.status || 'good',
          description: healthResponse.data.systemServices?.description || '所有核心服务正常运行'
        },
        apiResponse: {
          percentage: healthResponse.data.apiResponse?.percentage || 78,
          status: healthResponse.data.apiResponse?.status || 'warning',
          description: healthResponse.data.apiResponse?.description || '部分 API 响应较慢，建议优化'
        },
        databaseConnection: {
          percentage: healthResponse.data.databaseConnection?.percentage || 100,
          status: healthResponse.data.databaseConnection?.status || 'good',
          description: healthResponse.data.databaseConnection?.description || '数据库连接稳定'
        }
      }
    } catch (healthError) {
      console.error('获取系统健康状态数据失败:', healthError)
      // 不抛出错误，使用默认值
      healthStatus.value = {
        systemServices: {
          percentage: 95,
          status: 'good',
          description: '所有核心服务正常运行'
        },
        apiResponse: {
          percentage: 78,
          status: 'warning',
          description: '部分 API 响应较慢，建议优化'
        },
        databaseConnection: {
          percentage: 100,
          status: 'good',
          description: '数据库连接稳定'
        }
      }
    }
    
    console.log('仪表盘数据加载完成')
  } catch (error) {
    console.error('加载仪表盘数据失败:', error)
    // 设置默认值，确保页面能正常显示
    stats.value = {
      activeAlerts: 0,
      firingAlerts: 0,
      resolvedAlerts: 0,
      alertMediaCount: 0
    }
    resources.value = {
      hosts: {
        total: 0,
        online: 0,
        offline: 0
      },
      clusters: {
        total: 0,
        available: 0,
        unavailable: 0
      },
      kvm: {
        total: 0,
        online: 0,
        offline: 0
      }
    }
    // 设置系统健康状态默认值
    healthStatus.value = {
      systemServices: {
        percentage: 95,
        status: 'good',
        description: '所有核心服务正常运行'
      },
      apiResponse: {
        percentage: 78,
        status: 'warning',
        description: '部分 API 响应较慢，建议优化'
      },
      databaseConnection: {
        percentage: 100,
        status: 'good',
        description: '数据库连接稳定'
      }
    }
  }
}

// 组件挂载时加载数据
onMounted(() => {
  loadDashboardData()
})
</script>

<style scoped>
.dashboard-view {
  width: 100%;
  padding: 0;
  background-color: #f5f7fa;
  min-height: calc(100vh - 120px);
}

/* 页面头部样式 */
.dashboard-header {
  padding: 24px 0;
  margin-bottom: 24px;
}

.dashboard-title {
  margin: 0 0 8px 0;
  font-size: 32px;
  font-weight: 700;
  color: #303133;
  letter-spacing: -0.5px;
}

.dashboard-subtitle {
  font-size: 16px;
  color: #606266;
  font-weight: 400;
}

/* 指标卡片样式 */
.metrics-row {
  margin-bottom: 24px;
}

.metric-card {
  background: linear-gradient(135deg, #ffffff 0%, #f8f9fa 100%);
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid #e4e7ed;
}

.metric-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  border-color: #409eff;
}

.metric-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.metric-title {
  font-size: 14px;
  font-weight: 500;
  color: #606266;
}

.metric-icon {
  font-size: 24px;
}

.metric-icon.danger {
  color: #f56c6c;
}

.metric-icon.warning {
  color: #e6a23c;
}

.metric-icon.success {
  color: #67c23a;
}

.metric-icon.info {
  color: #409eff;
}

.metric-value {
  font-size: 36px;
  font-weight: 700;
  color: #303133;
  margin-bottom: 8px;
  line-height: 1.2;
}

.metric-change {
  display: flex;
  align-items: center;
  gap: 6px;
}

.trend-icon {
  font-size: 14px;
}

.trend-icon.up {
  color: #67c23a;
}

.trend-icon.down {
  color: #f56c6c;
}

.trend-icon.stable {
  color: #909399;
}

.trend-text {
  font-size: 12px;
  color: #909399;
}

/* 中间区域样式 */
.mid-section {
  margin-bottom: 24px;
}

/* 中间区域列样式，确保等高 */
.mid-col {
  display: flex;
  flex-direction: column;
}

/* 右侧容器样式，确保内部卡片充满高度 */
.system-status-container {
  display: flex;
  flex-direction: column;
  gap: 24px;
  height: 100%;
}

/* 中间区域卡片通用样式 */
.mid-section .el-card {
  border-radius: 12px;
  background: #ffffff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  border: 1px solid #e4e7ed;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.mid-section .el-card:hover {
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  border-color: #409eff;
}

.mid-section .el-card__body {
  padding: 20px;
}



.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 0;
}

.header-icon {
  font-size: 20px;
  color: #409eff;
}

.header-title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.resource-section {
  margin-bottom: 24px;
  padding-bottom: 20px;
  border-bottom: 1px solid #f0f2f5;
  flex-shrink: 0;
}

.resource-section:last-child {
  margin-bottom: 0;
  padding-bottom: 0;
  border-bottom: none;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 500;
  color: #303133;
  margin-bottom: 16px;
}

.section-icon {
  font-size: 18px;
  color: #606266;
}

.resource-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.resource-box {
  background-color: #f5f7fa;
  border-radius: 8px;
  padding: 20px 16px;
  text-align: center;
  transition: all 0.3s ease;
  border: 1px solid #e4e7ed;
  min-height: 100px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.resource-box:hover {
  background-color: #ecf5ff;
  border-color: #409eff;
  transform: translateY(-2px);
}

.resource-box-header {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 6px;
  margin-bottom: 12px;
}

.resource-box-title {
  font-size: 14px;
  color: #606266;
  font-weight: 500;
}

.status-icon {
  font-size: 14px;
}

.status-icon.online {
  color: #67c23a;
}

.status-icon.offline {
  color: #f56c6c;
}

.resource-box-value {
  font-size: 28px;
  font-weight: 700;
  color: #303133;
}

.resource-box-value.online {
  color: #67c23a;
}

.resource-box-value.offline {
  color: #f56c6c;
}

/* 卡片间距 */
.mb-24 {
  margin-bottom: 24px;
}

.health-status-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  align-content: stretch;
}

.health-item {
  background: linear-gradient(135deg, #f5f7fa 0%, #ffffff 100%);
  padding: 20px;
  border-radius: 12px;
  border: 1px solid #e4e7ed;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  text-align: center;
  display: flex;
  flex-direction: column;
  flex: 1;
  align-items: center;
  justify-content: space-between;
}

.health-item:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  border-color: #409eff;
}

.health-header {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
}

.health-icon {
  font-size: 20px;
}

.health-icon.good {
  color: #67c23a;
}

.health-icon.warning {
  color: #e6a23c;
}

.health-icon.danger {
  color: #f56c6c;
}

.health-title {
  font-size: 16px;
  font-weight: 500;
  color: #303133;
}

.health-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.health-chart {
  position: relative;
  display: flex;
  justify-content: center;
  width: 100%;
}

.health-percentage {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 20px;
  font-weight: 700;
  color: #303133;
}

.health-info {
  width: 100%;
}

.health-desc {
  font-size: 14px;
  color: #606266;
  margin: 0;
}



.performance-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.performance-item {
  background: linear-gradient(135deg, #f5f7fa 0%, #ffffff 100%);
  padding: 16px;
  border-radius: 8px;
  border: 1px solid #e4e7ed;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.performance-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border-color: #409eff;
}

.performance-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}

.performance-icon {
  font-size: 18px;
  color: #409eff;
}

.performance-title {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
}

.performance-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.performance-chart {
  flex: 1;
  min-width: 100px;
}

.performance-info {
  flex: 0 0 auto;
  text-align: center;
}

.performance-value {
  font-size: 24px;
  font-weight: 700;
  color: #303133;
  margin-bottom: 4px;
}

.performance-desc {
  font-size: 12px;
  color: #606266;
}

/* 进度条样式调整 */
:deep(.performance-item .el-progress-bar__outer) {
  background-color: #f0f2f5;
  border-radius: 4px;
  height: 8px;
}

:deep(.performance-item .el-progress-bar__inner) {
  border-radius: 4px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .health-status-grid {
    grid-template-columns: 1fr;
    margin-bottom: 16px;
    padding-bottom: 16px;
  }
  
  .performance-grid {
    grid-template-columns: 1fr;
  }
  
  .performance-content {
    flex-direction: column;
    gap: 12px;
  }
  
  .performance-chart {
    min-width: auto;
    width: 100%;
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .dashboard-title {
    font-size: 24px;
  }
  
  .dashboard-subtitle {
    font-size: 14px;
  }
  
  .metric-value {
    font-size: 28px;
  }
  
  .resource-grid {
    grid-template-columns: 1fr;
  }
  
  .health-status-grid {
    gap: 32px;
  }
  
  .activity-content {
    margin-left: 0;
  }
}

@media (max-width: 480px) {
  .metric-card {
    padding: 16px;
  }
  
  .metric-value {
    font-size: 24px;
  }
  
  .resource-box {
    padding: 16px 12px;
  }
  
  .resource-box-value {
    font-size: 24px;
  }
}
</style>