import { createRouter, createWebHistory } from 'vue-router'
import DashboardView from '../views/DashboardView.vue'
import AlertsView from '../views/AlertsView.vue'
import MediaView from '../views/MediaView.vue'
import TemplatesView from '../views/TemplatesView.vue'
import SilencesView from '../views/SilencesView.vue'
import InhibitionsView from '../views/InhibitionsView.vue'
import PrometheusView from '../views/PrometheusView.vue'
import ConfigView from '../views/ConfigView.vue'
import ProfileView from '../views/ProfileView.vue'
import UsersView from '../views/UsersView.vue'
import RulesView from '../views/RulesView.vue'
import LoginView from '../views/LoginView.vue'
import HostsView from '../views/HostsView.vue'
import ClustersView from '../views/ClustersView.vue'
import RouterTraceView from '../views/RouterTraceView.vue'
import PingTestView from '../views/PingTestView.vue'
import TelnetTestView from '../views/TelnetTestView.vue'
import KvmHostsView from '../views/KvmHostsView.vue'
import EmailConfigView from '../views/EmailConfigView.vue'
import OperationLogsView from '../views/OperationLogsView.vue'
import LoginLogsView from '../views/LoginLogsView.vue'
import NamespacesView from '../views/NamespacesView.vue'
import SslMonitorView from '../views/SslMonitorView.vue'
import JsonSerializerView from '../views/JsonSerializerView.vue'
import ImageCompressView from '../views/ImageCompressView.vue'
import MD5EncryptView from '../views/MD5EncryptView.vue'

