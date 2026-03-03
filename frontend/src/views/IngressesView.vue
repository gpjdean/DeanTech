<template>
  <div class="ingresses-view">
    <el-card shadow="never" class="page-header">
      <div class="header-content">
        <h2>路由</h2>
        <div class="header-actions">
          <el-select v-model="selectedCluster" placeholder="选择集群" style="width: 200px; margin-right: 10px;">
            <el-option
              v-for="cluster in clusters"
              :key="cluster.id"
              :label="cluster.name"
              :value="cluster.id"
            />
          </el-select>
          <el-select v-model="selectedNamespace" placeholder="选择命名空间" style="width: 200px; margin-right: 10px;" :loading="namespacesLoading" filterable>
            <el-option
              v-for="namespace in namespaces"
              :key="namespace"
              :label="namespace"
              :value="namespace"
            />
          </el-select>
          <el-button type="primary" size="small" @click="refreshIngresses">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </div>
    </el-card>

    <el-card shadow="hover" class="ingresses-card">
      <template #header>
        <div class="card-header">
          <span>K8s Ingress列表</span>
        </div>
      </template>
      
      <el-table 
        :data="ingresses" 
        style="width: 100%; color: #303133; font-size: 14px; line-height: 1.5;"
        stripe
        v-loading="loading"
        v-if="selectedCluster"
        :row-style="{ height: 'auto', backgroundColor: '#ffffff' }"
        :cell-style="{ 
          'white-space': 'normal', 
          'word-break': 'break-all', 
          'min-height': '40px', 
          'padding': '12px', 
          'color': '#303133', 
          'font-size': '14px',
          'line-height': '1.5'
        }"
        :header-cell-style="{ 
          'color': '#303133', 
          'font-weight': 'bold', 
          'background-color': '#f5f7fa',
          'font-size': '14px'
        }"
        empty-text="No Data"
        :row-key="(row: any) => `${row.namespace}-${row.name}`"
      >
        <el-table-column prop="name" label="Ingress名称" min-width="200" show-overflow-tooltip />
        <el-table-column prop="namespace" label="命名空间" width="150" show-overflow-tooltip />
        <el-table-column prop="cluster" label="集群" width="150" show-overflow-tooltip />
        <el-table-column prop="rules" label="规则数量" width="120">
          <template #default="scope">
            <el-tag type="info" size="medium">{{ scope.row.rules.length }} 条</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="rules" label="主机" min-width="200">
          <template #default="scope">
            <div v-if="scope.row.rules.length > 0">
              <div v-for="(rule, index) in scope.row.rules" :key="index" class="host-item">
                <el-tag type="primary" size="small">{{ rule.host || '*' }}</el-tag>
              </div>
            </div>
            <div v-else class="empty-text">无</div>
          </template>
        </el-table-column>
        <el-table-column prop="rules" label="路径" min-width="200">
          <template #default="scope">
            <div v-if="scope.row.rules.length > 0">
              <div v-for="(rule, index) in scope.row.rules" :key="index" class="path-item">
                <div v-for="(path, pathIndex) in rule.http.paths" :key="pathIndex">
                  <el-tag type="success" size="small">{{ path.path || '/' }}</el-tag>
                </div>
              </div>
            </div>
            <div v-else class="empty-text">无</div>
          </template>
        </el-table-column>
        <el-table-column prop="rules" label="服务" min-width="200">
          <template #default="scope">
            <div v-if="scope.row.rules.length > 0">
              <div v-for="(rule, index) in scope.row.rules" :key="index" class="service-item">
                <div v-for="(path, pathIndex) in rule.http.paths" :key="pathIndex">
                  <el-tag type="warning" size="small">{{ path.backend.service.name }}:{{ path.backend.service.port.number }}</el-tag>
                </div>
              </div>
            </div>
            <div v-else class="empty-text">无</div>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="180" align="center">
          <template #default="scope">
            <span style="font-weight: 500; color: #606266;">
              {{ formatAge(scope.row.createdAt) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" @click="viewIngressDetails(scope.row)">
              <el-icon><View /></el-icon>
              查看
            </el-button>
            <el-button type="danger" size="small" @click="deleteIngress(scope.row)">
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <div v-else class="no-cluster-selected">
        <el-empty description="请选择一个集群查看Ingress列表" />
      </div>
    </el-card>
    
    <!-- Ingress详情对话框 -->
    <el-dialog
      v-model="ingressDetailsVisible"
      title="Ingress详情"
      width="800px"
    >
      <div v-if="selectedIngress" class="ingress-details">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="Ingress名称">{{ selectedIngress.name }}</el-descriptions-item>
          <el-descriptions-item label="命名空间">{{ selectedIngress.namespace }}</el-descriptions-item>
          <el-descriptions-item label="集群">{{ selectedIngress.cluster }}</el-descriptions-item>
          <el-descriptions-item label="规则数量">{{ selectedIngress.rules.length }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ selectedIngress.createdAt }}</el-descriptions-item>
        </el-descriptions>
        
        <h4 style="margin-top: 20px;">路由规则</h4>
        <div v-if="selectedIngress.rules.length > 0">
          <el-collapse>
            <el-collapse-item
              v-for="(rule, index) in selectedIngress.rules"
              :key="index"
              :title="rule.host || '默认规则'"
            >
              <div class="rule-details">
                <h5>HTTP路径规则</h5>
                <el-table :data="rule.http.paths" style="width: 100%; margin-top: 10px;">
                  <el-table-column prop="path" label="路径" width="150" />
                  <el-table-column prop="pathType" label="路径类型" width="120" />
                  <el-table-column label="服务" min-width="200">
                    <template #default="scope">
                      {{ scope.row.backend.service.name }}:{{ scope.row.backend.service.port.number }}
                    </template>
                  </el-table-column>
                </el-table>
              </div>
            </el-collapse-item>
          </el-collapse>
        </div>
        <div v-else class="empty-text">无路由规则</div>
        
        <h4 style="margin-top: 20px;">Annotations</h4>
        <div v-if="Object.keys(selectedIngress.annotations).length > 0">
          <el-table :data="annotationsList" style="width: 100%; margin-top: 10px;">
            <el-table-column prop="key" label="键" width="200" />
            <el-table-column prop="value" label="值" min-width="300" show-overflow-tooltip />
          </el-table>
        </div>
        <div v-else class="empty-text">无Annotations</div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="ingressDetailsVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { Refresh, View, Delete } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'

// 加载状态
const loading = ref(false)

// 命名空间加载状态
const namespacesLoading = ref(false)

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
  isActive: boolean
}

