import { FC, useMemo } from 'react'
import { useSelector } from 'react-redux'

import { DevicesListPage as DevicesListPageCore, Props } from '@shared-ui/app/clientApp/Devices/List/DevicesListPage'
import { security } from '@shared-ui/common/services'
import { hasDifferentOwner } from '@shared-ui/common/services/api-utils'

import { CombinedStoreType } from '@/store/store'
import { useAppInitialization } from '@shared-ui/app/clientApp/Devices/hooks'

const DevicesListPage: FC<Props> = (props) => {
    const wellKnownConfig = security.getWellKnowConfig()
    const appStore = useSelector((state: CombinedStoreType) => state.app)

    const diffOwner = useMemo(
        () => hasDifferentOwner(wellKnownConfig, appStore.userWellKnownConfig),
        [wellKnownConfig, appStore.userWellKnownConfig]
    )

    const [initialize, initializationLoading] = useAppInitialization({
        wellKnownConfig,
        loading: diffOwner,
        clientData: appStore.userWellKnownConfig,
    })

    return (
        <DevicesListPageCore
            {...props}
            initializedByAnother={diffOwner || !initialize}
            loading={initializationLoading}
        />
    )
}

DevicesListPage.displayName = 'DevicesListPage'

export default DevicesListPage
