/** @jsx h */
const { h } = require('preact')

const Share = require('./../_shared/share')

const ShareAccounts = (props) => {
  return (
    <Share
      {...props}
      headline={__('Share all accounts')}
      subline={__('Share all your Offen accounts via email invitation. Invited users have full access to shared accounts.')}
    />
  )
}

module.exports = ShareAccounts