// 导入AI智能助手页面
import AiIntelligenceView from '../views/AiIntelligenceView.vue'
import AiCodeAssistantView from '../views/AiCodeAssistantView.vue'
import AiDataAnalysisView from '../views/AiDataAnalysisView.vue'
import AiContentGeneratorView from '../views/AiContentGeneratorView.vue'
import AiModelManagementView from '../views/AiModelManagementView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: '/dashboard'
    },
    {
      path: '/jenkins',
      name: 'jenkins',
      component: () => import('../views/JenkinsView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView,
      meta: { requiresAuth: false }
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: DashboardView,
      meta: { requiresAuth: true }
    },
    {
      path: '/alerts',
      name: 'alerts',
      component: AlertsView,
      meta: { requiresAuth: true }
    },
    {
      path: '/media',
      name: 'media',
      component: MediaView,
      meta: { requiresAuth: true }
    },
    {
      path: '/templates',
      name: 'templates',
      component: TemplatesView,
      meta: { requiresAuth: true }
    },
    {
      path: '/silences',
      name: 'silences',
      component: SilencesView,
      meta: { requiresAuth: true }
    },
    {
      path: '/inhibitions',
      name: 'inhibitions',
      component: InhibitionsView,
      meta: { requiresAuth: true }
    },
    {
      path: '/prometheus',
      name: 'prometheus',
      component: PrometheusView,
      meta: { requiresAuth: true }
    },
    {
      path: '/config',
      name: 'config',
      component: ConfigView,
      meta: { requiresAuth: true }
    },
    {
      path: '/email-config',
      name: 'email-config',
      component: EmailConfigView,
      meta: { requiresAuth: true }
    },
    {
      path: '/logs/operation',
      name: 'operation-logs',
      component: OperationLogsView,
      meta: { requiresAuth: true }
    },
    {
      path: '/logs/login',
      name: 'login-logs',
      component: LoginLogsView,
      meta: { requiresAuth: true }
    },
    {
      path: '/profile',
      name: 'profile',
      component: ProfileView,
      meta: { requiresAuth: true }
    },
    {
      path: '/users',
      name: 'users',
      component: UsersView,
      meta: { requiresAuth: true }
    },
    {
      path: '/rules',
      name: 'rules',
      component: RulesView,
      meta: { requiresAuth: true }
    },
    {
      path: '/hosts',
      name: 'hosts',
      component: HostsView,
      meta: { requiresAuth: true }
    },
    {
      path: '/cluster-dashboard',
      name: 'cluster-dashboard',
      component: () => import('../views/ClusterDashboardView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/clusters',
      name: 'clusters',
      component: ClustersView,
      meta: { requiresAuth: true }
    },
    {
      path: '/nodes',
      name: 'nodes',
      component: () => import('../views/NodesView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/pods',
      name: 'pods',
      component: () => import('../views/PodsView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/deployments',
      name: 'deployments',
      component: () => import('../views/DeploymentsView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/services',
      name: 'services',
      component: () => import('../views/ServicesView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/configmaps',
      name: 'configmaps',
      component: () => import('../views/ConfigMapsView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/secrets',
      name: 'secrets',
      component: () => import('../views/SecretsView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/ingresses',
      name: 'ingresses',
      component: () => import('../views/IngressesView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/pvcs',
      name: 'pvcs',
      component: () => import('../views/PVCsView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/pvs',
      name: 'pvs',
      component: () => import('../views/PVsView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/storageclasses',
      name: 'storageclasses',
      component: () => import('../views/StorageClassesView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/events',
      name: 'events',
      component: () => import('../views/EventsView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/resourcequotas',
      name: 'resourcequotas',
      component: () => import('../views/ResourceQuotasView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/limitranges',
      name: 'limitranges',
      component: () => import('../views/LimitRangesView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/namespaces',
      name: 'namespaces',
      component: NamespacesView,
      meta: { requiresAuth: true }
    },
    {
      path: '/sls-logs',
      name: 'sls-logs',
      component: () => import('../views/SlsLogsView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/sls-config',
      name: 'sls-config',
      component: () => import('../views/SlsConfigView.vue'),
      meta: { requiresAuth: true }
    },
    { path: '/router-trace',
      name: 'router-trace',
      component: RouterTraceView,
      meta: { requiresAuth: false }
    },
    { path: '/ping-test',
      name: 'ping-test',
      component: PingTestView,
      meta: { requiresAuth: false }
    },
    { path: '/telnet-test',
      name: 'telnet-test',
      component: TelnetTestView,
      meta: { requiresAuth: false }
    },
    {
      path: '/ssl-monitor',
      name: 'ssl-monitor',
      component: SslMonitorView,
      meta: { requiresAuth: true }
    },
    {
      path: '/json-serializer',
      name: 'json-serializer',
      component: JsonSerializerView,
      meta: { requiresAuth: true }
    },
    {
      path: '/image-compress',
      name: 'image-compress',
      component: ImageCompressView,
      meta: { requiresAuth: true }
    },
    {
      path: '/md5-encrypt',
      name: 'md5-encrypt',
      component: MD5EncryptView,
      meta: { requiresAuth: true }
    },
    {
      path: '/kvm-hosts',
      name: 'kvm-hosts',
      component: KvmHostsView,
      meta: { requiresAuth: true }
    },
    {
      path: '/ai-intelligence',
      name: 'ai-intelligence',
      component: AiIntelligenceView,
      meta: { requiresAuth: true }
    },
    {
      path: '/ai-code-assistant',
      name: 'ai-code-assistant',
      component: AiCodeAssistantView,
      meta: { requiresAuth: true }
    },
    {
      path: '/ai-data-analysis',
      name: 'ai-data-analysis',
      component: AiDataAnalysisView,
      meta: { requiresAuth: true }
    },
    {
      path: '/ai-content-generator',
      name: 'ai-content-generator',
      component: AiContentGeneratorView,
      meta: { requiresAuth: true }
    },
    {
      path: '/ai-model-management',
      name: 'ai-model-management',
      component: AiModelManagementView,
      meta: { requiresAuth: true }
    }
  ]
})

// 路由守卫
router.beforeEach((to, _from, next) => {
  // 检查路由是否需要认证
  if (to.matched.some(record => record.meta.requiresAuth)) {
    // 检查是否已登录
    const token = localStorage.getItem('token')
    if (!token) {
      // 未登录，重定向到登录页面
      next({
        path: '/login',
        query: { redirect: to.fullPath }
      })
    } else {
      // 已登录，继续访问
      next()
    }
  } else {
    // 不需要认证的路由，直接访问
    next()
  }
})

export default router