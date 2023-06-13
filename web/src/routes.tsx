import { Switch, Route, matchPath } from 'react-router-dom'
import { useIntl } from 'react-intl'

import NotFoundPage from '@shared-ui/components/Templates/NotFoundPage'
import { MenuItem } from '@shared-ui/components/Layout/LeftPanel/LeftPanel.types'

import DevicesListPage from '@/containers/Devices/List/DevicesListPage'
import DevicesDetailsPage from '@/containers/Devices/Detail/DevicesDetailsPage'
import MockApp from '@/containers/MockApp/MockApp'
import { messages as t } from './containers/App/App.i18n'
import { IconDevices } from '@shared-ui/components/Atomic'

export const menu = [
    {
        title: 'Main menu',
        items: [
            {
                icon: <IconDevices />,
                id: '1',
                title: 'Devices',
                link: '/',
                paths: ['/', '/devices/:id', '/devices/:id/:href*'],
            },
        ],
    },
]

export const mather = (location: string, item: MenuItem) =>
    matchPath(location, {
        path: item.paths,
        exact: false,
        strict: false,
    })

export const Routes = () => {
    const { formatMessage: _ } = useIntl()
    return (
        <Switch>
            <Route exact component={DevicesListPage} path='/' />
            <Route component={DevicesDetailsPage} path={['/devices/:id', '/devices/:id/:href*']} />
            <Route component={MockApp} path={['/devices', '/devices-code-redirect']} />
            <Route path='*'>
                <NotFoundPage message={_(t.notFoundPageDefaultMessage)} title={_(t.pageTitle)} />
            </Route>
        </Switch>
    )
}
