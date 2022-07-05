export const devicesStatuses = {
  ONLINE: 'ONLINE',
  OFFLINE: 'OFFLINE',
  REGISTERED: 'REGISTERED',
  UNREGISTERED: 'UNREGISTERED',
}

export const devicesOwnerships = {
  OWNED: 'OWNED',
  UNOWNED: 'UNOWNED',
}

export const DEVICE_TYPE_OIC_WK_D = 'oic.wk.d'

export const devicesApiEndpoints = {
  DEVICES: '/api/v1/devices',
  DEVICES_RESOURCES_SUFFIX: 'resource-links',
  DEVICES_WS: '/api/v1/ws/devices',
}

export const RESOURCES_DEFAULT_PAGE_SIZE = 5

export const DEVICES_DEFAULT_PAGE_SIZE = 10

export const RESOURCE_TREE_DEPTH_SIZE = 15 // px

export const errorCodes = {
  DEADLINE_EXCEEDED: 'DeadlineExceeded',
  INVALID_ARGUMENT: 'InvalidArgument',
}

export const resourceModalTypes = {
  UPDATE_RESOURCE: 'update',
  CREATE_RESOURCE: 'create',
}

export const resourceEventTypes = {
  ADDED: 'added',
  REMOVED: 'removed',
}

export const knownInterfaces = {
  OIC_IF_A: 'oic.if.a',
  OIC_IF_BASELINE: 'oic.if.baseline',
  OIC_IF_CREATE: 'oic.if.create',
}

export const knownResourceTypes = {
  OIC_WK_CON: 'oic.wk.con', // contains device name
  X_PLGD_DPS_CONF: 'x.plgd.dps.conf',
}

export const shadowSynchronizationStates = {
  UNSET: 'UNSET',
  ENABLED: 'ENABLED',
  DISABLED: 'DISABLED',
}

export const defaultNewResource = {
  rt: [],
  if: [knownInterfaces.OIC_IF_A, knownInterfaces.OIC_IF_BASELINE],
  rep: {},
  p: {
    bm: 3,
  },
}

export const commandTimeoutUnits = {
  INFINITE: 'infinite',
  MS: 'ms',
  S: 's',
  M: 'min',
  H: 'h',
  NS: 'ns',
}

export const MINIMAL_TTL_VALUE_MS = 100

export const NO_DEVICE_NAME = '<no-name>'

// Maximum amount of deviceIds filled into one delete request.
// (if ther is more deviceIds then the provided number, it creates more chunks of a maximum of this number)
export const DEVICE_DELETE_CHUNK_SIZE = 50

// Websocket keys
export const DEVICES_WS_KEY = 'devices'
export const STATUS_WS_KEY = 'status'
export const RESOURCE_WS_KEY = 'resource'
export const DEVICES_STATUS_WS_KEY = `${DEVICES_WS_KEY}.${STATUS_WS_KEY}`
export const DEVICES_RESOURCE_REGISTRATION_WS_KEY = `${DEVICES_WS_KEY}.${RESOURCE_WS_KEY}.registration`
export const DEVICES_RESOURCE_UPDATE_WS_KEY = `${DEVICES_WS_KEY}.${RESOURCE_WS_KEY}.update`

// Emitter Event keys
export const DEVICES_REGISTERED_UNREGISTERED_COUNT_EVENT_KEY =
  'devices-registered-unregistered-count'

// Constant used in the DEVICES_REGISTERED_UNREGISTERED_COUNT_EVENT_KEY when reseting the counter
export const RESET_COUNTER = 'reset-counter'

export const DISCOVERY_DEFAULT_TIMEOUT_RAW = 2000
export const TIMEOUT_UNIT_PRECISION = 1000000

export const DISCOVERY_DEFAULT_TIMEOUT =
  DISCOVERY_DEFAULT_TIMEOUT_RAW * TIMEOUT_UNIT_PRECISION

export const RESOURCE_DEFAULT_TTL_RAW = 1000

export const RESOURCE_DEFAULT_TTL =
  RESOURCE_DEFAULT_TTL_RAW * TIMEOUT_UNIT_PRECISION
