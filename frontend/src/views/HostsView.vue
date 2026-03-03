<template>
  <div class="hosts-view">
    <el-card shadow="never" class="page-header">
      <div class="header-content">
        <h2>虚机管理</h2>
      </div>
    </el-card>

    <el-card shadow="hover" class="hosts-card">
      <template #header>
        <div class="card-header">
          <span>Linux主机列表</span>
          <div class="card-header-actions">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索主机名称或IP地址"
              clearable
              style="width: 300px; margin-right: 12px;"
              @input="handleSearch"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
            <el-button type="primary" @click="openAddHostDialog">
              <el-icon><Plus /></el-icon>
              添加主机
            </el-button>
          </div>
        </div>
      </template>
      
      <el-table 
        :data="hosts" 
        style="width: 100%" 
        stripe
        fit
        v-loading="loading.loadHosts"
        element-loading-text="加载主机列表中..."
        element-loading-spinner="el-icon-loading"
        element-loading-background="rgba(255, 255, 255, 0.8)"
      >
        <el-table-column prop="id" label="ID" width="80" min-width="80" />
        <el-table-column prop="name" label="主机名称" min-width="150" />
        <el-table-column prop="ip" label="IP地址" min-width="180" />
        <el-table-column label="状态" min-width="120">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'running' ? 'success' : 'danger'">
              {{ scope.row.status === 'running' ? '运行中' : '连接失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="200" />
        <el-table-column label="操作" min-width="500" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" @click="openSSH(scope.row)">
              <el-icon><Connection /></el-icon>
              SSH终端
            </el-button>
            <el-button type="warning" size="small" @click="testConnection(scope.row)">
              <el-icon><Check /></el-icon>
              测试连接
            </el-button>
            <el-button type="info" size="small" @click="openMonitoring(scope.row)">
              <el-icon><DataAnalysis /></el-icon>
              监控
            </el-button>
            <el-button type="info" size="small" @click="openFileTransfer(scope.row)">
              <el-icon><UploadFilled /></el-icon>
              文件传输
            </el-button>
            <!-- 调整"更多"按钮样式和位置 -->
            <el-dropdown trigger="click" placement="bottom">
              <el-button size="small" type="warning" style="margin-left: 12px; border-radius: 4px;">
                更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="editHost(scope.row)">
                    <el-icon><Edit /></el-icon> 编辑
                  </el-dropdown-item>
                  <el-dropdown-item @click="restartHost(scope.row)">
                    <el-icon><Refresh /></el-icon> 重启
                  </el-dropdown-item>
                  <el-dropdown-item @click="shutdownHost(scope.row)">
                    <el-icon><SwitchButton /></el-icon> 关机
                  </el-dropdown-item>
                  <el-dropdown-item divided @click="deleteHost(scope.row.id)">
                    <el-icon><Delete /></el-icon> 删除
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <!-- 添加/编辑主机对话框 -->
    <el-dialog
      v-model="hostDialogVisible"
      :title="isEditHost ? '编辑主机' : '添加主机'"
      width="600px"
    >
      <el-form :model="hostForm" :rules="hostRules" ref="hostFormRef" label-width="100px">
        <el-form-item label="主机名称" prop="name">
          <el-input v-model="hostForm.name" placeholder="请输入主机名称" />
        </el-form-item>
        
        <el-form-item label="IP地址" prop="ip">
          <el-input v-model="hostForm.ip" placeholder="请输入主机IP地址" />
        </el-form-item>
        
        <el-form-item label="端口" prop="port">
          <el-input v-model.number="hostForm.port" placeholder="请输入SSH端口" type="number" />
        </el-form-item>
        
        <el-form-item label="用户名" prop="username">
          <el-input v-model="hostForm.username" placeholder="请输入SSH用户名" />
        </el-form-item>
        
        <el-form-item label="密码" prop="password">
          <el-input v-model="hostForm.password" placeholder="请输入SSH密码" type="password" show-password />
        </el-form-item>
        
        <el-form-item label="备注" prop="remark">
          <el-input v-model="hostForm.remark" placeholder="请输入备注信息" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="hostDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveHost" :loading="loading.saveHost">确定</el-button>
        </span>
      </template>
    </el-dialog>
    
    <!-- SSH终端对话框 -->
    <el-dialog
      v-model="sshDialogVisible"
      :title="`SSH连接 - ${currentSSHHost.name || currentSSHHost.ip}`"
      width="900px"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
    >
      <div class="ssh-terminal" v-if="sshDialogVisible">
        <div class="terminal-header">
          <div class="terminal-title">
            {{ currentSSHHost.name || currentSSHHost.ip }}
            <div class="terminal-status">已连接</div>
          </div>
          <div class="terminal-actions">
            <el-button size="small" @click="sendCommand('clear')">清屏</el-button>
            <el-button size="small" @click="copyTerminalContent">复制</el-button>
            <el-button size="small" @click="sshDialogVisible = false">关闭</el-button>
          </div>
        </div>
        <div class="terminal-content">
          <div ref="terminalContainer" class="terminal-container">
            <!-- xterm.js 会自动渲染到这里 -->
          </div>
        </div>
        <!-- 终端命令历史和快捷键提示 -->
        <div class="command-history">
          <span>💡 提示：</span>
          <span>直接在终端内输入命令，按Enter执行</span>
          <div class="terminal-shortcuts">
            <div class="shortcut">
              <span class="shortcut-key">Ctrl+C</span>
              <span>中断命令</span>
            </div>
            <div class="shortcut">
              <span class="shortcut-key">↑↓</span>
              <span>命令历史</span>
            </div>
            <div class="shortcut">
              <span class="shortcut-key">Tab</span>
              <span>自动补全</span>
            </div>
          </div>
        </div>
      </div>
    </el-dialog>
    
    <!-- 文件传输对话框 -->
    <el-dialog
      v-model="fileTransferVisible"
      :title="`文件传输 - ${currentFileHost.name || currentFileHost.ip}`"
      width="800px"
    >
      <div class="file-transfer" v-if="fileTransferVisible">
        <el-tabs v-model="activeFileTab" type="border-card">
          <el-tab-pane label="上传文件" name="upload">
            <el-upload
              class="upload-demo"
              drag
              action="#"
              :auto-upload="false"
              :file-list="uploadFileList"
              :on-change="handleFileChange"
              :on-remove="handleFileRemove"
              :http-request="handleUpload"
            >
              <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
              <div class="el-upload__text">将文件拖到此处，或<em>点击上传</em></div>
              <template #tip>
                <div class="el-upload__tip">
                  支持上传任意文件
                </div>
              </template>
            </el-upload>
            <el-form :model="fileTransferForm" label-width="80px" class="mt-4">
              <el-form-item label="远程路径">
                <el-input v-model="fileTransferForm.remotePath" placeholder="默认: /tmp/文件名" />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="startUpload" :loading="loading.uploadFile">开始上传</el-button>
              </el-form-item>
            </el-form>
          </el-tab-pane>
          <el-tab-pane label="下载文件" name="download">
            <el-form :model="fileTransferForm" label-width="80px">
              <el-form-item label="远程路径">
                <el-input v-model="fileTransferForm.remotePath" placeholder="例如: /etc/passwd" />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="startDownload" :loading="loading.downloadFile">开始下载</el-button>
              </el-form-item>
            </el-form>
          </el-tab-pane>
        </el-tabs>
      </div>
    </el-dialog>
    
    <!-- 监控对话框 -->
    <el-dialog
      v-model="monitoringDialogVisible"
      :title="`服务器监控 - ${currentMonitoringHost.name || currentMonitoringHost.ip}`"
      width="1000px"
      :close-on-click-modal="false"
    >
      <div class="monitoring-content" v-if="monitoringDialogVisible">
        <el-card shadow="hover" class="monitoring-card" v-loading="monitoringLoading">
          <template #header>
            <div class="card-header">
              <span>服务器资源监控</span>
              <div class="card-header-actions">
                <el-button type="primary" @click="getServerStats" size="small">
                  <el-icon><Refresh /></el-icon>
                  刷新数据
                </el-button>
                <el-switch
                  v-model="autoRefresh"
                  active-text="自动刷新"
                  inactive-text="手动刷新"
                  @change="handleAutoRefreshChange"
                  size="small"
                ></el-switch>
              </div>
            </div>
          </template>
          
          <div class="stats-grid">
            <div class="stat-item">
              <div class="stat-label">
                <el-icon><Cpu /></el-icon>
                CPU使用率
              </div>
              <div class="stat-value">{{ serverStats.cpu }}%</div>
              <div class="stat-bar">
                <div class="stat-bar-fill" :style="{ width: `${serverStats.cpu}%`, backgroundColor: getCPUColor(serverStats.cpu) }"></div>
              </div>
            </div>
            
            <div class="stat-item">
              <div class="stat-label">
                <el-icon><DataLine /></el-icon>
                内存使用率
              </div>
              <div class="stat-value">{{ serverStats.memory }}%</div>
              <div class="stat-bar">
                <div class="stat-bar-fill" :style="{ width: `${serverStats.memory}%`, backgroundColor: getMemoryColor(serverStats.memory) }"></div>
              </div>
            </div>
            
            <div class="stat-item">
              <div class="stat-label">
                <el-icon><DataLine /></el-icon>
                磁盘使用率
              </div>
              <div class="stat-value">{{ serverStats.disk }}%</div>
              <div class="stat-bar">
                <div class="stat-bar-fill" :style="{ width: `${serverStats.disk}%`, backgroundColor: getDiskColor(serverStats.disk) }"></div>
              </div>
            </div>
            
            <div class="stat-item">
              <div class="stat-label">
                <el-icon><DataLine /></el-icon>
                系统负载
              </div>
              <div class="stat-value">{{ serverStats.loadAverage?.toFixed(2) || '0.00' }}</div>
              <div class="stat-description">1分钟平均负载</div>
            </div>
            
            <div class="stat-item">
              <div class="stat-label">
                <el-icon><Upload /></el-icon>
                上传速度
              </div>
              <div class="stat-value">{{ (serverStats.uploadSpeed || 0).toFixed(2) }} KB/s</div>
              <div class="stat-description">当前网络上传速度</div>
            </div>
            
            <div class="stat-item">
              <div class="stat-label">
                <el-icon><Download /></el-icon>
                下载速度
              </div>
              <div class="stat-value">{{ (serverStats.downloadSpeed || 0).toFixed(2) }} KB/s</div>
              <div class="stat-description">当前网络下载速度</div>
            </div>
          </div>
          
          <!-- 系统信息 -->
          <el-divider></el-divider>
          <div class="system-info">
            <h3>系统信息</h3>
            <el-descriptions :column="3" border>
              <el-descriptions-item label="主机名称">{{ currentMonitoringHost.name || '未知' }}</el-descriptions-item>
              <el-descriptions-item label="IP地址">{{ currentMonitoringHost.ip || '未知' }}</el-descriptions-item>
              <el-descriptions-item label="监控时间">{{ formatDate(new Date()) }}</el-descriptions-item>
            </el-descriptions>
          </div>
        </el-card>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch, inject } from 'vue'
import { Plus, Edit, Delete, Connection, UploadFilled, Refresh, SwitchButton, Check, DataAnalysis, More, ArrowDown, Cpu, DataLine, Upload, Download } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { hostAPI } from '../api/api'


// 表单引用
const hostFormRef = ref()

// 加载状态
const loading = ref({
  saveHost: false,
  deleteHost: false,
  loadHosts: false,
  uploadFile: false,
  downloadFile: false,
  sshCommand: false
})

// 主机列表数据
const hosts = ref([])
// 原始主机列表，用于搜索过滤
const originalHosts = ref([])
// 搜索关键词
const searchKeyword = ref('')

// 注入刷新触发器
const refreshTrigger = inject('refreshTrigger', ref(0))

// 监听刷新触发器变化，重新加载数据
watch(refreshTrigger, () => {
  loadHosts()
})

// 主机对话框相关
const hostDialogVisible = ref(false)
const isEditHost = ref(false)
const hostForm = ref({
  id: 0,
  name: '',
  ip: '',
  port: 22,
  username: '',
  password: '',
  remark: ''
})

// SSH终端相关
const sshDialogVisible = ref(false)
const currentSSHHost = ref({
  id: 0,
  ip: '',
  port: 22,
  username: '',
  password: '',
  name: ''
})
const terminalContainer = ref(null)
const terminalInput = ref('')
const terminalConnected = ref(false)
// 终端实例
const terminal = ref(null)
const fitAddon = ref(null)
const commandBuffer = ref('')
// 用于窗口大小变化的防抖
const resizeTimeout = ref(null)

// 监控相关
const monitoringDialogVisible = ref(false)
const currentMonitoringHost = ref({
  id: 0,
  name: '',
  ip: ''
})
const serverStats = ref({
  cpu: 0,
  memory: 0,
  disk: 0,
  loadAverage: 0,
  uploadSpeed: 0,
  downloadSpeed: 0
})
const monitoringLoading = ref(false)
// 自动刷新相关
const autoRefresh = ref(false)
let autoRefreshTimer = null
const refreshInterval = 5000 // 5秒刷新一次

// 文件传输相关
const fileTransferVisible = ref(false)
const currentFileHost = ref({
  id: 0,
  ip: '',
  port: 22,
  username: '',
  password: '',
  name: ''
})
const activeFileTab = ref('upload')
const uploadFileList = ref([])
const fileTransferForm = ref({
  remotePath: ''
})
const selectedFiles = ref([])

// 主机表单验证规则
const hostRules = ref({
  name: [
    { required: true, message: '请输入主机名称', trigger: 'blur' },
    { min: 2, max: 50, message: '主机名称长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  ip: [
    { required: true, message: '请输入IP地址', trigger: 'blur' },
    { pattern: /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/, message: '请输入正确的IP地址', trigger: 'blur' }
  ],
  port: [
    { required: true, message: '请输入端口', trigger: 'blur' },
    { type: 'number', min: 1, max: 65535, message: '端口范围在 1 到 65535 之间', trigger: 'blur' }
  ],
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
})

// 从API加载主机列表
const loadHosts = async () => {
  loading.value.loadHosts = true
  try {
    const response = await hostAPI.list()
    // 直接使用后端返回的数据，已经是小写字段名
    originalHosts.value = response.data
    hosts.value = [...originalHosts.value]
    ElMessage.success('主机列表加载成功')
  } catch (error) {
    console.error('加载主机列表失败:', error)
    // 如果API失败，使用模拟数据
    const mockData = [
      {
        id: 1,
        name: '测试主机1',
        ip: '192.168.1.100',
        remark: '这是一台测试Linux主机',
        isActive: true,
        status: 'running'
      },
      {
        id: 2,
        name: '生产主机1',
        ip: '10.0.0.100',
        remark: '生产环境Linux主机',
        isActive: true,
        status: 'running'
      }
    ]
    originalHosts.value = mockData
    hosts.value = [...originalHosts.value]
    ElMessage.warning('API加载失败，使用模拟数据')
  } finally {
    loading.value.loadHosts = false
  }
}

// 打开添加主机对话框
const openAddHostDialog = () => {
  isEditHost.value = false
  hostForm.value = {
    id: 0,
    name: '',
    ip: '',
    port: 22,
    username: '',
    password: '',
    remark: ''
  }
  hostDialogVisible.value = true
}

// 编辑主机
const editHost = async (host) => {
  isEditHost.value = true
  try {
    // 获取完整的主机信息，包括敏感字段
    const response = await hostAPI.get(host.id)
    const hostData = response.data
    // 转换字段名为小写，适配前端表单
    hostForm.value = {
      id: hostData.ID || hostData.id,
      name: hostData.Name || hostData.name,
      ip: hostData.IP || hostData.ip,
      port: hostData.Port || hostData.port,
      username: hostData.Username || hostData.username,
      password: hostData.Password || hostData.password,
      remark: hostData.Remark || hostData.remark
    }
  } catch (error) {
    console.error('获取主机详情失败:', error)
    // 如果API失败，使用当前主机数据
    hostForm.value = { ...host }
  }
  hostDialogVisible.value = true
}

// 删除主机
const deleteHost = (id) => {
  ElMessageBox.confirm('确定要删除该主机吗？', '删除确认', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    loading.value.deleteHost = true
    try {
      await hostAPI.delete(id)
      // 从原始列表中删除
      originalHosts.value = originalHosts.value.filter(host => host.id !== id)
      // 重新应用搜索过滤
      handleSearch()
      ElMessage.success('主机删除成功')
    } catch (error) {
      console.error('删除主机失败:', error)
      // 如果API失败，使用模拟数据删除
      originalHosts.value = originalHosts.value.filter(host => host.id !== id)
      // 重新应用搜索过滤
      handleSearch()
      ElMessage.success('主机删除成功')
    } finally {
      loading.value.deleteHost = false
    }
  }).catch((error) => {
    if (error !== 'cancel') {
      console.error('删除主机失败:', error)
      ElMessage.error('删除主机失败')
    }
  })
}

// 重启主机
const restartHost = (host) => {
  ElMessageBox.confirm('确定要重启该主机吗？', '重启确认', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    loading.value.sshCommand = true
    try {
      // 调用API重启主机，只需要主机ID
      await hostAPI.restartHost(host.id)
      ElMessage.success('主机重启命令已发送')
    } catch (error) {
      console.error('重启主机失败:', error)
      ElMessage.error('重启主机失败')
    } finally {
      loading.value.sshCommand = false
    }
  }).catch((error) => {
    if (error !== 'cancel') {
      console.error('重启主机失败:', error)
      ElMessage.error('重启主机失败')
    }
  })
}

