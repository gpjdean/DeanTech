<template>
  <div class="image-compress-container">
    <el-card class="main-card">
      <template #header>
        <div class="card-header">
          <el-icon><Picture /></el-icon>
          <h2>图片压缩工具</h2>
        </div>
      </template>

      <!-- 操作区域 -->
      <el-row :gutter="20">
        <!-- 左侧上传和设置 -->
        <el-col :span="12">
          <!-- 图片上传 -->
          <el-card shadow="hover" class="upload-card">
            <template #header>
              <div class="sub-card-header">
                <el-icon><UploadFilled /></el-icon>
                <h3>上传图片</h3>
              </div>
            </template>
            <div class="upload-area">
              <el-upload
                v-model:file-list="fileList"
                class="upload-demo"
                drag
                action=""
                :auto-upload="false"
                :on-change="handleFileChange"
                :show-file-list="false"
                accept="image/*"
              >
                <el-icon class="el-icon--upload"><Plus /></el-icon>
                <div class="el-upload__text">
                  将图片拖到此处，或 <em>点击上传</em>
                </div>
                <template #tip>
                  <div class="el-upload__tip">
                    支持 JPG、PNG、WebP 等格式，单文件不超过 10MB
                  </div>
                </template>
              </el-upload>
            </div>
            <!-- 图片预览 -->
            <div v-if="originalImage" class="image-preview">
              <h4>原图预览</h4>
              <div class="preview-container">
                <img
                  :src="originalImage"
                  :alt="originalFileName"
                  class="preview-image"
                  @load="handleImageLoad"
                />
              </div>
              <div class="image-info">
                <p>文件名: {{ originalFileName }}</p>
                <p>尺寸: {{ originalWidth }} x {{ originalHeight }} px</p>
                <p>大小: {{ formatFileSize(originalSize) }}</p>
              </div>
            </div>
          </el-card>

          <!-- 压缩设置 -->
          <el-card shadow="hover" class="settings-card" v-if="originalImage">
            <template #header>
              <div class="sub-card-header">
                <el-icon><Setting /></el-icon>
                <h3>压缩设置</h3>
              </div>
            </template>
            <div class="settings-content">
              <!-- 质量调整 -->
              <div class="setting-item">
                <div class="setting-label">
                  <span>压缩质量</span>
                  <span class="value">{{ quality }}%</span>
                </div>
                <el-slider
                  v-model="quality"
                  :min="10"
                  :max="100"
                  :step="1"
                  @change="handleCompress"
                />
              </div>

              <!-- 输出格式 -->
              <div class="setting-item">
                <div class="setting-label">
                  <span>输出格式</span>
                </div>
                <el-select
                  v-model="outputFormat"
                  placeholder="请选择输出格式"
                  @change="handleCompress"
                >
                  <el-option label="JPEG" value="image/jpeg" />
                  <el-option label="PNG" value="image/png" />
                  <el-option label="WebP" value="image/webp" />
                </el-select>
              </div>

              <!-- 调整尺寸 -->
              <div class="setting-item">
                <div class="setting-label">
                  <span>调整尺寸</span>
                </div>
                <el-radio-group v-model="resizeMode" @change="handleResizeModeChange">
                  <el-radio-button label="original">保持原图尺寸</el-radio-button>
                  <el-radio-button label="scale">按比例缩放</el-radio-button>
                  <el-radio-button label="custom">自定义尺寸</el-radio-button>
                </el-radio-group>
              </div>

              <!-- 缩放比例 -->
              <div class="setting-item" v-if="resizeMode === 'scale'">
                <div class="setting-label">
                  <span>缩放比例</span>
                  <span class="value">{{ scale }}%</span>
                </div>
                <el-slider
                  v-model="scale"
                  :min="10"
                  :max="100"
                  :step="1"
                  @change="handleCompress"
                />
              </div>

              <!-- 自定义尺寸 -->
              <div class="setting-item custom-size" v-if="resizeMode === 'custom'">
                <div class="setting-label">
                  <span>自定义宽高</span>
                </div>
                <div class="size-inputs">
                  <el-input-number
                    v-model="customWidth"
                    :min="1"
                    :max="originalWidth * 2"
                    @change="handleCompress"
                    placeholder="宽度"
                  />
                  <span class="size-separator">x</span>
                  <el-input-number
                    v-model="customHeight"
                    :min="1"
                    :max="originalHeight * 2"
                    @change="handleCompress"
                    placeholder="高度"
                  />
                  <el-button type="primary" size="small" @click="lockAspectRatio = !lockAspectRatio" class="aspect-ratio-btn">
                    <el-icon><Link /></el-icon>
                    {{ lockAspectRatio ? '锁定比例' : '解锁比例' }}
                  </el-button>
                </div>
              </div>

              <!-- 压缩按钮 -->
              <div class="compress-btn-container">
                <el-button type="primary" size="large" @click="handleCompress" :loading="isCompressing">
                  <el-icon><Upload /></el-icon>
                  开始压缩
                </el-button>
                <el-button type="default" size="large" @click="resetAll">
                  <el-icon><Refresh /></el-icon>
                  重置
                </el-button>
              </div>
            </div>
          </el-card>
        </el-col>

        <!-- 右侧预览和结果 -->
        <el-col :span="12">
          <!-- 压缩结果预览 -->
          <el-card shadow="hover" class="result-card">
            <template #header>
              <div class="sub-card-header">
                <el-icon><Picture /></el-icon>
                <h3>压缩结果</h3>
              </div>
            </template>
            <div v-if="compressedImage" class="result-area">
              <h4>压缩后预览</h4>
              <div class="preview-container">
                <img
                  :src="compressedImage"
                  alt="压缩后的图片"
                  class="preview-image"
                />
              </div>
              <div class="image-info">
                <p>格式: {{ outputFormat.split('/')[1].toUpperCase() }}</p>
                <p>尺寸: {{ compressedWidth }} x {{ compressedHeight }} px</p>
                <p>大小: {{ formatFileSize(compressedSize) }}</p>
                <p v-if="originalSize > 0" class="compression-rate">
                  压缩率: {{ ((1 - compressedSize / originalSize) * 100).toFixed(1) }}%
                </p>
              </div>

              <!-- 下载和复制按钮 -->
              <div class="result-actions">
                <el-button type="primary" @click="downloadImage" :disabled="!compressedImage">
                  <el-icon><Download /></el-icon>
                  下载图片
                </el-button>
                <el-button type="success" @click="copyImage" :disabled="!compressedImage">
                  <el-icon><DocumentCopy /></el-icon>
                  复制链接
                </el-button>
              </div>
            </div>
            <div v-else class="no-result">
              <el-empty description="请上传图片并设置压缩参数" />
            </div>
          </el-card>

          <!-- 压缩历史 -->
          <el-card shadow="hover" class="history-card">
            <template #header>
              <div class="sub-card-header">
                <el-icon><Clock /></el-icon>
                <h3>压缩历史</h3>
                <el-button
                  type="text"
                  size="small"
                  @click="clearHistory"
                  :disabled="history.length === 0"
                >
                  清空
                </el-button>
              </div>
            </template>
            <div v-if="history.length > 0" class="history-list">
              <el-timeline>
                <el-timeline-item
                  v-for="(item, index) in history"
                  :key="index"
                  :timestamp="item.timestamp"
                >
                  <div class="history-item">
                    <div class="history-info">
                      <span class="history-filename">{{ item.filename }}</span>
                      <span class="history-size">{{ formatFileSize(item.originalSize) }} → {{ formatFileSize(item.compressedSize) }}</span>
                      <span class="history-rate">({{ item.compressionRate }}% 压缩率)</span>
                    </div>
                    <div class="history-actions">
                      <el-button
                        type="text"
                        size="small"
                        @click="previewHistory(item)"
                      >
                        <el-icon><View /></el-icon>
                        预览
                      </el-button>
                      <el-button
                        type="text"
                        size="small"
                        @click="downloadHistory(item)"
                      >
                        <el-icon><Download /></el-icon>
                        下载
                      </el-button>
                    </div>
                  </div>
                </el-timeline-item>
              </el-timeline>
            </div>
            <div v-else class="no-history">
              <el-empty description="暂无压缩历史" />
            </div>
          </el-card>
        </el-col>
      </el-row>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Picture,
  UploadFilled,
  Plus,
  Setting,
  Refresh,
  Upload,
  Download,
  DocumentCopy,
  Clock,
  View,
  Link
} from '@element-plus/icons-vue'
import { imageCompressAPI } from '../api/api'

