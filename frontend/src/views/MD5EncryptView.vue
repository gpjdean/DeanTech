<template>
  <div class="md5-encrypt-container">
    <el-card class="main-card">
      <template #header>
        <div class="card-header">
          <el-icon><Key /></el-icon>
          <h2>MD5加解密工具</h2>
        </div>
      </template>

      <div class="encrypt-content">
        <!-- 输入区域 -->
        <div class="input-section">
          <el-card shadow="hover" class="sub-card">
            <template #header>
              <div class="sub-card-header">
                <el-icon><EditPen /></el-icon>
                <h3>输入内容</h3>
              </div>
            </template>
            <div class="input-area">
              <el-input
                v-model="inputText"
                type="textarea"
                :rows="6"
                placeholder="请输入要加密的文本..."
                resize="vertical"
              />
            </div>
            <div class="action-buttons">
              <el-button type="primary" size="large" @click="handleEncrypt" :loading="isEncrypting">
                <el-icon><Upload /></el-icon>
                生成MD5
              </el-button>
              <el-button type="default" size="large" @click="clearInput">
                <el-icon><Delete /></el-icon>
                清空
              </el-button>
              <el-button type="success" size="large" @click="copyInput" :disabled="!inputText">
                <el-icon><DocumentCopy /></el-icon>
                复制输入
              </el-button>
            </div>
          </el-card>
        </div>

        <!-- 输出区域 -->
        <div class="output-section">
          <el-card shadow="hover" class="sub-card">
            <template #header>
              <div class="sub-card-header">
                <el-icon><DocumentCopy /></el-icon>
                <h3>MD5结果</h3>
              </div>
            </template>
            <div class="output-area">
              <div v-if="md5Results.length > 0" class="result-list">
                <div
                  v-for="(result, index) in md5Results"
                  :key="index"
                  class="result-item"
                >
                  <el-input
                    v-model="result.value"
                    type="textarea"
                    :rows="1"
                    readonly
                    class="result-input"
                  />
                  <div class="result-actions">
                    <el-button
                      type="success"
                      size="small"
                      @click="copyResult(result.value)"
                    >
                      <el-icon><DocumentCopy /></el-icon>
                      复制
                    </el-button>
                    <el-button
                      type="danger"
                      size="small"
                      @click="removeResult(index)"
                    >
                      <el-icon><Delete /></el-icon>
                      删除
                    </el-button>
                  </div>
                </div>
              </div>
              <div v-else class="no-result">
                <el-empty description="点击生成MD5按钮开始加密" />
              </div>
            </div>
            <div class="result-info" v-if="md5Results.length > 0">
              <p>已生成 {{ md5Results.length }} 个MD5结果</p>
              <el-button type="warning" size="small" @click="clearResults">
                <el-icon><Delete /></el-icon>
                清空所有结果
              </el-button>
            </div>
          </el-card>
        </div>
      </div>

      <!-- 选项区域 -->
      <el-card shadow="hover" class="options-card">
        <template #header>
          <div class="sub-card-header">
            <el-icon><Setting /></el-icon>
            <h3>加密选项</h3>
          </div>
        </template>
        <div class="options-content">
          <el-row :gutter="20">
            <el-col :span="8">
              <div class="option-item">
                <div class="option-label">
                  <span>输出格式</span>
                </div>
                <el-radio-group v-model="outputFormat">
                  <el-radio-button label="lowercase">小写</el-radio-button>
                  <el-radio-button label="uppercase">大写</el-radio-button>
                </el-radio-group>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="option-item">
                <div class="option-label">
                  <span>字符集</span>
                </div>
                <el-radio-group v-model="charset">
                  <el-radio-button label="utf8">UTF-8</el-radio-button>
                  <el-radio-button label="gbk">GBK</el-radio-button>
                </el-radio-group>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="option-item">
                <div class="option-label">
                  <span>自动生成</span>
                </div>
                <el-switch
                  v-model="autoGenerate"
                  active-text="开启"
                  inactive-text="关闭"
                />
              </div>
            </el-col>
          </el-row>
        </div>
      </el-card>

      <!-- 历史记录 -->
      <el-card shadow="hover" class="history-card">
        <template #header>
          <div class="sub-card-header">
            <el-icon><Clock /></el-icon>
            <h3>加密历史</h3>
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
                <div class="history-input">
                  <span class="history-label">输入:</span>
                  <span class="history-content">{{ item.input }}</span>
                </div>
                <div class="history-output">
                  <span class="history-label">输出:</span>
                  <span class="history-content md5-result">{{ item.output }}</span>
                </div>
                <div class="history-actions">
                  <el-button
                    type="text"
                    size="small"
                    @click="reuseInput(item.input)"
                  >
                    <el-icon><EditPen /></el-icon>
                    重新使用
                  </el-button>
                  <el-button
                    type="text"
                    size="small"
                    @click="copyResult(item.output)"
                  >
                    <el-icon><DocumentCopy /></el-icon>
                    复制结果
                  </el-button>
                </div>
              </div>
            </el-timeline-item>
          </el-timeline>
        </div>
        <div v-else class="no-history">
          <el-empty description="暂无加密历史" />
        </div>
      </el-card>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Key,
  EditPen,
  Upload,
  Delete,
  DocumentCopy,
  Setting,
  Clock
} from '@element-plus/icons-vue'
import { md5API } from '../api/api'

