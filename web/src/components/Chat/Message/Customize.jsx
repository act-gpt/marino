import { findAndReplace } from 'mdast-util-find-and-replace'
import { visit } from 'unist-util-visit'
import absolute from 'is-absolute-url'
// workflow
// markdown to mdast -> mdast to hset -> hset to html
const replaceDoubleAt = (tree, opts) => {
  // find @@
  return findAndReplace(tree, [
    /*
    [/\@\@(([\u4E00-\u9FA5？，])+)/g, function ($0, $1) {
        return {
            type: 'link', 
            url: '#' + $1,
            title: $1,
            data: {
                hProperties: {
                    className: ['query', 'inline-block'],
                    rel: ['nofollow'],
                    target: '_self',
                    "data-type":["query"],
                }
            },
            children: [
                {type: 'text', value: $1}
            ]
        }
    }],
    */
    [
      /\@\@(([^@\n])+)/g,
      function ($0, $1) {
        $1 = $1.trim()
        return {
          type: 'link',
          url: '#' + $1,
          title: $1,
          data: {
            hProperties: {
              className: ['query', 'inline-block'],
              rel: ['nofollow'],
              target: '_self',
              'data-type': ['query']
            }
          },
          children: [{ type: 'text', value: $1 }]
        }
      }
    ]
  ])
}

const transformLink = (tree) => {
  const target = '_blank'
  const defaultRel = ['nofollow', 'noopener', 'noreferrer']
  const protocols = ['http', 'https']
  visit(tree, (node) => {
    const data = node.data || (node.data = {})
    const props = data.hProperties || {}
    if (node.type === 'image') {
      props.className = ['img', 'inline-block', 'embed']
      node.data.hProperties = props
    }
    if (node.type === 'paragraph') {
      props['x-act-mark'] = ['ok']
      node.data.hProperties = props
    }
    if (node.type === 'link') {
      const protocol = node.url.slice(0, node.url.indexOf(':'))
      if (absolute(node.url) && protocols.includes(protocol)) {
        // copy hProperties
        props.target = target
        props.rel = [defaultRel[0]]
        props.className = ['opener', 'inline']
        props['data-type'] = ['opener']
        node.data.hProperties = props
      }
    }
  })
}

export default {
  ReplaceDoubleAt: () => replaceDoubleAt,
  TransformLink: () => transformLink
}
