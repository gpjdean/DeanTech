<template>
  <div class="profile-view">
    <el-card shadow="never" class="page-header">
      <div class="header-content">
        <h2>个人中心</h2>
      </div>
    </el-card>

    <!-- 个人中心内容 -->
    <el-card shadow="hover" class="config-card" v-loading="loading.profile" element-loading-text="加载个人信息中...">
      <template #header>
        <div class="card-header">
          <span>个人信息</span>
        </div>
      </template>
      
      <el-form ref="profileForm" :model="profileFormData" label-width="120px" class="profile-form">
        <!-- 头像设置 -->
        <el-form-item label="头像">
          <div class="avatar-uploader">
            <el-upload
              class="avatar-uploader"
              action="#"
              :show-file-list="false"
              :on-success="handleAvatarSuccess"
              :before-upload="beforeAvatarUpload"
              :auto-upload="true"
              :http-request="uploadAvatar"
            >
              <img v-if="profileFormData.avatar" :src="profileFormData.avatar" class="avatar">
              <el-icon v-else class="avatar-uploader-icon"><Plus /></el-icon>
            </el-upload>
            <div class="avatar-hint">支持 JPG、PNG 格式，大小不超过 2MB</div>
          </div>
        </el-form-item>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="用户名" prop="username">
              <el-input v-model="profileFormData.username" placeholder="请输入用户名" disabled /></el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="profileFormData.email" placeholder="请输入邮箱" /></el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="姓名" prop="name">
              <el-input v-model="profileFormData.name" placeholder="请输入姓名" /></el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="电话" prop="phone">
              <el-input v-model="profileFormData.phone" placeholder="请输入电话" /></el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="部门">
          <el-input v-model="profileFormData.department" placeholder="请输入部门" /></el-form-item>

        <el-form-item label="职位">
          <el-input v-model="profileFormData.position" placeholder="请输入职位" /></el-form-item>

        <el-form-item label="备注">
          <el-input v-model="profileFormData.remark" type="textarea" placeholder="请输入备注" :rows="3" /></el-form-item>

        <el-form-item>
          <el-button type="primary" @click="saveProfile">保存修改</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 密码修改 -->
    <el-card shadow="hover" class="config-card">
      <template #header>
        <div class="card-header">
          <span>修改密码</span>
        </div>
      </template>
      
      <el-form ref="passwordForm" :model="passwordFormData" label-width="120px" class="password-form">
        <el-form-item label="原密码" prop="oldPassword">
          <el-input v-model="passwordFormData.oldPassword" type="password" placeholder="请输入原密码" show-password /></el-form-item>

        <el-form-item label="新密码" prop="newPassword">
          <el-input v-model="passwordFormData.newPassword" type="password" placeholder="请输入新密码" show-password /></el-form-item>

        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input v-model="passwordFormData.confirmPassword" type="password" placeholder="请再次输入新密码" show-password /></el-form-item>

        <el-form-item>
          <el-button type="primary" @click="changePassword">保存密码</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { userAPI } from '../api/api'



// 表单引用
const profileForm = ref()
const passwordForm = ref()

// 加载状态
const loading = ref({
  profile: false,
  password: false
})

// 个人设置表单
const profileFormData = ref({
  id: 0,
  avatar: '',
  username: '',
  email: '',
  name: '',
  phone: '',
  department: '',
  position: '',
  remark: ''
})

// 密码修改表单
const passwordFormData = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})


// 组件挂载时加载当前用户信息
onMounted(() => {
  // 加载当前用户信息
  loadCurrentUser()
})


// 加载当前用户信息
const loadCurrentUser = async () => {
  try {
    loading.value.profile = true
    const userData = await userAPI.getCurrentUser()
    const avatar = userData.avatar || userData.Avatar || ''
    // 将后端返回的大写字段名转换为小写驼峰式
    profileFormData.value = {
      id: userData.id || userData.ID || 0, // 保存用户ID
      avatar: avatar,
      username: userData.username || userData.Username,
      email: userData.email || userData.Email,
      name: userData.name || userData.Name || '',
      phone: userData.phone || userData.Phone || '',
      department: userData.department || userData.Department || '',
      position: userData.position || userData.Position || '',
      remark: userData.remark || userData.Remark || ''
    }
    // 将用户头像保存到localStorage，避免刷新页面时出现默认头像
    localStorage.setItem('user.avatar', avatar)
    localStorage.setItem('username', userData.username || userData.Username || '管理员')
  } catch (error) {
    console.error('加载用户信息失败:', error)
    ElMessage.error('加载用户信息失败')
  } finally {
    loading.value.profile = false
  }
}

// 自定义头像上传方法
const uploadAvatar = (options: any) => {
  const file = options.file
  // 这里可以添加实际的上传逻辑，比如调用API上传到服务器
  // 目前我们只在客户端处理
  const reader = new FileReader()
  reader.onload = (e) => {
    const result = e.target?.result as string
    profileFormData.value.avatar = result
    options.onSuccess(result)
  }
  reader.readAsDataURL(file)
}

// 头像上传成功处理
const handleAvatarSuccess = () => {
  ElMessage.success('头像上传成功')
}

