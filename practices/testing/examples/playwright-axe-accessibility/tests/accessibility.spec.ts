import { test, expect } from '@playwright/test'
import AxeBuilder from '@axe-core/playwright'

test.describe('A11y Testing', () => {
    test('Tests the whole page with the accessibility analyzer', async ({ page }) => {
        await page.goto('https://digital.gov');
        const accessibilityScanResults = await new AxeBuilder({ page }).analyze();
        expect(accessibilityScanResults.violations.length).toEqual(0);
        expect(accessibilityScanResults.violations).toEqual([]);
    })

    test('Tests a specific component for accessibility errors and expects to pass', async ({ page }) => {
        await page.goto('https://digital.gov');
        await page.getByRole('banner').waitFor();

        const accessibilityScanResults = await new AxeBuilder({ page })
            .include('[role="banner"]')
            .analyze()
        expect(accessibilityScanResults.violations.length).toEqual(0);
        expect(accessibilityScanResults.violations).toEqual([]);
    })

    test('Tests a page against WCAG2 AAA standard', async ({ page }) => {
        await page.goto('https://digital.gov');
        const accessibilityScanResults = await new AxeBuilder({ page })
            .withTags(['wcag2aaa'])
            .analyze();
        expect(accessibilityScanResults.violations.length).toEqual(0);
        expect(accessibilityScanResults.violations).toEqual([]);

    })

    test('Tests a page against WCAG2.2 AA standard', async ({ page }) => {
        await page.goto('https://digital.gov');
        const accessibilityScanResults = await new AxeBuilder({ page })
            .withTags(['wcag22aa'])
            .analyze();
        expect(accessibilityScanResults.violations.length).toEqual(0);
        expect(accessibilityScanResults.violations).toEqual([]);

    })

    test('Tests another page against WCAG2.2 AAA standard', async ({ page }) => {
        await page.goto('https://digital.gov/resources/how-test-websites-for-accessibility/');
        const accessibilityScanResults = await new AxeBuilder({ page })
            .withTags(['wcag2aaa'])
            .analyze();
        expect(accessibilityScanResults.violations.length).toEqual(0);
        expect(accessibilityScanResults.violations).toEqual([]);

    })
   
})