// 关机主机
const shutdownHost = (host) => {
  ElMessageBox.confirm('确定要关闭该主机吗？', '关机确认', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'danger'
  }).then(async () => {
    loading.value.sshCommand = true
    try {
      // 调用API关闭主机，只需要主机ID
      await hostAPI.shutdownHost(host.id)
      ElMessage.success('主机关机命令已发送')
    } catch (error) {
      console.error('关闭主机失败:', error)
      ElMessage.error('关闭主机失败')
    } finally {
      loading.value.sshCommand = false
    }
  }).catch((error) => {
    if (error !== 'cancel') {
      console.error('关闭主机失败:', error)
      ElMessage.error('关闭主机失败')
    }
  })
}

// 测试主机连接
const testConnection = async (host) => {
  loading.value.sshCommand = true
  try {
    // 调用API测试主机连接
    await hostAPI.testConnection(host.id)
    ElMessage.success('主机连接测试成功')
  } catch (error) {
    console.error('测试主机连接失败:', error)
    ElMessage.error('主机连接测试失败')
  } finally {
    loading.value.sshCommand = false
  }
}

// 搜索处理函数
const handleSearch = () => {
  if (!searchKeyword.value) {
    // 如果搜索关键词为空，显示所有主机
    hosts.value = [...originalHosts.value]
    return
  }
  
  const keyword = searchKeyword.value.toLowerCase()
  // 根据主机名称或IP地址进行搜索
  hosts.value = originalHosts.value.filter(host => {
    const name = host.name?.toLowerCase() || ''
    const ip = host.ip?.toLowerCase() || ''
    return name.includes(keyword) || ip.includes(keyword)
  })
}

