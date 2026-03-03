<template>
  <div class="alerts-view">
    <el-card shadow="never" class="page-header">
      <div class="header-content">
        <h2>告警管理</h2>
        <div class="header-actions">
          <el-button type="primary" @click="showSilenceDialog = true">
            <el-icon><Bell /></el-icon>
            新增静默
          </el-button>
          <el-button type="success" @click="showInhibitDialog = true">
            <el-icon><Lock /></el-icon>
            新增抑制
          </el-button>
          <el-button type="info" @click="refreshAlerts">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <el-statistic-card title="总告警数" :value="alerts.length" :value-style="{ color: '#3f8600' }">
          <template #suffix>
            <el-icon><Warning /></el-icon>
          </template>
        </el-statistic-card>
      </el-col>
      <el-col :span="6">
        <el-statistic-card title="触发中" :value="firingAlerts" :value-style="{ color: '#cf1322' }">
          <template #suffix>
            <el-icon><CircleClose /></el-icon>
          </template>
        </el-statistic-card>
      </el-col>
      <el-col :span="6">
        <el-statistic-card title="已解决" :value="resolvedAlerts" :value-style="{ color: '#13c2c2' }">
          <template #suffix>
            <el-icon><CircleCheck /></el-icon>
          </template>
        </el-statistic-card>
      </el-col>
      <el-col :span="6">
        <el-statistic-card title="已静默" :value="silencedAlerts" :value-style="{ color: '#722ed1' }">
          <template #suffix>
            <el-icon><Bell /></el-icon>
          </template>
        </el-statistic-card>
      </el-col>
    </el-row>

    <!-- 筛选区域 -->
    <el-card shadow="never" class="filter-card">
      <div class="filter-content">
        <el-input
          v-model="searchQuery"
          placeholder="搜索告警名称或实例"
          style="width: 300px; margin-right: 15px"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-select
          v-model="statusFilter"
          placeholder="状态筛选"
          style="width: 120px; margin-right: 15px"
        >
          <el-option label="全部" value=""></el-option>
          <el-option label="触发" value="firing"></el-option>
          <el-option label="解决" value="resolved"></el-option>
        </el-select>
        <el-select
          v-model="severityFilter"
          placeholder="级别筛选"
          style="width: 120px; margin-right: 15px"
        >
          <el-option label="全部" value=""></el-option>
          <el-option label="紧急" value="critical"></el-option>
          <el-option label="警告" value="warning"></el-option>
          <el-option label="信息" value="info"></el-option>
        </el-select>
        <el-button type="primary" @click="clearFilters">
          清除筛选
        </el-button>
      </div>
    </el-card>

    <!-- 告警列表 -->
    <el-card shadow="never" class="content-card">
      <el-table :data="paginatedAlerts" style="width: 100%" stripe size="small" fit>
        <el-table-column type="selection" width="55"></el-table-column>
        <el-table-column prop="alertName" label="告警名称" min-width="180"></el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'firing' ? 'danger' : 'success'" size="small">
              {{ scope.row.status === 'firing' ? '触发' : '解决' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="severity" label="级别" width="100">
          <template #default="scope">
            <el-tag :type="getSeverityType(scope.row.severity)" size="small">
              {{ getSeverityName(scope.row.severity) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="instance" label="实例" min-width="150">
          <template #default="scope">
            {{ JSON.parse(scope.row.labels).instance || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="job" label="任务" width="120">
          <template #default="scope">
            {{ JSON.parse(scope.row.labels).job || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="annotations" label="描述" min-width="250">
          <template #default="scope">
            {{ JSON.parse(scope.row.annotations).description || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="startsAt" label="开始时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.startsAt) }}
          </template>
        </el-table-column>
        <el-table-column prop="endsAt" label="结束时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.endsAt) }}
          </template>
        </el-table-column>
        <el-table-column prop="silenced" label="静默状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.silenced ? 'warning' : 'success'" size="small">
              {{ scope.row.silenced ? '已静默' : '正常' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="scope">
            <el-dropdown trigger="click">
              <el-button type="primary" size="small">
                操作 <el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item v-if="scope.row.status === 'firing'" @click="resolveAlert(scope.row.id)">
                    <el-icon><Check /></el-icon> 解决告警
                  </el-dropdown-item>
                  <el-dropdown-item @click="silenceAlert(scope.row)">
                    <el-icon><Bell /></el-icon> 静默告警
                  </el-dropdown-item>
                  <el-dropdown-item @click="viewAlertDetails(scope.row)">
                    <el-icon><View /></el-icon> 查看详情
                  </el-dropdown-item>
                  <el-dropdown-item @click="copyAlertInfo(scope.row)">
                    <el-icon><DocumentCopy /></el-icon> 复制信息
                  </el-dropdown-item>
                  <el-dropdown-item divided @click="deleteAlert(scope.row.id)">
                    <el-icon><Delete /></el-icon> 删除告警
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>
      
      <div class="pagination-container">
        <el-pagination
          layout="total, sizes, prev, pager, next, jumper"
          :total="filteredAlerts.length"
          :page-size="pageSize"
          v-model:current-page="currentPage"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        ></el-pagination>
      </div>
    </el-card>

    <!-- 告警详情对话框 -->
    <el-dialog v-model="showAlertDetails" title="告警详情" width="800px">
      <div class="alert-details" v-if="selectedAlert">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="告警名称">{{ selectedAlert.alertName }}</el-descriptions-item>
          <el-descriptions-item label="状态">{{ selectedAlert.status === 'firing' ? '触发' : '解决' }}</el-descriptions-item>
          <el-descriptions-item label="级别">{{ getSeverityName(selectedAlert.severity) }}</el-descriptions-item>
          <el-descriptions-item label="指纹">{{ selectedAlert.fingerprint }}</el-descriptions-item>
          <el-descriptions-item label="开始时间" :span="2">{{ formatDate(selectedAlert.startsAt) }}</el-descriptions-item>
          <el-descriptions-item label="结束时间" :span="2">{{ formatDate(selectedAlert.endsAt) }}</el-descriptions-item>
          <el-descriptions-item label="生成链接" :span="2">
            <el-link :href="selectedAlert.generatorURL" target="_blank">{{ selectedAlert.generatorURL }}</el-link>
          </el-descriptions-item>
          <el-descriptions-item label="标签" :span="2">
            <pre class="alert-labels">{{ JSON.stringify(JSON.parse(selectedAlert.labels), null, 2) }}</pre>
          </el-descriptions-item>
          <el-descriptions-item label="注解" :span="2">
            <pre class="alert-annotations">{{ JSON.stringify(JSON.parse(selectedAlert.annotations), null, 2) }}</pre>
          </el-descriptions-item>
        </el-descriptions>
      </div>
      <template #footer>
        <el-button @click="showAlertDetails = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 静默对话框 -->
    <el-dialog v-model="showSilenceDialog" title="新增静默规则" width="600px">
      <el-form :model="silenceForm" label-width="100px">
        <el-form-item label="名称">
          <el-input v-model="silenceForm.name" placeholder="请输入静默规则名称"></el-input>
        </el-form-item>
        <el-form-item label="匹配标签">
          <el-input
            v-model="silenceForm.matchers"
            type="textarea"
            :rows="3"
            placeholder="请输入匹配标签，支持JSON格式"
          ></el-input>
        </el-form-item>
        <el-form-item label="开始时间">
          <el-date-picker
            v-model="silenceForm.startsAt"
            type="datetime"
            placeholder="选择开始时间"
            style="width: 100%"
          ></el-date-picker>
        </el-form-item>
        <el-form-item label="结束时间">
          <el-date-picker
            v-model="silenceForm.endsAt"
            type="datetime"
            placeholder="选择结束时间"
            style="width: 100%"
          ></el-date-picker>
        </el-form-item>
        <el-form-item label="注释">
          <el-input
            v-model="silenceForm.comment"
            type="textarea"
            :rows="2"
            placeholder="请输入静默规则注释"
          ></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showSilenceDialog = false">取消</el-button>
        <el-button type="primary" @click="saveSilence">保存</el-button>
      </template>
    </el-dialog>

    <!-- 抑制对话框 -->
    <el-dialog v-model="showInhibitDialog" title="新增抑制规则" width="600px">
      <el-form :model="inhibitForm" label-width="100px">
        <el-form-item label="名称">
          <el-input v-model="inhibitForm.name" placeholder="请输入抑制规则名称"></el-input>
        </el-form-item>
        <el-form-item label="源告警标签">
          <el-input
            v-model="inhibitForm.sourceMatchers"
            type="textarea"
            :rows="2"
            placeholder="请输入源告警匹配标签"
          ></el-input>
        </el-form-item>
        <el-form-item label="目标告警标签">
          <el-input
            v-model="inhibitForm.targetMatchers"
            type="textarea"
            :rows="2"
            placeholder="请输入目标告警匹配标签"
          ></el-input>
        </el-form-item>
        <el-form-item label="注释">
          <el-input
            v-model="inhibitForm.comment"
            type="textarea"
            :rows="2"
            placeholder="请输入抑制规则注释"
          ></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showInhibitDialog = false">取消</el-button>
        <el-button type="primary" @click="saveInhibit">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { 
  Search, Refresh, Bell, Lock, Warning, CircleClose, CircleCheck,
  Check, View, DocumentCopy, Delete, ArrowDown
} from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'

// 模拟告警数据
const alerts = ref([
  {
    id: 1,
    alertName: 'HighCPUUsage',
    status: 'firing',
    severity: 'critical',
    labels: JSON.stringify({ instance: 'server1:9100', job: 'node_exporter' }),
    annotations: JSON.stringify({ 
      description: 'CPU使用率超过80%，持续5分钟',
      summary: '高CPU使用率告警' 
    }),
    startsAt: '2024-01-01T12:00:00Z',
    endsAt: '0001-01-01T00:00:00Z',
    generatorURL: 'http://localhost:9090/graph?g0.expr=node_cpu_seconds_total',
    fingerprint: 'abc123',
    silenced: false
  },
  {
    id: 2,
    alertName: 'LowDiskSpace',
    status: 'firing',
    severity: 'warning',
    labels: JSON.stringify({ instance: 'server2:9100', job: 'node_exporter' }),
    annotations: JSON.stringify({ 
      description: '磁盘空间不足，剩余空间小于20%',
      summary: '低磁盘空间告警' 
    }),
    startsAt: '2024-01-01T11:30:00Z',
    endsAt: '0001-01-01T00:00:00Z',
    generatorURL: 'http://localhost:9090/graph?g0.expr=node_filesystem_free_bytes',
    fingerprint: 'def456',
    silenced: true
  },
  {
    id: 3,
    alertName: 'HighMemoryUsage',
    status: 'resolved',
    severity: 'info',
    labels: JSON.stringify({ instance: 'server3:9100', job: 'node_exporter' }),
    annotations: JSON.stringify({ 
      description: '内存使用率超过90%，已恢复',
      summary: '高内存使用率告警' 
    }),
    startsAt: '2024-01-01T10:00:00Z',
    endsAt: '2024-01-01T11:00:00Z',
    generatorURL: 'http://localhost:9090/graph?g0.expr=node_memory_MemTotal_bytes',
    fingerprint: 'ghi789',
    silenced: false
  }
])

// 筛选条件
const searchQuery = ref('')
const statusFilter = ref('')
const severityFilter = ref('')

// 分页
const currentPage = ref(1)
const pageSize = ref(10)

// 对话框状态
const showAlertDetails = ref(false)
const showSilenceDialog = ref(false)
const showInhibitDialog = ref(false)
// 定义告警类型接口
interface Alert {
  id: number;
  alertName: string;
  status: 'firing' | 'resolved';
  severity: 'critical' | 'warning' | 'info';
  labels: string;
  annotations: string;
  startsAt: string;
  endsAt: string;
  generatorURL: string;
  fingerprint: string;
  silenced: boolean;
}

const selectedAlert = ref<Alert | null>(null)

// 表单数据
const silenceForm = ref({
  name: '',
  matchers: '',
  startsAt: new Date(),
  endsAt: new Date(Date.now() + 3600000), // 1小时后
  comment: ''
})

const inhibitForm = ref({
  name: '',
  sourceMatchers: '',
  targetMatchers: '',
  comment: ''
})

// 计算属性
const firingAlerts = computed(() => {
  return alerts.value.filter(alert => alert.status === 'firing').length
})

const resolvedAlerts = computed(() => {
  return alerts.value.filter(alert => alert.status === 'resolved').length
})

const silencedAlerts = computed(() => {
  return alerts.value.filter(alert => alert.silenced).length
})

const filteredAlerts = computed(() => {
  return alerts.value.filter(alert => {
    const matchesSearch = 
      alert.alertName.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      JSON.parse(alert.labels).instance?.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      JSON.parse(alert.labels).job?.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesStatus = statusFilter.value ? alert.status === statusFilter.value : true
    const matchesSeverity = severityFilter.value ? alert.severity === severityFilter.value : true
    return matchesSearch && matchesStatus && matchesSeverity
  })
})

const paginatedAlerts = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredAlerts.value.slice(start, end)
})

// 方法
const getSeverityType = (severity: string) => {
  switch (severity) {
    case 'critical':
      return 'danger'
    case 'warning':
      return 'warning'
    case 'info':
      return 'info'
    default:
      return 'success'
  }
}

const getSeverityName = (severity: string) => {
  switch (severity) {
    case 'critical':
      return '紧急'
    case 'warning':
      return '警告'
    case 'info':
      return '信息'
    default:
      return severity
  }
}

const formatDate = (dateString: string) => {
  if (!dateString || dateString === '0001-01-01T00:00:00Z') return '-'
  const date = new Date(dateString)
  return date.toLocaleString()
}

const refreshAlerts = () => {
  ElMessage.info('正在刷新告警数据...')
  // 这里可以添加API调用，从后端获取最新告警数据
  setTimeout(() => {
    ElMessage.success('告警数据刷新成功')
  }, 1000)
}

const clearFilters = () => {
  searchQuery.value = ''
  statusFilter.value = ''
  severityFilter.value = ''
}

const resolveAlert = (id: number) => {
  const alert = alerts.value.find(a => a.id === id)
  if (alert) {
    alert.status = 'resolved'
    alert.endsAt = new Date().toISOString()
    ElMessage.success('告警已解决')
  }
}

const silenceAlert = (alert: any) => {
  alert.silenced = !alert.silenced
  ElMessage.success(alert.silenced ? '告警已静默' : '告警静默已取消')
}

const viewAlertDetails = (alert: any) => {
  selectedAlert.value = alert
  showAlertDetails.value = true
}

const copyAlertInfo = (alert: any) => {
  const info = `告警名称: ${alert.alertName}\n状态: ${alert.status === 'firing' ? '触发' : '解决'}\n级别: ${getSeverityName(alert.severity)}\n实例: ${JSON.parse(alert.labels).instance}\n开始时间: ${formatDate(alert.startsAt)}`
  navigator.clipboard.writeText(info)
  ElMessage.success('告警信息已复制到剪贴板')
}

const deleteAlert = (id: number) => {
  ElMessageBox.confirm('确定要删除这个告警吗？', '删除确认', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    alerts.value = alerts.value.filter(a => a.id !== id)
    ElMessage.success('告警已删除')
  }).catch(() => {
    // 取消删除
  })
}

const saveSilence = () => {
  // 这里可以添加API调用，保存静默规则
  showSilenceDialog.value = false
  ElMessage.success('静默规则已保存')
}

const saveInhibit = () => {
  // 这里可以添加API调用，保存抑制规则
  showInhibitDialog.value = false
  ElMessage.success('抑制规则已保存')
}

// 分页处理
const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
}
</script>

<style scoped>
.alerts-view {
  width: 100%;
  padding: 20px;
  background-color: #f5f7fa;
  min-height: 100vh;
}

.page-header {
  margin-bottom: 20px;
  background-color: #fff;
  border-radius: 8px;
  padding: 16px 24px;
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

.stats-row {
  margin-bottom: 20px;
}

.filter-card {
  margin-bottom: 20px;
  background-color: #fff;
  border-radius: 8px;
}

.filter-content {
  display: flex;
  align-items: center;
  padding: 16px 24px;
}

.content-card {
  margin-bottom: 20px;
  background-color: #fff;
  border-radius: 8px;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
  padding-right: 20px;
}

:deep(.el-table__header-wrapper) {
  background-color: #fafafa;
}

.alert-details {
  padding: 10px 0;
}

.alert-labels,
.alert-annotations {
  background-color: #f5f7fa;
  padding: 10px;
  border-radius: 4px;
  margin: 0;
  font-family: monospace;
  font-size: 13px;
  line-height: 1.5;
}

/* 统计卡片样式 */
:deep(.el-statistic-card) {
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;
}

:deep(.el-statistic-card:hover) {
  box-shadow: 0 4px 16px 0 rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}
</style>