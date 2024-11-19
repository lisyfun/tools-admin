import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw, RouteMeta } from 'vue-router'
import Layout from '@/layout/index.vue'

// 扩展 RouteMeta 类型
interface CustomRouteMeta extends RouteMeta {
  title?: string
  icon?: string
  affix?: boolean
}

// 扩展 RouteRecordRaw 类型
type CustomRouteRecordRaw = RouteRecordRaw & {
  hidden?: boolean
  children?: CustomRouteRecordRaw[]
  meta?: CustomRouteMeta
}

const routes: CustomRouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: { title: '登录' },
    hidden: true
  },
  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: { title: '数据面板', icon: 'DataLine' },
        children: [
          {
            path: 'sms-count',
            name: 'SmsCount',
            component: () => import('@/views/dashboard/sms-count.vue'),
            meta: { title: '短信数', icon: 'Message' }
          },
          {
            path: 'task-count',
            name: 'TaskCount',
            component: () => import('@/views/dashboard/task-count.vue'),
            meta: { title: '任务数', icon: 'List' }
          },
          {
            path: 'success-rate',
            name: 'SuccessRate',
            component: () => import('@/views/dashboard/success-rate.vue'),
            meta: { title: '执行成功率', icon: 'TrendCharts' }
          }
        ]
      }
    ]
  },
  {
    path: '/task',
    component: Layout,
    children: [
      {
        path: '',
        name: 'Task',
        component: () => import('@/views/task/index.vue'),
        meta: { title: '任务管理', icon: 'List' }
      },
      {
        path: 'generate',
        name: 'TaskGenerate',
        meta: { title: '任务生成', icon: 'Plus' },
        children: [
          {
            path: 'datax',
            name: 'DataxTask',
            component: () => import('@/views/task/generate/datax.vue'),
            meta: { title: 'DataX任务', icon: 'DataAnalysis' }
          },
          {
            path: 'shell',
            name: 'ShellTask',
            component: () => import('@/views/task/generate/shell.vue'),
            meta: { title: 'Shell任务', icon: 'Terminal' }
          },
          {
            path: 'http',
            name: 'HttpTask',
            component: () => import('@/views/task/generate/http.vue'),
            meta: { title: 'HTTP任务', icon: 'Link' }
          }
        ]
      },
      {
        path: 'log',
        name: 'TaskLog',
        component: () => import('@/views/task/log/index.vue'),
        meta: { title: '任务日志', icon: 'Document' }
      }
    ]
  },
  {
    path: '/sms',
    component: Layout,
    name: 'SMS',
    meta: {
      title: '短信管理',
      icon: 'Message'
    },
    children: [
      {
        path: 'index',
        name: 'SmsList',
        component: () => import('@/views/sms/index.vue'),
        meta: { title: '短信列表', icon: 'List' }
      },
      {
        path: 'template',
        name: 'SmsTemplate',
        component: () => import('@/views/sms/template/index.vue'),
        meta: { title: '短信模板', icon: 'Document' }
      },
      {
        path: 'recipient',
        name: 'SmsRecipient',
        component: () => import('@/views/sms/recipient/index.vue'),
        meta: { title: '收件人', icon: 'User' }
      }
    ]
  },
  {
    path: '/pipeline',
    component: Layout,
    name: 'Pipeline',
    meta: { title: '流水线管理', icon: 'Connection' },
    children: [
      {
        path: 'build',
        name: 'PipelineBuild',
        component: () => import('@/views/pipeline/build/index.vue'),
        meta: { title: '容器构建', icon: 'Box' }
      },
      {
        path: 'deploy',
        name: 'PipelineDeploy',
        component: () => import('@/views/pipeline/deploy/index.vue'),
        meta: { title: '容器部署', icon: 'Upload' }
      }
    ]
  },
  {
    path: '/server',
    component: Layout,
    name: 'Server',
    meta: { title: '服务器管理', icon: 'Monitor' },
    children: [
      {
        path: 'add',
        name: 'ServerAdd',
        component: () => import('@/views/server/add/index.vue'),
        meta: { title: '服务器新增', icon: 'Plus' }
      },
      {
        path: 'log',
        name: 'ServerLog',
        component: () => import('@/views/server/log/index.vue'),
        meta: { title: '服务器日志', icon: 'Document' }
      }
    ]
  },
  {
    path: '/database',
    component: Layout,
    name: 'Database',
    meta: { title: '数据库管理', icon: 'Grid' },
    children: [
      {
        path: 'connection',
        name: 'DatabaseConnection',
        component: () => import('@/views/database/connection/index.vue'),
        meta: { title: '数据库连接', icon: 'Link' }
      },
      {
        path: 'operation',
        name: 'DatabaseOperation',
        component: () => import('@/views/database/operation/index.vue'),
        meta: { title: '数据库操作', icon: 'Operation' }
      }
    ]
  },
  {
    path: '/system',
    component: Layout,
    name: 'System',
    meta: { title: '系统管理', icon: 'Setting' },
    children: [
      {
        path: 'permission',
        name: 'Permission',
        meta: { title: '权限管理', icon: 'Lock' },
        children: [
          {
            path: 'user',
            name: 'User',
            component: () => import('@/views/system/permission/user.vue'),
            meta: { title: '用户管理', icon: 'User' }
          },
          {
            path: 'menu',
            name: 'Menu',
            component: () => import('@/views/system/permission/menu.vue'),
            meta: { title: '菜单管理', icon: 'Menu' }
          },
          {
            path: 'role',
            name: 'Role',
            component: () => import('@/views/system/permission/role.vue'),
            meta: { title: '角色管理', icon: 'UserFilled' }
          }
        ]
      },
      {
        path: 'log',
        name: 'Log',
        meta: { title: '日志管理', icon: 'Document' },
        children: [
          {
            path: 'operation',
            name: 'OperationLog',
            component: () => import('@/views/system/log/operation.vue'),
            meta: { title: '操作日志', icon: 'List' }
          },
          {
            path: 'system',
            name: 'SystemLog',
            component: () => import('@/views/system/log/system.vue'),
            meta: { title: '系统日志', icon: 'Monitor' }
          }
        ]
      }
    ]
  }
]

// 动态路由，基于用户权限动态去加载
export const dynamicRoutes: CustomRouteRecordRaw[] = []

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
