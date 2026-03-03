<template>
  <div class="app-container">
    <!-- 检查当前路由是否为登录页面 -->
    <template v-if="$route.name === 'login'">
      <!-- 登录页面只显示路由内容 -->
      <router-view />
    </template>
    <template v-else>
      <!-- 非登录页面显示完整布局 -->
      <el-container class="layout-container" direction="vertical">
        <!-- 顶部导航 -->
        <el-header class="app-header">
          <div class="header-content">
            <div class="logo-area">
              <img v-if="homeIcon" :src="homeIcon" alt="系统图标" class="custom-icon" />
              <el-icon v-else size="28"><Monitor /></el-icon>
              <h1 class="app-title">{{ homeTitle || 'DeanTech' }}</h1>
            </div>
            <div class="header-actions">
              <!-- 中英文切换 -->
              <el-tooltip content="中英文切换" placement="bottom">
                <el-button type="text" size="small" @click="toggleLanguage" class="header-icon-btn">
                  <LanguageSwitchIcon class="highlight-icon" />
                </el-button>
              </el-tooltip>

              <!-- 通知中心 -->
              <el-tooltip content="通知中心" placement="bottom">
                <el-button type="text" size="small" @click="openNotification" class="header-icon-btn">
                  <el-icon size="18"><Bell /></el-icon>
                </el-button>
              </el-tooltip>
              <!-- 全屏切换 -->
              <el-tooltip content="全屏切换" placement="bottom">
                <el-button type="text" size="small" @click="toggleFullscreen" class="header-icon-btn">
                  <el-icon size="18" class="highlight-icon"><FullScreen /></el-icon>
                </el-button>
              </el-tooltip>
              <el-dropdown>
                <div class="user-info">
                  <el-avatar :size="32" :src="currentUser.avatar || '/images/03b0d39583f48206768a7534e55bcpng.png'" ></el-avatar>
                  <span class="user-name">{{ currentUser.username || '管理员' }}</span>
                  <el-icon class="arrow-icon"><ArrowDown /></el-icon>
                </div>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item @click="goToConfig('profile')">
                      <el-icon><User /></el-icon>
                      <span>个人中心</span>
                    </el-dropdown-item>
                    <el-dropdown-item divided @click="handleLogout">
                      <el-icon><SwitchButton /></el-icon>
                      <span class="logout-text">退出登录</span>
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </div>
        </el-header>
        
        <el-container class="main-container">
          <!-- 侧边菜单 -->
          <el-aside :width="isCollapse ? '64px' : '220px'" class="app-aside">
            <el-menu
              :default-active="activeMenu"
              class="aside-menu"
              text-color="#303133"
              active-text-color="#409eff"
              background-color="#ffffff"
              @select="handleMenuSelect"
              :collapse="isCollapse"
              collapse-transition
            >
              <el-menu-item index="1">
                <el-icon><DataBoard /></el-icon>
                <span>平台概览</span>
              </el-menu-item>
              
              <!-- 监控告警侧边栏 -->
              <el-sub-menu index="2">
                <template #title>
                  <el-icon><Bell /></el-icon>
                  <span>监控告警</span>
                </template>
                
                <el-menu-item index="2-1">
                  <el-icon><Bell /></el-icon>
                  <span>告警列表</span>
                </el-menu-item>
                <el-menu-item index="2-2">
                  <el-icon><DataAnalysis /></el-icon>
                  <span>告警规则</span>
                </el-menu-item>
                <el-menu-item index="2-3">
                  <el-icon><Message /></el-icon>
                  <span>告警介质</span>
                </el-menu-item>
                <el-menu-item index="2-4">
                  <el-icon><Document /></el-icon>
                  <span>告警模板</span>
                </el-menu-item>
                <el-menu-item index="2-5">
                  <el-icon><Timer /></el-icon>
                  <span>静默规则</span>
                </el-menu-item>
                <el-menu-item index="2-6">
                  <el-icon><Lock /></el-icon>
                  <span>抑制规则</span>
                </el-menu-item>
                <el-menu-item index="2-7">
                  <el-icon><DataAnalysis /></el-icon>
                  <span>Prometheus</span>
                </el-menu-item>
              </el-sub-menu>
              

              
              <!-- 集群管理侧边栏 -->
              <el-sub-menu index="3">
                <template #title>
                  <el-icon><DataBoard /></el-icon>
                  <span>集群管理</span>
                </template>
                
                <el-menu-item index="3-1">
                  <el-icon><DataBoard /></el-icon>
                  <span>仪表盘</span>
                </el-menu-item>
                <el-menu-item index="3-2">
                  <el-icon><DataBoard /></el-icon>
                  <span>集群列表</span>
                </el-menu-item>
                <el-menu-item index="3-3">
                  <el-icon><Connection /></el-icon>
                  <span>节点管理</span>
                </el-menu-item>
                <el-menu-item index="3-10">
                  <el-icon><Grid /></el-icon>
                  <span>命名空间</span>
                </el-menu-item>
                <el-menu-item index="3-4">
                  <el-icon><Box /></el-icon>
                  <span>Pod</span>
                </el-menu-item>
                <el-menu-item index="3-5">
                  <el-icon><Grid /></el-icon>
                  <span>工作负载</span>
                </el-menu-item>
                <el-menu-item index="3-6">
                  <el-icon><Link /></el-icon>
                  <span>服务</span>
                </el-menu-item>
                <el-menu-item index="3-7">
                  <el-icon><DocumentCopy /></el-icon>
                  <span>配置</span>
                </el-menu-item>
                <el-menu-item index="3-8">
                  <el-icon><Key /></el-icon>
                  <span>密钥</span>
                </el-menu-item>
                <el-menu-item index="3-9">
                  <el-icon><Link /></el-icon>
                  <span>路由</span>
                </el-menu-item>
                <el-menu-item index="3-11">
                  <el-icon><Box /></el-icon>
                  <span>存储卷声明</span>
                </el-menu-item>
                <el-menu-item index="3-12">
                  <el-icon><Box /></el-icon>
                  <span>存储卷</span>
                </el-menu-item>
                <el-menu-item index="3-13">
                  <el-icon><Box /></el-icon>
                  <span>存储类</span>
                </el-menu-item>
                <el-menu-item index="3-14">
                  <el-icon><Bell /></el-icon>
                  <span>事件记录</span>
                </el-menu-item>
                <el-menu-item index="3-15">
                  <el-icon><Document /></el-icon>
                  <span>资源配额</span>
                </el-menu-item>
                <el-menu-item index="3-16">
                  <el-icon><Document /></el-icon>
                  <span>限制范围</span>
                </el-menu-item>
              </el-sub-menu>
              
              <!-- 虚机管理侧边栏 -->
              <el-sub-menu index="4">
                <template #title>
                  <el-icon><Monitor /></el-icon>
                  <span>主机管理</span>
                </template>
                
                <el-menu-item index="4-1">
                  <el-icon><Monitor /></el-icon>
                  <span>虚机管理</span>
                </el-menu-item>
                <el-menu-item index="4-2">
                  <el-icon><Cpu /></el-icon>
                  <span>KVM管理</span>
                </el-menu-item>
              </el-sub-menu>
              <!-- 日志查询侧边栏 -->
              <el-sub-menu index="6">
                <template #title>
                  <el-icon><Message /></el-icon>
                  <span>日志查询</span>
                </template>
                
                <el-menu-item index="6-1">
                  <el-icon><Document /></el-icon>
                  <span>SLS日志查询</span>
                </el-menu-item>
                <el-menu-item index="6-2">
                  <el-icon><Setting /></el-icon>
                  <span>SLS配置管理</span>
                </el-menu-item>
              </el-sub-menu>
              
              <!-- 实用工具侧边栏 -->
              <el-sub-menu index="7">
                <template #title>
                  <el-icon><Tools /></el-icon>
                  <span>实用工具</span>
                </template>
                
                <el-menu-item index="7-1">
                  <el-icon><Connection /></el-icon>
                  <span>路由跟踪</span>
                </el-menu-item>
                <el-menu-item index="7-2">
                  <el-icon><Monitor /></el-icon>
                  <span>Ping测试</span>
                </el-menu-item>
                <el-menu-item index="7-3">
                  <el-icon><Link /></el-icon>
                  <span>Telnet测试</span>
                </el-menu-item>
                <el-menu-item index="7-4">
                  <el-icon><Key /></el-icon>
                  <span>SSL证书监控</span>
                </el-menu-item>
                <el-menu-item index="7-5">
                  <el-icon><DocumentCopy /></el-icon>
                  <span>JSON/Base64</span>
                </el-menu-item>
                <el-menu-item index="7-6">
                  <el-icon><Picture /></el-icon>
                  <span>图片压缩</span>
                </el-menu-item>
                <el-menu-item index="7-7">
                  <el-icon><Key /></el-icon>
                  <span>MD5加解密</span>
                </el-menu-item>
              </el-sub-menu>
              
              <!-- AI智能侧边栏 -->
              <el-sub-menu index="8">
                <template #title>
                  <el-icon><ChatLineRound /></el-icon>
                  <span>AI智能</span>
                </template>
                
                <el-menu-item index="8-1">
                  <el-icon><ChatLineRound /></el-icon>
                  <span>AI智能助手</span>
                </el-menu-item>
                <el-menu-item index="8-2">
                  <el-icon><DocumentCopy /></el-icon>
                  <span>AI代码助手</span>
                </el-menu-item>
                <el-menu-item index="8-3">
                  <el-icon><DataAnalysis /></el-icon>
                  <span>AI数据分析</span>
                </el-menu-item>
                <el-menu-item index="8-4">
                  <el-icon><Message /></el-icon>
                  <span>AI内容生成</span>
                </el-menu-item>
                <el-menu-item index="8-5">
                  <el-icon><Setting /></el-icon>
                  <span>AI模型管理</span>
                </el-menu-item>
              </el-sub-menu>
              
              <!-- Jenkins侧边栏 -->
              <el-sub-menu index="9">
                <template #title>
                  <el-icon><DataAnalysis /></el-icon>
                  <span>持续集成</span>
                </template>
                
                <el-menu-item index="9-1">
                  <el-icon><Timer /></el-icon>
                  <span>Jenkins</span>
                </el-menu-item>
              </el-sub-menu>
              
              <!-- 系统配置侧边栏 -->
              <el-sub-menu index="5">
                <template #title>
                  <el-icon><Setting /></el-icon>
                  <span>系统配置</span>
                </template>
                
                <el-menu-item index="5-1">
                  <el-icon><User /></el-icon>
                  <span>个人中心</span>
                </el-menu-item>
                <el-menu-item index="5-2">
                  <el-icon><UserFilled /></el-icon>
                  <span>用户管理</span>
                </el-menu-item>
                <el-menu-item index="5-3">
                  <el-icon><Brush /></el-icon>
                  <span>系统风格</span>
                </el-menu-item>
                <el-menu-item index="5-4">
                  <el-icon><Message /></el-icon>
                  <span>邮箱配置</span>
                </el-menu-item>
                <el-menu-item index="5-5">
                  <el-icon><Document /></el-icon>
                  <span>操作日志</span>
                </el-menu-item>
                <el-menu-item index="5-6">
                  <el-icon><User /></el-icon>
                  <span>登录日志</span>
                </el-menu-item>
              </el-sub-menu>
            </el-menu>
            <div class="aside-footer">
              <el-button type="text" @click="toggleCollapse" class="collapse-btn">
                <el-icon v-if="isCollapse"><Expand /></el-icon>
                <el-icon v-else><Fold /></el-icon>
              </el-button>
            </div>
          </el-aside>
          
          <!-- 主内容区域 -->
          <el-main class="app-main">
            <router-view />
          </el-main>
        </el-container>
        
        <!-- 页脚区域 -->
        <el-footer class="app-footer">
          <p class="footer-text">{{ homeFooter || '© 2026 Dean\'s Copyright. All rights reserved.' }}</p>
        </el-footer>
      </el-container>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, provide, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Bell,
  Message,
  Document,
  Timer,
  Lock,
  DataAnalysis,
  User,
  UserFilled,
  ArrowDown,
  Setting,
  ChatLineRound,
  DataBoard,
  Monitor,
  Connection,
  Box,
  Grid,
  Link,
  DocumentCopy,
  Key,
  Brush,
  Fold,
  Expand,
  Tools,
  Cpu,
  FullScreen,
  SwitchButton,
  Picture
} from '@element-plus/icons-vue'
import { userAPI, dashboardAPI, systemSettingAPI } from './api/api'
import LanguageSwitchIcon from './components/LanguageSwitchIcon.vue'
const router = useRouter()

