import { defineMessages } from '@formatjs/intl'

export const messages = defineMessages({
    online: {
        id: 'devices.online',
        defaultMessage: 'online',
    },
    offline: {
        id: 'devices.offline',
        defaultMessage: 'offline',
    },
    ownDevice: {
        id: 'devices.ownDevice',
        defaultMessage: 'Own Device',
    },
    disOwnDevice: {
        id: 'devices.disOwnDevice',
        defaultMessage: 'Disown Device',
    },
    deviceOwned: {
        id: 'devices.deviceOwned',
        defaultMessage: 'Device owned',
    },
    deviceDisOwned: {
        id: 'devices.deviceDisOwned',
        defaultMessage: 'Device disowned',
    },
    deviceWasOwned: {
        id: 'devices.deviceWasOwned',
        defaultMessage: 'device {name} was successfully owned.',
    },
    deviceWasDisOwned: {
        id: 'devices.deviceWasDisOwned',
        defaultMessage: 'device {name} was successfully disowned.',
    },
    deviceOwnError: {
        id: 'devices.deviceOwned',
        defaultMessage: 'Device own error',
    },
    owned: {
        id: 'devices.owned',
        defaultMessage: 'Owned',
    },
    unowned: {
        id: 'devices.unowned',
        defaultMessage: 'Unowned',
    },
    unsupported: {
        id: 'devices.unsupported',
        defaultMessage: 'Unsupported',
    },
    device: {
        id: 'devices.device',
        defaultMessage: 'Device',
    },
    deviceByIp: {
        id: 'devices.deviceByIp',
        defaultMessage: 'Device by IP',
    },
    deviceIp: {
        id: 'devices.deviceId',
        defaultMessage: 'Device IP',
    },
    enterDeviceIp: {
        id: 'devices.enterDeviceIp',
        defaultMessage: 'Enter the device IP',
    },
    invalidIp: {
        id: 'devices.invalidIp',
        defaultMessage: 'Invalid device IP',
    },
    name: {
        id: 'devices.name',
        defaultMessage: 'Name',
    },
    types: {
        id: 'devices.types',
        defaultMessage: 'Types',
    },
    endpoints: {
        id: 'devices.endpoints',
        defaultMessage: 'Endpoints',
    },
    supportedTypes: {
        id: 'devices.supportedTypes',
        defaultMessage: 'Supported Types',
    },
    interfaces: {
        id: 'devices.interfaces',
        defaultMessage: 'Interfaces',
    },
    deviceInterfaces: {
        id: 'devices.deviceInterfaces',
        defaultMessage: 'Device Interfaces',
    },
    status: {
        id: 'devices.status',
        defaultMessage: 'Status',
    },
    ownershipStatus: {
        id: 'devices.ownershipStatus',
        defaultMessage: 'Ownership status',
    },
    onboardingStatus: {
        id: 'devices.onboardingStatus',
        defaultMessage: 'Onboarding status',
    },
    deviceNotFound: {
        id: 'devices.deviceNotFound',
        defaultMessage: 'device not found',
    },
    deviceNotFoundMessage: {
        id: 'devices.deviceNotFoundMessage',
        defaultMessage: 'device with ID "{id}" does not exist.',
    },
    deviceResourcesNotFound: {
        id: 'devices.deviceResourcesNotFound',
        defaultMessage: 'device resources not found',
    },
    deviceResourcesNotFoundMessage: {
        id: 'devices.deviceResourcesNotFoundMessage',
        defaultMessage: 'device resources for device with ID "{id}" does not exist.',
    },
    href: {
        id: 'devices.href',
        defaultMessage: 'Href',
    },
    resources: {
        id: 'devices.resources',
        defaultMessage: 'Resources',
    },
    resourceInterfaces: {
        id: 'devices.resourceInterfaces',
        defaultMessage: 'Resource Interfaces',
    },
    deviceId: {
        id: 'devices.deviceId',
        defaultMessage: 'Device ID',
    },
    update: {
        id: 'devices.update',
        defaultMessage: 'Update',
    },
    updating: {
        id: 'devices.updating',
        defaultMessage: 'Updating',
    },
    create: {
        id: 'devices.create',
        defaultMessage: 'Create',
    },
    creating: {
        id: 'devices.creating',
        defaultMessage: 'Creating',
    },
    details: {
        id: 'devices.details',
        defaultMessage: 'Details',
    },
    retrieve: {
        id: 'devices.retrieve',
        defaultMessage: 'Retrieve',
    },
    retrieving: {
        id: 'devices.retrieving',
        defaultMessage: 'Retrieving',
    },
    delete: {
        id: 'devices.delete',
        defaultMessage: 'Delete',
    },
    select: {
        id: 'devices.select',
        defaultMessage: 'Select',
    },
    view: {
        id: 'devices.view',
        defaultMessage: 'View',
    },
    flushCache: {
        id: 'devices.flushCache',
        defaultMessage: 'Flush Cache',
    },
    flushDevices: {
        id: 'devices.flushDevices',
        defaultMessage: 'Flush Devices',
    },
    flushDevicesMessage: {
        id: 'devices.flushDevicesMessage',
        defaultMessage: 'Flush Devices',
    },
    deleteDevices: {
        id: 'devices.deleteDevices',
        defaultMessage: 'Delete devices',
    },
    deleting: {
        id: 'devices.deleting',
        defaultMessage: 'Deleting',
    },
    actions: {
        id: 'devices.actions',
        defaultMessage: 'Actions',
    },
    deleteResourceMessageSubtitle: {
        id: 'devices.deleteResourceMessageSubtitle',
        defaultMessage: 'This action cannot be undone.',
    },
    deleteResourceMessage: {
        id: 'devices.deleteResourceMessage',
        defaultMessage: 'Are you sure you want to delete this Resource?',
    },
    deleteDeviceMessage: {
        id: 'devices.deleteDeviceMessage',
        defaultMessage: 'Are you sure you want to delete this device? This action cannot be undone.',
    },
    deleteDevicesMessage: {
        id: 'devices.deleteDevicesMessage',
        defaultMessage: 'Are you sure you want to delete these {count} devices? This action cannot be undone.',
    },
    deleteAllDeviceMessage: {
        id: 'devices.deleteDeviceMessage',
        defaultMessage: 'Are you sure you want to delete all devices? This action cannot be undone.',
    },
    resourceWasUpdated: {
        id: 'devices.resourceWasUpdated',
        defaultMessage: 'The resource was updated successfully.',
    },
    resourceWasUpdatedOffline: {
        id: 'devices.resourceWasUpdatedOffline',
        defaultMessage: 'The resource update was scheduled, changes will be applied once the device is online.',
    },
    resourceWasDeletedOffline: {
        id: 'devices.resourceWasDeletedOffline',
        defaultMessage: 'Deleting of the resource was scheduled, it will be deleted once the device is online.',
    },
    resourceWasCreated: {
        id: 'devices.resourceWasCreated',
        defaultMessage: 'The resource was created successfully.',
    },
    resourceWasCreatedOffline: {
        id: 'devices.resourceWasCreatedOffline',
        defaultMessage: 'The resource creation was scheduled, changes will be applied once the device is online.',
    },
    invalidArgument: {
        id: 'devices.invalidArgument',
        defaultMessage: 'There was an invalid argument in the JSON structure.',
    },
    resourceUpdateSuccess: {
        id: 'devices.resourceUpdateSuccess',
        defaultMessage: 'Resource update successful',
        description: 'Title of the toast message on resource update success.',
    },
    resourceUpdate: {
        id: 'devices.resourceUpdate',
        defaultMessage: 'Resource update',
        description: 'Title of the toast message on resource update expired.',
    },
    resourceCreate: {
        id: 'devices.resourceCreate',
        defaultMessage: 'Resource creation',
        description: 'Title of the toast message on resource creation expired.',
    },
    resourceDelete: {
        id: 'devices.resourceDelete',
        defaultMessage: 'Resource deletion',
        description: 'Title of the toast message on resource deletion expired.',
    },
    commandOnResourceExpired: {
        id: 'devices.commandOnResourceExpired',
        defaultMessage: 'command on resource {deviceId}{href} has expired.',
        description: 'Continuos message for command expiration, keep the first letter lowercase!',
    },
    resourceUpdateError: {
        id: 'devices.resourceUpdateError',
        defaultMessage: 'Failed to update a resource',
        description: 'Title of the toast message on resource update error.',
    },
    resourceCreateSuccess: {
        id: 'devices.resourceCreateSuccess',
        defaultMessage: 'Resource created successfully',
        description: 'Title of the toast message on create resource success.',
    },
    resourceCreateError: {
        id: 'devices.resourceCreateError',
        defaultMessage: 'Failed to create a resource',
        description: 'Title of the toast message on resource create error.',
    },
    resourceRetrieveError: {
        id: 'devices.resourceRetrieveError',
        defaultMessage: 'Failed to retrieve a resource',
        description: 'Title of the toast message on resource retrieve error.',
    },
    resourceDeleteSuccess: {
        id: 'devices.resourceDeleteSuccess',
        defaultMessage: 'Resource delete scheduled',
        description: 'Title of the toast message on delete resource schedule success.',
    },
    resourceWasDeleted: {
        id: 'devices.resourceWasDeleted',
        defaultMessage: 'The resource delete was scheduled, you will be notified when the resource was deleted.',
    },
    resourceDeleteError: {
        id: 'devices.resourceDeleteError',
        defaultMessage: 'Failed to delete a resource',
        description: 'Title of the toast message on resource delete error.',
    },
    shadowSynchronizationError: {
        id: 'devices.shadowSynchronizationError',
        defaultMessage: 'Failed to set shadow synchronization',
        description: 'Title of the toast message on shadow synchronization set error.',
    },
    shadowSynchronizationWasSetOffline: {
        id: 'devices.shadowSynchronizationWasSetOffline',
        defaultMessage: 'Shadow synchronization was scheduled, changes will be applied once the device is online.',
    },
    deviceWentOnline: {
        id: 'devices.deviceWentOnline',
        defaultMessage: 'Device "{name}" went online.',
    },
    deviceWentOffline: {
        id: 'devices.deviceWentOffline',
        defaultMessage: 'Device "{name}" went offline.',
    },
    deviceWasUnregistered: {
        id: 'devices.deviceWasUnregistered',
        defaultMessage: 'Device "{name}" was unregistered.',
    },
    devicestatusChange: {
        id: 'devices.devicestatusChange',
        defaultMessage: 'Device status change',
    },
    notifications: {
        id: 'devices.notifications',
        defaultMessage: 'Notifications',
    },
    refresh: {
        id: 'devices.refresh',
        defaultMessage: 'Refresh',
    },
    discovery: {
        id: 'devices.discovery',
        defaultMessage: 'Discovery',
    },
    newResource: {
        id: 'devices.newResource',
        defaultMessage: 'New Resource',
    },
    resourceDeleted: {
        id: 'devices.resourceDeleted',
        defaultMessage: 'Resource Deleted',
    },
    newResources: {
        id: 'devices.newResources',
        defaultMessage: 'New Resources',
    },
    resourcesDeleted: {
        id: 'devices.resourcesDeleted',
        defaultMessage: 'Resources Deleted',
    },
    resourceWithHrefWasDeleted: {
        id: 'devices.resourceWithHrefWasDeleted',
        defaultMessage: 'Resource {href} was deleted from device {deviceName} ({deviceId}).',
    },
    resourceAdded: {
        id: 'devices.resourceAdded',
        defaultMessage: 'New resource {href} was added to the device {deviceName} ({deviceId}).',
    },
    resourcesAdded: {
        id: 'devices.resourcesAdded',
        defaultMessage: '{count} new resources were added to the device {deviceName} ({deviceId}).',
    },
    resourcesWereDeleted: {
        id: 'devices.resourcesWereDeleted',
        defaultMessage: '{count} resources were deleted from device {deviceName} ({deviceId}).',
    },
    resourceUpdated: {
        id: 'devices.resourceUpdated',
        defaultMessage: 'Resource Updated',
    },
    resourceUpdatedDesc: {
        id: 'devices.resourceUpdatedDesc',
        defaultMessage: 'Resource {href} on a device called {deviceName} was updated.',
    },
    treeView: {
        id: 'devices.treeView',
        defaultMessage: 'Tree view',
    },
    shadowSynchronization: {
        id: 'devices.shadowSynchronization',
        defaultMessage: 'Shadow synchronization',
    },
    save: {
        id: 'devices.save',
        defaultMessage: 'Save',
    },
    saving: {
        id: 'devices.saving',
        defaultMessage: 'Saving',
    },
    enterdeviceName: {
        id: 'devices.enterdeviceName',
        defaultMessage: 'Enter device name',
    },
    deviceNameChangeFailed: {
        id: 'devices.deviceNameChangeFailed',
        defaultMessage: 'device name change failed',
    },
    enabled: {
        id: 'devices.enabled',
        defaultMessage: 'Enabled',
    },
    disabled: {
        id: 'devices.disabled',
        defaultMessage: 'Disabled',
    },
    commandTimeout: {
        id: 'devices.commandTimeout',
        defaultMessage: 'Command Timeout',
    },
    minimalValueIs: {
        id: 'devices.minimalValueIs',
        defaultMessage: 'Minimal value is {minimalValue}.',
    },
    minimalValueIs2: {
        id: 'devices.minimalValueIs',
        defaultMessage: 'Minimal value is',
    },
    devicesDeleted: {
        id: 'devices.devicesDeleted',
        defaultMessage: 'Devices deleted',
        description: 'Title of the toast message on devices deleted success.',
    },
    devicesDeletedMessage: {
        id: 'devices.devicesDeletedMessage',
        defaultMessage: 'Devices were successfully deleted.',
    },
    deviceDeleted: {
        id: 'devices.deviceDeleted',
        defaultMessage: 'device deleted',
        description: 'Title of the toast message on device deleted success.',
    },
    deviceWasDeleted: {
        id: 'devices.deviceWasDeleted',
        defaultMessage: 'device {name} was successfully deleted.',
    },
    devicesDeletionError: {
        id: 'devices.devicesDeletion',
        defaultMessage: 'Failed to delete selected devices.',
        description: 'Title of the toast message on devices deleted failed.',
    },
    deviceDeletionError: {
        id: 'devices.deviceDeletionError',
        defaultMessage: 'Failed to delete this device.',
        description: 'Title of the toast message on devices deleted failed.',
    },
    default: {
        id: 'devices.default',
        defaultMessage: 'Default',
    },
    cancel: {
        id: 'devices.cancel',
        defaultMessage: 'Cancel',
    },
    close: {
        id: 'devices.close',
        defaultMessage: 'Close',
    },
    enterDeviceId: {
        id: 'devices.enterDeviceId',
        defaultMessage: 'Enter the device ID',
    },
    getCode: {
        id: 'devices.getCode',
        defaultMessage: 'Get the Code',
    },
    addDevice: {
        id: 'devices.addDevice',
        defaultMessage: 'Add device',
    },
    back: {
        id: 'devices.back',
        defaultMessage: 'Back',
    },
    provisionNewDevice: {
        id: 'devices.provisionNewDevice',
        defaultMessage: 'Provision a new device',
    },
    findDeviceByIp: {
        id: 'devices.findDeviceByIp',
        defaultMessage: 'Find a device by IP',
    },
    deviceAuthCodeError: {
        id: 'devices.deviceAuthCodeError',
        defaultMessage: 'Device Authorization Code Error',
    },
    deviceAddByIpError: {
        id: 'devices.deviceAddByIpError',
        defaultMessage: 'Device Error',
    },
    deviceAddByIpSuccess: {
        id: 'devices.deviceAddByIpSuccess',
        defaultMessage: 'Device was added successfully',
    },
    authorizationCode: {
        id: 'devices.authorizationCode',
        defaultMessage: 'Authorization Code',
    },
    authorizationProvider: {
        id: 'devices.authorizationProvider',
        defaultMessage: 'Authorization Provider',
    },
    deviceEndpoint: {
        id: 'devices.deviceEndpoint',
        defaultMessage: 'Device Endpoint',
    },
    hubId: {
        id: 'devices.hubId',
        defaultMessage: 'Hub ID',
    },
    certificateAuthorities: {
        id: 'devices.certificateAuthorities',
        defaultMessage: 'Certificate Authorities',
    },
    changeTimeout: {
        id: 'devices.changeTimeout',
        defaultMessage: 'Change timeout',
    },
    changeDiscoveryTimeout: {
        id: 'devices.changeDiscoveryTimeout',
        defaultMessage: 'Change discovery timeout',
    },
    discoveryTimeout: {
        id: 'devices.discoveryTimeout',
        defaultMessage: 'Discovery timeout',
    },
    enterDeviceName: {
        id: 'devices.enterDeviceName',
        defaultMessage: 'Enter device name',
    },
    setDpsEndpoint: {
        id: 'devices.setDpsEndpoint',
        defaultMessage: 'Set DPS Endpoint',
    },
    provisionNewDeviceTitle: {
        id: 'devices.provisionNewDeviceTitle',
        defaultMessage: 'Provision a new device using Device Provisioning Service',
    },
    deviceProvisioningServiceEndpoint: {
        id: 'devices.deviceProvisioningServiceEndpoint',
        defaultMessage: 'Device Provisioning Service Endpoint',
    },
    dpsStatus: {
        id: 'devices.dpsStatus',
        defaultMessage: 'DPS status',
    },
    notAvailable: {
        id: 'devices.notAvailable',
        defaultMessage: 'n/a',
    },
    onboardDevice: {
        id: 'devices.onboardDevice',
        defaultMessage: 'Onboard device',
    },
    offboardDevice: {
        id: 'devices.offboardDevice',
        defaultMessage: 'Offboard device',
    },
    onboardIncompleteModalTitle: {
        id: 'devices.onboardIncompleteModalTitle',
        defaultMessage: 'Incomplete data',
    },
    onboardingFieldCertificateAuthority: {
        id: 'devices.onboardingFieldCertificateAuthority',
        defaultMessage: 'Certificate Authorities',
    },
    onboardingFieldDeviceEndpoint: {
        id: 'devices.onboardingFieldDeviceEndpoint',
        defaultMessage: 'Device Endpoint',
    },
    onboardingFieldClientId: {
        id: 'devices.onboardingFieldClientId',
        defaultMessage: 'ClientId',
    },
    onboardingFieldAuthorizationProvider: {
        id: 'devices.onboardingFieldAuthorizationProvider',
        defaultMessage: 'Authorization Provider',
    },
    onboardingFieldAuthorizationCode: {
        id: 'devices.onboardingFieldAuthorizationCode',
        defaultMessage: 'Authorization Code',
    },
    onboardingFieldDeviceId: {
        id: 'devices.onboardingFieldDeviceId',
        defaultMessage: 'Device ID',
    },
    onboardingFieldHubId: {
        id: 'devices.onboardingFieldHubId',
        defaultMessage: 'Hub ID',
    },
    changeOnboardingData: {
        id: 'devices.changeOnboardingData',
        defaultMessage: 'Change onboarding data',
    },
    edit: {
        id: 'devices.edit',
        defaultMessage: 'Edit',
    },
    ok: {
        id: 'devices.ok',
        defaultMessage: 'Ok',
    },
    on: {
        id: 'devices.on',
        defaultMessage: 'On',
    },
    off: {
        id: 'devices.off',
        defaultMessage: 'Off',
    },
    firstTimeTitle: {
        id: 'devices.firstTimeTitle',
        defaultMessage: 'Pop-up Authorization Required for Device Onboarding',
    },
    firstTimeDescription: {
        id: 'devices.firstTimeDescription',
        defaultMessage: 'Please allow pop-ups for this website in order to proceed with device onboarding.',
    },
    recentTasks: {
        id: 'pendingCommands.recentTasks',
        defaultMessage: 'Recent tasks',
    },
    search: {
        id: 'devices.search',
        defaultMessage: 'Search',
    },
    editName: {
        id: 'devices.editName',
        defaultMessage: 'Edit name',
    },
    deviceName: {
        id: 'devices.deviceName',
        defaultMessage: 'Device Name',
    },
    reset: {
        id: 'devices.reset',
        defaultMessage: 'Reset',
    },
    saveChange: {
        id: 'devices.saveChange',
        defaultMessage: 'Save change',
    },
    savingChanges: {
        id: 'devices.savingChanges',
        defaultMessage: 'Saving change',
    },
    deviceInformation: {
        id: 'devices.deviceInformation',
        defaultMessage: 'Device information',
    },
    id: {
        id: 'devices.id',
        defaultMessage: 'ID',
    },
    duration: {
        id: 'devices.duration',
        defaultMessage: 'duration',
    },
    placeholder: {
        id: 'devices.placeholder',
        defaultMessage: 'placeholder',
    },
    unit: {
        id: 'devices.unit',
        defaultMessage: 'ID',
    },
    pasteAll: {
        id: 'devices.pasteAll',
        defaultMessage: 'Paste All',
    },
})
