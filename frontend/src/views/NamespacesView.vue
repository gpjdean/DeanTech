<template>
  <div class="namespaces-view">
    <el-card shadow="never" class="page-header">
      <div class="header-content">
        <h2>命名空间</h2>
        <div class="header-actions">
          <el-select v-model="selectedCluster" placeholder="选择集群" style="width: 200px; margin-right: 10px;">
            <el-option
              v-for="cluster in clusters"
              :key="cluster.id"
              :label="cluster.name"
              :value="cluster.id"
            />
          </el-select>
          <el-button type="primary" @click="openAddDialog">
            <el-icon><Plus /></el-icon>
            添加命名空间
          </el-button>
        </div>
      </div>
    </el-card>

    <el-card shadow="hover" class="namespaces-card">
      <template #header>
        <div class="card-header">
          <span>命名空间列表</span>
          <!-- 搜索区域 -->
          <div class="search-container" v-if="selectedCluster">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索命名空间名称"
              clearable
              style="width: 200px; margin-right: 10px;"
              @keyup.enter="handleSearch"
            />
            <el-select
              v-model="searchStatus"
              placeholder="筛选状态"
              style="width: 150px; margin-right: 10px;"
              clearable
              @change="handleSearch"
            >
              <el-option label="Active" value="Active" />
              <el-option label="Terminating" value="Terminating" />
            </el-select>
            <el-button size="small" @click="resetSearch">重置</el-button>
          </div>
        </div>
      </template>

      <!-- 命名空间列表 -->
      <div class="table-container">
        <el-table 
          v-if="selectedCluster" 
          :data="namespacesList" 
          stripe 
          style="width: 100%" 
          v-loading="listLoading"
        >
          <el-table-column prop="name" label="命名空间名称" width="200" />
          <el-table-column prop="status" label="状态" width="120">
            <template #default="scope">
              <el-tag :type="scope.row.status === 'Active' ? 'success' : 'warning'" size="small">
                {{ scope.row.status }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="createdAt" label="创建时间" width="200">
            <template #default="scope">
              {{ formatDate(scope.row.createdAt) }}
            </template>
          </el-table-column>
          <el-table-column prop="labels" label="标签" min-width="200">
            <template #default="scope">
              <el-tag size="small" v-for="(value, key) in scope.row.labels" :key="key" class="label-tag">
                {{ key }}: {{ value }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="180" fixed="right">
            <template #default="scope">
              <el-button type="primary" size="small" @click="openViewDialog(scope.row)">
                <el-icon><View /></el-icon>
                查看
              </el-button>
              <el-button type="danger" size="small" @click="deleteNamespace(scope.row)">
                <el-icon><Delete /></el-icon>
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
        
        <div v-else class="no-cluster-selected">
          <el-empty description="请选择一个集群查看命名空间列表" />
        </div>
      </div>

      <!-- 分页 -->
      <div class="pagination-container" v-if="selectedCluster">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 添加命名空间对话框 -->
    <el-dialog
      v-model="addDialogVisible"
      title="添加命名空间"
      width="600px"
      destroy-on-close
    >
      <el-form :model="namespaceForm" label-width="100px">
        <!-- 集群选择 -->
        <el-form-item label="选择集群" required>
          <el-select 
            v-model="selectedClusterForAdd" 
            placeholder="请选择集群" 
            style="width: 100%"
          >
            <el-option
              v-for="cluster in clusters"
              :key="cluster.id"
              :label="cluster.name"
              :value="cluster.id"
            />
          </el-select>
        </el-form-item>
        
        <!-- 基本信息 -->
        <el-form-item label="名称" required>
          <el-input 
            v-model="namespaceForm.name" 
            placeholder="请输入命名空间名称" 
            clearable 
            maxlength="63"
            show-word-limit
          />
          <div class="form-hint">命名空间名称只能包含小写字母、数字、连字符和点，并且不能以连字符或点开头或结尾</div>
        </el-form-item>
        
        <!-- 标签管理 -->
        <el-form-item label="标签">
          <div class="tags-container">
            <el-tag
              v-for="(value, key) in namespaceForm.labels"
              :key="key"
              closable
              @close="deleteLabel(key)"
              class="tag-item"
            >
              {{ key }}: {{ value }}
            </el-tag>
            <el-button 
              type="dashed" 
              size="small" 
              @click="showAddLabelDialog = true"
              class="add-tag-btn"
            >
              <el-icon><Plus /></el-icon> 添加标签
            </el-button>
          </div>
        </el-form-item>
        
        <!-- 注解管理 -->
        <el-form-item label="注解">
          <div class="tags-container">
            <el-tag
              v-for="(value, key) in namespaceForm.annotations"
              :key="key"
              closable
              @close="deleteAnnotation(key)"
              class="tag-item annotation-tag"
            >
              {{ key }}: {{ value }}
            </el-tag>
            <el-button 
              type="dashed" 
              size="small" 
              @click="showAddAnnotationDialog = true"
              class="add-tag-btn"
            >
              <el-icon><Plus /></el-icon> 添加注解
            </el-button>
          </div>
        </el-form-item>
      </el-form>
      
      <!-- 添加标签对话框 -->
      <el-dialog
        v-model="showAddLabelDialog"
        title="添加标签"
        width="400px"
      >
        <el-form :model="tagForm" label-width="80px">
          <el-form-item label="键" required>
            <el-input v-model="tagForm.key" placeholder="请输入标签键" maxlength="63" show-word-limit />
          </el-form-item>
          <el-form-item label="值" required>
            <el-input v-model="tagForm.value" placeholder="请输入标签值" maxlength="63" show-word-limit />
          </el-form-item>
        </el-form>
        <template #footer>
          <span class="dialog-footer">
            <el-button @click="showAddLabelDialog = false">取消</el-button>
            <el-button type="primary" @click="addLabel">确定</el-button>
          </span>
        </template>
      </el-dialog>
      
      <!-- 添加注解对话框 -->
      <el-dialog
        v-model="showAddAnnotationDialog"
        title="添加注解"
        width="400px"
      >
        <el-form :model="annotationForm" label-width="80px">
          <el-form-item label="键" required>
            <el-input v-model="annotationForm.key" placeholder="请输入注解键" maxlength="63" show-word-limit />
          </el-form-item>
          <el-form-item label="值">
            <el-input v-model="annotationForm.value" placeholder="请输入注解值" type="textarea" :rows="3" maxlength="256" show-word-limit />
          </el-form-item>
        </el-form>
        <template #footer>
          <span class="dialog-footer">
            <el-button @click="showAddAnnotationDialog = false">取消</el-button>
            <el-button type="primary" @click="addAnnotation">确定</el-button>
          </span>
        </template>
      </el-dialog>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="addDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="addNamespace">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 查看命名空间对话框 -->
    <el-dialog
      v-model="viewDialogVisible"
      title="查看命名空间"
      width="800px"
    >
      <div class="namespace-view-dialog">
        <!-- 基本信息 -->
        <el-card shadow="never" class="info-card">
          <template #header>
            <div class="card-header">
              <span>基本信息</span>
            </div>
          </template>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="命名空间名称">{{ namespaceForm.name }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="namespaceForm.status === 'Active' ? 'success' : 'warning'" size="small">
                {{ namespaceForm.status }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ formatDate(namespaceForm.createdAt) }}</el-descriptions-item>
            <el-descriptions-item label="标签数量">{{ Object.keys(namespaceForm.labels || {}).length }}</el-descriptions-item>
            <el-descriptions-item label="注解数量">{{ Object.keys(namespaceForm.annotations || {}).length }}</el-descriptions-item>
            <el-descriptions-item label="">
              <!-- 占位符，保持布局对齐 -->
            </el-descriptions-item>
          </el-descriptions>
        </el-card>

        <!-- 标签信息 -->
        <el-card shadow="never" class="info-card" v-if="Object.keys(namespaceForm.labels || {}).length > 0">
          <template #header>
            <div class="card-header">
              <span>标签</span>
            </div>
          </template>
          <div class="labels-content">
            <div 
              v-for="(value, key) in namespaceForm.labels" 
              :key="key" 
              class="label-item"
            >
              <div class="label-key">{{ key }}</div>
              <div class="label-value">{{ value }}</div>
            </div>
          </div>
        </el-card>
        
        <!-- 注解信息 -->
        <el-card shadow="never" class="info-card" v-if="Object.keys(namespaceForm.annotations || {}).length > 0">
          <template #header>
            <div class="card-header">
              <span>注解</span>
            </div>
          </template>
          <div class="annotation-content">
            <div 
              v-for="(value, key) in namespaceForm.annotations" 
              :key="key" 
              class="annotation-item"
            >
              <div class="annotation-key">{{ key }}</div>
              <div class="annotation-value">{{ value }}</div>
            </div>
          </div>
        </el-card>

        <!-- YAML信息 -->
        <el-card shadow="never" class="info-card">
          <template #header>
            <div class="card-header">
              <span>YAML配置</span>
            </div>
          </template>
          <el-input
            v-model="namespaceYaml"
            type="textarea"
            :rows="15"
            readonly
            monospace
            placeholder="加载中..."
            style="font-family: Consolas, Monaco, 'Courier New', monospace;"
          />
        </el-card>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="viewDialogVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Delete, View } from '@element-plus/icons-vue'
import axios from 'axios'

// 命名空间数据类型定义
interface Namespace {
  name: string
  status: string
  createdAt: string
  labels: Record<string, string>
  annotations: Record<string, string>
}

// 集群数据类型定义
interface Cluster {
  id: number
  name: string
  apiURL: string
  token: string
  kubeconfig: string
  description: string
  status: string
  version: string
  createdAt: string
  connectionType: string
  nodeCount?: string
}

// 命名空间列表
const namespacesList = ref<Namespace[]>([])
// 列表加载状态
const listLoading = ref(false)
// 当前页码
const currentPage = ref(1)
// 每页大小
const pageSize = ref(10)
// 总条数
const total = ref(0)

// 集群列表
const clusters = ref<Cluster[]>([])
// 选中的集群
const selectedCluster = ref<number | undefined>(undefined)

// 对话框状态
const addDialogVisible = ref(false)
const viewDialogVisible = ref(false)

// 命名空间YAML内容
const namespaceYaml = ref('')

// 表单数据
const namespaceForm = ref<Namespace>({
  name: '',
  status: 'Active',
  createdAt: '',
  labels: {},
  annotations: {}
})

// 标签表单
const tagForm = ref({
  key: '',
  value: ''
})

// 注解表单
const annotationForm = ref({
  key: '',
  value: ''
})

// 对话框状态
const showAddLabelDialog = ref(false)
const showAddAnnotationDialog = ref(false)

// 添加命名空间时选择的集群
const selectedClusterForAdd = ref<number | undefined>(undefined)

// 搜索相关
const searchKeyword = ref('')
const searchStatus = ref('')

// 格式化日期函数
const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  
  // 手动格式化日期，确保格式为 YYYY-MM-DD HH:MM:SS
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

// 生成命名空间YAML函数
const generateNamespaceYaml = (namespace: Namespace) => {
  const yamlObj = {
    apiVersion: 'v1',
    kind: 'Namespace',
    metadata: {
      name: namespace.name,
      labels: namespace.labels || {},
      annotations: namespace.annotations || {}
    },
    status: {
      phase: namespace.status
    }
  }

  // 简单的YAML转换函数
  const toYaml = (obj: any, indent: number = 0): string => {
    let yaml = ''
    const spaces = '  '.repeat(indent)

    for (const [key, value] of Object.entries(obj)) {
      if (value === null || value === undefined) {
        continue
      }

      if (typeof value === 'object' && value !== null) {
        if (Array.isArray(value)) {
          yaml += `${spaces}${key}:\n`
          for (const item of value) {
            yaml += `${spaces}- `
            if (typeof item === 'object' && item !== null) {
              yaml += `\n${toYaml(item, indent + 2)}`
            } else {
              yaml += `${item}\n`
            }
          }
        } else {
          const hasEntries = Object.keys(value).length > 0
          yaml += `${spaces}${key}:${hasEntries ? '\n' : ' {}\n'}`
          if (hasEntries) {
            yaml += toYaml(value, indent + 1)
          }
        }
      } else {
        yaml += `${spaces}${key}: ${value}\n`
      }
    }

    return yaml
  }

  return toYaml(yamlObj)
}

// 组件挂载时加载数据
onMounted(() => {
  loadClusters()
})

// 监听集群选择变化，加载对应的命名空间
watch(selectedCluster, (newClusterId) => {
  if (newClusterId) {
    loadNamespaces(newClusterId)
  } else {
    // 重置命名空间列表
    namespacesList.value = []
    total.value = 0
  }
})

// 移除自动选择集群的逻辑，让用户手动选择

// 加载集群列表
const loadClusters = async () => {
  try {
    // 调用真实的API获取集群列表，与集群管理页面保持一致
    const response = await axios.get('/clusters')
    clusters.value = response.data
    console.log('集群列表加载成功:', clusters.value)
  } catch (error) {
    console.error('加载集群列表失败:', error)
    // 与集群管理页面保持一致，加载失败时返回空数组
    clusters.value = []
    ElMessage.error('加载集群列表失败')
  }
}

// 取消令牌，用于取消之前的请求
let cancelTokenSource: any = null

// 加载命名空间列表
const loadNamespaces = async (clusterId: number) => {
  // 取消之前的请求
  if (cancelTokenSource) {
    cancelTokenSource.cancel('Operation canceled by the user.')
  }
  
  // 创建新的取消令牌
  cancelTokenSource = axios.CancelToken.source()
  
  listLoading.value = true
  try {
    // 调用API获取命名空间列表，与其他资源页面保持一致
    const response = await axios.get(`/clusters/${clusterId}/namespaces`, {
      params: {
        page: currentPage.value,
        pageSize: pageSize.value,
        keyword: searchKeyword.value,
        status: searchStatus.value,
        _t: Date.now() // 添加时间戳，防止浏览器缓存
      },
      headers: {
        'Cache-Control': 'no-cache' // 禁用缓存
      },
      cancelToken: cancelTokenSource.token
    })
    
    // 更新命名空间列表，确保 createdAt 字段是字符串格式
    namespacesList.value = response.data.data.map((ns: any) => ({
      ...ns,
      createdAt: typeof ns.createdAt === 'string' ? ns.createdAt : new Date(ns.createdAt).toISOString()
    }))
    total.value = response.data.total
    listLoading.value = false
    
    // 如果返回的命名空间列表为空，根据搜索条件显示不同的提示信息
    if (namespacesList.value.length === 0) {
      if (searchKeyword.value || searchStatus.value) {
        // 如果有搜索条件，提示没有匹配的命名空间
        ElMessage.info('没有匹配的命名空间')
      } else {
        // 如果没有搜索条件，提示集群内无命名空间
        ElMessage.info('当前集群内无命名空间')
      }
    }
  } catch (error: any) {
    // 如果是取消请求的错误，不显示错误信息
    if (axios.isCancel(error)) {
      console.log('Request canceled:', error.message)
      return
    }
    console.error('加载命名空间列表失败:', error)
    listLoading.value = false
    ElMessage.error('加载命名空间列表失败')
    // 确保在错误情况下，列表为空，显示正确的空状态
    namespacesList.value = []
    total.value = 0
  }
}

// 打开添加对话框
const openAddDialog = () => {
  // 重置表单数据
  namespaceForm.value = {
    name: '',
    status: 'Active',
    createdAt: '',
    labels: {},
    annotations: {}
  }
  // 重置集群选择
  selectedClusterForAdd.value = undefined
  addDialogVisible.value = true
}

// 添加标签
const addLabel = () => {
  if (!tagForm.value.key || !tagForm.value.value) {
    ElMessage.warning('请输入标签键和值')
    return
  }
  
  // 检查标签键是否已存在
  if (namespaceForm.value.labels[tagForm.value.key]) {
    ElMessage.warning('该标签键已存在')
    return
  }
  
  // 添加标签
  namespaceForm.value.labels[tagForm.value.key] = tagForm.value.value
  
  // 重置标签表单
  tagForm.value = {
    key: '',
    value: ''
  }
  
  // 关闭对话框
  showAddLabelDialog.value = false
}

// 删除标签
const deleteLabel = (key: string) => {
  delete namespaceForm.value.labels[key]
}

// 添加注解
const addAnnotation = () => {
  if (!annotationForm.value.key) {
    ElMessage.warning('请输入注解键')
    return
  }
  
  // 检查注解键是否已存在
  if (namespaceForm.value.annotations[annotationForm.value.key]) {
    ElMessage.warning('该注解键已存在')
    return
  }
  
  // 添加注解
  namespaceForm.value.annotations[annotationForm.value.key] = annotationForm.value.value
  
  // 重置注解表单
  annotationForm.value = {
    key: '',
    value: ''
  }
  
  // 关闭对话框
  showAddAnnotationDialog.value = false
}

// 删除注解
const deleteAnnotation = (key: string) => {
  delete namespaceForm.value.annotations[key]
}

// 搜索命名空间
const handleSearch = () => {
  // 重置当前页码
  currentPage.value = 1
  // 重新加载命名空间列表
  if (selectedCluster.value) {
    loadNamespaces(selectedCluster.value)
  }
}

// 防抖函数
const debounce = (fn: Function, delay: number) => {
  let timer: ReturnType<typeof setTimeout> | null = null
  return (...args: any[]) => {
    if (timer) {
      clearTimeout(timer)
    }
    timer = setTimeout(() => {
      fn(...args)
      timer = null
    }, delay)
  }
}

// 防抖搜索函数
const debouncedSearch = debounce(() => {
  if (selectedCluster.value) {
    currentPage.value = 1
    loadNamespaces(selectedCluster.value)
  }
}, 300)

// 监听搜索关键词变化，自动搜索
watch(searchKeyword, () => {
  debouncedSearch()
})

// 重置搜索条件
const resetSearch = () => {
  searchKeyword.value = ''
  searchStatus.value = ''
  currentPage.value = 1
  if (selectedCluster.value) {
    loadNamespaces(selectedCluster.value)
  }
}

// 打开查看对话框
const openViewDialog = (row: Namespace) => {
  // 复制行数据到表单
  namespaceForm.value = JSON.parse(JSON.stringify(row))
  // 生成YAML内容
  namespaceYaml.value = generateNamespaceYaml(row)
  // 打开对话框
  viewDialogVisible.value = true
}

// 添加命名空间
const addNamespace = async () => {
  if (!selectedClusterForAdd.value) {
    ElMessage.warning('请选择集群')
    return
  }
  
  const clusterId = selectedClusterForAdd.value

  try {
    // 调用API添加命名空间
    await axios.post(`/clusters/${clusterId}/namespaces`, namespaceForm.value)
    
    // 如果当前选中的集群与添加的集群相同，则重新加载命名空间列表
    if (selectedCluster.value === clusterId) {
      await loadNamespaces(clusterId)
    }
    ElMessage.success('命名空间添加成功')
    addDialogVisible.value = false
  } catch (error: any) {
    console.error('添加命名空间失败:', error)
    ElMessage.error('添加命名空间失败')
  }
}

// 删除命名空间
const deleteNamespace = async (row: Namespace) => {
  if (!selectedCluster.value) return
  
  const clusterId = selectedCluster.value

  ElMessageBox.confirm(`确定要删除命名空间 "${row.name}" 吗？`, '删除确认', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      // 调用API删除命名空间
      await axios.delete(`/clusters/${clusterId}/namespaces/${row.name}`)
      
      // 重新加载命名空间列表
      await loadNamespaces(clusterId)
      ElMessage.success('命名空间删除成功')
    } catch (error: any) {
      console.error('删除命名空间失败:', error)
      ElMessage.error('删除命名空间失败')
    }
  }).catch(() => {
    // 取消删除
  })
}

