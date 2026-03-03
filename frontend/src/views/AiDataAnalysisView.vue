<template>
  <div class="ai-data-analysis-view">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <el-icon class="header-icon"><DataAnalysis /></el-icon>
        <div class="header-text">
          <h2>AI数据分析</h2>
          <p>智能数据分析、可视化和洞察工具</p>
        </div>
        <div class="header-actions">
          <el-button type="primary" size="small" :icon="Refresh" @click="resetAnalysis">
            重置
          </el-button>
        </div>
      </div>
    </div>

    <!-- 快速操作卡片 -->
    <div class="quick-actions">
      <el-card shadow="hover" class="quick-action-card" @click="setAnalysisType('import')">
        <div class="quick-action-content">
          <el-icon class="action-icon"><Upload /></el-icon>
          <div class="action-text">
            <h3>数据导入</h3>
            <p>上传或输入数据进行分析</p>
          </div>
        </div>
      </el-card>
      
      <el-card shadow="hover" class="quick-action-card" @click="setAnalysisType('clean')">
        <div class="quick-action-content">
          <el-icon class="action-icon"><Cleanup /></el-icon>
          <div class="action-text">
            <h3>数据清洗</h3>
            <p>智能清理和预处理数据</p>
          </div>
        </div>
      </el-card>
      
      <el-card shadow="hover" class="quick-action-card" @click="setAnalysisType('visualize')">
        <div class="quick-action-content">
          <el-icon class="action-icon"><TrendCharts /></el-icon>
          <div class="action-text">
            <h3>数据可视化</h3>
            <p>生成交互式数据图表</p>
          </div>
        </div>
      </el-card>
      
      <el-card shadow="hover" class="quick-action-card" @click="setAnalysisType('insights')">
        <div class="quick-action-content">
          <el-icon class="action-icon"><Lightbulb /></el-icon>
          <div class="action-text">
            <h3>智能洞察</h3>
            <p>获取AI生成的数据洞察</p>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 主内容区域 -->
    <el-card shadow="hover" class="main-card">
      <template #header>
        <div class="card-header">
          <el-icon class="header-icon"><Setting /></el-icon>
          <span>数据分析工作台</span>
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
        <!-- 左侧：数据输入区域 -->
        <div class="workspace-panel">
          <div class="panel-header">
            <span>数据输入</span>
            <div class="panel-actions">
              <el-button type="text" size="small" @click="clearData" :icon="Delete">
                清空
              </el-button>
              <el-button type="text" size="small" @click="showDataUpload" :icon="Upload">
                上传
              </el-button>
            </div>
          </div>
          
          <!-- 数据输入区域 -->
          <div class="data-input-section">
            <el-tabs v-model="dataInputType" class="data-input-tabs">
              <el-tab-pane label="文本输入" name="text">
                <el-input
                  v-model="dataInput"
                  type="textarea"
                  :rows="15"
                  placeholder="请输入您的数据...\n例如：\n姓名,年龄,性别,城市\n张三,25,男,北京\n李四,30,女,上海\n王五,28,男,广州"
                  resize="vertical"
                  class="data-text-input"
                />
              </el-tab-pane>
              <el-tab-pane label="CSV上传" name="csv">
                <div class="file-upload-area">
                  <el-upload
                    class="upload-demo"
                    drag
                    action="#"
                    multiple
                    :auto-upload="false"
                    :on-change="handleFileChange"
                  >
                    <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
                    <div class="el-upload__text">
                      将文件拖到此处，或<em>点击上传</em>
                    </div>
                    <template #tip>
                      <div class="el-upload__tip">
                        支持上传 .csv, .txt, .json 格式文件
                      </div>
                    </template>
                  </el-upload>
                  <div v-if="uploadedFile" class="uploaded-file-info">
                    <el-icon><Document /></el-icon>
                    <span>{{ uploadedFile.name }}</span>
                    <el-button type="text" size="small" @click="removeFile" :icon="Delete">
                      移除
                    </el-button>
                  </div>
                </div>
              </el-tab-pane>
            </el-tabs>
          </div>
        </div>

        <!-- 右侧：AI分析结果区域 -->
        <div class="workspace-panel">
          <div class="panel-header">
            <span>AI分析结果</span>
            <el-button type="text" size="small" @click="copyAnalysis" :icon="DocumentCopy">
              复制
            </el-button>
          </div>
          <div class="ai-response" ref="aiResponse">
            <div v-if="analyzing" class="typing-indicator">
              <div class="typing-dots">
                <span></span>
                <span></span>
                <span></span>
              </div>
              <span class="typing-text">AI正在分析数据...</span>
            </div>
            <div v-else-if="analysisResult" class="response-content">
              <h3 class="analysis-title">分析结果</h3>
              <div class="analysis-text">{{ analysisResult }}</div>
            </div>
            <div v-else class="empty-response">
              <el-icon size="48"><DataAnalysis /></el-icon>
              <p>AI分析结果将显示在这里</p>
              <p class="hint-text">输入数据并点击分析按钮</p>
            </div>
          </div>
        </div>
      </div>

      <!-- 数据可视化区域 -->
      <div class="visualization-section" v-if="visualizationData">
        <div class="panel-header">
          <span>数据可视化</span>
          <el-select v-model="chartType" placeholder="选择图表类型" size="small" @change="updateVisualization">
            <el-option label="柱状图" value="bar" />
            <el-option label="折线图" value="line" />
            <el-option label="饼图" value="pie" />
            <el-option label="散点图" value="scatter" />
            <el-option label="雷达图" value="radar" />
          </el-select>
        </div>
        <div class="chart-container">
          <div class="chart-placeholder">
            <el-icon size="64"><TrendCharts /></el-icon>
            <p>图表可视化区域</p>
            <p class="hint-text">AI将根据分析结果生成交互式图表</p>
          </div>
        </div>
      </div>

      <!-- 操作按钮区域 -->
      <div class="action-buttons">
        <div class="prompt-input">
          <el-input
            v-model="analysisPrompt"
            placeholder="输入分析需求（可选）"
            size="large"
            :prefix-icon="Message"
            @keyup.enter.ctrl="analyzeData"
          />
        </div>
        <div class="button-group">
          <el-button type="primary" size="large" @click="analyzeData" :loading="analyzing" :icon="Search">
            分析数据
          </el-button>
          <el-button size="large" @click="saveAnalysis" :icon="Download">
            保存结果
          </el-button>
          <el-button size="large" @click="shareAnalysis" :icon="Share">
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
          <span>分析历史</span>
          <el-button type="text" size="small" @click="clearHistory" :icon="Delete">
            清空历史
          </el-button>
        </div>
      </template>
      <div class="history-list">
        <div v-if="history.length === 0" class="empty-history">
          <p>暂无分析历史</p>
        </div>
        <div v-for="(item, index) in history" :key="index" class="history-item" @click="loadHistory(item)">
          <div class="history-item-header">
            <span class="history-type">{{ getAnalysisTypeText(item.analysisType) }}</span>
            <span class="history-time">{{ item.timestamp }}</span>
          </div>
          <div class="history-content">
            <p class="history-preview">{{ item.dataInput.substring(0, 100) }}...</p>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 数据上传对话框 -->
    <el-dialog v-model="showUploadDialog" title="数据上传" width="500px">
      <el-upload
        class="upload-demo"
        drag
        action="#"
        multiple
        :auto-upload="false"
        :on-change="handleFileChange"
      >
        <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
        <div class="el-upload__text">
          将文件拖到此处，或<em>点击上传</em>
        </div>
        <template #tip>
          <div class="el-upload__tip">
            支持上传 .csv, .txt, .json 格式文件
          </div>
        </template>
      </el-upload>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showUploadDialog = false">取消</el-button>
          <el-button type="primary" @click="showUploadDialog = false">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage, ElNotification } from 'element-plus'
