<template>
  <div class="storageclasses-view">
    <el-card shadow="never" class="page-header">
      <div class="header-content">
        <h2>存储类</h2>
        <div class="header-actions">
          <el-select v-model="selectedCluster" placeholder="选择集群" style="width: 200px; margin-right: 10px;">
            <el-option
              v-for="cluster in clusters"
              :key="cluster.id"
              :label="cluster.name"
              :value="cluster.id"
            />
          </el-select>
          <el-button type="primary" size="small" @click="refreshStorageClasses">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </div>
    </el-card>

    <el-card shadow="hover" class="storageclasses-card">
      <template #header>
        <div class="card-header">
          <span>K8s StorageClass列表</span>
        </div>
      </template>
      
      <el-table 
        :data="storageClasses" 
        style="width: 100%; color: #303133; font-size: 14px; line-height: 1.5;"
        stripe
        v-loading="loading"
        v-if="selectedCluster"
        :row-style="{ height: 'auto', backgroundColor: '#ffffff' }"
        :cell-style="{ 
          'white-space': 'normal', 
          'word-break': 'break-all', 
          'min-height': '40px', 
          'padding': '8px 12px', 
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
        border
        empty-text="No Data"
        :row-key="(row: any) => row.name"
      >
        <el-table-column prop="name" label="存储类名称" min-width="200" show-overflow-tooltip />
        <el-table-column prop="provisioner" label="供应商" width="180" show-overflow-tooltip />
        <el-table-column prop="reclaimPolicy" label="回收策略" width="120">
          <template #default="scope">
            <el-tag :type="getReclaimPolicyTagType(scope.row.reclaimPolicy)" size="medium">
              {{ scope.row.reclaimPolicy }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="volumeBindingMode" label="卷绑定模式" width="150" show-overflow-tooltip />
        <el-table-column prop="allowVolumeExpansion" label="允许卷扩展" width="120">
          <template #default="scope">
            <el-tag :type="scope.row.allowVolumeExpansion ? 'success' : 'danger'" size="medium">
              {{ scope.row.allowVolumeExpansion ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="isDefault" label="默认存储类" width="120">
          <template #default="scope">
            <el-tag :type="scope.row.isDefault ? 'success' : 'danger'" size="medium">
              {{ scope.row.isDefault ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="180" :cell-style="{ 'white-space': 'nowrap' }">
          <template #default="scope">
            {{ formatAge(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" @click="viewStorageClassDetails(scope.row)">
              <el-icon><View /></el-icon>
              查看
            </el-button>
            <el-button type="danger" size="small" @click="deleteStorageClass(scope.row)">
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <div v-else class="no-cluster-selected">
        <el-empty description="请选择一个集群查看存储类列表" />
      </div>
    </el-card>
    
    <!-- 存储类详情对话框 -->
    <el-dialog
      v-model="storageClassDetailsVisible"
      title="存储类详情"
      width="800px"
    >
      <div v-if="selectedStorageClass" class="storageclass-details">
        <el-tabs v-model="activeTab">
          <el-tab-pane label="详情" name="details">
            <el-descriptions :column="2" border>
              <el-descriptions-item label="存储类名称">{{ selectedStorageClass.name }}</el-descriptions-item>
              <el-descriptions-item label="供应商">{{ selectedStorageClass.provisioner }}</el-descriptions-item>
              <el-descriptions-item label="回收策略">{{ selectedStorageClass.reclaimPolicy }}</el-descriptions-item>
              <el-descriptions-item label="卷绑定模式">{{ selectedStorageClass.volumeBindingMode }}</el-descriptions-item>
              <el-descriptions-item label="允许卷扩展">{{ selectedStorageClass.allowVolumeExpansion ? '是' : '否' }}</el-descriptions-item>
              <el-descriptions-item label="默认存储类">{{ selectedStorageClass.isDefault ? '是' : '否' }}</el-descriptions-item>
              <el-descriptions-item label="创建时间">{{ selectedStorageClass.createdAt }}</el-descriptions-item>
            </el-descriptions>
            
            <h4 style="margin-top: 20px;">Labels</h4>
            <div class="labels-section">
              <el-tag v-if="Object.keys(selectedStorageClass.labels).length === 0" type="info" size="small">
                无
              </el-tag>
              <el-tag
                v-for="(value, key) in selectedStorageClass.labels"
                :key="key"
                type="info"
                size="small"
                style="margin-right: 8px; margin-bottom: 8px;"
              >
                {{ key }}={{ value }}
              </el-tag>
            </div>
            
            <h4 style="margin-top: 20px;">Annotations</h4>
            <div class="annotations-section">
              <el-tag v-if="Object.keys(selectedStorageClass.annotations).length === 0" type="info" size="small">
                无
              </el-tag>
              <el-tag
                v-for="(value, key) in selectedStorageClass.annotations"
                :key="key"
                type="success"
                size="small"
                style="margin-right: 8px; margin-bottom: 8px;"
              >
                {{ key }}={{ value }}
              </el-tag>
            </div>
          </el-tab-pane>
          <el-tab-pane label="YAML" name="yaml">
            <div class="yaml-section">
              <div class="yaml-actions" style="margin-bottom: 10px;">
                <el-button type="primary" size="small" @click="copyYaml">
                  <el-icon><DocumentCopy /></el-icon>
                  复制YAML
                </el-button>
              </div>
              <el-scrollbar height="500px">
                <pre v-if="yamlContent" class="yaml-code">{{ yamlContent }}</pre>
                <div v-else class="loading-yaml">
                  <el-skeleton :rows="10" animated />
                </div>
              </el-scrollbar>
            </div>
          </el-tab-pane>
        </el-tabs>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="storageClassDetailsVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { Refresh, View, Delete, DocumentCopy } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import axios from 'axios'

// 加载状态
const loading = ref(false)

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

// 选择的集群
const selectedCluster = ref<number | undefined>(undefined)

// 存储类数据类型定义
interface StorageClass {
  name: string
  cluster: string
  provisioner: string
  reclaimPolicy: string
  volumeBindingMode: string
  allowVolumeExpansion: boolean
  isDefault: boolean
  labels: Record<string, string>
  annotations: Record<string, string>
  createdAt: string
}

// 存储类数据
const storageClasses = ref<StorageClass[]>([])

// 对话框相关
const storageClassDetailsVisible = ref(false)
const selectedStorageClass = ref<StorageClass | null>(null)
const activeTab = ref('details')
const yamlContent = ref('')

// 组件挂载时加载数据
onMounted(() => {
  loadClusters()
})

// 监听集群选择变化，自动刷新存储类列表
watch(selectedCluster, () => {
  if (selectedCluster.value) {
    loadStorageClasses()
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

// 加载存储类列表
const loadStorageClasses = async () => {
  if (!selectedCluster.value) {
    ElMessage.warning('请先选择一个集群')
    return
  }
  
  loading.value = true
  try {
    // 调用API获取存储类列表，axios会自动添加baseURL
    const response = await axios.get(`/clusters/${selectedCluster.value}/storageclasses`, {
      params: {
        _t: Date.now() // 添加时间戳，防止浏览器缓存
      },
      headers: {
        'Cache-Control': 'no-cache' // 禁用缓存
      }
    })
    // 为每个存储类添加集群名称
    const clusterName = clusters.value.find(c => c.id === selectedCluster.value)?.name || ''
    const newStorageClasses = response.data.map((sc: any) => ({
      ...sc,
      cluster: clusterName
    }))
    
    // 先更新数据，然后设置loading为false，确保数据正确显示
    storageClasses.value = newStorageClasses
    loading.value = false
    
    // 如果返回的存储类列表为空，显示提示信息
    if (newStorageClasses.length === 0) {
      ElMessage.info('当前集群内无存储类实例')
    }
  } catch (error) {
    console.error('加载存储类列表失败:', error)
    loading.value = false
    ElMessage.error('加载存储类列表失败')
    // 确保在错误情况下，列表为空，显示正确的空状态
    storageClasses.value = []
  }
}

// 刷新存储类列表
const refreshStorageClasses = () => {
  loadStorageClasses()
}

// 获取回收策略标签类型
const getReclaimPolicyTagType = (policy: string) => {
  switch (policy) {
    case 'Retain':
      return 'warning'
    case 'Delete':
      return 'success'
    case 'Recycle':
      return 'info'
    default:
      return 'info'
  }
}

// 格式化时间
const formatAge = (createdAt: string) => {
  if (!createdAt) return ''
  
  try {
    const created = new Date(createdAt)
    
    // 格式化时间为YYYY-MM-DD HH:MM:SS
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

// 查看存储类详情
const viewStorageClassDetails = async (storageClass: StorageClass) => {
  selectedStorageClass.value = storageClass
  storageClassDetailsVisible.value = true
  activeTab.value = 'details'
  
  // 加载YAML内容
  await loadYamlContent(storageClass.name)
}

// 加载YAML内容
const loadYamlContent = async (name: string) => {
  if (!selectedCluster.value) return
  
  try {
    // 调用API获取存储类YAML内容
    const response = await axios.get(`/clusters/${selectedCluster.value}/storageclasses/${name}/yaml`)
    yamlContent.value = response.data
  } catch (error) {
    console.error('加载YAML内容失败:', error)
    yamlContent.value = ''
    ElMessage.error('加载YAML内容失败')
  }
}

// 复制YAML内容
const copyYaml = () => {
  if (!yamlContent.value) {
    ElMessage.warning('YAML内容为空，无法复制')
    return
  }
  
  // 使用浏览器的Clipboard API复制内容
  navigator.clipboard.writeText(yamlContent.value)
    .then(() => {
      ElMessage.success('YAML内容已复制到剪贴板')
    })
    .catch(err => {
      console.error('复制失败:', err)
      ElMessage.error('复制失败，请手动复制')
    })
}

// 删除存储类
const deleteStorageClass = async (storageClass: StorageClass) => {
  // 显示确认对话框
  const confirmResult = await ElMessageBox.confirm(
    `确定要删除存储类 "${storageClass.name}" 吗？此操作无法撤销。`,
    '删除确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).catch(() => {
    // 用户取消删除
    return 'cancel'
  })
  
  if (confirmResult !== 'confirm') {
    return
  }
  
  try {
    // 调用API删除存储类
    await axios.delete(`/clusters/${selectedCluster.value}/storageclasses/${storageClass.name}`)
    
    // 从列表中移除该存储类
    const index = storageClasses.value.findIndex((sc) => sc.name === storageClass.name)
    if (index !== -1) {
      storageClasses.value.splice(index, 1)
    }
    
    ElMessage.success('删除存储类成功')
  } catch (error) {
    console.error('删除存储类失败:', error)
    ElMessage.error('删除存储类失败')
  }
}
</script>

<style scoped>
.storageclasses-view {
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

.storageclasses-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.storageclass-details {
  padding: 10px 0;
}

.labels-section,
.annotations-section {
  margin-top: 10px;
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

/* YAML样式 */
.yaml-section {
  margin-top: 10px;
}

.yaml-actions {
  display: flex;
  justify-content: flex-end;
}

.yaml-code {
  background-color: #f6f8fa;
  border: 1px solid #e1e4e8;
  border-radius: 6px;
  padding: 16px;
  font-family: 'Courier New', Courier, monospace;
  font-size: 14px;
  line-height: 1.5;
  overflow-x: auto;
  color: #24292e;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.loading-yaml {
  padding: 16px;
}
</style>