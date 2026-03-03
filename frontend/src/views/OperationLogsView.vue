<template>
  <div class="operation-logs-view">
    <el-card shadow="never" class="page-header">
      <div class="header-content">
        <h2>操作日志</h2>
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
            placeholder="搜索用户名、菜单、操作或API"
            prefix-icon="Search"
            style="width: 300px"
            @keyup.enter="refreshLogs"
            @input="handleSearch"
          />
          <el-select v-model="statusFilter" placeholder="筛选状态" style="width: 120px; margin-left: 10px;" clearable @change="handleSearch">
            <el-option label="成功" value="success" />
            <el-option label="失败" value="failed" />
          </el-select>
        </div>
      </template>
      
      <el-table 
        :data="operationLogs" 
        style="width: 100%"
        stripe
        v-loading="loading"
        :row-key="(row: any) => row.ID"
        @selection-change="handleSelectionChange"
        fit
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="ID" label="ID" width="80" sortable />
        <el-table-column prop="username" label="操作用户" width="120" sortable />
        <el-table-column prop="menu" label="菜单" width="150" sortable />
        <el-table-column prop="operation" label="操作" width="150" sortable />
        <el-table-column prop="api" label="API" min-width="200" show-overflow-tooltip />
        <el-table-column prop="method" label="方法" width="80" sortable>
          <template #default="scope">
            <el-tag :type="getMethodTagType(scope.row.method)">
              {{ scope.row.method }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" sortable>
          <template #default="scope">
            <el-tag :type="scope.row.status === 'success' ? 'success' : 'danger'">
              {{ scope.row.status === 'success' ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="timeCost" label="耗时(ms)" width="100" sortable />
        <el-table-column prop="ip" label="IP地址" width="150" sortable />
        <el-table-column prop="createdAt" label="操作时间" width="180" sortable>
          <template #default="scope">
            {{ formatDateTime(scope.row.createdAt) }}
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

// 操作日志列表
const operationLogs = ref<any[]>([])

// 分页信息
const pagination = ref({
  page: 1,
  pageSize: 20,
  total: 0
})

// 选中的日志
const selectedLogs = ref<any[]>([])

// 获取请求方法的标签类型
const getMethodTagType = (method: string) => {
  const methodMap: Record<string, 'success' | 'warning' | 'danger' | 'info'> = {
    'GET': 'success',
    'POST': 'warning',
    'PUT': 'info',
    'DELETE': 'danger'
  }
  return methodMap[method] || 'info'
}

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

// 加载操作日志
const loadOperationLogs = async () => {
  loading.value = true
  try {
    const response = await logAPI.getOperationLogs({
      page: pagination.value.page,
      pageSize: pagination.value.pageSize,
      search: searchKey.value,
      status: statusFilter.value
    })
    
    operationLogs.value = response.data.items
    pagination.value.total = response.data.total
  } catch (error) {
    console.error('加载操作日志失败:', error)
    ElMessage.error('加载操作日志失败')
  } finally {
    loading.value = false
  }
}

// 刷新日志
const refreshLogs = () => {
  pagination.value.page = 1
  loadOperationLogs()
}

// 处理搜索输入和筛选变化
const handleSearch = () => {
  pagination.value.page = 1
  loadOperationLogs()
}

// 处理选择变化
const handleSelectionChange = (selection: any[]) => {
  selectedLogs.value = selection
}

// 处理分页大小变化
const handleSizeChange = (size: number) => {
  pagination.value.pageSize = size
  pagination.value.page = 1
  loadOperationLogs()
}

// 处理当前页变化
const handleCurrentChange = (page: number) => {
  pagination.value.page = page
  loadOperationLogs()
}

// 删除单个日志
const handleDelete = async (id: number) => {
  try {
    await logAPI.deleteOperationLog(id)
    ElMessage.success('删除成功')
    loadOperationLogs()
  } catch (error) {
    console.error('删除操作日志失败:', error)
    ElMessage.error('删除操作日志失败')
  }
}

// 批量删除日志
const handleBatchDelete = async () => {
  const ids = selectedLogs.value.map(log => log.ID)
  try {
    await logAPI.batchDeleteOperationLogs(ids)
    ElMessage.success('批量删除成功')
    selectedLogs.value = []
    loadOperationLogs()
  } catch (error) {
    console.error('批量删除操作日志失败:', error)
    ElMessage.error('批量删除操作日志失败')
  }
}

// 组件挂载时加载日志
onMounted(() => {
  loadOperationLogs()
})
</script>

<style scoped>
.operation-logs-view {
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
  gap: 10px;
}

.logs-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
