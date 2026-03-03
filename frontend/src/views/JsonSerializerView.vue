<template>
  <div class="json-serializer-view">
    <el-card shadow="never" class="page-header">
      <div class="header-content">
        <h2>JSON/Base64</h2>
        <div class="header-actions">
          <el-button type="primary" size="small" @click="reset">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
        </div>
      </div>
    </el-card>

    <el-card shadow="hover" class="json-serializer-card">
      <div class="serializer-content">
        <div class="input-section">
          <div class="section-header">
            <h3>输入数据</h3>
            <div class="input-type-selector">
              <el-radio-group v-model="inputType" size="small">
              <el-radio label="json">JSON</el-radio>
              <el-radio label="text">文本</el-radio>
              <el-radio label="base64">Base64</el-radio>
            </el-radio-group>
            </div>
          </div>
          <el-input
            v-model="inputData"
            type="textarea"
            :rows="15"
            placeholder="请输入要转换的数据..."
            class="data-input"
          ></el-input>
        </div>

        <div class="action-section">
          <el-button
            type="primary"
            size="large"
            @click="convert"
            :loading="converting"
            class="convert-btn"
          >
            <el-icon><ArrowRight /></el-icon>
            转换
          </el-button>
          <el-button
            type="success"
            size="large"
            @click="copyOutput"
            :disabled="!outputData"
            class="copy-btn"
          >
            <el-icon><DocumentCopy /></el-icon>
            复制结果
          </el-button>
        </div>

        <div class="output-section">
          <div class="section-header">
            <h3>输出结果</h3>
            <div class="output-type-selector">
              <el-radio-group v-model="outputType" size="small">
              <el-radio label="json">JSON</el-radio>
              <el-radio label="text">文本</el-radio>
              <el-radio label="base64">Base64</el-radio>
            </el-radio-group>
            </div>
          </div>
          <el-input
            v-model="outputData"
            type="textarea"
            :rows="15"
            placeholder="转换结果将显示在这里..."
            class="data-output"
            readonly
          ></el-input>
          <div class="output-info" v-if="outputInfo">
            <el-tag type="info">{{ outputInfo }}</el-tag>
          </div>
        </div>
      </div>
    </el-card>

    <el-card shadow="hover" class="features-card">
      <template #header>
        <div class="card-header">
          <h3>功能特性</h3>
        </div>
      </template>
      <div class="features-grid">
        <div class="feature-item">
          <el-icon class="feature-icon"><Document /></el-icon>
          <h4>JSON格式化</h4>
          <p>美化或压缩JSON数据，支持语法高亮</p>
        </div>

        <div class="feature-item">
          <el-icon class="feature-icon"><Check /></el-icon>
          <h4>语法验证</h4>
          <p>自动验证JSON语法，提示错误位置</p>
        </div>
        <div class="feature-item">
          <el-icon class="feature-icon"><CopyDocument /></el-icon>
          <h4>一键复制</h4>
          <p>快速复制转换结果到剪贴板</p>
        </div>
        <div class="feature-item">
          <el-icon class="feature-icon"><Lock /></el-icon>
          <h4>Base64编码</h4>
          <p>将文本或JSON数据编码为Base64格式</p>
        </div>
        <div class="feature-item">
          <el-icon class="feature-icon"><Unlock /></el-icon>
          <h4>Base64解码</h4>
          <p>将Base64编码数据解码为文本或JSON</p>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Refresh,
  ArrowRight,
  DocumentCopy,
  Document,
  Check,
  CopyDocument,
  Lock,
  Unlock
} from '@element-plus/icons-vue'

// 输入输出类型
const inputType = ref('json')
const outputType = ref('json')

// 输入输出数据
const inputData = ref('')
const outputData = ref('')

// 转换状态
const converting = ref(false)
const outputInfo = ref('')

// Base64编码函数
const base64Encode = (str: string): string => {
  return btoa(unescape(encodeURIComponent(str)))
}

// Base64解码函数
const base64Decode = (str: string): string => {
  try {
    return decodeURIComponent(escape(atob(str)))
  } catch (error) {
    throw new Error('无效的Base64编码')
  }
}