// 路由到菜单索引的映射
  const routeToMenuMap: Record<string, string> = {
    '/dashboard': '1',
    '/alerts': '2-1',
    '/rules': '2-2',
    '/media': '2-3',
    '/templates': '2-4',
    '/silences': '2-5',
    '/inhibitions': '2-6',
    '/prometheus': '2-7',
    '/cluster-dashboard': '3-1',
    '/clusters': '3-2',
    '/nodes': '3-3',
    '/pods': '3-4',
    '/deployments': '3-5',
    '/services': '3-6',
    '/configmaps': '3-7',
    '/secrets': '3-8',
    '/ingresses': '3-9',
    '/hosts': '4-1',
    '/profile': '5-1',
    '/users': '5-2',
    '/config': '5-3',
    '/email-config': '5-4',
    '/logs/operation': '5-5',
    '/logs/login': '5-6',
    '/sls-logs': '6-1',
    '/sls-config': '6-2',
    '/router-trace': '7-1',
    '/ping-test': '7-2',
    '/telnet-test': '7-3',
    '/ssl-monitor': '7-4',
    '/json-serializer': '7-5',
    '/image-compress': '7-6',
    '/md5-encrypt': '7-7',
    '/ai-intelligence': '8-1',
    '/ai-code-assistant': '8-2',
    '/ai-data-analysis': '8-3',
    '/ai-content-generator': '8-4',
    '/ai-model-management': '8-5',
    '/namespaces': '3-10',
    '/kvm-hosts': '4-2',
    '/pvcs': '3-11',
    '/pvs': '3-12',
    '/storageclasses': '3-13',
    '/events': '3-14',
    '/resourcequotas': '3-15',
    '/limitranges': '3-16',
    '/jenkins': '9-1'
  }

