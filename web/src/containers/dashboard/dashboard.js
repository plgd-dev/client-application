import { useIntl } from 'react-intl'

import Layout from '@shared-ui/components/new/Layout'
import { messages as menuT } from '@shared-ui/components/new/Menu/Menu.i18n'

export const Dashboard = () => {
  const { formatMessage: _ } = useIntl()

  return (
    <Layout
      title={_(menuT.dashboard)}
      breadcrumbs={[
        {
          label: _(menuT.dashboard),
        },
      ]}
    >
      <div />
    </Layout>
  )
}
