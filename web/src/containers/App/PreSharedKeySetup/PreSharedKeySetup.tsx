import './PreSharedKeySetup.scss'
import LogoPlgd from './LogoPlgd'
import Label from '../../../../shared-ui/src/components/new/Label'
import TextField from '../../../../shared-ui/src/components/new/TextField'
import classNames from 'classnames'
import { useIntl } from 'react-intl'
import { messages as t } from './PreSharedKeySetup.i18n'
import { useState } from 'react'
import Button from '../../../../shared-ui/src/components/new/Button'

const PreSharedKeySetup = () => {
    const { formatMessage: _ } = useIntl()
    const [username, setUsername] = useState<string>('')
    const [password, setPassword] = useState<string>('')
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
                    <h2>Welcome to Start</h2>
                    <div className='fromWrapper'>
                        <form action=''>
                            <Label title={_(t.username)} onClick={(e) => e.preventDefault()}>
                                <TextField
                                    className={classNames({ error: false })}
                                    value={username}
                                    onChange={(e) => setUsername(e.target.value)}
                                />
                            </Label>
                            <Label title={_(t.password)} onClick={(e) => e.preventDefault()}>
                                <TextField
                                    className={classNames({ error: false })}
                                    value={password}
                                    onChange={(e) => setPassword(e.target.value)}
                                />
                            </Label>
                            <div className='buttons-wrapper'>
                                <Button icon='fa-bacon' variant='secondary' onClick={console.log}>
                                    {_(t.generate)}
                                </Button>
                                <Button
                                    icon='fa-bacon'
                                    variant='primary'
                                    disabled={!username || !password}
                                    className='m-l-10'
                                    onClick={console.log}
                                >
                                    {_(t.save)}
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
