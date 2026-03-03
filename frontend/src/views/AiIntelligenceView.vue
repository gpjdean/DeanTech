<template>
  <div class="ai-intelligence-view">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <el-icon class="header-icon"><ChatLineRound /></el-icon>
        <div class="header-text">
          <h2>AI智能助手</h2>
          <p>支持对接豆包、DeepSeek、千问和元宝等AI模型</p>
        </div>
        <div class="header-stats">
          <div class="stat-item">
            <span class="stat-value">{{ models.length }}</span>
            <span class="stat-label">可用模型</span>
          </div>
          <div class="stat-item">
            <span class="stat-value">{{ messages.length }}</span>
            <span class="stat-label">聊天记录</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 快速操作卡片 -->
    <div class="quick-actions">
      <el-card shadow="hover" class="quick-action-card" :class="{'active': activeTab === 'doubao'}">
        <div class="quick-action-content">
          <el-icon class="action-icon"><Cpu /></el-icon>
          <div class="action-text">
            <h3>模型配置</h3>
            <p>管理AI模型的API Key和参数设置</p>
          </div>
          <el-button type="primary" size="small" :icon="Setting" @click="switchTab('doubao')" circle>
            配置
          </el-button>
        </div>
      </el-card>
      
      <el-card shadow="hover" class="quick-action-card">
        <div class="quick-action-content">
          <el-icon class="action-icon"><Message /></el-icon>
          <div class="action-text">
            <h3>智能聊天</h3>
            <p>与AI模型进行实时对话交流</p>
          </div>
          <el-button type="success" size="small" :icon="ChatLineRound" @click="scrollToChat" circle>
            聊天
          </el-button>
        </div>
      </el-card>
      
      <el-card shadow="hover" class="quick-action-card">
        <div class="quick-action-content">
          <el-icon class="action-icon"><DataAnalysis /></el-icon>
          <div class="action-text">
            <h3>AI应用</h3>
            <p>探索AI的各种应用场景</p>
          </div>
          <el-button type="warning" size="small" :icon="MagicStick" @click="showAIApplications" circle>
            探索
          </el-button>
        </div>
      </el-card>
      
      <el-card shadow="hover" class="quick-action-card">
        <div class="quick-action-content">
          <el-icon class="action-icon"><RefreshRight /></el-icon>
          <div class="action-text">
            <h3>状态监控</h3>
            <p>查看AI模型的连接状态和使用情况</p>
          </div>
          <el-button type="info" size="small" :icon="Monitor" @click="checkAllConnections" circle>
            监控
          </el-button>
        </div>
      </el-card>
    </div>

    <!-- AI模型配置 -->
    <el-card shadow="hover" class="config-card">
      <template #header>
        <div class="card-header">
          <el-icon class="header-icon"><Setting /></el-icon>
          <span>AI模型配置</span>
          <el-tag v-if="activeConfigStatus" :type="activeConfigStatus" size="small" class="config-status-tag">
            {{ activeConfigStatus === 'success' ? '已配置' : '未配置' }}
          </el-tag>
        </div>
      </template>

      <el-tabs v-model="activeTab" class="ai-tabs" ref="aiTabs" @tab-change="handleTabChange">
        <el-tab-pane label="豆包" name="doubao">
          <el-form ref="doubaoForm" :model="aiConfig.doubao" label-width="120px" class="ai-form">
            <el-form-item label="API Key" prop="apiKey">
              <el-input v-model="aiConfig.doubao.apiKey" type="password" placeholder="请输入豆包API Key" show-password />
            </el-form-item>
            <el-form-item label="模型名称" prop="model">
              <el-select v-model="aiConfig.doubao.model" placeholder="请选择豆包模型">
                <el-option label="豆包-4" value="doubao-pro-4" />
                <el-option label="豆包-3" value="doubao-pro-3" />
                <el-option label="豆包-2" value="doubao-pro-2" />
              </el-select>
            </el-form-item>
            <el-form-item label="温度" prop="temperature">
              <div class="slider-container">
                <el-slider v-model="aiConfig.doubao.temperature" :min="0" :max="1" :step="0.1" />
                <span class="slider-value">{{ aiConfig.doubao.temperature }}</span>
              </div>
            </el-form-item>
            <el-form-item label="最大 tokens" prop="maxTokens">
              <el-input-number v-model="aiConfig.doubao.maxTokens" :min="100" :max="8192" :step="100" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveConfig('doubao')" :icon="Check">保存配置</el-button>
              <el-button @click="testConnection('doubao')" :icon="CircleCheck">测试连接</el-button>
              <el-button @click="resetConfig('doubao')" :icon="Refresh">重置</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="DeepSeek" name="deepseek">
          <el-form ref="deepseekForm" :model="aiConfig.deepseek" label-width="120px" class="ai-form">
            <el-form-item label="API Key" prop="apiKey">
              <el-input v-model="aiConfig.deepseek.apiKey" type="password" placeholder="请输入DeepSeek API Key" show-password />
            </el-form-item>
            <el-form-item label="模型名称" prop="model">
              <el-select v-model="aiConfig.deepseek.model" placeholder="请选择DeepSeek模型">
                <el-option label="DeepSeek-R1" value="deepseek-r1" />
                <el-option label="DeepSeek-V2" value="deepseek-v2" />
                <el-option label="DeepSeek-Coder" value="deepseek-coder" />
              </el-select>
            </el-form-item>
            <el-form-item label="温度" prop="temperature">
              <div class="slider-container">
                <el-slider v-model="aiConfig.deepseek.temperature" :min="0" :max="1" :step="0.1" />
                <span class="slider-value">{{ aiConfig.deepseek.temperature }}</span>
              </div>
            </el-form-item>
            <el-form-item label="最大 tokens" prop="maxTokens">
              <el-input-number v-model="aiConfig.deepseek.maxTokens" :min="100" :max="8192" :step="100" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveConfig('deepseek')" :icon="Check">保存配置</el-button>
              <el-button @click="testConnection('deepseek')" :icon="CircleCheck">测试连接</el-button>
              <el-button @click="resetConfig('deepseek')" :icon="Refresh">重置</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="千问" name="qianwen">
          <el-form ref="qianwenForm" :model="aiConfig.qianwen" label-width="120px" class="ai-form">
            <el-form-item label="API Key" prop="apiKey">
              <el-input v-model="aiConfig.qianwen.apiKey" type="password" placeholder="请输入千问API Key" show-password />
            </el-form-item>
            <el-form-item label="模型名称" prop="model">
              <el-select v-model="aiConfig.qianwen.model" placeholder="请选择千问模型">
                <el-option label="千问-4" value="qianwen-4" />
                <el-option label="千问-3" value="qianwen-3" />
                <el-option label="千问-2" value="qianwen-2" />
              </el-select>
            </el-form-item>
            <el-form-item label="温度" prop="temperature">
              <div class="slider-container">
                <el-slider v-model="aiConfig.qianwen.temperature" :min="0" :max="1" :step="0.1" />
                <span class="slider-value">{{ aiConfig.qianwen.temperature }}</span>
              </div>
            </el-form-item>
            <el-form-item label="最大 tokens" prop="maxTokens">
              <el-input-number v-model="aiConfig.qianwen.maxTokens" :min="100" :max="8192" :step="100" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveConfig('qianwen')" :icon="Check">保存配置</el-button>
              <el-button @click="testConnection('qianwen')" :icon="CircleCheck">测试连接</el-button>
              <el-button @click="resetConfig('qianwen')" :icon="Refresh">重置</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="元宝" name="yuanbao">
          <el-form ref="yuanbaoForm" :model="aiConfig.yuanbao" label-width="120px" class="ai-form">
            <el-form-item label="API Key" prop="apiKey">
              <el-input v-model="aiConfig.yuanbao.apiKey" type="password" placeholder="请输入元宝API Key" show-password />
            </el-form-item>
            <el-form-item label="模型名称" prop="model">
              <el-select v-model="aiConfig.yuanbao.model" placeholder="请选择元宝模型">
                <el-option label="元宝-4" value="yuanbao-4" />
                <el-option label="元宝-3" value="yuanbao-3" />
                <el-option label="元宝-2" value="yuanbao-2" />
              </el-select>
            </el-form-item>
            <el-form-item label="温度" prop="temperature">
              <div class="slider-container">
                <el-slider v-model="aiConfig.yuanbao.temperature" :min="0" :max="1" :step="0.1" />
                <span class="slider-value">{{ aiConfig.yuanbao.temperature }}</span>
              </div>
            </el-form-item>
            <el-form-item label="最大 tokens" prop="maxTokens">
              <el-input-number v-model="aiConfig.yuanbao.maxTokens" :min="100" :max="8192" :step="100" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveConfig('yuanbao')" :icon="Check">保存配置</el-button>
              <el-button @click="testConnection('yuanbao')" :icon="CircleCheck">测试连接</el-button>
              <el-button @click="resetConfig('yuanbao')" :icon="Refresh">重置</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- AI聊天助手 -->
    <el-card shadow="hover" class="chat-card">
      <template #header>
        <div class="card-header">
          <el-icon class="header-icon"><ChatLineRound /></el-icon>
          <span>AI聊天助手</span>
          <div class="header-right">
            <el-select v-model="selectedModel" placeholder="选择AI模型" size="small" class="model-select">
              <el-option label="豆包" value="doubao" />
              <el-option label="DeepSeek" value="deepseek" />
              <el-option label="千问" value="qianwen" />
              <el-option label="元宝" value="yuanbao" />
            </el-select>
            <el-button type="text" size="small" @click="showChatSettings" :icon="Setting" circle />
          </div>
        </div>
      </template>

      <div class="chat-container">
        <!-- 聊天消息区域 -->
        <div class="chat-messages" ref="chatMessages">
          <div v-for="(message, index) in messages" :key="index" :class="['message-item', message.role]">
            <div class="message-avatar">
              <el-avatar v-if="message.role === 'user'" :icon="User" class="user-avatar" />
              <el-avatar v-else :icon="ChatLineRound" class="ai-avatar" />
            </div>
            <div class="message-content">
              <div class="message-text">{{ message.content }}</div>
              <div class="message-time">{{ message.timestamp }}</div>
            </div>
          </div>
          <div v-if="sending" class="typing-indicator">
            <div class="typing-dots">
              <span></span>
              <span></span>
              <span></span>
            </div>
            <span class="typing-text">AI正在思考...</span>
          </div>
        </div>

        <!-- 输入区域 -->
        <div class="chat-input-area">
          <div class="input-toolbar">
            <el-button type="text" size="small" @click="showQuickPrompts" :class="{'active': showingQuickPrompts}">
              <el-icon><Document /></el-icon> 快捷提示
            </el-button>
            <el-button type="text" size="small" @click="clearChat" :icon="Delete">
              清空聊天
            </el-button>
            <el-button type="text" size="small" @click="exportChat" :icon="Download">
              导出聊天
            </el-button>
          </div>
          <el-input
            v-model="inputMessage"
            type="textarea"
            placeholder="请输入您的问题... (Ctrl+Enter 发送)"
            :rows="3"
            resize="none"
            @keyup.enter.ctrl="sendMessage"
            class="chat-input"
          />
          <div class="chat-actions">
            <el-button @click="clearInput" size="small" :icon="CircleClose">取消</el-button>
            <el-button type="primary" @click="sendMessage" :loading="sending" size="default">
              发送
            </el-button>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 快捷提示面板 -->
    <el-drawer
      v-model="showingQuickPrompts"
      title="快捷提示"
      direction="left"
      :with-header="true"
      :size="'300px'"
    >
      <div class="quick-prompts-container">
        <h3>常用提示词</h3>
        <div class="quick-prompts-grid">
          <el-button
            v-for="prompt in quickPrompts"
            :key="prompt.id"
            type="text"
            class="quick-prompt-item"
            @click="useQuickPrompt(prompt.text)"
          >
            <el-icon class="prompt-icon"><Document /></el-icon>
            <span>{{ prompt.text }}</span>
          </el-button>
        </div>
      </div>
    </el-drawer>

    <!-- 聊天设置对话框 -->
    <el-dialog v-model="showingChatSettings" title="聊天设置" width="400px">
      <el-form :model="chatSettings" label-width="100px">
        <el-form-item label="自动滚动">
          <el-switch v-model="chatSettings.autoScroll" />
        </el-form-item>
        <el-form-item label="消息时间">
          <el-switch v-model="chatSettings.showTimestamp" />
        </el-form-item>
        <el-form-item label="主题">
          <el-select v-model="chatSettings.theme" placeholder="选择主题">
            <el-option label="浅色" value="light" />
            <el-option label="深色" value="dark" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showingChatSettings = false">取消</el-button>
          <el-button type="primary" @click="saveChatSettings">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, nextTick, computed } from 'vue'