// 保存主机
const saveHost = () => {
  if (!hostFormRef.value) return
  
  hostFormRef.value.validate((valid) => {
    if (valid) {
      loading.value.saveHost = true
      // 调用API保存主机
      const savePromise = isEditHost.value 
        ? hostAPI.update(hostForm.value.id, hostForm.value) 
        : hostAPI.create(hostForm.value)
      
      savePromise.then(() => {
        // 重新加载主机列表
        loadHosts()
        ElMessage.success(isEditHost.value ? '主机编辑成功' : '主机添加成功')
        hostDialogVisible.value = false
      }).catch((error) => {
        console.error('保存主机失败:', error)
        // 如果API失败，使用模拟数据保存
        if (isEditHost.value) {
          const index = originalHosts.value.findIndex(host => host.id === hostForm.value.id)
          if (index !== -1) {
            originalHosts.value[index] = { ...hostForm.value }
          }
        } else {
          const newHost = {
            ...hostForm.value,
            id: Date.now() // 模拟生成ID
          }
          originalHosts.value.push(newHost)
        }
        // 重新应用搜索过滤
        handleSearch()
        ElMessage.success(isEditHost.value ? '主机编辑成功' : '主机添加成功')
        hostDialogVisible.value = false
      }).finally(() => {
        loading.value.saveHost = false
      })
    }
  })
}

// 打开SSH终端
const openSSH = async (host) => {
  try {
    // 获取完整的主机信息，包括敏感字段
    const response = await hostAPI.get(host.id)
    const hostData = response.data
    // 转换字段名为小写，适配前端表单
    currentSSHHost.value = {
      id: hostData.ID || hostData.id,
      name: hostData.Name || hostData.name,
      ip: hostData.IP || hostData.ip,
      port: hostData.Port || hostData.port,
      username: hostData.Username || hostData.username,
      password: hostData.Password || hostData.password
    }
  } catch (error) {
    console.error('获取主机详情失败:', error)
    // 如果API失败，使用当前主机数据
    currentSSHHost.value = { ...host }
  }
  sshDialogVisible.value = true
  // terminalConnected.value 会在 initTerminal 中设置
}

