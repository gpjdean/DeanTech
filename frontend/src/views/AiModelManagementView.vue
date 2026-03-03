<template>
  <div class="ai-model-management-view">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <el-icon class="header-icon"><Cpu /></el-icon>
        <div class="header-text">
          <h2>AI模型管理</h2>
          <p>管理和配置AI模型连接，支持多种AI服务提供商</p>
        </div>
        <div class="header-stats">
          <div class="stat-item">
            <span class="stat-value">{{ models.length }}</span>
            <span class="stat-label">可用模型</span>
          </div>
          <div class="stat-item">
            <span class="stat-value">{{ activeModelsCount }}</span>
            <span class="stat-label">已启用</span>
          </div>
          <div class="stat-item">
            <span class="stat-value">{{ inactiveModelsCount }}</span>
            <span class="stat-label">已禁用</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 操作按钮和搜索 -->
    <div class="toolbar">
      <div class="toolbar-left">
        <el-button type="primary" @click="showAddModelDialog" :icon="Plus">
          添加模型
        </el-button>
        <el-button @click="refreshModels" :icon="Refresh">
          刷新列表
        </el-button>
        <el-dropdown @command="handleBatchAction">
          <el-button>
            批量操作
            <el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="enable">批量启用</el-dropdown-item>
              <el-dropdown-item command="disable">批量禁用</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
      <div class="toolbar-right">
        <el-input
          v-model="searchQuery"
          placeholder="搜索模型名称或类型"
          prefix-icon="Search"
          clearable
          class="search-input"
        />
      </div>
    </div>

    <!-- 模型列表 -->
    <el-card shadow="hover" class="models-card">
      <template #header>
        <div class="card-header">
          <el-icon class="header-icon"><DataAnalysis /></el-icon>
          <span>模型列表</span>
        </div>
      </template>

      <div class="models-container">
        <el-table
          v-loading="loading"
          :data="filteredModels"
          @selection-change="handleSelectionChange"
          style="width: 100%"
          class="models-table"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column prop="name" label="模型名称" min-width="180">
            <template #default="scope">
              <div class="model-name-cell">
                <el-avatar :size="32" :src="scope.row.avatar || defaultAvatar">
                  {{ scope.row.name.charAt(0) }}
                </el-avatar>
                <div class="model-name-info">
                  <div class="model-name">{{ scope.row.name }}</div>
                  <div class="model-provider">{{ scope.row.provider }}</div>
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="type" label="模型类型" width="120">
            <template #default="scope">
              <el-tag :type="getModelTypeTagType(scope.row.type)">
                {{ scope.row.type }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100">
            <template #default="scope">
              <el-switch
                v-model="scope.row.status"
                active-text="启用"
                inactive-text="禁用"
                @change="handleModelStatusChange(scope.row)"
              />
            </template>
          </el-table-column>
          <el-table-column prop="version" label="版本" width="120" />
          <el-table-column prop="createdAt" label="创建时间" width="180" />
          <el-table-column prop="updatedAt" label="更新时间" width="180" />
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="scope">
              <div class="operation-buttons">
                <el-button type="primary" size="small" @click="testModelConnection(scope.row)" :icon="Connection">
                  测试
                </el-button>
                <el-button type="warning" size="small" @click="showEditModelDialog(scope.row)" :icon="Edit">
                  编辑
                </el-button>
                <el-button type="danger" size="small" @click="deleteModel(scope.row)" :icon="Delete">
                  删除
                </el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="filteredModels.length"
        />
      </div>
    </el-card>

    <!-- 添加/编辑模型对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="modelFormRef"
        :model="formData"
        :rules="formRules"
        label-position="top"
        label-width="120px"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="模型名称" prop="name">
              <el-input v-model="formData.name" placeholder="请输入模型名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="服务提供商" prop="provider">
              <el-select v-model="formData.provider" placeholder="请选择服务提供商">
                <el-option label="豆包" value="doubao" />
                <el-option label="DeepSeek" value="deepseek" />
                <el-option label="千问" value="qianwen" />
                <el-option label="元宝" value="yuanbao" />
                <el-option label="其他" value="other" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="模型类型" prop="type">
              <el-select v-model="formData.type" placeholder="请选择模型类型">
                <el-option label="聊天模型" value="chat" />
                <el-option label="代码模型" value="code" />
                <el-option label="图像模型" value="image" />
                <el-option label="音频模型" value="audio" />
                <el-option label="多模态模型" value="multimodal" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="模型版本" prop="version">
              <el-input v-model="formData.version" placeholder="请输入模型版本" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="API密钥" prop="apiKey">
          <el-input
            v-model="formData.apiKey"
            placeholder="请输入API密钥"
            type="password"
            show-password
          />
        </el-form-item>

        <el-form-item label="API地址" prop="apiUrl">
          <el-input
            v-model="formData.apiUrl"
            placeholder="请输入API地址"
            clearable
          />
        </el-form-item>

        <el-form-item label="模型描述">
          <el-input
            v-model="formData.description"
            type="textarea"
            placeholder="请输入模型描述"
            :rows="3"
          />
        </el-form-item>

        <el-form-item label="额外参数">
          <el-input
            v-model="formData.extraParams"
            type="textarea"
            placeholder="请输入额外参数（JSON格式）"
            :rows="4"
            monospaced
          />
        </el-form-item>

        <el-form-item label="状态">
          <el-switch v-model="formData.status" active-text="启用" inactive-text="禁用" />
        </el-form-item>
      </el-form>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitForm" :loading="submitting">
            确认
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 模型测试对话框 -->
    <el-dialog
      v-model="testDialogVisible"
      title="测试模型连接"
      width="500px"
      :close-on-click-modal="false"
    >
      <div class="test-dialog-content">
        <div v-if="testing" class="testing-indicator">
          <el-icon class="loading-icon"><Loading /></el-icon>
          <p>正在测试模型连接...</p>
        </div>
        <div v-else-if="testResult" class="test-result">
          <div :class="['result-icon', testResult.success ? 'success' : 'error']">
            <el-icon v-if="testResult.success" size="48"><svg viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg" width="48" height="48"><path d="M512 64a448 448 0 1 1 0 896 448 448 0 0 1 0-896zm195.2 316.8a32 32 0 0 0-44.8 0L480 604.8l-86.4-86.4a32 32 0 0 0-44.8 44.8l112 112a32 32 0 0 0 44.8 0l224-224a32 32 0 0 0 0-44.8z" fill="currentColor"/></svg></el-icon>
            <el-icon v-else size="48"><svg viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg" width="48" height="48"><path d="M512 64a448 448 0 1 1 0 896 448 448 0 0 1 0-896zm195.2 518.4a32 32 0 1 0-44.8-44.8L512 604.8l-140.8-140.8a32 32 0 0 0-44.8 44.8L467.2 649.6l188.8 188.8a32 32 0 0 0 44.8-44.8L512 649.6l195.2-195.2z" fill="currentColor"/></svg></el-icon>
          </div>
          <h3>{{ testResult.success ? '连接成功' : '连接失败' }}</h3>
          <div class="test-result-details">
            <pre>{{ testResult.message }}</pre>
          </div>
        </div>
        <div v-else class="test-instructions">
          <p>点击测试按钮开始测试模型连接</p>
        </div>
      </div>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="testDialogVisible = false">关闭</el-button>
          <el-button type="primary" @click="startModelTest" :loading="testing">
            开始测试
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Cpu,
  Plus,
  Refresh,
  ArrowDown,
  DataAnalysis,
  Connection,
  Edit,
  Delete,
  Loading
} from '@element-plus/icons-vue'

