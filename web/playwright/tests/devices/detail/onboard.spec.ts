import { test, expect } from '@playwright/test'
import testId from '../../../../src/testId'

test('onboard device', async ({ page }) => {
    const { REACT_APP_TEST_LOGIN_USERNAME, REACT_APP_TEST_LOGIN_PASSWORD } = process.env
    const { onboardButton, offboardButton, firstTimeModalButton, onboardTitleStatus, ownTitleStatus, ownButton } =
        testId.devices.detail

    // login to try.plgd.dev
    await page.goto('https://try.plgd.cloud/')

    // browser is not logged to try.plgd.dev
    if ((await page.title()) === 'Login | plgd.dev') {
        await page.locator('#email').fill(REACT_APP_TEST_LOGIN_USERNAME || '')
        await page.locator('#password').fill(REACT_APP_TEST_LOGIN_PASSWORD || '')
        await page.getByRole('button', { name: 'Sign In' }).click()
    }

    // back to page
    await page.goto('/')

    await page.getByTestId('devsim-server00').click()

    await page.waitForLoadState('networkidle')

    const onwButton = await page.locator(`[data-test-id="${ownButton}"]`)
    await page.waitForTimeout(3000)

    // device must be owned
    if (await onwButton.count()) {
        await onwButton.click()
    }

    await expect(page.getByTestId(ownTitleStatus)).toContainText('owned')

    // is onboarded -> exist offboard button
    const offboardButtonE = await page.locator(`[data-test-id="${offboardButton}"]`)
    await page.waitForTimeout(3000)

    if (await offboardButtonE.count()) {
        await offboardButtonE.click()
    }

    await expect(page.getByTestId(onboardTitleStatus)).toContainText('uninitialized')

    await page.getByTestId(onboardButton).click()

    await page.getByTestId(firstTimeModalButton).click()

    await expect(page.getByTestId(onboardTitleStatus)).not.toContainText('uninitialized')
})
