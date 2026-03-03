<template>
  <div class="login-logs-view">
    <el-card shadow="never" class="page-header">
      <div class="header-content">
        <h2>登录日志</h2>
        <div class="header-actions">
          <el-button type="danger" size="small" @click="handleBatchDelete" :disabled="selectedLogs.length === 0">
            <el-icon><Delete /></el-icon>
            批量删除
          </el-button>
          <el-button type="primary" size="small" @click="refreshLogs">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </div>
    </el-card>

    <el-card shadow="hover" class="logs-card">
      <template #header>
        <div class="card-header">
          <el-input
            v-model="searchKey"
            placeholder="搜索用户名、IP地址或位置"
            prefix-icon="Search"
            style="width: 300px"
            @keyup.enter="refreshLogs"
          />
          <el-select v-model="statusFilter" placeholder="筛选状态" style="width: 120px; margin-left: 10px;" clearable>
            <el-option label="成功" value="success" />
            <el-option label="失败" value="failed" />
          </el-select>
        </div>
      </template>
      
      <el-table 
        :data="loginLogs" 
        style="width: 100%"
        stripe
        v-loading="loading"
        :row-key="(row: any) => row.ID"
        @selection-change="handleSelectionChange"
        fit
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="ID" label="ID" width="80" sortable />
        <el-table-column prop="username" label="用户名" width="120" sortable />
        <el-table-column prop="ip" label="IP地址" width="150" sortable />
        <el-table-column prop="location" label="地理位置" min-width="180" sortable />
        <el-table-column prop="status" label="登录状态" width="120" sortable>
          <template #default="scope">
            <el-tag :type="scope.row.status === 'success' ? 'success' : 'danger'">
              {{ scope.row.status === 'success' ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="loginTime" label="登录时间" width="180" sortable>
          <template #default="scope">
            {{ formatDateTime(scope.row.loginTime) }}
          </template>
        </el-table-column>
        <el-table-column label="系统类型" min-width="150" sortable>
          <template #default="scope">
            {{ getSystemType(scope.row.userAgent) }}
          </template>
        </el-table-column>
        <el-table-column label="浏览器" min-width="150" sortable>
          <template #default="scope">
            {{ getBrowser(scope.row.userAgent) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="scope">
            <el-button type="danger" size="small" @click="handleDelete(scope.row.ID)">
              <el-icon><Delete /></el-icon>
              删除
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
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { logAPI } from '../api/api'
import { Delete, Refresh } from '@element-plus/icons-vue'

// 加载状态
const loading = ref(false)

// 搜索和筛选
const searchKey = ref('')
const statusFilter = ref('')

// 登录日志列表
const loginLogs = ref<any[]>([])

// 分页信息
const pagination = ref({
  page: 1,
  pageSize: 20,
  total: 0
})

// 选中的日志
const selectedLogs = ref<any[]>([])

// 格式化日期时间
const formatDateTime = (dateTime: string) => {
  if (!dateTime) return ''
  const date = new Date(dateTime)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

// 解析用户代理，获取系统类型
const getSystemType = (userAgent: string) => {
  if (!userAgent) return ''
  
  if (userAgent.includes('Windows')) {
    return 'Windows'
  } else if (userAgent.includes('Macintosh')) {
    return 'macOS'
  } else if (userAgent.includes('Linux')) {
    return 'Linux'
  } else if (userAgent.includes('Android')) {
    return 'Android'
  } else if (userAgent.includes('iPhone') || userAgent.includes('iPad')) {
    return 'iOS'
  } else {
    return '未知系统'
  }
}

// 解析用户代理，获取浏览器信息
const getBrowser = (userAgent: string) => {
  if (!userAgent) return ''
  
  if (userAgent.includes('Chrome')) {
    return 'Chrome'
  } else if (userAgent.includes('Firefox')) {
    return 'Firefox'
  } else if (userAgent.includes('Safari')) {
    return 'Safari'
  } else if (userAgent.includes('Edge')) {
    return 'Edge'
  } else if (userAgent.includes('Opera') || userAgent.includes('OPR')) {
    return 'Opera'
  } else if (userAgent.includes('MSIE') || userAgent.includes('Trident')) {
    return 'Internet Explorer'
  } else {
    return '未知浏览器'
  }
}

// 加载登录日志
const loadLoginLogs = async () => {
  loading.value = true
  try {
    const response = await logAPI.getLoginLogs({
      page: pagination.value.page,
      pageSize: pagination.value.pageSize,
      search: searchKey.value,
      status: statusFilter.value
    })
    
    loginLogs.value = response.data.items
    pagination.value.total = response.data.total
  } catch (error) {
    console.error('加载登录日志失败:', error)
    ElMessage.error('加载登录日志失败')
  } finally {
    loading.value = false
  }
}

// 刷新日志
const refreshLogs = () => {
  pagination.value.page = 1
  loadLoginLogs()
}

// 处理选择变化
const handleSelectionChange = (selection: any[]) => {
  selectedLogs.value = selection
}

// 处理分页大小变化
const handleSizeChange = (size: number) => {
  pagination.value.pageSize = size
  pagination.value.page = 1
  loadLoginLogs()
}

// 处理当前页变化
const handleCurrentChange = (page: number) => {
  pagination.value.page = page
  loadLoginLogs()
}

// 删除单个日志
const handleDelete = async (id: number) => {
  try {
    await logAPI.deleteLoginLog(id)
    ElMessage.success('删除成功')
    loadLoginLogs()
  } catch (error) {
    console.error('删除登录日志失败:', error)
    ElMessage.error('删除登录日志失败')
  }
}

// 批量删除日志
const handleBatchDelete = async () => {
  const ids = selectedLogs.value.map(log => log.ID)
  try {
    await logAPI.batchDeleteLoginLogs(ids)
    ElMessage.success('批量删除成功')
    selectedLogs.value = []
    loadLoginLogs()
  } catch (error) {
    console.error('批量删除登录日志失败:', error)
    ElMessage.error('批量删除登录日志失败')
  }
}

// 组件挂载时加载日志
onMounted(() => {
  loadLoginLogs()
})
</script>

<style scoped>
.login-logs-view {
  width: 100%;
  padding: 0;
  max-width: 100vw;
  box-sizing: border-box;
}

.page-header {
  margin-bottom: 20px;
  width: 100%;
  max-width: 100%;
  box-sizing: border-box;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  box-sizing: border-box;
}

.header-content h2 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: #333;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.logs-card {
  margin-bottom: 20px;
  width: 100%;
  max-width: 100%;
  box-sizing: border-box;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  width: 100%;
  box-sizing: border-box;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
  width: 100%;
  box-sizing: border-box;
}

/* 确保表格占满整个容器宽度 */
:deep(.el-table) {
  width: 100% !important;
  max-width: 100% !important;
  box-sizing: border-box !important;
}

/* 确保表格容器占满整个宽度 */
:deep(.el-card__body) {
  width: 100%;
  max-width: 100%;
  box-sizing: border-box;
  padding: 0 20px 20px;
}

/* 确保分页组件占满宽度 */
:deep(.el-pagination) {
  margin-top: 20px;
  width: 100%;
  box-sizing: border-box;
}
</style>