// 模拟数据
const defaultAvatar = '/images/03b0d39583f48206768a7534e55bcpng.png'

// 状态管理
const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const testDialogVisible = ref(false)
const testing = ref(false)
const selectedModels = ref<any[]>([])
const searchQuery = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const modelFormRef = ref()
const selectedModelForTest = ref<any>(null)
const testResult = ref<any>(null)

// 表单数据
const dialogTitle = ref('添加模型')
const isEditMode = ref(false)
const formData = reactive({
  id: '',
  name: '',
  provider: '',
  type: '',
  version: '',
  apiKey: '',
  apiUrl: '',
  description: '',
  extraParams: '',
  status: true,
  createdAt: '',
  updatedAt: ''
})

// 表单验证规则
const formRules = reactive({
  name: [{ required: true, message: '请输入模型名称', trigger: 'blur' }],
  provider: [{ required: true, message: '请选择服务提供商', trigger: 'change' }],
  type: [{ required: true, message: '请选择模型类型', trigger: 'change' }],
  version: [{ required: true, message: '请输入模型版本', trigger: 'blur' }],
  apiKey: [{ required: true, message: '请输入API密钥', trigger: 'blur' }]
})

// 模型数据
const models = ref([
  {
    id: '1',
    name: '豆包AI',
    provider: 'doubao',
    type: 'chat',
    version: 'v1.0',
    apiKey: 'sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx',
    apiUrl: 'https://api.doubao.com/v1/chat/completions',
    description: '字节跳动旗下的AI大模型',
    extraParams: '{}',
    status: true,
    createdAt: '2026-02-20 14:30:00',
    updatedAt: '2026-02-25 10:15:00'
  },
  {
    id: '2',
    name: 'DeepSeek-Coder',
    provider: 'deepseek',
    type: 'code',
    version: 'v2.0',
    apiKey: 'ds-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx',
    apiUrl: 'https://api.deepseek.com/v1/codex/completions',
    description: '专注于代码生成的AI模型',
    extraParams: '{"temperature": 0.7}',
    status: true,
    createdAt: '2026-02-21 09:45:00',
    updatedAt: '2026-02-24 16:20:00'
  },
  {
    id: '3',
    name: '通义千问',
    provider: 'qianwen',
    type: 'multimodal',
    version: 'v3.0',
    apiKey: 'qw-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx',
    apiUrl: 'https://api.qianwen.com/v1/chat/completions',
    description: '阿里云通义千问大模型',
    extraParams: '{}',
    status: false,
    createdAt: '2026-02-22 11:20:00',
    updatedAt: '2026-02-23 15:40:00'
  },
  {
    id: '4',
    name: '智谱元宝',
    provider: 'yuanbao',
    type: 'chat',
    version: 'v2.5',
    apiKey: 'yb-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx',
    apiUrl: 'https://api.yuanbao.com/v1/chat/completions',
    description: '智谱AI旗下的AI助手',
    extraParams: '{"max_tokens": 2048}',
    status: true,
    createdAt: '2026-02-23 14:10:00',
    updatedAt: '2026-02-25 09:30:00'
  }
])