// 菜单索引到路由的映射
  const menuToRouteMap: Record<string, string> = {
    '1': '/dashboard',
    '2-1': '/alerts',
    '2-2': '/rules',
    '2-3': '/media',
    '2-4': '/templates',
    '2-5': '/silences',
    '2-6': '/inhibitions',
    '2-7': '/prometheus',
    '3-1': '/cluster-dashboard',
    '3-2': '/clusters',
    '3-3': '/nodes',
    '3-4': '/pods',
    '3-5': '/deployments',
    '3-6': '/services',
    '3-7': '/configmaps',
    '3-8': '/secrets',
    '3-9': '/ingresses',
    '3-10': '/namespaces',
    '3-11': '/pvcs',
    '3-12': '/pvs',
    '3-13': '/storageclasses',
    '3-14': '/events',
    '3-15': '/resourcequotas',
    '3-16': '/limitranges',
    '4-1': '/hosts',
    '4-2': '/kvm-hosts',
    '5-1': '/profile',
    '5-2': '/users',
    '5-3': '/config?activeTab=system-style',
    '5-4': '/email-config',
    '5-5': '/logs/operation',
    '5-6': '/logs/login',
    '6-1': '/sls-logs',
    '6-2': '/sls-config',
    '7-1': '/router-trace',
    '7-2': '/ping-test',
    '7-3': '/telnet-test',
    '7-4': '/ssl-monitor',
    '7-5': '/json-serializer',
    '7-6': '/image-compress',
    '7-7': '/md5-encrypt',
    '8-1': '/ai-intelligence',
    '8-2': '/ai-code-assistant',
    '8-3': '/ai-data-analysis',
    '8-4': '/ai-content-generator',
    '8-5': '/ai-model-management',
    '9-1': '/jenkins'
  }