// 文件列表
const fileList = ref<any[]>([])

// 原图信息
const originalImage = ref<string>('')
const originalFileName = ref<string>('')
const originalSize = ref<number>(0)
const originalWidth = ref<number>(0)
const originalHeight = ref<number>(0)

// 压缩设置
const quality = ref<number>(80)
const outputFormat = ref<string>('image/jpeg')
const resizeMode = ref<string>('original')
const scale = ref<number>(100)
const customWidth = ref<number>(0)
const customHeight = ref<number>(0)
const lockAspectRatio = ref<boolean>(true)

// 压缩结果
const compressedImage = ref<string>('')
const compressedSize = ref<number>(0)
const compressedWidth = ref<number>(0)
const compressedHeight = ref<number>(0)
const isCompressing = ref<boolean>(false)

// 压缩历史
const history = ref<any[]>([])

// 获取当前用户名
const getCurrentUsername = (): string => {
  return localStorage.getItem('username') || 'unknown'
}

// 获取当前用户ID
const getCurrentUserId = (): number => {
  const userId = localStorage.getItem('userId')
  return userId ? parseInt(userId) : 0
}

// 从服务器获取压缩历史记录
const fetchHistories = async () => {
  try {
    // 获取本地存储的历史记录，包含完整图片数据
    const localHistory = JSON.parse(localStorage.getItem('imageCompressHistory') || '[]')
    
    // 如果本地有历史记录，直接使用本地记录，不依赖服务器返回的图片数据
    if (localHistory.length > 0) {
      history.value = localHistory
      return
    }
    
    // 如果本地没有历史记录，再从服务器获取
    const response = await imageCompressAPI.getHistories()
    if (response.data && response.data.data && response.data.data.histories) {
      // 服务器返回的历史记录没有图片数据，只显示元数据
      const serverHistories = response.data.data.histories.map((item: any) => ({
        id: item.id,
        filename: item.filename,
        originalSize: item.originalSize,
        compressedSize: item.compressedSize,
        compressionRate: item.compressionRate.toFixed(1),
        width: item.width,
        height: item.height,
        quality: item.quality,
        format: item.format,
        timestamp: new Date(item.createdAt).toLocaleString(),
        image: '' // 服务器返回的历史记录没有图片数据
      }))
      
      history.value = serverHistories
    }
  } catch (error) {
    console.error('获取压缩历史失败:', error)
    // 如果从服务器获取失败，使用本地存储的历史记录
    const localHistory = JSON.parse(localStorage.getItem('imageCompressHistory') || '[]')
    if (localHistory.length > 0) {
      history.value = localHistory
    } else {
      ElMessage.error('获取压缩历史失败')
    }
  }
}

