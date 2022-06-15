// This state holds information about devices.

import { createSlice } from '@reduxjs/toolkit'
import findIndex from 'lodash/findIndex'
import { devicesOwnerships } from '@/containers/devices/constants'
const { OWNED, UNOWNED } = devicesOwnerships

const initialState = {
  activeNotifications: [],
  devicesList: [],
}

const { reducer, actions } = createSlice({
  name: 'devices',
  initialState,
  reducers: {
    addActiveNotification(state, { payload }) {
      state.activeNotifications.push(payload)
    },
    deleteActiveNotification(state, { payload }) {
      state.activeNotifications.splice(
        state.activeNotifications.findIndex(
          notification => notification === payload
        ),
        1
      )
    },
    toggleActiveNotification(state, { payload }) {
      if (!state.activeNotifications.includes(payload)) {
        state.activeNotifications.push(payload)
      } else {
        state.activeNotifications.splice(
          state.activeNotifications.findIndex(
            notification => notification === payload
          ),
          1
        )
      }
    },
    setDevices(state, { payload }) {
      state.devicesList = payload
    },
    addDevice(state, { payload }) {
      state.devicesList.push(payload)
    },
    flushDevices(state) {
      state.devicesList = []
    },
    toggleOwnDevice(state, { payload }) {
      const index = findIndex(
        state.devicesList,
        device => device.id === payload.deviceId
      )

      if (index > -1) {
        state.devicesList[index].ownershipStatus = payload.ownState
          ? OWNED
          : UNOWNED
      }
    },
  },
})

// Actions
export const {
  addActiveNotification,
  deleteActiveNotification,
  toggleActiveNotification,
  setDevices,
  addDevice,
  flushDevices,
  toggleOwnDevice,
} = actions

// Reducer
export default reducer

// Selectors
export const selectActiveNotifications = state =>
  state.devices.activeNotifications

export const isNotificationActive = key => state =>
  state.devices.activeNotifications?.includes?.(key) || false

export const getDevices = state => state.devices.devicesList