import { 
  DataAnalysis, Refresh, Upload, TrendCharts, 
  Setting, Delete, DocumentCopy, Message, 
  Search, Download, Share, Document, 
  UploadFilled 
} from '@element-plus/icons-vue'
import { aiAPI } from '@/api/api'

// 数据分析类型
const analysisType = ref('import')
// 数据输入类型
const dataInputType = ref('text')
// 数据输入
const dataInput = ref('')
// AI分析结果
const analysisResult = ref('')
// 可视化数据
const visualizationData = ref(null)
// 图表类型
const chartType = ref('bar')
// 分析提示
const analysisPrompt = ref('')
// 分析状态
const analyzing = ref(false)
// 选择的AI模型
const selectedModel = ref('doubao')
// 上传的文件
interface UploadedFile {
  name: string
  size: number
  type: string
}
const uploadedFile = ref<UploadedFile | null>(null)
// 显示上传对话框
const showUploadDialog = ref(false)
// 历史记录
interface HistoryItem {
  dataInput: string
  analysisResult: string
  analysisType: string
  chartType: string
  model: string
  prompt: string
  timestamp: string
}
const history = ref<HistoryItem[]>([])
// AI响应区域引用
const aiResponse = ref<HTMLElement | null>(null)

// 设置分析类型
const setAnalysisType = (type: string) => {
  analysisType.value = type
  ElMessage.info(`已选择：${getAnalysisTypeText(type)}`)
}

// 获取分析类型文本
const getAnalysisTypeText = (type: string) => {
  const typeMap: Record<string, string> = {
    import: '数据导入',
    clean: '数据清洗',
    visualize: '数据可视化',
    insights: '智能洞察'
  }
  return typeMap[type] || type
}

// 显示数据上传
const showDataUpload = () => {
  showUploadDialog.value = true
}

// 处理文件上传
const handleFileChange = (file: any) => {
  uploadedFile.value = file
  // 读取文件内容
  const reader = new FileReader()
  reader.onload = (e) => {
    dataInput.value = e.target?.result as string
    dataInputType.value = 'text'
  }
  reader.readAsText(file.raw)
}