// 转换函数
const convert = async () => {
  if (!inputData.value.trim()) {
    ElMessage.warning('请输入要转换的数据')
    return
  }

  converting.value = true
  outputData.value = ''
  outputInfo.value = ''

  try {
    let result = ''

    if (inputType.value === 'json') {
      // JSON输入处理
      if (outputType.value === 'json') {
        // JSON -> JSON（格式化/压缩）
        const parsed = JSON.parse(inputData.value)
        result = JSON.stringify(parsed, null, 2)
        outputInfo.value = 'JSON格式化完成'
      } else if (outputType.value === 'text') {
        // JSON -> 文本
        const parsed = JSON.parse(inputData.value)
        result = JSON.stringify(parsed)
        outputInfo.value = 'JSON转文本完成'
      } else if (outputType.value === 'base64') {
        // JSON -> Base64
        const parsed = JSON.parse(inputData.value)
        const jsonStr = JSON.stringify(parsed)
        result = base64Encode(jsonStr)
        outputInfo.value = 'JSON转Base64完成'
      }
    } else if (inputType.value === 'text') {
      // 文本输入处理
      if (outputType.value === 'json') {
        // 文本 -> JSON（尝试解析为JSON）
        try {
          const parsed = JSON.parse(inputData.value)
          result = JSON.stringify(parsed, null, 2)
          outputInfo.value = '文本转JSON完成'
        } catch (e) {
          // 文本不是有效的JSON，直接输出
          result = inputData.value
          outputInfo.value = '文本原样输出'
        }
      } else if (outputType.value === 'text') {
        // 文本 -> 文本（原样输出）
        result = inputData.value
        outputInfo.value = '文本原样输出'
      } else if (outputType.value === 'base64') {
        // 文本 -> Base64
        result = base64Encode(inputData.value)
        outputInfo.value = '文本转Base64完成'
      }
    } else if (inputType.value === 'base64') {
      // Base64输入处理
      const decodedText = base64Decode(inputData.value)
      
      if (outputType.value === 'json') {
        // Base64 -> JSON
        try {
          const parsed = JSON.parse(decodedText)
          result = JSON.stringify(parsed, null, 2)
          outputInfo.value = 'Base64转JSON完成'
        } catch (e) {
          // 解码后的文本不是有效的JSON，输出解码后的文本
          result = decodedText
          outputInfo.value = 'Base64解码完成（非JSON）'
        }
      } else if (outputType.value === 'text') {
        // Base64 -> 文本
        result = decodedText
        outputInfo.value = 'Base64转文本完成'
      } else if (outputType.value === 'base64') {
        // Base64 -> Base64（重新编码）
        result = base64Encode(decodedText)
        outputInfo.value = 'Base64重新编码完成'
      }
    }

    outputData.value = result
    ElMessage.success('转换成功')
  } catch (error: any) {
    outputData.value = `转换失败: ${error.message}`
    outputInfo.value = '转换失败'
    ElMessage.error('转换失败，请检查输入数据格式')
  } finally {
    converting.value = false
  }
}

// 复制结果
const copyOutput = async () => {
  try {
    await navigator.clipboard.writeText(outputData.value)
    ElMessage.success('复制成功')
  } catch (error) {
    ElMessage.error('复制失败，请手动复制')
  }
}

// 重置
const reset = () => {
  inputData.value = ''
  outputData.value = ''
  outputInfo.value = ''
  inputType.value = 'json'
  outputType.value = 'json'
}
</script>

<style scoped>
.json-serializer-view {
  padding: 20px;
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

.json-serializer-card {
  margin-bottom: 20px;
}

.serializer-content {
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  gap: 20px;
  align-items: center;
}

.input-section,
.output-section {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.section-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 500;
  color: #333;
}

.data-input {
  font-family: 'Courier New', Courier, monospace;
  font-size: 14px;
  min-height: 300px;
}

.action-section {
  display: flex;
  flex-direction: column;
  gap: 15px;
  align-items: center;
  justify-content: center;
  padding: 10px 0;
}

.convert-btn,
.copy-btn {
  width: 140px;
  height: 40px;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0;
  border-radius: 4px;
  transition: all 0.3s ease;
  font-size: 14px;
}

.convert-btn {
  order: 1;
}

.copy-btn {
  order: 2;
}

/* 确保按钮内部图标和文字对齐 */
.convert-btn .el-icon,
.copy-btn .el-icon {
  margin-right: 6px;
  font-size: 16px;
}

/* 响应式布局调整 */
@media (max-width: 768px) {
  .action-section {
    flex-direction: row;
    gap: 10px;
  }
  
  .convert-btn,
  .copy-btn {
    width: 120px;
    height: 36px;
    font-size: 13px;
  }
}

.features-card {
  margin-top: 20px;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  margin-top: 20px;
}

.feature-item {
  text-align: center;
  padding: 20px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.feature-item:hover {
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

.feature-icon {
  font-size: 32px;
  color: #409eff;
  margin-bottom: 10px;
}

.feature-item h4 {
  margin: 0 0 10px 0;
  font-size: 16px;
  font-weight: 500;
}

.feature-item p {
  margin: 0;
  font-size: 14px;
  color: #606266;
}

.output-info {
  margin-top: 10px;
}

@media (max-width: 768px) {
  .serializer-content {
    grid-template-columns: 1fr;
  }

  .action-section {
    flex-direction: row;
    justify-content: center;
  }

  .features-grid {
    grid-template-columns: 1fr;
  }
}
</style>