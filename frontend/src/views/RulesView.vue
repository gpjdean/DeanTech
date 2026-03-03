<template>
  <div class="rules-view">
    <el-card shadow="never" class="page-header">
      <div class="header-content">
        <h2>告警规则</h2>
      </div>
    </el-card>

    <div class="rules-section">
      <div class="rules-header">
        <h3>告警规则管理</h3>
        <el-button type="primary" @click="openAddRuleDialog">
          <el-icon><Plus /></el-icon>
          添加告警规则
        </el-button>
      </div>

      <el-card shadow="hover" class="rules-card">
        <el-table :data="rules" style="width: 100%" stripe>
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="name" label="规则名称" min-width="150" />
          <el-table-column prop="expr" label="表达式" min-width="200" />
          <el-table-column prop="severity" label="级别" width="100">
            <template #default="scope">
              <el-tag :type="getSeverityType(scope.row.severity)" size="small">
                {{ getSeverityName(scope.row.severity) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="for" label="持续时间" width="120" />
          <el-table-column prop="labels" label="标签" width="150">
            <template #default="scope">
              <el-tag v-for="(value, key) in scope.row.labels" :key="key" size="small" class="label-tag">
                {{ key }}: {{ value }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="annotations" label="注释" width="150">
            <template #default="scope">
              <el-tag v-for="(value, key) in scope.row.annotations" :key="key" size="small" class="annotation-tag">
                {{ key }}: {{ value }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100">
            <template #default="scope">
              <el-switch v-model="scope.row.status" active-text="启用" inactive-text="禁用" />
            </template>
          </el-table-column>
          <el-table-column label="操作" width="150" fixed="right">
            <template #default="scope">
              <el-button type="primary" size="small" @click="editRule(scope.row)">
                <el-icon><Edit /></el-icon>
                编辑
              </el-button>
              <el-button type="danger" size="small" @click="deleteRule(scope.row.id)">
                <el-icon><Delete /></el-icon>
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </div>

    <!-- 添加/编辑告警规则对话框 -->
    <el-dialog
      v-model="ruleDialogVisible"
      :title="isEditRule ? '编辑告警规则' : '添加告警规则'"
      width="700px"
    >
      <el-form ref="ruleForm" :model="ruleFormData" label-width="100px">
        <el-form-item label="规则名称" prop="name">
          <el-input v-model="ruleFormData.name" placeholder="请输入规则名称" /></el-form-item>

        <el-form-item label="表达式" prop="expr">
          <el-input v-model="ruleFormData.expr" placeholder="请输入PromQL表达式" type="textarea" :rows="3" /></el-form-item>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="告警级别" prop="severity">
              <el-select v-model="ruleFormData.severity" placeholder="请选择告警级别">
                <el-option label="紧急" value="critical" />
                <el-option label="警告" value="warning" />
                <el-option label="信息" value="info" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="持续时间" prop="for">
              <el-input v-model="ruleFormData.for" placeholder="请输入持续时间，如5m" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="状态" prop="status">
          <el-switch v-model="ruleFormData.status" active-text="启用" inactive-text="禁用" />
        </el-form-item>

        <el-form-item label="标签">
          <el-input v-model="ruleFormData.labels" placeholder="请输入标签，格式：{key: 'value'}" type="textarea" :rows="2" />
        </el-form-item>

        <el-form-item label="注释">
          <el-input v-model="ruleFormData.annotations" placeholder="请输入注释，格式：{key: 'value'}" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="ruleDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveRule" :loading="loading.saveRule">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus, Edit, Delete } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'

// 表单引用
const ruleForm = ref()

// 加载状态
const loading = ref({
  rules: false,
  saveRule: false
})

// 监控规则数据
interface Rule {
  id: number;
  name: string;
  expr: string;
  severity: string;
  for: string;
  labels: Record<string, string>;
  annotations: Record<string, string>;
  status: boolean;
}

const rules = ref<Rule[]>([])

// 监控规则对话框相关
const ruleDialogVisible = ref(false)
const isEditRule = ref(false)
const ruleFormData = ref({
  id: 0,
  name: '',
  expr: '',
  severity: 'warning',
  for: '5m',
  labels: {},
  annotations: {},
  status: true
})

// 组件挂载时加载数据
onMounted(() => {
  loadRules()
})

// 加载监控规则列表
const loadRules = async () => {
  try {
    loading.value.rules = true
    // 这里应该调用真实的API获取监控规则
    // 目前使用模拟数据
    // const response = await ruleAPI.list()
    // rules.value = response.data
    
    // 模拟数据
    rules.value = [
      {
        id: 1,
        name: 'High CPU Usage',
        expr: 'sum(rate(node_cpu_seconds_total{mode="idle"}[5m])) by (instance) < 0.5',
        severity: 'critical',
        for: '5m',
        labels: { alertname: 'HighCPUUsage', service: 'node' },
        annotations: { description: 'CPU usage is too high', summary: 'High CPU Usage Alert' },
        status: true
      },
      {
        id: 2,
        name: 'Low Disk Space',
        expr: 'node_filesystem_avail_bytes{fstype="ext4"} / node_filesystem_size_bytes{fstype="ext4"} < 0.1',
        severity: 'warning',
        for: '10m',
        labels: { alertname: 'LowDiskSpace', service: 'node' },
        annotations: { description: 'Disk space is running low', summary: 'Low Disk Space Alert' },
        status: true
      }
    ]
    ElMessage.success('加载告警规则成功')
  } catch (error) {
    console.error('加载告警规则失败:', error)
    ElMessage.error('加载告警规则失败')
  } finally {
    loading.value.rules = false
  }
}

// 打开添加监控规则对话框
const openAddRuleDialog = () => {
  isEditRule.value = false
  ruleFormData.value = {
    id: 0,
    name: '',
    expr: '',
    severity: 'warning',
    for: '5m',
    labels: {},
    annotations: {},
    status: true
  }
  ruleDialogVisible.value = true
}

// 编辑监控规则
const editRule = (rule: any) => {
  isEditRule.value = true
  ruleFormData.value = { ...rule }
  ruleDialogVisible.value = true
}

// 删除监控规则
const deleteRule = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除该监控规则吗？', '删除确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    // 这里应该调用真实的API删除监控规则
    // await ruleAPI.delete(id)
    
    // 模拟删除
    rules.value = rules.value.filter(rule => rule.id !== id)
    ElMessage.success('监控规则删除成功')
  } catch (error) {
    // 取消删除或删除失败
    if (error !== 'cancel') {
      console.error('删除监控规则失败:', error)
      ElMessage.error('删除监控规则失败')
    }
  }
}

// 保存监控规则
const saveRule = async () => {
  try {
    loading.value.saveRule = true
    
    if (isEditRule.value) {
      // 编辑监控规则
      // await ruleAPI.update(ruleFormData.value.id, ruleFormData.value)
      
      // 模拟编辑
      const index = rules.value.findIndex(rule => rule.id === ruleFormData.value.id)
      if (index !== -1) {
        rules.value[index] = { ...ruleFormData.value }
      }
      ElMessage.success('监控规则编辑成功')
    } else {
      // 添加监控规则
      // await ruleAPI.create(ruleFormData.value)
      
      // 模拟添加
      const newRule = {
        ...ruleFormData.value,
        id: rules.value.length + 1
      }
      rules.value.push(newRule)
      ElMessage.success('监控规则添加成功')
    }
    
    ruleDialogVisible.value = false
  } catch (error) {
    console.error('保存监控规则失败:', error)
    ElMessage.error('保存监控规则失败')
  } finally {
    loading.value.saveRule = false
  }
}

// 获取告警级别类型
const getSeverityType = (severity: string) => {
  switch (severity) {
    case 'critical':
      return 'danger'
    case 'warning':
      return 'warning'
    case 'info':
      return 'info'
    default:
      return 'success'
  }
}

// 获取告警级别名称
const getSeverityName = (severity: string) => {
  switch (severity) {
    case 'critical':
      return '紧急'
    case 'warning':
      return '警告'
    case 'info':
      return '信息'
    default:
      return severity
  }
}
</script>

<style scoped>
.rules-view {
  width: 100%;
  padding: 0;
}

.page-header {
  margin-bottom: 20px;
}

.header-content h2 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: #333;
}

.rules-section {
  margin-bottom: 20px;
}

.rules-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.rules-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.rules-card {
  margin-bottom: 20px;
}

.label-tag {
  margin-right: 8px;
  margin-bottom: 8px;
  background-color: #ecf5ff;
  border-color: #d9ecff;
  color: #409eff;
}

.annotation-tag {
  margin-right: 8px;
  margin-bottom: 8px;
  background-color: #f0f9eb;
  border-color: #e1f3d8;
  color: #67c23a;
}
</style>