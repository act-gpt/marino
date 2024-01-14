import React, { useState, useEffect } from 'react'
import { Table, Row, Col, Tooltip, message } from 'antd'
import { useTranslation } from 'react-i18next'
import Markdown from 'react-markdown'
import { RightOutlined, DownOutlined, QuestionCircleOutlined } from '@ant-design/icons'
import { connect } from 'react-redux'
import Api from '@/apiService/Api'
import actions from '@/redux/actions'

import moment from 'moment'


const Messages = (props) => {
  const { bot } = props
  const { t } = useTranslation()
  const [messageApi, contextHolder] = message.useMessage()
  const [data, setData] = useState([])
  const [loading, setLoading] = useState(true)
  const [params, setParams] = useState({
    page: 1,
    total: 0,
    size: 20
  })

  const getData = async (page, size) => {
    if (page < 0) {
      return
    }
    const res = await Api.messages(bot.id, { page, size })
    const { success, meta, data } = res
    setLoading(false)
    if (!success) {
      return messageApi.open({
        type: 'error',
        content: t('process_fail')
      })
    }
    setData(data)
    setParams(meta)
  }

  useEffect(() => {
    if (!loading) {
      return
    }
    getData(params.page, params.size)
  }, [loading])

  const columns = [
    {
      title: t('message.columns.question'),
      dataIndex: 'question',
      key: 'source',
      ellipsis: true
    },
    {
      title: t('message.columns.answer'),
      dataIndex: 'answer',
      key: 'answer',
      ellipsis: true
    },
    {
      title: (
        <>
          {t('message.columns.feedback')}{' '}
          <Tooltip placement='right' title={t('message.columns.feedback_desc')}>
            <QuestionCircleOutlined />
          </Tooltip>
        </>
      ),
      dataIndex: 'feedback',
      key: 'feedback',
      align: 'center',
      width: 120,
      render: (r, item) => {
        return item.dislike + ' / ' + item.like
      }
    },
    {
        title: t("message.columns.lft"),
        dataIndex: 'llm_first_time',
        key: 'llm_first_time',
        align: 'left',
        width: 120,
        render: (r, item) => {
          return r.toFixed(2) + ' s'
        }
    },
    {
      title: t('message.columns.llm'),
      dataIndex: 'llm_time',
      key: 'llm_time',
      align: 'center',
      width: 120,
      render: (r, item) => {
        return r.toFixed(2) + ' s'
      }
    },
    {
      title: t('message.columns.time'),
      dataIndex: 'created_at',
      key: 'created_at',
      align: 'center',
      width: 160,
      render: (r) => {
        return r ? moment(r).format('YYYY-MM-DD HH:mm') : '-'
      }
    }
  ]
  return (
    <div className='settings-container px-5 pt-4'>
      {contextHolder}
      <div>
      {
        <Table
          className='textColor'
          rowKey='id'
          pagination={{
            current: params.page,
            defaultPageSize: params.size,
            total: params.total,
            hideOnSinglePage: true,
            onChange: (page, size) => {
              getData(page, size)
            }
          }}
          expandable={{
            expandedRowRender: (record) => (
              <div className='ml-2'>
                <Row wrap={false} className='py-2'>
                  <Col flex='80px'>
                    <h6>{t('message.columns.question')}</h6>
                  </Col>
                  <Col flex='auto'>{record.question}</Col>
                </Row>
                <Row wrap={false} className='py-2'>
                  <Col flex='80px'>
                    <h6>{t('message.columns.answer')}</h6>
                  </Col>
                  <Col flex='auto'><Markdown>{record.answer}</Markdown></Col>
                </Row>
              </div>
            ),
            expandIcon: ({ expanded, onExpand, record }) => (expanded ? <DownOutlined style={{color: "#bfbfbf"}} onClick={(e) => onExpand(record, e)} /> : <RightOutlined style={{color: "#bfbfbf"}}  onClick={(e) => onExpand(record, e)} />)
          }}
          columns={columns}
          dataSource={data}
        />
      }
      </div>
    </div>
  )
}

const mapStateToProps = (state) => {
  return {
    bot: state.bot,
    user: state.user,
    app: state.app
  }
}
const mapDispatchToProps = {
  ...actions
}

export default connect(mapStateToProps, mapDispatchToProps)(Messages)
