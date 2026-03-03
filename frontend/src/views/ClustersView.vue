<template>
  <div class="clusters-view">
    <div class="page-header">
      <div class="header-content">
        <h2>集群列表</h2>
      </div>
    </div>

    <div class="table-container">
      <div class="table-header">
        <span>K8s集群列表</span>
        <div class="header-actions">
          <el-button type="primary" @click="openAddClusterDialog">
            <el-icon><Plus /></el-icon>
            导入集群
          </el-button>
        </div>
      </div>
      
      <el-table 
        :data="clusters" 
        style="width: 100%" 
        stripe
        fit
        v-loading="loading"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="集群名称" min-width="150" />
        <el-table-column prop="status" label="集群状态" min-width="120">
          <template #default="scope">
            <el-tag :type="getStatusTagType(scope.row.status)" size="small">
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="version" label="K8s版本" min-width="100" />
        <el-table-column label="节点数" min-width="100">
          <template #header>
            <div class="column-header">
              <span>节点数</span>
              <el-tooltip content="正常节点/总节点" placement="top">
                <el-icon class="info-icon"><InfoFilled /></el-icon>
              </el-tooltip>
            </div>
          </template>
          <template #default="scope">
            {{ scope.row.nodeCount || '0/0' }}
          </template>
        </el-table-column>
        <el-table-column label="导入方式" min-width="120">
          <template #default="scope">
            <el-tag size="small">
              {{ getImportMethod(scope.row.connectionType) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="导入时间" min-width="160">
          <template #default="scope">
            {{ formatDate(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" min-width="280" fixed="right">
          <template #default="scope">
            <div style="display: flex; gap: 5px; white-space: nowrap;">
              <el-button type="primary" size="small" @click="testConnection(scope.row)" :loading="testConnectionLoading[scope.row.id]">
                <el-icon v-if="!testConnectionLoading[scope.row.id]"><Check /></el-icon>
                测试连接
              </el-button>
              <el-button type="success" size="small" @click="goToClusterDashboard(scope.row)">
                <el-icon><Monitor /></el-icon>
                监控
              </el-button>
              <el-button type="info" size="small" @click="openEditClusterDialog(scope.row)">
                <el-icon><Edit /></el-icon>
                编辑
              </el-button>
              <el-button type="danger" size="small" @click="deleteCluster(scope.row.id)">
                <el-icon><Delete /></el-icon>
                删除
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </div>
    
    <!-- 添加/编辑集群对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑集群' : '导入集群'"
      width="600px"
    >
      <el-form :model="clusterForm" :rules="clusterRules" ref="clusterFormRef" label-width="120px">
        <el-form-item label="集群名称" prop="name" required>
          <el-input v-model="clusterForm.name" placeholder="请输入集群名称" />
        </el-form-item>
        
        <el-form-item label="集群描述">
          <el-input v-model="clusterForm.description" type="textarea" :rows="3" placeholder="请输入集群描述" />
        </el-form-item>
        
        <el-form-item label="导入方式">
          <el-radio-group v-model="clusterForm.connectionType" @change="onConnectionTypeChange" size="large">
            <el-radio-button label="kubeconfig">kubeconfig 文件导入</el-radio-button>
            <el-radio-button label="token">Secret Token导入</el-radio-button>
          </el-radio-group>
        </el-form-item>
        
        <!-- Secret Token方式连接 -->
        <template v-if="clusterForm.connectionType === 'token'">
          <el-form-item label="API地址" prop="apiURL">
            <el-input v-model="clusterForm.apiURL" placeholder="请输入K8s API地址" />
          </el-form-item>
          
          <el-form-item label="Token" prop="token">
            <el-input v-model="clusterForm.token" placeholder="请输入集群访问Token" type="textarea" :rows="6" />
          </el-form-item>
        </template>
        
        <!-- KubeConfig文件方式连接 -->
        <template v-else>
          <el-form-item label="kubeconfig">
            <el-input
              v-model="clusterForm.kubeconfig"
              type="textarea"
              :rows="12"
              placeholder="请将您的 '/etc/kubernetes/admin.conf' 文件粘贴到此处，仅支持使用 kubeadm 安装集群时生成的 kubeconfig 文件（TKE、kubespary 等工具安装的集群也使用 kubeadm）。"
              style="fontFamily: 'monospace'"
            >
              <template #prepend>
                <div style="backgroundColor: #f5f7fa; color: #909399; padding: '0 10px'; borderRight: '1px solid #e4e7ed'">
                  1
                </div>
              </template>
            </el-input>
            <div style="margin-top: 10px; color: #666; font-size: 14px;">
              <p>您可以在集群控制节点上执行此命令获得 kubeconfig 文件的内容：</p>
              <p style="font-family: monospace; background-color: #f5f5f5; padding: 5px 10px; border-radius: 3px;">
                cat /etc/kubernetes/admin.conf
              </p>
            </div>
          </el-form-item>
        </template>
        

      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveCluster" :loading="saving">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus, Edit, Delete, Check, Monitor, InfoFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'

// 表单引用
const clusterFormRef = ref()

// 加载状态
const loading = ref(false)
const saving = ref(false)
// 测试连接的加载状态，使用对象存储每个集群的加载状态
const testConnectionLoading = ref<Record<number, boolean>>({})

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

// 集群数据
const clusters = ref<Cluster[]>([])

// 对话框相关
const dialogVisible = ref(false)
const isEdit = ref(false)
const clusterForm = ref({
  id: 0,
  name: '',
  apiURL: '',
  token: '',
  kubeconfig: '',
  description: '',
  connectionType: 'kubeconfig'
})

// 表单验证规则
const clusterRules = {
  name: [{ required: true, message: '请输入集群名称', trigger: 'blur' }],
  apiURL: [{ 
    required: true, 
    message: '请输入API地址', 
    trigger: 'blur',
    validator: (_rule: any, value: any, callback: any) => {
      if (clusterForm.value.connectionType === 'token' && !value) {
        callback(new Error('请输入API地址'))
      } else {
        callback()
      }
    }
  }],
  token: [{ 
    required: true, 
    message: '请输入Token', 
    trigger: 'blur',
    validator: (_rule: any, value: any, callback: any) => {
      if (clusterForm.value.connectionType === 'token' && !value) {
        callback(new Error('请输入Token'))
      } else {
        callback()
      }
    }
  }],
  kubeconfig: [{ 
    required: true, 
    message: '请输入kubeconfig内容', 
    trigger: 'blur',
    validator: (_rule: any, value: any, callback: any) => {
      if (clusterForm.value.connectionType === 'kubeconfig' && !value) {
        callback(new Error('请输入kubeconfig内容'))
      } else {
        callback()
      }
    }
  }]
}

// 处理连接方式变化
const onConnectionTypeChange = () => {
  if (clusterForm.value.connectionType === 'token') {
    clusterForm.value.kubeconfig = ''
  } else {
    clusterForm.value.apiURL = ''
    clusterForm.value.token = ''
  }
  // 重置表单验证
  if (clusterFormRef.value) {
    clusterFormRef.value.clearValidate()
  }
}

// 组件挂载时加载数据
onMounted(() => {
  loadClusters()
})

// 加载集群列表
const loadClusters = async () => {
  loading.value = true
  try {
    // 调用真实的API获取集群列表，添加超时配置
    const response = await axios.get('/clusters', { timeout: 10000 }) // 10秒超时
    const clustersData = response.data
    
    // 并行获取所有集群的节点数，使用Promise.all提高性能，并添加超时配置
    const nodeCountPromises = clustersData.map(async (cluster: any) => {
      try {
        const nodesResponse = await axios.get(`/clusters/${cluster.id}/nodes`, { timeout: 8000 }) // 8秒超时
        const nodes = nodesResponse.data
        const totalNodes = nodes.length
        const readyNodes = nodes.filter((node: any) => node.status === 'Ready').length
        return { id: cluster.id, nodeCount: `${readyNodes}/${totalNodes}` }
      } catch (error) {
        console.error(`获取集群${cluster.id}节点数失败:`, error)
        return { id: cluster.id, nodeCount: '0/0' }
      }
    })
    
    // 等待所有节点数请求完成，添加总超时保护
    const nodeCountResults = await Promise.race([
      Promise.all(nodeCountPromises),
      new Promise((_, reject) => {
        setTimeout(() => reject(new Error('获取节点数超时')), 15000) // 15秒总超时
      })
    ])
    
    // 将节点数结果合并到集群数据中
    for (const cluster of clustersData) {
      const nodeCountResult = nodeCountResults.find(result => result.id === cluster.id)
      if (nodeCountResult) {
        cluster.nodeCount = nodeCountResult.nodeCount
      }
    }
    
    clusters.value = clustersData
    loading.value = false
  } catch (error) {
    console.error('加载集群列表失败:', error)
    // 即使出错，也显示已获取的集群列表
    loading.value = false
    // 不将clusters.value设置为空数组，保留已获取的集群数据
  }
}

// 获取状态标签类型
const getStatusTagType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'connected':
    case 'ready':
    case 'running':
    case 'succeeded':
    case 'available':
      return 'success'
    case 'disconnected':
    case 'notready':
    case 'pending':
    case 'progressing':
      return 'warning'
    case 'error':
    case 'failed':
    case 'degraded':
      return 'danger'
    default:
      return 'info'
  }
}

// 获取中文状态文本
const getStatusText = (status: string) => {
  switch (status.toLowerCase()) {
    case 'connected':
      return '已连接'
    case 'disconnected':
      return '已断开'
    case 'error':
      return '错误'
    case 'ready':
      return '就绪'
    case 'notready':
      return '未就绪'
    case 'running':
      return '运行中'
    case 'pending':
      return '待处理'
    case 'succeeded':
      return '成功'
    case 'failed':
      return '失败'
    case 'unknown':
      return '未知'
    case 'available':
      return '可用'
    case 'progressing':
      return '进行中'
    case 'degraded':
      return '降级'
    default:
      return '未知'
  }
}

// 获取导入方式文本
const getImportMethod = (connectionType: string) => {
  switch (connectionType) {
    case 'kubeconfig':
      return 'KubeConfig 文件导入'
    case 'token':
      return 'Token 导入'
    default:
      return '未知'
  }
}

// 格式化日期为 YYYY-MM-DD HH:MM:SS 格式
const formatDate = (dateString: string) => {
  if (!dateString) return ''
  
  const date = new Date(dateString)
  
  // 检查日期是否有效
  if (isNaN(date.getTime())) {
    return ''
  }
  
  const year = date.getFullYear()
  // 月份和日期补零，格式化为两位数
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  
  // 使用 '-' 作为日期分隔符，空格作为日期和时间的分隔符
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

// 打开添加集群对话框
const openAddClusterDialog = () => {
  isEdit.value = false
  clusterForm.value = {
    id: 0,
    name: '',
    apiURL: '',
    token: '',
    kubeconfig: '',
    description: '',
    connectionType: 'kubeconfig'
  }
  dialogVisible.value = true
}

// 打开编辑集群对话框
const openEditClusterDialog = (cluster: any) => {
  isEdit.value = true
  clusterForm.value = {
    id: cluster.id,
    name: cluster.name,
    apiURL: cluster.apiURL,
    token: cluster.token,
    kubeconfig: cluster.kubeconfig,
    description: cluster.description,
    connectionType: cluster.kubeconfig ? 'kubeconfig' : 'token'
  }
  dialogVisible.value = true
}

// 保存集群
const saveCluster = async () => {
  if (!clusterFormRef.value) return
  
  clusterFormRef.value.validate(async (valid: boolean) => {
    if (valid) {
      saving.value = true
      try {
        if (isEdit.value) {
          // 编辑现有集群，使用PUT请求
          await axios.put(`/clusters/${clusterForm.value.id}`, clusterForm.value)
        } else {
          // 添加新集群，使用POST请求
          await axios.post('/clusters', clusterForm.value)
        }
        // 重新加载集群列表
        await loadClusters()
        dialogVisible.value = false
        saving.value = false
        ElMessage.success(isEdit.value ? '编辑集群成功' : '导入集群成功')
      } catch (error: any) {
        console.error('保存集群失败:', error)
        saving.value = false
        // 显示更具体的错误信息
        let errorMsg = isEdit.value ? '编辑集群失败' : '添加集群失败'
        if (error.response?.data?.message) {
          errorMsg = error.response.data.message
          // 提取更友好的错误信息，如集群名称已存在
          if (error.response.data.message.includes('name already exists') || error.response.data.message.includes('名称已存在')) {
            errorMsg = '集群名称已存在'
          }
        }
        ElMessage.error(errorMsg)
      }
    }
  })
}

// 删除集群
const deleteCluster = async (id: number) => {
  try {
    // 调用真实的API彻底删除集群，axios会自动添加baseURL
    await axios.delete(`/clusters/${id}`)
    // 从列表中移除删除的集群
    clusters.value = clusters.value.filter((c: any) => c.id !== id)
    ElMessage.success('删除集群成功')
  } catch (error: any) {
    console.error('删除集群失败:', error)
    // 显示更具体的错误信息
    let errorMsg = '删除集群失败'
    if (error.response?.data?.message) {
      errorMsg = error.response.data.message
    }
    ElMessage.error(errorMsg)
  }
}



// 测试集群连接
const testConnection = async (cluster: any) => {
  // 设置当前集群的测试连接加载状态为true
  testConnectionLoading.value[cluster.id] = true
  
  try {
    // 调用真实的API测试集群连接，axios会自动添加baseURL
    const response = await axios.post(`/clusters/${cluster.id}/test-connection`)
    const index = clusters.value.findIndex((c: any) => c.id === cluster.id)
    if (index !== -1) {
      clusters.value[index].status = 'connected'
      clusters.value[index].version = response.data.version
    }
    ElMessage.success('连接测试成功')
  } catch (error: any) {
    console.error('连接测试失败:', error)
    // 优化错误信息显示
    let errorMessage = '连接测试失败'
    const errorDetails = error.message || ''
    
    // 更新集群状态
    const index = clusters.value.findIndex((c: any) => c.id === cluster.id)
    if (index !== -1) {
      clusters.value[index].status = 'disconnected'
      clusters.value[index].version = ''
    }
    
    // 识别不同类型的错误并显示友好信息
    if (errorDetails.includes('timeout')) {
      errorMessage = '连接超时，请检查网络连接或集群地址'
    } else if (errorDetails.includes('connect: connection refused')) {
      errorMessage = '连接被拒绝，请检查集群地址和端口是否正确'
    } else if (errorDetails.includes('no such host')) {
      errorMessage = '无法解析主机名，请检查集群地址是否正确'
    } else if (errorDetails.includes('unauthorized') || errorDetails.includes('Unauthorized')) {
      errorMessage = '认证失败，请检查Token或KubeConfig是否正确'
    } else if (error.response?.data?.message) {
      // 服务器返回的错误信息
      const serverMsg = error.response.data.message
      if (serverMsg.includes('Failed to connect to K8s cluster')) {
        // 提取更友好的错误信息
        if (serverMsg.includes('timeout')) {
          errorMessage = '连接超时，请检查网络连接或集群地址'
        } else if (serverMsg.includes('connection refused')) {
          errorMessage = '连接被拒绝，请检查集群地址和端口是否正确'
        } else {
          errorMessage = '连接失败，请检查集群配置'
        }
      } else {
        errorMessage = serverMsg
      }
    }
    
    ElMessage.error(errorMessage)
  } finally {
    // 无论成功或失败，都将加载状态设置为false
    testConnectionLoading.value[cluster.id] = false
  }
}

// 跳转到集群仪表盘
const goToClusterDashboard = (cluster: any) => {
  // 跳转到集群仪表盘页面，传递集群ID或名称
  window.location.href = `/cluster-dashboard?clusterId=${cluster.id}&clusterName=${encodeURIComponent(cluster.name)}`
}
</script>

<style scoped>
/* 页面容器样式 */
.clusters-view {
  width: 100%;
  padding: 20px 20px 20px 16px;
  margin: 0;
  background-color: #f5f7fa;
  min-height: 100vh;
  box-sizing: border-box;
}

/* 页面头部样式 */
.page-header {
  width: 100%;
  margin: 0 auto 30px auto;
  padding: 20px;
  background-color: #ffffff;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  border-radius: 8px;
  box-sizing: border-box;
}

/* 页面头部内容样式 */
.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

/* 页面标题样式 */
.header-content h2 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: #333;
}

/* 头部操作按钮样式 */
.header-actions {
  display: flex;
  gap: 10px;
}

/* 表格容器样式 */
.table-container {
  width: 100%;
  margin: 0 auto 20px auto;
  background-color: #ffffff;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  border-radius: 8px;
  box-sizing: border-box;
  padding: 20px;
}

/* 表格头部样式 */
.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  padding-bottom: 15px;
  border-bottom: 1px solid #ebeef5;
}

/* 表格标题文字样式 */
.table-header span {
  font-size: 16px;
  font-weight: 400;
  color: #333;
}

/* 表格样式 */
.el-table {
  width: 100%;
  margin-top: 0;
  overflow-x: auto;
  box-sizing: border-box;
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

/* 表头样式优化 */
:deep(.el-table th.el-table__header-cell) {
  padding: 16px;
  font-weight: 600;
  font-size: 14px;
  background-color: #fafafa;
  color: #333;
}

/* 表行样式优化 */
:deep(.el-table__row) {
  height: 65px;
  transition: all 0.3s ease;
}

/* 表格单元格样式优化 */
:deep(.el-table__cell) {
  padding: 20px 16px;
  font-size: 14px;
  color: #606266;
}

/* 表格行悬停效果 */
:deep(.el-table__row:hover) {
  background-color: #f5f7fa;
}

/* 操作按钮样式优化 */
:deep(.el-button + .el-button) {
  margin-left: 12px;
}

/* 状态标签样式优化 */
:deep(.el-tag) {
  margin: 0 4px;
  font-size: 12px;
}

/* 确保表格列能横铺整个页面 */
:deep(.el-table__header-wrapper),
:deep(.el-table__body-wrapper) {
  width: 100%;
}

/* 调整操作列样式 */
:deep(.el-table__column--fixed-right) {
  right: 0;
  left: auto;
}
</style>
