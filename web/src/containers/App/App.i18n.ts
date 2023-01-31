import { defineMessages } from '@formatjs/intl'

export const messages = defineMessages({
    loading: {
        id: 'app.loading',
        defaultMessage: 'Loading',
    },
    authError: {
        id: 'app.authError',
        defaultMessage: 'Authorization server error',
    },
    mockPart1: {
        id: 'app.mockPark1',
        defaultMessage: 'You are currently logged in to',
    },
    mockPart2: {
        id: 'app.mockPark2',
        defaultMessage:
            ".  Being logged in allows you to onboard your device without having to re-enter authorization each time. If you would like to log out, please click the 'Logout' button.",
    },
})
