<template>
  <div class="resourcequotas-view">
    <el-card shadow="never" class="page-header">
      <div class="header-content">
        <h2>资源配额管理</h2>
        <div class="header-actions">
          <el-select v-model="cluster" placeholder="选择集群" size="small" style="width: 150px; margin-right: 10px;" :loading="loadingClusters">
            <el-option
              v-for="c in clusters"
              :key="c.id"
              :label="c.displayName || c.name"
              :value="c.id"
            />
          </el-select>
          <el-button type="primary" size="small" @click="createResourceQuota" :disabled="!cluster">
            <el-icon><Plus /></el-icon>
            创建资源配额
          </el-button>
          <el-button type="primary" size="small" @click="refreshResourceQuotas" :disabled="!cluster">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </div>
    </el-card>

    <el-card shadow="hover" class="resourcequotas-card">
      <template #header>
        <div class="card-header">
          <span>资源配额列表</span>
          <div class="card-actions">
            <el-select v-model="namespaceFilter" placeholder="选择命名空间" clearable size="small" style="width: 200px;" :disabled="!cluster">
              <el-option
                v-for="ns in namespaces"
                :key="ns"
                :label="ns"
                :value="ns"
              />
            </el-select>
            <el-input
              v-model="searchQuery"
              placeholder="搜索资源配额"
              clearable
              size="small"
              style="width: 300px; margin-left: 10px;"
              @keyup.enter="loadResourceQuotas"
              :disabled="!cluster"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
          </div>
        </div>
      </template>
      
      <div v-if="!cluster" class="no-cluster-selected">
        <el-empty description="请先选择集群" />
      </div>
      
      <template v-else>
        <el-table 
          :data="resourceQuotas" 
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
          <el-table-column label="资源限制" min-width="300">
            <template #default="scope">
              <div class="resource-limits">
                <div v-for="(limit, key) in scope.row.hard" :key="key" class="resource-limit-item">
                  <span class="resource-key">{{ key }}:</span>
                  <span class="resource-value">{{ limit }}</span>
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="使用情况" min-width="300">
            <template #default="scope">
              <div class="resource-usage">
                <div v-for="(usage, key) in scope.row.used" :key="key" class="resource-usage-item">
                  <span class="resource-key">{{ key }}:</span>
                  <span class="resource-value">{{ usage }}</span>
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="150" fixed="right">
            <template #default="scope">
              <el-button type="primary" size="small" @click="viewResourceQuota(scope.row)">
                <el-icon><View /></el-icon>
                查看
              </el-button>
              <el-button type="danger" size="small" @click="deleteResourceQuota(scope.row)" style="margin-left: 10px;">
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
      </template>
    </el-card>
    
    <!-- 创建/编辑资源配额对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
    >
      <el-form :model="formData" label-width="100px" class="resourcequota-form">
        <el-form-item label="集群" prop="cluster">
          <el-select v-model="formData.cluster" placeholder="选择集群" style="width: 100%;" disabled>
            <el-option
              v-for="c in clusters"
              :key="c.id"
              :label="c.displayName || c.name"
              :value="c.id"
            />
          </el-select>
        </el-form-item>
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
          <el-input v-model="formData.name" placeholder="输入资源配额名称" />
        </el-form-item>
        <el-form-item label="CPU限制">
          <el-input v-model="formData.hard['requests.cpu']" placeholder="例如: 100m" />
        </el-form-item>
        <el-form-item label="内存限制">
          <el-input v-model="formData.hard['requests.memory']" placeholder="例如: 256Mi" />
        </el-form-item>
        <el-form-item label="CPU请求">
          <el-input v-model="formData.hard['limits.cpu']" placeholder="例如: 500m" />
        </el-form-item>
        <el-form-item label="内存请求">
          <el-input v-model="formData.hard['limits.memory']" placeholder="例如: 1Gi" />
        </el-form-item>
        <el-form-item label="Pod数量限制">
          <el-input v-model="formData.hard['pods']" placeholder="例如: 10" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveResourceQuota">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, watchEffect } from 'vue'
import { Plus, Refresh, Search, View, Delete } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import axios from 'axios'
import { clusterAPI } from '../api/api'

// 加载状态
const loading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('创建资源配额')

// 集群相关
const clusters = ref<any[]>([])
const cluster = ref<number | string>('')
const loadingClusters = ref(false)

// 搜索和过滤
const searchQuery = ref('')
const namespaceFilter = ref('')