// 头像上传前校验
const beforeAvatarUpload = (file: File) => {
  const isJPG = file.type === 'image/jpeg' || file.type === 'image/png'
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isJPG) {
    ElMessage.error('请上传 JPG 或 PNG 格式的图片')
    return false
  }
  if (!isLt2M) {
    ElMessage.error('图片大小不能超过 2MB')
    return false
  }
  return true
}

// 保存个人设置
const saveProfile = async () => {
  try {
    loading.value.profile = true
    // 构建符合后端期望的请求体，包含头像和备注数据
    const requestData = {
      user_id: profileFormData.value.id,
      email: profileFormData.value.email,
      name: profileFormData.value.name,
      phone: profileFormData.value.phone,
      department: profileFormData.value.department,
      position: profileFormData.value.position,
      remark: profileFormData.value.remark, // 添加备注数据
      avatar: profileFormData.value.avatar // 添加头像数据
    }
    await userAPI.updateProfile(requestData)
    // 保存成功后重新加载用户信息，确保显示最新数据
    await loadCurrentUser()
    // 触发全局事件，通知App.vue更新用户信息
    window.dispatchEvent(new CustomEvent('userInfoUpdated'))
    ElMessage.success('个人信息保存成功')
  } catch (error: any) {
        console.error('保存个人信息失败:', error)
        console.error('错误详情:', error.response?.data || error.message)
        // 显示详细的错误信息
        const errorMsg = error.response?.data?.error || error.message || '保存个人信息失败'
        ElMessage.error(errorMsg)
      } finally {
        loading.value.profile = false
      }
}

// 修改密码
const changePassword = async () => {
  if (passwordFormData.value.newPassword !== passwordFormData.value.confirmPassword) {
    ElMessage.error('两次输入的密码不一致')
    return
  }
  
  try {
    loading.value.password = true
    // 构建符合后端期望的请求体
    await userAPI.changePassword({
      user_id: profileFormData.value.id,
      old_password: passwordFormData.value.oldPassword,
      new_password: passwordFormData.value.newPassword
    })
    ElMessage.success('密码修改成功')
    // 清空表单
    passwordFormData.value = {
      oldPassword: '',
      newPassword: '',
      confirmPassword: ''
    }
  } catch (error: any) {
    console.error('修改密码失败:', error)
    console.error('错误详情:', error.response?.data || error.message)
    // 优化错误提示，根据不同错误类型显示具体信息
    let errorMsg = '修改密码失败'
    
    // 检查错误响应结构
    if (error.response) {
      // 处理不同的错误响应格式
      const errorData = error.response.data
      
      // 情况1: error.response.data.error 是字符串
      if (errorData.error && typeof errorData.error === 'string') {
        errorMsg = errorData.error
        
        // 针对常见错误类型进行更友好的提示
        if (errorMsg.includes('旧密码') || errorMsg.includes('原密码') || 
            errorMsg.includes('old password') || errorMsg.includes('original password')) {
          errorMsg = '原密码不正确，请重新输入'
        } else if (errorMsg.includes('新密码') || errorMsg.includes('new password')) {
          errorMsg = '新密码不符合要求，请检查密码规则'
        }
      }
      // 情况2: error.response.data.error 是对象，包含字段级错误
      else if (errorData.error && typeof errorData.error === 'object') {
        if (errorData.error.oldPassword || errorData.error.old_password) {
          errorMsg = '原密码不正确，请重新输入'
        } else if (errorData.error.newPassword || errorData.error.new_password) {
          errorMsg = '新密码不符合要求，请检查密码规则'
        } else {
          // 其他对象类型错误，尝试转换为字符串
          errorMsg = JSON.stringify(errorData.error)
        }
      }
      // 情况3: error.response.data 直接是字符串
      else if (typeof errorData === 'string') {
        errorMsg = errorData
        if (errorMsg.includes('旧密码') || errorMsg.includes('原密码') || 
            errorMsg.includes('old password') || errorMsg.includes('original password')) {
          errorMsg = '原密码不正确，请重新输入'
        }
      }
      // 情况4: 根据HTTP状态码判断
      else if (error.response.status === 401 || error.response.status === 403) {
        errorMsg = '原密码不正确，请重新输入'
      }
    } 
    // 情况5: 没有response，直接使用error.message
    else if (error.message) {
      errorMsg = error.message
      if (errorMsg.includes('401') || errorMsg.includes('403') || 
          errorMsg.includes('unauthorized') || errorMsg.includes('forbidden')) {
        errorMsg = '原密码不正确，请重新输入'
      }
    }
    
    ElMessage.error(errorMsg)
  } finally {
    loading.value.password = false
  }
}
</script>

<style scoped>
.profile-view {
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

.profile-form {
  max-width: 600px;
}

.avatar-uploader {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.avatar {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  cursor: pointer;
  object-fit: cover;
}

.avatar-uploader-icon {
  width: 120px;
  height: 120px;
  background-color: #f5f7fa;
  border: 1px dashed #d9d9d9;
  border-radius: 50%;
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 28px;
  color: #8c8c8c;
  cursor: pointer;
  transition: all 0.3s;
}

.avatar-uploader-icon:hover {
  border-color: #409eff;
  color: #409eff;
}

.avatar-hint {
  margin-top: 8px;
  font-size: 12px;
  color: #909399;
}

.password-form {
  max-width: 600px;
}
</style>
