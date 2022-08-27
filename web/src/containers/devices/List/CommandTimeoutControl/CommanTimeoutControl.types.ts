export type Props = {
    defaultTtlValue?: number | string
    defaultValue?: number| string
    disabled: boolean
    isDelete?: boolean
    onChange: (value: number | string) => void
    onTtlHasError: (callback: boolean) => void
    title?: string
    ttlHasError: boolean
}

export const defaultProps = {
    defaultTtlValue: 0,
    defaultValue: 0,
    isDelete: false,
}