// 输入内容
const inputText = ref<string>('')

// 加密结果
interface MD5Result {
  value: string
}
const md5Results = ref<MD5Result[]>([])
const isEncrypting = ref<boolean>(false)

// 加密选项
const outputFormat = ref<string>('lowercase') // lowercase | uppercase
const charset = ref<string>('utf8') // utf8 | gbk
const autoGenerate = ref<boolean>(false)

// 加密历史
interface HistoryItem {
  id?: number
  input: string
  output: string
  timestamp: string
  format?: string
  charset?: string
}
const history = ref<HistoryItem[]>([])

// 监听输入变化，自动生成MD5
watch(inputText, (newValue) => {
  if (autoGenerate.value && newValue) {
    handleEncrypt()
  }
})

// 获取当前用户名
const getCurrentUsername = (): string => {
  return localStorage.getItem('username') || 'unknown'
}

// 获取当前用户ID
const getCurrentUserId = (): number => {
  const userId = localStorage.getItem('userId')
  return userId ? parseInt(userId) : 0
}

// 从服务器获取加密历史记录
const fetchHistories = async () => {
  try {
    const response = await md5API.getHistories()
    if (response.data && response.data.data && response.data.data.histories) {
      // 将服务器返回的历史记录转换为前端需要的格式
      history.value = response.data.data.histories.map((item: any) => ({
        id: item.id,
        input: item.input,
        output: item.output,
        timestamp: new Date(item.createdAt).toLocaleString(),
        format: item.format,
        charset: item.charset
      }))
    }
  } catch (error) {
    console.error('获取加密历史失败:', error)
    ElMessage.error('获取加密历史失败')
  }
}

// 组件挂载时获取加密历史记录
onMounted(() => {
  fetchHistories()
})