// 刷新状态，用于触发组件刷新
const refreshTrigger = ref(0)

// 侧边栏伸缩状态
const isCollapse = ref(false)

// 从localStorage获取保存的首页设置，避免初始加载闪烁
const getHomeSettingsFromLocalStorage = () => {
  const savedTitle = localStorage.getItem('home.title')
  const savedIcon = localStorage.getItem('home.icon')
  const savedFooter = localStorage.getItem('home.footer')
  return {
    title: savedTitle || 'DeanTech',
    icon: savedIcon || '',
    footer: savedFooter || ''
  }
}

// 首页标题、图标和页脚设置
const savedSettings = getHomeSettingsFromLocalStorage()
const homeTitle = ref(savedSettings.title)
const homeIcon = ref(savedSettings.icon)
const homeFooter = ref(savedSettings.footer)

// 切换侧边栏伸缩
const toggleCollapse = () => {
  isCollapse.value = !isCollapse.value
}

// 提供刷新功能给所有子组件
provide('refreshTrigger', refreshTrigger)

// 根据当前路由设置活动菜单
const setActiveMenuByRoute = () => {
  const currentPath = router.currentRoute.value.path
  return routeToMenuMap[currentPath] || '1'
}

// 当前用户信息
const currentUser = ref({
  username: localStorage.getItem('username') || '管理员',
  avatar: localStorage.getItem('user.avatar') || ''
})

// 初始化时直接根据当前路由计算activeMenu，避免初始加载闪烁
const activeMenu = ref(setActiveMenuByRoute())

const handleMenuSelect = (key: string) => {
  activeMenu.value = key
  // 根据菜单key跳转到对应路由
  if (menuToRouteMap[key]) {
    router.push(menuToRouteMap[key])
  }
}

// 监听路由变化，更新活动菜单和用户信息
router.afterEach(() => {
  activeMenu.value = setActiveMenuByRoute()
  
  // 如果从登录页面跳转过来，更新用户信息
  if (router.currentRoute.value.path !== '/login') {
    getCurrentUserInfo()
  }
})

// 创建BroadcastChannel用于跨标签页通信
let broadcastChannel: BroadcastChannel | null = null

