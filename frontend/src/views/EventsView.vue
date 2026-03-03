<template>
  <div class="events-view">
    <el-card shadow="never" class="page-header">
      <div class="header-content">
        <h2>K8s事件记录</h2>
        <div class="header-actions">
          <el-select v-model="selectedCluster" placeholder="选择集群" style="width: 200px; margin-right: 10px;">
            <el-option
              v-for="cluster in clusters"
              :key="cluster.id"
              :label="cluster.name"
              :value="cluster.id"
            />
          </el-select>
          <el-select v-model="namespaceFilter" placeholder="筛选命名空间" style="width: 200px; margin-right: 10px;" clearable>
            <el-option
              v-for="namespace in namespaces"
              :key="namespace"
              :label="namespace"
              :value="namespace"
            />
          </el-select>
          <el-select v-model="typeFilter" placeholder="筛选事件类型" style="width: 150px; margin-right: 10px;" clearable>
            <el-option label="Normal" value="Normal" />
            <el-option label="Warning" value="Warning" />
          </el-select>
          <el-input
            v-model="searchKey"
            placeholder="搜索事件信息"
            prefix-icon="Search"
            style="width: 250px; margin-right: 10px;"
            @input="handleSearch"
          />
          <el-button type="primary" size="small" @click="refreshEvents">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </div>
    </el-card>

    <el-card shadow="hover" class="events-card">
      <template #header>
        <div class="card-header">
          <span>K8s Events列表</span>
        </div>
      </template>
      
      <el-table 
        :data="events"
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
        :row-key="(row: any) => `${row.name}-${row.namespace}-${row.lastTimestamp}`"
      >
        <el-table-column prop="type" label="事件类型" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.type === 'Warning' ? 'danger' : 'success'" size="medium">
              {{ scope.row.type }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="reason" label="原因" width="120" show-overflow-tooltip />
        <el-table-column prop="message" label="事件信息" min-width="300" show-overflow-tooltip />
        <el-table-column prop="involvedObject.kind" label="资源类型" width="120" show-overflow-tooltip />
        <el-table-column prop="involvedObject.name" label="资源名称" width="180" show-overflow-tooltip />
        <el-table-column prop="namespace" label="命名空间" width="150" show-overflow-tooltip />
        <el-table-column prop="source.component" label="来源组件" width="150" show-overflow-tooltip />
        <el-table-column prop="count" label="事件次数" width="80" />
        <el-table-column prop="firstTimestamp" label="首次出现" width="180" sortable>
          <template #default="scope">
            {{ formatDateTime(scope.row.firstTimestamp) }}
          </template>
        </el-table-column>
        <el-table-column prop="lastTimestamp" label="最后出现" width="180" sortable>
          <template #default="scope">
            {{ formatDateTime(scope.row.lastTimestamp) }}
          </template>
        </el-table-column>
      </el-table>
      
      <div v-else class="no-cluster-selected">
        <el-empty description="请选择一个集群查看事件记录" />
      </div>
      
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
import { ref, onMounted, watch } from 'vue'
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

// 集群列表
const clusters = ref<Cluster[]>([])

// 选择的集群
const selectedCluster = ref<number | undefined>(undefined)

// 命名空间列表
const namespaces = ref<string[]>([])

// 筛选条件
const namespaceFilter = ref('')
const typeFilter = ref('')
const searchKey = ref('')

// 事件数据
const events = ref<any[]>([])

// 分页信息
const pagination = ref({
  page: 1,
  pageSize: 20,
  total: 0
})

// 组件挂载时加载数据
onMounted(() => {
  loadClusters()
})

// 监听集群选择变化，自动刷新事件列表
watch(selectedCluster, () => {
  if (selectedCluster.value) {
    loadNamespaces()
    loadEvents()
  }
})

// 监听筛选条件变化，自动刷新事件列表
watch([namespaceFilter, typeFilter], () => {
  if (selectedCluster.value) {
    handleSearch()
  }
})

// 加载集群列表
const loadClusters = async () => {
  try {
    const response = await axios.get('/clusters')
    clusters.value = response.data
  } catch (error) {
    console.error('加载集群列表失败:', error)
    clusters.value = []
    ElMessage.error('加载集群列表失败')
  }
}

// 加载命名空间列表
const loadNamespaces = async () => {
  if (!selectedCluster.value) return
  
  try {
    const response = await axios.get(`/clusters/${selectedCluster.value}/namespaces`)
    // 后端返回的是分页数据，命名空间列表在response.data.data中
    const namespaceData = response.data.data || []
    // 从返回的对象数组中提取name字段，转换为字符串数组
    namespaces.value = namespaceData.map((ns: any) => ns.name)
  } catch (error) {
    console.error('加载命名空间列表失败:', error)
    namespaces.value = []
  }
}

// 加载事件列表
const loadEvents = async () => {
  if (!selectedCluster.value) {
    ElMessage.warning('请先选择一个集群')
    return
  }
  
  loading.value = true
  try {
    const response = await axios.get(`/clusters/${selectedCluster.value}/events`, {
      params: {
        page: pagination.value.page,
        pageSize: pagination.value.pageSize,
        namespace: namespaceFilter.value,
        type: typeFilter.value,
        search: searchKey.value,
        _t: Date.now()
      }
    })
    
    events.value = response.data.items
    pagination.value.total = response.data.total
  } catch (error) {
    console.error('加载事件列表失败:', error)
    ElMessage.error('加载事件列表失败')
    events.value = []
    pagination.value.total = 0
  } finally {
    loading.value = false
  }
}

// 刷新事件列表
const refreshEvents = () => {
  pagination.value.page = 1
  loadEvents()
}

// 处理搜索输入
const handleSearch = () => {
  pagination.value.page = 1
  loadEvents()
}

// 处理分页大小变化
const handleSizeChange = (size: number) => {
  pagination.value.pageSize = size
  pagination.value.page = 1
  loadEvents()
}

// 处理当前页变化
const handleCurrentChange = (page: number) => {
  pagination.value.page = page
  loadEvents()
}

// 格式化日期时间
const formatDateTime = (dateTime: string) => {
  if (!dateTime) return ''
  
  try {
    const date = new Date(dateTime)
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
    return dateTime
  }
}
</script>

<style scoped>
.events-view {
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

.events-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
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

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>