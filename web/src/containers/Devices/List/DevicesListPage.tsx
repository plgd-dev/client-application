import { FC } from 'react'

import { DevicesListPage as DevicesListPageCore, Props } from '@shared-ui/app/clientApp/Devices/List/DevicesListPage'

const DevicesListPage: FC<Props> = (props) => <DevicesListPageCore {...props} />

DevicesListPage.displayName = 'DevicesListPage'

export default DevicesListPage
