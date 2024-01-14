import { useCallback } from 'react'
import { useEditor, EditorContent, BubbleMenu, FloatingMenu } from '@tiptap/react'
import StarterKit from '@tiptap/starter-kit'
// import Link from '@tiptap/extension-link'
// import Image from '@tiptap/extension-image'
import Placeholder from '@tiptap/extension-placeholder'
import UniqueID from '@tiptap-pro/extension-unique-id'
import { useTranslation } from 'react-i18next'
import './style.scss'

// Video see: https://www.codemzy.com/blog/tiptap-video-embed-extension
// https://github.com/sereneinserenade/tiptap-extension-video

const Tiptap = (props) => {
  const { t } = useTranslation()

  function uploadImage(file) {
    const data = new FormData()
    data.append('file', file)
    //return axios.post('/documents/image/upload', data);
  }

  const editor = useEditor({
    extensions: [
      StarterKit,
      Placeholder.configure({
        placeholder: t('knowledge.tiptap.placeholder')
      }),
      UniqueID.configure({
        types: ['heading', 'paragraph'],
      }),
    ],
    editorProps: {
      handleDrop: function (view, event, slice, moved) {
        /*
        if (!moved && event.dataTransfer && event.dataTransfer.files && event.dataTransfer.files[0]) {
          // if dropping external files
          let file = event.dataTransfer.files[0] // the dropped file
          let filesize = (file.size / 1024 / 1024).toFixed(4) // get the filesize in MB
          if ((file.type === 'image/jpeg' || file.type === 'image/png') && filesize < 10) {
            // check valid image type under 10MB
            // check the dimensions
            let _URL = window.URL || window.webkitURL
            let img = new Image() 
            img.src = _URL.createObjectURL(file)
            img.onload = function () {
              if (this.width > 5000 || this.height > 5000) {
                window.alert('Your images need to be less than 5000 pixels in height and width.') // display alert
              } else {
                // valid image so upload to server
                // uploadImage will be your function to upload the image to the server or s3 bucket somewhere
                uploadImage(file)
                  .then(function (response) {
                    // response is the image url for where it has been saved
                    // pre-load the image before responding so loading indicators can stay
                    // and swaps out smoothly when image is ready
                    let image = new Image()
                    image.src = response
                    image.onload = function () {
                      // place the now uploaded image in the editor where it was dropped
                      const { schema } = view.state
                      const coordinates = view.posAtCoords({ left: event.clientX, top: event.clientY })
                      const node = schema.nodes.image.create({ src: response }) // creates the image element
                      const transaction = view.state.tr.insert(coordinates.pos, node) // places it in the correct position
                      return view.dispatch(transaction)
                    }
                  })
                  .catch(function (error) {
                    if (error) {
                      window.alert('There was a problem uploading your image, please try again.')
                    }
                  })
              }
            }
          } else {
            window.alert('Images need to be in jpg or png format and less than 10mb in size.')
          }
          return true // handled
        }
        */
        return false // not handled use default behaviour
      }
    },
    content: props.html || '',
    onUpdate({ editor }) {
      var html = editor.getHTML()
      props.onChange && props.onChange(html)
    }
  })

  const setLink = useCallback(() => {
    const previousUrl = editor.getAttributes('link').href
    const url = window.prompt('URL', previousUrl)
    // cancelled
    if (url === null) {
      return
    }

    // empty
    if (url === '') {
      editor.chain().focus().extendMarkRange('link').unsetLink().run()
      return
    }
    // update link
    editor.chain().focus().extendMarkRange('link').setLink({ href: url }).run()
  }, [editor])

  return (
    <>
      {editor && (
        <>
          <BubbleMenu className='bubble-menu' tippyOptions={{ duration: 100 }} editor={editor}>
            <button onClick={() => editor.chain().focus().toggleBold().run()} className={editor.isActive('bold') ? 'is-active' : ''}>
              {t('knowledge.tiptap.bold')}
            </button>
            <button onClick={() => editor.chain().focus().toggleItalic().run()} className={editor.isActive('italic') ? 'is-active' : ''}>
              {t('knowledge.tiptap.italic')}
            </button>
            <button onClick={() => editor.chain().focus().unsetAllMarks().run()}>{t('knowledge.tiptap.clean')}</button>
          </BubbleMenu>
          <FloatingMenu className='floating-menu' tippyOptions={{ duration: 100 }} editor={editor}>
            <button onClick={() => editor.chain().focus().toggleHeading({ level: 1 }).run()} className={editor.isActive('heading', { level: 1 }) ? 'is-active' : ''}>
              {t('knowledge.tiptap.h1')}
            </button>
            <button onClick={() => editor.chain().focus().toggleHeading({ level: 2 }).run()} className={editor.isActive('heading', { level: 2 }) ? 'is-active' : ''}>
              {t('knowledge.tiptap.h2')}
            </button>
            <button onClick={() => editor.chain().focus().toggleHeading({ level: 4 }).run()} className={editor.isActive('heading', { level: 2 }) ? 'is-active' : ''}>
              {t('knowledge.tiptap.h4')}
            </button>
            <button onClick={() => editor.chain().focus().toggleBulletList().run()} className={editor.isActive('bulletList') ? 'is-active' : ''}>
              {t('knowledge.tiptap.list')}
            </button>
          </FloatingMenu>
        </>
      )}
      <EditorContent editor={editor} />
    </>
  )
}

export default Tiptap