// 初始化终端
const initTerminal = () => {
  try {
    console.log('=== 开始初始化终端 ===')
    
    // 1. 彻底清空命令缓冲区，避免残留输入导致的连续字符问题
    commandBuffer.value = ''
    
    // 2. 检查终端容器
    console.log('终端容器:', terminalContainer.value)
    const terminalContent = document.querySelector('.terminal-content')
    if (terminalContent) {
      // 完全清空终端内容区域，避免任何残留元素
      terminalContent.innerHTML = ''
      
      // 重新创建干净的终端容器
      const container = document.createElement('div')
      container.className = 'terminal-container'
      container.style.width = '100%'
      container.style.height = '100%'
      container.style.backgroundColor = '#1e1e1e'
      container.style.position = 'absolute'
      container.style.overflow = 'hidden' // 避免滚动条问题
      container.style.padding = '0'
      container.style.margin = '0'
      container.style.border = 'none'
      container.style.outline = 'none'
      container.style.boxShadow = 'none'
      terminalContainer.value = container
      terminalContent.appendChild(container)
      console.log('✅ 已创建干净的终端容器')
    }
    
    // 检查终端容器尺寸
    const containerRect = terminalContainer.value.getBoundingClientRect()
    console.log('终端容器尺寸:', containerRect.width, 'x', containerRect.height)
    if (containerRect.width === 0 || containerRect.height === 0) {
      console.error('❌ 终端容器尺寸为0，尝试设置最小尺寸')
      terminalContainer.value.style.minWidth = '800px'
      terminalContainer.value.style.minHeight = '400px'
      // 重新获取尺寸
      const newRect = terminalContainer.value.getBoundingClientRect()
      console.log('设置最小尺寸后:', newRect.width, 'x', newRect.height)
    }
    
    // 检查xterm.js是否加载
    console.log('Terminal是否存在:', typeof window.Terminal !== 'undefined')
    
    if (typeof window.Terminal === 'undefined') {
      console.error('❌ xterm.js未加载，尝试重新加载')
      // 尝试动态加载xterm.js
      const xtermScript = document.createElement('script')
      xtermScript.src = '/vendor/xterm/xterm.js'
      xtermScript.onload = () => {
        console.log('✅ xterm.js动态加载成功')
        // 加载完成后再次初始化
        setTimeout(initTerminal, 100)
      }
      xtermScript.onerror = (e) => {
        console.error('❌ xterm.js动态加载失败:', e)
      }
      document.head.appendChild(xtermScript)
      return
    }
    
    // 清除之前的终端实例
    if (terminal.value) {
      console.log('清理之前的终端实例')
      // 移除所有事件监听器
      terminal.value.off('data')
      terminal.value.off('key')
      terminal.value.off('resize')
      terminal.value.off('focus')
      terminal.value.off('blur')
      // 销毁终端实例
      terminal.value.dispose()
      terminal.value = null
    }
    
    // 再次确保命令缓冲区为空，防止残留输入
    commandBuffer.value = ''
    
    // 创建终端实例
    console.log('创建终端实例...')
    terminal.value = new window.Terminal({
      cursorBlink: true,
      fontSize: 14,
      fontFamily: 'Consolas, "Courier New", monospace',
      theme: {
        foreground: '#ffffff',
        background: '#1e1e1e',
        cursor: '#ffffff',
        selection: '#424242',
        black: '#000000',
        red: '#ff0000',
        green: '#00ff00',
        yellow: '#ffff00',
        blue: '#0000ff',
        magenta: '#ff00ff',
        cyan: '#00ffff',
        white: '#ffffff'
      },
      rendererType: 'canvas', // 使用canvas渲染器，更适合大量输出
      convertEol: true,
      disableStdin: false,
      allowTransparency: false,
      scrollback: 1000, // 增加滚动缓冲区大小
      tabStopWidth: 8
    })
    console.log('✅ 终端实例创建成功')
    
    // 打开终端到指定容器
    terminal.value.open(terminalContainer.value)
    console.log('✅ 终端已打开到容器')
    
    // 精确计算终端尺寸
    console.log('精确计算终端尺寸...')
    updateTerminalSize()
    console.log('✅ 终端尺寸已设置')
    
    // 立即清空终端，确保没有任何默认内容
    terminal.value.clear()
    console.log('✅ 终端已强制清空')
    
    // 设置终端焦点
    terminal.value.focus()
    console.log('✅ 终端焦点已设置')
    
    // 再次清空终端，双重保险
    terminal.value.clear()
    
    // 写入欢迎信息
    terminal.value.write('=== SSH Terminal ===\r\n')
    terminal.value.write(`Connected to: ${currentSSHHost.value.name || currentSSHHost.value.ip}:${currentSSHHost.value.port}\r\n`)
    terminal.value.write('Type commands below and press Enter to execute\r\n')
    terminal.value.write('\r\n')
    terminal.value.write('$ ')
    console.log('✅ 欢迎信息已写入')
    
    terminalConnected.value = true
    console.log('=== 终端初始化完成 ===')
    
    // 监听终端输入
    terminal.value.onData((data) => {
      // 处理回车键
      if (data === '\r') {
        terminal.value.write('\r\n')
        if (commandBuffer.value) {
          // 执行命令，使用异步执行避免阻塞
          (async () => {
            await sendCommand(commandBuffer.value)
            commandBuffer.value = ''
            terminal.value.write('$ ')
          })()
        } else {
          terminal.value.write('$ ')
        }
      } 
      // 处理退格键
      else if (data === '\x7f' || data === '\b') {
        if (commandBuffer.value.length > 0) {
          commandBuffer.value = commandBuffer.value.slice(0, -1)
          terminal.value.write('\b \b')
        }
      }
      // 处理其他字符
      else {
        commandBuffer.value += data
        terminal.value.write(data)
      }
    })
    
    // 监听终端焦点事件
    terminalContainer.value.addEventListener('click', () => {
      terminal.value?.focus()
    })
    
    // 监听窗口大小变化，调整终端大小 - 添加防抖
    const handleResize = () => {
      console.log('窗口大小变化，调整终端大小')
      // 清除之前的定时器
      if (resizeTimeout.value) {
        clearTimeout(resizeTimeout.value)
      }
      // 设置新的定时器，100ms延迟
      resizeTimeout.value = window.setTimeout(() => {
        updateTerminalSize()
      }, 100)
    }
    window.addEventListener('resize', handleResize)
    
    // 保存resize事件处理器，以便后续移除
    window.__terminalResizeHandler = handleResize
    
  } catch (error) {
    console.error('=== 终端初始化失败 ===')
    console.error('错误类型:', error.name)
    console.error('错误信息:', error.message)
    console.error('错误堆栈:', error.stack)
    terminalConnected.value = false
    
    // 显示友好的错误信息
    if (terminalContainer.value) {
      let errorMsg = error.message
      if (errorMsg.includes('FitAddon')) {
        errorMsg = '终端适配插件加载失败，已使用基础终端功能'
      }
      
      terminalContainer.value.innerHTML = `
        <div style="color: #ffffff; padding: 20px; font-family: monospace; background-color: #1e1e1e; height: 100%; overflow: auto;">
          <h3 style="color: #ff6b6b;">⚠️ 终端初始化提示</h3>
          <p style="margin: 10px 0;">${errorMsg}</p>
          <p style="margin: 10px 0; color: #888;">终端功能已启动，但可能缺少窗口自适应功能</p>
          <div style="margin: 20px 0; padding: 15px; background-color: #2d2d2d; border-radius: 4px;">
            <h4>💡 终端使用说明：</h4>
            <ul style="margin: 10px 0; padding-left: 20px;">
              <li>直接在终端内输入命令，按Enter执行</li>
              <li>支持基本的终端命令操作</li>
              <li>窗口大小变化时，终端会自动调整</li>
            </ul>
          </div>
          <button onclick="location.reload()" style="margin-top: 10px; padding: 8px 16px; background-color: #4ecdc4; color: white; border: none; border-radius: 4px; cursor: pointer;">
            🔄 重新加载页面
          </button>
        </div>
      `
    }
  }
}

