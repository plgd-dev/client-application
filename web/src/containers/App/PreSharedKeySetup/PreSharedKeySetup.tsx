import './PreSharedKeySetup.scss'
import LogoPlgd from './LogoPlgd'
import Label from '../../../../shared-ui/src/components/new/Label'
import TextField from '../../../../shared-ui/src/components/new/TextField'
import classNames from 'classnames'
import { useIntl } from 'react-intl'
import { messages as t } from './PreSharedKeySetup.i18n'
import { ChangeEvent, useState } from 'react'
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
    const [passwordType, setPasswordType] = useState('password')
    const [uuidError, setUUidError] = useState(false)

    const checkIfValidUUID = (str: string) => {
        // Regular expression to check if string is a valid UUID
        const regexExp = /^[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}$/gi

        return regexExp.test(str)
    }

    const handleSubmit = () => {
        if (checkIfValidUUID(uuid)) {
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
        } else {
            setUUidError(true)
        }
    }

    const handleCopy = (data: string) => {
        copyToClipboard(data)
        showSuccessToast({
            title: _(t.done),
            message: _(t.copied),
        })
    }

    const handleUuidChange = (e: ChangeEvent<HTMLInputElement>) => {
        if (uuidError) {
            setUUidError(false)
        }
        setUuid(e.target.value)
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
                                    className={classNames({ error: uuidError })}
                                    value={uuid}
                                    name='subjectId'
                                    autoComplete='subjectId'
                                    onChange={handleUuidChange}
                                />
                                <span className='copy' onClick={() => handleCopy(uuid)}>
                                    <i className={`fas fa-copy`} />
                                </span>
                            </Label>
                            {uuidError && <div className='m-b-10 error-message'>{_(t.subjectIdError)}</div>}
                            <Label title={_(t.key)} onClick={(e) => e.preventDefault()}>
                                <Form.Control
                                    className={classNames({ error: false })}
                                    type={passwordType}
                                    value={key}
                                    autoComplete='current-password'
                                    onChange={(e) => setKey(e.target.value)}
                                />
                                <span className='copy' onClick={() => handleCopy(key)}>
                                    <i className={`fas fa-copy`} />
                                </span>
                                <span
                                    className='show-password'
                                    onClick={() => setPasswordType(passwordType === 'password' ? 'text' : 'password')}
                                >
                                    <i className={`fas ${passwordType === 'password' ? 'fa-eye' : 'fa-eye-slash'}`} />
                                </span>
                            </Label>
                            <div className='buttons-wrapper'>
                                <Button
                                    variant='primary'
                                    disabled={!uuid || !key || uuidError}
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
