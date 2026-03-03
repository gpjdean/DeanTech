<template>
  <div class="ai-code-assistant-view">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <el-icon class="header-icon"><Code /></el-icon>
        <div class="header-text">
          <h2>AI代码助手</h2>
          <p>智能代码生成、解释和优化工具</p>
        </div>
        <div class="header-actions">
          <el-button type="primary" size="small" :icon="Refresh" @click="resetEditor">
            重置
          </el-button>
        </div>
      </div>
    </div>

    <!-- 快速操作卡片 -->
    <div class="quick-actions">
      <el-card shadow="hover" class="quick-action-card" @click="setAction('generate')">
        <div class="quick-action-content">
          <el-icon class="action-icon"><DocumentAdd /></el-icon>
          <div class="action-text">
            <h3>代码生成</h3>
            <p>根据需求生成高质量代码</p>
          </div>
        </div>
      </el-card>
      
      <el-card shadow="hover" class="quick-action-card" @click="setAction('explain')">
        <div class="quick-action-content">
          <el-icon class="action-icon"><Document /></el-icon>
          <div class="action-text">
            <h3>代码解释</h3>
            <p>详细解释代码功能和实现原理</p>
          </div>
        </div>
      </el-card>
      
      <el-card shadow="hover" class="quick-action-card" @click="setAction('optimize')">
        <div class="quick-action-content">
          <el-icon class="action-icon"><Top /></el-icon>
          <div class="action-text">
            <h3>代码优化</h3>
            <p>优化代码性能和可读性</p>
          </div>
        </div>
      </el-card>
      
      <el-card shadow="hover" class="quick-action-card" @click="setAction('debug')">
        <div class="quick-action-content">
          <el-icon class="action-icon"><Warning /></el-icon>
          <div class="action-text">
            <h3>代码调试</h3>
            <p>帮助查找和修复代码错误</p>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 主内容区域 -->
    <el-card shadow="hover" class="main-card">
      <template #header>
        <div class="card-header">
          <el-icon class="header-icon"><Setting /></el-icon>
          <span>代码助手工作台</span>
          <div class="header-right">
            <el-select v-model="selectedLanguage" placeholder="选择编程语言" size="small" class="language-select">
              <el-option label="JavaScript" value="javascript" />
              <el-option label="Python" value="python" />
              <el-option label="Go" value="go" />
              <el-option label="Java" value="java" />
              <el-option label="TypeScript" value="typescript" />
              <el-option label="C++" value="cpp" />
              <el-option label="Ruby" value="ruby" />
              <el-option label="PHP" value="php" />
            </el-select>
            <el-select v-model="selectedModel" placeholder="选择AI模型" size="small" class="model-select">
              <el-option label="豆包" value="doubao" />
              <el-option label="DeepSeek" value="deepseek" />
              <el-option label="千问" value="qianwen" />
              <el-option label="元宝" value="yuanbao" />
            </el-select>
          </div>
        </div>
      </template>

      <div class="workspace">
        <!-- 左侧：代码输入区域 -->
        <div class="workspace-panel">
          <div class="panel-header">
            <span>代码输入</span>
            <el-button type="text" size="small" @click="clearCode" :icon="Delete">
              清空
            </el-button>
          </div>
          <el-input
            v-model="codeInput"
            type="textarea"
            :rows="20"
            placeholder="请输入您的代码或需求..."
            resize="vertical"
            class="code-input"
          />
        </div>

        <!-- 右侧：AI响应区域 -->
        <div class="workspace-panel">
          <div class="panel-header">
            <span>AI响应</span>
            <el-button type="text" size="small" @click="copyResponse" :icon="DocumentCopy">
              复制
            </el-button>
          </div>
          <div class="ai-response" ref="aiResponse">
            <div v-if="sending" class="typing-indicator">
              <div class="typing-dots">
                <span></span>
                <span></span>
                <span></span>
              </div>
              <span class="typing-text">AI正在思考...</span>
            </div>
            <div v-else-if="aiResponseText" class="response-content">
              <pre><code>{{ aiResponseText }}</code></pre>
            </div>
            <div v-else class="empty-response">
              <el-icon size="48"><ChatLineRound /></el-icon>
              <p>AI响应将显示在这里</p>
              <p class="hint-text">选择一个操作并输入代码或需求</p>
            </div>
          </div>
        </div>
      </div>

      <!-- 操作按钮区域 -->
      <div class="action-buttons">
        <div class="prompt-input">
          <el-input
            v-model="prompt"
            placeholder="输入额外提示（可选）"
            size="large"
            :prefix-icon="Message"
            @keyup.enter.ctrl="executeAction"
          />
        </div>
        <div class="button-group">
          <el-button type="primary" size="large" @click="executeAction" :loading="sending">
            执行
          </el-button>
          <el-button size="large" @click="saveCode" :icon="Download">
            保存代码
          </el-button>
          <el-button size="large" @click="shareCode" :icon="Share">
            分享
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- 历史记录 -->
    <el-card shadow="hover" class="history-card">
      <template #header>
        <div class="card-header">
          <el-icon class="header-icon"><History /></el-icon>
          <span>历史记录</span>
          <el-button type="text" size="small" @click="clearHistory" :icon="Delete">
            清空历史
          </el-button>
        </div>
      </template>
      <div class="history-list">
        <div v-if="history.length === 0" class="empty-history">
          <p>暂无历史记录</p>
        </div>
        <div v-for="(item, index) in history" :key="index" class="history-item" @click="loadHistory(item)">
          <div class="history-item-header">
            <span class="history-action">{{ item.action === 'generate' ? '代码生成' : item.action === 'explain' ? '代码解释' : item.action === 'optimize' ? '代码优化' : '代码调试' }}</span>
            <span class="history-time">{{ item.timestamp }}</span>
          </div>
          <div class="history-content">
            <p class="history-preview">{{ item.codeInput.substring(0, 100) }}...</p>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage, ElNotification } from 'element-plus'