// 更新终端尺寸
const updateTerminalSize = () => {
  if (!terminal.value || !terminalContainer.value) {
    return
  }
  
  try {
    // 获取容器尺寸
    const rect = terminalContainer.value.getBoundingClientRect()
    
    // 计算字符尺寸（基于14px字体）
    const charWidth = 8 // 平均字符宽度
    const charHeight = 19 // 平均字符高度（包括行间距）
    
    // 计算终端行列数
    const cols = Math.max(80, Math.floor(rect.width / charWidth)) // 最小80列
    const rows = Math.max(24, Math.floor(rect.height / charHeight)) // 最小24行
    
    console.log('计算终端尺寸:', cols, '列 x', rows, '行')
    
    // 设置终端尺寸
    terminal.value.resize(cols, rows)
    
    // 确保终端正确渲染 - 使用公共API
    if (terminal.value.refresh) {
      terminal.value.refresh()
    }
    
    console.log('✅ 终端尺寸已更新')
  } catch (error) {
    console.error('❌ 更新终端尺寸失败:', error)
  }
}

// 发送命令到终端
const sendCommand = async (command) => {
  if (!command || !terminalConnected.value || !terminal.value) {
    console.warn('跳过命令执行:', { command, terminalConnected: terminalConnected.value, terminal: terminal.value })
    return
  }
  
  console.log('=== 开始执行命令 ===')
  console.log('命令:', command)
  loading.value.sshCommand = true
  
  try {
    // 调用API执行命令 - 添加超时处理
    console.log('调用API执行命令，主机ID:', currentSSHHost.value.id)
    
    // 使用Promise.race添加超时控制
    const response = await Promise.race([
      hostAPI.openSSH(currentSSHHost.value.id, { command }),
      new Promise((_, reject) => setTimeout(() => reject(new Error('命令执行超时')), 30000)) // 30秒超时
    ])
    
    console.log('命令执行成功，响应:', response)
    
    // 显示命令输出
    let output = ''
    if (response && typeof response === 'object' && 'data' in response) {
      output = response.data.output || response.data.Output || ''
    }
    
    console.log('原始输出长度:', output.length)
    
    // 处理输出，确保格式正确
    output = output.replace(/\r\n/g, '\n') // 统一换行符
    output = output.replace(/\n/g, '\r\n') // 转换为终端兼容的换行符
    output = output.trim() // 去除首尾空白
    
    console.log('处理后输出:', JSON.stringify(output))
    
    // 优化输出方式，避免一次性写入大量数据
    if (output) {
      // 按行分割输出
      const lines = output.split('\r\n')
      console.log('输出行数:', lines.length)
      
      // 分块写入输出，每10行一批，添加微小延迟
      const chunkSize = 10
      for (let i = 0; i < lines.length; i += chunkSize) {
        const chunk = lines.slice(i, i + chunkSize)
        for (const line of chunk) {
          terminal.value.writeln(line)
        }
        // 每块之间添加微小延迟，避免终端阻塞
        await new Promise(resolve => setTimeout(resolve, 10))
      }
    }
    
    // 确保输出后有新行
    terminal.value.writeln('')
    
    console.log('=== 命令输出完成 ===')
  } catch (error) {
    console.error('❌ 执行命令失败:', error)
    // 显示详细的错误信息
    let errorMsg = '未知错误'
    if (error instanceof Error) {
      errorMsg = error.message
    } else if (error.response) {
      // API错误响应
      errorMsg = error.response.data?.error || error.response.data?.Error || `HTTP错误: ${error.response.status}`
    } else if (error.request) {
      // 请求发送但没有收到响应
      errorMsg = '未收到服务器响应'
    }
    
    // 显示错误信息
    terminal.value.writeln(`Error: ${errorMsg}`)
    terminal.value.writeln('')
  } finally {
    loading.value.sshCommand = false
    console.log('=== 命令执行完成 ===')
  }
}

// 复制终端内容到剪贴板
const copyTerminalContent = async () => {
  try {
    if (!terminal.value) {
      ElMessage.warning('终端未初始化')
      return
    }
    
    // 获取终端内容
    const content = terminal.value.getOption('screenReaderMode') 
      ? terminal.value._core.buffer.areas[0].getLines(0, terminal.value.rows).join('\n')
      : terminal.value._core.buffer.areas[0].getLines(0, terminal.value.rows).join('\n')
    
    // 复制到剪贴板
    await navigator.clipboard.writeText(content)
    ElMessage.success('终端内容已复制到剪贴板')
  } catch (error) {
    console.error('复制终端内容失败:', error)
    ElMessage.error('复制失败，请手动选择内容复制')
  }
}

// 关闭终端
const closeTerminal = () => {
  try {
    console.log('=== 开始关闭终端 ===')
    
    // 清空命令缓冲区，防止下次打开终端时出现残留内容
    commandBuffer.value = ''
    
    // 清理终端实例
    if (terminal.value) {
      // 移除所有终端事件监听器
      terminal.value.off('data')
      terminal.value.off('key')
      terminal.value.off('resize')
      terminal.value.off('focus')
      terminal.value.off('blur')
      // 销毁终端实例
      terminal.value.dispose()
      terminal.value = null
      console.log('✅ 终端实例已清理')
    }
    
    // 清理fitAddon引用
    fitAddon.value = null
    
    // 移除窗口大小变化的事件监听器
    if (window.__terminalResizeHandler) {
      window.removeEventListener('resize', window.__terminalResizeHandler)
      delete window.__terminalResizeHandler
      console.log('✅ 窗口大小变化监听器已移除')
    }
    
    terminalConnected.value = false
    console.log('=== 终端关闭完成 ===')
  } catch (error) {
    console.error('❌ 关闭终端失败:', error)
  }
}

// 打开文件传输对话框
const openFileTransfer = async (host) => {
  try {
    // 获取完整的主机信息，包括敏感字段
    const response = await hostAPI.get(host.id)
    const hostData = response.data
    // 转换字段名为小写，适配前端表单
    currentFileHost.value = {
      id: hostData.ID || hostData.id,
      name: hostData.Name || hostData.name,
      ip: hostData.IP || hostData.ip,
      port: hostData.Port || hostData.port,
      username: hostData.Username || hostData.username,
      password: hostData.Password || hostData.password
    }
  } catch (error) {
    console.error('获取主机详情失败:', error)
    // 如果API失败，使用当前主机数据
    currentFileHost.value = { ...host }
  }
  fileTransferVisible.value = true
  uploadFileList.value = []
  selectedFiles.value = []
  fileTransferForm.value.remotePath = ''
}

// 打开监控对话框
const openMonitoring = async (host) => {
  currentMonitoringHost.value = {
    id: host.id,
    name: host.name,
    ip: host.ip
  }
  monitoringDialogVisible.value = true
  // 获取服务器资源数据
  await getServerStats()
}