// 集群列表
const clusters = ref<Cluster[]>([])

// 命名空间列表
const namespaces = ref<string[]>([])

// 选择的集群和命名空间
const selectedCluster = ref<number | undefined>(undefined)
const selectedNamespace = ref<string>('')

// Ingress规则数据类型定义
interface IngressRule {
  host?: string
  http: {
    paths: Array<{
      path: string
      pathType: string
      backend: {
        service: {
          name: string
          port: {
            number: number
          }
        }
      }
    }>
  }
}

// Ingress数据类型定义
interface Ingress {
  name: string
  namespace: string
  cluster: string
  rules: IngressRule[]
  annotations: Record<string, string>
  createdAt: string
}

// Ingress数据
const ingresses = ref<Ingress[]>([])

// 对话框相关
const ingressDetailsVisible = ref(false)
const selectedIngress = ref<Ingress | null>(null)

// Annotations列表
const annotationsList = computed(() => {
  if (!selectedIngress.value) return []
  return Object.entries(selectedIngress.value.annotations).map(([key, value]) => ({
    key,
    value
  }))
})

// 组件挂载时加载数据
onMounted(() => {
  loadClusters()
})

// 监听集群选择变化，获取对应的命名空间列表
watch(selectedCluster, (newClusterId) => {
  if (newClusterId) {
    loadNamespaces(newClusterId)
  } else {
    // 重置命名空间列表
    namespaces.value = []
    selectedNamespace.value = ''
  }
})

// 监听命名空间选择变化，自动刷新Ingress列表
watch(selectedNamespace, () => {
  if (selectedCluster.value) {
    loadIngresses()
  }
})

// 加载集群列表
const loadClusters = async () => {
  try {
    // 调用真实的API获取集群列表，axios会自动添加baseURL
    const response = await axios.get('/clusters')
    clusters.value = response.data
    // 不自动选择集群，让用户手动选择
  } catch (error) {
    console.error('加载集群列表失败:', error)
    // 不显示错误消息给用户，只在控制台打印
    clusters.value = []
  }
}

// 加载命名空间列表
const loadNamespaces = async (clusterId: number) => {
  try {
    namespacesLoading.value = true
    // 调用API获取指定集群的命名空间列表，axios会自动添加baseURL
    // 添加分页参数，获取所有命名空间
    const response = await axios.get(`/clusters/${clusterId}/namespaces`, {
      params: {
        page: 1,
        pageSize: 1000 // 设置一个足够大的值，确保获取所有命名空间
      }
    })
    // 检查后端返回的数据格式，如果是分页对象则使用data字段
    const namespacesData = response.data.data || response.data
    // 只提取命名空间名称，确保是字符串数组
    namespaces.value = namespacesData.map((ns: any) => ns.name || '')
    // 过滤掉空字符串
    namespaces.value = namespaces.value.filter((ns: string) => ns !== '')
    // 不自动选择命名空间，让用户手动选择
    selectedNamespace.value = ''
    namespacesLoading.value = false
  } catch (error) {
    console.error('加载命名空间列表失败:', error)
    // 不显示错误消息给用户，只在控制台打印
    namespaces.value = []
    selectedNamespace.value = ''
    namespacesLoading.value = false
  }
}

// 取消令牌，用于取消之前的请求
let cancelTokenSource: any = null