import {
  Refresh, DocumentAdd, Document, Top, Warning,
  Setting, Delete, DocumentCopy, Message,
  Download, Share, ChatLineRound
} from '@element-plus/icons-vue'
import { aiAPI } from '@/api/api'

// 代码输入
const codeInput = ref('')
// AI响应
const aiResponseText = ref('')
// 发送状态
const sending = ref(false)
// 选择的编程语言
const selectedLanguage = ref('javascript')
// 选择的AI模型
const selectedModel = ref('doubao')
// 额外提示
const prompt = ref('')
// 当前操作
const currentAction = ref('generate')
// 历史记录
interface HistoryItem {
  codeInput: string
  aiResponse: string
  action: string
  language: string
  model: string
  prompt: string
  timestamp: string
}

const history = ref<HistoryItem[]>([])
// AI响应区域引用
const aiResponse = ref<HTMLElement | null>(null)

// 设置当前操作
const setAction = (action: string) => {
  currentAction.value = action
  ElMessage.info(`已选择：${action === 'generate' ? '代码生成' : action === 'explain' ? '代码解释' : action === 'optimize' ? '代码优化' : '代码调试'}`)
}

// 执行操作
const executeAction = async () => {
  if (!codeInput.value.trim()) {
    ElMessage.warning('请输入代码或需求')
    return
  }

  sending.value = true
  try {
    // 准备请求参数
    const requestData = {
      model: selectedModel.value,
      code: codeInput.value,
      action: currentAction.value,
      language: selectedLanguage.value,
      prompt: prompt.value
    }

    // 调用AI API
    const response = await aiAPI.chat(requestData)
    aiResponseText.value = response.data.result

    // 保存到历史记录
    saveToHistory()

    // 滚动到底部
    scrollToBottom()
  } catch (error: any) {
    console.error('执行操作失败:', error)
    ElMessage.error(error.response?.data?.error || '执行操作失败')
  } finally {
    sending.value = false
  }
}

// 保存到历史记录
const saveToHistory = () => {
  const historyItem = {
    codeInput: codeInput.value,
    aiResponse: aiResponseText.value,
    action: currentAction.value,
    language: selectedLanguage.value,
    model: selectedModel.value,
    prompt: prompt.value,
    timestamp: new Date().toLocaleString()
  }
  history.value.unshift(historyItem)
  // 只保留最近20条记录
  if (history.value.length > 20) {
    history.value.pop()
  }
}

// 加载历史记录
const loadHistory = (item: any) => {
  codeInput.value = item.codeInput
  aiResponseText.value = item.aiResponse
  currentAction.value = item.action
  selectedLanguage.value = item.language
  selectedModel.value = item.model
  prompt.value = item.prompt
  scrollToBottom()
}

// 清空历史记录
const clearHistory = () => {
  history.value = []
  ElMessage.success('历史记录已清空')
}

// 清空代码输入
const clearCode = () => {
  codeInput.value = ''
}

// 复制响应
const copyResponse = () => {
  if (aiResponseText.value) {
    navigator.clipboard.writeText(aiResponseText.value)
      .then(() => {
        ElMessage.success('已复制到剪贴板')
      })
      .catch(() => {
        ElMessage.error('复制失败')
      })
  }
}

// 重置编辑器
const resetEditor = () => {
  codeInput.value = ''
  aiResponseText.value = ''
  prompt.value = ''
  selectedLanguage.value = 'javascript'
  selectedModel.value = 'doubao'
  currentAction.value = 'generate'
  ElMessage.success('编辑器已重置')
}

