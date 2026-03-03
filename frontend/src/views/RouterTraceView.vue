<template>
  <div class="trace-route-test-view">
    <el-card shadow="never" class="page-header">
      <div class="header-content">
        <h2>路由跟踪测试</h2>
      </div>
    </el-card>

    <el-card shadow="hover" class="trace-route-card">
      <template #header>
        <div class="card-header">
          <span>路由跟踪工具</span>
        </div>
      </template>

      <el-form :model="traceRouteForm" :rules="traceRouteRules" ref="traceRouteFormRef" label-width="120px" class="trace-route-form">
        <el-form-item label="目标地址" prop="target">
          <el-input
            v-model="traceRouteForm.target"
            placeholder="请输入目标IP地址或域名，例如：www.example.com 或 192.168.1.1"
            clearable
          />
        </el-form-item>

        <el-form-item label="最大跳数" prop="maxHops">
          <el-input-number
            v-model="traceRouteForm.maxHops"
            :min="1"
            :max="30"
            :step="1"
            placeholder="默认值：30"
          />
        </el-form-item>

        <el-form-item label="超时时间" prop="timeout">
          <el-input-number
            v-model="traceRouteForm.timeout"
            :min="100"
            :max="5000"
            :step="100"
            placeholder="默认值：1000毫秒"
          />
        </el-form-item>

        <el-form-item>
          <div class="form-actions">
            <el-button type="primary" @click="startTraceRoute" :loading="traceRouteLoading">
              <el-icon><Search /></el-icon>
              开始测试
            </el-button>
            <el-button @click="resetForm">
              <el-icon><Refresh /></el-icon>
              重置
            </el-button>
            <el-button @click="clearResult">
              <el-icon><Delete /></el-icon>
              清空结果
            </el-button>
          </div>
        </el-form-item>
      </el-form>

      <el-card shadow="hover" class="result-card" v-if="traceRouteResult.hops.length > 0 || traceRouteLoading">
        <template #header>
          <div class="card-header">
            <span>测试结果</span>
            <div class="result-stats" v-if="traceRouteResult.stats">
              <el-tag type="success" size="small" class="stat-tag" v-if="traceRouteResult.stats.reached">目标可达</el-tag>
              <el-tag type="danger" size="small" class="stat-tag" v-else>目标不可达</el-tag>
              <el-tag type="success" size="small" class="stat-tag">平均延迟: {{ traceRouteResult.stats.avgRTT }} ms</el-tag>
              <el-tag type="info" size="small" class="stat-tag">最小延迟: {{ traceRouteResult.stats.minRTT }} ms</el-tag>
              <el-tag type="warning" size="small" class="stat-tag">最大延迟: {{ traceRouteResult.stats.maxRTT }} ms</el-tag>
              <el-tag type="primary" size="small" class="stat-tag">总跳数: {{ traceRouteResult.stats.totalHops }}</el-tag>
            </div>
          </div>
        </template>

        <div class="result-content" v-loading="traceRouteLoading">
          <div class="trace-route-hop" v-for="hop in traceRouteResult.hops" :key="hop.hop">
            <div class="hop-index">{{ hop.hop }}</div>
            <div class="hop-info">
              <div class="hop-status" :class="hop.success ? 'success' : 'failed'">
                <el-icon v-if="hop.success"><Check /></el-icon>
                <el-icon v-else><Close /></el-icon>
                {{ hop.success ? '成功' : '超时' }}
              </div>
              <div class="hop-details" v-if="hop.success">
                <span>{{ hop.host }} ({{ hop.ip }})</span>
                <span class="hop-time">时间: {{ hop.time }} ms</span>
                <span class="hop-ttl">TTL: {{ hop.ttl }}</span>
              </div>
            </div>
          </div>
          <div class="trace-route-summary" v-if="traceRouteResult.hops.length > 0">
            <el-divider />
            <div class="summary-content">
              <span class="summary-label">路由跟踪 {{ traceRouteResult.stats.success ? '成功' : '失败' }}</span>
              <span class="summary-label">目标地址：{{ traceRouteForm.target }}</span>
              <span class="summary-label">状态：{{ traceRouteResult.stats.reached ? '目标可达' : '目标不可达' }}</span>
            </div>
          </div>
        </div>
      </el-card>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Refresh, Delete, Check, Close } from '@element-plus/icons-vue'

