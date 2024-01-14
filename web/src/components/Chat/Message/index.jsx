import { memo, createElement } from 'react'
import 'katex/dist/katex.min.css'
import RemarkMath from 'remark-math'
import RemarkBreaks from 'remark-breaks'
import MathJax from 'rehype-mathjax'
import RemarkGfm from 'remark-gfm'
import RemarkGemoji from 'remark-gemoji'
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'
import { oneDark } from 'react-syntax-highlighter/dist/cjs/styles/prism'
import Customize from './Customize'
import { MemoizedReactMarkdown } from './MemoizedReactMarkdown'

const { ReplaceDoubleAt, TransformLink } = Customize

const Component = memo((props) => {
  const { streaming, content, html, loading, source, onClick } = props
  const reg = /(https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,4}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*))/gi
  return loading ? (
        <div className="px-5" style={{paddingTop: "0.6rem", paddingBottom: "0.6rem"}}>
          <div className='spin spin-small'>
          </div>
      </div>
  ) : (
    <div>
      <div onClick={onClick}>
        {html ? (
          <div dangerouslySetInnerHTML={{ __html: html }}></div>
        ) : (
          <MemoizedReactMarkdown
            className=''
            remarkPlugins={[RemarkGfm, RemarkMath, RemarkBreaks, ReplaceDoubleAt, TransformLink, RemarkGemoji]}
            rehypePlugins={[MathJax]}
            components={{
              code({ node, inline, className, children, ...props }) {
                if (children.length) {
                  if (children[0] == '▍') {
                    return <span className='animate-pulse cursor-default mt-1'>▍</span>
                  }
                  children[0] = children[0].replace('`▍`', '▍')
                }
                const match = /language-(\w+)/.exec(className || '')
                return !inline && match ? (
                  <SyntaxHighlighter language={match[1]} style={oneDark} showLineNumbers PreTag='div' value={String(children).replace(/\n$/, '')} customStyle={{ margin: 0 }}></SyntaxHighlighter>
                ) : (
                  <code {...props} className={className}>
                    {children}
                  </code>
                )
              },
              table({ children }) {
                return <table className='border-collapse border border-black px-3 py-1 dark:border-white'>{children}</table>
              },
              th({ children }) {
                return <th className='break-words border border-black bg-gray-500 px-3 py-1 text-white dark:border-white'>{children}</th>
              },
              td({ children }) {
                return <td className='break-words border border-black px-3 py-1 dark:border-white'>{children}</td>
              },
              img({ node, ...props }) {
                return (
                  <>
                    <a className='block embed' data-type='embed' href={props.src}>
                      <img {...props} />
                    </a>
                  </>
                )
              },

              a({ node, children }) {
                const props = {}
                const items = node.properties
                Object.keys(items).forEach((key) => {
                  props[key] = Array.isArray(items[key]) ? items[key].join(' ') : items[key]
                })
                return createElement(node.tagName, props, children)
              }
            }}
          >
            {`${(content || '').replace(reg, ' $1 ')}${streaming ? '`▍`' : ''}`}
          </MemoizedReactMarkdown>
        )}
      </div>
      {source ? '' : ''}
    </div>
  )
})

export default Component
