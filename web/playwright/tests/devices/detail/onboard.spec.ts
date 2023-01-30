import { test, expect } from '@playwright/test'
import testId from '../../../../src/testId'

test('onboard device', async ({ page }) => {
    const { WELL_KNOWN_CONFIG, REACT_APP_TEST_LOGIN_USERNAME, REACT_APP_TEST_LOGIN_PASSWORD } = process.env
    const { onboardButton, offboardButton, firstTimeModalButton, onboardTitleStatus } = testId.devices.detail

    // login to try.plgd.dev
    await page.goto('https://try.plgd.cloud/')
    await page.locator('#email').fill(REACT_APP_TEST_LOGIN_USERNAME || '')
    await page.locator('#password').fill(REACT_APP_TEST_LOGIN_PASSWORD || '')
    await page.getByRole('button', { name: 'Sign In' }).click()

    // back to page
    await page.goto('/')

    await page.getByTestId('devsim-server00').click()

    // is onboarded -> exist offboard button
    if (await page.getByTestId(offboardButton).isVisible()) {
        await page.getByTestId(offboardButton).click()
    }

    await page.getByTestId(onboardButton).click()

    await page.getByTestId(firstTimeModalButton).click()

    await expect(page.getByTestId(onboardTitleStatus)).not.toContainText('uninitialized')
})
