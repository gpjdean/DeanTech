<template>
  <div class="limitranges-view">
    <el-card shadow="never" class="page-header">
      <div class="header-content">
        <h2>限制范围管理</h2>
        <div class="header-actions">
          <el-button type="primary" size="small" @click="createLimitRange">
            <el-icon><Plus /></el-icon>
            创建限制范围
          </el-button>
          <el-button type="primary" size="small" @click="refreshLimitRanges">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </div>
    </el-card>

    <el-card shadow="hover" class="limitranges-card">
      <template #header>
        <div class="card-header">
          <span>限制范围列表</span>
          <div class="card-actions">
            <el-select v-model="namespaceFilter" placeholder="选择命名空间" clearable size="small" style="width: 200px;">
              <el-option
                v-for="ns in namespaces"
                :key="ns"
                :label="ns"
                :value="ns"
              />
            </el-select>
            <el-input
              v-model="searchQuery"
              placeholder="搜索限制范围"
              clearable
              size="small"
              style="width: 300px; margin-left: 10px;"
              @keyup.enter="loadLimitRanges"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
          </div>
        </div>
      </template>
      
      <el-table 
        :data="limitRanges" 
        style="width: 100%"
        stripe
        v-loading="loading"
        :row-key="(row: any) => `${row.namespace}-${row.name}`"
      >
        <el-table-column prop="namespace" label="命名空间" width="150" sortable />
        <el-table-column prop="name" label="名称" width="200" sortable />
        <el-table-column prop="createdAt" label="创建时间" width="180" sortable>
          <template #default="scope">
            {{ formatDateTime(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="默认限制" min-width="250">
          <template #default="scope">
            <div class="limit-details">
              <div v-if="scope.row.spec.limits.default" class="limit-item">
                <span class="limit-type">{{ scope.row.spec.limits[0].type }}:</span>
                <div v-for="(value, key) in scope.row.spec.limits[0].default" :key="key" class="limit-value">
                  {{ key }}: {{ value }}
                </div>
              </div>
              <span v-else class="no-data">无</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="默认请求" min-width="250">
          <template #default="scope">
            <div class="limit-details">
              <div v-if="scope.row.spec.limits.defaultRequest" class="limit-item">
                <span class="limit-type">{{ scope.row.spec.limits[0].type }}:</span>
                <div v-for="(value, key) in scope.row.spec.limits[0].defaultRequest" :key="key" class="limit-value">
                  {{ key }}: {{ value }}
                </div>
              </div>
              <span v-else class="no-data">无</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" @click="viewLimitRange(scope.row)">
              <el-icon><View /></el-icon>
              查看
            </el-button>
            <el-button type="danger" size="small" @click="deleteLimitRange(scope.row)" style="margin-left: 10px;">
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="pagination.total"
          page-text="页"
          size-text="/页"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
    
    <!-- 创建/编辑限制范围对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
    >
      <el-form :model="formData" label-width="100px" class="limitrange-form">
        <el-form-item label="命名空间" prop="namespace">
          <el-select v-model="formData.namespace" placeholder="选择命名空间" style="width: 100%;">
            <el-option
              v-for="ns in namespaces"
              :key="ns"
              :label="ns"
              :value="ns"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <el-input v-model="formData.name" placeholder="输入限制范围名称" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="formData.spec.limits[0].type" placeholder="选择类型" style="width: 100%;">
            <el-option label="Container" value="Container" />
            <el-option label="Pod" value="Pod" />
            <el-option label="PersistentVolumeClaim" value="PersistentVolumeClaim" />
          </el-select>
        </el-form-item>
        <el-divider content-position="left">默认限制</el-divider>
        <el-form-item label="CPU限制">
          <el-input v-model="formData.spec.limits[0].default.cpu" placeholder="例如: 500m" />
        </el-form-item>
        <el-form-item label="内存限制">
          <el-input v-model="formData.spec.limits[0].default.memory" placeholder="例如: 512Mi" />
        </el-form-item>
        <el-divider content-position="left">默认请求</el-divider>
        <el-form-item label="CPU请求">
          <el-input v-model="formData.spec.limits[0].defaultRequest.cpu" placeholder="例如: 100m" />
        </el-form-item>
        <el-form-item label="内存请求">
          <el-input v-model="formData.spec.limits[0].defaultRequest.memory" placeholder="例如: 256Mi" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveLimitRange">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus, Refresh, Search, View, Delete } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'

// 加载状态
const loading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('创建限制范围')

// 搜索和过滤
const searchQuery = ref('')
const namespaceFilter = ref('')

// 命名空间列表
const namespaces = ref(['default', 'kube-system', 'kube-public', 'monitoring', 'dev', 'test', 'prod'])

// 限制范围列表
const limitRanges = ref<any[]>([])

// 分页信息
const pagination = ref({
  page: 1,
  pageSize: 20,
  total: 0
})

// 表单数据
const formData = ref({
  namespace: 'default',
  name: '',
  spec: {
    limits: [
      {
        type: 'Container',
        default: {
          cpu: '',
          memory: ''
        },
        defaultRequest: {
          cpu: '',
          memory: ''
        }
      }
    ]
  }
})

// 组件挂载时加载数据
onMounted(() => {
  loadLimitRanges()
})

// 格式化日期时间
const formatDateTime = (timestamp: any) => {
  if (!timestamp) return ''
  try {
    const date = new Date(timestamp)
    return date.toLocaleString('zh-CN')
  } catch (error) {
    return ''
  }
}

// 加载限制范围列表
const loadLimitRanges = async () => {
  try {
    loading.value = true
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // 模拟数据
    const mockData = [
      {
        namespace: 'default',
        name: 'default-limit-range',
        createdAt: Date.now() - 3600000,
        spec: {
          limits: [
            {
              type: 'Container',
              default: {
                cpu: '500m',
                memory: '512Mi'
              },
              defaultRequest: {
                cpu: '100m',
                memory: '256Mi'
              },
              max: {
                cpu: '1',
                memory: '1Gi'
              },
              min: {
                cpu: '50m',
                memory: '128Mi'
              }
            }
          ]
        }
      },
      {
        namespace: 'dev',
        name: 'dev-limit-range',
        createdAt: Date.now() - 7200000,
        spec: {
          limits: [
            {
              type: 'Container',
              default: {
                cpu: '1',
                memory: '1Gi'
              },
              defaultRequest: {
                cpu: '200m',
                memory: '512Mi'
              }
            }
          ]
        }
      },
      {
        namespace: 'test',
        name: 'test-limit-range',
        createdAt: Date.now() - 10800000,
        spec: {
          limits: [
            {
              type: 'Container',
              default: {
                cpu: '500m',
                memory: '512Mi'
              },
              defaultRequest: {
                cpu: '100m',
                memory: '256Mi'
              }
            },
            {
              type: 'Pod',
              max: {
                cpu: '4',
                memory: '4Gi'
              },
              min: {
                cpu: '200m',
                memory: '512Mi'
              }
            }
          ]
        }
      }
    ]
    
    // 应用过滤
    let filtered = mockData
    if (searchQuery.value) {
      filtered = filtered.filter(item => item.name.includes(searchQuery.value))
    }
    if (namespaceFilter.value) {
      filtered = filtered.filter(item => item.namespace === namespaceFilter.value)
    }
    
    limitRanges.value = filtered
    pagination.value.total = filtered.length
  } catch (error) {
    console.error('加载限制范围失败:', error)
    ElMessage.error('加载限制范围失败')
    limitRanges.value = []
    pagination.value.total = 0
  } finally {
    loading.value = false
  }
}

// 刷新限制范围
const refreshLimitRanges = () => {
  loadLimitRanges()
}

// 打开创建限制范围对话框
const createLimitRange = () => {
  dialogTitle.value = '创建限制范围'
  formData.value = {
    namespace: 'default',
    name: '',
    spec: {
      limits: [
        {
          type: 'Container',
          default: {
            cpu: '',
            memory: ''
          },
          defaultRequest: {
            cpu: '',
            memory: ''
          }
        }
      ]
    }
  }
  dialogVisible.value = true
}

// 保存限制范围
const saveLimitRange = async () => {
  try {
    if (!formData.value.name) {
      ElMessage.warning('请输入限制范围名称')
      return
    }
    
    loading.value = true
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    dialogVisible.value = false
    ElMessage.success('限制范围保存成功')
    loadLimitRanges()
  } catch (error) {
    console.error('保存限制范围失败:', error)
    ElMessage.error('保存限制范围失败')
  } finally {
    loading.value = false
  }
}

// 查看限制范围详情
const viewLimitRange = (_limitRange: any) => {
  ElMessage.info('查看限制范围详情功能开发中')
}

// 删除限制范围
const deleteLimitRange = async (limitRange: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除限制范围 ${limitRange.name} 吗？`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    loading.value = true
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    ElMessage.success('限制范围删除成功')
    loadLimitRanges()
  } catch (error) {
    // 用户取消删除
  } finally {
    loading.value = false
  }
}

// 处理分页大小变化
const handleSizeChange = (size: number) => {
  pagination.value.pageSize = size
  pagination.value.page = 1
  loadLimitRanges()
}

// 处理当前页变化
const handleCurrentChange = (page: number) => {
  pagination.value.page = page
  loadLimitRanges()
}
</script>

<style scoped>
.limitranges-view {
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

.header-actions {
  display: flex;
  align-items: center;
}

.limitranges-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.card-actions {
  display: flex;
  align-items: center;
}

.limit-details {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.limit-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.limit-type {
  font-weight: 500;
  color: #606266;
  font-size: 13px;
}

.limit-value {
  font-size: 13px;
  color: #303133;
  margin-left: 10px;
}

.no-data {
  color: #909399;
  font-style: italic;
  font-size: 13px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.limitrange-form {
  margin-top: 20px;
}
</style>