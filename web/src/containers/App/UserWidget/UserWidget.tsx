import { FC, memo } from 'react'
import { useAuth } from 'oidc-react'

import UserWidgetShared from '@shared-ui/components/Layout/Header/UserWidget'
import { Props as UserWidgetProps } from '@shared-ui/components/Layout/Header/UserWidget/UserWidget.types'

type Props = Omit<UserWidgetProps, 'name'>

const UserWidget: FC<Props> = memo((props) => {
    const { userData } = useAuth()

    return (
        <UserWidgetShared
            {...props}
            description={userData?.profile?.family_name || ''}
            image={userData?.profile?.picture}
            name={userData?.profile?.name || ''}
        />
    )
})

UserWidget.displayName = 'UserWidget'

export default UserWidget