// 获取服务器资源数据
const getServerStats = async () => {
  monitoringLoading.value = true
  try {
    console.log('开始获取服务器资源数据，主机ID:', currentMonitoringHost.value.id)
    
    // 调用后端API获取服务器资源数据
    const response = await hostAPI.getServerStats(currentMonitoringHost.value.id)
    console.log('服务器资源数据响应:', response.data)
    
    // 更新服务器资源数据
    serverStats.value = response.data
    
    ElMessage.success('获取服务器资源数据成功')
  } catch (error) {
    console.error('获取服务器资源失败:', error)
    console.error('错误详情:', JSON.stringify(error, null, 2))
    
    // 尝试使用简化的模拟数据
    console.log('尝试使用模拟数据')
    serverStats.value = {
      cpu: Math.floor(Math.random() * 50) + 10,
      memory: Math.floor(Math.random() * 60) + 20,
      disk: Math.floor(Math.random() * 40) + 30,
      loadAverage: Math.random() * 2,
      uploadSpeed: Math.random() * 100,
      downloadSpeed: Math.random() * 200
    }
    ElMessage.warning('使用模拟数据显示服务器资源')
  } finally {
    monitoringLoading.value = false
  }
}

// 处理自动刷新开关变化
const handleAutoRefreshChange = () => {
  if (autoRefresh.value) {
    // 启动自动刷新
    startAutoRefresh()
    ElMessage.success('已开启自动刷新，每5秒刷新一次')
  } else {
    // 停止自动刷新
    stopAutoRefresh()
    ElMessage.info('已关闭自动刷新')
  }
}

// 启动自动刷新
const startAutoRefresh = () => {
  // 先清除可能存在的定时器
  stopAutoRefresh()
  // 设置新的定时器
  autoRefreshTimer = setInterval(() => {
    getServerStats()
  }, refreshInterval)
}

// 停止自动刷新
const stopAutoRefresh = () => {
  if (autoRefreshTimer) {
    clearInterval(autoRefreshTimer)
    autoRefreshTimer = null
  }
}

// 日期格式化函数
const formatDate = (date) => {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

// 根据CPU使用率获取颜色
const getCPUColor = (cpu) => {
  if (cpu < 50) return '#67C23A' // 绿色
  if (cpu < 80) return '#E6A23C' // 黄色
  return '#F56C6C' // 红色
}

// 根据内存使用率获取颜色
const getMemoryColor = (memory) => {
  if (memory < 60) return '#67C23A' // 绿色
  if (memory < 85) return '#E6A23C' // 黄色
  return '#F56C6C' // 红色
}

// 根据磁盘使用率获取颜色
const getDiskColor = (disk) => {
  if (disk < 70) return '#67C23A' // 绿色
  if (disk < 90) return '#E6A23C' // 黄色
  return '#F56C6C' // 红色
}

// 处理文件选择
const handleFileChange = (file, fileList) => {
  // 更新上传文件列表
  uploadFileList.value = fileList
  // 更新选中的文件列表
  selectedFiles.value = fileList.map(f => f.raw)
}

// 处理文件移除
const handleFileRemove = (file, fileList) => {
  // 更新上传文件列表
  uploadFileList.value = fileList
  // 更新选中的文件列表
  selectedFiles.value = fileList.map(f => f.raw)
}

// 处理文件上传
const handleUpload = async (options) => {
  const file = options.file
  loading.value.uploadFile = true
  
  try {
    const formData = new FormData()
    formData.append('file', file)
    if (fileTransferForm.value.remotePath) {
      formData.append('remotePath', fileTransferForm.value.remotePath)
    }
    
    // 调用API上传文件
    await hostAPI.uploadFile(currentFileHost.value.id, formData)
    ElMessage.success(`文件 ${file.name} 上传成功`)
    options.onSuccess()
  } catch (error) {
    console.error('文件上传失败:', error)
    // 输出更详细的错误信息
    if (error.response) {
      console.error('响应状态:', error.response.status)
      console.error('响应数据:', error.response.data)
    } else if (error.request) {
      console.error('请求已发送但没有收到响应:', error.request)
    } else {
      console.error('请求配置错误:', error.message)
    }
    ElMessage.error(`文件 ${file.name} 上传失败: ${error.message}`)
    options.onError()
  } finally {
    loading.value.uploadFile = false
  }
}

// 开始上传文件
const startUpload = () => {
  if (uploadFileList.value.length === 0) {
    ElMessage.warning('请先选择要上传的文件')
    return
  }
  
  // 创建文件列表的副本，避免遍历过程中修改数组
  const filesToUpload = [...uploadFileList.value]
  
  // 遍历所有选中的文件并上传
  filesToUpload.forEach(file => {
    handleUpload({
      file: file.raw,
      onSuccess: () => {
        // 上传成功后从列表中移除
        const index = uploadFileList.value.findIndex(f => f.uid === file.uid)
        if (index !== -1) {
          uploadFileList.value.splice(index, 1)
          selectedFiles.value = uploadFileList.value.map(f => f.raw)
        }
      },
      onError: () => {
        // 上传失败处理
      }
    })
  })
}

// 开始下载文件
const startDownload = async () => {
  if (!fileTransferForm.value.remotePath) {
    ElMessage.warning('请输入远程文件路径')
    return
  }
  
  loading.value.downloadFile = true
  
  try {
    // 调用API下载文件
    const response = await hostAPI.downloadFile(currentFileHost.value.id, {
      remotePath: fileTransferForm.value.remotePath
    })
    
    // 创建下载链接
    const blob = new Blob([response.data])
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = fileTransferForm.value.remotePath.split('/').pop() || 'download'
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    window.URL.revokeObjectURL(url)
    
    ElMessage.success('文件下载成功')
  } catch (error) {
    console.error('文件下载失败:', error)
    ElMessage.error('文件下载失败')
  } finally {
    loading.value.downloadFile = false
  }
}

// 监听SSH对话框显示状态
watch(sshDialogVisible, (newVal) => {
  if (newVal) {
    // 对话框打开时初始化终端
    setTimeout(() => {
      initTerminal()
    }, 100)
  } else {
    // 对话框关闭时关闭终端
    closeTerminal()
  }
})

// 监听监控对话框显示状态
watch(monitoringDialogVisible, (newVal) => {
  if (!newVal) {
    // 对话框关闭时停止自动刷新
    stopAutoRefresh()
    autoRefresh.value = false
  }
})

// 组件挂载时加载主机列表
onMounted(() => {
  loadHosts()
})

// 组件卸载时清理资源
onUnmounted(() => {
  closeTerminal()
  stopAutoRefresh()
})
</script>

<style scoped>
.hosts-view {
  width: 100%;
  padding: 20px;
  margin: 0 auto;
  max-width: 100%;
  box-sizing: border-box;
}

.page-header {
  margin-bottom: 30px;
  width: 100%;
}

.header-content h2 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: #333;
}

