<template>
  <div :class="['login-container', `login-style-${loginStyle}`]">

    <!-- 风格1: 左右布局 -->
    <div class="login-style-1" v-if="loginStyle === '1'">
      <!-- 左侧装饰区域 -->
      <div class="login-decoration">
        <div class="decoration-content">
          <div v-if="systemSettings.logo" class="decoration-logo">
                <img :src="systemSettings.logo" alt="系统Logo" class="logo-image" />
              </div>
              <el-icon v-else class="decoration-icon"><Setting /></el-icon>
          <h1 class="decoration-title">{{ systemSettings.name }}</h1>
          <p class="decoration-subtitle">{{ systemSettings.subtitle }}</p>
          <div class="decoration-features">
            <div class="feature-item">
              <el-icon class="feature-icon"><DataAnalysis /></el-icon>
              <span>监控告警</span>
            </div>
            <div class="feature-item">
              <el-icon class="feature-icon"><Cpu /></el-icon>
              <span>K8s管理</span>
            </div>
            <div class="feature-item">
              <el-icon class="feature-icon"><OfficeBuilding /></el-icon>
              <span>主机管理</span>
            </div>
            <div class="feature-item">
              <el-icon class="feature-icon"><Setting /></el-icon>
              <span>自动化运维</span>
            </div>
            <div class="feature-item">
              <el-icon class="feature-icon"><Monitor /></el-icon>
              <span>SSH终端</span>
            </div>
            <div class="feature-item">
              <el-icon class="feature-icon"><Upload /></el-icon>
              <span>文件传输</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧登录表单区域 -->
      <div class="login-form-section">
        <div class="form-wrapper">
          <el-card class="login-card" shadow="hover">
            <template #header>
              <div class="card-header">
                <h2 class="login-title">欢迎登录</h2>
                <p class="login-subtitle">请输入您的用户名和密码</p>
              </div>
            </template>
            
            <el-form :model="loginForm" :rules="loginRules" ref="loginFormRef" label-position="top">
              <el-form-item prop="username">
                <el-input 
                  v-model="loginForm.username" 
                  placeholder="用户名" 
                  size="large"
                  autocomplete="username"
                >
                  <template #prepend>
                    <el-icon><User /></el-icon>
                  </template>
                </el-input>
              </el-form-item>
              
              <el-form-item prop="password">
                <el-input 
                  v-model="loginForm.password" 
                  placeholder="密码" 
                  type="password" 
                  show-password 
                  size="large"
                  autocomplete="current-password"
                  @keyup.enter="handleLogin"
                >
                  <template #prepend>
                    <el-icon><Lock /></el-icon>
                  </template>
                </el-input>
              </el-form-item>
              
              <el-button 
                type="primary" 
                :loading="loading" 
                @click="handleLogin" 
                block
                size="large"
                class="login-button"
              >
                <span v-if="!loading">登录</span>
                <span v-else>登录中...</span>
              </el-button>
              
              <div class="form-actions">
                <el-checkbox v-model="loginForm.rememberMe" class="remember-me">
                  <span>记住我</span>
                </el-checkbox>
                <el-link type="primary" @click="showForgotPassword = true" class="forgot-link">
                  忘记密码
                </el-link>
              </div>
              
              <div class="divider-line">
                <span class="divider-text">或</span>
              </div>
              
              <el-button 
                type="default" 
                @click="showRegister = true" 
                block
                size="large"
                class="register-button"
              >
                注册账号
              </el-button>
            </el-form>

            <div class="login-footer">
          <p class="footer-text">{{ systemSettings.footer || '© 2026 Dean\'s Copyright. All rights reserved.' }}</p>
        </div>
          </el-card>
        </div>
      </div>
    </div>

    <!-- 风格2: 居中卡片布局 -->
    <div class="login-style-2" v-else-if="loginStyle === '2'">
      <div class="style-2-background">
        <div class="style-2-gradient"></div>
        <div class="style-2-pattern"></div>
      </div>
      <div class="style-2-content">
        <div class="style-2-logo">
          <el-icon class="style-2-logo-icon"><Setting /></el-icon>
          <h1 class="style-2-title">{{ systemSettings.name }}</h1>
        </div>
        <el-card class="style-2-card" shadow="lg">
          <template #header>
            <div class="style-2-card-header">
              <h2 class="style-2-login-title">登录</h2>
              <p class="style-2-login-subtitle">欢迎回来，请登录您的账号</p>
            </div>
          </template>
          
          <el-form :model="loginForm" :rules="loginRules" ref="loginFormRef" label-position="top">
            <el-form-item prop="username">
              <el-input 
                v-model="loginForm.username" 
                placeholder="用户名" 
                size="large"
                autocomplete="username"
                class="style-2-input"
              >
                <template #prepend>
                  <el-icon><User /></el-icon>
                </template>
              </el-input>
            </el-form-item>
            
            <el-form-item prop="password">
              <el-input 
                v-model="loginForm.password" 
                placeholder="密码" 
                type="password" 
                show-password 
                size="large"
                autocomplete="current-password"
                class="style-2-input"
                @keyup.enter="handleLogin"
              >
                <template #prepend>
                  <el-icon><Lock /></el-icon>
                </template>
              </el-input>
            </el-form-item>
            
            <el-button 
              type="primary" 
              :loading="loading" 
              @click="handleLogin" 
              block
              size="large"
              class="style-2-button"
            >
              <span v-if="!loading">登录</span>
              <span v-else>登录中...</span>
            </el-button>
            
            <div class="style-2-actions">
              <el-checkbox v-model="loginForm.rememberMe" class="style-2-remember">
                <span>记住我</span>
              </el-checkbox>
              <el-link type="primary" @click="showForgotPassword = true" class="forgot-link">
                忘记密码
              </el-link>
            </div>
            
            <div class="divider-line">
              <span class="divider-text">或</span>
            </div>
            
            <el-button 
              type="default" 
              @click="showRegister = true" 
              block
              size="large"
              class="style-2-register-button"
            >
              注册账号
            </el-button>
          </el-form>
        </el-card>
        <div class="style-2-footer">
          <p class="style-2-copyright">{{ systemSettings.footer || '© 2026 Dean\'s Copyright. All rights reserved.' }}</p>
        </div>
      </div>
    </div>

    <!-- 风格3: 深色主题 -->
    <div class="login-style-3" v-else-if="loginStyle === '3'">
      <div class="style-3-background">
        <div class="style-3-gradient"></div>
        <div class="style-3-grid"></div>
      </div>
      <div class="style-3-content">
        <div class="style-3-logo">
          <el-icon class="style-3-logo-icon"><Setting /></el-icon>
          <h1 class="style-3-title">{{ systemSettings.name }}</h1>
        </div>
        <div class="style-3-card">
          <div class="style-3-card-header">
            <h2 class="style-3-login-title">登录</h2>
            <p class="style-3-login-subtitle">欢迎回来，请登录您的账号</p>
          </div>
          <el-form :model="loginForm" :rules="loginRules" ref="loginFormRef" label-position="top" class="style-3-form">
            <el-form-item prop="username">
              <el-input 
                v-model="loginForm.username" 
                placeholder="用户名" 
                size="large"
                autocomplete="username"
                class="style-3-input"
              >
                <template #prepend>
                  <el-icon color="#4facfe"><User /></el-icon>
                </template>
              </el-input>
            </el-form-item>
            
            <el-form-item prop="password">
              <el-input 
                v-model="loginForm.password" 
                placeholder="密码" 
                type="password" 
                show-password 
                size="large"
                autocomplete="current-password"
                class="style-3-input"
                @keyup.enter="handleLogin"
              >
                <template #prepend>
                  <el-icon color="#4facfe"><Lock /></el-icon>
                </template>
              </el-input>
            </el-form-item>
            
            <el-button 
              type="primary" 
              :loading="loading" 
              @click="handleLogin" 
              block
              size="large"
              class="style-3-button"
            >
              <span v-if="!loading">登录</span>
              <span v-else>登录中...</span>
            </el-button>
            
            <div class="style-3-actions">
              <el-checkbox v-model="loginForm.rememberMe" class="style-3-remember" text-color="#909399">
                <span>记住我</span>
              </el-checkbox>
              <el-link type="primary" @click="showForgotPassword = true" class="forgot-link">
                忘记密码
              </el-link>
            </div>
            
            <div class="divider-line">
              <span class="divider-text">或</span>
            </div>
            
            <el-button 
              type="default" 
              @click="showRegister = true" 
              block
              size="large"
              class="style-3-register-button"
            >
              注册账号
            </el-button>
          </el-form>
        </div>
        <div class="style-3-footer">
          <p class="style-3-copyright">{{ systemSettings.footer || '© 2026 Dean\'s Copyright. All rights reserved.' }}</p>
        </div>
      </div>
    </div>

    <!-- 风格4: 简约设计 -->
    <div class="login-style-4" v-else-if="loginStyle === '4'">
      <div class="style-4-background">
        <div class="style-4-shape style-4-shape-1"></div>
        <div class="style-4-shape style-4-shape-2"></div>
        <div class="style-4-shape style-4-shape-3"></div>
      </div>
      <div class="style-4-content">
        <el-card class="style-4-card" shadow="none">
          <div class="style-4-logo">
            <el-icon class="style-4-logo-icon"><Setting /></el-icon>
          </div>
          <h1 class="style-4-title">欢迎登录</h1>
          <p class="style-4-subtitle">请输入您的账号信息</p>
          
          <el-form :model="loginForm" :rules="loginRules" ref="loginFormRef" label-position="top" class="style-4-form">
            <el-form-item prop="username">
              <el-input 
                v-model="loginForm.username" 
                placeholder="用户名" 
                size="large"
                autocomplete="username"
              >
                <template #prepend>
                  <el-icon><User /></el-icon>
                </template>
              </el-input>
            </el-form-item>
            
            <el-form-item prop="password">
              <el-input 
                v-model="loginForm.password" 
                placeholder="密码" 
                type="password" 
                show-password 
                size="large"
                autocomplete="current-password"
                @keyup.enter="handleLogin"
              >
                <template #prepend>
                  <el-icon><Lock /></el-icon>
                </template>
              </el-input>
            </el-form-item>
            
            <el-button 
              type="primary" 
              :loading="loading" 
              @click="handleLogin" 
              block
              size="large"
              class="style-4-button"
            >
              <span v-if="!loading">登录</span>
              <span v-else>登录中...</span>
            </el-button>
            
            <div class="style-4-actions">
              <el-checkbox v-model="loginForm.rememberMe" class="style-4-remember">
                <span>记住我</span>
              </el-checkbox>
              <el-link type="primary" @click="showForgotPassword = true" class="forgot-link">
                忘记密码
              </el-link>
            </div>
            
            <div class="divider-line">
              <span class="divider-text">或</span>
            </div>
            
            <el-button 
              type="default" 
              @click="showRegister = true" 
              block
              size="large"
              class="style-4-register-button"
            >
              注册
            </el-button>
          </el-form>
        </el-card>
        <div class="style-4-footer">
          <p class="style-4-copyright">{{ systemSettings.footer || '© 2026 Dean\'s Copyright' }}</p>
        </div>
      </div>
    </div>

    <!-- 风格5: 渐变卡片 -->
    <div class="login-style-5" v-else>
      <div class="style-5-background">
        <div class="style-5-gradient"></div>
      </div>
      <div class="style-5-content">
        <div class="style-5-logo">
          <el-icon class="style-5-logo-icon"><Setting /></el-icon>
        </div>
        <h1 class="style-5-title">{{ systemSettings.name }}</h1>
        <p class="style-5-subtitle">专业的运维管理解决方案</p>
        <div class="style-5-card">
          <div class="style-5-card-inner">
            <div class="style-5-card-header">
              <h2 class="style-5-login-title">账户登录</h2>
            </div>
            <el-form :model="loginForm" :rules="loginRules" ref="loginFormRef" label-position="top" class="style-5-form">
              <el-form-item prop="username">
                <el-input 
                  v-model="loginForm.username" 
                  placeholder="用户名" 
                  size="large"
                  autocomplete="username"
                  class="style-5-input"
                >
                  <template #prepend>
                    <el-icon><User /></el-icon>
                  </template>
                </el-input>
              </el-form-item>
              
              <el-form-item prop="password">
                <el-input 
                  v-model="loginForm.password" 
                  placeholder="密码" 
                  type="password" 
                  show-password 
                  size="large"
                  autocomplete="current-password"
                  class="style-5-input"
                  @keyup.enter="handleLogin"
                >
                  <template #prepend>
                    <el-icon><Lock /></el-icon>
                  </template>
                </el-input>
              </el-form-item>
              
              <el-button 
                type="primary" 
                :loading="loading" 
                @click="handleLogin" 
                block
                size="large"
                class="style-5-button"
              >
                <span v-if="!loading">登录</span>
                <span v-else>登录中...</span>
              </el-button>
              
              <div class="style-5-actions">
                <el-checkbox v-model="loginForm.rememberMe" class="style-5-remember">
                  <span>记住登录状态</span>
                </el-checkbox>
                <el-link type="primary" @click="showForgotPassword = true" class="forgot-link">
                  忘记密码
                </el-link>
              </div>
              
              <div class="divider-line">
                <span class="divider-text">或</span>
              </div>
              
              <el-button 
                type="default" 
                @click="showRegister = true" 
                block
                size="large"
                class="style-5-register-button"
              >
                立即注册
              </el-button>
            </el-form>
          </div>
        </div>
        <div class="style-5-footer">
          <p class="style-5-copyright">{{ systemSettings.footer || '© 2026 Dean\'s Copyright. All rights reserved.' }}</p>
        </div>
      </div>
    </div>
    
    <!-- 注册对话框 -->
    <el-dialog v-model="showRegister" title="注册新用户" width="500px" center>
      <el-form :model="registerForm" :rules="registerRules" ref="registerFormRef" label-position="top">
        <el-form-item prop="username">
          <el-input 
            v-model="registerForm.username" 
            placeholder="请输入用户名"
            :prefix-icon="User"
            size="large"
          ></el-input>
        </el-form-item>
        
        <el-form-item prop="name">
          <el-input 
            v-model="registerForm.name" 
            placeholder="请输入姓名"
            :prefix-icon="User"
            size="large"
          ></el-input>
        </el-form-item>
        
        <el-form-item prop="email">
          <el-input 
            v-model="registerForm.email" 
            placeholder="请输入邮箱"
            :prefix-icon="Message"
            size="large"
          ></el-input>
        </el-form-item>
        
        <el-form-item prop="phone">
          <el-input 
            v-model="registerForm.phone" 
            placeholder="请输入电话"
            :prefix-icon="Phone"
            size="large"
          ></el-input>
        </el-form-item>
        
        <el-form-item prop="department">
          <el-input 
            v-model="registerForm.department" 
            placeholder="请输入部门"
            :prefix-icon="OfficeBuilding"
            size="large"
          ></el-input>
        </el-form-item>
        
        <el-form-item prop="position">
          <el-input 
            v-model="registerForm.position" 
            placeholder="请输入职位"
            :prefix-icon="Setting"
            size="large"
          ></el-input>
        </el-form-item>
        
        <el-form-item prop="password">
          <el-input 
            v-model="registerForm.password" 
            placeholder="请输入密码（至少6位，包含字母和数字）"
            type="password"
            show-password
            :prefix-icon="Lock"
            size="large"
          ></el-input>
          <div class="password-strength" v-if="registerForm.password">
            <span class="strength-label">密码强度：</span>
            <span 
              class="strength-bar" 
              :class="getPasswordStrength()"
            ></span>
            <span class="strength-text">{{ getPasswordStrengthText() }}</span>
          </div>
        </el-form-item>
        
        <el-form-item prop="confirmPassword">
          <el-input 
            v-model="registerForm.confirmPassword" 
            placeholder="请再次输入密码"
            type="password"
            show-password
            :prefix-icon="Lock"
            size="large"
          ></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showRegister = false" size="large">取消</el-button>
          <el-button type="primary" :loading="registerLoading" @click="handleRegister" size="large">
            注册
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 忘记密码对话框 -->
    <el-dialog 
      v-model="showForgotPassword" 
      :title="forgotPasswordTitle" 
      width="500px" 
      center
      :show-close="forgotPasswordStep < 3"
      :close-on-click-modal="forgotPasswordStep < 3"
    >
      <!-- 步骤1: 输入邮箱 -->
      <div v-if="forgotPasswordStep === 1">
        <el-form :model="forgotPasswordForm" :rules="forgotPasswordRules" ref="forgotPasswordFormRef" label-position="top">
          <el-form-item prop="email">
            <el-input 
              v-model="forgotPasswordForm.email" 
              placeholder="请输入您的邮箱"
              :prefix-icon="Message"
              size="large"
            ></el-input>
          </el-form-item>
        </el-form>
      </div>
      
      <!-- 步骤2: 输入用户名 -->
      <div v-else-if="forgotPasswordStep === 2">
        <el-form :model="forgotPasswordForm" :rules="forgotPasswordRules" ref="forgotPasswordFormRef" label-position="top">
          <el-form-item prop="username">
            <el-input 
              v-model="forgotPasswordForm.username" 
              placeholder="请输入您的用户名"
              :prefix-icon="User"
              size="large"
            ></el-input>
          </el-form-item>
        </el-form>
      </div>
      
      <!-- 步骤3: 提示成功 -->
      <div v-else-if="forgotPasswordStep === 3" class="forgot-password-success">
        <el-icon class="success-icon"><CircleCheck /></el-icon>
        <h3>密码重置成功</h3>
        <p>我们已向您的邮箱发送了新的6位随机密码，请查收并登录。</p>
      </div>
      
      <template #footer>
        <span class="dialog-footer">
          <!-- 只在步骤1和2显示取消按钮 -->
          <el-button 
            v-if="forgotPasswordStep < 3" 
            @click="handleForgotPasswordCancel" 
            size="large"
          >
            取消
          </el-button>
          <!-- 步骤1和2显示下一步/确认按钮 -->
          <el-button 
            v-if="forgotPasswordStep < 3" 
            type="primary" 
            :loading="forgotPasswordLoading" 
            @click="handleForgotPasswordNext" 
            size="large"
          >
            {{ forgotPasswordStep === 1 ? '下一步' : '确认' }}
          </el-button>
          <!-- 步骤3只显示关闭按钮 -->
          <el-button 
            v-else 
            type="primary" 
            @click="handleForgotPasswordClose" 
            size="large"
          >
            关闭
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { User, Lock, Message, Phone, DataAnalysis, Setting, OfficeBuilding, Cpu, Upload, Monitor, CircleCheck } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'
import { systemSettingAPI, authAPI } from '../api/api'

