import PageLoader from '@shared-ui/components/Atomic/PageLoader'
import { useIntl } from 'react-intl'

import { messages as t } from '@/containers/App/App.i18n'

const AppLoader = () => {
    const { formatMessage: _ } = useIntl()

    return (
        <>
            <PageLoader loading className='auth-loader' />
            <div className='page-loading-text'>{`${_(t.loading)}...`}</div>
        </>
    )
}

AppLoader.displayName = 'AppLoader'

export default AppLoader
