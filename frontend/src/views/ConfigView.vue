<template>
  <div class="config-view">
    <el-card shadow="never" class="page-header">
      <div class="header-content">
        <h2>系统风格</h2>
      </div>
    </el-card>

    <!-- 统一加载遮罩 -->
    <div v-if="loading.systemStyle" class="global-loading">
      <el-icon class="loading-icon"><Loading /></el-icon>
      <span class="loading-text">加载中...</span>
    </div>

    <el-card shadow="hover" class="config-card">
      <template #header>
        <div class="card-header">
          <span>登录页面设置</span>
        </div>
      </template>
      
      <el-form ref="systemStyleForm" :model="systemStyleFormData" label-width="120px" class="system-style-form">
        <el-form-item label="系统名称" prop="systemName">
          <el-input v-model="systemStyleFormData.systemName" placeholder="请输入登录页面显示的系统名称" />
        </el-form-item>

        <el-form-item label="系统副标题" prop="systemSubtitle">
          <el-input v-model="systemStyleFormData.systemSubtitle" placeholder="请输入登录页面显示的系统副标题" />
        </el-form-item>

        <el-form-item label="登录Logo" prop="loginLogo">
          <el-input v-model="systemStyleFormData.loginLogo" placeholder="请输入Logo图片URL" />
          <div class="form-hint">支持PNG、JPG格式，建议尺寸：120x120px</div>
        </el-form-item>

        <el-form-item label="主题色" prop="primaryColor">
          <div class="color-picker-row">
            <el-color-picker v-model="systemStyleFormData.primaryColor" show-alpha />
            <el-input v-model="systemStyleFormData.primaryColor" placeholder="#409eff" class="color-input" />
          </div>
        </el-form-item>

        <el-form-item label="登录风格" prop="loginStyle">
          <el-select v-model="systemStyleFormData.loginStyle" placeholder="请选择登录风格">
            <el-option label="左右布局" value="style1" />
            <el-option label="居中卡片" value="style2" />
            <el-option label="深色主题" value="style3" />
            <el-option label="简约设计" value="style4" />
            <el-option label="渐变卡片" value="style5" />
          </el-select>
        </el-form-item>

        <el-form-item label="登录页脚" prop="loginFooter">
          <el-input v-model="systemStyleFormData.loginFooter" placeholder="请输入登录页面显示的页脚信息" />
          <div class="form-hint">设置后将显示在登录页面底部</div>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="saveSystemStyle" :loading="loading.systemStyle">保存设置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
    
    <el-card shadow="hover" class="config-card">
      <template #header>
        <div class="card-header">
          <span>首页信息设置</span>
        </div>
      </template>
      
      <el-form :model="systemStyleFormData" label-width="120px" class="system-style-form">
        <el-form-item label="首页标题" prop="homeTitle">
          <el-input v-model="systemStyleFormData.homeTitle" placeholder="请输入首页左上角显示的系统标题" />
        </el-form-item>

        <el-form-item label="首页图标" prop="homeIcon">
          <el-input v-model="systemStyleFormData.homeIcon" placeholder="请输入首页左上角显示的图标URL" />
          <div class="form-hint">支持PNG、JPG、SVG格式，建议尺寸：32x32px</div>
        </el-form-item>

        <el-form-item label="首页Title" prop="homeBrowserTitle">
          <el-input v-model="systemStyleFormData.homeBrowserTitle" placeholder="请输入浏览器标签页显示的标题" />
          <div class="form-hint">设置后将显示在浏览器标签页上</div>
        </el-form-item>

        <el-form-item label="首页页脚" prop="homeFooter">
          <el-input v-model="systemStyleFormData.homeFooter" placeholder="请输入首页页脚显示的信息" />
          <div class="form-hint">设置后将显示在页面底部页脚处</div>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="saveSystemStyle" :loading="loading.systemStyle">保存设置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Loading } from '@element-plus/icons-vue'

import { ElMessage } from 'element-plus'
import { systemSettingAPI } from '../api/api'

// 表单引用
const systemStyleForm = ref()

// 加载状态
const loading = ref({
  systemStyle: false
})

// 系统风格表单数据
const systemStyleFormData = ref({
  systemName: '',
  systemSubtitle: '',
  loginLogo: '',
  primaryColor: '#409eff',
  loginStyle: 'style1',
  loginFooter: '',
  homeTitle: '',
  homeIcon: '',
  homeBrowserTitle: '',
  homeFooter: ''
})

// 组件挂载时加载数据
onMounted(() => {
  loadSystemSettings()
})