// 移除上传的文件
const removeFile = () => {
  uploadedFile.value = null
}

// 执行数据分析
const analyzeData = async () => {
  if (!dataInput.value.trim()) {
    ElMessage.warning('请输入或上传数据')
    return
  }

  analyzing.value = true
  try {
    // 准备请求参数
    const requestData = {
      model: selectedModel.value,
      data: dataInput.value,
      analysisType: analysisType.value,
      chartType: chartType.value,
      prompt: analysisPrompt.value
    }

    // 调用AI API
    const response = await aiAPI.chat(requestData)
    analysisResult.value = response.data.result
    visualizationData.value = response.data.visualization || null

    // 保存到历史记录
    saveToHistory()

    // 滚动到底部
    scrollToBottom()
  } catch (error: any) {
    console.error('数据分析失败:', error)
    ElMessage.error(error.response?.data?.error || '数据分析失败')
  } finally {
    analyzing.value = false
  }
}

// 更新可视化
const updateVisualization = () => {
  // 这里可以实现图表类型切换逻辑
  ElMessage.info(`图表类型已切换为：${chartType.value}`)
}

// 保存到历史记录
const saveToHistory = () => {
  const historyItem = {
    dataInput: dataInput.value,
    analysisResult: analysisResult.value,
    analysisType: analysisType.value,
    chartType: chartType.value,
    model: selectedModel.value,
    prompt: analysisPrompt.value,
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
  dataInput.value = item.dataInput
  analysisResult.value = item.analysisResult
  analysisType.value = item.analysisType
  chartType.value = item.chartType
  selectedModel.value = item.model
  analysisPrompt.value = item.prompt
  visualizationData.value = null // 重置可视化数据
  scrollToBottom()
}

// 清空历史记录
const clearHistory = () => {
  history.value = []
  ElMessage.success('历史记录已清空')
}

// 清空数据输入
const clearData = () => {
  dataInput.value = ''
  uploadedFile.value = null
}

// 复制分析结果
const copyAnalysis = () => {
  if (analysisResult.value) {
    navigator.clipboard.writeText(analysisResult.value)
      .then(() => {
        ElMessage.success('已复制到剪贴板')
      })
      .catch(() => {
        ElMessage.error('复制失败')
      })
  }
}

// 重置分析
const resetAnalysis = () => {
  dataInput.value = ''
  analysisResult.value = ''
  analysisPrompt.value = ''
  analysisType.value = 'import'
  dataInputType.value = 'text'
  selectedModel.value = 'doubao'
  chartType.value = 'bar'
  uploadedFile.value = null
  visualizationData.value = null
  ElMessage.success('分析已重置')
}

// 保存分析结果
const saveAnalysis = () => {
  if (!analysisResult.value) {
    ElMessage.warning('没有可保存的分析结果')
    return
  }
  
  const content = `# AI数据分析结果\n\n## 分析类型\n${getAnalysisTypeText(analysisType.value)}\n\n## 分析模型\n${selectedModel.value}\n\n## 分析时间\n${new Date().toLocaleString()}\n\n## 分析结果\n${analysisResult.value}`
  
  const blob = new Blob([content], { type: 'text/plain;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `ai-analysis-${Date.now()}.txt`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
  
  ElMessage.success('分析结果已保存')
}

// 分享分析结果
const shareAnalysis = () => {
  if (!analysisResult.value) {
    ElMessage.warning('没有可分享的分析结果')
    return
  }
  
  navigator.clipboard.writeText(analysisResult.value)
    .then(() => {
      ElNotification({
        title: '分享成功',
        message: '分析结果已复制到剪贴板，您可以粘贴分享给他人',
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
.ai-data-analysis-view {
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

/* 数据输入区域 */
.data-input-section {
  padding: 16px;
}

.data-input-tabs {
  margin-bottom: 0;
}

.data-text-input {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
  line-height: 1.6;
  min-height: 300px;
  resize: vertical;
}

/* 文件上传区域 */
.file-upload-area {
  padding: 20px;
  text-align: center;
}

.uploaded-file-info {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px;
  background: #f5f7fa;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  margin-top: 16px;
}

/* AI响应区域 */
.ai-response {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
  line-height: 1.6;
  padding: 16px;
  height: 300px;
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

.analysis-title {
  margin: 0 0 16px 0;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  border-bottom: 2px solid #ecf5ff;
  padding-bottom: 8px;
}

.analysis-text {
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

/* 可视化区域 */
.visualization-section {
  margin: 0 24px 24px 24px;
  background: white;
  border-radius: 8px;
  border: 1px solid #e4e7ed;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.chart-container {
  padding: 24px;
  background: #fafafa;
  min-height: 300px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.chart-placeholder {
  text-align: center;
  color: #909399;
}

.chart-placeholder .el-icon {
  margin-bottom: 16px;
  color: #c0c4cc;
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
  
  .visualization-section {
    margin: 0 16px 16px 16px;
  }
  
  .chart-container {
    padding: 16px;
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