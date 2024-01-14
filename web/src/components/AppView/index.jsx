import React from 'react'
import { Avatar, Tag } from 'antd'
import chatImage from '@/assets/images/chat.png'
import style from './index.module.scss'


const AppView = (props) => {
    const { data } = props

    const split = (text) => {
        const items = text.split(/\n/).filter((i) => i.length)
        return items.map((t, i) => <div key={i} className={i !== items.length - 1 ? "mb-4" : ''}>{t}</div>)
    }


    return <React.Fragment>
        <div className={style.container}>
            <div className={style.layout} style={{ backgroundImage: "url(" + chatImage + ")" }}>
                <div className={style.header}>
                    {props.title}
                </div>
                <div className={style.body}>
                    {(data || []).map((item, i) => (
                        item.info.tip ? (
                            <div className={style.tip} key={i}>
                                {item.info.tip}
                            </div>
                        ) : (
                            <div className={style.item} key={i}>
                                <div className={style.avatar}>
                                    <span><Avatar src={item.info.avatar} /></span>
                                </div>
                                <div>
                                    <div className={[style.feed, style.name].join(' ')}>
                                        <span style={{ paddingRight: '6px' }}>{item.info.name}</span>
                                        {item.info.tag ? <Tag color="orange" style={{ paddingRight: '6px' }}>{item.info.tag}</Tag> : ""}
                                        <span>{item.info.signature}</span>
                                    </div>
                                    {(item.message || []).map((t, i) =>
                                        <div key={i} className={[style.feed, style.c].join(' ')}>
                                            <div>{t.type === 'text' ? split(t.text) : t.text}</div>
                                            {
                                                t.type === 'tags' ? (<div className={style.tags}>
                                                    {(t.tags || []).map((item, i) => <Tag color="#108ee9" key={i}>{item.name}</Tag>)}
                                                </div>) :
                                                    t.type === 'recommend' ? (t.recommend.map((t, i) => <div className={style.link} key={i}>{t}</div>)) : ''
                                            }
                                        </div>
                                    )}
                                </div>
                            </div>
                        )
                    ))}
                </div>
            </div>
        </div>
    </React.Fragment>
}
export default AppView