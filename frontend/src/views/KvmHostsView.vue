<template>
  <div class="kvm-hosts-view">
    <el-card shadow="never" class="page-header">
      <div class="header-content">
        <h2>KVM管理</h2>
      </div>
    </el-card>

    <!-- KVM主机列表 -->
    <el-card shadow="hover" class="config-card" v-loading="loading.hosts" element-loading-text="加载KVM主机列表中...">
      <template #header>
        <div class="card-header">
          <span>KVM主机列表</span>
          <el-button type="primary" @click="showAddKvmHostDialog = true">
            <el-icon><Plus /></el-icon>
            添加KVM主机
          </el-button>
        </div>
      </template>
      
      <el-table :data="kvmHosts" stripe style="width: 100%" fit>
        <el-table-column prop="id" label="ID" width="80" min-width="80"></el-table-column>
        <el-table-column prop="name" label="主机名称" min-width="150"></el-table-column>
        <el-table-column prop="ip" label="IP地址" min-width="180"></el-table-column>
        <el-table-column prop="port" label="端口" width="100" min-width="100"></el-table-column>
        <el-table-column prop="connectionType" label="连接方式" width="120" min-width="120">
          <template #default="scope">
            <el-tag :type="scope.row.connectionType === 'ssh' ? 'success' : 'primary'">
              {{ scope.row.connectionType === 'ssh' ? 'SSH' : 'virsh' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="120" min-width="120">
          <template #default="scope">
            <el-tag :type="(scope.row.status === 'active' || scope.row.status === 'running') ? 'success' : 'danger'">
              {{ (scope.row.status === 'active' || scope.row.status === 'running') ? '在线' : '离线' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="180" min-width="180">
          <template #default="scope">
            {{ formatDateTime(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="240" min-width="240" fixed="right">
          <template #default="scope">
            <el-space size="small" direction="horizontal">
              <el-button type="primary" size="small" @click="viewVmList(scope.row)">
                <el-icon><View /></el-icon>
                虚拟机列表
              </el-button>
              <el-dropdown>
                <el-button type="info" size="small">
                  更多 <el-icon class="el-icon--right"><ArrowDown /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item @click="openEditKvmHostDialog(scope.row)">
                      <el-icon><Edit /></el-icon>
                      编辑
                    </el-dropdown-item>
                    <el-dropdown-item @click="testConnection(scope.row)">
                      <el-icon><Connection /></el-icon>
                      测试连接
                    </el-dropdown-item>
                    <el-dropdown-item divided @click="deleteKvmHost(scope.row)">
                      <el-icon><Delete /></el-icon>
                      删除
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </el-space>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 虚拟机列表弹窗 -->
    <el-dialog
      v-model="showVmListDialog"
      :title="`虚拟机列表 - ${selectedKvmHost.name}`"
      width="80%"
    >
      <el-table :data="vmList" stripe style="width: 100%" v-loading="loading.vms" element-loading-text="加载虚拟机列表中..." :align="'center'" :header-align="'center'" fit>
        <el-table-column prop="name" label="虚拟机名称" min-width="150" :header-align="'left'" :align="'left'">
          <template #default="scope">
            <div class="vm-name-cell">
              {{ scope.row.name }}
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" min-width="80">
          <template #default="scope">
            <el-tag :type="getStatusTagType(scope.row.status)" size="small" effect="light">
              {{ scope.row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="cpu" label="CPU" min-width="60">
          <template #default="scope">
            <span class="vm-metric">{{ scope.row.cpu }} vCPU</span>
          </template>
        </el-table-column>
        <el-table-column prop="memory" label="内存" min-width="90">
          <template #default="scope">
            <span class="vm-metric">{{ scope.row.memory }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="ip" label="IP地址" min-width="120">
          <template #default="scope">
            <el-tag v-if="scope.row.ip" size="small" type="info" effect="light">
              {{ scope.row.ip }}
            </el-tag>
            <span v-else class="no-ip">-</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" min-width="200">
          <template #default="scope">
            <el-space size="mini">
              <el-button type="primary" size="small" @click="startVm(scope.row)" :disabled="scope.row.status === 'running'">
                <el-icon><VideoPlay /></el-icon>
                启动
              </el-button>
              <el-button type="warning" size="small" @click="shutdownVm(scope.row)" :disabled="scope.row.status !== 'running'">
                <el-icon><VideoPause /></el-icon>
                关机
              </el-button>
              <el-button type="danger" size="small" @click="forceShutdownVm(scope.row)" :disabled="scope.row.status !== 'running'">
                <el-icon><CircleClose /></el-icon>
                强制关机
              </el-button>
              <el-button type="info" size="small" @click="rebootVm(scope.row)" :disabled="scope.row.status !== 'running'">
                <el-icon><Refresh /></el-icon>
                重启
              </el-button>
            </el-space>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>

    <!-- 添加/编辑KVM主机弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="currentEditHostId > 0 ? '编辑KVM主机' : '添加KVM主机'"
      width="600px"
      @close="resetDialog"
    >
      <el-form ref="kvmHostFormRef" :model="kvmHostForm" :rules="kvmHostRules" label-width="140px" class="kvm-host-form">
        <el-form-item label="主机名称" prop="name">
          <el-input v-model="kvmHostForm.name" placeholder="请输入KVM主机名称" clearable />
        </el-form-item>

        <el-form-item label="连接方式" prop="connectionType">
          <el-radio-group v-model="kvmHostForm.connectionType">
            <el-radio label="ssh">SSH连接</el-radio>
            <el-radio label="virsh">virsh直接连接</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="IP地址" prop="ip">
          <el-input v-model="kvmHostForm.ip" placeholder="请输入KVM主机IP地址" clearable />
        </el-form-item>

        <el-form-item label="端口" prop="port">
          <el-input-number v-model="kvmHostForm.port" :min="1" :max="65535" placeholder="请输入端口" />
        </el-form-item>

        <!-- SSH连接方式的额外字段 -->
        <el-form-item v-if="kvmHostForm.connectionType === 'ssh'" label="用户名" prop="username">
          <el-input v-model="kvmHostForm.username" placeholder="请输入SSH用户名" clearable />
        </el-form-item>

        <el-form-item v-if="kvmHostForm.connectionType === 'ssh'" label="密码" prop="password">
          <el-input v-model="kvmHostForm.password" type="password" placeholder="请输入SSH密码（留空则不修改）" show-password clearable />
        </el-form-item>

        <el-form-item label="描述" prop="description">
          <el-input v-model="kvmHostForm.description" type="textarea" :rows="3" placeholder="请输入描述" clearable />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="saveKvmHost" :loading="loading.save">保存</el-button>
          <el-button @click="resetKvmHostForm">重置</el-button>
        </el-form-item>
      </el-form>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watchEffect } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, View, Connection, Delete, VideoPlay, VideoPause, CircleClose, Refresh, Edit, ArrowDown } from '@element-plus/icons-vue'
import { kvmAPI } from '../api/api'

// 加载状态
const loading = ref({
  hosts: false,
  vms: false,
  test: false,
  save: false
})

// KVM主机列表
const kvmHosts = ref<any[]>([])

// 虚拟机列表
const vmList = ref<any[]>([])

// 选中的KVM主机
const selectedKvmHost = ref({
  id: 0,
  name: ''
})

// 弹窗状态
const showAddKvmHostDialog = ref(false)
const showVmListDialog = ref(false)
const showEditKvmHostDialog = ref(false)
const dialogVisible = ref(false)

// 当前编辑的KVM主机ID
const currentEditHostId = ref(0)

// 表单引用
const kvmHostFormRef = ref()

// KVM主机表单
const kvmHostForm = reactive({
  name: '',
  ip: '',
  port: 22,
  connectionType: 'ssh', // 默认连接方式为SSH
  username: '',
  password: '',
  description: ''
})

// KVM主机表单验证规则
const kvmHostRules = ref({
  name: [
    { required: true, message: '请输入KVM主机名称', trigger: 'blur' },
    { min: 2, max: 50, message: '主机名称长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  ip: [
    { required: true, message: '请输入KVM主机IP地址', trigger: 'blur' },
    { pattern: /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/, message: '请输入正确的IP地址', trigger: 'blur' }
  ],
  port: [
    { required: true, message: '请输入端口', trigger: 'blur' },
    { type: 'number', min: 1, max: 65535, message: '端口范围在 1 到 65535 之间', trigger: 'blur' }
  ],
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur', validator: (_rule: any, value: any, callback: any) => {
      if (kvmHostForm.connectionType === 'ssh' && !value) {
        callback(new Error('请输入用户名'))
      } else {
        callback(null)
      }
    } }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur', validator: (_rule: any, value: any, callback: any) => {
      if (kvmHostForm.connectionType === 'ssh' && !value) {
        callback(new Error('请输入密码'))
      } else {
        callback(null)
      }
    } }
  ]
})

// 监听对话框状态变化
watchEffect(() => {
  dialogVisible.value = showAddKvmHostDialog.value || showEditKvmHostDialog.value
})

// 重置对话框状态
const resetDialog = () => {
  showAddKvmHostDialog.value = false
  showEditKvmHostDialog.value = false
  currentEditHostId.value = 0
  // 重置表单
  Object.assign(kvmHostForm, {
    name: '',
    ip: '',
    port: 22,
    connectionType: 'ssh',
    username: '',
    password: '',
    description: ''
  })
}

// 格式化日期时间为 YYYY-MM-DD HH:mm:ss 格式
const formatDateTime = (dateTime: string | Date) => {
  const date = typeof dateTime === 'string' ? new Date(dateTime) : dateTime
  
  // 确保日期有效
  if (isNaN(date.getTime())) {
    return ''
  }
  
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

// 生命周期钩子
onMounted(() => {
  // 加载KVM主机列表
  loadKvmHosts()
})

// 加载KVM主机列表
const loadKvmHosts = async () => {
  try {
    loading.value.hosts = true
    const response = await kvmAPI.list()
    kvmHosts.value = response.data || []
  } catch (error) {
    console.error('加载KVM主机列表失败:', error)
    kvmHosts.value = []
  } finally {
    loading.value.hosts = false
  }
}

// 查看虚拟机列表
const viewVmList = async (host: any) => {
  selectedKvmHost.value = host
  showVmListDialog.value = true
  try {
    loading.value.vms = true
    const response = await kvmAPI.getVmList(host.id)
    // 处理后端返回的数据格式，axios会将响应数据包装在response.data中
    vmList.value = response.data || []
  } catch (error) {
    console.error('获取虚拟机列表失败:', error)
    vmList.value = []
  } finally {
    loading.value.vms = false
  }
}

// 显示编辑KVM主机弹窗
const openEditKvmHostDialog = (host: any) => {
  // 填充表单数据
  Object.assign(kvmHostForm, {
    name: host.name,
    ip: host.ip,
    port: host.port,
    connectionType: host.connectionType,
    username: host.username,
    password: '', // 密码不显示在编辑表单中
    description: host.description
  })
  currentEditHostId.value = host.id
  showEditKvmHostDialog.value = true
}

// 测试连接
const testConnection = async (host: any) => {
  try {
    loading.value.test = true
    await kvmAPI.testConnection(host.id)
    ElMessage.success('测试连接成功')
  } catch (error: any) {
    console.error('测试连接失败:', error)
    // 显示更详细的错误信息
    const errorMsg = error.response?.data?.error || '测试连接失败'
    ElMessage.error(errorMsg)
  } finally {
    loading.value.test = false
    // 无论测试成功还是失败，都重新加载主机列表，更新状态
    await loadKvmHosts()
  }
}

// 删除KVM主机
const deleteKvmHost = async (host: any) => {
  try {
    await ElMessageBox.confirm('确定要删除该KVM主机吗？', '删除确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await kvmAPI.delete(host.id)
    ElMessage.success('删除KVM主机成功')
    loadKvmHosts()
  } catch (error) {
    console.error('删除KVM主机失败:', error)
    if (error !== 'cancel') {
      ElMessage.error('删除KVM主机失败')
    }
  }
}

// 获取状态标签类型
const getStatusTagType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'running':
      return 'success'
    case 'stopped':
      return 'danger'
    case 'paused':
      return 'warning'
    default:
      return 'info'
  }
}

// 启动虚拟机
const startVm = async (vm: any) => {
  try {
    await kvmAPI.startVm(selectedKvmHost.value.id, vm.name)
    ElMessage.success(`启动虚拟机${vm.name}成功`)
    // 重新加载虚拟机列表
    await viewVmList(selectedKvmHost.value)
  } catch (error) {
    console.error('启动虚拟机失败:', error)
    ElMessage.error(`启动虚拟机${vm.name}失败`)
  }
}

// 关闭虚拟机
const shutdownVm = async (vm: any) => {
  try {
    await kvmAPI.shutdownVm(selectedKvmHost.value.id, vm.name)
    ElMessage.success(`关闭虚拟机${vm.name}成功`)
    // 重新加载虚拟机列表
    await viewVmList(selectedKvmHost.value)
  } catch (error) {
    console.error('关闭虚拟机失败:', error)
    ElMessage.error(`关闭虚拟机${vm.name}失败`)
  }
}

// 重启虚拟机
const rebootVm = async (vm: any) => {
  try {
    await kvmAPI.rebootVm(selectedKvmHost.value.id, vm.name)
    ElMessage.success(`重启虚拟机${vm.name}成功`)
    // 重新加载虚拟机列表
    await viewVmList(selectedKvmHost.value)
  } catch (error) {
    console.error('重启虚拟机失败:', error)
    ElMessage.error(`重启虚拟机${vm.name}失败`)
  }
}

// 强制关机虚拟机
const forceShutdownVm = async (vm: any) => {
  try {
    // 显示确认对话框
    await ElMessageBox.confirm(`确定要强制关机虚拟机${vm.name}吗？这个操作可能会导致数据丢失！`, {
      title: '强制关机确认',
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'error'
    })
    await kvmAPI.forceShutdownVm(selectedKvmHost.value.id, vm.name)
    ElMessage.success(`强制关机虚拟机${vm.name}成功`)
    // 重新加载虚拟机列表
    await viewVmList(selectedKvmHost.value)
  } catch (error) {
    console.error('强制关机虚拟机失败:', error)
    if (error !== 'cancel') {
      ElMessage.error(`强制关机虚拟机${vm.name}失败`)
    }
  }
}

// 保存KVM主机（创建或编辑）
const saveKvmHost = async () => {
  if (!kvmHostFormRef.value) return
  
  try {
    await kvmHostFormRef.value.validate()
    loading.value.save = true
    console.log('开始调用API，表单数据:', kvmHostForm)
    
    let response
    let successMsg
    
    if (currentEditHostId.value > 0) {
      // 编辑操作
      response = await kvmAPI.update(currentEditHostId.value, kvmHostForm)
      successMsg = '编辑KVM主机成功'
      showEditKvmHostDialog.value = false
    } else {
      // 创建操作
      response = await kvmAPI.create(kvmHostForm)
      successMsg = '添加KVM主机成功'
      showAddKvmHostDialog.value = false
    }
    
    console.log('API调用成功，响应:', response)
    ElMessage.success(successMsg)
    resetKvmHostForm()
    loadKvmHosts()
  } catch (error) {
    console.error('保存KVM主机失败:', error)
    // 解析错误信息，提供更详细的提示
    let errorMsg = '保存KVM主机失败'
    if (typeof error === 'object' && error !== null) {
      if ('response' in error) {
        // axios错误
        const axiosError = error as any
        if (axiosError.response) {
          // 服务器返回了错误响应
          errorMsg = `保存KVM主机失败: ${axiosError.response.data?.message || axiosError.response.statusText || '未知错误'}`
          console.error('服务器错误响应:', axiosError.response.data)
        } else if (axiosError.request) {
          // 请求已发送但没有收到响应
          errorMsg = '保存KVM主机失败: 服务器无响应，请检查网络连接'
        } else {
          // 请求配置错误
          errorMsg = `保存KVM主机失败: ${axiosError.message}`
        }
      } else if ('message' in error) {
        // 其他类型的错误
        errorMsg = `保存KVM主机失败: ${(error as Error).message}`
      }
    }
    ElMessage.error(errorMsg)
  } finally {
    loading.value.save = false
  }
}

// 重置KVM主机表单
const resetKvmHostForm = () => {
  // 重置表单数据
  Object.assign(kvmHostForm, {
    name: '',
    ip: '',
    port: 22,
    connectionType: 'ssh',
    username: '',
    password: '',
    description: ''
  })
  
  // 重置编辑状态
  currentEditHostId.value = 0
  
  if (kvmHostFormRef.value) {
    kvmHostFormRef.value.resetFields()
  }
}
</script>

<style scoped>
.kvm-hosts-view {
  width: 100%;
  padding: 0;
}

.page-header {
  margin-bottom: 30px;
}

.header-content h2 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: #333;
}

.config-card {
  margin-bottom: 30px;
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.kvm-host-form {
  width: 100%;
}

.kvm-host-form .el-form-item {
  margin-bottom: 20px;
}

/* 虚拟机列表样式 */
:deep(.el-table__body tr:hover) {
  background-color: #f5f7fa;
}

.vm-name-cell {
  font-weight: 500;
  color: #333;
}

.vm-metric {
  color: #606266;
  font-size: 14px;
}

.no-ip {
  color: #909399;
  font-size: 12px;
}

:deep(.el-tag) {
  margin-right: 0;
}

:deep(.el-button--small) {
  padding: 4px 10px;
  font-size: 12px;
}

:deep(.el-button[disabled]) {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>