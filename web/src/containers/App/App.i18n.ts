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
    pageTitle: {
        id: 'not-found-page.pageTitle',
        defaultMessage: 'Page not found',
    },
    notFoundPageDefaultMessage: {
        id: 'not-found-page.notFoundPageDefaultMessage',
        defaultMessage: 'The page you are looking for does not exist.',
    },
    logOut: {
        id: 'app.logOut',
        defaultMessage: 'Log Out',
    },
    version: {
        id: 'app.version',
        defaultMessage: 'Version',
    },
    newUpdateIsAvailable: {
        id: 'app.newUpdateIsAvailable',
        defaultMessage: 'New update is available.',
    },
    clickHere: {
        id: 'app.clickHere',
        defaultMessage: 'Click here!',
    },
})
