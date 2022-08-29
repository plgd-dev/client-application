export type Props = {
  defaultTtlValue?: number
  defaultValue?: number
  disabled: boolean
  isDelete?: boolean
  onChange: (value: number) => void
  onTtlHasError: (callback: boolean) => void
  title?: string
  ttlHasError: boolean
}

export const defaultProps = {
  defaultTtlValue: 0,
  defaultValue: 0,
  isDelete: false,
}
