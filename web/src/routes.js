import { Switch, Route } from 'react-router-dom'
import DevicesListPage  from '@/containers/devices/List/DevicesListPage'
import DevicesDetailsPage from '@/containers/devices/Detail/DevicesDetailsPage'
import { NotFoundPage } from '@/containers/not-found-page'

export const Routes = () => (
  <Switch>
    <Route exact path="/" component={DevicesListPage} />
    <Route
      path={['/devices/:id', '/devices/:id/:href*']}
      component={DevicesDetailsPage}
    />
    <Route path="*">
      <NotFoundPage />
    </Route>
  </Switch>
)