const router = useRouter()
const loginFormRef = ref()
const registerFormRef = ref()
const forgotPasswordFormRef = ref()
const loading = ref(false)
const registerLoading = ref(false)
const forgotPasswordLoading = ref(false)
const showRegister = ref(false)
const showForgotPassword = ref(false)
const loginStyle = ref('1') // 1: 风格1, 2: 风格2

// 忘记密码相关变量
const forgotPasswordStep = ref(1) // 1: 输入邮箱, 2: 输入用户名, 3: 成功提示
const forgotPasswordForm = ref({
  email: '',
  username: ''
})

const forgotPasswordTitle = computed(() => {
  switch (forgotPasswordStep.value) {
    case 1: return '忘记密码 - 验证邮箱'
    case 2: return '忘记密码 - 验证用户名'
    case 3: return '密码重置成功'
    default: return '忘记密码'
  }
})

// 系统风格设置 - 从系统设置中获取
const systemSettings = ref({
  name: '运维管理平台',
  subtitle: '监控告警 · K8s管理 · 主机管理 · 自动化运维',
  logo: '',
  primaryColor: '#409eff',
  footer: ''
})

// 加载系统设置
const loadSystemSettings = async () => {
  try {
    const response = await systemSettingAPI.list()
    const settings = response.data
    
    // 查找系统名称设置
    const nameSetting = settings.find((setting: any) => setting.key === 'system.name')
    if (nameSetting) {
      systemSettings.value.name = nameSetting.value
    }
    
    // 查找系统副标题设置
    const subtitleSetting = settings.find((setting: any) => setting.key === 'system.subtitle')
    if (subtitleSetting) {
      systemSettings.value.subtitle = subtitleSetting.value
    }
    
    // 查找登录Logo设置
    const logoSetting = settings.find((setting: any) => setting.key === 'login.logo')
    if (logoSetting) {
      systemSettings.value.logo = logoSetting.value
    }
    
    // 查找主题色设置
    const colorSetting = settings.find((setting: any) => setting.key === 'theme.primaryColor')
    if (colorSetting) {
      systemSettings.value.primaryColor = colorSetting.value
    }
    
    // 查找登录风格设置
    const loginStyleSetting = settings.find((setting: any) => setting.key === 'login.style')
    if (loginStyleSetting) {
      // 将'style1'转换为'1'，'style2'转换为'2'，以此类推
      const styleValue = loginStyleSetting.value
      if (styleValue.startsWith('style')) {
        loginStyle.value = styleValue.charAt(5)
      }
    }
    
    // 查找登录页脚设置
    const loginFooterSetting = settings.find((setting: any) => setting.key === 'login.footer')
    if (loginFooterSetting) {
      systemSettings.value.footer = loginFooterSetting.value
    }
    
    // 应用主题色
    applyThemeColor()
  } catch (error) {
    console.error('加载系统设置失败:', error)
    // 使用默认设置，不影响页面正常显示
    console.log('使用默认系统设置')
  }
}

