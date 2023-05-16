import { createContext } from 'react'
import { AppContextType } from '@/containers/App/AppContext.types'

const AppContext = createContext<AppContextType>({})

export default AppContext
