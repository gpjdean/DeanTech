<template>
  <div class="users-view">
    <el-card shadow="never" class="page-header">
      <div class="header-content">
        <h2>用户管理</h2>
      </div>
    </el-card>

    <!-- 用户管理内容 -->
    <div class="users-section">
      <div class="users-header">
        <h3>用户列表</h3>
        <el-button type="primary" @click="openAddUserDialog">
          <el-icon><Plus /></el-icon>
          添加用户
        </el-button>
      </div>

      <el-card shadow="hover" class="users-card" v-loading="loading.users" element-loading-text="加载用户列表中...">
        <el-table :data="users" style="width: 100%" stripe fit>
          <el-table-column prop="id" label="ID" width="80" min-width="80" />
          <el-table-column prop="username" label="用户名" min-width="150" />
          <el-table-column prop="name" label="姓名" min-width="120" />
          <el-table-column prop="email" label="邮箱" min-width="200" />
          <el-table-column prop="phone" label="电话" min-width="120" />
          <el-table-column prop="department" label="部门" min-width="120" />
          <el-table-column prop="position" label="职位" min-width="120" />
          <el-table-column prop="createdAt" label="创建时间" min-width="180">
            <template #default="scope">
              {{ formatDate(scope.row.createdAt) }}
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="120">
            <template #default="scope">
              <div class="status-switch-container">
                <el-switch 
                  :model-value="scope.row.status" 
                  active-text="启用" 
                  inactive-text="禁用"
                  active-color="#67c23a" 
                  inactive-color="#f56c6c"
                  @update:model-value="updateUserStatus(scope.row, $event)"
                  size="large"
                  class="custom-status-switch"
                />
              </div>
            </template>
          </el-table-column>
          <el-table-column label="操作" min-width="180" fixed="right">
            <template #default="scope">
              <el-button type="primary" size="small" @click="editUser(scope.row)">
                <el-icon><Edit /></el-icon>
                编辑
              </el-button>
              <el-button type="danger" size="small" @click="deleteUser(scope.row.id)">
                <el-icon><Delete /></el-icon>
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </div>

    <!-- 添加/编辑用户对话框 -->
    <el-dialog
      v-model="userDialogVisible"
      :title="isEditUser ? '编辑用户' : '添加用户'"
      width="500px"
    >
      <el-form ref="userForm" :model="userFormData" :rules="userRules" label-width="100px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="userFormData.username" placeholder="请输入用户名" :disabled="isEditUser" /></el-form-item>

        <el-form-item label="密码" prop="password">
          <el-input v-model="userFormData.password" type="password" placeholder="如需修改密码请输入，否则留空" show-password /></el-form-item>

        <el-form-item label="姓名" prop="name">
          <el-input v-model="userFormData.name" placeholder="请输入姓名" /></el-form-item>

        <el-form-item label="邮箱" prop="email">
          <el-input v-model="userFormData.email" placeholder="请输入邮箱" /></el-form-item>

        <el-form-item label="电话" prop="phone">
          <el-input v-model="userFormData.phone" placeholder="请输入电话" /></el-form-item>

        <el-form-item label="部门" prop="department">
          <el-input v-model="userFormData.department" placeholder="请输入部门" /></el-form-item>

        <el-form-item label="职位" prop="position">
          <el-input v-model="userFormData.position" placeholder="请输入职位" /></el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="userDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveUser" :loading="loading.saveUser">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus, Edit, Delete } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { userAPI } from '../api/api'



// 表单引用
const userForm = ref()

// 加载状态
const loading = ref({
  users: false,
  saveUser: false
})

// 用户列表数据
interface User {
  id: number;
  username: string;
  email: string;
  name: string;
  phone: string;
  department: string;
  position: string;
  createdAt: string;
  updatedAt: string;
  deletedAt: string | null;
  status: boolean;
  rawStatus: string;
  role: string;
}

const users = ref<User[]>([])

