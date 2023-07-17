import { FC } from 'react'
import { useIntl } from 'react-intl'

import Button from '@shared-ui/components/Atomic/Button'
import { IconInfo } from '@shared-ui/components/Atomic'

import { Props } from './InitializedByAnother.types'
import { messages as t } from './InitializedByAnother.i18n'
import * as styles from './InitializedByAnother.styles'

const InitializedByAnother: FC<Props> = (props) => {
    const { show, logout } = props
    const { formatMessage: _ } = useIntl()

    if (!show) {
        return null
    }

    return (
        <div css={styles.initializedByAnother}>
            <div css={styles.infoBox}>
                <div css={styles.icon}>
                    <IconInfo />
                </div>
                <h1 css={styles.headline}>{_(t.headline)}</h1>
                <div css={styles.description}>{_(t.description)}</div>
                <div css={styles.buttonWrapper}>
                    <Button onClick={logout}>{_(t.logout)}</Button>
                </div>
            </div>
        </div>
    )
}

InitializedByAnother.displayName = 'InitializedByAnother'

export default InitializedByAnother
