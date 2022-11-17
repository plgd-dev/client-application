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

const PreSharedKeySetup = (props: Props) => {
    const { setInitialize } = props
    const { formatMessage: _ } = useIntl()
    const [uuid, setUuid] = useState<string>('')
    const [key, setKey] = useState<string>('')

    const handleSubmit = () => {
        initializedByPreShared(uuid, key).then((r) => {
            if (r.status === 200) {
                setInitialize(true)
            }
        })
    }

    return (
        <div className='preSharedKeySetupPage'>
            <div className='colLeft'>
                <div className='top'>
                    <h1>Start.</h1>
                    <p className='claim'>Lorem Ipsum claim text</p>
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
                            <Label title={_(t.uuid)} onClick={(e) => e.preventDefault()}>
                                <TextField
                                    className={classNames({ error: false })}
                                    value={uuid}
                                    onChange={(e) => setUuid(e.target.value)}
                                />
                            </Label>
                            <Label title={_(t.key)} onClick={(e) => e.preventDefault()}>
                                <Form.Control
                                    className={classNames({ error: false })}
                                    type='password'
                                    value={key}
                                    onChange={(e) => setKey(e.target.value)}
                                />
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