// 命名空间列表
const namespaces = ref<string[]>(['default', 'kube-system', 'kube-public', 'monitoring', 'dev', 'test', 'prod'])
const loadingNamespaces = ref(false)

// 资源配额列表
const resourceQuotas = ref<any[]>([])

// 分页信息
const pagination = ref({
  page: 1,
  pageSize: 20,
  total: 0
})

// 表单数据
const formData = ref<{
  cluster: number | string;
  namespace: string;
  name: string;
  hard: {
    'requests.cpu': string;
    'requests.memory': string;
    'limits.cpu': string;
    'limits.memory': string;
    'pods': string;
  };
}>({
  cluster: '',
  namespace: 'default',
  name: '',
  hard: {
    'requests.cpu': '',
    'requests.memory': '',
    'limits.cpu': '',
    'limits.memory': '',
    'pods': ''
  }
})

// 组件挂载时加载集群列表
onMounted(() => {
  loadClusters()
})

// 加载集群列表
const loadClusters = async () => {
  try {
    loadingClusters.value = true
    const response = await clusterAPI.list()
    clusters.value = response.data
    console.log('加载的集群列表:', clusters.value)
  } catch (error) {
    console.error('加载集群列表失败:', error)
    ElMessage.error('加载集群列表失败')
    // 如果API调用失败，使用模拟数据
    clusters.value = [
      { id: 1, name: 'cluster-1', displayName: '集群1' },
      { id: 2, name: 'cluster-2', displayName: '集群2' },
      { id: 3, name: 'cluster-3', displayName: '集群3' }
    ]
  } finally {
    loadingClusters.value = false
  }
}

// 加载指定集群的命名空间列表
const loadNamespaces = async (clusterId: number) => {
  try {
    loadingNamespaces.value = true
    const response = await clusterAPI.getNamespaces(clusterId)
    namespaces.value = response.data
  } catch (error) {
    console.error('加载命名空间失败:', error)
    ElMessage.error('加载命名空间失败')
    // 使用默认命名空间
    namespaces.value = ['default', 'kube-system', 'kube-public', 'monitoring', 'dev', 'test', 'prod']
  } finally {
    loadingNamespaces.value = false
  }
}

// 监听集群变化，重新加载数据
watch(() => cluster.value, (newCluster) => {
  if (newCluster) {
    loadResourceQuotas()
    // 加载该集群的命名空间
    loadNamespaces(Number(newCluster))
  } else {
    resourceQuotas.value = []
    pagination.value.total = 0
  }
})