// 应用主题色
const applyThemeColor = () => {
  document.documentElement.style.setProperty('--primary-color', systemSettings.value.primaryColor)
}

// 组件挂载时加载系统设置和保存的登录信息
onMounted(() => {
  loadSystemSettings()
  loadSavedLoginInfo()
})

// 加载保存的登录信息
const loadSavedLoginInfo = () => {
  // 从本地存储读取保存的登录信息
  const savedUsername = localStorage.getItem('savedUsername')
  const savedPassword = localStorage.getItem('savedPassword')
  const savedRememberMe = localStorage.getItem('rememberMe')
  
  if (savedRememberMe === 'true' && savedUsername && savedPassword) {
    loginForm.value.username = savedUsername
    loginForm.value.password = savedPassword
    loginForm.value.rememberMe = true
  }
}

const loginForm = ref({
  username: '',
  password: '',
  rememberMe: false
})

const registerForm = ref({
  username: '',
  name: '',
  email: '',
  phone: '',
  department: '',
  position: '',
  password: '',
  confirmPassword: ''
})

const loginRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const registerRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_]*$/, message: '用户名只能包含字母、数字和下划线', trigger: 'blur' }
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
  ],
  department: [
    { max: 50, message: '部门名称不能超过 50 个字符', trigger: 'blur' }
  ],
  position: [
    { max: 50, message: '职位名称不能超过 50 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6个字符', trigger: 'blur' },
    { pattern: /^(?=.*[a-zA-Z])(?=.*\d)/, message: '密码必须包含字母和数字', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    ({ getFieldValue }: { getFieldValue: (field: string) => any }) => ({
      validator: (_: any, value: string) => {
        if (!value || getFieldValue('password') === value) {
          return Promise.resolve()
        }
        return Promise.reject(new Error('两次输入的密码不一致'))
      },
      trigger: 'blur, change' as const
    })
  ]
}

