import { createElement, memo } from 'react'
import { useIntl } from 'react-intl'
import PropTypes from 'prop-types'
import classNames from 'classnames'
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import { Badge } from '@/components/badge'
import { Label } from '@/components/label'
import { getValue } from '@/common/utils'
import { DEVICE_TYPE_OIC_WK_D, devicesStatuses } from './constants'
import { deviceShape } from './shapes'
import { messages as t } from './devices-i18n'

export const DevicesDetails = memo(({ data, loading, isOwned }) => {
  const { formatMessage: _ } = useIntl()
  const deviceStatus = data?.metadata?.status?.value
  const isUnregistered = devicesStatuses.UNREGISTERED === deviceStatus
  const LabelWithLoading = p =>
    createElement(Label, {
      ...p,
      inline: true,
      className: classNames({
        shimmering: loading,
        'grayed-out': isUnregistered,
      }),
    })

  return (
    <Row>
      <Col>
        <LabelWithLoading title="ID">{getValue(data?.id)}</LabelWithLoading>
        <LabelWithLoading title={_(t.types)}>
          <div className="align-items-end badges-box-vertical">
            {data?.types?.map?.(type =>
              type !== DEVICE_TYPE_OIC_WK_D ? (
                <Badge key={type}>{type}</Badge>
              ) : null
            )}
          </div>
        </LabelWithLoading>
      </Col>
      <Col>
        <LabelWithLoading title={_(t.ownershipStatus)}>
          <Badge className={isOwned ? 'green' : 'red'}>
            {isOwned ? _(t.owned) : _(t.unowned)}
          </Badge>
        </LabelWithLoading>
        <LabelWithLoading title={_(t.endpoints)}>
          <div className="align-items-end badges-box-vertical">
            {data?.endpoints?.map?.(endpoint => (
              <Badge key={endpoint}>{endpoint}</Badge>
            ))}
          </div>
        </LabelWithLoading>
      </Col>
    </Row>
  )
})

DevicesDetails.propTypes = {
  data: deviceShape,
  loading: PropTypes.bool.isRequired,
}

DevicesDetails.defaultProps = {
  data: null,
}
