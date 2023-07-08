import { ChangeEvent, useState } from 'react'
import classNames from 'classnames'
import { useIntl } from 'react-intl'
import Form from 'react-bootstrap/Form'

import Label from '@shared-ui/components/Atomic/Label'
import TextField from '@shared-ui/components/Atomic/TextField'
import { copyToClipboard } from '@shared-ui/common/utils'
import Notification from '@shared-ui/components/Atomic/Notification/Toast'
import Button from '@shared-ui/components/Atomic/Button'
import { IconCopy, IconHidePassword, IconShowPassword } from '@shared-ui/components/Atomic/Icon'

import './PreSharedKeySetup.scss'
import LogoPlgd from './LogoPlgd'
import { messages as t } from './PreSharedKeySetup.i18n'
import { initializedByPreShared } from '@/containers/App/AppRest'
import { Props } from './PreSharedKeySetup.types'

const validate = require('validate.js')

type ValidationResult = {
    subjectId?: string[]
    ket?: string[]
}

const PreSharedKeySetup = (props: Props) => {
    const { setInitialize } = props
    const { formatMessage: _ } = useIntl()
    const [uuid, setUuid] = useState<string>('')
    const [key, setKey] = useState<string>('')
    const [passwordType, setPasswordType] = useState('password')
    const [validationResult, setValidationResult] = useState<undefined | ValidationResult>(undefined)

    const handleSubmit = () => {
        const result = validate(
            { subjectId: uuid, key },
            {
                subjectId: {
                    presence: true,
                },
                key: {
                    presence: true,
                },
            }
        )

        if (result === undefined) {
            initializedByPreShared(uuid, key)
                .then((r) => {
                    if (r.status === 200) {
                        setInitialize(true)
                    }
                })
                .catch((e) => {
                    Notification.error({
                        title: _(t.error),
                        message: e.response.data.message,
                    })
                })
        } else {
            setValidationResult(result)
        }
    }

    const handleCopy = (data: string) => {
        copyToClipboard(data)
        Notification.success({
            title: _(t.done),
            message: _(t.copied),
        })
    }

    const handleUuidChange = (e: ChangeEvent<HTMLInputElement>) => {
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
                    <h2>{_(t.headline)}</h2>
                    <div className='fromWrapper'>
                        <form action=''>
                            <Label onClick={(e) => e.preventDefault()} title={_(t.subjectId)}>
                                <TextField
                                    autoComplete='subjectId'
                                    className={classNames({ error: validationResult?.subjectId })}
                                    name='subjectId'
                                    onChange={handleUuidChange}
                                    value={uuid}
                                />
                                <span className='copy' onClick={() => handleCopy(uuid)}>
                                    <IconCopy />
                                </span>
                            </Label>
                            {validationResult?.subjectId && (
                                <div className='m-b-10 error-message'>{_(t.subjectIdError)}</div>
                            )}
                            <Label onClick={(e) => e.preventDefault()} title={_(t.key)}>
                                <Form.Control
                                    autoComplete='current-password'
                                    className={classNames({ error: false })}
                                    onChange={(e) => setKey(e.target.value)}
                                    type={passwordType}
                                    value={key}
                                />
                                <span className='copy' onClick={() => handleCopy(key)}>
                                    <IconCopy />
                                </span>
                                <span
                                    className='show-password'
                                    onClick={() => setPasswordType(passwordType === 'password' ? 'text' : 'password')}
                                >
                                    {passwordType === 'password' ? <IconShowPassword /> : <IconHidePassword />}
                                </span>
                            </Label>
                            <div className='buttons-wrapper'>
                                <Button
                                    className='m-l-10'
                                    disabled={!uuid || !key}
                                    onClick={handleSubmit}
                                    variant='primary'
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
