import loadable from '@/utils/loadable'

const Knowledges = loadable(() => import('@/pages/Knowledges'))
const Setting = loadable(() => import('@/pages/Setting'))
const Chat = loadable(() => import('@/pages/Chat'))
const Data = loadable(() => import('@/pages/Data'))
const Publish = loadable(() => import('@/pages/Publish'))
const Messages = loadable(() => import('@/pages/Messages'))

const router = [
  { path: '/admin/:bot/knowledges', exact: false, name: 'Knowledge', component: Knowledges, auth: [1] },
  { path: '/admin/:bot/setting', exact: false, name: 'setting', component: Setting, auth: [1] },
  { path: '/admin/:bot/chat', exact: false, name: 'chat', component: Chat, auth: [1] },
  { path: '/admin/:bot/data', exact: false, name: 'data', component: Data, auth: [1] },
  { path: '/admin/:bot/publish', exact: false, name: 'publish', component: Publish, auth: [1] },
  { path: '/admin/:bot/messages', exact: false, name: 'massages', component: Messages, auth: [1] }
]

export default router
