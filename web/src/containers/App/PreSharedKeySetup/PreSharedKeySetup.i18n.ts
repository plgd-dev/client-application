import { defineMessages } from '@formatjs/intl'

export const messages = defineMessages({
    reminder: {
        id: 'preSharedKeySetup.reminder',
        defaultMessage: 'Reminder',
    },
    reminderDescription: {
        id: 'preSharedKeySetup.reminderDescription',
        defaultMessage:
            "Please copy your pre-shared key and store it securely. If you reinstall your application and you won't initialize it with the same Pre-shared Key, devices you previously owned won't be accessible and the factory reset on these devices will need to be executed.",
    },
    headline: {
        id: 'preSharedKeySetup.headline',
        defaultMessage: 'Pre shared key setup',
    },
    subjectId: {
        id: 'preSharedKeySetup.subjectId',
        defaultMessage: 'Subject ID',
    },
    subjectIdError: {
        id: 'preSharedKeySetup.subjectIdError',
        defaultMessage: 'Must be in UUD format: XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX',
    },
    key: {
        id: 'preSharedKeySetup.key',
        defaultMessage: 'Key',
    },
    initialize: {
        id: 'preSharedKeySetup.initialize',
        defaultMessage: 'Initialize',
    },
    generate: {
        id: 'preSharedKeySetup.generate',
        defaultMessage: 'Generate',
    },
    done: {
        id: 'preSharedKeySetup.done',
        defaultMessage: 'Done',
    },
    copied: {
        id: 'preSharedKeySetup.copied',
        defaultMessage: 'Value is in your clipboard.',
    },
    error: {
        id: 'preSharedKeySetup.error',
        defaultMessage: 'Error',
    },
})