// 组件挂载时获取压缩历史记录
onMounted(() => {
  fetchHistories()
})

// 处理文件上传
const handleFileChange = (file: any) => {
  const reader = new FileReader()
  reader.onload = (e) => {
    originalImage.value = e.target?.result as string
    originalFileName.value = file.name
    originalSize.value = file.size
  }
  reader.readAsDataURL(file.raw)
}

// 处理图片加载，获取原图尺寸
const handleImageLoad = (e: Event) => {
  const img = e.target as HTMLImageElement
  originalWidth.value = img.width
  originalHeight.value = img.height
  customWidth.value = img.width
  customHeight.value = img.height
}

// 处理压缩
const handleCompress = async () => {
  if (!originalImage.value) {
    ElMessage.warning('请先上传图片')
    return
  }

  isCompressing.value = true

  try {
    // 创建图片对象
    const img = new Image()
    img.src = originalImage.value
    await new Promise((resolve) => {
      img.onload = resolve
    })

    // 计算压缩后的尺寸
    let targetWidth = originalWidth.value
    let targetHeight = originalHeight.value

    if (resizeMode.value === 'scale') {
      targetWidth = Math.round(originalWidth.value * scale.value / 100)
      targetHeight = Math.round(originalHeight.value * scale.value / 100)
    } else if (resizeMode.value === 'custom') {
      if (customWidth.value && customHeight.value) {
        targetWidth = customWidth.value
        targetHeight = customHeight.value
      } else if (customWidth.value) {
        targetWidth = customWidth.value
        targetHeight = Math.round(targetWidth * originalHeight.value / originalWidth.value)
        customHeight.value = targetHeight
      } else if (customHeight.value) {
        targetHeight = customHeight.value
        targetWidth = Math.round(targetHeight * originalWidth.value / originalHeight.value)
        customWidth.value = targetWidth
      }
    }

    // 创建Canvas
    const canvas = document.createElement('canvas')
    canvas.width = targetWidth
    canvas.height = targetHeight
    const ctx = canvas.getContext('2d')

    if (ctx) {
      // 如果输出格式是JPEG，设置白色背景
      if (outputFormat.value === 'image/jpeg') {
        ctx.fillStyle = '#ffffff'
        ctx.fillRect(0, 0, canvas.width, canvas.height)
      }
      // 绘制图片
      ctx.drawImage(img, 0, 0, targetWidth, targetHeight)

      // 压缩图片
      const qualityValue = quality.value / 100
      const dataUrl = canvas.toDataURL(outputFormat.value, qualityValue)

      // 计算压缩后的大小
      const byteString = atob(dataUrl.split(',')[1])
      const ab = new ArrayBuffer(byteString.length)
      const ia = new Uint8Array(ab)
      for (let i = 0; i < byteString.length; i++) {
        ia[i] = byteString.charCodeAt(i)
      }
      const blob = new Blob([ab], { type: outputFormat.value })
      compressedSize.value = blob.size

      // 更新压缩结果
      compressedImage.value = dataUrl
      compressedWidth.value = targetWidth
      compressedHeight.value = targetHeight

      // 添加到压缩历史
      const historyItem: any = {
        filename: originalFileName.value,
        originalSize: originalSize.value,
        compressedSize: compressedSize.value,
        compressionRate: ((1 - compressedSize.value / originalSize.value) * 100).toFixed(1),
        image: dataUrl,
        width: targetWidth,
        height: targetHeight,
        quality: quality.value,
        format: outputFormat.value,
        timestamp: new Date().toLocaleString()
      }
      
      // 保存到服务器
      try {
        const response = await imageCompressAPI.createHistory({
          UserID: getCurrentUserId(),
          Username: getCurrentUsername(),
          Filename: originalFileName.value,
          OriginalSize: originalSize.value,
          CompressedSize: compressedSize.value,
          CompressionRate: (1 - compressedSize.value / originalSize.value) * 100,
          Width: targetWidth,
          Height: targetHeight,
          Quality: quality.value,
          Format: outputFormat.value
        })
        
        // 如果服务器返回了ID，添加到历史记录
        if (response.data && response.data.data && response.data.data.id) {
          historyItem.id = response.data.data.id
        }
        
        // 添加到本地历史记录
      history.value.unshift(historyItem)
      if (history.value.length > 10) {
        history.value.pop()
      }
      
      // 保存到本地存储，保留图片数据
      localStorage.setItem('imageCompressHistory', JSON.stringify(history.value))
      } catch (error) {
        console.error('保存压缩历史失败:', error)
        // 即使保存失败，也将历史记录添加到本地
        history.value.unshift(historyItem)
        if (history.value.length > 10) {
          history.value.pop()
        }
        
        // 保存到本地存储，保留图片数据
        localStorage.setItem('imageCompressHistory', JSON.stringify(history.value))
      }

      ElMessage.success('图片压缩成功')
    }
  } catch (error) {
    ElMessage.error('图片压缩失败')
    console.error('压缩失败:', error)
  } finally {
    isCompressing.value = false
  }
}

