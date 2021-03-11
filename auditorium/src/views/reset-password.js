/**
 * Copyright 2020 - Offen Authors <hioffen@posteo.de>
 * SPDX-License-Identifier: Apache-2.0
 */

/** @jsx h */
const { h } = require('preact')
const { connect } = require('react-redux')

const authentication = require('./../action-creators/authentication')
const withLayout = require('./components/_shared/with-layout')
const useAutofocus = require('./components/_shared/use-autofocus')
const Form = require('./components/reset-password/form')

const ResetPasswordView = (props) => {
  const { matches } = props
  const { token } = matches
  const autofocusRef = useAutofocus()
  return (
    <div class='mw8 center mt4 mb2 br0 br2-ns'>
      <Form
        onResetPassword={props.handleReset}
        ref={autofocusRef}
        token={token}
      />
    </div>
  )
}

const mapDispatchToProps = {
  handleReset: authentication.resetPassword
}

module.exports = connect(null, mapDispatchToProps)(
  withLayout()(
    ResetPasswordView
  )
)