// 忘记密码验证规则
const forgotPasswordRules = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur, change' }
  ],
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  loginFormRef.value.validate(async (valid: boolean) => {
    if (valid) {
      loading.value = true
      try {
        // 调用登录API
        const response = await authAPI.login(loginForm.value.username, loginForm.value.password)
        
        const userData = response.data
        loading.value = false
        
        ElMessage.success('登录成功')
        // 存储用户信息到本地存储
        localStorage.setItem('token', 'mock-token-' + Date.now())
        localStorage.setItem('username', userData.username)
        localStorage.setItem('userId', userData.id.toString())
        
        // 处理记住我功能
        if (loginForm.value.rememberMe) {
          // 保存用户名和密码到本地存储
          localStorage.setItem('savedUsername', loginForm.value.username)
          localStorage.setItem('savedPassword', loginForm.value.password)
          localStorage.setItem('rememberMe', 'true')
        } else {
          // 清除之前保存的登录信息
          localStorage.removeItem('savedUsername')
          localStorage.removeItem('savedPassword')
          localStorage.removeItem('rememberMe')
        }
        
        // 触发用户信息更新事件，通知App组件更新用户信息
        window.dispatchEvent(new Event('userInfoUpdated'))
        
        // 跳转到仪表盘首页
        router.push('/dashboard')
      } catch (error: any) {
          loading.value = false
          // 根据后端返回的错误信息显示不同提示
          const errorMessage = error.response?.data?.error || '登录失败，请检查用户名和密码'
          // 如果是用户已被禁用，显示更友好的提示
          if (errorMessage === '用户已被禁用') {
            ElMessage.error('用户已被管理员禁用')
          } else {
            ElMessage.error(errorMessage)
          }
        }
    }
  })
}

// 计算密码强度
const getPasswordStrength = () => {
  const password = registerForm.value.password
  let strength = 'weak'
  
  if (password.length >= 8) {
    strength = 'medium'
  }
  
  if (password.length >= 10 && /[a-zA-Z]/.test(password) && /\d/.test(password) && /[^a-zA-Z0-9]/.test(password)) {
    strength = 'strong'
  }
  
  return strength
}

// 获取密码强度文本
const getPasswordStrengthText = () => {
  const strength = getPasswordStrength()
  const strengthMap: Record<string, string> = {
    weak: '弱',
    medium: '中',
    strong: '强'
  }
  return strengthMap[strength] || '弱'
}

const handleRegister = async () => {
  if (!registerFormRef.value) return
  
  registerFormRef.value.validate(async (valid: boolean) => {
    if (valid) {
      registerLoading.value = true
      try {
        // 调用注册API
        await axios.post('/users/register', {
          username: registerForm.value.username,
          name: registerForm.value.name,
          email: registerForm.value.email,
          phone: registerForm.value.phone,
          department: registerForm.value.department,
          position: registerForm.value.position,
          password: registerForm.value.password
        })
        
        registerLoading.value = false
        ElMessage.success('注册成功！您可以使用新账号登录了')
        // 清空注册表单
        registerForm.value = {
          username: '',
          name: '',
          email: '',
          phone: '',
          department: '',
          position: '',
          password: '',
          confirmPassword: ''
        }
        // 关闭注册对话框
        showRegister.value = false
      } catch (error: any) {
        registerLoading.value = false
        ElMessage.error(error.response?.data?.error || '注册失败，请稍后重试')
      }
    }
  })
}

// 忘记密码处理函数
const handleForgotPasswordNext = async () => {
  if (!forgotPasswordFormRef.value) return
  
  forgotPasswordFormRef.value.validate(async (valid: boolean) => {
    if (valid) {
      forgotPasswordLoading.value = true
      try {
        if (forgotPasswordStep.value === 1) {
          // 步骤1: 验证邮箱
          await authAPI.verifyEmail(forgotPasswordForm.value.email)
          ElMessage.success('邮箱验证通过')
          forgotPasswordStep.value = 2
        } else if (forgotPasswordStep.value === 2) {
          // 步骤2: 验证用户名并重置密码
          await authAPI.resetPassword(forgotPasswordForm.value.email, forgotPasswordForm.value.username)
          ElMessage.success('密码重置成功，新密码已发送到您的邮箱')
          forgotPasswordStep.value = 3
        }
      } catch (error: any) {
        ElMessage.error(error.response?.data?.error || '操作失败，请稍后重试')
      } finally {
        forgotPasswordLoading.value = false
      }
    }
  })
}