import { ElMessage, ElNotification, ElMessageBox } from 'element-plus'
import { 
  User, ChatLineRound, Setting, Cpu, 
  DataAnalysis, Document, Delete, Check, CircleCheck, 
  Refresh, MagicStick, RefreshRight, Monitor, Download, 
  CircleClose 
} from '@element-plus/icons-vue'
import { systemSettingAPI, aiAPI } from '@/api/api'

// 表单引用
const doubaoForm = ref()
const deepseekForm = ref()
const qianwenForm = ref()
const yuanbaoForm = ref()
const chatMessages = ref<HTMLElement | null>(null)
const aiTabs = ref<any>()

// 选项卡激活状态
const activeTab = ref('doubao')

// 发送状态
const sending = ref(false)

// 选择的AI模型
const selectedModel = ref('doubao')

// 聊天消息
const messages = ref([
  { role: 'assistant', content: '您好，我是AI智能助手，请问有什么可以帮助您的？', timestamp: new Date().toLocaleString() }
])

// 输入消息
const inputMessage = ref('')

// 快捷提示面板显示状态
const showingQuickPrompts = ref(false)

// 聊天设置对话框显示状态
const showingChatSettings = ref(false)

// 可用模型列表
const models = ref(['doubao', 'deepseek', 'qianwen', 'yuanbao'])

