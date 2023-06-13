import { FC, memo } from 'react'
import { useAuth } from 'oidc-react'
import { useIntl } from 'react-intl'

import UserWidgetShared from '@shared-ui/components/Layout/Header/UserWidget'

import { messages as t } from '@/containers/App/App.i18n'

type Props = {
    logout: () => void
}

const UserWidget: FC<Props> = memo((props) => {
    const { logout } = props
    const { formatMessage: _ } = useIntl()
    const { userData } = useAuth()
    return (
        <UserWidgetShared
            description={userData?.profile?.family_name || ''}
            dropdownItems={[{ title: _(t.logOut), onClick: logout }]}
            image={userData?.profile?.picture}
            name={userData?.profile?.name || ''}
        />
    )
})

UserWidget.displayName = 'UserWidget'

export default UserWidget
