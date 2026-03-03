<template>
  <div class="ai-content-generator-view">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <el-icon class="header-icon"><EditPen /></el-icon>
        <div class="header-text">
          <h2>AI内容生成</h2>
          <p>智能生成文章、标题、摘要、广告语等多种内容</p>
        </div>
        <div class="header-actions">
          <el-button type="primary" size="small" :icon="Refresh" @click="resetGenerator">
            重置
          </el-button>
        </div>
      </div>
    </div>

    <!-- 快速操作卡片 -->
    <div class="quick-actions">
      <el-card shadow="hover" class="quick-action-card" @click="setContentType('article')">
        <div class="quick-action-content">
          <el-icon class="action-icon"><Document /></el-icon>
          <div class="action-text">
            <h3>文章生成</h3>
            <p>生成完整的文章内容</p>
          </div>
        </div>
      </el-card>
      
      <el-card shadow="hover" class="quick-action-card" @click="setContentType('title')">
        <div class="quick-action-content">
          <el-icon class="action-icon"><Position /></el-icon>
          <div class="action-text">
            <h3>标题生成</h3>
            <p>生成吸引人的标题</p>
          </div>
        </div>
      </el-card>
      
      <el-card shadow="hover" class="quick-action-card" @click="setContentType('summary')">
        <div class="quick-action-content">
          <el-icon class="action-icon"><DocumentCopy /></el-icon>
          <div class="action-text">
            <h3>摘要生成</h3>
            <p>生成内容摘要</p>
          </div>
        </div>
      </el-card>
      
      <el-card shadow="hover" class="quick-action-card" @click="setContentType('ad')">
        <div class="quick-action-content">
          <el-icon class="action-icon"><Promotion /></el-icon>
          <div class="action-text">
            <h3>广告语生成</h3>
            <p>生成吸引人的广告语</p>
          </div>
        </div>
      </el-card>
      
      <el-card shadow="hover" class="quick-action-card" @click="setContentType('email')">
        <div class="quick-action-content">
          <el-icon class="action-icon"><Message /></el-icon>
          <div class="action-text">
            <h3>邮件生成</h3>
            <p>生成专业邮件内容</p>
          </div>
        </div>
      </el-card>
      
      <el-card shadow="hover" class="quick-action-card" @click="setContentType('social')">
        <div class="quick-action-content">
          <el-icon class="action-icon"><Share /></el-icon>
          <div class="action-text">
            <h3>社交媒体</h3>
            <p>生成社交媒体内容</p>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 主内容区域 -->
    <el-card shadow="hover" class="main-card">
      <template #header>
        <div class="card-header">
          <el-icon class="header-icon"><Setting /></el-icon>
          <span>内容生成工作台</span>
          <div class="header-right">
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
        <!-- 左侧：内容输入和参数设置 -->
        <div class="workspace-panel">
          <div class="panel-header">
            <span>内容设置</span>
            <el-button type="text" size="small" @click="clearContent" :icon="Delete">
              清空
            </el-button>
          </div>
          
          <!-- 内容类型选择 -->
          <div class="setting-section">
            <label class="setting-label">内容类型</label>
            <el-select v-model="contentType" placeholder="选择内容类型" class="setting-select">
              <el-option label="文章" value="article" />
              <el-option label="标题" value="title" />
              <el-option label="摘要" value="summary" />
              <el-option label="广告语" value="ad" />
              <el-option label="邮件" value="email" />
              <el-option label="社交媒体" value="social" />
            </el-select>
          </div>
          
          <!-- 主题输入 -->
          <div class="setting-section">
            <label class="setting-label">主题/关键词</label>
            <el-input
              v-model="contentTopic"
              placeholder="请输入内容主题或关键词"
              size="large"
              :prefix-icon="Search"
            />
          </div>
          
          <!-- 详细要求 -->
          <div class="setting-section">
            <label class="setting-label">详细要求</label>
            <el-input
              v-model="contentRequirements"
              type="textarea"
              :rows="6"
              placeholder="请输入更详细的要求..."
              resize="vertical"
            />
          </div>
          
          <!-- 生成参数 -->
          <div class="setting-section">
            <label class="setting-label">生成参数</label>
            <div class="param-grid">
              <div class="param-item">
                <label>长度</label>
                <el-select v-model="contentLength" placeholder="选择长度">
                  <el-option label="短" value="short" />
                  <el-option label="中" value="medium" />
                  <el-option label="长" value="long" />
                </el-select>
              </div>
              <div class="param-item">
                <label>风格</label>
                <el-select v-model="contentStyle" placeholder="选择风格">
                  <el-option label="正式" value="formal" />
                  <el-option label="随意" value="casual" />
                  <el-option label="专业" value="professional" />
                  <el-option label="创意" value="creative" />
                </el-select>
              </div>
              <div class="param-item">
                <label>语言</label>
                <el-select v-model="contentLanguage" placeholder="选择语言">
                  <el-option label="中文" value="chinese" />
                  <el-option label="英文" value="english" />
                </el-select>
              </div>
            </div>
          </div>
        </div>

        <!-- 右侧：AI响应区域 -->
        <div class="workspace-panel">
          <div class="panel-header">
            <span>AI生成结果</span>
            <div class="panel-actions">
              <el-button type="text" size="small" @click="copyContent" :icon="DocumentCopy">
                复制
              </el-button>
              <el-button type="text" size="small" @click="saveContent" :icon="Download">
                保存
              </el-button>
              <el-button type="text" size="small" @click="shareContent" :icon="Share">
                分享
              </el-button>
            </div>
          </div>
          <div class="ai-response" ref="aiResponse">
            <div v-if="generating" class="typing-indicator">
              <div class="typing-dots">
                <span></span>
                <span></span>
                <span></span>
              </div>
              <span class="typing-text">AI正在生成内容...</span>
            </div>
            <div v-else-if="generatedContent" class="response-content">
              <h3 class="content-title">生成结果</h3>
              <div class="content-body">{{ generatedContent }}</div>
            </div>
            <div v-else class="empty-response">
              <el-icon size="48"><EditPen /></el-icon>
              <p>AI生成结果将显示在这里</p>
              <p class="hint-text">设置内容参数并点击生成按钮</p>
            </div>
          </div>
        </div>
      </div>

      <!-- 操作按钮区域 -->
      <div class="action-buttons">
        <div class="button-group">
          <el-button type="primary" size="large" @click="generateContent" :loading="generating">
            生成内容
          </el-button>
          <el-button size="large" @click="regenerateContent" :loading="generating" :icon="Refresh">
            重新生成
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- 历史记录 -->
    <el-card shadow="hover" class="history-card">
      <template #header>
        <div class="card-header">
          <el-icon class="header-icon"><History /></el-icon>
          <span>生成历史</span>
          <el-button type="text" size="small" @click="clearHistory" :icon="Delete">
            清空历史
          </el-button>
        </div>
      </template>
      <div class="history-list">
        <div v-if="history.length === 0" class="empty-history">
          <p>暂无生成历史</p>
        </div>
        <div v-for="(item, index) in history" :key="index" class="history-item" @click="loadHistory(item)">
          <div class="history-item-header">
            <span class="history-type">{{ getContentTypeText(item.contentType) }}</span>
            <span class="history-time">{{ item.timestamp }}</span>
          </div>
          <div class="history-content">
            <p class="history-preview">{{ item.contentTopic }}</p>
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
  EditPen, Refresh, Document, Position, DocumentCopy, 
  Promotion, Message, Share, Setting, Delete, Search, 
  Download, 
} from '@element-plus/icons-vue'
import { aiAPI } from '@/api/api'

