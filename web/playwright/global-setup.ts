import { chromium, expect, FullConfig } from '@playwright/test'

async function globalSetup(config: FullConfig) {
    const httpGatewayAddress = process.env.REACT_APP_HTTP_GATEWAY_ADDRESS || window.location.origin
    const browser = await chromium.launch()
    const page = await browser.newPage()

    const response = await page.request.get(`${httpGatewayAddress}/.well-known/configuration`)
    await expect(response).toBeOK()

    const data = await response.json()

    process.env.WELL_KNOWN_CONFIG = JSON.stringify(data)

    await page.goto('http://localhost:3000/')

    if (data?.deviceAuthenticationMode === 'X509') {
        await page.locator('#email').fill(process.env.REACT_APP_TEST_LOGIN_USERNAME || '')
        await page.locator('#password').fill(process.env.REACT_APP_TEST_LOGIN_PASSWORD || '')

        await page.getByRole('button', { name: 'Sign In' }).click()
    }

    await expect(page).toHaveTitle(/Devices | plgd Dashboard/)

    await page.context().storageState({ path: 'storageState.json' })
    await browser.close()
}

export default globalSetup