// 处理忘记密码取消
const handleForgotPasswordCancel = () => {
  // 重置表单和步骤
  resetForgotPasswordForm()
  showForgotPassword.value = false
}

// 处理忘记密码关闭
const handleForgotPasswordClose = () => {
  // 重置表单和步骤
  resetForgotPasswordForm()
  showForgotPassword.value = false
}

// 重置忘记密码表单
const resetForgotPasswordForm = () => {
  forgotPasswordStep.value = 1
  forgotPasswordForm.value = {
    email: '',
    username: ''
  }
  if (forgotPasswordFormRef.value) {
    forgotPasswordFormRef.value.resetFields()
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  padding: 0;
  margin: 0;
  overflow: hidden;
  position: relative;
  transition: all 0.3s ease;
  width: 100vw;
  height: 100vh;
  display: flex;
  align-items: stretch;
}

/* 风格切换开关 */
.style-switcher {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 1000;
  background: rgba(255, 255, 255, 0.9);
  padding: 10px;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
}

/* 主题色变量 */
:root {
  --primary-color: #409eff;
}

/* 风格1: 左右布局 */
.login-style-1 {
  display: flex;
  align-items: stretch;
  background: linear-gradient(135deg, var(--primary-color) 0%, #764ba2 100%);
  width: 100%;
  height: 100vh;
  overflow: hidden;
}

.login-style-1::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: radial-gradient(circle at 50% 50%, rgba(255, 255, 255, 0.1) 0%, rgba(255, 255, 255, 0) 70%);
  pointer-events: none;
}

/* 左侧装饰区域 */
.login-decoration {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.1) 0%, rgba(255, 255, 255, 0.05) 100%);
  backdrop-filter: blur(10px);
  border-right: 1px solid rgba(255, 255, 255, 0.1);
  flex: 1;
  animation: fadeInLeft 1s ease-out;
  width: 50%;
}

.decoration-content {
  color: white;
  text-align: center;
  max-width: 100%;
  padding: 0 20px;
}

.decoration-icon {
  font-size: 64px;
  margin-bottom: 20px;
  animation: pulse 2s infinite;
}

/* Logo图片样式 */
.decoration-logo {
  margin-bottom: 20px;
  animation: pulse 2s infinite;
}

.logo-image {
  width: 120px;
  height: 120px;
  object-fit: contain;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  padding: 10px;
  backdrop-filter: blur(10px);
}

.decoration-title {
  font-size: 48px;
  font-weight: 700;
  margin-bottom: 16px;
  background: linear-gradient(135deg, #ffffff 0%, #f0f0f0 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.decoration-subtitle {
  font-size: 18px;
  margin-bottom: 40px;
  opacity: 0.9;
  line-height: 1.6;
}

.decoration-features {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  align-items: center;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 16px;
  opacity: 0.9;
  transition: all 0.3s ease;
}

.feature-item:hover {
  transform: translateX(10px);
  opacity: 1;
}

.feature-icon {
  font-size: 20px;
}

/* 右侧表单区域 */
.login-form-section {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
  background: white;
  flex: 1;
  animation: fadeInRight 1s ease-out;
  width: 50%;
}

.form-wrapper {
  width: 100%;
  max-width: 500px;
  padding: 0 20px;
}

.login-card {
  border-radius: 16px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  border: none;
  transition: all 0.3s ease;
}

.login-card:hover {
  box-shadow: 0 30px 80px rgba(0, 0, 0, 0.15);
  transform: translateY(-5px);
}

.card-header {
  padding: 32px 32px 20px;
  text-align: center;
}

.login-title {
  font-size: 28px;
  font-weight: 700;
  color: #333;
  margin: 0 0 8px 0;
}

.login-subtitle {
  font-size: 14px;
  color: #909399;
  margin: 0;
}

:deep(.el-card__body) {
  padding: 0 32px 32px;
}

:deep(.el-form) {
  margin-top: 32px;
}

:deep(.el-form-item) {
  margin-bottom: 28px;
  transition: all 0.3s ease;
}

:deep(.el-form-item--error) {
  animation: shake 0.5s ease-in-out;
}

.form-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding: 0;
}

:deep(.el-form-item__label) {
  font-weight: 600;
  color: #606266;
  font-size: 14px;
  margin-bottom: 8px;
}

:deep(.el-input__wrapper) {
  border-radius: 12px;
  transition: all 0.3s ease;
  border: 2px solid transparent;
  background-color: #f5f7fa;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

:deep(.el-input__wrapper:hover) {
  background-color: #ecf5ff;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
}

:deep(.el-input__wrapper.is-focus) {
  background-color: #ffffff;
  border-color: #409eff;
  box-shadow: 0 0 0 4px rgba(64, 158, 255, 0.15);
}

:deep(.el-input__inner) {
  padding: 14px 16px;
  font-size: 15px;
  border: none;
  background-color: transparent;
}

:deep(.el-input__prefix-inner),
:deep(.el-input__suffix-inner) {
  padding: 0 8px;
}

:deep(.el-input-group__prepend) {
  background-color: #409eff;
  border-right: none;
  border-top-left-radius: 10px;
  border-bottom-left-radius: 10px;
  border: 2px solid #409eff;
  border-right: none;
  color: white;
}

:deep(.el-input-group__prepend .el-icon) {
  font-size: 18px;
  color: white;
}

:deep(.el-input__outer-wrapper:focus-within .el-input-group__prepend) {
  background-color: #667eea;
  border-color: #667eea;
}

.form-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding: 0;
}

.remember-me {
  font-size: 14px;
  color: #606266;
}

:deep(.el-checkbox__label) {
  color: #606266;
}

.register-link {
  font-size: 14px;
  font-weight: 500;
}

.login-button {
  padding: 18px 24px;
  font-size: 18px;
  font-weight: 700;
  border-radius: 12px;
  background-color: var(--primary-color);
  border: none;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  color: white;
  letter-spacing: 1px;
  margin-top: 28px;
  height: 56px;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  width: 100%;
  text-align: center;
}

.login-button:hover:not(:disabled) {
  background-color: #667eea;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
  transform: translateY(-1px);
}

.login-button:active:not(:disabled) {
  background-color: #337ecc;
  box-shadow: 0 2px 4px rgba(64, 158, 255, 0.2);
  transform: translateY(0);
}

.login-button:disabled {
  opacity: 0.6;
  transform: none;
  box-shadow: none;
}

.login-button::before {
  content: none;
}

.login-footer {
  margin-top: 24px;
  text-align: center;
}

.footer-text {
  font-size: 12px;
  color: #909399;
  margin: 0;
}

/* 风格2: 居中卡片布局 */
.login-style-2 {
  display: flex;
  justify-content: center;
  align-items: center;
  background: #f0f2f5;
  position: relative;
  overflow: hidden;
  width: 100vw;
  height: 100vh;
  margin: 0;
  padding: 0;
}