// 快捷提示词列表
const quickPrompts = ref([
  { id: 1, text: '请解释一下Kubernetes的工作原理' },
  { id: 2, text: '如何优化Linux系统性能？' },
  { id: 3, text: '请生成一个Python脚本，用于监控CPU使用率' },
  { id: 4, text: '什么是容器化技术？' },
  { id: 5, text: '如何设计高可用的分布式系统？' },
  { id: 6, text: '请解释一下RESTful API设计原则' },
  { id: 7, text: '如何进行代码性能优化？' },
  { id: 8, text: '请推荐一些学习Go语言的资源' },
  { id: 9, text: '如何排查网络连接问题？' },
  { id: 10, text: '什么是微服务架构？' }
])

// AI配置数据
const aiConfig = reactive({
  doubao: {
    apiKey: '',
    model: 'doubao-pro-4',
    temperature: 0.7,
    maxTokens: 2048
  },
  deepseek: {
    apiKey: '',
    model: 'deepseek-r1',
    temperature: 0.7,
    maxTokens: 2048
  },
  qianwen: {
    apiKey: '',
    model: 'qianwen-4',
    temperature: 0.7,
    maxTokens: 2048
  },
  yuanbao: {
    apiKey: '',
    model: 'yuanbao-4',
    temperature: 0.7,
    maxTokens: 2048
  }
})