// 处理分页大小变化
const handleSizeChange = (size: number) => {
  pageSize.value = size
  if (selectedCluster.value) {
    loadNamespaces(selectedCluster.value)
  }
}

// 处理页码变化
const handleCurrentChange = (page: number) => {
  currentPage.value = page
  if (selectedCluster.value) {
    loadNamespaces(selectedCluster.value)
  }
}
</script>

<style scoped>
.namespaces-view {
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

.namespaces-card {
  margin-bottom: 20px;
}

.search-container {
  display: flex;
  align-items: center;
  gap: 10px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.card-header span {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.table-container {
  margin-bottom: 20px;
}

.no-cluster-selected {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 300px;
  background-color: #fafafa;
  border-radius: 8px;
  margin-top: 20px;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

/* 表格样式 */
:deep(.el-table) {
  width: 100%;
}

:deep(.el-table__header-wrapper th.el-table__header-cell) {
  color: #333;
  font-weight: 600;
}

:deep(.el-table__body-wrapper td.el-table__cell) {
  color: #606266;
}

.label-tag {
  margin-right: 8px;
  margin-bottom: 8px;
  background-color: #f0f9ff;
  border-color: #91d5ff;
  color: #1890ff;
}

.annotation-tag {
  margin-right: 8px;
  margin-bottom: 8px;
  background-color: #f6ffed;
  border-color: #b7eb8f;
  color: #52c41a;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

/* 查看对话框样式 */
.namespace-view-dialog {
  padding: 10px 0;
}

.info-card {
  margin-bottom: 20px;
  border-radius: 8px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header span {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.label-list {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 10px;
}

.annotation-list {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 10px;
}

/* 确保标签和注解在卡片中正确显示 */
:deep(.el-card__body) {
  padding: 20px;
}

/* 确保描述列表样式正确 */
:deep(.el-descriptions) {
  margin-bottom: 0;
}

:deep(.el-descriptions__cell) {
  padding: 12px;
}

:deep(.el-descriptions__label) {
  font-weight: 600;
  background-color: #fafafa;
}

/* 添加命名空间对话框样式 */
.form-hint {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
  margin-left: 10px;
}

.tags-container {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 5px;
  padding: 10px;
  background-color: #fafafa;
  border-radius: 4px;
  min-height: 40px;
  align-items: center;
}

.tag-item {
  margin: 5px 0;
}

.add-tag-btn {
  margin: 5px 0;
}

:deep(.el-form-item__label) {
  font-weight: 500;
}

/* 标签和注解的标签样式 */
.annotation-tag {
  background-color: #f6ffed;
  border-color: #b7eb8f;
  color: #52c41a;
}

/* 标签内容样式 */
.labels-content {
  margin-top: 10px;
  padding: 10px 0;
}

/* 标签项样式 */
.label-item {
  display: flex;
  margin-bottom: 15px;
  padding: 12px;
  background-color: #f9fafb;
  border-left: 4px solid #1890ff;
  border-radius: 0 8px 8px 0;
  transition: all 0.3s ease;
}

.label-item:hover {
  background-color: #f0f9ff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

/* 标签键样式 */
.label-key {
  font-weight: 600;
  color: #333;
  min-width: 150px;
  max-width: 200px;
  flex-shrink: 0;
  padding-right: 15px;
  border-right: 1px solid #e5e7eb;
  word-break: break-all;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
}

/* 标签值样式 */
.label-value {
  flex: 1;
  padding-left: 15px;
  color: #606266;
  line-height: 1.6;
  word-break: break-all;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  white-space: pre-wrap;
  background-color: #ffffff;
  padding: 8px 12px;
  border-radius: 4px;
  border: 1px solid #e5e7eb;
}

/* 注解内容样式 */
.annotation-content {
  margin-top: 10px;
  padding: 10px 0;
}

/* 注解项样式 */
.annotation-item {
  display: flex;
  margin-bottom: 15px;
  padding: 12px;
  background-color: #f9fafb;
  border-left: 4px solid #52c41a;
  border-radius: 0 8px 8px 0;
  transition: all 0.3s ease;
}

.annotation-item:hover {
  background-color: #f0f9ff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

/* 注解键样式 */
.annotation-key {
  font-weight: 600;
  color: #333;
  min-width: 150px;
  max-width: 200px;
  flex-shrink: 0;
  padding-right: 15px;
  border-right: 1px solid #e5e7eb;
  word-break: break-all;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
}

/* 注解值样式 */
.annotation-value {
  flex: 1;
  padding-left: 15px;
  color: #606266;
  line-height: 1.6;
  word-break: break-all;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  white-space: pre-wrap;
  background-color: #ffffff;
  padding: 8px 12px;
  border-radius: 4px;
  border: 1px solid #e5e7eb;
}
</style>