// 加载系统设置
const loadSystemSettings = async () => {
  try {
    loading.value.systemStyle = true
    
    const response = await systemSettingAPI.list()
    const settings = response.data
    
    // 查找系统名称设置
    const systemNameSetting = settings.find((setting: any) => setting.key === 'system.name')
    if (systemNameSetting) {
      systemStyleFormData.value.systemName = systemNameSetting.value
    }
    
    // 查找系统副标题设置
    const systemSubtitleSetting = settings.find((setting: any) => setting.key === 'system.subtitle')
    if (systemSubtitleSetting) {
      systemStyleFormData.value.systemSubtitle = systemSubtitleSetting.value
    }
    
    // 查找登录Logo设置
    const loginLogoSetting = settings.find((setting: any) => setting.key === 'login.logo')
    if (loginLogoSetting) {
      systemStyleFormData.value.loginLogo = loginLogoSetting.value
    }
    
    // 查找主题色设置
    const primaryColorSetting = settings.find((setting: any) => setting.key === 'theme.primaryColor')
    if (primaryColorSetting) {
      systemStyleFormData.value.primaryColor = primaryColorSetting.value
    }
    
    // 查找登录风格设置
    const loginStyleSetting = settings.find((setting: any) => setting.key === 'login.style')
    if (loginStyleSetting) {
      systemStyleFormData.value.loginStyle = loginStyleSetting.value
    }
    
    // 查找登录页脚设置
    const loginFooterSetting = settings.find((setting: any) => setting.key === 'login.footer')
    if (loginFooterSetting) {
      systemStyleFormData.value.loginFooter = loginFooterSetting.value
    }
    
    // 查找首页标题设置
    const homeTitleSetting = settings.find((setting: any) => setting.key === 'home.title')
    if (homeTitleSetting) {
      systemStyleFormData.value.homeTitle = homeTitleSetting.value
    }
    
    // 查找首页图标设置
    const homeIconSetting = settings.find((setting: any) => setting.key === 'home.icon')
    if (homeIconSetting) {
      systemStyleFormData.value.homeIcon = homeIconSetting.value
    }
    
    // 查找首页浏览器标题设置
    const homeBrowserTitleSetting = settings.find((setting: any) => setting.key === 'home.browserTitle')
    if (homeBrowserTitleSetting) {
      systemStyleFormData.value.homeBrowserTitle = homeBrowserTitleSetting.value
    }
    
    // 查找首页页脚设置
    const homeFooterSetting = settings.find((setting: any) => setting.key === 'home.footer')
    if (homeFooterSetting) {
      systemStyleFormData.value.homeFooter = homeFooterSetting.value
    }
  } catch (error) {
    console.error('加载系统设置失败:', error)
    ElMessage.error('加载系统设置失败')
  } finally {
    loading.value.systemStyle = false
  }
}

// 保存系统风格设置
const saveSystemStyle = async () => {
  try {
    loading.value.systemStyle = true
    
    // 获取现有设置
    const response = await systemSettingAPI.list()
    const settings = response.data
    
    // 定义需要保存的设置项
    const settingItems = [
      { key: 'system.name', value: systemStyleFormData.value.systemName, description: '登录页面系统名称' },
      { key: 'system.subtitle', value: systemStyleFormData.value.systemSubtitle, description: '登录页面系统副标题' },
      { key: 'login.logo', value: systemStyleFormData.value.loginLogo, description: '登录页面Logo' },
      { key: 'theme.primaryColor', value: systemStyleFormData.value.primaryColor, description: '系统主题色' },
      { key: 'login.style', value: systemStyleFormData.value.loginStyle, description: '登录页面风格' },
      { key: 'login.footer', value: systemStyleFormData.value.loginFooter, description: '登录页面底部页脚信息' },
      { key: 'home.title', value: systemStyleFormData.value.homeTitle, description: '首页左上角系统标题' },
      { key: 'home.icon', value: systemStyleFormData.value.homeIcon, description: '首页左上角系统图标' },
      { key: 'home.browserTitle', value: systemStyleFormData.value.homeBrowserTitle, description: '浏览器标签页标题' },
      { key: 'home.footer', value: systemStyleFormData.value.homeFooter, description: '页面底部页脚信息' }
    ]
    
    // 保存每个设置项
    for (const item of settingItems) {
      const existingSetting = settings.find((setting: any) => setting.key === item.key)
      if (existingSetting) {
        // 更新现有设置
        await systemSettingAPI.update(item.key, {
          value: item.value,
          description: item.description
        })
      } else {
        // 创建新设置
        await systemSettingAPI.create({
          key: item.key,
          value: item.value,
          description: item.description
        })
      }
    }
    
    // 立即更新LocalStorage，避免刷新页面时出现默认值
    localStorage.setItem('login.footer', systemStyleFormData.value.loginFooter)
    localStorage.setItem('home.title', systemStyleFormData.value.homeTitle)
    localStorage.setItem('home.icon', systemStyleFormData.value.homeIcon)
    localStorage.setItem('home.browserTitle', systemStyleFormData.value.homeBrowserTitle)
    localStorage.setItem('home.footer', systemStyleFormData.value.homeFooter)
    
    // 发送事件通知App.vue更新状态
    window.dispatchEvent(new Event('homeSettingsUpdated'))
    
    ElMessage.success('系统设置保存成功')
  } catch (error) {
    console.error('保存系统设置失败:', error)
    ElMessage.error('保存系统设置失败')
  } finally {
    loading.value.systemStyle = false
  }
}
</script>

<style scoped>
.config-view {
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

.config-tabs {
  margin-bottom: 30px;
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

.users-section {
  margin-bottom: 30px;
}

.users-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding: 0 10px;
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
}

/* 调整标签页样式 */
:deep(.el-tabs__content) {
  padding-top: 20px;
}

/* 系统风格表单样式 */
.system-style-form {
  max-width: 600px;
}

.form-hint {
  margin-top: 8px;
  font-size: 12px;
  color: #909399;
}

.color-picker-row {
  display: flex;
  align-items: center;
  gap: 16px;
}

.color-input {
  flex: 1;
  width: 120px;
}

/* 颜色选择器样式调整 */
:deep(.el-color-picker) {
  margin-right: 16px;
}

/* 全局加载遮罩样式 */
.global-loading {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(255, 255, 255, 0.8);
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  z-index: 9999;
  backdrop-filter: blur(2px);
}

.loading-icon {
  font-size: 48px;
  color: #409eff;
  margin-bottom: 16px;
  animation: rotate 1s linear infinite;
}

.loading-text {
  font-size: 16px;
  color: #606266;
}

@keyframes rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

</style>
