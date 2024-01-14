import React, { useState, useEffect } from 'react'
import { message, Form, Input, Button, Breadcrumb } from 'antd'
import { connect } from 'react-redux'
import { useTranslation } from 'react-i18next'
import CommonModal from '@/commonComponents/commonModal'

import Tree from '@/components/Tree'
import actions from '@/redux/actions'
import Api from '@/apiService/Api'
import Tiptap from './tiptap'
import TagsInput from 'react-tagsinput'
import '@/style/react-tagsinput.css'

const msg = message
const layout = {
  labelCol: { span: 4 },
  wrapperCol: { span: 24 }
}

const Edit = (props) => {
  const { folders, knowledgeChange, knowledge, folder, bot, close, breadcrumb } = props
  const { t } = useTranslation()
  const [form] = Form.useForm()

  const [names, setNames] = useState([])
  const [tags, setTags] = useState([])
  const [visible, setVisible] = useState(false)
  const [path, setPath] = useState('')
  const [html, setHtml] = useState('')

  const formValue = {
    id: knowledge.id || '',
    folder_id: knowledge.folder_id || folder.id,
    bot_id: bot.id
  }

  if (folders.length && !folders.length) {
    //setExpandedKeys([folders[0].key])
  }

  useEffect(() => {
    const item = (knowledge.path || folder.id || '').split(',')
    const names = breadcrumb(item[item.length - 1])
    setNames(names)
    const path = names.map((n) => n.id).join(',')
    setPath(path)
    setHtml(knowledge.content)
    setTags(knowledge.tags || [])
    form.setFieldsValue({
      path
    })
  }, [])

  let key = null
  const selectFolder = (keys) => {
    if (keys.length) {
      key = keys[0]
    } else {
      key = null
    }
  }

  const setFolder = () => {
    if (key) {
      const names = breadcrumb(key)
      setNames(names)
      const path = '' + names.map((n) => n.id).join(',')
      setPath(path)
      form.setFieldsValue({
        path
      })
      setVisible(false)
      key = null
    } else {
      const text = t('knowledge.modify_folder_alert')
      message.warning({
        content: text,
        key: text
      })
    }
  }

  const handleTagChange = (tags) => {
    setTags(tags)
  }

  const confirm = () => {
    form
      .validateFields()
      .then(async (values) => {
        values.tags = tags
        values.content = html
        if (values.id) {
          const res = await Api.put_knowledge(values.id, values)
          const { success, data } = res
          if (success) {
            msg.success(t('knowledge.success'))
            knowledgeChange(data)
          } else {
            msg.error(t('process_fail'))
          }
        } else {
          const res = await Api.add_knowledge(values)
          const { success,  data } = res
          if (success) {
            msg.success(t('knowledge.success'))
            knowledgeChange(data, '+')
          } else {
            msg.error(t('process_fail'))
          }
        }
        close()
      })
      .catch((info) => {
        console.error('info=====>>>', info)
      })
  }

  const rules = [{ required: true, message: t('required') }]

  return (
    <div className='mainContentInner'>
      <CommonModal
        title={t('knowledge.modify_folder')}
        width={600}
        visible={visible}
        common_cancel={() => {
          setVisible(false)
        }}
        common_confirm={setFolder}
        children={
          <div>
            <Tree onSelect={selectFolder} data={folders} />
          </div>
        }
      />
      <div className='py-3'>
        <span className='hand back-link' onClick={() => close()}>
          <span style={{ paddingRight: '8px' }}> &lt;</span>
          {t('back')}
        </span>
      </div>
      <Form
        {...layout}
        form={form}
        onFinish={confirm}
        initialValues={{
          ...formValue,
          name: knowledge.name,
          path
        }}
      >
        <Form.Item name='id' noStyle>
          <Input type='hidden' />
        </Form.Item>
        <Form.Item name='bot_id' noStyle>
          <Input type='hidden' />
        </Form.Item>
        <Form.Item name='folder_id' noStyle>
          <Input type='hidden' />
        </Form.Item>
        <Form.Item name='path' noStyle>
          <Input type='hidden' style={{ width: '100%' }} />
        </Form.Item>
        <Form.Item labelAlign='left'>
          <Breadcrumb separator='>' style={{ marginRight: 16, display: 'inline-block' }}>
            {names.map((n) => (
              <Breadcrumb.Item key={n.name}>{n.name}</Breadcrumb.Item>
            ))}
          </Breadcrumb>
          <Button className='btn' type='link' onClick={() => setVisible(true)}>
            {t('knowledge.modify_folder')}
          </Button>
        </Form.Item>
        <Form.Item labelAlign='left' rules={rules} name='name'>
          <Input placeholder={t('knowledge.knowledge_p')} />
        </Form.Item>
        <Form.Item labelAlign='left'>
          <div style={{ border: '1px solid #d9d9d9', borderRadius: '2px', padding: '0 6px' }}>
            <Tiptap
              html={knowledge.content}
              onChange={(html) => {
                setHtml(html)
              }}
            />
          </div>
        </Form.Item>
        <Form.Item labelAlign='left'>
          <TagsInput
            style={{ width: '100%', borderRadius: '2px' }}
            value={tags}
            onChange={handleTagChange}
            addOnBlur={true}
            inputProps={{ placeholder: !tags.length ? t('knowledge.relate_p') : '' }}
          />
        </Form.Item>
        <Form.Item className='py-4'>
          <Button size="large" className="mr-4" onClick={() => close()}>
            {t('cancel')}
          </Button>
          <Button type='primary' size="large" htmlType='submit'>
            {t('ok')}
          </Button>
        </Form.Item>
      </Form>
    </div>
  )
}

const mapStateToProps = (state) => {
  return {
    bot: state.bot,
    app: state.app
  }
}
const mapDispatchToProps = {
  ...actions
}

export default connect(mapStateToProps, mapDispatchToProps)(React.memo(Edit))
