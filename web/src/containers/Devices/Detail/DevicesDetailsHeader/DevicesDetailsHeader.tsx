import { FC, memo, useMemo, useRef } from 'react'
import { useIntl } from 'react-intl'
import { useSelector } from 'react-redux'

import SplitButton from '@shared-ui/components/new/SplitButton'
import Button from '@shared-ui/components/new/Button'
import { Icon } from '@shared-ui/components/new/Icon'

import { canChangeDeviceName, canSetDPSEndpoint, getDeviceNotificationKey } from '../../utils'
import { isNotificationActive } from '../../slice'
import { messages as t } from '../../Devices.i18n'
import { Props } from './DevicesDetailsHeader.types'
import { devicesOnboardingStatuses } from '@/containers/Devices/constants'
import testId from '@/testId'
import * as styles from './DevicesDetailsHeader.styles'

export const DevicesDetailsHeader: FC<Props> = memo((props) => {
    const {
        deviceId,
        isUnregistered,
        onOwnChange,
        isOwned,
        resources,
        openDpsModal,
        onboardResourceLoading,
        onboardButtonCallback,
        deviceOnboardingResourceData,
        incompleteOnboardingData,
        openOnboardingModal,
        onboarding,
        handleOpenEditDeviceNameModal,
    } = props
    const { formatMessage: _ } = useIntl()
    const deviceNotificationKey = getDeviceNotificationKey(deviceId)
    const notificationsEnabled = useRef(false)
    notificationsEnabled.current = useSelector(isNotificationActive(deviceNotificationKey))

    const hasDPS = useMemo(() => canSetDPSEndpoint(resources), [resources])
    const canUpdate = useMemo(() => canChangeDeviceName(resources) && isOwned, [resources, isOwned])

    const hasOnboardButton = deviceOnboardingResourceData?.content?.cps
    const isOnboarded = hasOnboardButton !== devicesOnboardingStatuses.UNINITIALIZED
    const { offboardButton, onboardButton, onboardButtonDropdown } = testId.devices.detail

    return (
        <div css={styles.header}>
            {canUpdate && (
                <Button
                    disabled={isUnregistered}
                    icon={<Icon icon='edit' />}
                    onClick={handleOpenEditDeviceNameModal}
                    style={{ marginLeft: 8 }}
                    variant='tertiary'
                >
                    {_(t.editName)}
                </Button>
            )}
            {hasOnboardButton && (incompleteOnboardingData || isOnboarded) && (
                <Button
                    dataTestId={isOnboarded ? offboardButton : onboardButton}
                    disabled={!isOwned || onboardResourceLoading || onboarding}
                    icon={<Icon icon={isOnboarded ? 'close' : 'plus'} />}
                    loading={onboardResourceLoading || onboarding}
                    onClick={onboardButtonCallback}
                    variant='tertiary'
                >
                    {isOnboarded ? _(t.offboardDevice) : _(t.onboardDevice)}
                </Button>
            )}
            {hasOnboardButton &&
                !incompleteOnboardingData &&
                hasOnboardButton === devicesOnboardingStatuses.UNINITIALIZED && (
                    <div className='m-r-10'>
                        <SplitButton
                            dataTestId={isOnboarded ? offboardButton : onboardButton}
                            dataTestIdDropdown={onboardButtonDropdown}
                            disabled={onboardResourceLoading || onboarding}
                            icon='fa-plus'
                            items={[
                                {
                                    onClick: openOnboardingModal,
                                    label: _(t.changeOnboardingData),
                                    icon: 'edit',
                                },
                            ]}
                            loading={onboardResourceLoading || onboarding}
                            menuProps={{
                                placement: 'bottom-end',
                            }}
                            onClick={onboardButtonCallback}
                            variant='tertiary'
                        >
                            {_(t.onboardDevice)}
                        </SplitButton>
                    </div>
                )}
            <Button
                disabled={isUnregistered}
                icon={<Icon icon={isOwned ? 'close' : 'plus'} />}
                onClick={onOwnChange}
                variant='tertiary'
            >
                {isOwned ? _(t.disOwnDevice) : _(t.ownDevice)}
            </Button>
            {hasDPS && (
                <Button
                    className='m-l-10'
                    disabled={!isOwned}
                    icon={<Icon icon='network' />}
                    onClick={openDpsModal}
                    variant='tertiary'
                >
                    {_(t.setDpsEndpoint)}
                </Button>
            )}
        </div>
    )
})

DevicesDetailsHeader.displayName = 'DevicesDetailsHeader'

export default DevicesDetailsHeader