// 纯JavaScript实现的MD5算法
const md5 = (message: string): string => {
  // 辅助函数
  const rotateLeft = (lValue: number, iShiftBits: number): number => {
    return (lValue << iShiftBits) | (lValue >>> (32 - iShiftBits))
  }

  const addUnsigned = (lX: number, lY: number): number => {
    const lX4 = lX & 0x40000000
    const lY4 = lY & 0x40000000
    const lX8 = lX & 0x80000000
    const lY8 = lY & 0x80000000
    const lResult = (lX & 0x3FFFFFFF) + (lY & 0x3FFFFFFF)
    if (lX4 & lY4) {
      return (lResult ^ 0x80000000 ^ lX8 ^ lY8)
    }
    if (lX4 | lY4) {
      if (lResult & 0x40000000) {
        return (lResult ^ 0xC0000000 ^ lX8 ^ lY8)
      } else {
        return (lResult ^ 0x40000000 ^ lX8 ^ lY8)
      }
    } else {
      return (lResult ^ lX8 ^ lY8)
    }
  }

  const F = (x: number, y: number, z: number): number => {
    return (x & y) | ((~x) & z)
  }

  const G = (x: number, y: number, z: number): number => {
    return (x & z) | (y & (~z))
  }

  const H = (x: number, y: number, z: number): number => {
    return x ^ y ^ z
  }

  const I = (x: number, y: number, z: number): number => {
    return y ^ (x | (~z))
  }

  const FF = (a: number, b: number, c: number, d: number, x: number, s: number, ac: number): number => {
    a = addUnsigned(a, addUnsigned(addUnsigned(F(b, c, d), x), ac))
    return addUnsigned(rotateLeft(a, s), b)
  }

  const GG = (a: number, b: number, c: number, d: number, x: number, s: number, ac: number): number => {
    a = addUnsigned(a, addUnsigned(addUnsigned(G(b, c, d), x), ac))
    return addUnsigned(rotateLeft(a, s), b)
  }

  const HH = (a: number, b: number, c: number, d: number, x: number, s: number, ac: number): number => {
    a = addUnsigned(a, addUnsigned(addUnsigned(H(b, c, d), x), ac))
    return addUnsigned(rotateLeft(a, s), b)
  }

  const II = (a: number, b: number, c: number, d: number, x: number, s: number, ac: number): number => {
    a = addUnsigned(a, addUnsigned(addUnsigned(I(b, c, d), x), ac))
    return addUnsigned(rotateLeft(a, s), b)
  }

  const convertToWordArray = (input: string): number[] => {
    const lWordCount = Math.ceil(input.length / 4)
    const lWordArray = new Array(lWordCount)
    for (let i = 0; i < lWordCount; i++) {
      lWordArray[i] = 0
      for (let j = 0; j < 4; j++) {
        const charCode = input.charCodeAt(i * 4 + j)
        if (charCode !== undefined) {
          lWordArray[i] |= charCode << (8 * j)
        }
      }
    }
    const lMessageLength = input.length * 8
    const lNumberOfBitsBeforePadding = lMessageLength
    lWordArray[lWordCount] = lNumberOfBitsBeforePadding & 0xFFFFFFFF
    lWordArray[lWordCount + 1] = lNumberOfBitsBeforePadding >>> 32
    return lWordArray
  }

  const wordToHex = (lValue: number): string => {
    let lHexValue = ''
    const lHexChars = '0123456789abcdef'
    for (let i = 0; i <= 3; i++) {
      lHexValue += lHexChars.charAt((lValue >>> (i * 8 + 4)) & 0x0F) + lHexChars.charAt((lValue >>> (i * 8)) & 0x0F)
    }
    return lHexValue
  }

  let x = convertToWordArray(message)
  let a = 0x67452301
  let b = 0xEFCDAB89
  let c = 0x98BADCFE
  let d = 0x10325476
  let olda: number = a, oldb: number = b, oldc: number = c, oldd: number = d

  // 第一轮
  a = FF(a, b, c, d, x[0], 7, 0xD76AA478)
  d = FF(d, a, b, c, x[1], 12, 0xE8C7B756)
  c = FF(c, d, a, b, x[2], 17, 0x242070DB)
  b = FF(b, c, d, a, x[3], 22, 0xC1BDCEEE)
  a = FF(a, b, c, d, x[4], 7, 0xF57C0FAF)
  d = FF(d, a, b, c, x[5], 12, 0x4787C62A)
  c = FF(c, d, a, b, x[6], 17, 0xA8304613)
  b = FF(b, c, d, a, x[7], 22, 0xFD469501)
  a = FF(a, b, c, d, x[8], 7, 0x698098D8)
  d = FF(d, a, b, c, x[9], 12, 0x8B44F7AF)
  c = FF(c, d, a, b, x[10], 17, 0xFFFF5BB1)
  b = FF(b, c, d, a, x[11], 22, 0x895CD7BE)
  a = FF(a, b, c, d, x[12], 7, 0x6B901122)
  d = FF(d, a, b, c, x[13], 12, 0xFD987193)
  c = FF(c, d, a, b, x[14], 17, 0xA679438E)
  b = FF(b, c, d, a, x[15], 22, 0x49B40821)

  // 第二轮
  a = GG(a, b, c, d, x[1], 5, 0xF61E2562)
  d = GG(d, a, b, c, x[6], 9, 0xC040B340)
  c = GG(c, d, a, b, x[11], 14, 0x265E5A51)
  b = GG(b, c, d, a, x[0], 20, 0xE9B6C7AA)
  a = GG(a, b, c, d, x[5], 5, 0xD62F105D)
  d = GG(d, a, b, c, x[10], 9, 0x02441453)
  c = GG(c, d, a, b, x[15], 14, 0xD8A1E681)
  b = GG(b, c, d, a, x[4], 20, 0xE7D3FBC8)
  a = GG(a, b, c, d, x[9], 5, 0x21E1CDE6)
  d = GG(d, a, b, c, x[14], 9, 0xC33707D6)
  c = GG(c, d, a, b, x[3], 14, 0xF4D50D87)
  b = GG(b, c, d, a, x[8], 20, 0x455A14ED)
  a = GG(a, b, c, d, x[13], 5, 0xA9E3E905)
  d = GG(d, a, b, c, x[2], 9, 0xFCEFA3F8)
  c = GG(c, d, a, b, x[7], 14, 0x676F02D9)
  b = GG(b, c, d, a, x[12], 20, 0x8D2A4C8A)

  // 第三轮
  a = HH(a, b, c, d, x[5], 4, 0xFFFA3942)
  d = HH(d, a, b, c, x[8], 11, 0x8771F681)
  c = HH(c, d, a, b, x[11], 16, 0x6D9D6122)
  b = HH(b, c, d, a, x[14], 23, 0xFDE5380C)
  a = HH(a, b, c, d, x[1], 4, 0xA4BEEA44)
  d = HH(d, a, b, c, x[4], 11, 0x4BDECFA9)
  c = HH(c, d, a, b, x[7], 16, 0xF6BB4B60)
  b = HH(b, c, d, a, x[10], 23, 0xBEBFBC70)
  a = HH(a, b, c, d, x[13], 4, 0x289B7EC6)
  d = HH(d, a, b, c, x[0], 11, 0xEAA127FA)
  c = HH(c, d, a, b, x[3], 16, 0xD4EF3085)
  b = HH(b, c, d, a, x[6], 23, 0x04881D05)
  a = HH(a, b, c, d, x[9], 4, 0xD9D4D039)
  d = HH(d, a, b, c, x[12], 11, 0xE6DB99E5)
  c = HH(c, d, a, b, x[15], 16, 0x1FA27CF8)
  b = HH(b, c, d, a, x[2], 23, 0xC4AC5665)

  // 第四轮
  a = II(a, b, c, d, x[0], 6, 0xF4292244)
  d = II(d, a, b, c, x[7], 10, 0x432AFF97)
  c = II(c, d, a, b, x[14], 15, 0xAB9423A7)
  b = II(b, c, d, a, x[5], 21, 0xFC93A039)
  a = II(a, b, c, d, x[12], 6, 0x655B59C3)
  d = II(d, a, b, c, x[3], 10, 0x8F0CCC92)
  c = II(c, d, a, b, x[10], 15, 0xFFEFF47D)
  b = II(b, c, d, a, x[1], 21, 0x85845DD1)
  a = II(a, b, c, d, x[8], 6, 0x6FA87E4F)
  d = II(d, a, b, c, x[15], 10, 0xFE2CE6E0)
  c = II(c, d, a, b, x[6], 15, 0xA3014314)
  b = II(b, c, d, a, x[13], 21, 0x4E0811A1)
  a = II(a, b, c, d, x[4], 6, 0xF7537E82)
  d = II(d, a, b, c, x[11], 10, 0xBD3AF235)
  c = II(c, d, a, b, x[2], 15, 0x2AD7D2BB)
  b = II(b, c, d, a, x[9], 21, 0xEB86D391)

  a = addUnsigned(a, olda || 0x67452301)
  b = addUnsigned(b, oldb || 0xEFCDAB89)
  c = addUnsigned(c, oldc || 0x98BADCFE)
  d = addUnsigned(d, oldd || 0x10325476)

  return wordToHex(a) + wordToHex(b) + wordToHex(c) + wordToHex(d)
}

