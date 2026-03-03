<template>
  <div class="sls-config-view">
    <el-card shadow="never" class="page-header">
      <div class="header-content">
        <h2>SLS配置管理</h2>
        <div class="header-actions">
          <el-button type="primary" @click="handleAddConfig">
            <el-icon><Plus /></el-icon>
            新增配置
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- 配置列表 -->
    <el-card shadow="hover" class="sls-config-card">
      <template #header>
        <div class="card-header">
          <span>配置列表</span>
        </div>
      </template>
      
      <div class="config-list-container">
        <el-table 
          :data="configs" 
          style="width: 100%" 
          v-loading="loading"
          @row-dblclick="handleEditConfig"
        >
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="accessKeyId" label="AccessKey ID" min-width="200" />
          <el-table-column prop="regionId" label="地域" width="150" />
          <el-table-column prop="defaultProject" label="日志仓库" width="180" />
          <el-table-column prop="defaultLogstore" label="日志主题" width="180" />
          <el-table-column prop="timeout" label="超时时间" width="120" />
          <el-table-column prop="isActive" label="是否活跃" width="120">
            <template #default="scope">
              <el-switch 
                v-model="scope.row.isActive" 
                active-text="活跃" 
                inactive-text="非活跃"
                @change="handleActiveChange(scope.row)"
              />
            </template>
          </el-table-column>
          <el-table-column prop="createdAt" label="创建时间" width="200" />
          <el-table-column prop="updatedAt" label="更新时间" width="200" />
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="scope">
              <div class="table-actions">
                <el-button 
                  type="primary" 
                  size="small" 
                  @click="handleEditConfig(scope.row)"
                >
                  <el-icon><Edit /></el-icon>
                  编辑
                </el-button>
                <el-button 
                  type="danger" 
                  size="small" 
                  @click="handleDeleteConfig(scope.row.id)"
                >
                  <el-icon><Delete /></el-icon>
                  删除
                </el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>

    <!-- 配置表单对话框 -->
    <el-dialog 
      v-model="dialogVisible" 
      :title="dialogTitle" 
      width="600px" 
      center
    >
      <el-form 
        ref="configFormRef" 
        :model="configForm" 
        label-width="140px" 
        class="config-form"
        :rules="configRules"
      >
        <el-form-item label="AccessKey ID" prop="accessKeyId">
          <el-input 
            v-model="configForm.accessKeyId" 
            placeholder="请输入阿里云AccessKey ID"
            style="width: 100%"
          />
        </el-form-item>
        
        <el-form-item label="AccessKey Secret" prop="accessKeySecret">
          <el-input 
            v-model="configForm.accessKeySecret" 
            type="password" 
            placeholder="请输入阿里云AccessKey Secret"
            show-password
            style="width: 100%"
          />
        </el-form-item>
        
        <el-form-item label="地域ID" prop="regionId">
          <el-select 
            v-model="configForm.regionId" 
            placeholder="选择或输入地域ID"
            filterable
            allow-create
            default-first-option
            style="width: 100%"
          >
            <el-option
              v-for="region in regions"
              :key="region.value"
              :label="region.label"
              :value="region.value"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="默认日志仓库" prop="defaultProject">
          <el-input 
            v-model="configForm.defaultProject" 
            placeholder="请输入默认日志仓库名称"
            style="width: 100%"
          />
        </el-form-item>
        
        <el-form-item label="默认日志主题" prop="defaultLogstore">
          <el-input 
            v-model="configForm.defaultLogstore" 
            placeholder="请输入默认日志主题名称"
            style="width: 100%"
          />
        </el-form-item>
        
        <el-form-item label="超时时间(秒)" prop="timeout">
          <el-input-number 
            v-model="configForm.timeout" 
            :min="1" 
            :max="60" 
            :step="1" 
            placeholder="请输入超时时间"
            style="width: 100%"
          />
        </el-form-item>
        
        <el-form-item label="是否活跃" prop="isActive">
          <el-switch v-model="configForm.isActive" active-text="是" inactive-text="否" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSaveConfig" :loading="formLoading">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, Edit, Delete } from '@element-plus/icons-vue'
import axios from 'axios'