// 表单引用
const traceRouteFormRef = ref()

// 路由跟踪表单
const traceRouteForm = ref({
  target: '',
  maxHops: 30,
  timeout: 1000
})

// 路由跟踪规则
const traceRouteRules = ref({
  target: [
    { required: true, message: '请输入目标地址', trigger: 'blur' }
  ],
  maxHops: [
    { required: true, message: '请输入最大跳数', trigger: 'blur' },
    { type: 'number', min: 1, max: 30, message: '最大跳数范围在 1 到 30 之间', trigger: 'blur' }
  ],
  timeout: [
    { required: true, message: '请输入超时时间', trigger: 'blur' },
    { type: 'number', min: 100, max: 5000, message: '超时时间范围在 100 到 5000 毫秒之间', trigger: 'blur' }
  ]
})

// 路由跟踪结果
const traceRouteResult = ref({
  target: '',
  hops: [],
  stats: {
    success: false,
    totalHops: 0,
    reached: false,
    maxRTT: 0,
    minRTT: 0,
    avgRTT: 0
  }
})
// 路由跟踪加载状态
const traceRouteLoading = ref(false)

// 开始路由跟踪测试
const startTraceRoute = async () => {
  if (!traceRouteFormRef.value) return

  traceRouteFormRef.value.validate(async (valid) => {
    if (valid) {
      traceRouteLoading.value = true
      // 初始化结果
      traceRouteResult.value = {
        target: traceRouteForm.value.target,
        hops: [],
        stats: {
          success: false,
          totalHops: 0,
          reached: false,
          maxRTT: 0,
          minRTT: 0,
          avgRTT: 0
        }
      }

      try {
        // 调用后端API执行路由跟踪测试
        const response = await fetch('/api/trace-route/test', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            target: traceRouteForm.value.target,
            maxHops: traceRouteForm.value.maxHops,
            timeout: traceRouteForm.value.timeout
          })
        })

        const result = await response.json()
        traceRouteResult.value = result
        ElMessage.success('路由跟踪测试完成')
      } catch (error) {
        console.error('路由跟踪测试失败:', error)
        ElMessage.error('路由跟踪测试失败')
      } finally {
        traceRouteLoading.value = false
      }
    }
  })
}

// 重置表单
const resetForm = () => {
  if (traceRouteFormRef.value) {
    traceRouteFormRef.value.resetFields()
  }
}

// 清空结果
const clearResult = () => {
  traceRouteResult.value = {
    target: '',
    hops: [],
    stats: {
      success: false,
      totalHops: 0,
      reached: false,
      maxRTT: 0,
      minRTT: 0,
      avgRTT: 0
    }
  }
}
</script>

<style scoped>
.trace-route-test-view {
  width: 100%;
  padding: 0;
}

.page-header {
  margin-bottom: 20px;
}

.header-content h2 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: #333;
}

.trace-route-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.card-header span {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.result-stats {
  display: flex;
  gap: 10px;
}

.stat-tag {
  margin-right: 8px;
}

.trace-route-form {
  margin-bottom: 20px;
}

.form-actions {
  display: flex;
  gap: 12px;
}

.result-card {
  margin-top: 20px;
}

.result-content {
  padding: 10px 0;
}

.trace-route-hop {
  display: flex;
  align-items: flex-start;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
  transition: all 0.3s ease;
}

.trace-route-hop:hover {
  background-color: #f5f7fa;
}

.trace-route-hop:last-child {
  border-bottom: none;
}

.hop-index {
  width: 40px;
  font-weight: 600;
  color: #409eff;
  text-align: center;
  margin-top: 2px;
}

.hop-info {
  flex: 1;
  margin-left: 20px;
}

.hop-status {
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 4px;
  display: flex;
  align-items: center;
  gap: 4px;
}

.hop-status.success {
  color: #67c23a;
}

.hop-status.failed {
  color: #f56c6c;
}

.hop-details {
  font-size: 12px;
  color: #606266;
  margin-left: 22px;
  display: flex;
  flex-wrap: wrap;
  gap: 15px;
}

.hop-time {
  color: #409eff;
}

.hop-ttl {
  color: #e6a23c;
}

.trace-route-summary {
  margin-top: 20px;
}

.summary-content {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  padding: 15px;
  background-color: #f0f9ff;
  border-radius: 8px;
}

.summary-label {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}
</style>
