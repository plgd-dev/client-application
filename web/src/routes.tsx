import { Routes as RoutesGroup, Route, matchPath } from 'react-router-dom'
import { useIntl } from 'react-intl'

import NotFoundPage from '@shared-ui/components/Templates/NotFoundPage'
import { IconDevices } from '@shared-ui/components/Atomic'
import DevicesListPage from '@shared-ui/app/clientApp/Devices/List/DevicesListPage'
import DevicesDetailsPage from '@shared-ui/app/clientApp/Devices/Detail/DevicesDetailsPage'
import MockApp from '@shared-ui/app/clientApp/MockApp'

import { messages as t } from './containers/App/App.i18n'

export const menu = [
    {
        title: 'Main menu',
        items: [
            {
                icon: <IconDevices />,
                id: '1',
                title: 'Devices',
                link: '/',
                paths: ['/', '/devices/:id', '/devices/:id/resources', '/devices/:id/resources/:href'],
            },
        ],
    },
]

export const mather = (pathname: string, pattern: string) => matchPath(pattern, pathname)

export const Routes = () => {
    const { formatMessage: _ } = useIntl()
    return (
        <RoutesGroup>
            <Route element={<DevicesListPage defaultActiveTab={0} />} path='/' />
            <Route element={<DevicesListPage defaultActiveTab={1} />} path='/configuration' />
            <Route element={<DevicesDetailsPage defaultActiveTab={0} />} path='/devices/:id' />
            <Route element={<DevicesDetailsPage defaultActiveTab={1} />} path='/devices/:id/resources' />
            <Route element={<DevicesDetailsPage defaultActiveTab={1} />} path='/devices/:id/resources/*' />
            <Route element={<MockApp />} path='/devices' />
            <Route element={<MockApp />} path='/devices-code-redirect' />
            <Route
                element={<NotFoundPage message={_(t.notFoundPageDefaultMessage)} title={_(t.pageTitle)} />}
                path='*'
            />
        </RoutesGroup>
    )
}