// 聊天设置
const chatSettings = reactive({
  autoScroll: true,
  showTimestamp: true,
  theme: 'light'
})

// 计算当前活跃配置的状态
const activeConfigStatus = computed(() => {
  const config = aiConfig[activeTab.value as keyof typeof aiConfig]
  return config.apiKey ? 'success' : 'warning'
})

// 组件挂载时加载配置
onMounted(() => {
  loadConfig()
  loadChatSettings()
})

// 加载配置
const loadConfig = async () => {
  try {
    // 从后端加载所有系统设置
    const response = await systemSettingAPI.list()
    const settings = response.data
    
    // 查找AI配置
    const aiSetting = settings.find((setting: any) => setting.key === 'aiConfig')
    if (aiSetting) {
      const config = JSON.parse(aiSetting.value)
      // 合并配置，确保新字段有默认值
      Object.keys(aiConfig).forEach(key => {
        if (config[key]) {
          aiConfig[key as keyof typeof aiConfig] = {
            ...aiConfig[key as keyof typeof aiConfig],
            ...config[key]
          }
        }
      })
    }
  } catch (error) {
    console.error('加载AI配置失败:', error)
    ElMessage.error('加载AI配置失败')
  }
}

// 保存配置
const saveConfig = async (model: string) => {
  try {
    // 准备配置数据
    const configData = {
      key: 'aiConfig',
      value: JSON.stringify(aiConfig),
      description: 'AI智能助手配置'
    }
    
      try {
      // 尝试更新现有配置
      await systemSettingAPI.update('aiConfig', configData)
    } catch (error: any) {
      // 如果配置不存在，则创建新配置
      if (error.response && error.response.status === 404) {
        await systemSettingAPI.create(configData)
      } else {
        throw error
      }
    }
    
    ElMessage.success(`${model}配置保存成功`)
  } catch (error) {
    console.error('保存AI配置失败:', error)
    ElMessage.error('保存AI配置失败')
  }
}