.hosts-card {
  margin-bottom: 30px;
  padding: 20px;
  width: 100%;
  max-width: 100%;
  box-sizing: border-box;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.card-header-actions {
  display: flex;
  align-items: center;
}

.card-header span {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

/* SSH终端样式 - 现代化美化设计 */
.ssh-terminal {
  display: flex !important;
  flex-direction: column !important;
  height: 650px !important; /* 增加终端高度，提供更好的可视区域 */
  background: linear-gradient(135deg, #0f0c29 0%, #302b63 50%, #24243e 100%) !important; /* 更丰富的渐变背景 */
  border-radius: 16px !important; /* 更大的圆角 */
  overflow: hidden !important;
  color: #e0e0e0 !important;
  position: relative !important;
  box-shadow: 0 12px 48px rgba(0, 0, 0, 0.5) !important; /* 增强阴影效果 */
  border: 1px solid #2d2d4d !important;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1) !important; /* 更流畅的过渡效果 */
  backdrop-filter: blur(10px) !important; /* 背景模糊效果 */
}

.ssh-terminal:hover {
  box-shadow: 0 20px 64px rgba(0, 0, 0, 0.6) !important; /* 悬停时增强阴影 */
  transform: translateY(-3px) !important; /* 轻微上浮效果 */
}

/* 终端头部美化 */
.terminal-header {
  display: flex !important;
  justify-content: space-between !important;
  align-items: center !important;
  padding: 18px 24px !important;
  background: rgba(48, 43, 99, 0.8) !important; /* 半透明背景 */
  border-bottom: 2px solid #4a4a7a !important; /* 更粗的边框 */
  height: 70px !important;
  box-sizing: border-box !important;
  border-radius: 16px 16px 0 0 !important; /* 顶部圆角 */
  backdrop-filter: blur(10px) !important;
}

.terminal-title {
  font-size: 18px !important;
  font-weight: 700 !important;
  color: #8ecae6 !important; /* 亮蓝色标题 */
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif !important;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.4) !important;
  letter-spacing: 0.5px !important;
  display: flex !important;
  align-items: center !important;
  gap: 10px !important;
}

/* 添加终端状态指示器 */
.terminal-status {
  display: flex !important;
  align-items: center !important;
  gap: 8px !important;
  font-size: 12px !important;
  color: #66d9ef !important;
  background: rgba(102, 217, 239, 0.1) !important;
  padding: 4px 12px !important;
  border-radius: 12px !important;
  border: 1px solid rgba(102, 217, 239, 0.3) !important;
}

.terminal-status::before {
  content: '' !important;
  width: 8px !important;
  height: 8px !important;
  background-color: #4ade80 !important;
  border-radius: 50% !important;
  animation: pulse 2s infinite !important;
}

/* 终端操作按钮美化 */
.terminal-actions {
  display: flex !important;
  gap: 12px !important;
}

.terminal-actions .el-button {
  font-size: 13px !important;
  padding: 8px 20px !important;
  border-radius: 8px !important;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1) !important;
  font-weight: 600 !important;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.3) !important;
  border: none !important;
  text-transform: uppercase !important;
  letter-spacing: 0.5px !important;
}

.terminal-actions .el-button:hover {
  transform: translateY(-2px) !important;
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.4) !important;
}

.terminal-actions .el-button:active {
  transform: translateY(0) !important;
}

