import { useIntl } from 'react-intl'
import { messages as t } from '@/containers/App/App.i18n'
import { security } from '@shared-ui/common/services'
import { WellKnownConfigType } from '@shared-ui/common/hooks'

const MockApp = () => {
    const { formatMessage: _ } = useIntl()
    const wellKnowConfig = security.getWellKnowConfig() as WellKnownConfigType

    console.log({ wellKnowConfig })

    return (
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100%' }}>
            <div style={{ fontSize: 16 }}>
                {_(t.mockPart1)}&nbsp;
                <strong>{wellKnowConfig.remoteProvisioning?.authority}</strong>
                {_(t.mockPart2)}
            </div>
        </div>
    )
}

MockApp.displayName = 'MockApp'

export default MockApp
