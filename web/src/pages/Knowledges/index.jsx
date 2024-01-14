import React, { useState, useEffect, useRef } from 'react'
import { MoreOutlined, UploadOutlined, DeleteOutlined, SolutionOutlined, CheckOutlined, CloseOutlined, LockOutlined, ExceptionOutlined, QuestionCircleOutlined } from '@ant-design/icons'
import { Avatar, Breadcrumb, Tooltip, message, Table, Button, Row, Col, Dropdown, Menu, Upload, Alert, Space, Tour } from 'antd'
import { connect } from 'react-redux'
import { useTranslation } from 'react-i18next'
import moment from 'moment'
import _ from 'lodash'
import Tree from '@/components/Tree'
import arrayToTree from 'array-to-tree'
import actions from '@/redux/actions'
import Api from '@/apiService/Api'
import Classification from './Classification'
import CommonModal from '@/commonComponents/commonModal'
import Edit from './edit'

const { Dragger } = Upload
const msg = message

const Knowledge = (props) => {
  const { bot } = props
  const { t } = useTranslation()
  const defaultQuery = { p: 0 }
  const defaultKey = '0'
  const defaultFolder = [
    {
      title: t('knowledge.folder.all'),
      key: '0',
      children: []
    }
  ]

  const [open, setOpen] = useState(false)
  const [folder, setFolder] = useState({})
  const [folders, setFolders] = useState(defaultFolder)
  const [knowledge, setKnowledge] = useState({})
  const [knowledges, setKnowledges] = useState([])
  const [delVisible, setDelVisible] = useState(false)
  const [selectedRowKeys, setSelectedRowKeys] = useState([])
  const [modifyFolderVisible, setModifyFolderVisible] = useState(false)
  const [visible, setVisible] = useState(false)
  const [editVisible, setEditVisible] = useState(false)
  const [fileList, setFileList] = useState([])
  const [isMounted, setIsMounted] = useState(true)
  const [pagination, setPagination] = useState({
    page: 1,
    size: 20,
    total: 0
  })

  const ref1 = useRef(null);
  const ref2 = useRef(null);
  const steps = [
    {
      title: t('knowledge.tour.ref1_title'),
      description: t('knowledge.tour.ref1_desc'),
      target: () => ref1.current,
      placement:"rightTop",
    },
    {
      title: t('knowledge.upload'),
      description: t('knowledge.support_file'),
      target: () => ref2.current,
      placement:"top",
    },
  ]

  const closeTour = () => {
    localStorage.setItem('knowledges.tour', 1)
    setOpen(false)
  }

  useEffect(() => {
    if (!localStorage.getItem('knowledges.tour') && ref1.current){
      setOpen(true)
    }
  }, [])

  const upload = {
    name: 'files',
    multiple: true,
    showUploadList: true,
    fileList: fileList,
    accept: '.doc, .docx, .txt, text/markdown, application/msword, application/vnd.openxmlformats-officedocument.wordprocessingml.document, .pdf',
    data: (file) => {
      const items = breadcrumb(folder.key)
      const path = items.map((item) => item.id).join(',')
      return {
        bot_id: bot.id,
        folder_id: folder.id,
        path: path
      }
    },
    action: Api.UPLOAD_URL,
    onChange(info) {
      const { status, response, name } = info.file
      if (status !== 'uploading') {
        //console.log(info.file, info.fileList)
      }
      if (status === 'done') {
        const { success, data } = response
        if (!success) {
          return message.error(t('knowledge.error', { file: name }), 6)
        }
        message.success(t('knowledge.success', { file: name }), 6)
        for (let i = 0; i < data.length; i++) {
          knowledgeChange(data[i], '+')
        }
      } else if (status === 'error') {
        message.error(t('knowledge.error', { file: name }), 6)
      }
      setFileList(info.fileList)
    },
    onDrop(e) {
      console.error('Dropped files', e.dataTransfer.files)
    }
  }

  const delKnowledge = async () => {
    if (selectedRowKeys.length) {
      return deleteKnowledges()
    }
    const { id } = knowledge
    if (!id) {
      return
    }
    const res = await Api.del_knowledge(id)
    const { success } = res
    if (success) {
      msg.success(t('knowledge.success'))
      let items = [...knowledges]
      items = _.filter(items, (item) => {
        return item.id !== id
      })
      setKnowledges(items)
      selectChange([id])
    } else {
      msg.error(t('process_fail'))
    }
    setDelVisible(false)
  }

  const deleteKnowledges = async () => {
    const items = _.cloneDeep(selectedRowKeys)
    const res = await Api.batch_del_knowledge(items)
    const { success, data } = res
    if (success) {
      msg.success(t('knowledge.success'))
      let items = [...knowledges]
      items = _.filter(items, (item) => {
        return !data.includes(item.id)
      })
      setKnowledges(items)
      setSelectedRowKeys([])
    } else {
      msg.error(t('process_fail'))
    }
    setDelVisible(false)
  }

  const edit = (node, action) => {
    if (action === 'delete') {
      setKnowledge(node)
      return setDelVisible(true)
    }
    setKnowledge(action === 'add' ? {} : node)
    setEditVisible(true)
  }

  const selectChange = (list) => {
    const cpItems = _.cloneDeep(selectedRowKeys)
    _.remove(cpItems, (v, i) => list.includes(v))
    setSelectedRowKeys(cpItems)
  }

  // breadcrumb for category name
  const breadcrumb = (key) => {
    const info = folders[0]
    const treePath = (tree, func, path = []) => {
      if (!tree) return []
      for (const data of tree) {
        path.push(data)
        if (func(data)) return path
        if (data.children) {
          const children = treePath(data.children, func, path)
          if (children.length) return children
        }
        path.pop()
      }
      return []
    }
    const items = treePath(info.children, (node) => node.key === key)
    return items
  }

  let key = null
  const selectFolder = (keys) => {
    if (keys.length) {
      key = keys[0]
    } else {
      key = null
    }
  }

  const refresh = () => {
    if (isMounted) {
      const items = breadcrumb(folder.key)
      const key = items.map((item) => item.id).join(',')
      get(key ? key : defaultKey)
    }
  }

  const knowledgeChange = (c, action) => {
    const items = _.cloneDeep(knowledges)
    // remove
    if (action === '-') {
      _.remove(items, (v, i) => v.id === c.id)
    }
    if (action === '+') {
      items.push(c)
    }
    // update
    if (!action) {
      const i = items.findIndex((i) => i.id === c.id)
      items.splice(i, 1, c)
    }
    setKnowledges(items)
    return items
  }

  // recive folder change
  const folderChange = (c, action, add) => {
    setSelectedRowKeys([])
    // select
    if (action === 'select') {
      setFolder(c)
      const items = breadcrumb(c.key)
      const key = items.map((item) => item.id).join(',')
      return get(key ? key : defaultKey)
    }
    const loop = (items) => {
      if (items.key === c.key && action === '+') {
        if (!items.children) {
          items.children = [add]
        } else {
          items.children.push(add)
        }
        return
      }
      items.children.forEach((item, index) => {
        if (item.key === c.key) {
          // add
          if (action === '+') {
            if (!item.children) {
              item.children = [add]
            } else {
              item.children.push(add)
            }
          }
          // remove
          if (action === '-') {
            _.remove(items.children, (v, i) => i === index)
          }
          // update
          if (!action) {
            item = c
            setFolder(c)
          }
        } else {
          item.children && loop(item)
        }
      })
      if (!items.children.length) {
        delete items.children
      }
    }
    let items = [...folders]
    loop(items[0])
    setFolders(items)
  }

  const closeEdit = () => {
    setEditVisible(false)
  }

  const onShowSizeChange = (current, pageSize) => {
    setPagination({
      page: current,
      size: pageSize
    })
    get(defaultKey)
  }

  const get = (path, search) => {
    const q = { ...defaultQuery }
    if (path && path === defaultKey) {
      q.path = '0'
    } else {
      q.path = path
    }
    Api.get_knowledge(bot.id, q).then((res) => {
      if (res.data) {
        const data = _.map(res.data, (item) => {
          item.key = item.id
          return item
        })
        setKnowledges([...data])
      }
    })
  }

  const treeFolders = (folders) => {
    _.map(folders, (item) => {
      item.key = item.id
      item.title = item.name
      return item
    })
    const children = arrayToTree(folders, {
      parentProperty: 'parent'
    })
    const info = defaultFolder[0]
    setFolders([
      {
        key: info.key,
        title: info.title,
        children
      }
    ])
  }

  useEffect(() => {
    Api.get_folders(bot.id).then((res) => {
      const folders = res.data
      treeFolders(folders)
      get(defaultKey)
    })
  }, [folders.prop])

  useEffect(() => {
    return () => {
      setIsMounted(false)
    }
  }, [])

  useEffect(() => {
    const item = knowledges.find((item) => item.status === 2)
    if (item) {
      setTimeout(() => {
        refresh()
      }, 10 * 1000)
    }
  }, [knowledges])

  const status = t('knowledge.table.status_list', { returnObjects: true })

  const columns = [
    {
      title: t('knowledge.table.name'),
      dataIndex: 'name',
      key: 'name',
      align: 'left',
      render: (r, record) => {
        return r ? (
          <span
            className='faq-link hand'
            onClick={(e) => {
              edit(record)
            }}
          >
            {' '}
            {r}{' '}
          </span>
        ) : (
          '-'
        )
      }
    },
    {
      title: t('knowledge.table.folder_name'),
      dataIndex: 'folder_name',
      key: 'folder_name',
      align: 'left',
      render: (r, node) => {
        let item = node.path.split(',')
        item = breadcrumb(item[item.length - 1])
        item = item.map((n) => n.name)
        if (item.length === 0) {
          return '-'
        }
        return item.join(' / ')
      }
    },
    {
      title: t('knowledge.table.relate'),
      dataIndex: 'tags',
      key: 'tags',
      align: 'left',
      render: (r) => {
        return r && r.length ? r.join(', ') : '-'
      }
    },
    {
      title: (
        <>
          {t('knowledge.table.status')}{' '}
        </>
      ),
      dataIndex: 'status',
      key: 'status',
      align: 'left',
      render: (r) => {
        const ColorList = ['#ff4d4f', '#52c41a', '#1677ff', '#f9c000', '#f9c000']
        const IconList = [<CloseOutlined />, <CheckOutlined />, <SolutionOutlined />, <LockOutlined />, <ExceptionOutlined />]
        return (
          <>
            <Avatar
              size='small'
              style={{
                backgroundColor: ColorList[r],
                width: '16px',
                height: '16px',
                lineHeight: '16px',
                fontSize: '10px'
              }}
              icon={IconList[r]}
            />{' '}
            {status[r]}
          </>
        )
      }
    },
    {
      title: t('knowledge.table.update'),
      dataIndex: 'created_at',
      key: 'created_at',
      align: 'left',
      render: (r) => {
        return r ? moment(r).format('MMM DD h:m') : '-'
      }
    },
    {
      title: '',
      align: 'left',
      key: 'op',
      width: 80,
      defaultSortOrder: 'descend',
      render: (text, record, idx) => {
        return (
          <Dropdown overlay={() => dropdown(record)} placement='bottomLeft'>
            <MoreOutlined style={{ cursor: 'pointer', fontSize: '180%', fontWeight: 'bolder' }} />
          </Dropdown>
        )
      }
    }
  ]

  const dropdown = (node) => (
    <Menu>
      <Menu.Item onClick={() => edit(node, 'edit')}>{t('dropdown.edit')}</Menu.Item>
      <Menu.Item onClick={() => edit(node, 'delete')}>{t('dropdown.delete')}</Menu.Item>
    </Menu>
  )

  return (
    <div className='settings-container'>
      {
        <div className='d-flex' style={{ height: '100%' }}>
          <div style={{ flexBasis: '250px', borderRight: '1px solid #ddd' }} >
            <Classification folderChange={folderChange} breadcrumb={breadcrumb} folders={folders} />
          </div>
          <div style={{ padding: '0', flexBasis: 'calc(100% - 250px)' }}>
            {selectedRowKeys && selectedRowKeys.length > 0 && (
              <Alert
                message={<strong>{t('knowledge.selected', { num: selectedRowKeys.length })}</strong>}
                type='warning'
                action={
                  <Space>
                    <Button
                      icon={<DeleteOutlined />}
                      size='small'
                      type='link'
                      onClick={() => {
                        setDelVisible(true)
                      }}
                    >
                      {t('knowledge.table.delete_title')}
                    </Button>
                  </Space>
                }
              />
            )}
            {editVisible ? (
              <Edit close={closeEdit} knowledgeChange={knowledgeChange} knowledge={knowledge} folders={folders} breadcrumb={breadcrumb} folder={folder} />
            ) : (
              <>
                <Row className='row'>
                  <div ref={ref1} ></div>
                  <Col flex='1 0 25%' className='column'>
                    <Space align='center' style={{ margin: 16 }}>
                      <Breadcrumb separator='>' style={{ margin: '0 16px 0 0', display: 'inline-block' }}>
                        {breadcrumb(folder.key).map((n) => (
                          <Breadcrumb.Item key={n.name}>{n.name}</Breadcrumb.Item>
                        ))}
                      </Breadcrumb>
                    </Space>
                  </Col>
                  <Col flex='' className='column'>
                    <Space align='center' style={{ margin: 16 }}>
                      <Button icon={<DeleteOutlined />} size='small' type='link' onClick={() => edit(null, 'add')}>
                        {t('knowledge.add')}
                      </Button>
                      <div  ref={ref2} >
                        <Button icon={<DeleteOutlined />} size='small'type='link' onClick={() => setVisible(true)}>
                          {t('knowledge.upload')}
                        </Button>
                      </div>
                      
                    </Space>
                  </Col>
                </Row>
                <Table
                  dataSource={knowledges}
                  columns={columns}
                  rowSelection={{
                    type: 'checkbox',
                    onChange: (selectedRowKeys) => {
                      setSelectedRowKeys(selectedRowKeys)
                    }
                  }}
                  pagination={{
                    defaultCurrent: 1,
                    defaultPageSize: 20,
                    current: pagination.page,
                    pageSize: pagination.size,
                    total: pagination.total,
                    hideOnSinglePage: true,
                    showTotal: (total) => t('knowledge.table.total', { total }),
                    onShowSizeChange: { onShowSizeChange }
                  }}
                />
              </>
            )}

            <CommonModal
              title={t('knowledge.modify_folder')}
              width={600}
              visible={modifyFolderVisible}
              common_cancel={() => {
                setModifyFolderVisible(false)
              }}
              common_confirm={setFolders}
              children={
                <div>
                  <Tree onSelect={selectFolder} data={folders} />
                </div>
              }
            />
          </div>
        </div>
      }

      <CommonModal
        title={t('knowledge.table.delete_title')}
        width={600}
        danger={true}
        visible={delVisible}
        common_cancel={() => {
          setDelVisible(false)
        }}
        common_confirm={delKnowledge}
        children={<div style={{ textAlign: 'center', paddingTop: 10, paddingBottom: 10 }}>{t('knowledge.table.delete_desc')}</div>}
      />

      <CommonModal
        visible={visible}
        title={<h3 className='text-center mt-2'>{t('knowledge.upload_title')}</h3>}
        common_cancel={() => {
          setVisible(false)
          setFileList([])
        }}
        width={600}
        footer={false}
      >
        <Dragger {...upload}>
          <div className='py-5'>
            <p className='ant-upload-drag-icon'>
              <UploadOutlined />
            </p>
            <p className='ant-upload-text'>{t('knowledge.support_file')}</p>
            <p className='ant-upload-hint'>{t('knowledge.upload_hint')}</p>
          </div>
        </Dragger>
      </CommonModal>
      <Tour open={open} 
        mask={false} 
        type="primary" 
        onClose={closeTour} 
        steps={steps}
        indicatorsRender={(current, total) => (
          <span>
            {current + 1} / {total}
          </span>
        )}
        />
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
export default connect(mapStateToProps, mapDispatchToProps)(React.memo(Knowledge))
