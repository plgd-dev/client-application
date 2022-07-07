/* eslint-disable react/jsx-no-target-blank */
import { memo } from 'react'
import { useIntl } from 'react-intl'
import OverlayTrigger from 'react-bootstrap/OverlayTrigger'
import Tooltip from 'react-bootstrap/Tooltip'

import { messages as t } from './footer-i18n'
import './footer.scss'
import { useAppConfig } from '@/containers/app'
import { copyToClipboard } from '@/common/utils'
import { showSuccessToast } from '@/components/toast'

export const Footer = memo(() => {
  const { formatMessage: _ } = useIntl()
  const { version } = useAppConfig()

  const copyVersion = () => {
    copyToClipboard(version)
    showSuccessToast({
      title: _(t.done),
      message: _(t.copied),
    })
  }

  return (
    <footer id="footer">
      <div className="left">
        <OverlayTrigger
          placement="top"
          overlay={
            <Tooltip className="plgd-tooltip" id="footer-version">
              {_(t.clickToCopy)}
            </Tooltip>
          }
        >
          <span className="copy" onClick={copyVersion}>
            <i className={`fas fa-copy m-r-10`} />
          </span>
        </OverlayTrigger>
        <span className="copy" onClick={copyVersion}>
          {version}
        </span>
      </div>
      <div className="right">
        <a
          href="https://github.com/plgd-dev/client-application/blob/main/pb/service.swagger.json"
          target="_blank"
          rel="noopener"
        >
          {_(t.API)}
        </a>
        <a href="https://docs.plgd.dev/" target="_blank" rel="noopener">
          {_(t.docs)}
        </a>
        <a href="https://discord.gg/Pcusx938kg" target="_blank" rel="noopener">
          {_(t.contribute)}
        </a>
      </div>
    </footer>
  )
})
