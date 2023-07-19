import { defineMessages } from '@formatjs/intl'

export const messages = defineMessages({
    headline: {
        id: 'app.headline',
        defaultMessage: 'Application is initialized by a different user.',
    },
    description: {
        id: 'app.description',
        defaultMessage:
            'Application Initialization Restricted. Please ensure the current user logs out before proceeding. Only after the different user has logged out, will you be able to utilize the application.',
    },
    logout: {
        id: 'app.logout',
        defaultMessage: 'Log Out',
    },
})