// 组件挂载时获取当前用户信息和系统配置
onMounted(() => {
  getCurrentUserInfo()
  getSystemInfo()
  
  // 监听用户信息更新事件，用于个人中心修改信息后更新头部用户信息
  window.addEventListener('userInfoUpdated', getCurrentUserInfo)
  // 监听首页设置更新事件，用于系统设置修改后更新头部标题和图标
  window.addEventListener('homeSettingsUpdated', getSystemInfo)
  // 监听用户状态变更事件，用于处理用户被禁用时的强制退出
  window.addEventListener('userStatusChanged', handleUserStatusChanged as EventListener)
  // 监听全局状态检查事件，用于立即检查用户状态
  window.addEventListener('checkUserStatus', checkUserStatus as EventListener)
  
  // 初始化BroadcastChannel用于跨标签页通信
  broadcastChannel = new BroadcastChannel('user-status-channel')
  
  // 监听BroadcastChannel消息
  broadcastChannel.addEventListener('message', (event) => {
    console.log('收到BroadcastChannel消息:', event.data)
    const { type, data } = event.data
    if (type === 'userStatusChanged') {
      handleUserStatusChanged({ detail: data } as any)
    } else if (type === 'checkUserStatus') {
      checkUserStatus()
    }
  })
  
  // 定期检查当前用户状态，确保被禁用的用户能被强制退出
  // 每30秒检查一次，平衡检测及时性和服务器压力
  userStatusCheckInterval = window.setInterval(() => {
    checkUserStatus()
  }, 30000)
})

// 设置浏览器favicon
const setFavicon = (iconUrl: string) => {
  let link = document.querySelector<HTMLLinkElement>('link[rel~="icon"]')
  if (!link) {
    link = document.createElement('link')
    link.rel = 'icon'
    document.head.appendChild(link)
  }
  if (iconUrl) {
    link.href = iconUrl
  } else {
    // 如果没有设置图标，移除favicon
    link.remove()
  }
}

// 获取系统信息，用于更新浏览器标题、首页标题和图标
const getSystemInfo = async () => {
  let browserTitle = ''
  
  // 获取首页设置
  try {
    const response = await systemSettingAPI.list()
    const settings = response.data
    
    // 查找首页标题设置
    const homeTitleSetting = settings.find((setting: any) => setting.key === 'home.title')
    if (homeTitleSetting) {
      homeTitle.value = homeTitleSetting.value
      // 保存到localStorage
      localStorage.setItem('home.title', homeTitleSetting.value)
    }
    
    // 查找首页图标设置
    const homeIconSetting = settings.find((setting: any) => setting.key === 'home.icon')
    if (homeIconSetting) {
      homeIcon.value = homeIconSetting.value
      // 保存到localStorage
      localStorage.setItem('home.icon', homeIconSetting.value)
      // 设置浏览器favicon
      setFavicon(homeIconSetting.value)
    } else {
      // 移除favicon
      setFavicon('')
    }
    
    // 查找浏览器标题设置
    const homeBrowserTitleSetting = settings.find((setting: any) => setting.key === 'home.browserTitle')
    if (homeBrowserTitleSetting) {
      browserTitle = homeBrowserTitleSetting.value
      // 保存到localStorage
      localStorage.setItem('home.browserTitle', homeBrowserTitleSetting.value)
    }
    
    // 查找首页页脚设置
    const homeFooterSetting = settings.find((setting: any) => setting.key === 'home.footer')
    if (homeFooterSetting) {
      homeFooter.value = homeFooterSetting.value
      // 保存到localStorage
      localStorage.setItem('home.footer', homeFooterSetting.value)
    }
  } catch (error) {
    console.error('获取首页设置失败:', error)
    // 使用localStorage中的值或默认值
    const savedTitle = localStorage.getItem('home.title')
    const savedIcon = localStorage.getItem('home.icon')
    const savedFooter = localStorage.getItem('home.footer')
    browserTitle = localStorage.getItem('home.browserTitle') || ''
    homeTitle.value = savedTitle || 'DeanTech'
    homeIcon.value = savedIcon || ''
    homeFooter.value = savedFooter || ''
    // 设置浏览器favicon
    setFavicon(savedIcon || '')
  }
  
  // 更新浏览器标题
  if (browserTitle) {
    document.title = browserTitle
  } else {
    try {
      const response = await dashboardAPI.getSystemInfo()
      const systemData = response.data
      document.title = systemData.title
    } catch (error) {
      console.error('获取系统信息失败:', error)
      // 使用默认标题
      document.title = 'DeanTech 集群管理平台'
    }
  }
}

// 用户状态检查定时器
let userStatusCheckInterval: number | null = null

// 处理用户状态变更事件
const handleUserStatusChanged = (event: any) => {
  const { userId, status } = event.detail || { userId: null, status: null }
  // 如果是当前登录用户被禁用，强制退出
  const currentUserId = localStorage.getItem('userId')
  
  // 如果事件中包含userId，并且与当前登录用户ID匹配，执行强制退出
  if (currentUserId && userId !== null && parseInt(currentUserId) === userId && (status === 'inactive' || status === 'disabled')) {
    ElMessage.error('该用户已被禁用')
    handleLogout()
    return
  }
  
  // 如果事件中没有包含userId，或者与当前登录用户ID不匹配，执行全局状态检查
  checkUserStatus()
}

