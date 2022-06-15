import { Switch, Route } from 'react-router-dom'
import { DevicesListPage, DevicesDetailsPage } from '@/containers/devices'
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
