import { test, expect } from '@playwright/test'

test('login action', async ({ page }) => {
    await page.goto('http://localhost:3000')

    await page.locator('#email').fill(process.env.REACT_APP_TEST_LOGIN_USERNAME || '')
    await page.locator('#password').fill(process.env.REACT_APP_TEST_LOGIN_PASSWORD || '')

    await page.getByRole('button', { name: 'Sign In' }).click()

    // Expect a title "to contain" a substring.
    await expect(page).toHaveTitle(/Devices | plgd Dashboard/)
})