// 重置配置
const resetConfig = (model: string) => {
  ElMessageBox.confirm('确定要重置该模型的配置吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    aiConfig[model as keyof typeof aiConfig] = {
      apiKey: '',
      model: `${model}-pro-4`,
      temperature: 0.7,
      maxTokens: 2048
    }
    ElMessage.success(`${model}配置已重置`)
  }).catch(() => {
    // 取消重置操作
  })
}

// 测试连接
const testConnection = async (model: string) => {
  try {
    ElMessage.info(`正在测试${model}连接...`)
    
    // 调用后端测试连接API
    const response = await aiAPI.testConnection({
      model: model,
      apiKey: aiConfig[model as keyof typeof aiConfig].apiKey
    })
    
    ElNotification({
      title: '连接测试',
      message: response.data.message,
      type: 'success',
      duration: 3000
    })
  } catch (error: any) {
    console.error('测试连接失败:', error)
    ElMessage.error(error.response?.data?.error || '测试连接失败')
  }
}

// 检查所有模型连接状态
const checkAllConnections = async () => {
  try {
    const results = await Promise.all(
      models.value.map(async (model) => {
        const apiKey = aiConfig[model as keyof typeof aiConfig].apiKey
        if (!apiKey) {
          return { model, status: '未配置', message: 'API Key未设置' }
        }
        try {
          const response = await aiAPI.testConnection({
            model: model,
            apiKey: apiKey
          })
          return { model, status: '成功', message: response.data.message }
        } catch (error: any) {
          return { model, status: '失败', message: error.response?.data?.error || '连接失败' }
        }
      })
    )
    
    // 显示结果
    const successCount = results.filter(r => r.status === '成功').length
    const failedCount = results.filter(r => r.status === '失败').length
    const unconfiguredCount = results.filter(r => r.status === '未配置').length
    
    ElMessageBox.alert(
      `<div style="text-align: left;">
        <h4>模型连接状态检查结果：</h4>
        <ul>
          <li>成功：${successCount}个</li>
          <li>失败：${failedCount}个</li>
          <li>未配置：${unconfiguredCount}个</li>
        </ul>
        <h5>详细信息：</h5>
        <ul>
          ${results.map(r => 
            `<li style="color: ${r.status === '成功' ? '#67c23a' : r.status === '失败' ? '#f56c6c' : '#e6a23c'}">
              ${r.model}：${r.status} - ${r.message}
            </li>`
          ).join('')}
        </ul>
      </div>`,
      '模型连接状态',
      { dangerouslyUseHTMLString: true }
    )
  } catch (error) {
    console.error('检查模型连接状态失败:', error)
    ElMessage.error('检查模型连接状态失败')
  }
}

// 发送消息
const sendMessage = async () => {
  if (!inputMessage.value.trim()) {
    ElMessage.warning('请输入消息内容')
    return
  }

  // 添加用户消息
  const userMessage = {
    role: 'user',
    content: inputMessage.value,
    timestamp: new Date().toLocaleString()
  }
  messages.value.push(userMessage)
  inputMessage.value = ''

  // 滚动到底部
  await nextTick()
  scrollToBottom()

  // 发送消息到后端获取AI回复
  sending.value = true
  try {
    // 准备聊天历史
    const chatHistory = messages.value.map(msg => ({
      role: msg.role,
      content: msg.content
    }))
    
    // 调用后端AI聊天API
    const response = await aiAPI.chat({
      model: selectedModel.value,
      message: userMessage.content,
      history: chatHistory
    })
    
    // 添加AI回复
    const aiMessage = {
      role: 'assistant',
      content: response.data.response,
      timestamp: new Date().toLocaleString()
    }
    messages.value.push(aiMessage)
    sending.value = false
    scrollToBottom()
  } catch (error: any) {
    console.error('发送消息失败:', error)
    ElMessage.error(error.response?.data?.error || '发送消息失败')
    sending.value = false
  }
}

