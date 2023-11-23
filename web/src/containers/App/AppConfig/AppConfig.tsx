import React, { FC, useMemo } from 'react'
import { useIntl } from 'react-intl'
import { useDispatch, useSelector } from 'react-redux'

import IconSettings from '@shared-ui/components/Atomic/Icon/components/IconSettings'
import { convertSize } from '@shared-ui/components/Atomic'
import FloatingPanel from '@shared-ui/components/Atomic/FloatingPanel/FloatingPanel'
import FormGroup from '@shared-ui/components/Atomic/FormGroup'
import FormLabel from '@shared-ui/components/Atomic/FormLabel'
import FormSelect from '@shared-ui/components/Atomic/FormSelect'
import Notification from '@shared-ui/components/Atomic/Notification/Toast'

import { messages as t } from '../App.i18n'
import { CombinedStoreType } from '@/store/store'
import { setTheme } from '@/containers/App/slice'

const AppConfig: FC<any> = () => {
    const { formatMessage: _ } = useIntl()
    const appStore = useSelector((state: CombinedStoreType) => state.app)
    const dispatch = useDispatch()
    const options = useMemo(
        () => appStore.configuration.themes.map((t) => ({ value: t, label: t })),
        [appStore.configuration.themes]
    )
    const defTheme = useMemo(
        () => options.find((o) => o.value === appStore.configuration?.theme) || options[0],
        [appStore, options]
    )

    const onChange = (v: any) => {
        if (v !== defTheme) {
            dispatch(setTheme(v.value))

            Notification.success({ title: _(t.configurationUpdated), message: _(t.configurationUpdatedMessage) })
        }
    }

    return (
        <FloatingPanel reference={<IconSettings {...convertSize(20)} />} title='Config'>
            <div style={{ padding: '16px 0' }}>
                <FormGroup id='form-group-1' marginBottom={false}>
                    <FormLabel text={_(t.theme)} />
                    <FormSelect defaultValue={defTheme} onChange={onChange} options={options} />
                </FormGroup>
            </div>
        </FloatingPanel>
    )
}

export default AppConfig
