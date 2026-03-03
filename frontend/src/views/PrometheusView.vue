<template>
  <div class="prometheus-view">
    <el-card shadow="never" class="page-header">
      <div class="header-content">
        <h2>Prometheus配置</h2>
        <div class="header-actions">
          <el-button type="primary" @click="showConfigDialog = true">
            <el-icon><Plus /></el-icon>
            新增Prometheus
          </el-button>
          <el-button type="info" @click="refreshData">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </div>
    </el-card>

    <el-card shadow="hover" class="config-card">
      <template #header>
        <div class="card-header">
          <span>Prometheus实例管理</span>
        </div>
      </template>
      <el-table :data="prometheusConfigs" style="width: 100%" stripe border>
        <el-table-column prop="name" label="名称" min-width="150" />
        <el-table-column prop="address" label="地址" min-width="200" />
        <el-table-column prop="alertManagerURL" label="AlertManager地址" min-width="200" />
        <el-table-column prop="timeout" label="超时时间" width="120">
          <template #default="scope">
            <el-tag type="info" size="small">{{ scope.row.timeout }}s</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'active' ? 'success' : 'danger'" size="small">
              {{ scope.row.status === 'active' ? '活跃' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="scope">
            <el-button size="small" type="primary" @click="editConfig(scope.row)">
              <el-icon><Edit /></el-icon>
              编辑
            </el-button>
            <el-button size="small" type="danger" @click="deleteConfig(scope.row.id)">
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
            <el-switch 
              v-model="scope.row.status" 
              :active-value="'active'" 
              :inactive-value="'inactive'"
              @change="updateStatus(scope.row)"
              active-text="活跃"
              inactive-text="禁用"
            />
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Prometheus配置对话框 -->
    <el-dialog v-model="showConfigDialog" :title="isEditing ? '编辑Prometheus配置' : '新增Prometheus配置'" width="500px">
      <el-form :model="configForm" :rules="configRules" ref="configFormRef" label-width="100px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="configForm.name" placeholder="请输入Prometheus名称"></el-input>
        </el-form-item>
        <el-form-item label="地址" prop="address">
          <el-input v-model="configForm.address" placeholder="请输入Prometheus地址，如：http://localhost:9090"></el-input>
        </el-form-item>
        <el-form-item label="AlertManager地址" prop="alertManagerURL">
          <el-input v-model="configForm.alertManagerURL" placeholder="请输入AlertManager地址，如：http://localhost:9093"></el-input>
        </el-form-item>
        <el-form-item label="超时时间">
          <el-input-number v-model="configForm.timeout" :min="5" :max="60" :step="5" placeholder="请输入超时时间"></el-input-number>
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="configForm.status" active-value="active" inactive-value="inactive"></el-switch>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="showConfigDialog = false">取消</el-button>
          <el-button type="primary" @click="saveConfig">保存</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { 
  Refresh, Plus, Edit, Delete
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const showConfigDialog = ref(false)
const isEditing = ref(false)
const configFormRef = ref()

// Prometheus配置表单
const configForm = ref({
  id: '',
  name: '',
  address: 'http://localhost:9090',
  alertManagerURL: 'http://localhost:9093',
  timeout: 10,
  status: 'active'
})

// 表单验证规则
const configRules = {
  name: [{ required: true, message: '请输入Prometheus名称', trigger: 'blur' }],
  address: [{ required: true, message: '请输入Prometheus地址', trigger: 'blur' }],
  alertManagerURL: [{ required: true, message: '请输入AlertManager地址', trigger: 'blur' }]
}

// 模拟Prometheus配置数据
const prometheusConfigs = ref([
  {
    id: '1',
    name: '默认Prometheus',
    address: 'http://localhost:9090',
    alertManagerURL: 'http://localhost:9093',
    timeout: 10,
    status: 'active'
  }
])

// 刷新数据
const refreshData = () => {
  // 这里可以添加API调用，从Prometheus获取最新数据
  ElMessage.success('数据已刷新')
  console.log('Refreshing Prometheus data...')
}

// 打开编辑对话框
const editConfig = (row: any) => {
  isEditing.value = true
  configForm.value = { ...row }
  showConfigDialog.value = true
}

// 删除配置
const deleteConfig = (id: string) => {
  prometheusConfigs.value = prometheusConfigs.value.filter(config => config.id !== id)
  ElMessage.success('Prometheus配置已删除')
  console.log('删除Prometheus配置:', id)
}

// 更新状态
const updateStatus = (row: any) => {
  ElMessage.success(`Prometheus配置${row.status === 'active' ? '已激活' : '已禁用'}`)
  console.log('更新Prometheus配置状态:', row)
}

// 保存配置
const saveConfig = () => {
  if (!configFormRef.value) return
  
  configFormRef.value.validate((valid: boolean) => {
    if (valid) {
      if (isEditing.value) {
        // 编辑现有配置
        const index = prometheusConfigs.value.findIndex(config => config.id === configForm.value.id)
        if (index !== -1) {
          prometheusConfigs.value[index] = { ...configForm.value }
        }
        ElMessage.success('Prometheus配置已更新')
      } else {
        // 新增配置
        const newConfig = {
          ...configForm.value,
          id: Date.now().toString(),
        }
        prometheusConfigs.value.push(newConfig)
        ElMessage.success('Prometheus配置已添加')
      }
      showConfigDialog.value = false
      resetForm()
      console.log('保存Prometheus配置:', configForm.value)
    }
  })
}

// 重置表单
const resetForm = () => {
  if (configFormRef.value) {
    configFormRef.value.resetFields()
  }
  isEditing.value = false
  configForm.value = {
    id: '',
    name: '',
    address: 'http://localhost:9090',
    alertManagerURL: 'http://localhost:9093',
    timeout: 10,
    status: 'active'
  }
}
</script>

<style scoped>
.prometheus-view {
  width: 100%;
  padding: 0;
  background-color: #f5f7fa;
  min-height: 100vh;
}

.page-header {
  margin-bottom: 20px;
  background-color: #ffffff;
  border-radius: 0;
  padding: 16px 24px;
  border-bottom: 1px solid #ebeef5;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.08);
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
  color: #303133;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.config-card {
  margin: 20px;
  border-radius: 8px;
  overflow: hidden;
  background-color: #ffffff;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;
}

.config-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

:deep(.el-table) {
  border-radius: 8px;
  overflow: hidden;
}

:deep(.el-table__header-wrapper th) {
  background-color: #f7f8fa;
  font-weight: 600;
  color: #303133;
}

:deep(.el-table__body-wrapper tr:hover) {
  background-color: #f5f7fa;
}

:deep(.el-button) {
  margin-right: 8px;
}

:deep(.el-switch) {
  margin-left: 8px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px;
  background-color: #fafafa;
  border-top: 1px solid #ebeef5;
}

:deep(.el-dialog__header) {
  background-color: #f7f8fa;
  border-bottom: 1px solid #ebeef5;
  padding: 20px;
}

:deep(.el-dialog__title) {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

:deep(.el-dialog__body) {
  padding: 24px;
}

:deep(.el-form-item) {
  margin-bottom: 20px;
}

:deep(.el-form-item__label) {
  font-weight: 500;
  color: #303133;
}

:deep(.el-button--primary) {
  background-color: #409eff;
  border-color: #409eff;
}

:deep(.el-button--primary:hover) {
  background-color: #66b1ff;
  border-color: #66b1ff;
}

:deep(.el-button--danger) {
  background-color: #f56c6c;
  border-color: #f56c6c;
}

:deep(.el-button--danger:hover) {
  background-color: #f78989;
  border-color: #f78989;
}
</style>