// 用户表单验证规则
const userRules = ref({
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    {
      validator: (_rule: any, value: string, callback: Function) => {
        if (!isEditUser.value && !value) {
          callback(new Error('请输入密码'))
        } else if (value && (value.length < 6 || value.length > 20)) {
          callback(new Error('密码长度在 6 到 20 个字符'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  name: [
    { required: true, message: '请输入姓名', trigger: 'blur' },
    { min: 2, max: 20, message: '姓名长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur, change' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
  ]
})

// 用户对话框相关
const userDialogVisible = ref(false)
const isEditUser = ref(false)
const userFormData = ref({
  id: 0,
  username: '',
  password: '',
  name: '',
  email: '',
  phone: '',
  department: '',
  position: '',
  status: false
})


// 组件挂载时加载用户列表
onMounted(() => {
  // 加载用户列表
  loadUsers()
})


// 加载用户列表
const loadUsers = async () => {
  try {
    loading.value.users = true
    const usersData = await userAPI.list()
    // 调试信息：打印从后端获取的原始用户数据
    console.log('从后端获取的用户数据:', usersData)
    // 将后端返回的字段转换为统一格式，支持大小写字段名
      users.value = usersData.map((user: any) => {
        // 注意：API返回的是小写的status，而我们发送的是大写的Status
        const userStatus = user.status || user.Status;
        console.log(`用户${user.username}的原始状态:`, userStatus);
        return {
          id: user.id || user.ID,
          username: user.username || user.Username,
          email: user.email || user.Email,
          name: user.name || user.Name,
          phone: user.phone || user.Phone,
          department: user.department || user.Department,
          position: user.position || user.Position,
          createdAt: user.createdAt || user.CreatedAt,
          updatedAt: user.updatedAt || user.UpdatedAt,
          deletedAt: user.deletedAt || user.DeletedAt,
          status: userStatus === 'active', // 只有状态明确为active时，才显示启用
          rawStatus: userStatus, // 保存原始状态用于调试
          role: user.role || user.Role
        };
      })
    // 调试信息：打印转换后的用户数据
    console.log('转换后的用户数据:', users.value)
  } catch (error: any) {
    console.error('加载用户列表失败:', error)
    ElMessage.error('加载用户列表失败')
  } finally {
    loading.value.users = false
  }
}

// 打开添加用户对话框
const openAddUserDialog = () => {
  isEditUser.value = false
  userFormData.value = {
    id: 0,
    username: '',
    password: '',
    name: '',
    email: '',
    phone: '',
    department: '',
    position: '',
    status: false
  }
  userDialogVisible.value = true
}

// 编辑用户
const editUser = (user: any) => {
  isEditUser.value = true
  userFormData.value = {
    ...user,
    password: '' // 编辑时密码字段清空，方便用户输入新密码
  }
  userDialogVisible.value = true
}

// 删除用户
const deleteUser = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除该用户吗？', '删除确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await userAPI.delete(id)
    await loadUsers()
    ElMessage.success('用户删除成功')
  } catch (error) {
    // 取消删除或删除失败
    if (error !== 'cancel') {
      console.error('删除用户失败:', error)
      ElMessage.error('删除用户失败')
    }
  }
}

// 更新用户状态
const updateUserStatus = async (user: User, newStatus: boolean) => {
  // 确保newStatus是boolean类型
  const status = Boolean(newStatus)

  console.log('=== 开始更新用户状态 ===')
  console.log('用户信息:', user)
  console.log('当前状态:', user.status)
  console.log('目标状态:', status)
  
  // 计算API需要的状态值
  const apiStatus = status ? 'active' : 'inactive'
  
  // 先本地更新状态，提供即时反馈
  const userIndex = users.value.findIndex(u => u.id === user.id)
  if (userIndex !== -1) {
    users.value[userIndex].status = status
    users.value[userIndex].rawStatus = apiStatus
  }
  
  try {
    // 保存原始状态，用于API调用失败时恢复
    const originalStatus = user.status
    const originalRawStatus = user.rawStatus
    
    // 调用API更新用户状态，只使用大写Status字段，确保后端能正确处理
    console.log('准备调用API，用户ID:', user.id, '目标状态:', apiStatus)
    const response = await userAPI.update(user.id, {
      Status: apiStatus
    })
    console.log('API调用成功，响应:', response)
    console.log('API调用成功，响应数据:', response.data)
    
    // API调用成功后，显示操作结果
    ElMessage.success(`用户已${status ? '启用' : '禁用'}`)
    
    // 无论是否禁用用户，都触发全局状态检查，确保所有客户端都能立即检查状态
    console.log('触发全局状态检查...')
    window.dispatchEvent(new Event('checkUserStatus'))
    
    // 如果是禁用用户，显示额外提示并触发用户状态变更事件
    if (!status) {
      ElMessage.info('被禁用的用户将立即被强制退出')
      console.log('触发用户状态变更事件，禁用用户ID:', user.id)
      
      // 触发用户状态变更事件，通知所有客户端
      window.dispatchEvent(new CustomEvent('userStatusChanged', {
        detail: {
          userId: user.id,
          status: 'inactive'
        }
      }))
      
      // 使用BroadcastChannel发送跨标签页消息，确保同一浏览器的其他标签页能收到通知
      try {
        const bc = new BroadcastChannel('user-status-channel')
        bc.postMessage({
          type: 'checkUserStatus'
        })
        bc.postMessage({
          type: 'userStatusChanged',
          data: {
            userId: user.id,
            status: 'inactive'
          }
        })
        bc.close()
      } catch (error) {
        console.error('BroadcastChannel发送消息失败:', error)
      }
    }
    
    // API调用成功后，重新加载用户列表，确保状态一致
    console.log('重新加载用户列表...')
    await loadUsers()
    console.log('用户列表加载完成')
    console.log('=== 用户状态更新完成 ===')
  } catch (error: any) {
    console.error('=== 更新用户状态失败 ===')
    console.error('错误信息:', error.message)
    console.error('错误详情:', error.response?.data || error)
    console.error('错误响应:', error.response)
    ElMessage.error('更新用户状态失败')
    
    // 重新加载用户列表，恢复正确状态
    console.log('重新加载用户列表，恢复正确状态...')
    await loadUsers()
    console.log('=== 错误处理完成 ===')
  }
}

// 保存用户
const saveUser = async () => {
  try {
    // 表单验证
    if (!userForm.value) return
    
    const valid = await userForm.value.validate()
    if (!valid) return
    
    loading.value.saveUser = true
    
    // 准备发送的数据，只包含需要的字段
    const baseData = {
      username: userFormData.value.username,
      name: userFormData.value.name,
      email: userFormData.value.email,
      phone: userFormData.value.phone,
      department: userFormData.value.department,
      position: userFormData.value.position
    }
    
    let requestData: any = baseData
    
    if (isEditUser.value) {
      // 编辑用户时，如果密码为空则不发送密码字段
      if (userFormData.value.password !== '') {
        requestData = {
          ...baseData,
          password: userFormData.value.password
        }
      }
      
      // 编辑用户
      await userAPI.update(userFormData.value.id, requestData)
      ElMessage.success('用户编辑成功')
    } else {
      // 添加用户，必须包含密码
      requestData = {
        ...baseData,
        password: userFormData.value.password
      }
      
      // 添加用户
      await userAPI.create(requestData)
      ElMessage.success('用户添加成功')
    }
    
    // 重新加载用户列表
    await loadUsers()
    userDialogVisible.value = false
  } catch (error) {
    console.error('保存用户失败:', error)
    const errorObj = error as any
    console.error('错误详情:', errorObj.response?.data || errorObj.message)
    // 显示详细的错误信息
    const errorMsg = errorObj.response?.data?.error || errorObj.message || '保存用户失败'
    ElMessage.error(errorMsg)
  } finally {
    loading.value.saveUser = false
  }
}

// 格式化日期
const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleString()
}
</script>

<style scoped>
.users-view {
  width: 100%;
  padding: 0;
  min-height: 100vh;
  background-color: #f5f7fa;
}

/* 主容器，设置最大宽度和居中 */
.users-view {
  display: flex;
  flex-direction: column;
  max-width: 100%;
  margin: 0 auto;
}

.page-header {
  margin: 0 20px 30px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.header-content h2 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: #333;
}

.users-section {
  margin: 0 20px 30px;
}

.users-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding: 0;
}

.users-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.users-card {
  margin-bottom: 30px;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;
}

.users-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

/* 表格容器，确保占满宽度 */
:deep(.el-table) {
  width: 100% !important;
}

/* 调整表格行高和内边距 */
:deep(.el-table__row) {
  height: 60px;
  transition: all 0.2s ease;
}

:deep(.el-table__row:hover) {
  background-color: rgba(0, 0, 0, 0.02);
}

:deep(.el-table__cell) {
  padding: 16px 12px;
  font-size: 14px;
  color: #303133;
}

/* 调整表格头部样式 */
:deep(.el-table__header-wrapper .el-table__header) {
  background-color: #fafafa;
}

:deep(.el-table__header-wrapper th.el-table__cell) {
  background-color: #fafafa;
  font-weight: 600;
  color: #303133;
  font-size: 14px;
}

/* 调整操作按钮间距 */
:deep(.el-button + .el-button) {
  margin-left: 12px;
}

/* 状态开关自定义样式 */
.status-switch-container {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 4px 0;
}

:deep(.custom-status-switch) {
  --el-switch-on-color: #67c23a;
  --el-switch-off-color: #f56c6c;
  --el-switch-on-text-color: #fff;
  --el-switch-off-text-color: #fff;
  --el-switch-border-radius: 20px;
  --el-switch-on-border-color: #67c23a;
  --el-switch-off-border-color: #f56c6c;
  font-size: 14px;
  font-weight: 500;
}

:deep(.custom-status-switch .el-switch__core) {
  border-radius: 20px;
  transition: all 0.3s cubic-bezier(0.645, 0.045, 0.355, 1);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

:deep(.custom-status-switch .el-switch__core::after) {
  border-radius: 50%;
  transition: all 0.3s cubic-bezier(0.645, 0.045, 0.355, 1);
  width: 20px;
  height: 20px;
}

:deep(.custom-status-switch.is-checked .el-switch__core::after) {
  transform: translateX(20px);
}

:deep(.custom-status-switch .el-switch__label) {
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s ease;
}

:deep(.custom-status-switch .el-switch__label--left) {
  color: #f56c6c;
  margin-right: 8px;
}

:deep(.custom-status-switch.is-checked .el-switch__label--left) {
  color: rgba(103, 194, 58, 0.6);
}

:deep(.custom-status-switch .el-switch__label--right) {
  color: rgba(103, 194, 58, 0.6);
  margin-left: 8px;
}

:deep(.custom-status-switch.is-checked .el-switch__label--right) {
  color: #67c23a;
}

/* 状态指示器样式 */
:deep(.el-table__cell) {
  position: relative;
}

/* 响应式设计，适配不同屏幕尺寸 */
@media (max-width: 1200px) {
  .users-view {
    padding: 0 10px;
  }
  
  .page-header,
  .users-section {
    margin: 0 10px 20px;
  }
  
  .users-card {
    padding: 15px;
  }
}

@media (max-width: 768px) {
  .users-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  
  .users-card {
    padding: 10px;
  }
  
  :deep(.el-table__cell) {
    padding: 12px 8px;
    font-size: 13px;
  }
  
  :deep(.el-button) {
    padding: 4px 12px;
    font-size: 12px;
  }
  
  :deep(.custom-status-switch) {
    font-size: 12px;
  }
  
  :deep(.custom-status-switch .el-switch__label--left),
  :deep(.custom-status-switch .el-switch__label--right) {
    margin: 0 4px;
  }
}
</style>