// 计算属性
const filteredModels = computed(() => {
  if (!searchQuery.value) {
    return models.value
  }
  return models.value.filter(model => {
    const query = searchQuery.value.toLowerCase()
    return (
      model.name.toLowerCase().includes(query) ||
      model.type.toLowerCase().includes(query) ||
      model.provider.toLowerCase().includes(query)
    )
  })
})

const activeModelsCount = computed(() => {
  return models.value.filter(model => model.status).length
})

const inactiveModelsCount = computed(() => {
  return models.value.filter(model => !model.status).length
})

// 方法
const getModelTypeTagType = (type: string) => {
  const typeMap: Record<string, string> = {
    chat: 'primary',
    code: 'success',
    image: 'warning',
    audio: 'info',
    multimodal: 'danger'
  }
  return typeMap[type] || 'info'
}

const handleSelectionChange = (selection: any[]) => {
  selectedModels.value = selection
}

const handleBatchAction = (command: string) => {
  if (selectedModels.value.length === 0) {
    ElMessage.warning('请选择要操作的模型')
    return
  }

  switch (command) {
    case 'enable':
      selectedModels.value.forEach(model => {
        model.status = true
      })
      ElMessage.success(`已成功启用 ${selectedModels.value.length} 个模型`)
      break
    case 'disable':
      selectedModels.value.forEach(model => {
        model.status = false
      })
      ElMessage.success(`已成功禁用 ${selectedModels.value.length} 个模型`)
      break
  }
}

const handleModelStatusChange = (model: any) => {
  ElMessage.success(`模型 ${model.name} 已${model.status ? '启用' : '禁用'}`)
}

const showAddModelDialog = () => {
  resetForm()
  dialogTitle.value = '添加模型'
  isEditMode.value = false
  dialogVisible.value = true
}

const showEditModelDialog = (model: any) => {
  resetForm()
  dialogTitle.value = '编辑模型'
  isEditMode.value = true
  
  // 复制模型数据到表单
  Object.assign(formData, { ...model })
  dialogVisible.value = true
}

const resetForm = () => {
  Object.assign(formData, {
    id: '',
    name: '',
    provider: '',
    type: '',
    version: '',
    apiKey: '',
    apiUrl: '',
    description: '',
    extraParams: '',
    status: true
  })
  if (modelFormRef.value) {
    modelFormRef.value.resetFields()
  }
}

const submitForm = async () => {
  if (!modelFormRef.value) return
  
  try {
    await modelFormRef.value.validate()
    submitting.value = true
    
    // 模拟API请求
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    if (isEditMode.value) {
      // 更新现有模型
      const index = models.value.findIndex(m => m.id === formData.id)
      if (index !== -1) {
        models.value[index] = { ...formData, updatedAt: new Date().toLocaleString() }
      }
      ElMessage.success('模型更新成功')
    } else {
      // 添加新模型
      const newModel = {
        ...formData,
        id: Date.now().toString(),
        createdAt: new Date().toLocaleString(),
        updatedAt: new Date().toLocaleString()
      }
      models.value.unshift(newModel)
      ElMessage.success('模型添加成功')
    }
    
    dialogVisible.value = false
  } catch (error: any) {
    console.error('表单验证失败:', error)
  } finally {
    submitting.value = false
  }
}