// 清空聊天
const clearChat = () => {
  ElMessageBox.confirm('确定要清空所有聊天记录吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    messages.value = [
      { role: 'assistant', content: '您好，我是AI智能助手，请问有什么可以帮助您的？', timestamp: new Date().toLocaleString() }
    ]
    ElMessage.success('聊天记录已清空')
  }).catch(() => {
    // 取消清空操作
  })
}

// 导出聊天记录
const exportChat = () => {
  if (messages.value.length <= 1) {
    ElMessage.warning('没有可导出的聊天记录')
    return
  }
  
  const chatContent = messages.value.map(msg => {
    return `${msg.role === 'user' ? '用户' : 'AI'} ${msg.timestamp}:\n${msg.content}\n`
  }).join('\n')
  
  const blob = new Blob([chatContent], { type: 'text/plain;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `ai-chat-${new Date().toISOString().slice(0, 19).replace(/:/g, '-')}.txt`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
  
  ElMessage.success('聊天记录已导出')
}

// 清空输入框
const clearInput = () => {
  inputMessage.value = ''
}

// 滚动到底部
const scrollToBottom = () => {
  if (chatMessages.value && chatSettings.autoScroll) {
    chatMessages.value.scrollTop = chatMessages.value.scrollHeight
  }
}

// 滚动到聊天区域
const scrollToChat = () => {
  const chatCard = document.querySelector('.chat-card')
  if (chatCard) {
    chatCard.scrollIntoView({ behavior: 'smooth' })
  }
}

// 显示AI应用
const showAIApplications = () => {
  ElMessage.info('AI应用功能正在开发中，敬请期待！')
}

// 显示快捷提示
const showQuickPrompts = () => {
  showingQuickPrompts.value = true
}

// 显示聊天设置
const showChatSettings = () => {
  showingChatSettings.value = true
}

// 保存聊天设置
const saveChatSettings = () => {
  // 保存聊天设置到localStorage
  localStorage.setItem('aiChatSettings', JSON.stringify(chatSettings))
  ElMessage.success('聊天设置已保存')
  showingChatSettings.value = false
}

// 加载聊天设置
const loadChatSettings = () => {
  const savedSettings = localStorage.getItem('aiChatSettings')
  if (savedSettings) {
    Object.assign(chatSettings, JSON.parse(savedSettings))
  }
}

// 使用快捷提示
const useQuickPrompt = (prompt: string) => {
  inputMessage.value = prompt
  showingQuickPrompts.value = false
}

// 切换标签页
const switchTab = (tabName: string) => {
  if (aiTabs.value) {
    aiTabs.value.setActiveName(tabName)
  }
}

// 处理选项卡切换
const handleTabChange = (tabName: string) => {
  activeTab.value = tabName
}
</script>

<style scoped>
.ai-intelligence-view {
  width: 100%;
  padding: 0;
  background-color: #f5f7fa;
  min-height: calc(100vh - 120px);
}

/* 页面头部样式 */
.page-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 40px 0;
  margin-bottom: 30px;
  border-radius: 0;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  position: relative;
  overflow: hidden;
}

.page-header::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -10%;
  width: 400px;
  height: 400px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 50%;
  filter: blur(40px);
  animation: float 6s ease-in-out infinite;
}

.page-header::after {
  content: '';
  position: absolute;
  bottom: -30%;
  left: -10%;
  width: 300px;
  height: 300px;
  background: rgba(255, 255, 255, 0.15);
  border-radius: 50%;
  filter: blur(30px);
  animation: float 8s ease-in-out infinite reverse;
}

@keyframes float {
  0%, 100% { transform: translateY(0px); }
  50% { transform: translateY(-20px); }
}

.header-content {
  display: flex;
  align-items: center;
  gap: 40px;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
  position: relative;
  z-index: 1;
}

