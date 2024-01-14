import React, { useState }  from 'react'
import {Tree} from 'antd'
import { CaretDownOutlined } from '@ant-design/icons'

const { TreeNode } = Tree

const TreeComponent = (props) => {
    const [expandedKeys, setExpandedKeys] = useState([])
    const empty = () => {}
    const render = (node) => {
      return (
          <span style={{ width: '100%' }}>
          {
            <div className='node'>
              <span className='node-name'>{node.title}</span>
            </div>
          }
          </span>
      )
    }
    const renderTreeNodes = (data) => {
        return data.map((item) => {
          if (item.children) {
            return (
              <TreeNode title={item.title} key={item.key} dataRef={item}>
                {renderTreeNodes(item.children)}
              </TreeNode>
            )
          }
          return <TreeNode {...item} />
        })
    }
    return (<Tree
              showLine
              onSelect={ empty }
              expandedKeys={expandedKeys}
              onExpand={(expanded) => setExpandedKeys(expanded)}
              switcherIcon={<CaretDownOutlined />}
              titleRender={render}
              {...props}
              >
              {renderTreeNodes(props.data || [])}
        </Tree>)
}

export default TreeComponent