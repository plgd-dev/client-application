import { test, expect } from '@playwright/test'

test('login action', async ({ page }) => {
    await page.goto('http://localhost:3000')

    await page.locator('#email').fill('test.user.ocfcloud@gmail.com')
    await page.locator('#password').fill('@br@k@d@br@')

    await page.locator('text="Sign In"').click()

    // Expect a title "to contain" a substring.
    await expect(page).toHaveTitle(/Devices | plgd Dashboard/)
})