.header-icon {
  font-size: 48px;
  background: rgba(255, 255, 255, 0.2);
  padding: 16px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  backdrop-filter: blur(10px);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.header-text {
  flex: 1;
}

.header-text h2 {
  margin: 0 0 8px 0;
  font-size: 32px;
  font-weight: 700;
  letter-spacing: -0.5px;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.header-text p {
  margin: 0;
  font-size: 16px;
  opacity: 0.9;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

/* 头部统计信息 */
.header-stats {
  display: flex;
  gap: 24px;
  background: rgba(255, 255, 255, 0.15);
  padding: 16px 24px;
  border-radius: 12px;
  backdrop-filter: blur(10px);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: white;
}

.stat-label {
  font-size: 12px;
  opacity: 0.9;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

/* 快速操作卡片样式 */
.quick-actions {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 24px;
  margin-bottom: 30px;
  padding: 0 24px;
  max-width: 1200px;
  margin: 0 auto 30px auto;
}

.quick-action-card {
  border-radius: 12px;
  background: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid #e4e7ed;
  overflow: hidden;
  position: relative;
}

.quick-action-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
  transform: scaleX(0);
  transition: transform 0.3s ease;
}

.quick-action-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  border-color: #409eff;
}

.quick-action-card:hover::before {
  transform: scaleX(1);
}

.quick-action-card.active {
  border-color: #409eff;
  box-shadow: 0 4px 16px rgba(64, 158, 255, 0.2);
}

.quick-action-card.active::before {
  transform: scaleX(1);
}

.quick-action-content {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 24px;
  position: relative;
  z-index: 1;
}

.action-icon {
  font-size: 32px;
  color: #409eff;
  background: #ecf5ff;
  padding: 16px;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.quick-action-card:hover .action-icon {
  transform: scale(1.1);
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.action-text {
  flex: 1;
}

.action-text h3 {
  margin: 0 0 4px 0;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.action-text p {
  margin: 0;
  font-size: 14px;
  color: #606266;
}

/* 卡片通用样式 */
.config-card,
.chat-card {
  max-width: 1200px;
  margin: 0 auto 30px auto;
  border-radius: 12px;
  background: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  border: 1px solid #e4e7ed;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
}

.config-card:hover,
.chat-card:hover {
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  border-color: #409eff;
}

/* 卡片头部样式 */
.card-header {
  display: flex;
  justify-content: flex-start;
  align-items: center;
  gap: 12px;
  margin-bottom: 0;
  padding: 20px 24px;
  border-bottom: 1px solid #f0f2f5;
  border-radius: 12px 12px 0 0;
  background: linear-gradient(135deg, #fafafa 0%, #f0f2f5 100%);
  position: relative;
}

.card-header::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(90deg, transparent 0%, #409eff 50%, transparent 100%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.config-card:hover .card-header::after,
.chat-card:hover .card-header::after {
  opacity: 1;
}

.card-header .header-icon {
  font-size: 20px;
  color: #409eff;
  background: transparent;
  padding: 0;
}

.card-header span {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

/* 配置状态标签 */
.config-status-tag {
  margin-left: 12px;
}

/* 卡片头部右侧区域 */
.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-left: auto;
}

/* AI模型选择样式 */
.model-select {
  margin-left: auto;
}

/* 选项卡样式 */
.ai-tabs {
  max-width: 100%;
}

/* 表单样式 */
.ai-form {
  margin-top: 20px;
  padding: 0 24px 24px 24px;
}

/* 滑块容器 */
.slider-container {
  display: flex;
  align-items: center;
  gap: 16px;
}

.slider-value {
  color: #606266;
  font-size: 14px;
  min-width: 30px;
  text-align: right;
}

/* 聊天容器样式 */
.chat-container {
  display: flex;
  flex-direction: column;
  height: 600px;
  padding: 24px;
}

/* 聊天消息区域样式 */
.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
  background-color: #fafafa;
  border-radius: 12px;
  margin-bottom: 20px;
  border: 1px solid #e4e7ed;
  scrollbar-width: thin;
  scrollbar-color: #c0c4cc #f0f2f5;
}

.chat-messages::-webkit-scrollbar {
  width: 8px;
}

.chat-messages::-webkit-scrollbar-track {
  background: #f0f2f5;
  border-radius: 4px;
}

.chat-messages::-webkit-scrollbar-thumb {
  background: #c0c4cc;
  border-radius: 4px;
  transition: background 0.3s ease;
}

.chat-messages::-webkit-scrollbar-thumb:hover {
  background: #909399;
}

.message-item {
  display: flex;
  margin-bottom: 24px;
  align-items: flex-start;
  animation: slideIn 0.3s ease;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.message-item.user {
  flex-direction: row-reverse;
}

.message-avatar {
  margin: 0 12px;
  flex-shrink: 0;
}

.user-avatar {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  transition: all 0.3s ease;
}

.ai-avatar {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  color: white;
  transition: all 0.3s ease;
}

.message-item:hover .user-avatar,
.message-item:hover .ai-avatar {
  transform: scale(1.1);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.message-content {
  max-width: 75%;
}

.message-text {
  padding: 14px 18px;
  border-radius: 18px;
  line-height: 1.6;
  word-wrap: break-word;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  position: relative;
}

.message-item.user .message-text {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-bottom-right-radius: 4px;
}

.message-item.assistant .message-text {
  background-color: white;
  color: #333;
  border: 1px solid #e4e7ed;
  border-bottom-left-radius: 4px;
}

.message-time {
  font-size: 12px;
  color: #909399;
  margin-top: 6px;
  text-align: right;
  opacity: 0.7;
  transition: opacity 0.3s ease;
}

.message-item:hover .message-time {
  opacity: 1;
}

.message-item.assistant .message-time {
  text-align: left;
}

/* 打字指示器 */
.typing-indicator {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 24px;
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.typing-dots {
  display: flex;
  gap: 6px;
}

.typing-dots span {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: #409eff;
  animation: typing 1.4s infinite ease-in-out both;
}

.typing-dots span:nth-child(1) {
  animation-delay: -0.32s;
}

.typing-dots span:nth-child(2) {
  animation-delay: -0.16s;
}

@keyframes typing {
  0%, 80%, 100% {
    transform: scale(0);
    opacity: 0.5;
  }
  40% {
    transform: scale(1);
    opacity: 1;
  }
}

.typing-text {
  color: #909399;
  font-size: 14px;
  font-style: italic;
}

/* 聊天输入区域样式 */
.chat-input-area {
  display: flex;
  flex-direction: column;
  gap: 16px;
  background: white;
  padding: 20px;
  border-radius: 12px;
  border: 1px solid #e4e7ed;
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.05);
}

.input-toolbar {
  display: flex;
  justify-content: flex-start;
  gap: 12px;
  border-bottom: 1px solid #f0f2f5;
  padding-bottom: 12px;
}

.input-toolbar .el-button {
  color: #606266;
  transition: all 0.3s ease;
}

.input-toolbar .el-button:hover {
  color: #409eff;
  background-color: rgba(64, 158, 255, 0.1);
}

.input-toolbar .el-button.active {
  color: #409eff;
  background-color: rgba(64, 158, 255, 0.1);
  border-color: #409eff;
}

.chat-input {
  border-radius: 8px;
  border: 1px solid #dcdfe6;
  transition: all 0.3s ease;
  resize: vertical;
  min-height: 80px;
}

.chat-input:focus {
  border-color: #409eff;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.2);
}

.chat-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

/* 快捷提示面板样式 */
.quick-prompts-container {
  padding: 20px;
}

.quick-prompts-container h3 {
  margin: 0 0 20px 0;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.quick-prompts-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 12px;
}

.quick-prompt-item {
  text-align: left;
  padding: 14px 18px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  transition: all 0.3s ease;
  color: #303133;
  background: white;
  text-align: left;
  display: flex;
  align-items: center;
  gap: 10px;
}

.quick-prompt-item .prompt-icon {
  font-size: 16px;
  color: #409eff;
  flex-shrink: 0;
}

.quick-prompt-item:hover {
  background: linear-gradient(135deg, #ecf5ff 0%, #f0f5ff 100%);
  border-color: #409eff;
  color: #409eff;
  transform: translateX(6px);
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.15);
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .header-content {
    flex-wrap: wrap;
    gap: 24px;
  }
  
  .header-stats {
    order: 3;
    width: 100%;
    justify-content: center;
  }
}

@media (max-width: 768px) {
  .header-content {
    flex-direction: column;
    text-align: center;
    gap: 16px;
  }

  .header-text h2 {
    font-size: 24px;
  }

  .quick-actions {
    grid-template-columns: 1fr;
    padding: 0 16px;
  }

  .quick-action-content {
    flex-direction: column;
    text-align: center;
  }

  .chat-messages {
    padding: 16px;
  }

  .message-content {
    max-width: 85%;
  }

  .message-item {
    margin-bottom: 20px;
  }

  .chat-container {
    padding: 16px;
    height: 500px;
  }

  .config-card,
  .chat-card {
    margin: 0 16px 24px 16px;
  }

  .ai-form {
    padding: 0 16px 16px 16px;
  }
  
  .header-stats {
    gap: 16px;
    padding: 12px 16px;
  }
  
  .stat-value {
    font-size: 20px;
  }
}
</style>
