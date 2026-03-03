<template>
  <div class="templates-view">
    <el-card shadow="never" class="page-header">
      <div class="header-content">
        <h2>告警模板</h2>
        <div class="header-actions">
          <el-button type="primary" @click="showAddDialog = true">
            <el-icon><Plus /></el-icon>
            新增模板
          </el-button>
        </div>
      </div>
    </el-card>

    <el-card shadow="never" class="content-card">
      <el-table :data="templates" style="width: 100%">
        <el-table-column prop="name" label="模板名称" min-width="150"></el-table-column>
        <el-table-column prop="type" label="模板类型" width="120"></el-table-column>
        <el-table-column prop="isDefault" label="默认模板" width="100">
          <template #default="scope">
            <el-tag v-if="scope.row.isDefault" type="success">是</el-tag>
            <el-tag v-else type="info">否</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" @click="showEditDialog(scope.row)">
              编辑
            </el-button>
            <el-button type="danger" size="small" @click="deleteTemplate(scope.row.id)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新增/编辑模板对话框 -->
    <el-dialog
      v-model="showDialog"
      :title="dialogTitle"
      width="700px"
      @close="resetForm"
    >
      <el-form :model="formData" label-width="120px">
        <el-form-item label="模板名称" required>
          <el-input v-model="formData.name" placeholder="请输入模板名称"></el-input>
        </el-form-item>
        <el-form-item label="模板类型" required>
          <el-select v-model="formData.type" placeholder="请选择模板类型">
            <el-option label="邮件" value="email"></el-option>
            <el-option label="短信" value="sms"></el-option>
            <el-option label="Webhook" value="webhook"></el-option>
            <el-option label="卡片" value="card"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="模板内容" required>
          <el-input
            v-model="formData.content"
            type="textarea"
            :rows="10"
            placeholder="请输入模板内容"
          ></el-input>
        </el-form-item>
        <el-form-item label="设为默认模板">
          <el-switch v-model="formData.isDefault"></el-switch>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showDialog = false">取消</el-button>
          <el-button type="primary" @click="saveTemplate">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Plus } from '@element-plus/icons-vue'

// 模拟数据
const templates = ref([
  {
    id: 1,
    name: '默认邮件模板',
    content: `告警名称: {{ alertName }}\n状态: {{ status }}\n级别: {{ severity }}\n开始时间: {{ startsAt }}\n详情: {{ annotations.description }}`,
    type: 'email',
    isDefault: true
  },
  {
    id: 2,
    name: 'Webhook模板',
    content: `{"alert":"{{ alertName }}","status":"{{ status }}","severity":"{{ severity }}"}`,
    type: 'webhook',
    isDefault: false
  },
  {
    id: 3,
    name: '短信模板',
    content: `【DeanTech】{{ alertName }}告警触发，级别：{{ severity }}，请及时处理！`,
    type: 'sms',
    isDefault: false
  }
])

const showDialog = ref(false)
const showAddDialog = ref(false)
const dialogTitle = ref('新增模板')
const formData = ref({
  id: 0,
  name: '',
  content: '',
  type: '',
  isDefault: false
})

// 显示编辑对话框
const showEditDialog = (row: any) => {
  formData.value = { ...row }
  dialogTitle.value = '编辑模板'
  showDialog.value = true
}

// 重置表单
const resetForm = () => {
  formData.value = {
    id: 0,
    name: '',
    content: '',
    type: '',
    isDefault: false
  }
}

// 保存模板
const saveTemplate = () => {
  if (formData.value.isDefault) {
    // 如果设为默认模板，将其他模板的默认状态取消
    templates.value.forEach(t => {
      t.isDefault = false
    })
  }

  if (formData.value.id) {
    // 编辑
    const index = templates.value.findIndex(item => item.id === formData.value.id)
    if (index !== -1) {
      templates.value[index] = { ...formData.value }
    }
  } else {
    // 新增
    const newTemplate = {
      ...formData.value,
      id: Date.now()
    }
    templates.value.push(newTemplate)
  }
  showDialog.value = false
  resetForm()
}

// 删除模板
const deleteTemplate = (id: number) => {
  templates.value = templates.value.filter(item => item.id !== id)
}

// 监听新增对话框
const watchAddDialog = (newVal: boolean) => {
  if (newVal) {
    dialogTitle.value = '新增模板'
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
.templates-view {
  width: 100%;
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
  font-size: 20px;
}

.content-card {
  margin-bottom: 20px;
}
</style>