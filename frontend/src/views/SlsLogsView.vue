<template>
  <div class="sls-logs-view">
    <el-card shadow="never" class="page-header">
      <div class="header-content">
        <h2>SLS日志查询</h2>
      </div>
    </el-card>

    <el-card shadow="hover" class="sls-logs-card">
      <template #header>
        <div class="card-header">
          <span>日志查询</span>
          <div class="header-actions">
            <el-button type="primary" @click="queryLogs">
              <el-icon><Search /></el-icon>
              查询
            </el-button>
          </div>
        </div>
      </template>
      
      <el-form :model="queryForm" label-width="120px" class="query-form">
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="日志仓库">
              <el-select v-model="queryForm.project" placeholder="选择日志仓库">
                <el-option
                  v-for="project in projects"
                  :key="project"
                  :label="project"
                  :value="project"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="日志主题">
              <el-select v-model="queryForm.logstore" placeholder="选择日志主题">
                <el-option
                  v-for="logstore in logstores"
                  :key="logstore"
                  :label="logstore"
                  :value="logstore"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="查询时间">
              <el-date-picker
                v-model="queryForm.timeRange"
                type="datetimerange"
                range-separator="至"
                start-placeholder="开始时间"
                end-placeholder="结束时间"
                format="YYYY-MM-DD HH:mm:ss"
                value-format="YYYY-MM-DD HH:mm:ss"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="24">
            <el-form-item label="查询语句">
              <el-input
                v-model="queryForm.query"
                type="textarea"
                placeholder="请输入查询语句，例如：status:500 AND message:error"
                :rows="3"
              />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      
      <el-table 
        :data="logs" 
        style="width: 100%" 
        stripe
        v-loading="loading"
      >
        <el-table-column prop="@timestamp" label="时间" width="200" />
        <el-table-column prop="message" label="日志内容" min-width="400" />
        <el-table-column prop="level" label="日志级别" width="100" />
        <el-table-column prop="hostname" label="主机名" width="150" />
        <el-table-column prop="service" label="服务名称" width="150" />
        <el-table-column prop="status" label="状态码" width="80" />
      </el-table>
      
      <div class="pagination-container">
        <el-pagination
          background
          layout="prev, pager, next, jumper, ->, total"
          :total="total"
          :page-size="10"
          :current-page="currentPage"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Search } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

// 加载状态
const loading = ref(false)

// 查询表单
const queryForm = ref({
  project: '',
  logstore: '',
  timeRange: [],
  query: ''
})

// 日志仓库列表
const projects = ref(['project1', 'project2', 'project3'])

// 日志主题列表
const logstores = ref(['app-log', 'access-log', 'error-log'])

// 定义日志对象类型
interface LogItem {
  '@timestamp': string
  message: string
  level: string
  hostname: string
  service: string
  status: number
}

// 日志数据
const logs = ref<LogItem[]>([])

// 分页数据
const total = ref(0)
const currentPage = ref(1)

// 查询日志
const queryLogs = () => {
  if (!queryForm.value.project || !queryForm.value.logstore) {
    ElMessage.warning('请选择日志仓库和日志主题')
    return
  }
  
  loading.value = true
  // 这里可以添加实际的SLS日志查询逻辑
  setTimeout(() => {
    logs.value = [
      {
        '@timestamp': '2026-02-05 10:00:00',
        message: '请求处理成功',
        level: 'INFO',
        hostname: 'web-01',
        service: 'app-service',
        status: 200
      },
      {
        '@timestamp': '2026-02-05 10:00:01',
        message: '数据库连接失败',
        level: 'ERROR',
        hostname: 'db-01',
        service: 'db-service',
        status: 500
      }
    ]
    total.value = 100
    loading.value = false
    ElMessage.success('日志查询成功')
  }, 1000)
}

// 分页变化处理
const handleCurrentChange = (page: number) => {
  currentPage.value = page
  queryLogs()
}
</script>

<style scoped>
.sls-logs-view {
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

.sls-logs-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.query-form {
  margin-bottom: 20px;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>