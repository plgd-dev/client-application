import { useCallback, useEffect, useMemo, useState } from 'react'
import ReactDOM from 'react-dom'
import { useIntl } from 'react-intl'
import { useParams } from 'react-router-dom'
import { useDispatch } from 'react-redux'

import Footer from '@shared-ui/components/new/Layout/Footer'
import NotFoundPage from '@shared-ui/components/templates/NotFoundPage'
import PageLayout from '@shared-ui/components/new/PageLayout'
import { useIsMounted, WellKnownConfigType } from '@shared-ui/common/hooks'
import { messages as menuT } from '@shared-ui/components/new/Menu/Menu.i18n'
import { showErrorToast, showSuccessToast } from '@shared-ui/components/new/Toast/Toast'
import { BreadcrumbItem } from '@shared-ui/components/new/Breadcrumbs/Breadcrumbs.types'
import { security } from '@shared-ui/common/services'
import StatusTag from '@shared-ui/components/new/StatusTag'
import Breadcrumbs from '@shared-ui/components/new/Layout/Header/Breadcrumbs'
import EditDeviceNameModal from '@shared-ui/components/organisms/EditDeviceNameModal'
import Tabs from '@shared-ui/components/new/Tabs'
import { getApiErrorMessage } from '@shared-ui/common/utils'

import { history } from '@/store'
import { devicesStatuses, NO_DEVICE_NAME, devicesOwnerships, devicesOnboardingStatuses } from '../../constants'
import { handleDeleteDevicesErrors, getDeviceChangeResourceHref } from '../../utils'
import {
    ownDeviceApi,
    disownDeviceApi,
    getDeviceAuthCode,
    onboardDeviceApi,
    offboardDeviceApi,
    PLGD_BROWSER_USED,
    updateDevicesResourceApi,
} from '../../rest'
import DevicesDetailsHeader from '../DevicesDetailsHeader'
import { messages as t } from '../../Devices.i18n'
import { useDeviceDetails, useDevicesResources, useOnboardingButton } from '../../hooks'
import './DevicesDetailsPage.scss'
import { disOwnDevice, ownDevice } from '@/containers/Devices/slice'
import IncompleteOnboardingDataModal, {
    getOnboardingDataFromConfig,
} from '@/containers/Devices/Detail/IncompleteOnboardingDataModal'
import {
    OnboardingDataType,
    onboardingDataDefault,
} from '../IncompleteOnboardingDataModal/IncompleteOnboardingDataModal.types'
import FirstTimeOnboardingModal from '@/containers/Devices/Detail/FirstTimeOnboardingModal/FirstTimeOnboardingModal'
import Tab1 from './Tabs/Tab1'
import Tab2 from './Tabs/Tab2'

