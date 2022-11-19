import './PreSharedKeySetup.scss'
import LogoPlgd from './LogoPlgd'
import Label from '../../../../shared-ui/src/components/new/Label'
import TextField from '../../../../shared-ui/src/components/new/TextField'
import classNames from 'classnames'
import { useIntl } from 'react-intl'
import { messages as t } from './PreSharedKeySetup.i18n'
import { useState } from 'react'
import Button from '../../../../shared-ui/src/components/new/Button'
import Form from 'react-bootstrap/Form'
import { initializedByPreShared } from '@/containers/App/AppRest'
import { Props } from './PreSharedKeySetup.types'
import { copyToClipboard } from '../../../../shared-ui/src/common/utils'
import { showSuccessToast, showErrorToast } from '../../../../shared-ui/src/components/old/toast'

const PreSharedKeySetup = (props: Props) => {
    const { setInitialize } = props
    const { formatMessage: _ } = useIntl()
    const [uuid, setUuid] = useState<string>('')
    const [key, setKey] = useState<string>('')

    const handleSubmit = () => {
        initializedByPreShared(uuid, key)
            .then((r) => {
                if (r.status === 200) {
                    setInitialize(true)
                }
            })
            .catch((e) => {
                showErrorToast({
                    title: _(t.error),
                    message: e.response.data.message,
                })
            })
    }

    const handleCopy = (data: string) => {
        copyToClipboard(data)
        showSuccessToast({
            title: _(t.done),
            message: _(t.copied),
        })
    }

    return (
        <div className='preSharedKeySetupPage'>
            <div className='colLeft'>
                <div className='top'>
                    <h1>{_(t.reminder)}</h1>
                    <p className='claim'>{_(t.reminderDescription)}</p>
                </div>
                <div className='bottom'>
                    <LogoPlgd />
                </div>
            </div>
            <div className='colRight'>
                <div className='formContainer'>
                    <h2>Pre shared key setup</h2>
                    <div className='fromWrapper'>
                        <form action=''>
                            <Label title={_(t.subjectId)} onClick={(e) => e.preventDefault()}>
                                <TextField
                                    className={classNames({ error: false })}
                                    value={uuid}
                                    name='subjectId'
                                    autoComplete='subjectId'
                                    onChange={(e) => setUuid(e.target.value)}
                                />
                                <span className='copy' onClick={() => handleCopy(uuid)}>
                                    <i className={`fas fa-copy`} />
                                </span>
                            </Label>
                            <Label title={_(t.key)} onClick={(e) => e.preventDefault()}>
                                <Form.Control
                                    className={classNames({ error: false })}
                                    type='password'
                                    value={key}
                                    autoComplete='current-password'
                                    onChange={(e) => setKey(e.target.value)}
                                />
                                <span className='copy' onClick={() => handleCopy(key)}>
                                    <i className={`fas fa-copy`} />
                                </span>
                            </Label>
                            <div className='buttons-wrapper'>
                                <Button
                                    variant='primary'
                                    disabled={!uuid || !key}
                                    className='m-l-10'
                                    onClick={handleSubmit}
                                >
                                    {_(t.initialize)}
                                </Button>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    )
}

export default PreSharedKeySetup
