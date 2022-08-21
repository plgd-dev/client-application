import { createElement, FC, memo, useEffect, useMemo, useState } from 'react'
import { useIntl } from 'react-intl'
import classNames from 'classnames'
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import Badge from '@shared-ui/components/new/Badge'
import Label from '@shared-ui/components/new/Label'
import { getValue } from '@shared-ui/common/utils'
import {
  DEVICE_PROVISION_STATUS_DELAY_MS,
  DEVICE_TYPE_OIC_WK_D,
  devicesStatuses,
} from '../../constants'
import { messages as t } from '../../Devices.i18n'
import {
  getColorByProvisionStatus,
  getDPSEndpoint,
  handleFetchResourceErrors,
} from '@/containers/devices/utils'
import { getDevicesResourcesApi } from '@/containers/devices/rest'
import omit from 'lodash/omit'
import Display from '@shared-ui/components/new/Display'
import { Props } from './DevicesDetails.types'
import { useIsMounted } from '@shared-ui/common/hooks/use-is-mounted'

const DevicesDetails: FC<Props> = memo(
  ({ data, loading, isOwned, resources, deviceId }) => {
    const { formatMessage: _ } = useIntl()
    const [resourceLoading, setResourceLoading] = useState(false)
    const [deviceResourceData, setDeviceResourceData] = useState<any>(undefined)
    const deviceStatus = data?.metadata?.status?.value
    const isUnregistered = devicesStatuses.UNREGISTERED === deviceStatus
    const isMounted = useIsMounted()
    const LabelWithLoading = (p: any) =>
      createElement(Label, {
        ...omit(p, 'loading'),
        inline: true,
        className: classNames({
          shimmering: loading || p.loading,
          'grayed-out': isUnregistered,
        }),
      } as any)

    const dpsEndpoint = useMemo(() => getDPSEndpoint(resources), [resources])

    const loadResourceData = async (href: string) => {
      try {
        const { data: deviceData } = await getDevicesResourcesApi({
          deviceId,
          href,
        })

        isMounted.current && setResourceLoading(false)

        return deviceData.data
      } catch (error) {
        if (error && isMounted.current) {
          handleFetchResourceErrors(error, _)
          setResourceLoading(false)
        }
      }
    }

    useEffect(() => {
      if (dpsEndpoint && isOwned && !deviceResourceData) {
        setResourceLoading(true)
        setTimeout(() => {
          loadResourceData(dpsEndpoint.href).then(rData => {
            setDeviceResourceData(rData)
            setResourceLoading(false)
          })
        }, DEVICE_PROVISION_STATUS_DELAY_MS)
      }
      // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [dpsEndpoint, isOwned, deviceResourceData])

    const provisionStatus = deviceResourceData?.content?.provisionStatus

    return (
      <Row>
        <Col>
          <LabelWithLoading title="ID">{getValue(data?.id)}</LabelWithLoading>
          <LabelWithLoading title={_(t.types)}>
            <div className="align-items-end badges-box-vertical">
              {data?.types
                ?.filter((type: string) => type !== DEVICE_TYPE_OIC_WK_D)
                .map?.((type: string) => (
                  <Badge key={type}>{type}</Badge>
                ))}
            </div>
          </LabelWithLoading>
          <LabelWithLoading title={_(t.ownershipStatus)}>
            <Badge
              className={classNames({
                green: isOwned,
                red: !isOwned,
              })}
            >
              {isOwned ? _(t.owned) : _(t.unowned)}
            </Badge>
          </LabelWithLoading>
          <Display when={!!dpsEndpoint}>
            <LabelWithLoading title={_(t.dpsStatus)} loading={resourceLoading}>
              <Badge
                className={
                  isOwned ? getColorByProvisionStatus(provisionStatus) : 'grey'
                }
              >
                {isOwned && deviceResourceData
                  ? provisionStatus
                  : _(t.notAvailable)}
              </Badge>
            </LabelWithLoading>
          </Display>
        </Col>
        <Col>
          <LabelWithLoading title={_(t.endpoints)}>
            <div className="align-items-end badges-box-vertical">
              {data?.endpoints?.map?.((endpoint: string) => (
                <Badge key={endpoint}>{endpoint}</Badge>
              ))}
            </div>
          </LabelWithLoading>
        </Col>
      </Row>
    )
  }
)

DevicesDetails.displayName = 'DevicesDetails'

export default DevicesDetails
