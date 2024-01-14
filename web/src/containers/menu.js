import i18n from 'i18next'

const menu = [
  {
    title: i18n.t('menu.faq'),
    key: '/admin/:bot/knowledges',
    icon: 'icon-knowledge',
    auth: [1]
  },
  {
    title: i18n.t('menu.chat'),
    key: '/admin/:bot/chat',
    icon: 'icon-chat',
    auth: [1]
  },
  {
    title: i18n.t('menu.publish'),
    key: '/admin/:bot/publish',
    icon: 'icon-launch',
    auth: [1]
  },
  {
    title: i18n.t('menu.messages'),
    key: '/admin/:bot/messages',
    icon: 'icon-history',
    auth: [1]
  },
  {
    title: i18n.t('menu.setting'),
    key: '/admin/:bot/setting',
    icon: 'icon-setting',
    auth: [1]
  }

  /*
  {
    title: i18n.t('menu.data'),
    key: '/admin/:bot/data',
    icon: 'icon-Data',
    auth: [1]
  },
  */
]

export default menu