// 内容类型
const contentType = ref('article')
// 内容主题
const contentTopic = ref('')
// 内容要求
const contentRequirements = ref('')
// 内容长度
const contentLength = ref('medium')
// 内容风格
const contentStyle = ref('professional')
// 内容语言
const contentLanguage = ref('chinese')
// 生成的内容
const generatedContent = ref('')
// 生成状态
const generating = ref(false)
// 选择的AI模型
const selectedModel = ref('doubao')
// 历史记录
interface HistoryItem {
  contentType: string
  contentTopic: string
  contentRequirements: string
  contentLength: string
  contentStyle: string
  contentLanguage: string
  generatedContent: string
  model: string
  timestamp: string
}

const history = ref<HistoryItem[]>([])
// AI响应区域引用
const aiResponse = ref<HTMLElement | null>(null)

// 设置内容类型
const setContentType = (type: string) => {
  contentType.value = type
  ElMessage.info(`已选择：${getContentTypeText(type)}生成`)
}

// 获取内容类型文本
const getContentTypeText = (type: string) => {
  const typeMap: Record<string, string> = {
    article: '文章',
    title: '标题',
    summary: '摘要',
    ad: '广告语',
    email: '邮件',
    social: '社交媒体'
  }
  return typeMap[type] || type
}

// 生成内容
const generateContent = async () => {
  if (!contentTopic.value.trim()) {
    ElMessage.warning('请输入内容主题或关键词')
    return
  }

  generating.value = true
  try {
    // 准备请求参数
    const requestData = {
      model: selectedModel.value,
      contentType: contentType.value,
      topic: contentTopic.value,
      requirements: contentRequirements.value,
      length: contentLength.value,
      style: contentStyle.value,
      language: contentLanguage.value
    }

    // 调用AI API
    const response = await aiAPI.chat(requestData)
    generatedContent.value = response.data.content

    // 保存到历史记录
    saveToHistory()

    // 滚动到底部
    scrollToBottom()
  } catch (error: any) {
    console.error('生成内容失败:', error)
    ElMessage.error(error.response?.data?.error || '生成内容失败')
  } finally {
    generating.value = false
  }
}