// 处理调整尺寸模式变化
const handleResizeModeChange = () => {
  if (resizeMode.value === 'scale') {
    scale.value = 100
  } else if (resizeMode.value === 'custom') {
    customWidth.value = originalWidth.value
    customHeight.value = originalHeight.value
  }
  handleCompress()
}

// 重置所有设置
const resetAll = () => {
  fileList.value = []
  originalImage.value = ''
  originalFileName.value = ''
  originalSize.value = 0
  originalWidth.value = 0
  originalHeight.value = 0
  compressedImage.value = ''
  compressedSize.value = 0
  compressedWidth.value = 0
  compressedHeight.value = 0
  quality.value = 80
  outputFormat.value = 'image/jpeg'
  resizeMode.value = 'original'
  scale.value = 100
  customWidth.value = 0
  customHeight.value = 0
  lockAspectRatio.value = true
}

// 格式化文件大小
const formatFileSize = (size: number): string => {
  if (size === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(size) / Math.log(k))
  return parseFloat((size / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 下载图片
const downloadImage = () => {
  if (!compressedImage.value) return

  const link = document.createElement('a')
  link.href = compressedImage.value
  const ext = outputFormat.value.split('/')[1]
  link.download = `${originalFileName.value.split('.')[0]}_compressed.${ext}`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  ElMessage.success('图片下载成功')
}

// 复制图片链接
const copyImage = () => {
  if (!compressedImage.value) return

  navigator.clipboard.writeText(compressedImage.value)
    .then(() => {
      ElMessage.success('图片链接已复制到剪贴板')
    })
    .catch(err => {
      console.error('复制失败:', err)
      ElMessage.error('复制失败，请手动复制')
    })
}

// 预览历史记录
const previewHistory = (item: any) => {
  compressedImage.value = item.image
  compressedSize.value = item.compressedSize
  compressedWidth.value = item.width
  compressedHeight.value = item.height
  outputFormat.value = item.format
  // 设置原图信息，方便用户查看对比
  originalFileName.value = item.filename
  originalSize.value = item.originalSize
}

// 下载历史记录
const downloadHistory = (item: any) => {
  if (!item.image) {
    ElMessage.warning('该历史记录没有图片数据，无法下载')
    return
  }
  
  const link = document.createElement('a')
  link.href = item.image
  // 使用历史记录中的格式来确定文件扩展名
  const ext = item.format.split('/')[1]
  link.download = `${item.filename.split('.')[0]}_compressed.${ext}`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  ElMessage.success('历史图片下载成功')
}

// 清空历史记录
const clearHistory = async () => {
  try {
    // 调用API清空服务器上的历史记录
    await imageCompressAPI.clearHistories()
    // 清空本地历史记录
    history.value = []
    // 清空本地存储中的历史记录
    localStorage.removeItem('imageCompressHistory')
    ElMessage.success('压缩历史已清空')
  } catch (error) {
    console.error('清空压缩历史失败:', error)
    ElMessage.error('清空压缩历史失败，请重试')
  }
}
</script>

<style scoped>
.image-compress-container {
  padding: 20px;
  background-color: #f5f7fa;
  min-height: calc(100vh - 64px);
}

.main-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 10px;
}

.card-header h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #303133;
}

.sub-card-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.sub-card-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 500;
  color: #303133;
}