// 保存代码
const saveCode = () => {
  if (!aiResponseText.value) {
    ElMessage.warning('没有可保存的代码')
    return
  }
  
  const blob = new Blob([aiResponseText.value], { type: 'text/plain;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `ai-code-${Date.now()}.${getFileExtension(selectedLanguage.value)}`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
  
  ElMessage.success('代码已保存')
}

// 获取文件扩展名
const getFileExtension = (language: string) => {
  const extensions: Record<string, string> = {
    javascript: 'js',
    typescript: 'ts',
    python: 'py',
    go: 'go',
    java: 'java',
    cpp: 'cpp',
    ruby: 'rb',
    php: 'php'
  }
  return extensions[language] || 'txt'
}

// 分享代码
const shareCode = () => {
  if (!aiResponseText.value) {
    ElMessage.warning('没有可分享的代码')
    return
  }
  
  navigator.clipboard.writeText(aiResponseText.value)
    .then(() => {
      ElNotification({
        title: '分享成功',
        message: '代码已复制到剪贴板，您可以粘贴分享给他人',
        type: 'success',
        duration: 3000
      })
    })
    .catch(() => {
      ElMessage.error('分享失败')
    })
}

// 滚动到底部
const scrollToBottom = () => {
  if (aiResponse.value) {
    aiResponse.value.scrollTop = aiResponse.value.scrollHeight
  }
}
</script>

<style scoped>
.ai-code-assistant-view {
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
}

.header-content {
  display: flex;
  align-items: center;
  gap: 20px;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
  justify-content: space-between;
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

/* 快速操作卡片样式 */
.quick-actions {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 20px;
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
  cursor: pointer;
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

.quick-action-content {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  position: relative;
  z-index: 1;
}

.action-icon {
  font-size: 28px;
  color: #409eff;
  background: #ecf5ff;
  padding: 12px;
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
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.action-text p {
  margin: 0;
  font-size: 12px;
  color: #606266;
}

/* 主卡片样式 */
.main-card {
  max-width: 1200px;
  margin: 0 auto 30px auto;
  border-radius: 12px;
  background: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  border: 1px solid #e4e7ed;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
}

.main-card:hover {
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

.main-card:hover .card-header::after {
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

/* 卡片头部右侧区域 */
.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-left: auto;
}

.language-select,
.model-select {
  min-width: 160px;
}

/* 工作区样式 */
.workspace {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
  padding: 24px;
  background-color: #fafafa;
}

.workspace-panel {
  display: flex;
  flex-direction: column;
  background: white;
  border-radius: 8px;
  border: 1px solid #e4e7ed;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #f5f7fa;
  border-bottom: 1px solid #e4e7ed;
  font-weight: 600;
  color: #303133;
}

.code-input {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
  line-height: 1.6;
  padding: 16px;
  resize: vertical;
  min-height: 400px;
  border: none;
  outline: none;
}

.code-input:focus {
  box-shadow: none;
}

/* AI响应区域 */
.ai-response {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
  line-height: 1.6;
  padding: 16px;
  height: 400px;
  overflow-y: auto;
  background: #fafafa;
  border: none;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.response-content {
  background: white;
  padding: 16px;
  border-radius: 8px;
  border: 1px solid #e4e7ed;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.response-content pre {
  margin: 0;
  overflow-x: auto;
}

.response-content code {
  background: transparent;
  padding: 0;
}

.empty-response {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #909399;
  text-align: center;
}

.hint-text {
  font-size: 12px;
  opacity: 0.7;
  margin-top: 8px;
}

/* 操作按钮区域 */
.action-buttons {
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 24px;
  background: white;
  border-top: 1px solid #f0f2f5;
}

.prompt-input {
  width: 100%;
}

.button-group {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

/* 历史记录卡片 */
.history-card {
  max-width: 1200px;
  margin: 0 auto 30px auto;
  border-radius: 12px;
  background: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  border: 1px solid #e4e7ed;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
}

.history-card:hover {
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  border-color: #409eff;
}

/* 历史记录列表 */
.history-list {
  padding: 0 24px 24px 24px;
  max-height: 300px;
  overflow-y: auto;
}

.history-item {
  padding: 16px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  margin-bottom: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  background: #fafafa;
}

.history-item:hover {
  background: #ecf5ff;
  border-color: #409eff;
  transform: translateX(4px);
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.15);
}

.history-item-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.history-action {
  font-weight: 600;
  color: #303133;
}

.history-time {
  font-size: 12px;
  color: #909399;
}

.history-content {
  max-height: 60px;
  overflow: hidden;
}

.history-preview {
  margin: 0;
  font-size: 13px;
  color: #606266;
  line-height: 1.4;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.empty-history {
  text-align: center;
  padding: 40px 0;
  color: #909399;
}

/* 打字指示器 */
.typing-indicator {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 24px;
  animation: fadeIn 0.3s ease;
  background: white;
  padding: 16px;
  border-radius: 8px;
  border: 1px solid #e4e7ed;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
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

/* 响应式设计 */
@media (max-width: 1024px) {
  .workspace {
    grid-template-columns: 1fr;
  }
  
  .code-input,
  .ai-response {
    min-height: 300px;
  }
  
  .button-group {
    flex-direction: column;
  }
  
  .button-group .el-button {
    width: 100%;
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

  .main-card,
  .history-card {
    margin: 0 16px 24px 16px;
  }

  .workspace {
    padding: 16px;
    gap: 16px;
  }

  .action-buttons {
    padding: 16px;
  }
  
  .language-select,
  .model-select {
    min-width: 120px;
  }
  
  .header-right {
    flex-wrap: wrap;
    justify-content: center;
  }
}
</style>