const DevicesDetailsPage = () => {
    const { formatMessage: _ } = useIntl()
    const {
        id,
    }: {
        id: string
    } = useParams()
    const [showDpsModal, setShowDpsModal] = useState(false)
    const [showIncompleteOnboardingModal, setShowIncompleteOnboardingModal] = useState(false)
    const [showFirstTimeOnboardingModal, setShowFirstTimeOnboardingModal] = useState(false)
    const [onboardingData, setOnboardingData] = useState<OnboardingDataType>(onboardingDataDefault)
    const [onboarding, setOnboarding] = useState(false)
    const [showEditNameModal, setShowEditNameModal] = useState(false)
    const [domReady, setDomReady] = useState(false)
    const [deviceNameLoading, setDeviceNameLoading] = useState(false)
    const [activeTabItem, setActiveTabItem] = useState(0)

    const isMounted = useIsMounted()
    const { data, updateData, loading, error: deviceError } = useDeviceDetails(id)
    const {
        data: resourcesData,
        loading: loadingResources,
        error: resourcesError,
        refresh: refreshResources,
    } = useDevicesResources(id)
    const dispatch = useDispatch()

    const isOwned = useMemo(() => data?.ownershipStatus === devicesOwnerships.OWNED, [data])
    const resources = useMemo(() => resourcesData?.resources || [], [resourcesData])

    const [
        incompleteOnboardingData,
        onboardResourceLoading,
        deviceOnboardingResourceData,
        refetchDeviceOnboardingData,
    ] = useOnboardingButton({
        resources,
        isOwned,
        deviceId: id,
    })

    const wellKnownConfig = security.getWellKnowConfig() as WellKnownConfigType
    const parseOnboardingData = useCallback(() => getOnboardingDataFromConfig(wellKnownConfig), [wellKnownConfig])
    const handleOpenEditDeviceNameModal = useCallback(() => setShowEditNameModal(true), [])

    useEffect(() => {
        setDomReady(true)
    }, [])

    // check onboarding status evert 1s if onboarding process running
    useEffect(() => {
        const { UNINITIALIZED, REGISTERED, FAILED } = devicesOnboardingStatuses

        if (
            deviceOnboardingResourceData?.content?.cps &&
            ![UNINITIALIZED, REGISTERED, FAILED].includes(deviceOnboardingResourceData.content.cps)
        ) {
            const interval = setInterval(() => {
                refetchDeviceOnboardingData()
            }, 1000)
            return () => clearInterval(interval)
        }
    }, [deviceOnboardingResourceData, refetchDeviceOnboardingData])

    const openDpsModal = useCallback(() => setShowDpsModal(true), [])

    const handleOnboardCallback = useCallback(() => {
        if (deviceOnboardingResourceData.content.cps === devicesOnboardingStatuses.UNINITIALIZED) {
            if (incompleteOnboardingData) {
                setShowIncompleteOnboardingModal(true)
            } else {
                onboardDevice({ ...parseOnboardingData(), authorizationCode: '' }).then()
            }
        } else {
            offboardDeviceApi(id).then(() => {
                setOnboardingData({ ...onboardingData, authorizationCode: '' })
                refetchDeviceOnboardingData()
            })
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [deviceOnboardingResourceData, id, incompleteOnboardingData, onboardingData])

    const deviceName = data?.data?.content?.n || NO_DEVICE_NAME

    const handleOwnChange = useCallback(() => {
        try {
            if (isOwned) {
                disownDeviceApi(id).then(() => {
                    if (isMounted.current) {
                        // @ts-ignore
                        dispatch(disOwnDevice(id))
                        history.push('/')

                        showSuccessToast({
                            title: _(t.deviceDisOwned),
                            message: _(t.deviceWasDisOwned, { name: deviceName }),
                        })
                    }
                })
            } else {
                ownDeviceApi(id).then(() => {
                    if (isMounted.current) {
                        // @ts-ignore
                        dispatch(ownDevice(id))

                        showSuccessToast({
                            title: _(t.deviceOwned),
                            message: _(t.deviceWasOwned, { name: deviceName }),
                        })
                    }
                })
            }
        } catch (error) {
            handleDeleteDevicesErrors(error, _, true)
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [_, id, isMounted, isOwned, deviceName])

    const openOnboardingModal = useCallback(() => {
        toggleOnboardingModal(true)
    }, [])

    if (deviceError) {
        return <NotFoundPage message={_(t.deviceNotFoundMessage, { id })} title={_(t.deviceNotFound)} />
    }

    if (resourcesError) {
        return (
            <NotFoundPage message={_(t.deviceResourcesNotFoundMessage, { id })} title={_(t.deviceResourcesNotFound)} />
        )
    }

    const deviceStatus = data?.metadata?.status?.value
    const isOnline = true
    const isUnregistered = devicesStatuses.UNREGISTERED === deviceStatus

    const breadcrumbs: BreadcrumbItem[] = [
        {
            to: '/',
            label: _(menuT.devices),
        },
    ]

    if (deviceName) {
        breadcrumbs.push({ label: deviceName })
    }

    // Update the device name in the data object
    const updateDeviceNameInData = (name: string) => {
        updateData({
            ...data,
            data: {
                ...data.data,
                content: {
                    ...data.data.content,
                    n: name,
                },
            },
        })
    }

    const onboardDevice = async (onboardingData: OnboardingDataType) => {
        try {
            setOnboarding(true)

            const wasBrowserUsed = localStorage.getItem(PLGD_BROWSER_USED)

            if (!wasBrowserUsed) {
                localStorage.setItem(PLGD_BROWSER_USED, '1')
                setShowFirstTimeOnboardingModal(true)
            }

            const code =
                onboardingData.authorizationCode !== '' ? onboardingData.authorizationCode : await getDeviceAuthCode(id)

            const cleanUpOnboardData = (d: string) => d.replace(/\\n/g, '\n')

            onboardDeviceApi(id, {
                coapGatewayAddress: onboardingData.coapGatewayAddress || '',
                authorizationCode: code as string,
                authorizationProviderName: onboardingData.authorizationProviderName || '',
                hubId: onboardingData.hubId || '',
                certificateAuthorities: cleanUpOnboardData(onboardingData.certificateAuthorities || ''),
            })
                .then((r) => {
                    setOnboarding(false)
                    refetchDeviceOnboardingData()
                })
                .catch((e) => {
                    showErrorToast(JSON.parse(e?.request?.response)?.message || e.message)
                    setOnboardingData(onboardingData)
                    toggleOnboardingModal(true)
                    setOnboarding(false)
                })
        } catch (e: any) {
            if (e !== 'user-cancel') {
                showErrorToast(e.message)
                console.error(e)
            }

            setOnboarding(false)
        }
    }

    function toggleOnboardingModal(state = false) {
        setShowIncompleteOnboardingModal(state)
    }

    const updateDeviceName = async (name: string) => {
        if (name.trim() !== '' && name !== deviceName) {
            const href = getDeviceChangeResourceHref(resources)

            setDeviceNameLoading(true)

            try {
                const { data } = await updateDevicesResourceApi(
                    { deviceId: id, href: href! },
                    {
                        n: name,
                    }
                )

                if (isMounted.current) {
                    setDeviceNameLoading(false)
                    updateDeviceNameInData(data?.n || name)
                }
            } catch (error) {
                if (error && isMounted.current) {
                    showErrorToast({
                        title: _(t.deviceNameChangeFailed),
                        message: getApiErrorMessage(error),
                    })
                    setDeviceNameLoading(false)
                    setShowEditNameModal(false)
                }
            }
        } else {
            setDeviceNameLoading(false)
            setShowEditNameModal(false)
        }
    }

    return (
        <PageLayout
            breadcrumbs={breadcrumbs}
            footer={<Footer footerExpanded={false} paginationComponent={<div id='paginationPortalTarget'></div>} />}
            header={
                <DevicesDetailsHeader
                    deviceId={id}
                    deviceName={deviceName}
                    deviceOnboardingResourceData={deviceOnboardingResourceData}
                    handleOpenEditDeviceNameModal={handleOpenEditDeviceNameModal}
                    incompleteOnboardingData={incompleteOnboardingData}
                    isOwned={isOwned}
                    isUnregistered={isUnregistered}
                    onOwnChange={handleOwnChange}
                    onboardButtonCallback={handleOnboardCallback}
                    onboardResourceLoading={onboardResourceLoading}
                    onboarding={onboarding}
                    openDpsModal={openDpsModal}
                    openOnboardingModal={openOnboardingModal}
                    resources={resources}
                />
            }
            headlineStatusTag={
                <StatusTag variant={isOwned ? 'success' : 'error'}>
                    {isOwned ? _(t.owned) : _(t.disOwnDevice)}
                </StatusTag>
            }
            loading={loading}
            title={deviceName || ''}
        >
            {domReady &&
                ReactDOM.createPortal(
                    <Breadcrumbs items={[{ label: _(menuT.devices), link: '/' }, { label: deviceName }]} />,
                    document.querySelector('#breadcrumbsPortalTarget') as Element
                )}

            <Tabs
                fullHeight={true}
                onItemChange={(i) => setActiveTabItem(i)}
                tabs={[
                    {
                        name: _(t.deviceInformation),
                        content: (
                            <Tab1
                                data={data}
                                deviceId={id}
                                deviceOnboardingResourceData={deviceOnboardingResourceData}
                                isActiveTab={activeTabItem === 0}
                                isOwned={isOwned}
                                onboardResourceLoading={onboardResourceLoading}
                                resources={resources}
                            />
                        ),
                    },
                    {
                        name: _(t.resources),
                        content: (
                            <Tab2
                                closeDpsModal={() => setShowDpsModal(false)}
                                deviceName={deviceName}
                                deviceStatus={deviceStatus}
                                isActiveTab={activeTabItem === 1}
                                isOnline={isOnline}
                                isOwned={isOwned}
                                isUnregistered={isUnregistered}
                                loadingResources={loadingResources}
                                refreshResources={refreshResources}
                                resourcesData={resourcesData}
                                showDpsModal={showDpsModal}
                            />
                        ),
                    },
                ]}
            />

            <IncompleteOnboardingDataModal
                deviceId={id}
                onClose={() => toggleOnboardingModal(false)}
                onSubmit={(onboardingData) => {
                    setOnboardingData(onboardingData)
                    onboardDevice(onboardingData).then()
                }}
                onboardingData={onboardingData}
                show={showIncompleteOnboardingModal}
            />

            <FirstTimeOnboardingModal
                onClose={() => {
                    setShowFirstTimeOnboardingModal(false)
                }}
                onSubmit={() => {
                    setShowFirstTimeOnboardingModal(false)
                }}
                show={showFirstTimeOnboardingModal}
            />

            <EditDeviceNameModal
                deviceName={deviceName}
                deviceNameLoading={deviceNameLoading}
                handleClose={() => setShowEditNameModal(false)}
                handleSubmit={updateDeviceName}
                i18n={{
                    close: _(t.close),
                    deviceName: _(t.deviceName),
                    edit: _(t.edit),
                    name: _(t.name),
                    reset: _(t.reset),
                    saveChange: _(t.saveChange),
                    savingChanges: _(t.savingChanges),
                }}
                show={showEditNameModal}
            />
        </PageLayout>
    )
}

export default DevicesDetailsPage