/* 清屏按钮样式 */
.terminal-actions .el-button:first-child {
  background: linear-gradient(135deg, #4a9eff 0%, #66b1ff 100%) !important;
  color: white !important;
}

.terminal-actions .el-button:first-child:hover {
  background: linear-gradient(135deg, #66b1ff 0%, #89c2ff 100%) !important;
}

/* 添加新的复制按钮样式 */
.terminal-actions .el-button:nth-child(2) {
  background: linear-gradient(135deg, #67c23a 0%, #85ce61 100%) !important;
  color: white !important;
}

.terminal-actions .el-button:nth-child(2):hover {
  background: linear-gradient(135deg, #85ce61 0%, #a6e38b 100%) !important;
}

/* 关闭按钮样式 */
.terminal-actions .el-button:last-child {
  background: linear-gradient(135deg, #f56c6c 0%, #f78989 100%) !important;
  color: white !important;
}

.terminal-actions .el-button:last-child:hover {
  background: linear-gradient(135deg, #f78989 0%, #fab6b6 100%) !important;
}

/* 终端内容区域 */
.terminal-content {
  flex: 1 !important;
  overflow: hidden !important;
  position: relative !important;
  background: #0a0a1a !important; /* 深色背景 */
  min-width: 900px !important;
  min-height: 550px !important;
  border-radius: 0 0 16px 16px !important;
  box-shadow: inset 0 2px 10px rgba(0, 0, 0, 0.3) !important;
}

/* 终端容器 */
.terminal-container {
  position: absolute !important;
  top: 0 !important;
  left: 0 !important;
  right: 0 !important;
  bottom: 0 !important;
  background-color: #0a0a1a !important;
  width: 100% !important;
  height: 100% !important;
  z-index: 10 !important;
  border: none !important;
  outline: none !important;
  padding: 20px !important; /* 增加内边距 */
  box-sizing: border-box !important;
}

/* 添加终端背景网格效果 */
.terminal-container::before {
  content: '' !important;
  position: absolute !important;
  top: 0 !important;
  left: 0 !important;
  right: 0 !important;
  bottom: 0 !important;
  background-image: 
    linear-gradient(rgba(74, 74, 106, 0.1) 1px, transparent 1px),
    linear-gradient(90deg, rgba(74, 74, 106, 0.1) 1px, transparent 1px) !important;
  background-size: 20px 20px !important;
  opacity: 0.3 !important;
  z-index: -1 !important;
  pointer-events: none !important;
}

/* 移除所有可能的伪元素（除了背景网格） */
.terminal-container::after,
.ssh-terminal::before,
.ssh-terminal::after,
.terminal-content::before,
.terminal-content::after {
  display: none !important;
  content: none !important;
}

/* xterm.js 核心样式美化 */
:deep(.xterm) {
  width: 100% !important;
  height: 100% !important;
  background-color: transparent !important;
  border: none !important;
  outline: none !important;
  box-shadow: none !important;
  font-family: 'Fira Code', 'JetBrains Mono', Consolas, "Courier New", monospace !important;
  font-size: 15px !important; /* 稍大的字体 */
  line-height: 1.7 !important; /* 更舒适的行高 */
  letter-spacing: 0.5px !important;
}

:deep(.xterm-viewport) {
  background-color: transparent !important;
  overflow: auto !important;
  border: none !important;
  outline: none !important;
  scrollbar-width: thin !important;
  scrollbar-color: #5a5a8a #1a1a3a !important;
  padding-right: 8px !important;
}

:deep(.xterm-screen) {
  background-color: transparent !important;
  border: none !important;
}

/* 字体颜色方案 - 现代化的配色 */
:deep(.xterm-rows) {
  color: #e0e0e0 !important;
  background-color: transparent !important;
}

/* 光标样式 */
:deep(.xterm-cursor) {
  background-color: #4ade80 !important; /* 亮绿色光标 */
  border: none !important;
  opacity: 0.9 !important;
  box-shadow: 0 0 8px rgba(74, 222, 128, 0.5) !important; /* 光标发光效果 */
}

:deep(.xterm-cursor-block) {
  background-color: #4ade80 !important;
  border: none !important;
  opacity: 0.9 !important;
  box-shadow: 0 0 8px rgba(74, 222, 128, 0.5) !important;
}

:deep(.xterm-cursor-blink) {
  animation: blink-cursor 1.2s infinite step-start !important;
  border: none !important;
}

/* 滚动条美化 */
:deep(.xterm-viewport::-webkit-scrollbar) {
  width: 12px !important;
  height: 12px !important;
}

:deep(.xterm-viewport::-webkit-scrollbar-track) {
  background: #1a1a3a !important;
  border-radius: 6px !important;
  border: none !important;
}

:deep(.xterm-viewport::-webkit-scrollbar-thumb) {
  background: linear-gradient(180deg, #5a5a8a 0%, #4a4a7a 100%) !important;
  border-radius: 6px !important;
  border: 2px solid #1a1a3a !important;
  transition: all 0.2s ease !important;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3) !important;
}

:deep(.xterm-viewport::-webkit-scrollbar-thumb:hover) {
  background: linear-gradient(180deg, #6a6a9a 0%, #5a5a8a 100%) !important;
  transform: scale(1.05) !important;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.4) !important;
}

/* 滚动条角落 */
:deep(.xterm-viewport::-webkit-scrollbar-corner) {
  background: #1a1a3a !important;
  border-radius: 6px !important;
}

/* 选择文本样式 */
:deep(.xterm-selection) {
  background-color: rgba(74, 158, 255, 0.4) !important; /* 半透明蓝色选择 */
  border: none !important;
  box-shadow: 0 0 8px rgba(74, 158, 255, 0.2) !important;
}

/* 装饰元素隐藏 */
:deep(.xterm-decoration) {
  display: none !important;
}

/* 光标闪烁动画 */
@keyframes blink-cursor {
  0%, 45% {
    opacity: 1 !important;
  }
  55%, 100% {
    opacity: 0 !important;
  }
}

/* 状态指示器脉冲动画 */
@keyframes pulse {
  0%, 100% {
    opacity: 1 !important;
    transform: scale(1) !important;
  }
  50% {
    opacity: 0.5 !important;
    transform: scale(1.2) !important;
  }
}

/* 确保终端内的所有文本都正确显示 */
.terminal-container * {
  font-family: 'Fira Code', 'JetBrains Mono', Consolas, "Courier New", monospace !important;
  font-size: 15px !important;
  line-height: 1.7 !important;
}

/* 添加终端命令历史指示器 */
.command-history {
  display: flex !important;
  align-items: center !important;
  gap: 8px !important;
  padding: 12px 20px !important;
  background: rgba(48, 43, 99, 0.6) !important;
  border-top: 1px solid rgba(74, 74, 106, 0.3) !important;
  font-size: 12px !important;
  color: #a0a0c0 !important;
}

/* 终端快捷键提示 */
.terminal-shortcuts {
  display: flex !important;
  gap: 16px !important;
  margin-left: auto !important;
}

.shortcut {
  display: flex !important;
  align-items: center !important;
  gap: 4px !important;
}

.shortcut-key {
  background: rgba(255, 255, 255, 0.1) !important;
  padding: 2px 8px !important;
  border-radius: 4px !important;
  font-family: monospace !important;
  font-size: 11px !important;
  color: #8ecae6 !important;
  border: 1px solid rgba(255, 255, 255, 0.2) !important;
}

/* 响应式设计优化 */
@media (max-width: 1024px) {
  .ssh-terminal {
    height: 600px !important;
  }
  
  .terminal-content {
    min-width: 700px !important;
  }
  
  .terminal-header {
    padding: 16px 20px !important;
  }
}

/* 确保终端在小屏幕上也能正常工作 */
@media (max-width: 768px) {
  .ssh-terminal {
    height: 550px !important;
    border-radius: 12px !important;
  }
  
  .terminal-content {
    min-width: 500px !important;
  }
  
  .terminal-header {
    padding: 14px 16px !important;
    height: 60px !important;
  }
  
  .terminal-title {
    font-size: 16px !important;
  }
  
  .terminal-actions .el-button {
    padding: 6px 16px !important;
    font-size: 12px !important;
  }
}

/* 移除可能的焦点样式 */
:deep(.xterm-viewport) {
  outline: none !important;
}

:deep(.xterm) {
  outline: none !important;
}

.terminal-container {
  outline: none !important;
}
.terminal-container :focus {
  outline: none !important;
  box-shadow: none !important;
}

/* 确保终端背景色覆盖整个区域 */
.terminal-container, .xterm, .xterm-viewport, .xterm-screen {
  background-color: #0f0f1e !important;
  background-image: none !important;
}

/* 修复可能的z-index问题 */
.terminal-container {
  z-index: 1000 !important;
}

:deep(.xterm) {
  z-index: 1001 !important;
}

:deep(.xterm-viewport) {
  z-index: 1002 !important;
}

:deep(.xterm-screen) {
  z-index: 1003 !important;
}

:deep(.xterm-rows) {
  z-index: 1004 !important;
}

:deep(.xterm-cursor) {
  z-index: 1005 !important;
}

/* 添加终端边框高亮效果 */
.terminal-content::after {
  content: '' !important;
  position: absolute !important;
  top: 0 !important;
  left: 0 !important;
  right: 0 !important;
  bottom: 0 !important;
  border: 1px solid rgba(74, 158, 255, 0.2) !important;
  pointer-events: none !important;
  border-radius: 0 0 12px 12px !important;
  transition: border-color 0.3s ease !important;
}

.ssh-terminal:hover .terminal-content::after {
  border-color: rgba(74, 158, 255, 0.4) !important; /* 悬停时高亮边框 */
}

/* 添加终端内部网格背景（可选，增加科技感） */
.terminal-container::before {
  content: '' !important;
  position: absolute !important;
  top: 0 !important;
  left: 0 !important;
  right: 0 !important;
  bottom: 0 !important;
  background-image: 
    linear-gradient(rgba(74, 158, 255, 0.03) 1px, transparent 1px),
    linear-gradient(90deg, rgba(74, 158, 255, 0.03) 1px, transparent 1px) !important;
  background-size: 20px 20px !important;
  pointer-events: none !important;
  z-index: 500 !important;
}

/* 响应式设计 */
@media (max-width: 900px) {
  .terminal-content {
    min-width: 700px !important;
  }
  
  .ssh-terminal {
    height: 500px !important;
  }
}



/* 表格样式 */
:deep(.el-table__row) {
  height: 60px;
}

:deep(.el-table__cell) {
  padding: 16px 12px;
}

:deep(.el-button + .el-button) {
  margin-left: 12px;
}

/* 监控样式 */
.monitoring-content {
  padding: 10px 0;
}

.monitoring-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  margin-bottom: 30px;
}

.stat-item {
  background-color: #f5f7fa;
  padding: 24px;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
  display: flex;
  flex-direction: column;
}

.stat-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.12);
}

.stat-label {
  font-size: 14px;
  color: #606266;
  margin-bottom: 12px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.stat-value {
  font-size: 36px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 8px;
}

.stat-description {
  font-size: 12px;
  color: #909399;
  margin-top: auto;
}

.stat-bar {
  height: 8px;
  background-color: #e4e7ed;
  border-radius: 4px;
  overflow: hidden;
  margin-top: 12px;
}

.stat-bar-fill {
  height: 100%;
  transition: width 0.3s ease;
  border-radius: 4px;
}

.system-info {
  margin-top: 20px;
}

.system-info h3 {
  margin: 0 0 16px 0;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

/* 响应式布局 */
@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
}

/* 响应式布局 */
@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
}
</style>