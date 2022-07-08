/* eslint-disable react/jsx-no-target-blank */
import { memo, useState } from 'react'
import { useIntl } from 'react-intl'
import OverlayTrigger from 'react-bootstrap/OverlayTrigger'
import Tooltip from 'react-bootstrap/Tooltip'
import Popover from 'react-bootstrap/Popover'
import PopoverContent from 'react-bootstrap/PopoverContent'
import { format } from 'date-fns'

import { messages as t } from './footer-i18n'
import './footer.scss'
import { useAppConfig } from '@/containers/app'
import { copyToClipboard } from '@/common/utils'
import { showSuccessToast } from '@/components/toast'
import { Display } from '@/components/display'

export const Footer = memo(() => {
  const { formatMessage: _ } = useIntl()
  const { buildInformation } = useAppConfig()
  const [show, setShow] = useState(false)
  const handleOnMouseEnter = () => {
    setShow(true)
  }
  const handleOnMouseLeave = () => {
    setShow(false)
  }

  const copyVersion = () => {
    copyToClipboard(buildInformation.version)
    showSuccessToast({
      title: _(t.done),
      message: _(t.copied),
    })
  }

  const BuildInformation = () => {
    if (!buildInformation?.buildDate) {
      return null
    }
    return (
      <Display when={!!buildInformation}>
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

        <OverlayTrigger
          placement="top-start"
          show={show}
          delay={{ hide: 1000 }}
          overlay={
            <Popover
              id="popover-positioned-top"
              className="plgd-popover interactive"
              onMouseEnter={handleOnMouseEnter}
              onMouseLeave={handleOnMouseLeave}
            >
              <PopoverContent>
                <div className="footerBuildInfo">
                  <div className="line">
                    <div className="attr">Build time:</div>
                    <div className="val">
                      {format(new Date(buildInformation.buildDate), 'Pp')}
                    </div>
                  </div>
                  <div className="line">
                    <div className="attr">Commit hash:</div>
                    <a
                      className="val"
                      href={`https://github.com/plgd-dev/client-application/commit/${buildInformation.commitHash}`}
                      target="_blank"
                    >
                      {buildInformation.commitHash}
                    </a>
                  </div>
                </div>
              </PopoverContent>
            </Popover>
          }
        >
          <span
            onMouseEnter={handleOnMouseEnter}
            onMouseLeave={handleOnMouseLeave}
          >
            {buildInformation.version}
          </span>
        </OverlayTrigger>
      </Display>
    )
  }

  return (
    <footer id="footer">
      <div className="left">
        <BuildInformation />
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