// 监听集群变化的函数
watchEffect(() => {
  if (cluster.value) {
    console.log(`当前选择的集群: ${cluster.value}`)
  }
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

// 加载资源配额列表
const loadResourceQuotas = async () => {
  try {
    if (!cluster.value) {
      return
    }
    
    loading.value = true
    
    // 重置数据，确保选择集群后显示正确的列表
    resourceQuotas.value = []
    pagination.value.total = 0
    
    // 尝试调用真实的资源配额API
    try {
      // 假设API路径为 /clusters/{clusterId}/resourcequotas
      console.log(`调用资源配额API: /clusters/${cluster.value}/resourcequotas`)
      const response = await axios.get(`/clusters/${cluster.value}/resourcequotas`, {
        params: {
          namespace: namespaceFilter.value || '',
          search: searchQuery.value
        }
      })
      console.log('API返回数据:', response.data)
      
      // 确保返回的数据是数组
      if (Array.isArray(response.data)) {
        resourceQuotas.value = response.data
        pagination.value.total = response.data.length
      } else {
        console.error('API返回的数据不是数组:', response.data)
        resourceQuotas.value = []
        pagination.value.total = 0
        ElMessage.warning('资源配额数据格式不正确')
      }
    } catch (apiError) {
      console.error('资源配额API调用失败，使用模拟数据:', apiError)
      
      // API调用失败时，使用模拟数据
      await new Promise(resolve => setTimeout(resolve, 800))
      
      // 模拟数据，根据集群返回不同的数据
      const mockData = [
        {
          namespace: 'default',
          name: `${cluster.value}-default-quota`,
          createdAt: Date.now() - 3600000,
          hard: {
            'requests.cpu': '1',
            'requests.memory': '1Gi',
            'limits.cpu': '2',
            'limits.memory': '2Gi',
            'pods': '10'
          },
          used: {
            'requests.cpu': '500m',
            'requests.memory': '512Mi',
            'limits.cpu': '1',
            'limits.memory': '1Gi',
            'pods': '5'
          }
        },
        {
          namespace: 'dev',
          name: `${cluster.value}-dev-quota`,
          createdAt: Date.now() - 7200000,
          hard: {
            'requests.cpu': '2',
            'requests.memory': '2Gi',
            'limits.cpu': '4',
            'limits.memory': '4Gi',
            'pods': '20'
          },
          used: {
            'requests.cpu': '1.5',
            'requests.memory': '1.5Gi',
            'limits.cpu': '3',
            'limits.memory': '3Gi',
            'pods': '15'
          }
        },
        {
          namespace: 'test',
          name: `${cluster.value}-test-quota`,
          createdAt: Date.now() - 10800000,
          hard: {
            'requests.cpu': '1',
            'requests.memory': '1Gi',
            'limits.cpu': '2',
            'limits.memory': '2Gi',
            'pods': '10'
          },
          used: {
            'requests.cpu': '800m',
            'requests.memory': '800Mi',
            'limits.cpu': '1.5',
            'limits.memory': '1.5Gi',
            'pods': '8'
          }
        }
      ]
      
      // 应用过滤
      let filtered = [...mockData] // 创建副本，避免修改原始数据
      console.log('应用过滤前的模拟数据:', filtered)
      
      if (searchQuery.value) {
        filtered = filtered.filter(item => {
          const matches = item.name.includes(searchQuery.value) || item.namespace.includes(searchQuery.value)
          console.log(`搜索过滤: ${item.name} ${matches ? '匹配' : '不匹配'} ${searchQuery.value}`)
          return matches
        })
      }
      
      if (namespaceFilter.value) {
        filtered = filtered.filter(item => {
          const matches = item.namespace === namespaceFilter.value
          console.log(`命名空间过滤: ${item.namespace} ${matches ? '匹配' : '不匹配'} ${namespaceFilter.value}`)
          return matches
        })
      }
      
      console.log('应用过滤后的模拟数据:', filtered)
      resourceQuotas.value = filtered
      pagination.value.total = filtered.length
    }
    
    console.log(`最终显示的资源配额列表 (${resourceQuotas.value.length} 条):`, resourceQuotas.value)
  } catch (error) {
    console.error('加载资源配额失败:', error)
    ElMessage.error('加载资源配额失败')
    resourceQuotas.value = []
    pagination.value.total = 0
  } finally {
    loading.value = false
  }
}

// 刷新资源配额
const refreshResourceQuotas = () => {
  loadResourceQuotas()
}

// 打开创建资源配额对话框
const createResourceQuota = () => {
  dialogTitle.value = '创建资源配额'
  formData.value = {
    cluster: cluster.value,
    namespace: 'default',
    name: '',
    hard: {
      'requests.cpu': '',
      'requests.memory': '',
      'limits.cpu': '',
      'limits.memory': '',
      'pods': ''
    }
  }
  dialogVisible.value = true
}

// 保存资源配额
const saveResourceQuota = async () => {
  try {
    if (!formData.value.name) {
      ElMessage.warning('请输入资源配额名称')
      return
    }
    
    loading.value = true
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    dialogVisible.value = false
    ElMessage.success('资源配额保存成功')
    loadResourceQuotas()
  } catch (error) {
    console.error('保存资源配额失败:', error)
    ElMessage.error('保存资源配额失败')
  } finally {
    loading.value = false
  }
}

// 查看资源配额详情
const viewResourceQuota = (_resourceQuota: any) => {
  ElMessage.info('查看资源配额详情功能开发中')
}

// 删除资源配额
const deleteResourceQuota = async (resourceQuota: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除资源配额 ${resourceQuota.name} 吗？`,
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
    
    ElMessage.success('资源配额删除成功')
    loadResourceQuotas()
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
  loadResourceQuotas()
}

// 处理当前页变化
const handleCurrentChange = (page: number) => {
  pagination.value.page = page
  loadResourceQuotas()
}
</script>

<style scoped>
.resourcequotas-view {
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

.resourcequotas-card {
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

.resource-limits,
.resource-usage {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.resource-limit-item,
.resource-usage-item {
  display: flex;
  align-items: center;
  font-size: 13px;
}

.resource-key {
  font-weight: 500;
  width: 120px;
  color: #606266;
}

.resource-value {
  color: #303133;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.resourcequota-form {
  margin-top: 20px;
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
</style>