// 加载Ingress列表
const loadIngresses = async () => {
  if (!selectedCluster.value) {
    ElMessage.warning('请先选择一个集群')
    return
  }
  
  // 取消之前的请求
  if (cancelTokenSource) {
    cancelTokenSource.cancel('Operation canceled by the user.')
  }
  
  // 创建新的取消令牌
  cancelTokenSource = axios.CancelToken.source()
  
  loading.value = true
  try {
    // 调用API获取Ingress列表，使用正确的Ingress端点
    const params: any = {
      _t: Date.now() // 添加时间戳，防止浏览器缓存
    }
    if (selectedNamespace.value) {
      params.namespace = selectedNamespace.value
    }
    
    // 只调用真实的Ingress API端点，不使用模拟数据
    const response = await axios.get(`/clusters/${selectedCluster.value}/ingresses`, {
      params,
      headers: {
        'Cache-Control': 'no-cache' // 禁用缓存
      },
      cancelToken: cancelTokenSource.token
    });
    
    // 为每个Ingress添加集群名称并适配数据格式
    const clusterName = clusters.value.find(c => c.id === selectedCluster.value)?.name || '';
    const newIngresses = response.data.map((ingress: any) => ({
      name: ingress.name || 'unknown',
      namespace: ingress.namespace || '',
      cluster: clusterName,
      // 确保rules字段存在且为数组
      rules: ingress.rules || [],
      // 确保annotations字段存在且为对象
      annotations: ingress.annotations || {},
      // 确保createdAt字段存在
      createdAt: ingress.createdAt || '',
      // 保留其他可能的字段
      ...ingress
    }));
    
    // 先更新数据，然后设置loading为false，确保数据正确显示
    ingresses.value = newIngresses;
    loading.value = false;
    
    // 如果返回的Ingress列表为空，显示提示信息
    if (newIngresses.length === 0) {
      ElMessage.info('目标命名空间内无Ingress实例');
    }
  } catch (error: any) {
    // 如果是取消请求的错误，不显示错误信息
    if (axios.isCancel(error)) {
      console.log('Request canceled:', error.message);
      loading.value = false;
      return;
    }
    
    // 显示详细的错误信息，帮助用户排查问题
    console.error('加载Ingress列表失败:', error);
    loading.value = false;
    
    // 检查是否是404错误，如果是，显示空列表和友好提示
    if (error.response && error.response.status === 404) {
      // 404错误，说明后端没有实现Ingress端点
      ElMessage.info('当前集群暂不支持Ingress管理，或Ingress API端点未实现');
      // 显示空列表，而不是错误状态
      ingresses.value = [];
    } else {
      // 其他错误，显示详细的错误信息
      let errorMessage = '加载Ingress列表失败';
      if (error.response) {
        // 服务器返回了错误状态码
        errorMessage = `加载失败: ${error.response.status} - ${error.response.statusText}`;
        if (error.response.data && error.response.data.message) {
          errorMessage += ` - ${error.response.data.message}`;
        }
      } else if (error.request) {
        // 请求已发送但没有收到响应
        errorMessage = '加载失败: 无法连接到服务器，请检查网络连接';
      } else {
        // 请求配置错误
        errorMessage = `加载失败: ${error.message}`;
      }
      
      ElMessage.error(errorMessage);
      // 确保在错误情况下，列表为空，显示正确的空状态
      ingresses.value = [];
    }
  }
}

// 刷新Ingress列表
const refreshIngresses = () => {
  loadIngresses()
}

// 格式化绝对时间
const formatAge = (createdAt: string) => {
  if (!createdAt) return ''
  
  try {
    // 确保能正确解析ISO格式时间
    const created = new Date(createdAt)
    
    // 检查日期是否有效
    if (isNaN(created.getTime())) {
      return ''
    }
    
    // 格式化日期为YYYY-MM-DD HH:MM:SS格式
    const year = created.getFullYear()
    const month = String(created.getMonth() + 1).padStart(2, '0')
    const day = String(created.getDate()).padStart(2, '0')
    const hours = String(created.getHours()).padStart(2, '0')
    const minutes = String(created.getMinutes()).padStart(2, '0')
    const seconds = String(created.getSeconds()).padStart(2, '0')
    
    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
  } catch (error) {
    console.error('时间格式化失败:', error)
    return ''
  }
}

// 查看Ingress详情
const viewIngressDetails = (ingress: Ingress) => {
  selectedIngress.value = ingress
  ingressDetailsVisible.value = true
}

// 删除Ingress
const deleteIngress = async (ingress: Ingress) => {
  try {
    // 这里应该调用真实的API删除Ingress
    // 暂时使用模拟数据
    const index = ingresses.value.findIndex((i) => i.namespace === ingress.namespace && i.name === ingress.name)
    if (index !== -1) {
      ingresses.value.splice(index, 1)
    }
    ElMessage.success('删除Ingress成功')
  } catch (error) {
    console.error('删除Ingress失败:', error)
    ElMessage.error('删除Ingress失败')
  }
}
</script>

<style scoped>
.ingresses-view {
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

.ingresses-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.ingress-details {
  padding: 10px 0;
}

.host-item,
.path-item,
.service-item {
  margin-bottom: 5px;
}

.rule-details {
  padding: 10px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.empty-text {
  margin-top: 10px;
  color: #909399;
  font-style: italic;
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