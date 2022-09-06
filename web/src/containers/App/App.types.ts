export type BuildInformationType = {
  version: string
  buildDate: string
  commitHash: string
}

export type Props = {
  config: {
    httpGatewayAddress: string
    openTelemetry: boolean
  }
}
