import { test, expect } from '@playwright/test'

test('login action', async ({ page }) => {
    const httpGatewayAddress = process.env.REACT_APP_HTTP_GATEWAY_ADDRESS || window.location.origin

    const response = await page.request.get(`${httpGatewayAddress}/.well-known/configuration`)
    await expect(response).toBeOK()

    const data = await response.json()

    if (data.wellKnownConfig?.deviceAuthenticationMode === 'X509') {
        await page.locator('#email').fill(process.env.REACT_APP_TEST_LOGIN_USERNAME || '')
        await page.locator('#password').fill(process.env.REACT_APP_TEST_LOGIN_PASSWORD || '')

        await page.getByRole('button', { name: 'Sign In' }).click()
    } else {
        await page.goto('/')
    }

    // Expect a title "to contain" a substring.
    await expect(page).toHaveTitle(/Devices | plgd Dashboard/)
})