const deleteModel = (model: any) => {
  ElMessageBox.confirm(
    `确定要删除模型 "${model.name}" 吗？`,
    '删除确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  )
    .then(() => {
      const index = models.value.findIndex(m => m.id === model.id)
      if (index !== -1) {
        models.value.splice(index, 1)
        ElMessage.success('模型删除成功')
      }
    })
    .catch(() => {
      // 取消删除
    })
}

const refreshModels = () => {
  loading.value = true
  // 模拟刷新数据
  setTimeout(() => {
    ElMessage.success('模型列表已刷新')
    loading.value = false
  }, 1000)
}

const testModelConnection = (model: any) => {
  selectedModelForTest.value = model
  testResult.value = null
  testDialogVisible.value = true
}

const startModelTest = async () => {
  testing.value = true
  testResult.value = null
  
  try {
    // 模拟测试请求
    await new Promise(resolve => setTimeout(resolve, 1500))
    
    // 模拟测试结果
    const success = Math.random() > 0.2
    testResult.value = {
      success,
      message: success 
        ? `成功连接到模型 ${selectedModelForTest.value.name}\n响应时间: ${Math.floor(Math.random() * 500) + 100}ms\n版本: ${selectedModelForTest.value.version}`
        : `无法连接到模型 ${selectedModelForTest.value.name}\n错误: API密钥无效或网络连接失败`
    }
  } catch (error: any) {
    testResult.value = {
      success: false,
      message: `测试失败: ${error.message}`
    }
  } finally {
    testing.value = false
  }
}
</script>

<style scoped>
.ai-model-management-view {
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

.header-stats {
  display: flex;
  gap: 24px;
}

.stat-item {
  text-align: center;
  background: rgba(255, 255, 255, 0.15);
  padding: 16px 24px;
  border-radius: 8px;
  backdrop-filter: blur(10px);
  min-width: 100px;
}

.stat-value {
  display: block;
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 4px;
}

.stat-label {
  display: block;
  font-size: 14px;
  opacity: 0.9;
}

/* 工具栏样式 */
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  max-width: 1200px;
  margin: 0 auto 24px auto;
  padding: 0 24px;
  gap: 24px;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.toolbar-right {
  flex: 0 0 auto;
}

.search-input {
  width: 300px;
}

/* 模型卡片样式 */
.models-card {
  max-width: 1200px;
  margin: 0 auto 30px auto;
  border-radius: 12px;
  background: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  border: 1px solid #e4e7ed;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
}

.models-card:hover {
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
}

.card-header .header-icon {
  font-size: 20px;
  color: #409eff;
  background: transparent;
  padding: 0;
  box-shadow: none;
}

.card-header span {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

/* 模型表格样式 */
.models-container {
  overflow-x: auto;
  padding: 0 24px;
}

.models-table {
  font-size: 14px;
}

.model-name-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.model-name-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.model-name {
  font-weight: 600;
  color: #303133;
}

.model-provider {
  font-size: 12px;
  color: #909399;
}

.operation-buttons {
  display: flex;
  gap: 8px;
}

/* 分页样式 */
.pagination-container {
  display: flex;
  justify-content: flex-end;
  padding: 20px 24px;
  border-top: 1px solid #f0f2f5;
}

/* 测试对话框样式 */
.test-dialog-content {
  min-height: 200px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.testing-indicator {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  text-align: center;
}

.loading-icon {
  font-size: 48px;
  color: #409eff;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.test-result {
  text-align: center;
  width: 100%;
}

.result-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.result-icon.success {
  color: #67c23a;
}

.result-icon.error {
  color: #f56c6c;
}

.test-result-details {
  margin-top: 16px;
  background: #f5f7fa;
  padding: 16px;
  border-radius: 8px;
  text-align: left;
  width: 100%;
  overflow-x: auto;
}

.test-result-details pre {
  margin: 0;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 13px;
  line-height: 1.5;
  color: #303133;
}

.test-instructions {
  text-align: center;
  color: #909399;
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .header-content {
    flex-direction: column;
    text-align: center;
    gap: 20px;
  }
  
  .toolbar {
    flex-direction: column;
    align-items: stretch;
    gap: 16px;
  }
  
  .toolbar-left,
  .toolbar-right {
    width: 100%;
  }
  
  .search-input {
    width: 100%;
  }
  
  .models-container {
    padding: 0 16px;
  }
}

@media (max-width: 768px) {
  .page-header {
    padding: 30px 0;
  }
  
  .header-text h2 {
    font-size: 24px;
  }
  
  .header-stats {
    flex-direction: column;
    gap: 12px;
    width: 100%;
  }
  
  .stat-item {
    width: 100%;
  }
  
  .toolbar {
    padding: 0 16px;
  }
  
  .toolbar-left {
    flex-wrap: wrap;
  }
  
  .models-card {
    margin: 0 16px 24px 16px;
  }
  
  .models-container {
    padding: 0 8px;
  }
  
  .pagination-container {
    padding: 16px;
  }
}
</style>