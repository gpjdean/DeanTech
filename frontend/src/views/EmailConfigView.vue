<template>
  <div class="email-config-view">
    <el-card shadow="never" class="page-header">
      <div class="header-content">
        <h2>邮箱配置</h2>
      </div>
    </el-card>

    <el-card shadow="hover" class="config-card" v-loading="loading.emailConfig" element-loading-text="加载邮箱配置中...">
      <template #header>
        <div class="card-header">
          <span>邮箱配置</span>
        </div>
      </template>
      
      <el-alert
        type="info"
        show-icon
        class="email-hint-alert"
      >
        <ul class="email-hint-list">
          <li><strong>端口 25</strong>：选“无”加密，仅用于服务器间转发，易被屏蔽。</li>
          <li><strong>端口 465</strong>：选“SSL”加密，连接即安全，但部分服务已弃用。</li>
          <li><strong>端口 587</strong>：选“TLS”加密，推荐使用。先普通连接后升级，最通用可靠。</li>
        </ul>
      </el-alert>
      
      <el-form ref="emailConfigForm" :model="emailConfigFormData" label-width="120px" class="system-style-form">
        <el-form-item label="SMTP服务器" prop="smtpHost">
          <el-input v-model="emailConfigFormData.smtpHost" placeholder="请输入SMTP服务器地址" />
        </el-form-item>

        <el-form-item label="SMTP端口" prop="smtpPort">
          <el-input-number v-model="emailConfigFormData.smtpPort" :min="1" :max="65535" placeholder="请输入SMTP端口" style="width: 100%" />
        </el-form-item>

        <el-form-item label="发件人邮箱" prop="fromEmail">
          <el-input v-model="emailConfigFormData.fromEmail" placeholder="请输入发件人邮箱地址" />
        </el-form-item>

        <el-form-item label="SMTP用户名" prop="smtpUsername">
          <el-input v-model="emailConfigFormData.smtpUsername" placeholder="请输入SMTP用户名" />
        </el-form-item>

        <el-form-item label="SMTP密码" prop="smtpPassword">
          <el-input v-model="emailConfigFormData.smtpPassword" type="password" placeholder="请输入SMTP密码或授权码" show-password />
        </el-form-item>

        <el-form-item label="加密方式" prop="smtpEncryption">
          <el-select v-model="emailConfigFormData.smtpEncryption" style="width: 100%">
            <el-option label="无" value="" />
            <el-option label="SSL" value="ssl" />
            <el-option label="TLS" value="tls" />
          </el-select>
        </el-form-item>

        <el-form-item label="收件人邮箱" prop="testEmail" style="margin-bottom: 30px;">
          <el-input v-model="emailConfigFormData.testEmail" placeholder="请输入测试邮件收件人地址" />
          <div class="form-hint">用于测试邮件发送功能</div>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="saveEmailConfig" :loading="loading.emailConfig">保存配置</el-button>
          <el-button type="success" @click="testEmailSend" :loading="loading.testEmail">测试发送邮件</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { emailAPI, emailSettingAPI } from '../api/api'

// 表单引用
const emailConfigForm = ref()

// 加载状态
const loading = ref({
  emailConfig: false,
  testEmail: false
})

// 邮箱配置表单数据
const emailConfigFormData = ref({
  smtpHost: '',
  smtpPort: 465,
  fromEmail: '',
  smtpUsername: '',
  smtpPassword: '',
  smtpEncryption: 'ssl',
  testEmail: ''
})

// 组件挂载时加载数据
onMounted(() => {
  loadEmailConfig()
})

// 加载邮箱配置
const loadEmailConfig = async () => {
  try {
    loading.value.emailConfig = true
    const response = await emailSettingAPI.get()
    const emailConfig = response.data
    
    // 直接赋值邮箱配置
    emailConfigFormData.value.smtpHost = emailConfig.smtpServer || ''
    emailConfigFormData.value.smtpPort = emailConfig.smtpPort || 465
    emailConfigFormData.value.fromEmail = emailConfig.fromEmail || ''
    emailConfigFormData.value.smtpUsername = emailConfig.smtpUsername || ''
    emailConfigFormData.value.smtpPassword = emailConfig.smtpPassword || ''
    emailConfigFormData.value.smtpEncryption = emailConfig.smtpEncryption || 'ssl'
  } catch (error) {
    console.error('加载邮箱配置失败:', error)
    ElMessage.error('加载邮箱配置失败')
  } finally {
    loading.value.emailConfig = false
  }
}

// 保存邮箱配置
const saveEmailConfig = async () => {
  try {
    loading.value.emailConfig = true
    
    // 准备要保存的数据
    const emailConfigData = {
      smtpServer: emailConfigFormData.value.smtpHost,
      smtpPort: emailConfigFormData.value.smtpPort,
      fromEmail: emailConfigFormData.value.fromEmail,
      smtpUsername: emailConfigFormData.value.smtpUsername,
      smtpPassword: emailConfigFormData.value.smtpPassword,
      smtpEncryption: emailConfigFormData.value.smtpEncryption
    }
    
    // 保存邮箱配置
    await emailSettingAPI.update(emailConfigData)
    
    ElMessage.success('邮箱配置保存成功')
  } catch (error) {
    console.error('保存邮箱配置失败:', error)
    ElMessage.error('保存邮箱配置失败')
  } finally {
    loading.value.emailConfig = false
  }
}

// 测试发送邮件
const testEmailSend = async () => {
  try {
    if (!emailConfigFormData.value.testEmail) {
      ElMessage.warning('请输入测试邮件收件人地址')
      return
    }
    
    loading.value.testEmail = true
    
    // 调用邮件测试API
    await emailAPI.testSend(emailConfigFormData.value.testEmail)
    
    ElMessage.success('测试邮件发送成功，请查收邮箱')
  } catch (error) {
    console.error('发送测试邮件失败:', error)
    ElMessage.error('发送测试邮件失败')
  } finally {
    loading.value.testEmail = false
  }
}
</script>

<style scoped>
.email-config-view {
  width: 100%;
  padding: 0;
}

.page-header {
  margin-bottom: 30px;
}

.header-content h2 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: #333;
}

.config-card {
  margin-bottom: 30px;
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

/* 系统风格表单样式 */
.system-style-form {
  max-width: 600px;
}

.form-hint {
  margin-top: 8px;
  font-size: 12px;
  color: #909399;
}

.email-hint-alert {
  margin-bottom: 20px;
}

.email-hint-list {
  margin: 10px 0 0 20px;
  padding: 0;
  font-size: 14px;
  color: #606266;
}

.email-hint-list li {
  margin-bottom: 8px;
  line-height: 1.5;
}

.email-hint-list strong {
  color: #303133;
}
</style>