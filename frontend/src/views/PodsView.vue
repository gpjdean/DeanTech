<template>
  <div class="pods-view">
    <el-card shadow="never" class="page-header">
      <div class="header-content">
        <h2>Pod</h2>
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
          <el-button type="primary" size="small" @click="refreshPods">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </div>
    </el-card>

    <el-card shadow="hover" class="pods-card">
      <template #header>
        <div class="card-header">
          <span>K8s Pod列表</span>
        </div>
      </template>
      
      <el-table 
        :data="pods" 
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
        :row-key="(row: any) => `${row.namespace}-${row.name}`"
      >
        <el-table-column prop="name" label="Pod名称" min-width="200" show-overflow-tooltip />
        <el-table-column prop="namespace" label="命名空间" width="150" show-overflow-tooltip />

        <el-table-column prop="status" label="状态" width="120">
          <template #default="scope">
            <el-tag :type="getStatusTagType(scope.row.status)" size="medium">
              {{ scope.row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="nodeName" label="节点" width="150" show-overflow-tooltip />
        <el-table-column prop="restartCount" label="重启次数" width="120" />
        <el-table-column prop="podIP" label="Pod IP" width="150" show-overflow-tooltip />
        <el-table-column prop="hostIP" label="节点IP" width="150" show-overflow-tooltip />
        <el-table-column prop="createdAt" label="创建时间" width="180">
          <template #default="scope">
            {{ formatAge(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="scope">
            <el-space>
              <el-button type="primary" size="small" @click="viewPodDetails(scope.row)">
                <el-icon><View /></el-icon>
                查看
              </el-button>
              <el-dropdown @command="(command: string) => handleDropdownCommand(command, scope.row)">
                <el-button size="small">
                  更多
                  <el-icon class="el-icon--right"><ArrowDown /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="logs">
                      <el-icon><Message /></el-icon>
                      日志
                    </el-dropdown-item>
                    <el-dropdown-item command="terminal">
                      <el-icon><View /></el-icon>
                      终端
                    </el-dropdown-item>
                    <el-dropdown-item command="delete" type="danger">
                      <el-icon><Delete /></el-icon>
                      删除
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </el-space>
          </template>
        </el-table-column>
      </el-table>
      
      <div v-else class="no-cluster-selected">
        <el-empty description="请选择一个集群查看Pod列表" />
      </div>
    </el-card>
    
    <!-- Pod详情对话框 -->
    <el-dialog
      v-model="podDetailsVisible"
      title="Pod详情"
      width="800px"
    >
      <div v-if="selectedPod" class="pod-details">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="Pod名称">{{ selectedPod.name }}</el-descriptions-item>
          <el-descriptions-item label="命名空间">{{ selectedPod.namespace }}</el-descriptions-item>
          <el-descriptions-item label="状态">{{ selectedPod.status }}</el-descriptions-item>
          <el-descriptions-item label="节点">{{ selectedPod.nodeName }}</el-descriptions-item>
          <el-descriptions-item label="重启次数">{{ selectedPod.restartCount }}</el-descriptions-item>
          <el-descriptions-item label="Pod IP">{{ selectedPod.podIP }}</el-descriptions-item>
          <el-descriptions-item label="节点IP">{{ selectedPod.hostIP }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ formatAge(selectedPod.createdAt) }}</el-descriptions-item>
          <el-descriptions-item label="CPU请求" v-if="selectedPod.resources">
            {{ selectedPod.resources.requests?.cpu || 'N/A' }}
          </el-descriptions-item>
          <el-descriptions-item label="内存请求" v-if="selectedPod.resources">
            {{ selectedPod.resources.requests?.memory || 'N/A' }}
          </el-descriptions-item>
          <el-descriptions-item label="CPU限制" v-if="selectedPod.resources">
            {{ selectedPod.resources.limits?.cpu || 'N/A' }}
          </el-descriptions-item>
          <el-descriptions-item label="内存限制" v-if="selectedPod.resources">
            {{ selectedPod.resources.limits?.memory || 'N/A' }}
          </el-descriptions-item>
        </el-descriptions>
        
        <h4 style="margin-top: 20px;">容器列表</h4>
        <el-table :data="selectedPod.containers" style="width: 100%; margin-top: 10px;">
          <el-table-column prop="name" label="容器名称" width="180" />
          <el-table-column prop="image" label="镜像" min-width="250" />
          <el-table-column prop="status" label="状态" width="120" />
          <el-table-column prop="ready" label="就绪" width="80">
            <template #default="scope">
              <el-icon v-if="scope.row.ready"><Check /></el-icon>
              <el-icon v-else><Close /></el-icon>
            </template>
          </el-table-column>
        </el-table>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="podDetailsVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>
    
    <!-- Pod日志对话框 -->
    <el-dialog
      v-model="podLogsVisible"
      title="Pod日志"
      width="800px"
      height="600px"
      @close="stopLogRefresh"
    >
      <div v-if="selectedPod" class="pod-logs">
        <div class="logs-header">
          <h4>{{ selectedPod.name }} - 容器日志</h4>
          <div class="logs-controls">
            <el-select v-model="selectedContainer" placeholder="选择容器" style="margin-right: 10px;">
              <el-option
                v-for="container in selectedPod.containers"
                :key="container.name"
                :label="container.name"
                :value="container.name"
              />
            </el-select>
            <el-select v-model="logLines" placeholder="选择行数" style="margin-right: 10px; width: 120px;">
              <el-option label="50行" :value="50" />
              <el-option label="100行" :value="100" />
              <el-option label="200行" :value="200" />
              <el-option label="500行" :value="500" />
              <el-option label="1000行" :value="1000" />
            </el-select>
            <el-button 
              type="info" 
              size="small" 
              @click="toggleLogRefresh" 
              :icon="logRefreshing ? 'Stop' : 'Refresh'"
              style="margin-right: 10px;"
            >
              {{ logRefreshing ? '停止刷新' : '实时刷新' }}
            </el-button>
            <el-button 
              type="success" 
              size="small" 
              @click="downloadLogs"
              icon="Download"
            >
              下载日志
            </el-button>
          </div>
        </div>
        <div class="logs-content">
          <pre v-loading="logsLoading">{{ podLogs }}</pre>
        </div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="podLogsVisible = false">关闭</el-button>
          <el-button type="primary" @click="loadPodLogs">刷新日志</el-button>
        </span>
      </template>
    </el-dialog>
    
    <!-- Pod终端对话框 -->
    <el-dialog
      v-model="podTerminalVisible"
      title="Pod终端"
      width="800px"
      height="600px"
      @close="closeTerminal"
    >
      <div v-if="selectedPod" class="pod-terminal">
        <div class="terminal-header">
          <h4>{{ selectedPod.name }} - {{ selectedContainer }} 容器终端</h4>
          <div class="terminal-controls">
            <el-select v-model="selectedContainer" placeholder="选择容器" style="margin-right: 10px;" @change="restartTerminal">
              <el-option
                v-for="container in selectedPod.containers"
                :key="container.name"
                :label="container.name"
                :value="container.name"
              />
            </el-select>
            <el-button 
              type="success" 
              size="small" 
              @click="restartTerminal"
              icon="Refresh"
            >
              重启终端
            </el-button>
          </div>
        </div>
        <div class="terminal-content" ref="terminalContentRef" v-loading="terminalLoading">
          <pre>{{ terminalContent }}</pre>
          <div class="terminal-input-area">
            <span class="terminal-prompt">{{ selectedPod.name }}:~$</span>
            <el-input
              v-model="terminalInput"
              placeholder="输入命令..."
              @keyup.enter="sendTerminalCommand"
              @keyup.up="navigateHistory(-1)"
              @keyup.down="navigateHistory(1)"
              size="small"
              style="width: calc(100% - 120px); margin-left: 10px;"
              auto-complete="off"
            />
          </div>
        </div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="podTerminalVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { Refresh, View, Message, Delete, Check, Close, ArrowDown } from '@element-plus/icons-vue'
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

// 命名空间列表
const namespaces = ref<string[]>([])
const namespacesLoading = ref(false)

// 选择的集群和命名空间
const selectedCluster = ref<number | undefined>(undefined)
const selectedNamespace = ref<string>('')

// Pod数据类型定义
interface PodContainer {
  name: string
  image: string
  status: string
  ready: boolean
}

interface PodResources {
  requests: {
    cpu: string
    memory: string
  }
  limits: {
    cpu: string
    memory: string
  }
}

interface Pod {
  clusterID: number
  name: string
  namespace: string
  cluster: string
  status: string
  nodeName: string
  restartCount: number
  podIP: string
  hostIP: string
  createdAt: string
  containers: PodContainer[]
  resources?: PodResources
}

// Pod数据
const pods = ref<Pod[]>([])

// 对话框相关
const podDetailsVisible = ref(false)
const podLogsVisible = ref(false)
const podTerminalVisible = ref(false)
const selectedPod = ref<Pod | null>(null)
const selectedContainer = ref('')
const podLogs = ref('')
const logLines = ref(100) // 默认显示100行日志
const logsLoading = ref(false) // 日志加载状态
const logRefreshing = ref(false) // 实时刷新标志
let refreshInterval: number | null = null // 刷新间隔ID

// 终端相关
const terminalContent = ref('')
const terminalInput = ref('')
const terminalLoading = ref(false)
const terminalContentRef = ref<HTMLElement | null>(null)
let ws: WebSocket | null = null
let terminalHistory: string[] = []
let historyIndex = -1

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

// 监听命名空间选择变化，自动刷新Pod列表
watch(selectedNamespace, () => {
  if (selectedCluster.value) {
    loadPods()
  }
})

// 监听容器选择变化，自动重新加载日志
watch(selectedContainer, () => {
  if (selectedPod.value && podLogsVisible.value) {
    loadPodLogs()
  }
})

// 监听日志行数变化，自动重新加载日志
watch(logLines, () => {
  if (selectedPod.value && podLogsVisible.value) {
    loadPodLogs()
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

// 加载Pod列表
const loadPods = async () => {
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
    // 调用API获取Pod列表，包含命名空间参数，axios会自动添加baseURL
    const response = await axios.get(`/clusters/${selectedCluster.value}/pods`, {
      params: {
        namespace: selectedNamespace.value || '',
        _t: Date.now() // 添加时间戳，防止浏览器缓存
      },
      headers: {
        'Cache-Control': 'no-cache' // 禁用缓存
      },
      cancelToken: cancelTokenSource.token
    })
    // 为每个Pod添加集群名称
    const clusterName = clusters.value.find(c => c.id === selectedCluster.value)?.name || ''
    const newPods = response.data.map((pod: any) => ({
      ...pod,
      cluster: clusterName
    }))
    
    // 先更新数据，然后设置loading为false，确保数据正确显示
    pods.value = newPods
    loading.value = false
    
    // 如果返回的Pod列表为空，显示提示信息
    if (newPods.length === 0) {
      ElMessage.info('目标命名空间内无Pod实例')
    }
  } catch (error) {
    // 如果是取消请求的错误，不显示错误信息
    if (axios.isCancel(error)) {
      console.log('Request canceled:', error.message)
      return
    }
    console.error('加载Pod列表失败:', error)
    loading.value = false
    ElMessage.error('加载Pod列表失败')
    // 确保在错误情况下，列表为空，显示正确的空状态
    pods.value = []
  }
}

// 刷新Pod列表
const refreshPods = () => {
  loadPods()
}

// 获取状态标签类型
const getStatusTagType = (status: string) => {
  switch (status) {
    case 'Running':
      return 'success'
    case 'Pending':
      return 'warning'
    case 'Failed':
      return 'danger'
    case 'Succeeded':
      return 'info'
    case 'Unknown':
      return 'info'
    default:
      return 'info'
  }
}

// 将ISO时间转换为绝对时间格式 YYYY-MM-DD HH:MM:SS
const formatAge = (createdAt: string) => {
  if (!createdAt) return ''
  
  try {
    const created = new Date(createdAt)
    if (isNaN(created.getTime())) return ''
    
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

// 查看Pod详情
const viewPodDetails = (pod: any) => {
  selectedPod.value = pod
  podDetailsVisible.value = true
}

// 查看Pod日志
const viewPodLogs = (pod: any) => {
  selectedPod.value = pod
  selectedContainer.value = pod.containers[0]?.name || ''
  podLogsVisible.value = true
  loadPodLogs()
}

// 加载Pod日志
const loadPodLogs = async () => {
  if (!selectedPod.value || !selectedCluster.value || !selectedContainer.value) {
    return
  }
  
  logsLoading.value = true
  try {
    // 调用API获取Pod日志，包含命名空间、容器名称和行数参数
    const response = await axios.get(`/clusters/${selectedCluster.value}/pods/${selectedPod.value.name}/logs`, {
      params: {
        namespace: selectedPod.value.namespace,
        container: selectedContainer.value,
        lines: logLines.value
      }
    })
    podLogs.value = response.data
    logsLoading.value = false
  } catch (error) {
    console.error('加载Pod日志失败:', error)
    logsLoading.value = false
    ElMessage.error('加载Pod日志失败')
  }
}

// 切换日志实时刷新
const toggleLogRefresh = () => {
  if (logRefreshing.value) {
    stopLogRefresh()
  } else {
    startLogRefresh()
  }
}

// 开始日志实时刷新
const startLogRefresh = () => {
  if (!selectedPod.value || !selectedCluster.value || !selectedContainer.value) {
    return
  }
  
  logRefreshing.value = true
  // 立即加载一次日志
  loadPodLogs()
  // 设置5秒刷新一次
  refreshInterval = window.setInterval(() => {
    loadPodLogs()
  }, 5000)
}

// 停止日志实时刷新
const stopLogRefresh = () => {
  logRefreshing.value = false
  if (refreshInterval) {
    clearInterval(refreshInterval)
    refreshInterval = null
  }
}

// 下载日志
const downloadLogs = () => {
  if (!selectedPod.value || !selectedContainer.value || !podLogs.value) {
    return
  }
  
  // 创建Blob对象
  const blob = new Blob([podLogs.value], { type: 'text/plain;charset=utf-8' })
  // 创建下载链接
  const link = document.createElement('a')
  link.href = URL.createObjectURL(blob)
  // 设置文件名
  const fileName = `${selectedPod.value.name}-${selectedContainer.value}-${new Date().toISOString().slice(0, 19).replace(/:/g, '-')}.log`
  link.download = fileName
  // 触发下载
  document.body.appendChild(link)
  link.click()
  // 清理
  document.body.removeChild(link)
  URL.revokeObjectURL(link.href)
  
  ElMessage.success('日志下载成功')
}

// 删除Pod
const deletePod = async (pod: any) => {
  try {
    // 这里应该调用真实的API删除Pod
    // 暂时使用模拟数据
    const index = pods.value.findIndex((p: any) => p.namespace === pod.namespace && p.name === pod.name)
    if (index !== -1) {
      pods.value.splice(index, 1)
    }
    ElMessage.success('删除Pod成功')
  } catch (error) {
    console.error('删除Pod失败:', error)
    ElMessage.error('删除Pod失败')
  }
}

// 处理下拉框命令
const handleDropdownCommand = (command: string, pod: any) => {
  switch (command) {
    case 'logs':
      viewPodLogs(pod)
      break
    case 'terminal':
      viewPodTerminal(pod)
      break
    case 'delete':
      deletePod(pod)
      break
    default:
      break
  }
}

// 查看Pod终端
const viewPodTerminal = (pod: any) => {
  selectedPod.value = pod
  selectedContainer.value = pod.containers[0]?.name || ''
  podTerminalVisible.value = true
  connectTerminal()
}

// 连接终端
const connectTerminal = () => {
  if (!selectedPod.value || !selectedCluster.value || !selectedContainer.value) {
    return
  }
  
  terminalLoading.value = true
  terminalContent.value = '正在连接到终端...\n'
  
  try {
    // 构建WebSocket URL
    if (!selectedCluster.value || !selectedPod.value) {
      terminalLoading.value = false
      terminalContent.value += '连接失败：缺少集群或Pod信息\n'
      ws = null
      scrollToBottom()
      return
    }
    const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const wsURL = `${wsProtocol}//${window.location.host}/api/clusters/${selectedCluster.value}/pods/${selectedPod.value.name}/exec`
    
    // 创建WebSocket连接
    ws = new WebSocket(wsURL, ['terminal'])
    
    // 发送初始数据
    ws.onopen = () => {
      terminalLoading.value = false
      terminalContent.value += '连接成功！\n'
      if (selectedPod.value) {
        terminalContent.value += `${selectedPod.value.name}:~$ `
        scrollToBottom()
        
        // 发送容器信息
        ws?.send(JSON.stringify({
          namespace: selectedPod.value.namespace,
          container: selectedContainer.value
        }))
      }
    }
    
    // 接收消息
    ws.onmessage = (event) => {
      terminalContent.value += event.data
      scrollToBottom()
    }
    
    // 连接关闭
    ws.onclose = () => {
      if (terminalLoading.value) {
        terminalLoading.value = false
        terminalContent.value += '连接失败！\n'
      } else {
        terminalContent.value += '\n连接已关闭\n'
      }
      ws = null
      scrollToBottom()
    }
    
    // 连接错误
    ws.onerror = (error) => {
      terminalLoading.value = false
      const errorObj = error as any
      terminalContent.value += `连接错误: ${errorObj.type || 'WebSocket Error'} - ${errorObj.message || '未知错误'}\n`
      ws = null
      scrollToBottom()
    }
    
    // 连接关闭
    ws.onclose = (event) => {
      if (terminalLoading.value) {
        terminalLoading.value = false
        terminalContent.value += `连接失败！原因: ${event.code} - ${event.reason}\n`
      } else {
        terminalContent.value += '\n连接已关闭\n'
      }
      ws = null
      scrollToBottom()
    }
  } catch (error) {
    terminalLoading.value = false
    terminalContent.value += `连接失败: ${error}\n`
    ws = null
    scrollToBottom()
  }
}

// 发送终端命令
const sendTerminalCommand = () => {
  if (!ws || ws.readyState !== WebSocket.OPEN) {
    terminalContent.value += '\n终端未连接\n'
    terminalContent.value += `${selectedPod.value?.name}:~$ `
    scrollToBottom()
    return
  }
  
  const command = terminalInput.value.trim()
  if (!command) {
    terminalContent.value += '\n'
    terminalContent.value += `${selectedPod.value?.name}:~$ `
    scrollToBottom()
    return
  }
  
  // 保存到历史记录
  terminalHistory.push(command)
  historyIndex = terminalHistory.length
  
  // 发送命令
  ws.send(command + '\n')
  
  // 清空输入
  terminalInput.value = ''
  scrollToBottom()
}

// 重启终端
const restartTerminal = () => {
  closeTerminal()
  connectTerminal()
}

// 关闭终端
const closeTerminal = () => {
  if (ws) {
    ws.close()
    ws = null
  }
  terminalContent.value = ''
  terminalInput.value = ''
  terminalLoading.value = false
}

// 导航历史命令
const navigateHistory = (direction: number) => {
  if (terminalHistory.length === 0) return
  
  historyIndex += direction
  
  // 确保索引在有效范围内
  if (historyIndex < 0) {
    historyIndex = 0
  } else if (historyIndex >= terminalHistory.length) {
    historyIndex = terminalHistory.length
    terminalInput.value = ''
    return
  }
  
  // 设置当前输入为历史命令
  terminalInput.value = terminalHistory[historyIndex]
}

// 滚动到终端底部
const scrollToBottom = () => {
  setTimeout(() => {
    if (terminalContentRef.value) {
      const preElement = terminalContentRef.value.querySelector('pre')
      if (preElement) {
        preElement.scrollTop = preElement.scrollHeight
      }
    }
  }, 50)
}
</script>

<style scoped>
.pods-view {
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

.pods-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pod-details {
  padding: 10px 0;
}

.pod-logs {
  padding: 10px 0;
}

.logs-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
  flex-wrap: wrap;
  gap: 10px;
}

.logs-controls {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.logs-content {
  height: 400px;
  overflow-y: auto;
  background-color: #f5f7fa;
  padding: 10px;
  border-radius: 4px;
}

.logs-content pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
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

/* 终端样式 */
.pod-terminal {
  padding: 10px 0;
}

.terminal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
  flex-wrap: wrap;
  gap: 10px;
}

.terminal-controls {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.terminal-content {
  height: 400px;
  overflow-y: auto;
  background-color: #1e1e1e;
  color: #d4d4d4;
  padding: 15px;
  border-radius: 4px;
  font-family: 'Courier New', Courier, monospace;
  font-size: 14px;
  line-height: 1.5;
  margin-bottom: 15px;
}

.terminal-content pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.terminal-input-area {
  display: flex;
  align-items: center;
  background-color: #1e1e1e;
  padding: 10px 15px;
  border-radius: 4px;
  font-family: 'Courier New', Courier, monospace;
}

.terminal-prompt {
  color: #4ec9b0;
  font-weight: bold;
  margin-right: 5px;
}

.terminal-content :deep(.el-loading-mask) {
  background-color: rgba(30, 30, 30, 0.8);
}

.terminal-content :deep(.el-loading-spinner .el-loading-text) {
  color: #d4d4d4;
}

.terminal-content :deep(.el-select),
.terminal-content :deep(.el-button) {
  margin-bottom: 10px;
}
</style>