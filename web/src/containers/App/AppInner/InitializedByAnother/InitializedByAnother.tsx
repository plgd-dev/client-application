import { FC } from 'react'
import { Props } from './InitializedByAnother.types'
import './InitializedByAnother.scss'
import { useIntl } from 'react-intl'
import { messages as t } from './InitializedByAnother.i18n'
import Button from '@shared-ui/components/new/Button'

const InitializedByAnother: FC<Props> = (props) => {
    const { show, logout } = props
    const { formatMessage: _ } = useIntl()

    if (!show) {
        return null
    }

    return (
        <div className='initialized-by-another'>
            <div className='info-box'>
                <div className='icon'>
                    <i className='fas fa-info-circle'></i>
                </div>
                <h1>{_(t.headline)}</h1>
                <div className='description'>{_(t.description)}</div>
                <div className='button-wrapper'>
                    <Button onClick={logout}>{_(t.logout)}</Button>
                </div>
            </div>
        </div>
    )
}

InitializedByAnother.displayName = 'InitializedByAnother'

export default InitializedByAnother
