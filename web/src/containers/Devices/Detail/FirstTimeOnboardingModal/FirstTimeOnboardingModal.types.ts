export type Props = {
    onClose?: () => void
    onSubmit?: () => void
    show?: boolean
}

export const defaultProps = {
    show: false,
}
