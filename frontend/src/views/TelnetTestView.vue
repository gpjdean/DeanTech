<template>
  <div class="telnet-test-view">
    <el-card shadow="never" class="page-header">
      <div class="header-content">
        <h2>Telnet测试</h2>
      </div>
    </el-card>

    <el-card shadow="hover" class="telnet-card">
      <template #header>
        <div class="card-header">
          <span>Telnet测试工具</span>
        </div>
      </template>

      <el-form :model="telnetForm" :rules="telnetRules" ref="telnetFormRef" label-width="120px" class="telnet-form">
        <el-form-item label="目标地址" prop="target">
          <el-input
            v-model="telnetForm.target"
            placeholder="请输入目标IP地址或域名，例如：www.example.com 或 192.168.1.1"
            clearable
          />
        </el-form-item>

        <el-form-item label="端口号" prop="port">
          <el-input-number
            v-model="telnetForm.port"
            :min="1"
            :max="65535"
            :step="1"
            placeholder="默认值：23"
          />
        </el-form-item>

        <el-form-item label="超时时间" prop="timeout">
          <el-input-number
            v-model="telnetForm.timeout"
            :min="100"
            :max="5000"
            :step="100"
            placeholder="默认值：1000毫秒"
          />
        </el-form-item>

        <el-form-item>
          <div class="form-actions">
            <el-button type="primary" @click="startTelnet" :loading="telnetLoading">
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

      <el-card shadow="hover" class="result-card" v-if="telnetResult || telnetLoading">
        <template #header>
          <div class="card-header">
            <span>测试结果</span>
          </div>
        </template>

        <div class="result-content" v-loading="telnetLoading">
          <div class="result-status" :class="telnetResult.success ? 'success' : 'failed'">
            <el-icon v-if="telnetResult.success"><Check /></el-icon>
            <el-icon v-else><Close /></el-icon>
            {{ telnetResult.message }}
          </div>
          <div class="result-details" v-if="telnetResult.details">
            <el-divider />
            <h4>详细信息：</h4>
            <pre class="result-pre">{{ telnetResult.details }}</pre>
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
const telnetFormRef = ref()

// Telnet测试表单
const telnetForm = ref({
  target: '',
  port: 23,
  timeout: 1000
})

// Telnet测试规则
const telnetRules = ref({
  target: [
    { required: true, message: '请输入目标地址', trigger: 'blur' }
  ],
  port: [
    { required: true, message: '请输入端口号', trigger: 'blur' },
    { type: 'number', min: 1, max: 65535, message: '端口号范围在 1 到 65535 之间', trigger: 'blur' }
  ],
  timeout: [
    { required: true, message: '请输入超时时间', trigger: 'blur' },
    { type: 'number', min: 100, max: 5000, message: '超时时间范围在 100 到 5000 毫秒之间', trigger: 'blur' }
  ]
})

// Telnet测试结果
const telnetResult = ref({
  success: false,
  message: '',
  details: ''
})
// Telnet测试加载状态
const telnetLoading = ref(false)

// 开始Telnet测试
const startTelnet = async () => {
  if (!telnetFormRef.value) return

  telnetFormRef.value.validate(async (valid) => {
    if (valid) {
      telnetLoading.value = true
      // 初始化结果
      telnetResult.value = {
        success: false,
        message: '',
        details: ''
      }

      try {
        // 调用后端API执行Telnet测试
        const response = await fetch('/api/telnet/test', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            target: telnetForm.value.target,
            port: telnetForm.value.port,
            timeout: telnetForm.value.timeout
          })
        })

        const result = await response.json()
        telnetResult.value = {
          success: result.success,
          message: result.message,
          details: result.details
        }
        ElMessage.success('Telnet测试完成')
      } catch (error) {
        console.error('Telnet测试失败:', error)
        telnetResult.value = {
          success: false,
          message: '测试请求失败',
          details: '无法连接到服务器或服务器内部错误'
        }
        ElMessage.error('Telnet测试失败')
      } finally {
        telnetLoading.value = false
      }
    }
  })
}

// 重置表单
const resetForm = () => {
  if (telnetFormRef.value) {
    telnetFormRef.value.resetFields()
  }
  // 清空结果
  telnetResult.value = {
    success: false,
    message: '',
    details: ''
  }
}

// 清空结果
const clearResult = () => {
  telnetResult.value = {
    success: false,
    message: '',
    details: ''
  }
}
</script>

<style scoped>
.telnet-test-view {
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

.telnet-card {
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

.telnet-form {
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

.result-status {
  font-size: 16px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
}

.result-status.success {
  color: #67c23a;
}

.result-status.failed {
  color: #f56c6c;
}

.result-details {
  margin-top: 15px;
}

.result-details h4 {
  margin: 10px 0;
  font-size: 14px;
  font-weight: 600;
  color: #303133;
}

.result-pre {
  background-color: #f5f7fa;
  border: 1px solid #ebeef5;
  border-radius: 4px;
  padding: 15px;
  font-family: 'Courier New', Courier, monospace;
  font-size: 12px;
  line-height: 1.5;
  color: #606266;
  overflow-x: auto;
  white-space: pre-wrap;
  word-wrap: break-word;
}
</style>