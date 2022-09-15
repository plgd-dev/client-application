import { fetchApi, security } from '@shared-ui/common/services'
import { SecurityConfig } from '@/containers/App/App.types'

const getConfig = () => security.getGeneralConfig() as SecurityConfig

export const getAppWellKnownConfiguration = () =>
  fetchApi(`${getConfig().httpGatewayAddress}/.well-known/configuration`, {
    useToken: false,
  })

export const getOpenIdConfiguration = (authority: string) =>
  fetchApi(`${authority}/.well-known/openid-configuration`, {
    useToken: false,
  })

export const getJwksData = (url: string) => fetchApi(url, { useToken: false })

export const initializeJwksData = (data: any) =>
  fetchApi(`${getConfig().httpGatewayAddress}/api/v1/initialize`, {
    method: 'POST',
    body: {
      jwks: data,
    },
  })

export const signIdentityCsr = (
  certificateAuthorityAddress: string,
  certificateSigningRequest: string
) =>
  fetchApi(`https://${certificateAuthorityAddress}/api/v1/sign/identity-csr`, {
    method: 'POST',
    body: {
      certificateSigningRequest,
    },
  })

export const initializeFinal = (state: string, certificate: string) =>
  fetchApi(`${getConfig().httpGatewayAddress}/api/v1/initialize/${state}`, {
    method: 'POST',
    body: {
      certificate,
    },
  })

export const reset = () =>
  fetchApi(`${getConfig().httpGatewayAddress}/api/v1/reset`, {
    method: 'POST',
    body: {},
  })