import { Switch, Route } from 'react-router-dom'
import DevicesListPage from '@/containers/Devices/List/DevicesListPage'
import DevicesDetailsPage from '@/containers/Devices/Detail/DevicesDetailsPage'
import NotFoundPage from '@/containers/NotFoundPage'

export const Routes = () => (
    <Switch>
        <Route exact path='/' component={DevicesListPage} />
        <Route path={['/devices/:id', '/devices/:id/:href*']} component={DevicesDetailsPage} />
        <Route path='*'>
            <NotFoundPage />
        </Route>
    </Switch>
)