.style-2-background {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 1;
}

.style-2-gradient {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, var(--primary-color) 0%, #00f2fe 100%);
  opacity: 0.8;
}

.style-2-pattern {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-image: 
    radial-gradient(circle at 25% 25%, rgba(255, 255, 255, 0.1) 0%, transparent 50%),
    radial-gradient(circle at 75% 75%, rgba(255, 255, 255, 0.1) 0%, transparent 50%);
  background-size: 100px 100px;
  animation: patternMove 20s linear infinite;
}

@keyframes patternMove {
  0% {
    background-position: 0 0;
  }
  100% {
    background-position: 100px 100px;
  }
}

.style-2-content {
  position: relative;
  z-index: 2;
  text-align: center;
  padding: 40px;
  width: 100%;
  max-width: 450px;
}

.style-2-logo {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  margin-bottom: 32px;
  animation: fadeInDown 0.8s ease-out;
}

.style-2-logo-icon {
  font-size: 48px;
  color: white;
  background: rgba(255, 255, 255, 0.2);
  padding: 16px;
  border-radius: 50%;
  backdrop-filter: blur(10px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
}

.style-2-title {
  font-size: 32px;
  font-weight: 700;
  color: white;
  margin: 0;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.style-2-card {
  background: white;
  border-radius: 20px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
  border: none;
  overflow: hidden;
  transition: all 0.3s ease;
  animation: fadeInUp 0.8s ease-out 0.2s both;
}

.style-2-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 25px 70px rgba(0, 0, 0, 0.2);
}

.style-2-card-header {
  text-align: center;
  padding: 32px 0 16px;
}

.style-2-login-title {
  font-size: 28px;
  font-weight: 700;
  color: #303133;
  margin: 0 0 8px;
}

.style-2-login-subtitle {
  font-size: 14px;
  color: #909399;
  margin: 0;
}

.style-2-input {
  border-radius: 12px;
  transition: all 0.3s ease;
}

.style-2-input :deep(.el-input__wrapper) {
  border-radius: 12px;
  border: 2px solid #f0f2f5;
  background: #fafafa;
  transition: all 0.3s ease;
}

.style-2-input :deep(.el-input__wrapper:hover) {
  border-color: var(--primary-color);
  background: #ffffff;
  box-shadow: 0 4px 12px rgba(var(--primary-color), 0.1);
}

.style-2-input :deep(.el-input__wrapper.is-focus) {
  border-color: var(--primary-color);
  box-shadow: 0 0 0 4px rgba(var(--primary-color), 0.15);
  background: #ffffff;
  transform: translateY(-2px);
}

.style-2-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding: 0;
}

.style-2-remember {
  font-size: 14px;
  color: #606266;
}

.style-2-register {
  font-size: 14px;
  font-weight: 500;
}

.style-2-button {
  padding: 16px 20px;
  font-size: 18px;
  font-weight: 700;
  border-radius: 12px;
  background: linear-gradient(135deg, var(--primary-color) 0%, #00f2fe 100%);
  border: none;
  transition: all 0.3s ease;
  box-shadow: 0 6px 18px rgba(var(--primary-color), 0.4);
  color: white;
  letter-spacing: 0.5px;
  margin-top: 24px;
  height: 56px;
  position: relative;
  overflow: hidden;
}

.style-2-button::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.5s ease;
}

.style-2-button:hover:not(:disabled)::before {
  left: 100%;
}

.style-2-button:hover:not(:disabled) {
  transform: translateY(-3px);
  box-shadow: 0 10px 25px rgba(79, 172, 254, 0.5);
  background: linear-gradient(135deg, #00f2fe 0%, var(--primary-color) 100%);
}

.style-2-button:active:not(:disabled) {
  transform: translateY(0);
  box-shadow: 0 4px 12px rgba(79, 172, 254, 0.4);
}

.style-2-footer {
  margin-top: 24px;
  text-align: center;
  animation: fadeInUp 0.8s ease-out 0.4s both;
}

.style-2-copyright {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.8);
  margin: 0;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

/* 风格3: 深色主题 */
.login-style-3 {
  display: flex;
  justify-content: center;
  align-items: center;
  background: #0f172a;
  position: relative;
  overflow: hidden;
  width: 100vw;
  height: 100vh;
  margin: 0;
  padding: 0;
}

.style-3-background {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 1;
}

.style-3-gradient {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: radial-gradient(circle at 50% 50%, var(--primary-color) 0%, #0f172a 70%);
  opacity: 0.3;
}

.style-3-grid {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-image: 
    linear-gradient(rgba(79, 172, 254, 0.1) 1px, transparent 1px),
    linear-gradient(90deg, rgba(79, 172, 254, 0.1) 1px, transparent 1px);
  background-size: 50px 50px;
  animation: gridMove 20s linear infinite;
}

@keyframes gridMove {
  0% {
    background-position: 0 0;
  }
  100% {
    background-position: 50px 50px;
  }
}

.style-3-content {
  position: relative;
  z-index: 2;
  text-align: center;
  padding: 40px;
  width: 100%;
  max-width: 450px;
}

.style-3-logo {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  margin-bottom: 32px;
  animation: fadeInDown 0.8s ease-out;
}

.style-3-logo-icon {
  font-size: 48px;
  color: var(--primary-color);
  background: rgba(var(--primary-color), 0.1);
  padding: 16px;
  border-radius: 50%;
  backdrop-filter: blur(10px);
  box-shadow: 0 8px 24px rgba(var(--primary-color), 0.2);
  border: 2px solid rgba(var(--primary-color), 0.3);
}

.style-3-title {
  font-size: 32px;
  font-weight: 700;
  color: #ffffff;
  margin: 0;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

.style-3-card {
  background: rgba(15, 23, 42, 0.8);
  border: 1px solid rgba(79, 172, 254, 0.3);
  border-radius: 20px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(10px);
  overflow: hidden;
  transition: all 0.3s ease;
  animation: fadeInUp 0.8s ease-out 0.2s both;
}

.style-3-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 25px 70px rgba(0, 0, 0, 0.4);
  border-color: rgba(79, 172, 254, 0.5);
}

.style-3-card-header {
  text-align: center;
  padding: 32px 0 16px;
}

.style-3-login-title {
  font-size: 28px;
  font-weight: 700;
  color: #ffffff;
  margin: 0 0 8px;
}

.style-3-login-subtitle {
  font-size: 14px;
  color: #94a3b8;
  margin: 0;
}

.style-3-form :deep(.el-form-item__label) {
  color: #cbd5e1;
}

.style-3-input {
  border-radius: 12px;
  transition: all 0.3s ease;
}

.style-3-input :deep(.el-input__wrapper) {
  border-radius: 12px;
  border: 2px solid rgba(79, 172, 254, 0.3);
  background: rgba(30, 41, 59, 0.8);
  transition: all 0.3s ease;
}

.style-3-input :deep(.el-input__inner) {
  color: #ffffff;
  background: transparent;
}

.style-3-input :deep(.el-input__wrapper:hover) {
  border-color: var(--primary-color);
  background: rgba(30, 41, 59, 1);
  box-shadow: 0 4px 12px rgba(var(--primary-color), 0.2);
}

.style-3-input :deep(.el-input__wrapper.is-focus) {
  border-color: #4facfe;
  box-shadow: 0 0 0 4px rgba(var(--primary-color), 0.2);
  background: rgba(30, 41, 59, 1);
  transform: translateY(-2px);
}

.style-3-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding: 0;
}

.style-3-remember {
  font-size: 14px;
  color: #94a3b8;
}

.style-3-remember :deep(.el-checkbox__label) {
  color: #94a3b8;
}

.style-3-register {
  font-size: 14px;
  font-weight: 500;
}

.style-3-button {
  padding: 16px 20px;
  font-size: 18px;
  font-weight: 700;
  border-radius: 12px;
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  border: none;
  transition: all 0.3s ease;
  box-shadow: 0 6px 18px rgba(79, 172, 254, 0.4);
  color: white;
  letter-spacing: 0.5px;
  margin-top: 24px;
  height: 56px;
  position: relative;
  overflow: hidden;
}

.style-3-button::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.5s ease;
}

