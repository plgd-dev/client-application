// This state holds information about devices.

import { createSlice } from '@reduxjs/toolkit'
import findIndex from 'lodash/findIndex'
import {
  devicesOwnerships,
  DISCOVERY_DEFAULT_TIMEOUT,
} from '@/containers/devices/constants'
import isEqual from 'lodash/isEqual'

const { OWNED } = devicesOwnerships

const initialState = {
  activeNotifications: [],
  devicesList: [],
  discoveryTimeout: DISCOVERY_DEFAULT_TIMEOUT,
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
    updateDevices(state, { payload }) {
      if (state.devicesList.length === 0) {
        state.devicesList = payload
      } else {
        payload.forEach(device => {
          const index = findIndex(state.devicesList, d => d.id === device.id)
          if (index > -1) {
            if (!isEqual(state.devicesList[index], device)) {
              state.devicesList[index] = device
            }
          } else {
            state.devicesList.push(device)
          }
        })
      }
    },
    addDevice(state, { payload }) {
      const index = findIndex(
        state.devicesList,
        device => device.id === payload.id
      )

      if (index > -1) {
        if (!isEqual(state.devicesList[index], payload)) {
          state.devicesList[index] = payload
        }
      } else {
        state.devicesList.push(payload)
      }
    },
    flushDevices(state) {
      state.devicesList = []
    },
    ownDevice(state, { payload }) {
      const index = findIndex(
        state.devicesList,
        device => device.id === payload
      )

      if (index > -1) {
        state.devicesList[index].ownershipStatus = OWNED
      }
    },
    disOwnDevice(state, { payload }) {
      state.devicesList.splice(
        state.devicesList.findIndex(device => device.id === payload),
        1
      )
    },
    setDiscoveryTimeout(state, { payload }) {
      state.discoveryTimeout = payload
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
  ownDevice,
  disOwnDevice,
  updateDevices,
  setDiscoveryTimeout,
} = actions

// Reducer
export default reducer

// Selectors
export const selectActiveNotifications = state =>
  state.devices.activeNotifications

export const isNotificationActive = key => state =>
  state.devices.activeNotifications?.includes?.(key) || false

export const getDevices = state => state.devices.devicesList

export const getDevicesDiscoveryTimeout = state =>
  state.devices.discoveryTimeout
