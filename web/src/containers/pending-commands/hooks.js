import { useStreamApi } from '@shared-ui/common/hooks'
import { useAppConfig } from '@/containers/app'
import { useEmitter } from '@shared-ui/common/hooks'

import {
  pendingCommandsApiEndpoints,
  NEW_PENDING_COMMAND_WS_KEY,
  UPDATE_PENDING_COMMANDS_WS_KEY,
} from './constants'
import {
  convertPendingCommandsList,
  updatePendingCommandsDataStatus,
} from './utils'

export const usePendingCommandsList = deviceId => {
  const { httpGatewayAddress } = useAppConfig()

  // Fetch the data
  const { data, updateData, ...rest } = useStreamApi(
    `${httpGatewayAddress}${pendingCommandsApiEndpoints.PENDING_COMMANDS}${
      deviceId ? `?deviceIdFilter=${deviceId}` : ''
    }`
  )

  // Add a new pending command when a WS event is emitted
  useEmitter(NEW_PENDING_COMMAND_WS_KEY, newCommand => {
    updateData((data || []).concat(newCommand))
  })

  useEmitter(UPDATE_PENDING_COMMANDS_WS_KEY, updated => {
    updateData(updatePendingCommandsDataStatus(data, updated))
  })

  return {
    data: convertPendingCommandsList(data),
    updateData,
    ...rest,
  }
}
