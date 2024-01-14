/* eslint-disable react-hooks/exhaustive-deps */
import React, { useState, useEffect } from 'react'
import { message, Breadcrumb, Menu, Dropdown, Form, Input, Button } from 'antd'
import { PlusOutlined, MoreOutlined, CaretDownOutlined } from '@ant-design/icons'
import CommonModal from '@/commonComponents/commonModal'
import '@/style/classification.scss'
import { connect } from "react-redux"
import actions from '@/redux/actions'
import { useTranslation } from 'react-i18next'
import Api from '@/apiService/Api'
import Tree from '@/components/Tree'
import { keys } from 'lodash'

const msg = message
const layout = {
    labelCol: { span: 8 },
    wrapperCol: { span: 16 },
}

const Classification = (props) => {

    const { folderChange, folders, defaultKey, breadcrumb, bot, app } = props
    const { t, i18n } = useTranslation()
    const [folder, setFolder] = useState({})
    const [visible, setVisible] = useState(false)
    const [delVisible, setDelVisible] = useState(false)
    const [selectKey, setSelectKey] = useState(defaultKey)
    const [name, setName] = useState([])

    const [formValue, setFormValue] = useState({
        id: "",
        bot_id: bot.id,
        parent: defaultKey || "0"
    })

    const [form] = Form.useForm()

    const edit = (node, action, e) => {
        //e && e.preventDefault()
        setFolder(node)
        if (action === "delete") {
            return setDelVisible(true)
        }
        const key = node.key 
        const val = { ...formValue, action }
        val.name = null
        if (action === "add") {
            val.parent = key == "0" ? '' : key
        } else {
            val.id = key == "0" ? '' : key
            val.parent = node.parent || ''
            val.name = node.name
        }
        setFormValue(val)
        form.setFieldsValue(val)
        const items = breadcrumb(node.key)
        setName(items)
        setVisible(true)
    }

    const confirm = () => {
        form
            .validateFields()
            .then((values) => {
                const add = formValue.action == "add"
                const act = add ? Api.add_folder(values) :
                    formValue.action == "edit" ? Api.put_folder(formValue.id, values) : null
                act && act.then((res) => {
                    const {success, message} = res
                    if (success) {
                        msg.success(formValue.action == "add" ? "success" : "success")
                        const obj = folder_info(selectKey)
                        const data = res.data
                        if (add) {
                            data.title = data.name
                            data.key = data.id
                        } else {
                            obj.name = obj.title = values.name
                        }
                        const item = folderChange(obj, add ? "+" : undefined, add ? data : null)
                    }else{
                        msg.error(message)
                    }
                })
                setVisible(false)
            })
            .catch((info) => {
                console.error('info=====>>>', info)
            })
    }


    const folder_info = (key) => {
        let val = {}
        const loop = (items) => {
            if (items.key == key) {
                return val = items
            }
            items.children.forEach((item) => {
                if (item.key === key) {
                    val = item
                } else {
                    item.children && loop(item)
                }
            })
        }
        loop(folders[0])
        return val
    }

    const delete_folder = async () => {
        const res = await Api.del_folder(folder.id)
        const {success, message} = res
        if(success){
            const obj = folder_info(selectKey)
            folderChange(obj, "-")
            msg.success('Deleteed')
        }else{
            msg.error(message)
        }
        setDelVisible(false)
    }

    const onSelect = (keys, e) => {
        e.nativeEvent.preventDefault()
        if (keys.length) {
            const key = keys[0]
            setSelectKey(key)
            const obj = folder_info(key)
            Object.keys(obj).length && folderChange(obj, "select")
        }
    }

    const dropdown = (node) => (
        <Menu>
            <Menu.Item onClick={(e) => { edit(node, 'edit', e) }}>{t("knowledge.folder.edit")}</Menu.Item>
            <Menu.Item onClick={(e) => { edit(node, 'add', e) }}>{t("knowledge.folder.add")}</Menu.Item>
            <Menu.Item onClick={(e) => { edit(node, 'delete', e) }}>{t("knowledge.folder.delete")}</Menu.Item>
        </Menu>
    )
    
    const all_command = (node) => {
        return (
            <div className="btn-wrapper" onClick={(e) => { edit(node, 'add'); }}>
                <Button icon={<PlusOutlined />} />
            </div>
        )
    }

    const more_command = (node) => {
        return (
            <Dropdown overlay={() => dropdown(node)} placement='bottomLeft' trigger={['click']}>
                <MoreOutlined style={{ cursor: 'pointer', fontSize: "180%", fontWeight: "bolder" }} />
            </Dropdown>
        )
    }


    const rules = [{ required: true, message: t('required') }]

    return (
        <div style={{ width: '100%', paddingTop: "12px" }}>
            <Tree
                virtual
                onSelect={onSelect}
                titleRender={(node) => {
                    return (<span style={{ width: "100%" }} onClick={(e) => {
                        if (node.key === selectKey) {
                            e.preventDefault()
                            e.stopPropagation()
                        }
                    }}>
                    {node.key != defaultKey ?
                        !node.children ? <span className="tree-node-noop"><span className="leaf-icon"></span></span> : "" : ""}
                        <div className="node">
                            <span className="node-name">{node.title}</span>
                            {node.key === "0" ? all_command(node) : more_command(node)}
                        </div></span>
                    )
                }}
                data={folders} />
            <CommonModal
                title={formValue.action === "add" ? t("knowledge.dialog.add_title") : t("knowledge.dialog.edit_title")}
                width={600}
                visible={visible}
                common_cancel={() => {
                    setVisible(false)
                }}
                common_confirm={confirm}
                children={
                    <Form
                        {...layout}
                        form={form}
                        initialValues={formValue}>
                        <Form.Item name="id" noStyle>
                            <Input type="hidden" />
                        </Form.Item>
                        <Form.Item name="bot_id" noStyle>
                            <Input type="hidden" />
                        </Form.Item>
                        <Form.Item name="action" noStyle>
                            <Input type="hidden" />
                        </Form.Item>
                        <Form.Item name="parent" noStyle>
                            <Input type="hidden" />
                        </Form.Item>
                        <Form.Item label={t("knowledge.dialog.clable")}>
                            <Breadcrumb separator=">" style={{ marginLeft: 16 }}>
                                {name.map((n) => <Breadcrumb.Item key={n.name}>{n.name}</Breadcrumb.Item>)}
                            </Breadcrumb>
                        </Form.Item>
                        <Form.Item rules={rules} label={t("knowledge.dialog.cname")} name='name' >
                            <Input autoComplete="false" placeholder={t("knowledge.dialog.cname_palceholder")} />
                        </Form.Item>
                    </Form>
                }
            />
            <CommonModal
                title={t("knowledge.dialog.delete_title")}
                danger={true}
                width={600}
                visible={delVisible}
                common_cancel={() => {
                    setDelVisible(false)
                }}
                common_confirm={delete_folder}
                children={<div style={{ textAlign: 'center', paddingTop: 10, paddingBottom: 10 }}>{t('knowledge.dialog.delete_desc')}</div>}
            />
        </div>
    )
}

const mapStateToProps = (state) => {
    return {
        bot: state.bot,
        app: state.app
    };
}
const mapDispatchToProps = {
    ...actions
};
export default connect(
    mapStateToProps,
    mapDispatchToProps
)(React.memo(Classification))