.style-3-button:hover:not(:disabled)::before {
  left: 100%;
}

.style-3-button:hover:not(:disabled) {
  transform: translateY(-3px);
  box-shadow: 0 10px 25px rgba(79, 172, 254, 0.5);
  background: linear-gradient(135deg, #00f2fe 0%, var(--primary-color) 100%);
}

.style-3-button:active:not(:disabled) {
  transform: translateY(0);
  box-shadow: 0 4px 12px rgba(79, 172, 254, 0.4);
}

.style-3-footer {
  margin-top: 24px;
  text-align: center;
  animation: fadeInUp 0.8s ease-out 0.4s both;
}

.style-3-copyright {
  font-size: 12px;
  color: rgba(148, 163, 184, 0.7);
  margin: 0;
}

/* 风格4: 简约设计 */
.login-style-4 {
  display: flex;
  justify-content: center;
  align-items: center;
  background: #f8fafc;
  position: relative;
  overflow: hidden;
  width: 100vw;
  height: 100vh;
  margin: 0;
  padding: 0;
}

.style-4-background {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 1;
}

.style-4-shape {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.6;
  animation: shapeFloat 20s ease-in-out infinite;
}

.style-4-shape-1 {
  width: 300px;
  height: 300px;
  background: var(--primary-color);
  top: -100px;
  left: -100px;
  animation-delay: 0s;
}

.style-4-shape-2 {
  width: 250px;
  height: 250px;
  background: #00f2fe;
  bottom: -100px;
  right: -100px;
  animation-delay: 5s;
}

.style-4-shape-3 {
  width: 200px;
  height: 200px;
  background: #667eea;
  top: 50%;
  right: 10%;
  animation-delay: 10s;
}

@keyframes shapeFloat {
  0%, 100% {
    transform: translate(0, 0) rotate(0deg);
  }
  33% {
    transform: translate(30px, -30px) rotate(120deg);
  }
  66% {
    transform: translate(-30px, 30px) rotate(240deg);
  }
}

.style-4-content {
  position: relative;
  z-index: 2;
  text-align: center;
  padding: 40px;
  width: 100%;
  max-width: 420px;
}

.style-4-card {
  background: white;
  border-radius: 16px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.08);
  border: none;
  overflow: hidden;
  transition: all 0.3s ease;
  animation: fadeInUp 0.8s ease-out;
}

.style-4-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 15px 50px rgba(0, 0, 0, 0.12);
}

.style-4-logo {
  margin: 40px 0 20px;
  animation: pulse 2s infinite;
}

.style-4-logo-icon {
  font-size: 64px;
  color: #4facfe;
}

.style-4-title {
  font-size: 32px;
  font-weight: 700;
  color: #1e293b;
  margin: 0 0 8px;
}

.style-4-subtitle {
  font-size: 14px;
  color: #64748b;
  margin: 0 0 40px;
}

.style-4-form :deep(.el-form-item) {
  margin-bottom: 24px;
}

.style-4-form :deep(.el-input__wrapper) {
  border-radius: 10px;
  border: 1px solid #e2e8f0;
  background: #f8fafc;
  transition: all 0.3s ease;
}

.style-4-form :deep(.el-input__wrapper:hover) {
  border-color: #4facfe;
  background: #ffffff;
  box-shadow: 0 2px 8px rgba(79, 172, 254, 0.1);
}

.style-4-form :deep(.el-input__wrapper.is-focus) {
  border-color: #4facfe;
  box-shadow: 0 0 0 3px rgba(79, 172, 254, 0.1);
  background: #ffffff;
  transform: translateY(-1px);
}

.style-4-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding: 0;
}

.style-4-remember {
  font-size: 14px;
  color: #64748b;
}

.style-4-register {
  font-size: 14px;
  font-weight: 500;
}

.style-4-button {
  padding: 16px 20px;
  font-size: 16px;
  font-weight: 600;
  border-radius: 10px;
  background: #4facfe;
  border: none;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(var(--primary-color), 0.3);
  color: white;
  margin-top: 20px;
  height: 52px;
}

.style-4-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(79, 172, 254, 0.4);
  background: var(--primary-color);
}

.style-4-button:active:not(:disabled) {
  transform: translateY(0);
  box-shadow: 0 2px 8px rgba(79, 172, 254, 0.3);
}

.style-4-footer {
  margin-top: 32px;
  text-align: center;
}

.style-4-copyright {
  font-size: 12px;
  color: #94a3b8;
  margin: 0;
}

/* 风格5: 渐变卡片 */
.login-style-5 {
  display: flex;
  justify-content: center;
  align-items: center;
  background: #f0f9ff;
  position: relative;
  overflow: hidden;
  width: 100vw;
  height: 100vh;
  margin: 0;
  padding: 0;
}

.style-5-background {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 1;
}

.style-5-gradient {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, #f0f9ff 0%, #e0f2fe 50%, #bae6fd 100%);
}

.style-5-content {
  position: relative;
  z-index: 2;
  text-align: center;
  padding: 40px;
  width: 100%;
  max-width: 480px;
}

.style-5-logo {
  margin-bottom: 16px;
  animation: fadeInDown 0.8s ease-out;
}

.style-5-logo-icon {
  font-size: 48px;
  color: #0ea5e9;
}

.style-5-title {
  font-size: 36px;
  font-weight: 700;
  color: #0c4a6e;
  margin: 0 0 8px;
  animation: fadeInDown 0.8s ease-out 0.1s both;
}

.style-5-subtitle {
  font-size: 16px;
  color: #0ea5e9;
  margin: 0 0 40px;
  animation: fadeInDown 0.8s ease-out 0.2s both;
}

.style-5-card {
  position: relative;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.1) 0%, rgba(255, 255, 255, 0.05) 100%);
  border-radius: 24px;
  padding: 8px;
  backdrop-filter: blur(10px);
  animation: fadeInUp 0.8s ease-out 0.3s both;
}

