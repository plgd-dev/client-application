export type BuildInformationType = {
    buildDate: string
    commitDate: string
    commitHash: string
    releaseUrl: string
    version: string
}

export type RemoteProvisioningDataType = {
    authorization: {
        audience: string
        authority: string
        clientId: string
        ownerClaim: string
        scopes: string[]
    }
    mode: string
    userAgent: {
        certificateAuthorityAddress: string
        csrChallengeStateExpiration: string
    }
}

export type WellKnownConfigType = {
    deviceAuthenticationMode: string
    isInitialized: boolean
    remoteProvisioning: RemoteProvisioningDataType
} & BuildInformationType

export type SecurityConfig = {
    httpGatewayAddress: string
}
