import { createContext } from 'react'
import { AppContextType } from '@/containers/App/AppContext.types'

const AppContext = createContext<AppContextType>({
  collapsed: false,
  httpGatewayAddress: '',
})

export default AppContext