.upload-card, .settings-card, .result-card, .history-card {
  margin-bottom: 20px;
}

.upload-area {
  margin-bottom: 20px;
}

.image-preview {
  margin-top: 20px;
}

.image-preview h4 {
  margin: 0 0 15px 0;
  font-size: 14px;
  font-weight: 500;
  color: #606266;
}

.preview-container {
  width: 100%;
  height: 300px;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #fafafa;
  border: 1px dashed #dcdfe6;
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 15px;
}

.preview-image {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.image-info {
  background-color: #f0f2f5;
  padding: 15px;
  border-radius: 4px;
  font-size: 13px;
}

.image-info p {
  margin: 5px 0;
  color: #606266;
}

.compression-rate {
  color: #67c23a !important;
  font-weight: 500;
}

.settings-content {
  padding: 10px 0;
}

.setting-item {
  margin-bottom: 20px;
}

.setting-label {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.setting-label span:first-child {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
}

.setting-label .value {
  font-size: 13px;
  color: #67c23a;
  font-weight: 500;
}

.custom-size {
  margin-top: 15px;
}

.size-inputs {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.size-separator {
  font-size: 16px;
  color: #909399;
  font-weight: 500;
}

.aspect-ratio-btn {
  margin-left: 10px;
}

.compress-btn-container {
  display: flex;
  gap: 10px;
  margin-top: 20px;
  justify-content: center;
}

.no-result {
  padding: 60px 0;
  text-align: center;
}

.result-area {
  padding: 20px 0;
}

.result-actions {
  display: flex;
  gap: 10px;
  margin-top: 20px;
  justify-content: center;
}

.history-list {
  max-height: 300px;
  overflow-y: auto;
}

.history-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px;
  background-color: #f0f2f5;
  border-radius: 4px;
  margin-bottom: 10px;
}

.history-info {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.history-filename {
  font-weight: 500;
  color: #303133;
}

.history-size {
  font-size: 13px;
  color: #606266;
}

.history-rate {
  font-size: 12px;
  color: #67c23a;
}

.history-actions {
  display: flex;
  gap: 5px;
}

.no-history {
  padding: 40px 0;
  text-align: center;
}
</style>