// 生成MD5
const handleEncrypt = async () => {
  if (!inputText.value.trim()) {
    ElMessage.warning('请输入要加密的文本')
    return
  }

  isEncrypting.value = true

  try {
    // 使用纯JavaScript实现的MD5算法
    let hashHex = md5(inputText.value)

    // 根据输出格式转换大小写
    if (outputFormat.value === 'uppercase') {
      hashHex = hashHex.toUpperCase()
    }

    // 添加到结果列表
    const result: MD5Result = {
      value: hashHex
    }
    md5Results.value.unshift(result)
    if (md5Results.value.length > 10) {
      md5Results.value.pop()
    }

    // 创建历史记录对象
    const historyItem: HistoryItem = {
      input: inputText.value,
      output: hashHex,
      timestamp: new Date().toLocaleString(),
      format: outputFormat.value,
      charset: charset.value
    }

    // 保存到服务器
    try {
      const response = await md5API.createHistory({
        UserID: getCurrentUserId(),
        Username: getCurrentUsername(),
        Input: inputText.value,
        Output: hashHex,
        Format: outputFormat.value,
        Charset: charset.value
      })
      
      // 如果服务器返回了ID，添加到历史记录
      if (response.data && response.data.data && response.data.data.id) {
        historyItem.id = response.data.data.id
      }
      
      // 添加到本地历史记录
      history.value.unshift(historyItem)
      if (history.value.length > 20) {
        history.value.pop()
      }
    } catch (error) {
      console.error('保存加密历史失败:', error)
      // 即使保存失败，也将历史记录添加到本地
      history.value.unshift(historyItem)
      if (history.value.length > 20) {
        history.value.pop()
      }
    }

    ElMessage.success('MD5生成成功')
  } catch (error) {
    ElMessage.error('MD5生成失败')
    console.error('MD5生成失败:', error)
  } finally {
    isEncrypting.value = false
  }
}

