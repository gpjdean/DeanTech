<template>
  <div class="inhibitions-view">
    <el-card shadow="never" class="page-header">
      <div class="header-content">
        <h2>抑制规则</h2>
        <div class="header-actions">
          <el-button type="primary" @click="showAddDialog = true">
            <el-icon><Plus /></el-icon>
            新增抑制规则
          </el-button>
        </div>
      </div>
    </el-card>

    <el-card shadow="never" class="content-card">
      <el-table :data="inhibitions" style="width: 100%">
        <el-table-column prop="sourceMatch" label="源匹配规则" min-width="200">
          <template #default="scope">
            <el-popover
              placement="top"
              title="源匹配规则"
              width="300"
              trigger="hover"
            >
              <pre>{{ JSON.parse(scope.row.sourceMatch) }}</pre>
              <template #reference>
                <el-button type="text" size="small">查看</el-button>
              </template>
            </el-popover>
          </template>
        </el-table-column>
        <el-table-column prop="targetMatch" label="目标匹配规则" min-width="200">
          <template #default="scope">
            <el-popover
              placement="top"
              title="目标匹配规则"
              width="300"
              trigger="hover"
            >
              <pre>{{ JSON.parse(scope.row.targetMatch) }}</pre>
              <template #reference>
                <el-button type="text" size="small">查看</el-button>
              </template>
            </el-popover>
          </template>
        </el-table-column>
        <el-table-column prop="equal" label="相等标签" min-width="150">
          <template #default="scope">
            <el-popover
              placement="top"
              title="相等标签"
              width="300"
              trigger="hover"
            >
              <pre>{{ JSON.parse(scope.row.equal) }}</pre>
              <template #reference>
                <el-button type="text" size="small">查看</el-button>
              </template>
            </el-popover>
          </template>
        </el-table-column>
        <el-table-column prop="createdBy" label="创建人" width="120"></el-table-column>
        <el-table-column prop="comment" label="备注" min-width="150"></el-table-column>
        <el-table-column prop="isEnabled" label="状态" width="100">
          <template #default="scope">
            <el-switch v-model="scope.row.isEnabled" @change="toggleInhibitionStatus(scope.row)"></el-switch>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" @click="showEditDialog(scope.row)">
              编辑
            </el-button>
            <el-button type="danger" size="small" @click="deleteInhibition(scope.row.id)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新增/编辑抑制规则对话框 -->
    <el-dialog
      v-model="showDialog"
      :title="dialogTitle"
      width="600px"
      @close="resetForm"
    >
      <el-form :model="formData" label-width="120px">
        <el-form-item label="源匹配规则" required>
          <el-input
            v-model="formData.sourceMatch"
            type="textarea"
            :rows="3"
            placeholder='请输入源匹配规则，JSON格式，例如：[{"name":"severity","value":"critical","isRegex":false}]'
          ></el-input>
        </el-form-item>
        <el-form-item label="目标匹配规则" required>
          <el-input
            v-model="formData.targetMatch"
            type="textarea"
            :rows="3"
            placeholder='请输入目标匹配规则，JSON格式，例如：[{"name":"severity","value":"warning","isRegex":false}]'
          ></el-input>
        </el-form-item>
        <el-form-item label="相等标签">
          <el-input
            v-model="formData.equal"
            type="textarea"
            :rows="2"
            placeholder='请输入相等标签，JSON格式，例如：["instance","job"]'
          ></el-input>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="formData.comment" placeholder="请输入备注"></el-input>
        </el-form-item>
        <el-form-item label="创建人" required>
          <el-input v-model="formData.createdBy" placeholder="请输入创建人"></el-input>
        </el-form-item>
        <el-form-item label="启用状态">
          <el-switch v-model="formData.isEnabled"></el-switch>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showDialog = false">取消</el-button>
          <el-button type="primary" @click="saveInhibition">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Plus } from '@element-plus/icons-vue'

// 模拟数据
const inhibitions = ref([
  {
    id: 1,
    sourceMatch: JSON.stringify([{ name: 'severity', value: 'critical', isRegex: false }]),
    targetMatch: JSON.stringify([{ name: 'severity', value: 'warning', isRegex: false }]),
    equal: JSON.stringify(['instance', 'job']),
    createdBy: 'admin',
    comment: '当存在critical告警时，抑制warning告警',
    isEnabled: true
  },
  {
    id: 2,
    sourceMatch: JSON.stringify([{ name: 'alertname', value: 'NodeDown', isRegex: false }]),
    targetMatch: JSON.stringify([{ name: 'job', value: 'node_exporter', isRegex: false }]),
    equal: JSON.stringify(['instance']),
    createdBy: 'admin',
    comment: '当节点宕机时，抑制该节点的所有其他告警',
    isEnabled: true
  }
])

const showDialog = ref(false)
const showAddDialog = ref(false)
const dialogTitle = ref('新增抑制规则')
const formData = ref({
  id: 0,
  sourceMatch: '',
  targetMatch: '',
  equal: '',
  createdBy: 'admin',
  comment: '',
  isEnabled: true
})

// 显示编辑对话框
const showEditDialog = (row: any) => {
  formData.value = { ...row }
  dialogTitle.value = '编辑抑制规则'
  showDialog.value = true
}

// 重置表单
const resetForm = () => {
  formData.value = {
    id: 0,
    sourceMatch: '',
    targetMatch: '',
    equal: '',
    createdBy: 'admin',
    comment: '',
    isEnabled: true
  }
}

// 保存抑制规则
const saveInhibition = () => {
  if (formData.value.id) {
    // 编辑
    const index = inhibitions.value.findIndex(item => item.id === formData.value.id)
    if (index !== -1) {
      inhibitions.value[index] = { ...formData.value }
    }
  } else {
    // 新增
    const newInhibition = {
      ...formData.value,
      id: Date.now()
    }
    inhibitions.value.push(newInhibition)
  }
  showDialog.value = false
  resetForm()
}

// 切换抑制规则状态
const toggleInhibitionStatus = (row: any) => {
  // 这里可以添加API调用，更新抑制规则状态
  console.log('Toggle inhibition status:', row.id, row.isEnabled)
}

// 删除抑制规则
const deleteInhibition = (id: number) => {
  inhibitions.value = inhibitions.value.filter(item => item.id !== id)
}

// 监听新增对话框
const watchAddDialog = (newVal: boolean) => {
  if (newVal) {
    dialogTitle.value = '新增抑制规则'
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
.inhibitions-view {
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