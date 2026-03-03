<template>
  <div class="silences-view">
    <el-card shadow="never" class="page-header">
      <div class="header-content">
        <h2>静默规则</h2>
        <div class="header-actions">
          <el-button type="primary" @click="showAddDialog = true">
            <el-icon><Plus /></el-icon>
            新增静默规则
          </el-button>
        </div>
      </div>
    </el-card>

    <el-card shadow="never" class="content-card">
      <el-table :data="silences" style="width: 100%">
        <el-table-column prop="matchers" label="匹配规则" min-width="200">
          <template #default="scope">
            <el-popover
              placement="top"
              title="匹配规则"
              width="300"
              trigger="hover"
            >
              <pre>{{ JSON.parse(scope.row.matchers) }}</pre>
              <template #reference>
                <el-button type="text" size="small">查看</el-button>
              </template>
            </el-popover>
          </template>
        </el-table-column>
        <el-table-column prop="startsAt" label="开始时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.startsAt) }}
          </template>
        </el-table-column>
        <el-table-column prop="endsAt" label="结束时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.endsAt) }}
          </template>
        </el-table-column>
        <el-table-column prop="createdBy" label="创建人" width="120"></el-table-column>
        <el-table-column prop="comment" label="备注" min-width="150"></el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'active' ? 'success' : 'info'">
              {{ scope.row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" @click="showEditDialog(scope.row)">
              <el-icon><Edit /></el-icon>
              编辑
            </el-button>
            <el-button type="danger" size="small" @click="deleteSilence(scope.row.id)">
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </template>
        </el-table-column>

      </el-table>
    </el-card>

    <!-- 新增/编辑静默规则对话框 -->
    <el-dialog
      v-model="showDialog"
      :title="dialogTitle"
      width="600px"
      @close="resetForm"
    >
      <el-form :model="formData" label-width="120px">
        <el-form-item label="匹配规则" required>
          <el-input
            v-model="formData.matchers"
            type="textarea"
            :rows="4"
            placeholder='请输入匹配规则，JSON格式，例如：[{"name":"job","value":"node_exporter","isRegex":false}]'
          ></el-input>
        </el-form-item>
        <el-form-item label="开始时间" required>
          <el-date-picker
            v-model="formData.startsAt"
            type="datetime"
            placeholder="选择开始时间"
            style="width: 100%"
          ></el-date-picker>
        </el-form-item>
        <el-form-item label="结束时间" required>
          <el-date-picker
            v-model="formData.endsAt"
            type="datetime"
            placeholder="选择结束时间"
            style="width: 100%"
          ></el-date-picker>
        </el-form-item>
        <el-form-item label="创建人" required>
          <el-input v-model="formData.createdBy" placeholder="请输入创建人"></el-input>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="formData.comment" placeholder="请输入备注"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showDialog = false">取消</el-button>
          <el-button type="primary" @click="saveSilence">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Plus } from '@element-plus/icons-vue'

// 模拟数据
const silences = ref([
  {
    id: 1,
    matchers: JSON.stringify([{ name: 'job', value: 'node_exporter', isRegex: false }]),
    startsAt: '2024-01-01T10:00:00Z',
    endsAt: '2024-01-02T10:00:00Z',
    createdBy: 'admin',
    comment: '维护期间静默',
    status: 'active'
  },
  {
    id: 2,
    matchers: JSON.stringify([{ name: 'instance', value: 'server1', isRegex: false }]),
    startsAt: '2024-01-01T08:00:00Z',
    endsAt: '2024-01-01T09:00:00Z',
    createdBy: 'admin',
    comment: '临时维护',
    status: 'expired'
  }
])

const showDialog = ref(false)
const showAddDialog = ref(false)
const dialogTitle = ref('新增静默规则')
const formData = ref({
  id: 0,
  matchers: '',
  startsAt: new Date(),
  endsAt: new Date(),
  createdBy: 'admin',
  comment: '',
  status: 'active'
})

// 显示编辑对话框
const showEditDialog = (row: any) => {
  formData.value = { 
    ...row,
    startsAt: new Date(row.startsAt),
    endsAt: new Date(row.endsAt)
  }
  dialogTitle.value = '编辑静默规则'
  showDialog.value = true
}

// 重置表单
const resetForm = () => {
  formData.value = {
    id: 0,
    matchers: '',
    startsAt: new Date(),
    endsAt: new Date(Date.now() + 3600 * 1000), // 默认1小时后结束
    createdBy: 'admin',
    comment: '',
    status: 'active'
  }
}

// 保存静默规则
const saveSilence = () => {
  if (formData.value.id) {
    // 编辑
    const index = silences.value.findIndex(item => item.id === formData.value.id)
    if (index !== -1) {
      silences.value[index] = { 
        ...formData.value,
        startsAt: formData.value.startsAt.toISOString(),
        endsAt: formData.value.endsAt.toISOString()
      }
    }
  } else {
    // 新增
    const newSilence = {
      ...formData.value,
      id: Date.now(),
      startsAt: formData.value.startsAt.toISOString(),
      endsAt: formData.value.endsAt.toISOString()
    }
    silences.value.push(newSilence)
  }
  showDialog.value = false
  resetForm()
}

// 删除静默规则
const deleteSilence = (id: number) => {
  silences.value = silences.value.filter(item => item.id !== id)
}

// 格式化日期
const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleString()
}

// 监听新增对话框
const watchAddDialog = (newVal: boolean) => {
  if (newVal) {
    dialogTitle.value = '新增静默规则'
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
.silences-view {
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