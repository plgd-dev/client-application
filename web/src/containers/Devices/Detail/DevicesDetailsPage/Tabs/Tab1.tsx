import { FC, useEffect, useMemo, useState } from 'react'
import { useIntl } from 'react-intl'
import classNames from 'classnames'

import SimpleStripTable from '@shared-ui/components/Atomic/SimpleStripTable'
import TagGroup from '@shared-ui/components/Atomic/TagGroup'
import Tag from '@shared-ui/components/Atomic/Tag'
import Badge from '@shared-ui/components/Atomic/Badge'

import { messages as t } from '@/containers/Devices/Devices.i18n'
import { Props } from './Tab1.types'
import {
    getColorByOnboardingStatus,
    getColorByProvisionStatus,
    getDPSEndpoint,
    loadResourceData,
} from '@/containers/Devices/utils'
import testId from '@/testId'
import { DEVICE_PROVISION_STATUS_DELAY_MS } from '@/containers/Devices/constants'
import { IconLoader } from '@shared-ui/components/Atomic/Loader'

const Tab1: FC<Props> = (props) => {
    const { data, deviceId, deviceOnboardingResourceData, isOwned, isUnsupported, resources, onboardResourceLoading } =
        props
    const { formatMessage: _ } = useIntl()

    const { onboardTitleStatus } = testId.devices.detail

    const [resourceLoading, setResourceLoading] = useState(false)
    const [deviceResourceData, setDeviceResourceData] = useState<any>(undefined)

    const provisionStatus = deviceResourceData?.content?.provisionStatus
    const onboardingStatus = deviceOnboardingResourceData?.content?.cps || 'n/a'
    const dpsEndpoint = useMemo(() => getDPSEndpoint(resources), [resources])

    useEffect(() => {
        if (dpsEndpoint && isOwned && !deviceResourceData) {
            setResourceLoading(true)
            setTimeout(() => {
                loadResourceData({
                    href: dpsEndpoint.href,
                    deviceId,
                    errorCallback: () => {
                        setResourceLoading(false)
                    },
                }).then((rData) => {
                    setDeviceResourceData(rData)
                    setResourceLoading(false)
                })
            }, DEVICE_PROVISION_STATUS_DELAY_MS)
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [dpsEndpoint, isOwned, deviceResourceData])

    const rows = [
        { attribute: _(t.id), value: deviceId },
        {
            attribute: _(t.types),
            value: data?.types ? (
                <TagGroup>
                    {data?.types.map((t, key) => (
                        <Tag key={t}>{t}</Tag>
                    ))}
                </TagGroup>
            ) : (
                <div>-</div>
            ),
        },
        {
            attribute: _(t.ownershipStatus),
            value: (
                <Badge
                    className={classNames({
                        green: isOwned,
                        red: !isOwned,
                    })}
                >
                    {isOwned ? _(t.owned) : _(t.unowned)}
                </Badge>
            ),
            hidden: isUnsupported,
        },
        {
            attribute: _(t.onboardingStatus),
            value: onboardResourceLoading ? (
                <IconLoader size={20} type='secondary' />
            ) : (
                <Badge className={getColorByOnboardingStatus(onboardingStatus)} data-test-id={onboardTitleStatus}>
                    {onboardingStatus}
                </Badge>
            ),
        },
        {
            attribute: _(t.dpsStatus),
            value: resourceLoading ? (
                <IconLoader size={20} type='secondary' />
            ) : (
                <Badge className={isOwned ? getColorByProvisionStatus(provisionStatus) : 'grey'}>
                    {isOwned && deviceResourceData ? provisionStatus : _(t.notAvailable)}
                </Badge>
            ),
            hidden: !dpsEndpoint,
        },
        {
            attribute: _(t.endpoints),
            value: data?.endpoints ? (
                <TagGroup>
                    {data?.endpoints?.map?.((endpoint: string) => (
                        <Tag key={endpoint}>{endpoint}</Tag>
                    ))}
                </TagGroup>
            ) : (
                <div>-</div>
            ),
        },
    ]

    return (
        <div
            style={{
                paddingTop: 8,
                overflow: 'hidden',
            }}
        >
            <SimpleStripTable rows={rows.filter((r) => r?.hidden === false || r?.hidden === undefined)} />
        </div>
    )
}

Tab1.displayName = 'Tab1'

export default Tab1
