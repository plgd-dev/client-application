export type Props = {
    deviceId: string
    disabled: boolean
    endpointInformations?: { endpoint: string }[]
    href: string
    interfaces?: string[]
    isOwned: boolean
    onCreate: (href: string) => void
    onDelete: (href: string) => void
    onUpdate: ({ deviceId, href }: { deviceId: string; href: string }) => void
}

export const defaultProps = {
    interfaces: [],
}