.style-5-card-inner {
  background: white;
  border-radius: 20px;
  padding: 40px;
  box-shadow: 0 20px 60px rgba(12, 74, 110, 0.15);
  transition: all 0.3s ease;
}

.style-5-card:hover .style-5-card-inner {
  transform: translateY(-2px);
  box-shadow: 0 25px 70px rgba(12, 74, 110, 0.2);
}

.style-5-card-header {
  margin-bottom: 32px;
}

.style-5-login-title {
  font-size: 24px;
  font-weight: 700;
  color: #0c4a6e;
  margin: 0;
}

.style-5-form :deep(.el-form-item) {
  margin-bottom: 28px;
}

.style-5-input {
  border-radius: 12px;
  transition: all 0.3s ease;
}

.style-5-input :deep(.el-input__wrapper) {
  border-radius: 12px;
  border: 2px solid #e0f2fe;
  background: #f0f9ff;
  transition: all 0.3s ease;
}

.style-5-input :deep(.el-input__wrapper:hover) {
  border-color: #0ea5e9;
  background: #ffffff;
  box-shadow: 0 4px 12px rgba(14, 165, 233, 0.1);
}

.style-5-input :deep(.el-input__wrapper.is-focus) {
  border-color: #0ea5e9;
  box-shadow: 0 0 0 4px rgba(14, 165, 233, 0.1);
  background: #ffffff;
  transform: translateY(-2px);
}

.style-5-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 28px;
}

.style-5-remember {
  font-size: 14px;
  color: #0c4a6e;
}

.style-5-register {
  font-size: 14px;
  font-weight: 500;
  color: #0ea5e9;
}

.style-5-button {
  padding: 18px 24px;
  font-size: 18px;
  font-weight: 700;
  border-radius: 12px;
  background: linear-gradient(135deg, var(--primary-color) 0%, #0284c7 100%);
  border: none;
  transition: all 0.3s ease;
  box-shadow: 0 6px 18px rgba(var(--primary-color), 0.4);
  color: white;
  letter-spacing: 0.5px;
  height: 56px;
  position: relative;
  overflow: hidden;
}

.style-5-button::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.5s ease;
}

.style-5-button:hover:not(:disabled)::before {
  left: 100%;
}

.style-5-button:hover:not(:disabled) {
  transform: translateY(-3px);
  box-shadow: 0 10px 25px rgba(14, 165, 233, 0.5);
  background: linear-gradient(135deg, #0284c7 0%, var(--primary-color) 100%);
}

.style-5-button:active:not(:disabled) {
  transform: translateY(0);
  box-shadow: 0 4px 12px rgba(14, 165, 233, 0.4);
}

.style-5-footer {
  margin-top: 32px;
  text-align: center;
  animation: fadeInUp 0.8s ease-out 0.4s both;
}

.style-5-copyright {
  font-size: 12px;
  color: #64748b;
  margin: 0;
}

/* 动画效果 */
@keyframes fadeInLeft {
  from {
    opacity: 0;
    transform: translateX(-50px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

@keyframes fadeInRight {
  from {
    opacity: 0;
    transform: translateX(50px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fadeInDown {
  from {
    opacity: 0;
    transform: translateY(-30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
    opacity: 1;
  }
  50% {
    transform: scale(1.1);
    opacity: 0.9;
  }
}

@keyframes shake {
  0%, 100% {
    transform: translateX(0);
  }
  10%, 30%, 50%, 70%, 90% {
    transform: translateX(-5px);
  }
  20%, 40%, 60%, 80% {
    transform: translateX(5px);
  }
}

/* 密码强度样式 */
.password-strength {
  display: flex;
  align-items: center;
  margin-top: 8px;
  font-size: 12px;
}

.strength-label {
  color: #606266;
  margin-right: 8px;
}

.strength-bar {
  flex: 1;
  height: 4px;
  margin: 0 8px;
  border-radius: 2px;
  transition: all 0.3s ease;
}

.strength-bar.weak {
  background-color: #f56c6c;
  width: 30%;
}

.strength-bar.medium {
  background-color: #e6a23c;
  width: 60%;
}

.strength-bar.strong {
  background-color: #67c23a;
  width: 100%;
}

.strength-text {
  color: #606266;
  font-weight: 500;
  margin-left: 8px;
}

/* 对话框样式 */
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 20px;
}

:deep(.el-dialog__header) {
  padding: 24px;
  border-bottom: 1px solid #ebeef5;
}

:deep(.el-dialog__title) {
  font-size: 20px;
  font-weight: 600;
  color: #303133;
}

:deep(.el-dialog__body) {
  padding: 24px;
}

:deep(.el-dialog__footer) {
  padding: 16px 24px;
  border-top: 1px solid #ebeef5;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .login-style-1 {
    flex-direction: column;
  }
  
  .login-decoration {
    border-right: none;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    padding: 30px 20px;
  }
  
  .decoration-title {
    font-size: 36px;
  }
  
  .decoration-icon {
    font-size: 48px;
  }
  
  .decoration-features {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .login-form-section {
    padding: 30px 20px;
  }
  
  .form-wrapper {
    max-width: 100%;
  }
  
  :deep(.el-card__body) {
    padding: 0 20px 20px;
  }
  
  .style-2-content,
  .style-3-content,
  .style-4-content,
  .style-5-content {
    padding: 20px;
  }
  
  .style-5-card-inner {
    padding: 30px 20px;
  }
}
/* 分隔线样式 */
.divider-line {
  display: flex;
  align-items: center;
  margin: 24px 0;
  position: relative;
}

.divider-line::before,
.divider-line::after {
  content: '';
  flex: 1;
  height: 1px;
  background-color: #ebeef5;
}

.divider-text {
  padding: 0 16px;
  color: #909399;
  font-size: 14px;
}

/* 注册按钮样式 */
.register-button {
  margin-top: 16px;
  height: 56px;
  font-size: 18px;
  font-weight: 600;
  border-radius: 12px;
  transition: all 0.3s ease;
  background-color: white;
  border: 2px solid var(--primary-color);
  color: var(--primary-color);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  text-align: center;
  width: 100%;
}

.register-button:hover:not(:disabled) {
  background-color: #ecf5ff;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.2);
  transform: translateY(-1px);
}

.register-button:active:not(:disabled) {
  background-color: #d9ecff;
  box-shadow: 0 2px 4px rgba(64, 158, 255, 0.15);
  transform: translateY(0);
}

/* 其他风格的注册按钮样式保持不变 */
.style-2-register-button,
.style-3-register-button,
.style-4-register-button,
.style-5-register-button {
  margin-top: 16px;
  height: 56px;
  font-size: 18px;
  font-weight: 600;
  border-radius: 12px;
  transition: all 0.3s ease;
}

/* 登录按钮样式调整 */
.login-button,
.style-2-button,
.style-3-button,
.style-4-button,
.style-5-button {
  margin-bottom: 16px;
}
</style>