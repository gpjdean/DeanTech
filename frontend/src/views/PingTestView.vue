<template>
  <div class="ping-test-view">
    <el-card shadow="never" class="page-header">
      <div class="header-content">
        <h2>Ping测试</h2>
      </div>
    </el-card>

    <el-card shadow="hover" class="ping-card">
      <template #header>
        <div class="card-header">
          <span>Ping测试工具</span>
        </div>
      </template>

      <el-form :model="pingForm" :rules="pingRules" ref="pingFormRef" label-width="120px" class="ping-form">
        <el-form-item label="目标地址" prop="target">
          <el-input
            v-model="pingForm.target"
            placeholder="请输入目标IP地址或域名，例如：www.example.com 或 192.168.1.1"
            clearable
          />
        </el-form-item>

        <el-form-item label="测试次数" prop="count">
          <el-input-number
            v-model="pingForm.count"
            :min="1"
            :max="10"
            :step="1"
            placeholder="默认值：4"
          />
        </el-form-item>

        <el-form-item label="超时时间" prop="timeout">
          <el-input-number
            v-model="pingForm.timeout"
            :min="100"
            :max="5000"
            :step="100"
            placeholder="默认值：1000毫秒"
          />
        </el-form-item>

        <el-form-item label="数据包大小" prop="size">
          <el-input-number
            v-model="pingForm.size"
            :min="32"
            :max="1500"
            :step="8"
            placeholder="默认值：32字节"
          />
        </el-form-item>

        <el-form-item>
          <div class="form-actions">
            <el-button type="primary" @click="startPing" :loading="pingLoading">
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

      <el-card shadow="hover" class="result-card" v-if="pingResult.length > 0 || pingLoading">
        <template #header>
          <div class="card-header">
            <span>测试结果</span>
            <div class="result-stats" v-if="pingStats">
              <el-tag type="success" size="small" class="stat-tag">平均延迟: {{ pingStats.avgDelay }} ms</el-tag>
              <el-tag type="info" size="small" class="stat-tag">最小延迟: {{ pingStats.minDelay }} ms</el-tag>
              <el-tag type="warning" size="small" class="stat-tag">最大延迟: {{ pingStats.maxDelay }} ms</el-tag>
              <el-tag type="primary" size="small" class="stat-tag">丢包率: {{ pingStats.packetLoss }}%</el-tag>
            </div>
          </div>
        </template>

        <div class="result-content" v-loading="pingLoading">
          <div class="ping-row" v-for="(result, index) in pingResult" :key="index">
            <div class="ping-index">{{ index + 1 }}</div>
            <div class="ping-info">
              <div class="ping-status" :class="result.success ? 'success' : 'failed'">
                <el-icon v-if="result.success"><Check /></el-icon>
                <el-icon v-else><Close /></el-icon>
                {{ result.message }}
              </div>
              <div class="ping-details" v-if="result.success">
                <span>来自 {{ result.ip }} 的回复: 字节={{ result.bytes }} 时间={{ result.time }} ms TTL={{ result.ttl }}</span>
              </div>
            </div>
          </div>
          <div class="ping-summary" v-if="pingResult.length > 0">
            <el-divider />
            <div class="summary-content">
              <span class="summary-label">Ping测试 {{ pingStats.success ? '成功' : '失败' }}</span>
              <span class="summary-label">目标地址：{{ pingForm.target }}</span>
              <span class="summary-label">发送：{{ pingStats.sent }} 接收：{{ pingStats.received }} 丢失：{{ pingStats.lost }}</span>
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
const pingFormRef = ref()

// Ping测试表单
const pingForm = ref({
  target: '',
  count: 4,
  timeout: 1000,
  size: 32
})

// Ping测试规则
const pingRules = ref({
  target: [
    { required: true, message: '请输入目标地址', trigger: 'blur' }
  ],
  count: [
    { required: true, message: '请输入测试次数', trigger: 'blur' },
    { type: 'number', min: 1, max: 10, message: '测试次数范围在 1 到 10 之间', trigger: 'blur' }
  ],
  timeout: [
    { required: true, message: '请输入超时时间', trigger: 'blur' },
    { type: 'number', min: 100, max: 5000, message: '超时时间范围在 100 到 5000 毫秒之间', trigger: 'blur' }
  ],
  size: [
    { required: true, message: '请输入数据包大小', trigger: 'blur' },
    { type: 'number', min: 32, max: 1500, message: '数据包大小范围在 32 到 1500 字节之间', trigger: 'blur' }
  ]
})

// Ping测试结果
const pingResult = ref([])
// Ping测试加载状态
const pingLoading = ref(false)
// Ping测试统计信息
const pingStats = ref({
  success: false,
  sent: 0,
  received: 0,
  lost: 0,
  minDelay: 0,
  avgDelay: 0,
  maxDelay: 0,
  packetLoss: 0
})

// 开始Ping测试
const startPing = async () => {
  if (!pingFormRef.value) return

  pingFormRef.value.validate(async (valid) => {
    if (valid) {
      pingLoading.value = true
      pingResult.value = []
      // 初始化统计信息
      pingStats.value = {
        success: false,
        sent: 0,
        received: 0,
        lost: 0,
        minDelay: 0,
        avgDelay: 0,
        maxDelay: 0,
        packetLoss: 0
      }

      try {
        // 调用后端API执行Ping测试
        const response = await fetch('/api/ping/test', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            target: pingForm.value.target,
            count: pingForm.value.count,
            timeout: pingForm.value.timeout,
            size: pingForm.value.size
          })
        })

        const result = await response.json()
        pingResult.value = result.results
        pingStats.value = result.stats
        ElMessage.success('Ping测试完成')
      } catch (error) {
        console.error('Ping测试失败:', error)
        ElMessage.error('Ping测试失败')
      } finally {
        pingLoading.value = false
      }
    }
  })
}

// 重置表单
const resetForm = () => {
  if (pingFormRef.value) {
    pingFormRef.value.resetFields()
  }
}

// 清空结果
const clearResult = () => {
  pingResult.value = []
  pingStats.value = {
    success: false,
    sent: 0,
    received: 0,
    lost: 0,
    minDelay: 0,
    avgDelay: 0,
    maxDelay: 0,
    packetLoss: 0
  }
}
</script>

<style scoped>
.ping-test-view {
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

.ping-card {
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

.ping-form {
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

.ping-row {
  display: flex;
  align-items: flex-start;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
  transition: all 0.3s ease;
}

.ping-row:hover {
  background-color: #f5f7fa;
}

.ping-row:last-child {
  border-bottom: none;
}

.ping-index {
  width: 40px;
  font-weight: 600;
  color: #409eff;
  text-align: center;
  margin-top: 2px;
}

.ping-info {
  flex: 1;
  margin-left: 20px;
}

.ping-status {
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 4px;
  display: flex;
  align-items: center;
  gap: 4px;
}

.ping-status.success {
  color: #67c23a;
}

.ping-status.failed {
  color: #f56c6c;
}

.ping-details {
  font-size: 12px;
  color: #606266;
  margin-left: 22px;
}

.ping-summary {
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