// 清空输入
const clearInput = () => {
  inputText.value = ''
  ElMessage.success('输入已清空')
}

// 复制输入
const copyInput = () => {
  if (!inputText.value) return

  navigator.clipboard.writeText(inputText.value)
    .then(() => {
      ElMessage.success('输入内容已复制到剪贴板')
    })
    .catch(err => {
      console.error('复制失败:', err)
      ElMessage.error('复制失败，请手动复制')
    })
}

// 复制结果
const copyResult = (result: string) => {
  navigator.clipboard.writeText(result)
    .then(() => {
      ElMessage.success('MD5结果已复制到剪贴板')
    })
    .catch(err => {
      console.error('复制失败:', err)
      ElMessage.error('复制失败，请手动复制')
    })
}

// 移除结果
const removeResult = (index: number) => {
  md5Results.value.splice(index, 1)
  ElMessage.success('结果已移除')
}

// 清空所有结果
const clearResults = () => {
  md5Results.value = []
  ElMessage.success('所有结果已清空')
}

// 清空历史记录
const clearHistory = async () => {
  try {
    // 调用API清空服务器上的历史记录
    await md5API.clearHistories()
    // 清空本地历史记录
    history.value = []
    ElMessage.success('加密历史已清空')
  } catch (error) {
    console.error('清空加密历史失败:', error)
    // 即使API调用失败，也清空本地历史记录
    history.value = []
    ElMessage.success('加密历史已清空')
  }
}

// 重新使用历史输入
const reuseInput = (input: string) => {
  inputText.value = input
  ElMessage.success('已重新使用历史输入')
}
</script>

<style scoped>
.md5-encrypt-container {
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

.encrypt-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  margin-bottom: 20px;
}

.input-section, .output-section {
  display: flex;
  flex-direction: column;
}

.sub-card {
  margin-bottom: 20px;
}

.input-area, .output-area {
  margin-bottom: 20px;
}

.action-buttons {
  display: flex;
  gap: 10px;
  justify-content: center;
  flex-wrap: wrap;
}

.result-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.result-item {
  display: flex;
  gap: 10px;
  align-items: flex-start;
}

.result-input {
  flex: 1;
  min-height: 40px;
}

.result-actions {
  display: flex;
  gap: 5px;
  flex-shrink: 0;
}

.no-result {
  padding: 40px 0;
  text-align: center;
}

.result-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
  border-top: 1px solid #ebeef5;
}

.result-info p {
  margin: 0;
  color: #606266;
  font-size: 14px;
}

.options-card {
  margin-bottom: 20px;
}

.options-content {
  padding: 10px 0;
}

.option-item {
  margin-bottom: 20px;
}

.option-label {
  margin-bottom: 10px;
  font-size: 14px;
  font-weight: 500;
  color: #303133;
}

.history-card {
  margin-bottom: 20px;
}

.history-list {
  max-height: 400px;
  overflow-y: auto;
}

.history-item {
  background-color: #f0f2f5;
  padding: 15px;
  border-radius: 8px;
  margin-bottom: 10px;
}

.history-input, .history-output {
  margin-bottom: 8px;
  word-break: break-all;
}

.history-label {
  font-weight: 500;
  color: #303133;
  margin-right: 8px;
}

.history-content {
  color: #606266;
}

.md5-result {
  font-family: 'Courier New', Courier, monospace;
  color: #67c23a;
  font-weight: 500;
}

.history-actions {
  display: flex;
  gap: 5px;
  margin-top: 10px;
}

.no-history {
  padding: 40px 0;
  text-align: center;
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .encrypt-content {
    grid-template-columns: 1fr;
  }
}
</style>