// SLS配置类型定义
interface SLSConfig {
  id: number | null;
  accessKeyId: string;
  accessKeySecret: string;
  regionId: string;
  defaultProject: string;
  defaultLogstore: string;
  timeout: number;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

// 表单引用
const configFormRef = ref()

// 加载状态
const loading = ref(false)
const formLoading = ref(false)

// 对话框状态
const dialogVisible = ref(false)
const dialogTitle = ref('新增SLS配置')

// 配置列表
const configs = ref<SLSConfig[]>([])

// 当前编辑的配置ID
const currentConfigId = ref<number | null>(null)

// 配置表单
const configForm = ref<SLSConfig>({
  id: null,
  accessKeyId: '',
  accessKeySecret: '',
  regionId: 'cn-hangzhou',
  defaultProject: '',
  defaultLogstore: '',
  timeout: 30,
  isActive: false,
  createdAt: '',
  updatedAt: ''
})

// 地域列表
const regions = ref([
  { label: '华东1(杭州)', value: 'cn-hangzhou' },
  { label: '华东2(上海)', value: 'cn-shanghai' },
  { label: '华南1(深圳)', value: 'cn-shenzhen' },
  { label: '华北1(青岛)', value: 'cn-qingdao' },
  { label: '华北2(北京)', value: 'cn-beijing' },
  { label: '华北3(张家口)', value: 'cn-zhangjiakou' },
  { label: '华北5(呼和浩特)', value: 'cn-huhehaote' },
  { label: '西南1(成都)', value: 'cn-chengdu' },
  { label: '中国香港', value: 'cn-hongkong' },
  { label: '新加坡', value: 'ap-southeast-1' },
  { label: '美国西部1(硅谷)', value: 'us-west-1' },
  { label: '美国东部1(弗吉尼亚)', value: 'us-east-1' },
  { label: '欧洲中部1(法兰克福)', value: 'eu-central-1' },
  { label: '中东东部1(迪拜)', value: 'me-east-1' }
])

// 表单验证规则
const configRules = {
  accessKeyId: [
    { required: true, message: '请输入AccessKey ID', trigger: 'blur' }
  ],
  accessKeySecret: [
    { required: true, message: '请输入AccessKey Secret', trigger: 'blur' }
  ],
  regionId: [
    { required: true, message: '请选择地域ID', trigger: 'change' }
  ]
}

// 组件挂载时加载配置列表
onMounted(() => {
  loadConfigs()
})

// 加载配置列表
const loadConfigs = async () => {
  try {
    loading.value = true
    // 从后端API加载配置列表，axios会自动添加baseURL
    const response = await axios.get('/sls/configs')
    configs.value = response.data
    loading.value = false
  } catch (error) {
    console.error('加载SLS配置列表失败:', error)
    loading.value = false
    ElMessage.error('加载SLS配置列表失败')
  }
}

// 打开新增配置对话框
const handleAddConfig = () => {
  dialogTitle.value = '新增SLS配置'
  currentConfigId.value = null
  resetConfigForm()
  dialogVisible.value = true
}

// 打开编辑配置对话框
const handleEditConfig = (row: SLSConfig) => {
  dialogTitle.value = '编辑SLS配置'
  currentConfigId.value = row.id
  // 深拷贝对象，避免直接修改表格数据
  configForm.value = { ...row }
  dialogVisible.value = true
}

// 保存配置
const handleSaveConfig = async () => {
  if (!configFormRef.value) return
  
  configFormRef.value.validate(async (valid: boolean) => {
    if (valid) {
      try {
        formLoading.value = true
        
        let response
        if (currentConfigId.value) {
          // 更新配置
          response = await axios.put(`/sls/configs/${currentConfigId.value}`, configForm.value)
        } else {
          // 创建配置
          response = await axios.post('/sls/configs', configForm.value)
        }
        
        formLoading.value = false
        dialogVisible.value = false
        ElMessage.success(response.data.message)
        // 重新加载配置列表
        loadConfigs()
      } catch (error) {
        console.error('保存SLS配置失败:', error)
        formLoading.value = false
        ElMessage.error('保存SLS配置失败')
      }
    }
  })
}

// 删除配置
const handleDeleteConfig = async (id: number) => {
  try {
    await axios.delete(`/sls/configs/${id}`)
    ElMessage.success('删除SLS配置成功')
    // 重新加载配置列表
    loadConfigs()
  } catch (error) {
    console.error('删除SLS配置失败:', error)
    ElMessage.error('删除SLS配置失败')
  }
}

// 处理活跃状态变更
const handleActiveChange = async (row: SLSConfig) => {
  try {
    await axios.put(`/sls/configs/${row.id}/active`)
    ElMessage.success('活跃状态更新成功')
    // 重新加载配置列表，确保只有一个活跃配置
    loadConfigs()
  } catch (error) {
    console.error('更新活跃状态失败:', error)
    ElMessage.error('更新活跃状态失败')
    // 恢复原状态
    loadConfigs()
  }
}

// 重置配置表单
const resetConfigForm = () => {
  configForm.value = {
    id: null,
    accessKeyId: '',
    accessKeySecret: '',
    regionId: 'cn-hangzhou',
    defaultProject: '',
    defaultLogstore: '',
    timeout: 30,
    isActive: false,
    createdAt: '',
    updatedAt: ''
  }
}
</script>

<style scoped>
.sls-config-view {
  width: 100%;
  padding: 0;
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

.sls-config-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.config-list-container {
  overflow-x: auto;
}

.config-form {
  margin: 0;
}

.table-actions {
  display: flex;
  gap: 8px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>