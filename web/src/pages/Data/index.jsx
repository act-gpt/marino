import React, { useState, useEffect } from 'react'
import { connect } from 'react-redux'
import actions from '@/redux/actions'

const Date = (props) => {
    return (<div>Comming soon</div>)
}



const mapStateToProps = (state) => {
    return {
      user: state.user,
      app: state.app,
    }
}
const mapDispatchToProps = {
    ...actions
}

export default connect(mapStateToProps, mapDispatchToProps)(Date)