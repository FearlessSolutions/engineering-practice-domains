
# Playwright + axe-core Accessibility Testing Demo

This project demonstrates how to integrate `axe-core` with Playwright for automated accessibility testing. It includes basic tests to showcase how you can leverage these tools to maintain accessibility standards in your web applications.

## Prerequisites

Before you begin, make sure you have the following installed:
- Node.js (v12 or higher)
- npm (v6 or higher)

## Installation

Follow these steps to set up the project on your local machine:

1. Clone this repository

2. Navigate to the project directory:
   ```
   cd playwright-axe-accessibility
   ```
3. Install dependencies:
   ```
   npm install
   ```

## Setting Up Playwright and axe-core

This project uses Playwright for testing along with the `@axe-core/playwright` module to integrate `axe-core` for accessibility checks.

1. **Playwrightg Installation**: Playwright is installed as part of the npm dependencies.
2. **axe-core Installation**: The `@axe-core/playwright` plugin is also included in the npm dependencies. This plugin adds specific `axe-core` functions to Playwright.

## Writing Tests

Tests are located in the `tests` folder. Here's a quick overview of how an accessibility test is structured using `axe-core`:

```javascript
    test('Has no detectable accessibility violations on load', async ({ page }) => {
        await page.goto('url-of-your-choice');
        // build the axe analyzer for the page and analyze
        const accessibilityScanResults = await new AxeBuilder({ page }).analyze();
        // Checks that there are no violations in the results
        expect(accessibilityScanResults.violations).toEqual([]);
    })
```

## Running Tests

To run your tests, use the following command:

```
npx playwright test accessibility.spec.ts
```

This command opens the Playwright Test Runner to execute tests and displays the test report summary.

## Original Installation and Configuration Guide

This repo was installed with the most basic functions from `Playwright` and `@axe-core/playwright`.  No additional changes were made to the configurations.

### Importing axe-core

After including `@axe-core/playwright` you must be sure to also import the library in each spec that you want to utilize `axe-core` functionality

Include the following import in your tests: 
```
import AxeBuilder from '@axe-core/playwright'

```

## Acknowledgments

- [Playwright](https://playwright.dev/)
- [axe-core](https://github.com/dequelabs/axe-core)
