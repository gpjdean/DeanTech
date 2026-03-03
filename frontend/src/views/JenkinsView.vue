<template>
  <div class="jenkins-view">
    <el-card shadow="never" class="page-header">
      <div class="header-content">
        <h2>Jenkins管理</h2>
        <div class="header-actions">
          <el-button type="primary" size="small" @click="connectJenkins">
            <el-icon><Connection /></el-icon>
            连接Jenkins
          </el-button>
          <el-button type="primary" size="small" @click="refreshProjects">
            <el-icon><Refresh /></el-icon>
            刷新项目
          </el-button>
        </div>
      </div>
    </el-card>

    <el-card shadow="hover" class="jenkins-card">
      <template #header>
        <div class="card-header">
          <span>Jenkins项目列表</span>
        </div>
      </template>
      
      <div v-if="!isConnected" class="no-connection">
        <el-empty description="请先连接Jenkins服务器" />
      </div>
      
      <div v-else>
        <el-table 
          :data="projects"
          style="width: 100%; color: #303133; font-size: 14px; line-height: 1.5;"
          stripe
          v-loading="loading"
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
          <el-table-column prop="name" label="项目名称" min-width="200" show-overflow-tooltip />
          <el-table-column prop="url" label="项目URL" min-width="300" show-overflow-tooltip>
            <template #default="scope">
              <a :href="scope.row.url" target="_blank" class="jenkins-link">{{ scope.row.url }}</a>
            </template>
          </el-table-column>
          <el-table-column prop="lastBuildNumber" label="最后构建号" width="120" />
          <el-table-column prop="lastBuildStatus" label="最后构建状态" width="120">
            <template #default="scope">
              <el-tag :type="getBuildStatusTagType(scope.row.lastBuildStatus)" size="medium">
                {{ scope.row.lastBuildStatus || '无构建记录' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="lastBuildTime" label="最后构建时间" width="180" sortable>
            <template #default="scope">
              {{ formatDateTime(scope.row.lastBuildTime) }}
            </template>
          </el-table-column>
          <el-table-column prop="description" label="项目描述" min-width="200" show-overflow-tooltip />
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="scope">
              <el-button type="primary" size="small" @click="buildProject(scope.row)" :disabled="loadingBuilds.includes(scope.row.name)">
                <el-icon><CirclePlus /></el-icon>
                立即构建
              </el-button>
              <el-button type="info" size="small" @click="viewBuildHistory(scope.row)" style="margin-left: 10px;">
                <el-icon><View /></el-icon>
                构建历史
              </el-button>
              <el-button type="warning" size="small" @click="stopProject(scope.row)" style="margin-left: 10px;" :disabled="loadingBuilds.includes(scope.row.name)">
                <el-icon><Close /></el-icon>
                停止构建
              </el-button>
            </template>
          </el-table-column>
        </el-table>
        
        <div class="pagination-container">
          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            :total="pagination.total"
            page-text="页"
            size-text="/页"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </div>
    </el-card>
    
    <!-- Jenkins连接配置对话框 -->
    <el-dialog
      v-model="connectionDialogVisible"
      title="Jenkins连接配置"
      width="500px"
    >
      <el-form :model="jenkinsConfig" label-width="80px" class="connection-form">
        <el-form-item label="Jenkins URL">
          <el-input v-model="jenkinsConfig.url" placeholder="http://jenkins.example.com" />
        </el-form-item>
        <el-form-item label="用户名">
          <el-input v-model="jenkinsConfig.username" placeholder="admin" />
        </el-form-item>
        <el-form-item label="API令牌">
          <el-input v-model="jenkinsConfig.apiToken" type="password" placeholder="Jenkins API Token" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="connectionDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveJenkinsConfig">保存配置</el-button>
        </span>
      </template>
    </el-dialog>
    
    <!-- 构建历史对话框 -->
    <el-dialog
      v-model="buildHistoryVisible"
      :title="`${selectedProject?.name || ''} 构建历史`"
      width="800px"
    >
      <div v-if="buildHistory.length > 0">
        <el-table 
          :data="buildHistory"
          style="width: 100%"
          stripe
          :row-key="(row: any) => row.buildNumber"
        >
          <el-table-column prop="buildNumber" label="构建号" width="100" sortable />
          <el-table-column prop="status" label="状态" width="120">
            <template #default="scope">
              <el-tag :type="getBuildStatusTagType(scope.row.status)" size="medium">
                {{ scope.row.status }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="duration" label="持续时间" width="120">
            <template #default="scope">
              {{ formatDuration(scope.row.duration) }}
            </template>
          </el-table-column>
          <el-table-column prop="timestamp" label="构建时间" width="180" sortable>
            <template #default="scope">
              {{ formatDateTime(scope.row.timestamp) }}
            </template>
          </el-table-column>
          <el-table-column prop="builtBy" label="构建者" width="120" />
          <el-table-column label="操作" width="150">
            <template #default="scope">
              <el-button type="primary" size="small" @click="viewBuildDetails(scope.row)">
                <el-icon><View /></el-icon>
                查看详情
              </el-button>
              <el-button type="danger" size="small" @click="deleteBuild(scope.row)" style="margin-left: 10px;">
                <el-icon><Delete /></el-icon>
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
      <div v-else>
        <el-empty description="暂无构建记录" />
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="buildHistoryVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Connection, Refresh, View, CirclePlus, Close, Delete } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'

// 加载状态
const loading = ref(false)
const isConnected = ref(false)
const connectionDialogVisible = ref(false)
const buildHistoryVisible = ref(false)
const selectedProject = ref<any>(null)
const buildHistory = ref<any[]>([])
const loadingBuilds = ref<string[]>([])

// Jenkins配置
const jenkinsConfig = ref({
  url: '',
  username: '',
  apiToken: ''
})

// Jenkins项目列表
const projects = ref<any[]>([])

// 分页信息
const pagination = ref({
  page: 1,
  pageSize: 20,
  total: 0
})

// 组件挂载时加载数据
onMounted(() => {
  // 检查本地存储中是否有Jenkins配置
  const savedConfig = localStorage.getItem('jenkinsConfig')
  if (savedConfig) {
    jenkinsConfig.value = JSON.parse(savedConfig)
    isConnected.value = true
    loadProjects()
  }
})

// 获取构建状态的标签类型
const getBuildStatusTagType = (status: string) => {
  switch (status) {
    case 'SUCCESS':
      return 'success'
    case 'FAILURE':
      return 'danger'
    case 'ABORTED':
      return 'warning'
    case 'IN_PROGRESS':
      return 'info'
    default:
      return 'info'
  }
}

// 格式化日期时间
const formatDateTime = (timestamp: any) => {
  if (!timestamp) return ''
  
  try {
    const date = new Date(timestamp)
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    })
  } catch (error) {
    console.error('时间格式化失败:', error)
    return ''
  }
}

// 格式化持续时间
const formatDuration = (milliseconds: number) => {
  if (!milliseconds) return '0s'
  
  const seconds = Math.floor(milliseconds / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)
  
  if (hours > 0) {
    return `${hours}h ${minutes % 60}m ${seconds % 60}s`
  } else if (minutes > 0) {
    return `${minutes}m ${seconds % 60}s`
  } else {
    return `${seconds}s`
  }
}

// 连接Jenkins
const connectJenkins = () => {
  connectionDialogVisible.value = true
}

// 保存Jenkins配置
const saveJenkinsConfig = async () => {
  if (!jenkinsConfig.value.url || !jenkinsConfig.value.username || !jenkinsConfig.value.apiToken) {
    ElMessage.warning('请填写完整的Jenkins配置')
    return
  }
  
  // 验证连接
  try {
    loading.value = true
    // 这里应该调用实际的API来验证Jenkins连接
    // 目前先模拟验证成功
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // 保存配置到本地存储
    localStorage.setItem('jenkinsConfig', JSON.stringify(jenkinsConfig.value))
    isConnected.value = true
    connectionDialogVisible.value = false
    ElMessage.success('Jenkins连接配置保存成功')
    
    // 加载项目列表
    loadProjects()
  } catch (error) {
    console.error('Jenkins连接验证失败:', error)
    ElMessage.error('Jenkins连接验证失败，请检查配置')
  } finally {
    loading.value = false
  }
}

// 加载Jenkins项目列表
const loadProjects = async () => {
  try {
    loading.value = true
    // 这里应该调用实际的Jenkins API来获取项目列表
    // 目前先模拟返回数据
    await new Promise(resolve => setTimeout(resolve, 1500))
    
    // 模拟Jenkins项目数据
    projects.value = [
      {
        name: 'demo-project-1',
        url: 'http://jenkins.example.com/job/demo-project-1/',
        lastBuildNumber: 123,
        lastBuildStatus: 'SUCCESS',
        lastBuildTime: Date.now() - 3600000,
        description: '演示项目1，用于测试Jenkins集成'
      },
      {
        name: 'demo-project-2',
        url: 'http://jenkins.example.com/job/demo-project-2/',
        lastBuildNumber: 456,
        lastBuildStatus: 'FAILURE',
        lastBuildTime: Date.now() - 7200000,
        description: '演示项目2，用于测试Jenkins构建失败场景'
      },
      {
        name: 'demo-project-3',
        url: 'http://jenkins.example.com/job/demo-project-3/',
        lastBuildNumber: 789,
        lastBuildStatus: 'ABORTED',
        lastBuildTime: Date.now() - 10800000,
        description: '演示项目3，用于测试Jenkins构建中止场景'
      }
    ]
    
    pagination.value.total = projects.value.length
  } catch (error) {
    console.error('加载Jenkins项目失败:', error)
    ElMessage.error('加载Jenkins项目失败')
    projects.value = []
    pagination.value.total = 0
  } finally {
    loading.value = false
  }
}

// 刷新项目列表
const refreshProjects = () => {
  loadProjects()
}

// 构建项目
const buildProject = async (project: any) => {
  try {
    loadingBuilds.value.push(project.name)
    // 这里应该调用实际的Jenkins API来触发构建
    // 目前先模拟构建过程
    await new Promise(resolve => setTimeout(resolve, 2000))
    
    ElMessage.success(`项目 ${project.name} 构建触发成功`)
    // 刷新项目列表
    loadProjects()
  } catch (error) {
    console.error('构建项目失败:', error)
    ElMessage.error(`项目 ${project.name} 构建触发失败`)
  } finally {
    loadingBuilds.value = loadingBuilds.value.filter(name => name !== project.name)
  }
}

// 停止构建
const stopProject = async (project: any) => {
  try {
    loadingBuilds.value.push(project.name)
    // 这里应该调用实际的Jenkins API来停止构建
    // 目前先模拟停止构建过程
    await new Promise(resolve => setTimeout(resolve, 1500))
    
    ElMessage.success(`项目 ${project.name} 构建停止成功`)
    // 刷新项目列表
    loadProjects()
  } catch (error) {
    console.error('停止项目构建失败:', error)
    ElMessage.error(`项目 ${project.name} 构建停止失败`)
  } finally {
    loadingBuilds.value = loadingBuilds.value.filter(name => name !== project.name)
  }
}

// 查看构建历史
const viewBuildHistory = (project: any) => {
  selectedProject.value = project
  // 这里应该调用实际的Jenkins API来获取构建历史
  // 目前先模拟构建历史数据
  buildHistory.value = [
    {
      buildNumber: 123,
      status: 'SUCCESS',
      duration: 300000,
      timestamp: Date.now() - 3600000,
      builtBy: 'admin'
    },
    {
      buildNumber: 122,
      status: 'FAILURE',
      duration: 250000,
      timestamp: Date.now() - 7200000,
      builtBy: 'user1'
    },
    {
      buildNumber: 121,
      status: 'SUCCESS',
      duration: 320000,
      timestamp: Date.now() - 10800000,
      builtBy: 'admin'
    }
  ]
  buildHistoryVisible.value = true
}

// 查看构建详情
const viewBuildDetails = (_build: any) => {
  ElMessage.info('查看构建详情功能开发中')
}

// 删除构建记录
const deleteBuild = async (build: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除构建记录 #${build.buildNumber} 吗？此操作无法撤销。`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 这里应该调用实际的Jenkins API来删除构建记录
    // 目前先模拟删除过程
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // 从列表中移除
    const index = buildHistory.value.findIndex(item => item.buildNumber === build.buildNumber)
    if (index !== -1) {
      buildHistory.value.splice(index, 1)
    }
    
    ElMessage.success('构建记录删除成功')
  } catch (error) {
    // 用户取消删除
  }
}

// 处理分页大小变化
const handleSizeChange = (size: number) => {
  pagination.value.pageSize = size
  pagination.value.page = 1
  loadProjects()
}

// 处理当前页变化
const handleCurrentChange = (page: number) => {
  pagination.value.page = page
  loadProjects()
}
</script>

<style scoped>
.jenkins-view {
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

.jenkins-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.no-connection {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 300px;
  background-color: #fafafa;
  border-radius: 8px;
  margin-top: 20px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.jenkins-link {
  color: #409eff;
  text-decoration: none;
}

.jenkins-link:hover {
  text-decoration: underline;
}

.connection-form {
  margin-top: 20px;
}
</style>