// 定期检查当前用户状态
  const checkUserStatus = async () => {
    // 只有在非登录页面且有token时才检查
    if (router.currentRoute.value.path !== '/login') {
      const token = localStorage.getItem('token')
      
      if (token) {
        try {
          // 直接调用用户信息API，确保每次都检查用户状态
          const userData = await userAPI.getCurrentUser()
          
          // 主动检查用户状态，如果返回的状态是inactive或被禁用，执行强制退出
          const userStatus = userData.status || userData.Status
          
          // 获取用户ID
          const userId = userData.id || userData.ID || ''
          
          // 只要检测到用户被禁用，就立即退出登录
          if (userStatus === 'inactive' || userStatus === 'disabled' || userStatus === '0') {
            ElMessage.error('该用户已被禁用')
            handleLogout()
            return // 退出函数，避免后续更新
          }
          
          // 如果状态正常，更新用户信息
          const newAvatar = userData.avatar || userData.Avatar || ''
          const username = userData.username || userData.Username || '管理员'
          
          currentUser.value = {
            username: username,
            avatar: newAvatar
          }
          
          // 更新localStorage
          localStorage.setItem('user.avatar', newAvatar)
          localStorage.setItem('username', username)
          localStorage.setItem('userId', userId.toString())
        } catch (error: any) {
          // 如果是401或403错误，直接退出登录
          if (error.response?.status === 401 || error.response?.status === 403 || error.response?.data?.error?.includes('禁用') || error.response?.data?.error?.includes('不存在')) {
            ElMessage.error('该用户已被禁用')
            handleLogout()
          }
        }
      } else {
        // 如果没有token，直接退出登录
        handleLogout()
      }
    }
  }

// 组件卸载时移除事件监听器和定时器
onBeforeUnmount(() => {
  window.removeEventListener('userInfoUpdated', getCurrentUserInfo)
  window.removeEventListener('homeSettingsUpdated', getSystemInfo)
  window.removeEventListener('userStatusChanged', handleUserStatusChanged as EventListener)
  window.removeEventListener('checkUserStatus', checkUserStatus as EventListener)
  
  // 关闭BroadcastChannel，避免内存泄漏
  if (broadcastChannel) {
    broadcastChannel.close()
    broadcastChannel = null
  }
  
  // 清除用户状态检查定时器
  if (userStatusCheckInterval) {
    clearInterval(userStatusCheckInterval)
  }
})

// 获取当前用户信息
const getCurrentUserInfo = async () => {
  try {
    const userData = await userAPI.getCurrentUser()
    const newAvatar = userData.avatar || userData.Avatar || ''
    const userId = userData.id || userData.ID || ''
    const username = userData.username || userData.Username || '管理员'
    
    currentUser.value = {
      username: username,
      avatar: newAvatar
    }
    
    // 将用户信息保存到localStorage，避免下次刷新页面时出现默认头像
    localStorage.setItem('user.avatar', newAvatar)
    localStorage.setItem('username', username)
    localStorage.setItem('userId', userId.toString())
  } catch (error: any) {
    console.error('获取当前用户信息失败:', error)
    
    // 检查是否是用户被禁用错误
    if (error.response?.data?.error === '用户已被禁用') {
      ElMessage.error('您的账号已被禁用，请联系管理员')
      // 执行登出操作
      localStorage.removeItem('token')
      localStorage.removeItem('username')
      localStorage.removeItem('user.avatar')
      localStorage.removeItem('userId')
      // 跳转到登录页面
      router.push('/login')
    }
  }
}



// 跳转到系统设置页面
const goToConfig = (tab: string) => {
  // 跳转到对应的独立页面
  if (tab === 'profile') {
    router.push('/profile')
  } else if (tab === 'users') {
    router.push('/users')
  }
}

// 退出登录
const handleLogout = () => {
  // 清除本地存储的用户信息和token
  localStorage.removeItem('token')
  localStorage.removeItem('username')
  localStorage.removeItem('user.avatar')
  localStorage.removeItem('userId')
  // 跳转到登录页面
  router.push('/login')
}

// 中英文切换
const toggleLanguage = () => {
  // 这里可以实现中英文切换逻辑
  ElMessage.info('中英文切换功能开发中')
}



// 通知中心
const openNotification = () => {
  // 这里可以实现通知中心逻辑
  ElMessage.info('通知中心功能开发中')
}

// 全屏切换
const toggleFullscreen = () => {
  // 实现全屏切换逻辑
  const element = document.documentElement
  if (!document.fullscreenElement) {
    element.requestFullscreen().catch(err => {
      ElMessage.error(`全屏切换失败: ${err.message}`)
    })
  } else {
    if (document.exitFullscreen) {
      document.exitFullscreen()
    }
  }
}
</script>

<style scoped>
.app-container {
  height: 100vh;
  width: 100vw;
  overflow: hidden;
  background-color: #f0f2f5;
}