// 重新生成内容
const regenerateContent = () => {
  if (generatedContent.value) {
    generateContent()
  }
}

// 保存到历史记录
const saveToHistory = () => {
  const historyItem = {
    contentType: contentType.value,
    contentTopic: contentTopic.value,
    contentRequirements: contentRequirements.value,
    contentLength: contentLength.value,
    contentStyle: contentStyle.value,
    contentLanguage: contentLanguage.value,
    generatedContent: generatedContent.value,
    model: selectedModel.value,
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
  contentType.value = item.contentType
  contentTopic.value = item.contentTopic
  contentRequirements.value = item.contentRequirements
  contentLength.value = item.contentLength
  contentStyle.value = item.contentStyle
  contentLanguage.value = item.contentLanguage
  generatedContent.value = item.generatedContent
  selectedModel.value = item.model
  scrollToBottom()
}

// 清空历史记录
const clearHistory = () => {
  history.value = []
  ElMessage.success('历史记录已清空')
}

// 清空内容设置
const clearContent = () => {
  contentTopic.value = ''
  contentRequirements.value = ''
  generatedContent.value = ''
}

// 复制生成的内容
const copyContent = () => {
  if (generatedContent.value) {
    navigator.clipboard.writeText(generatedContent.value)
      .then(() => {
        ElMessage.success('已复制到剪贴板')
      })
      .catch(() => {
        ElMessage.error('复制失败')
      })
  }
}

// 保存生成的内容
const saveContent = () => {
  if (!generatedContent.value) {
    ElMessage.warning('没有可保存的内容')
    return
  }
  
  const content = `# AI生成内容\n\n## 内容类型\n${getContentTypeText(contentType.value)}\n\n## 主题\n${contentTopic.value}\n\n## 生成模型\n${selectedModel.value}\n\n## 生成时间\n${new Date().toLocaleString()}\n\n## 生成结果\n${generatedContent.value}`
  
  const blob = new Blob([content], { type: 'text/plain;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `ai-content-${contentType.value}-${Date.now()}.txt`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
  
  ElMessage.success('内容已保存')
}

// 分享生成的内容
const shareContent = () => {
  if (!generatedContent.value) {
    ElMessage.warning('没有可分享的内容')
    return
  }
  
  navigator.clipboard.writeText(generatedContent.value)
    .then(() => {
      ElNotification({
        title: '分享成功',
        message: '内容已复制到剪贴板，您可以粘贴分享给他人',
        type: 'success',
        duration: 3000
      })
    })
    .catch(() => {
      ElMessage.error('分享失败')
    })
}

// 重置生成器
const resetGenerator = () => {
  contentTopic.value = ''
  contentRequirements.value = ''
  generatedContent.value = ''
  contentType.value = 'article'
  contentLength.value = 'medium'
  contentStyle.value = 'professional'
  contentLanguage.value = 'chinese'
  selectedModel.value = 'doubao'
  ElMessage.success('生成器已重置')
}

// 滚动到底部
const scrollToBottom = () => {
  if (aiResponse.value) {
    aiResponse.value.scrollTop = aiResponse.value.scrollHeight
  }
}
</script>

<style scoped>
.ai-content-generator-view {
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
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
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

.panel-actions {
  display: flex;
  gap: 8px;
}

/* 设置区域样式 */
.setting-section {
  padding: 16px;
  border-bottom: 1px solid #f0f2f5;
}

.setting-section:last-child {
  border-bottom: none;
}

.setting-label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  color: #303133;
}

.setting-select {
  width: 100%;
}

/* 参数网格 */
.param-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
  gap: 16px;
}

.param-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.param-item label {
  font-size: 13px;
  color: #606266;
  font-weight: 500;
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
  padding: 20px;
  border-radius: 8px;
  border: 1px solid #e4e7ed;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.content-title {
  margin: 0 0 16px 0;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  border-bottom: 2px solid #ecf5ff;
  padding-bottom: 8px;
}

.content-body {
  color: #606266;
  line-height: 1.8;
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
  justify-content: center;
  gap: 16px;
  padding: 24px;
  background: white;
  border-top: 1px solid #f0f2f5;
}

.button-group {
  display: flex;
  gap: 12px;
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

.history-type {
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
  
  .ai-response {
    min-height: 300px;
  }
  
  .button-group {
    flex-direction: column;
  }
  
  .button-group .el-button {
    width: 100%;
  }
  
  .quick-actions {
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
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
    grid-template-columns: 1fr 1fr;
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
  
  .param-grid {
    grid-template-columns: 1fr;
  }
  
  .model-select {
    min-width: 120px;
  }
  
  .header-right {
    flex-wrap: wrap;
    justify-content: center;
  }
}
</style>