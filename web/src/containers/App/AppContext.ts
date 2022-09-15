import { createContext } from 'react'
import { AppContextType } from '@/containers/App/AppContext.types'

const AppContext = createContext<AppContextType>({
  collapsed: false,
})

export default AppContext