.layout-container {
  height: 100%;
}

.app-header {
  display: flex;
  align-items: center;
  background-color: #ffffff;
  border-bottom: 1px solid #ebeef5;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  padding: 0 24px;
  height: 64px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.logo-area {
  display: flex;
  align-items: center;
  gap: 12px;
}

.app-title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #303133;
}

.custom-icon {
  width: 28px;
  height: 28px;
  object-fit: contain;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 4px;
}

/* 头部图标按钮样式 */
.header-icon-btn {
  --el-button-text-color: #606266;
  --el-button-hover-text-color: #3b82f6;
  --el-button-text-hover-bg-color: rgba(59, 130, 246, 0.12);
  padding: 8px;
  border-radius: 6px;
  transition: all 0.2s ease;
  min-width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.header-icon-btn:hover {
  background-color: rgba(59, 130, 246, 0.12);
}

.header-icon-btn .el-icon {
  color: #606266;
  transition: color 0.2s ease;
}

.header-icon-btn:hover .el-icon {
  color: #3b82f6;
}

/* 特定图标优化 */
.header-icon-btn .highlight-icon {
  color: #409eff !important;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
  transition: all 0.3s ease;
}

.user-info:hover {
  background-color: #f5f7fa;
}

.user-name {
  font-size: 14px;
  color: #303133;
}

.arrow-icon {
  font-size: 12px;
}

.main-container {
  flex: 1;
  background-color: #f0f2f5;
  overflow: hidden;
}

.app-aside {
  background-color: #ffffff;
  border-right: 1px solid #ebeef5;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.05);
  transition: width 0.3s ease;
  display: flex;
  flex-direction: column;
}

.aside-menu {
  background-color: #ffffff;
  border-right: none;
  flex: 1;
  overflow-y: auto;
}

.aside-footer {
  display: flex;
  justify-content: center;
  padding: 12px 8px;
  border-top: 1px solid #ebeef5;
}

.collapse-btn {
  font-size: 16px;
  color: #606266;
  transition: color 0.3s ease;
}

.collapse-btn:hover {
  color: #3b82f6;
}

/* 修复监控告警选项的文字对齐问题 */
.aside-menu .el-sub-menu__icon-arrow {
  margin-left: auto;
}

/* 确保子菜单和普通菜单项内部内容对齐 */
.aside-menu .el-sub-menu__title,
.aside-menu .el-menu-item {
  display: flex;
  align-items: center;
}

/* 统一菜单项内部内容布局 */
.aside-menu .el-menu-item .el-icon,
.aside-menu .el-sub-menu__title .el-icon {
  margin-right: 12px;
  font-size: 18px;
  vertical-align: middle;
}

/* 确保菜单项内部元素对齐 */
.aside-menu .el-menu-item span,
.aside-menu .el-sub-menu__title span {
  flex: 1;
  line-height: inherit;
}

/* 统一菜单样式，确保子菜单与普通菜单项对齐 */
.aside-menu .el-menu-item,
.aside-menu .el-sub-menu__title {
  height: 52px;
  line-height: 52px;
  margin: 8px 16px !important;
  border-radius: 8px;
  transition: all 0.3s cubic-bezier(0.645, 0.045, 0.355, 1);
  font-size: 14px;
  font-weight: 500;
  color: #303133;
  padding: 0 20px !important;
  /* 重置所有可能影响对齐的样式 */
  text-indent: 0 !important;
  box-sizing: border-box;
  position: relative;
  left: 0 !important;
}

/* 强制重置el-menu-item的默认padding和margin - 一级菜单 */
.aside-menu :deep(.el-menu > .el-menu-item) {
  margin: 8px 16px !important;
  padding: 0 20px !important;
}

/* 强制重置el-sub-menu__title的默认padding和margin - 一级菜单 */
.aside-menu :deep(.el-menu > .el-sub-menu > .el-sub-menu__title) {
  margin: 8px 16px !important;
  padding: 0 20px !important;
}

/* 二级菜单样式 - 添加缩进效果 */
.aside-menu :deep(.el-menu .el-sub-menu .el-menu-item) {
  margin: 4px 16px 4px 24px !important;
  padding: 0 28px !important;
  font-size: 13px !important;
  font-weight: 400 !important;
}

/* 确保二级菜单图标也有合适的间距 */
.aside-menu :deep(.el-menu .el-sub-menu .el-menu-item .el-icon) {
  margin-right: 10px !important;
  font-size: 16px !important;
}

/* 确保所有菜单图标左对齐 */
.aside-menu :deep(.el-menu > .el-menu-item .el-icon),
.aside-menu :deep(.el-menu > .el-sub-menu > .el-sub-menu__title .el-icon) {
  margin-right: 12px !important;
  margin-left: 0 !important;
  font-size: 18px;
  vertical-align: middle;
}

/* 确保平台概览与其他一级菜单左对齐 */
.aside-menu :deep(.el-menu > .el-menu-item) {
  margin: 8px 16px !important;
  padding: 0 20px !important;
  left: 0 !important;
  text-indent: 0 !important;
  display: flex !important;
  align-items: center !important;
  justify-content: flex-start !important;
}

/* 确保所有一级菜单标题对齐一致 */
.aside-menu :deep(.el-menu > .el-sub-menu > .el-sub-menu__title),
.aside-menu :deep(.el-menu > .el-menu-item) {
  margin: 8px 16px !important;
  padding: 0 20px !important;
  display: flex !important;
  align-items: center !important;
  justify-content: flex-start !important;
}

/* 统一一级菜单图标间距 */
.aside-menu :deep(.el-menu > .el-sub-menu > .el-sub-menu__title .el-icon),
.aside-menu :deep(.el-menu > .el-menu-item .el-icon) {
  margin-right: 12px !important;
  margin-left: 0 !important;
  font-size: 18px;
  vertical-align: middle;
}

/* 统一菜单悬停样式 */
.aside-menu .el-menu-item:hover,
.aside-menu .el-sub-menu__title:hover {
  background-color: #f0f9ff;
  color: #3b82f6;
  box-shadow: 0 2px 8px rgba(59, 130, 246, 0.15);
  transform: translateX(4px);
}

/* 统一激活菜单样式 */
.aside-menu .el-menu-item.is-active,
.aside-menu .el-sub-menu__title.is-active {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #ffffff;
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
  border-left: none;
  transform: translateX(4px);
}

/* 激活菜单的左侧指示器 */
.aside-menu .el-menu-item.is-active::before,
.aside-menu .el-sub-menu__title.is-active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 4px;
  height: 24px;
  background-color: #ffffff;
  border-radius: 0 4px 4px 0;
}

