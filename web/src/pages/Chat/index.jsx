import { connect } from 'react-redux'
import actions from '@/redux/actions'
import Component from '@/components/Chat'

const Chat = (props) => {
  return(<div className="pt-4">
    <Component {...props} />
  </div>)
}

const mapStateToProps = (state) => {
  return {
    id: state.bot.id,
    token: state.bot.access_token,
    user: state.user.name,
    app: state.app,
    dev: true
  }
}

const mapDispatchToProps = {
  ...actions
}

export default connect(mapStateToProps, mapDispatchToProps)(Chat)
