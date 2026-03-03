<template>
  <div class="media-view">
    <!-- 页面标题和操作 -->
    <div class="page-header">
      <div class="header-left">
        <h2 class="page-title">告警介质管理</h2>
        <p class="page-desc">配置和管理各种告警通知方式</p>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="showAddDialog = true" :icon="Plus">
          新增介质
        </el-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :span="4">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-info">
              <span class="stat-label">总介质数</span>
              <span class="stat-value">{{ mediaList.length }}</span>
            </div>
            <el-icon class="stat-icon"><Message /></el-icon>
          </div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-info">
              <span class="stat-label">启用中</span>
              <span class="stat-value">{{ enabledMediaCount }}</span>
            </div>
            <el-icon class="stat-icon success"><Check /></el-icon>
          </div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-info">
              <span class="stat-label">邮件介质</span>
              <span class="stat-value">{{ emailMediaCount }}</span>
            </div>
            <el-icon class="stat-icon email"><Message /></el-icon>
          </div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-info">
              <span class="stat-label">钉钉介质</span>
              <span class="stat-value">{{ dingtalkMediaCount }}</span>
            </div>
            <el-icon class="stat-icon dingtalk"><Message /></el-icon>
          </div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-info">
              <span class="stat-label">飞书介质</span>
              <span class="stat-value">{{ feishuMediaCount }}</span>
            </div>
            <el-icon class="stat-icon feishu"><Message /></el-icon>
          </div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-info">
              <span class="stat-label">企业微信</span>
              <span class="stat-value">{{ wechatMediaCount }}</span>
            </div>
            <el-icon class="stat-icon wechat"><Message /></el-icon>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 介质列表 -->
    <el-card shadow="hover" class="content-card">
      <div class="table-header">
        <div class="table-filters">
          <el-input
            v-model="searchQuery"
            placeholder="搜索介质名称"
            :prefix-icon="Search"
            style="width: 240px; margin-right: 16px"
          ></el-input>
          <el-select
            v-model="typeFilter"
            placeholder="筛选介质类型"
            style="width: 160px"
          >
            <el-option label="全部" value=""></el-option>
            <el-option label="邮件" value="email"></el-option>
            <el-option label="钉钉" value="dingtalk"></el-option>
            <el-option label="飞书" value="feishu"></el-option>
            <el-option label="企业微信" value="wechat"></el-option>
            <el-option label="Webhook" value="webhook"></el-option>
            <el-option label="卡片" value="card"></el-option>
          </el-select>
        </div>
      </div>
      <el-table :data="filteredMediaList" style="width: 100%" stripe>
        <el-table-column prop="name" label="介质名称" min-width="180"></el-table-column>
        <el-table-column prop="type" label="介质类型" width="120">
          <template #default="scope">
            <el-tag :type="getTypeTagType(scope.row.type)">
              {{ getTypeName(scope.row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="template.name" label="关联模板" min-width="180"></el-table-column>
        <el-table-column prop="isEnabled" label="状态" width="100">
          <template #default="scope">
            <el-switch v-model="scope.row.isEnabled" @change="toggleMediaStatus(scope.row)"></el-switch>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" @click="showEditDialog(scope.row)">
              <el-icon><Edit /></el-icon>
              编辑
            </el-button>
            <el-button type="success" size="small" @click="testMedia(scope.row)">
              <el-icon><Check /></el-icon>
              测试
            </el-button>
            <el-button type="danger" size="small" @click="deleteMedia(scope.row.id)">
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <div class="pagination-container">
        <el-pagination
          layout="total, sizes, prev, pager, next, jumper"
          :total="filteredMediaList.length"
          :page-size="pageSize"
          v-model:current-page="currentPage"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        ></el-pagination>
      </div>
    </el-card>

    <!-- 新增/编辑介质对话框 -->
    <el-dialog
      v-model="showDialog"
      :title="dialogTitle"
      width="700px"
      @close="resetForm"
    >
      <el-form :model="formData" label-width="120px" :rules="formRules" ref="formRef">
        <el-form-item label="介质名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入介质名称"></el-input>
        </el-form-item>
        <el-form-item label="介质类型" prop="type">
          <el-select v-model="formData.type" placeholder="请选择介质类型" @change="handleTypeChange">
            <el-option label="邮件" value="email"></el-option>
            <el-option label="钉钉" value="dingtalk"></el-option>
            <el-option label="飞书" value="feishu"></el-option>
            <el-option label="企业微信" value="wechat"></el-option>
            <el-option label="Webhook" value="webhook"></el-option>
            <el-option label="卡片" value="card"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="关联模板">
          <el-select v-model="formData.templateId" placeholder="请选择关联模板">
            <el-option
              v-for="template in templates"
              :key="template.id"
              :label="template.name"
              :value="template.id"
            ></el-option>
          </el-select>
        </el-form-item>
        
        <!-- 动态配置表单 -->
        <el-divider content-position="left">配置信息</el-divider>
        
        <!-- 邮件配置 -->
        <template v-if="formData.type === 'email'">
          <el-form-item label="SMTP服务器">
            <el-input v-model="emailConfig.smtp" placeholder="请输入SMTP服务器地址"></el-input>
          </el-form-item>
          <el-form-item label="SMTP端口">
            <el-input-number v-model="emailConfig.port" :min="1" :max="65535" placeholder="请输入SMTP端口"></el-input-number>
          </el-form-item>
          <el-form-item label="发件人邮箱">
            <el-input v-model="emailConfig.from" placeholder="请输入发件人邮箱"></el-input>
          </el-form-item>
          <el-form-item label="发件人名称">
            <el-input v-model="emailConfig.fromName" placeholder="请输入发件人名称"></el-input>
          </el-form-item>
          <el-form-item label="邮箱密码">
            <el-input v-model="emailConfig.password" type="password" placeholder="请输入邮箱密码或授权码"></el-input>
          </el-form-item>
          <el-form-item label="启用SSL">
            <el-switch v-model="emailConfig.ssl"></el-switch>
          </el-form-item>
        </template>
        
        <!-- 钉钉配置 -->
        <template v-else-if="formData.type === 'dingtalk'">
          <el-form-item label="Webhook地址">
            <el-input v-model="dingConfig.webhook" placeholder="请输入钉钉机器人Webhook地址"></el-input>
          </el-form-item>
          <el-form-item label="加签密钥">
            <el-input v-model="dingConfig.secret" placeholder="请输入钉钉机器人加签密钥"></el-input>
          </el-form-item>
          <el-form-item label="@所有人">
            <el-switch v-model="dingConfig.atAll"></el-switch>
          </el-form-item>
        </template>
        
        <!-- 飞书配置 -->
        <template v-else-if="formData.type === 'feishu'">
          <el-form-item label="Webhook地址">
            <el-input v-model="feishuConfig.webhook" placeholder="请输入飞书机器人Webhook地址"></el-input>
          </el-form-item>
          <el-form-item label="加签密钥">
            <el-input v-model="feishuConfig.secret" placeholder="请输入飞书机器人加签密钥"></el-input>
          </el-form-item>
        </template>
        
        <!-- 企业微信配置 -->
        <template v-else-if="formData.type === 'wechat'">
          <el-form-item label="Webhook地址">
            <el-input v-model="wechatConfig.webhook" placeholder="请输入企业微信机器人Webhook地址"></el-input>
          </el-form-item>
          <el-form-item label="@所有人">
            <el-switch v-model="wechatConfig.atAll"></el-switch>
          </el-form-item>
        </template>
        
        <!-- Webhook配置 -->
        <template v-else-if="formData.type === 'webhook'">
          <el-form-item label="Webhook地址">
            <el-input v-model="webhookConfig.url" placeholder="请输入Webhook地址"></el-input>
          </el-form-item>
          <el-form-item label="请求方法">
            <el-select v-model="webhookConfig.method">
              <el-option label="GET" value="GET"></el-option>
              <el-option label="POST" value="POST"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="请求头">
            <el-input
              v-model="webhookConfig.headers"
              type="textarea"
              :rows="3"
              placeholder="请输入请求头，JSON格式"
            ></el-input>
          </el-form-item>
        </template>
        
        <!-- 卡片配置 -->
        <template v-else-if="formData.type === 'card'">
          <el-form-item label="卡片标题">
            <el-input v-model="cardConfig.title" placeholder="请输入卡片标题"></el-input>
          </el-form-item>
          <el-form-item label="卡片颜色">
            <el-input v-model="cardConfig.color" placeholder="请输入卡片颜色"></el-input>
          </el-form-item>
        </template>
        
        <!-- 通用配置 -->
        <template v-else>
          <el-form-item label="配置信息">
            <el-input
              v-model="formData.config"
              type="textarea"
              :rows="6"
              placeholder="请输入配置信息，JSON格式"
            ></el-input>
          </el-form-item>
        </template>
        
        <el-form-item label="启用状态">
          <el-switch v-model="formData.isEnabled"></el-switch>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showDialog = false">取消</el-button>
          <el-button type="primary" @click="saveMedia">确定</el-button>
          <el-button type="success" @click="testCurrentMedia">测试连接</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { 
  Plus, 
  Message, 
  Check, 
  Edit, 
  Delete, 
  Search 
} from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'

// 模拟数据
const mediaList = ref([
  {
    id: 1,
    name: '邮件告警',
    type: 'email',
    config: JSON.stringify({ smtp: 'smtp.example.com', port: 587, from: 'alert@example.com', fromName: 'DeanTech' }),
    templateId: 1,
    template: { id: 1, name: '默认邮件模板' },
    isEnabled: true,
    createdAt: '2024-01-01T10:00:00Z'
  },
  {
    id: 2,
    name: '钉钉告警',
    type: 'dingtalk',
    config: JSON.stringify({ webhook: 'https://oapi.dingtalk.com/robot/send' }),
    templateId: 2,
    template: { id: 2, name: '钉钉模板' },
    isEnabled: true,
    createdAt: '2024-01-02T14:30:00Z'
  },
  {
    id: 3,
    name: '飞书告警',
    type: 'feishu',
    config: JSON.stringify({ webhook: 'https://open.feishu.cn/open-apis/bot/v2/hook/' }),
    templateId: 3,
    template: { id: 3, name: '飞书模板' },
    isEnabled: false,
    createdAt: '2024-01-03T09:15:00Z'
  },
  {
    id: 4,
    name: '企业微信告警',
    type: 'wechat',
    config: JSON.stringify({ webhook: 'https://qyapi.weixin.qq.com/cgi-bin/webhook/send' }),
    templateId: 4,
    template: { id: 4, name: '企业微信模板' },
    isEnabled: true,
    createdAt: '2024-01-04T16:45:00Z'
  }
])

const templates = ref([
  { id: 1, name: '默认邮件模板' },
  { id: 2, name: '钉钉模板' },
  { id: 3, name: '飞书模板' },
  { id: 4, name: '企业微信模板' },
  { id: 5, name: 'Webhook模板' }
])

// 表单数据
const showDialog = ref(false)
const showAddDialog = ref(false)
const dialogTitle = ref('新增介质')
const formRef = ref()
const formRules = {
  name: [{ required: true, message: '请输入介质名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择介质类型', trigger: 'change' }]
}

const formData = ref({
  id: 0,
  name: '',
  type: '',
  templateId: 0,
  config: '',
  isEnabled: true,
  createdAt: new Date().toISOString()
})

// 动态配置表单数据
const emailConfig = ref({
  smtp: '',
  port: 587,
  from: '',
  fromName: '',
  password: '',
  ssl: true
})

const dingConfig = ref({
  webhook: '',
  secret: '',
  atAll: false
})

const feishuConfig = ref({
  webhook: '',
  secret: ''
})

const wechatConfig = ref({
  webhook: '',
  atAll: false
})

const webhookConfig = ref({
  url: '',
  method: 'POST',
  headers: '{}'
})

const cardConfig = ref({
  title: '',
  color: '#409eff'
})

// 分页和筛选
const currentPage = ref(1)
const pageSize = ref(10)
const searchQuery = ref('')
const typeFilter = ref('')

// 计算属性
const enabledMediaCount = computed(() => {
  return mediaList.value.filter(item => item.isEnabled).length
})

const emailMediaCount = computed(() => {
  return mediaList.value.filter(item => item.type === 'email').length
})

const dingtalkMediaCount = computed(() => {
  return mediaList.value.filter(item => item.type === 'dingtalk').length
})

const feishuMediaCount = computed(() => {
  return mediaList.value.filter(item => item.type === 'feishu').length
})

const wechatMediaCount = computed(() => {
  return mediaList.value.filter(item => item.type === 'wechat').length
})

const filteredMediaList = computed(() => {
  return mediaList.value.filter(item => {
    const matchesSearch = item.name.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesType = typeFilter.value ? item.type === typeFilter.value : true
    return matchesSearch && matchesType
  })
})

// 方法
const getTypeTagType = (type: string) => {
  const typeMap: Record<string, string> = {
    'email': 'info',
    'dingtalk': 'success',
    'feishu': 'warning',
    'wechat': 'primary',
    'webhook': 'danger',
    'card': 'info'
  }
  return typeMap[type] || 'info'
}

const getTypeName = (type: string) => {
  const typeMap: Record<string, string> = {
    'email': '邮件',
    'dingtalk': '钉钉',
    'feishu': '飞书',
    'wechat': '企业微信',
    'webhook': 'Webhook',
    'card': '卡片'
  }
  return typeMap[type] || type
}

const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleString()
}

const handleTypeChange = () => {
  // 清空配置
  formData.value.config = ''
  // 重置所有配置表单
  Object.assign(emailConfig.value, { smtp: '', port: 587, from: '', fromName: '', password: '', ssl: true })
  Object.assign(dingConfig.value, { webhook: '', secret: '', atAll: false })
  Object.assign(feishuConfig.value, { webhook: '', secret: '' })
  Object.assign(wechatConfig.value, { webhook: '', atAll: false })
  Object.assign(webhookConfig.value, { url: '', method: 'POST', headers: '{}' })
  Object.assign(cardConfig.value, { title: '', color: '#409eff' })
}

// 显示编辑对话框
const showEditDialog = (row: any) => {
  formData.value = { ...row }
  dialogTitle.value = '编辑介质'
  showDialog.value = true
  
  // 解析配置到对应的表单
  const config = JSON.parse(row.config || '{}')
  if (row.type === 'email') {
    Object.assign(emailConfig.value, config)
  } else if (row.type === 'dingtalk') {
    Object.assign(dingConfig.value, config)
  } else if (row.type === 'feishu') {
    Object.assign(feishuConfig.value, config)
  } else if (row.type === 'wechat') {
    Object.assign(wechatConfig.value, config)
  } else if (row.type === 'webhook') {
    Object.assign(webhookConfig.value, config)
  } else if (row.type === 'card') {
    Object.assign(cardConfig.value, config)
  }
}

// 重置表单
const resetForm = () => {
  formData.value = {
    id: 0,
    name: '',
    type: '',
    templateId: 0,
    config: '',
    isEnabled: true,
    createdAt: new Date().toISOString()
  }
  
  // 重置所有配置表单
  Object.assign(emailConfig.value, { smtp: '', port: 587, from: '', fromName: '', password: '', ssl: true })
  Object.assign(dingConfig.value, { webhook: '', secret: '', atAll: false })
  Object.assign(feishuConfig.value, { webhook: '', secret: '' })
  Object.assign(wechatConfig.value, { webhook: '', atAll: false })
  Object.assign(webhookConfig.value, { url: '', method: 'POST', headers: '{}' })
  Object.assign(cardConfig.value, { title: '', color: '#409eff' })
}

// 保存介质
const saveMedia = () => {
  // 根据类型生成配置
  let config = ''
  if (formData.value.type === 'email') {
    config = JSON.stringify(emailConfig.value)
  } else if (formData.value.type === 'dingtalk') {
    config = JSON.stringify(dingConfig.value)
  } else if (formData.value.type === 'feishu') {
    config = JSON.stringify(feishuConfig.value)
  } else if (formData.value.type === 'wechat') {
    config = JSON.stringify(wechatConfig.value)
  } else if (formData.value.type === 'webhook') {
    config = JSON.stringify(webhookConfig.value)
  } else if (formData.value.type === 'card') {
    config = JSON.stringify(cardConfig.value)
  }
  
  formData.value.config = config
  
  if (formData.value.id) {
    // 编辑
    const index = mediaList.value.findIndex(item => item.id === formData.value.id)
    if (index !== -1) {
      const updatedMedia = {
        ...formData.value,
        template: templates.value.find(t => t.id === formData.value.templateId) || { id: 0, name: '' }
      }
      mediaList.value[index] = updatedMedia
    }
  } else {
    // 新增
    const newMedia = {
      ...formData.value,
      id: Date.now(),
      template: templates.value.find(t => t.id === formData.value.templateId) || { id: 0, name: '' }
    }
    mediaList.value.push(newMedia)
  }
  showDialog.value = false
  resetForm()
  ElMessage.success('保存成功')
}

// 切换介质状态
const toggleMediaStatus = (row: any) => {
  // 这里可以添加API调用，更新介质状态
  console.log('Toggle media status:', row.id, row.isEnabled)
  ElMessage.success(`介质${row.isEnabled ? '启用' : '禁用'}成功`)
}

// 删除介质
const deleteMedia = (id: number) => {
  ElMessageBox.confirm('确定要删除这个介质吗？', '删除确认', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    mediaList.value = mediaList.value.filter(item => item.id !== id)
    ElMessage.success('删除成功')
  }).catch(() => {
    // 取消删除
  })
}

// 测试介质
const testMedia = (row: any) => {
  ElMessage.info(`正在测试 ${row.name} 连接...`)
  // 这里可以添加API调用，测试介质连接
  setTimeout(() => {
    ElMessage.success(`${row.name} 测试成功`)
  }, 1000)
}

// 测试当前介质
const testCurrentMedia = () => {
  ElMessage.info('正在测试连接...')
  // 这里可以添加API调用，测试介质连接
  setTimeout(() => {
    ElMessage.success('测试成功')
  }, 1000)
}

// 分页处理
const handleSizeChange = (size: number) => {
  pageSize.value = size
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
}

// 监听新增对话框
const watchAddDialog = (newVal: boolean) => {
  if (newVal) {
    dialogTitle.value = '新增介质'
    resetForm()
    showDialog.value = true
  }
}

Object.defineProperty(showAddDialog, 'value', {
  get: () => showDialog.value && !formData.value.id,
  set: watchAddDialog
})
</script>

<style scoped>
.media-view {
  width: 100%;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
}

.page-title {
  margin: 0 0 8px 0;
  font-size: 24px;
  font-weight: 600;
  color: #303133;
}

.page-desc {
  margin: 0;
  font-size: 14px;
  color: #606266;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
}

.stat-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stat-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.stat-label {
  font-size: 14px;
  color: #606266;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
}

.stat-icon {
  font-size: 32px;
  color: #409eff;
}

.stat-icon.success {
  color: #67c23a;
}

.stat-icon.email {
  color: #409eff;
}

.stat-icon.dingtalk {
  color: #13c2c2;
}

.stat-icon.feishu {
  color: #722ed1;
}

.stat-icon.wechat {
  color: #67c23a;
}

.stat-icon.warning {
  color: #e6a23c;
}

.content-card {
  margin-bottom: 20px;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.table-filters {
  display: flex;
  gap: 12px;
  align-items: center;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>