/* 统一图标样式 */
.aside-menu .el-icon {
  margin-right: 12px;
  font-size: 18px;
  vertical-align: middle;
  width: 18px;
  height: 18px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

/* 确保所有一级菜单完全对齐 - 平台概览和监控告警等 */
.aside-menu :deep(.el-menu > .el-menu-item),
.aside-menu :deep(.el-menu > .el-sub-menu > .el-sub-menu__title) {
  height: 52px !important;
  line-height: 52px !important;
  margin: 8px 16px !important;
  padding: 0 20px !important;
  display: flex !important;
  align-items: center !important;
  justify-content: flex-start !important;
  box-sizing: border-box !important;
  position: relative !important;
  left: 0 !important;
  text-indent: 0 !important;
}

/* 确保所有一级菜单图标完全对齐 */
.aside-menu :deep(.el-menu > .el-menu-item .el-icon),
.aside-menu :deep(.el-menu > .el-sub-menu > .el-sub-menu__title .el-icon) {
  margin-right: 12px !important;
  margin-left: 0 !important;
  font-size: 18px !important;
  width: 18px !important;
  height: 18px !important;
  display: inline-flex !important;
  align-items: center !important;
  justify-content: center !important;
  vertical-align: middle !important;
}

/* 确保所有一级菜单文字对齐 */
.aside-menu :deep(.el-menu > .el-menu-item span),
.aside-menu :deep(.el-menu > .el-sub-menu > .el-sub-menu__title span) {
  font-size: 14px !important;
  font-weight: 500 !important;
  color: #303133 !important;
  vertical-align: middle !important;
}

/* 折叠状态下的样式调整 */
.aside-menu.el-menu--collapse .el-menu-item,
.aside-menu.el-menu--collapse .el-sub-menu__title {
  margin: 8px 8px;
  justify-content: center;
}

.aside-menu.el-menu--collapse .el-menu-item:hover,
.aside-menu.el-menu--collapse .el-sub-menu__title:hover {
  transform: none;
}

.aside-menu.el-menu--collapse .el-menu-item.is-active,
.aside-menu.el-menu--collapse .el-sub-menu__title.is-active {
  transform: none;
}

.aside-menu.el-menu--collapse .el-menu-item.is-active::before,
.aside-menu.el-menu--collapse .el-sub-menu__title.is-active::before {
  display: none;
}

.app-main {
  background-color: #f5f7fa;
  padding: 24px;
  overflow-y: auto;
  flex: 1;
  margin: 0 auto;
  max-width: 100%;
  box-sizing: border-box;
}

.app-footer {
  background-color: #ffffff;
  border-top: 1px solid #ebeef5;
  padding: 16px 24px;
  text-align: center;
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.05);
}

.footer-text {
  margin: 0;
  font-size: 12px;
  color: #909399;
}

/* 退出登录样式 */
.logout-text {
  color: #ff4d4f;
}
</style>