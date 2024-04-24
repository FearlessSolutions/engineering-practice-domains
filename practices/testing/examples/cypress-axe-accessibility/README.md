
# Cypress + axe-core Accessibility Testing Demo

This project demonstrates how to integrate `axe-core` with Cypress for automated accessibility testing. It includes basic tests to showcase how you can leverage these tools to maintain accessibility standards in your web applications.

## Prerequisites

Before you begin, make sure you have the following installed:
- Node.js (v12 or higher)
- npm (v6 or higher)

## Installation

Follow these steps to set up the project on your local machine:

1. Clone this repository

2. Navigate to the project directory:
   ```
   cd cypress-axe-accessibility
   ```
3. Install dependencies:
   ```
   npm install
   ```

## Setting Up Cypress and axe-core

This project uses Cypress for testing along with the `cypress-axe` plugin to integrate `axe-core` for accessibility checks.

1. **Cypress Installation**: Cypress is installed as part of the npm dependencies.
2. **cypress-axe Installation**: The `cypress-axe` plugin is also included in the npm dependencies. This plugin adds specific `axe-core` commands to Cypress.

## Writing Tests

Tests are located in the `cypress/e2e` folder. Here's a quick overview of how an accessibility test is structured using `axe-core`:

```javascript
describe('Accessibility checks', () => {
  beforeEach(() => {
    cy.visit('url-of-your-choice');
    cy.injectAxe();  // Injects the axe-core library
  });

  it('Has no detectable accessibility violations on load', () => {
    cy.checkA11y(); // Checks the entire page for accessibility issues
  });
});
```

## Running Tests

To run your tests, use the following command:

```
npx cypress open
```

This command opens the Cypress Test Runner, which provides a visual interface for running the tests.

## Command Line / Continuous Integration

To run your tests via command line or for continuous integration pipelines, use the following command:

```
npx cypress run
```

This command opens the Cypress command line runner, which provides a non-visual, command line based interface for running the tests.

## Original Installation and Configuration Guide

This repo was installed with the most basic functions from `Cypress` and `axe-core`.  Minimal configurations were made to accomodate the functions here.  The following are some of the most important changes that would need to be made if you want to add these changes in to your own repo that already exists.

### Importing cypress-axe

After including `cypress-axe` you must be sure to also import the library in the `cypress/support/e2e.js` for use in the end to end tests.

Include the following import: 
```
import 'cypress-axe'
```

### Cypress log and table tasks
In the `cypress.config.js` file, two functions need to be added to the `setupNodeEvents` listener.  A `log()` and `table()` tasks are created for use in the custom `terminalLog` function below. Add the following `on('task')` events to the `cypress.config.js`.

```
on('task', {
        log(message) {
          console.log(message)

          return null
        },
        table(message) {
          console.table(message)

          return null
        }
      })
```

### Custom terminalLog function
A custom logging function is created for use in one of the example tests to demonstrate terminal loging and tabular logging of accessibilty errors.  This custom function can be found in the `cypress/support/functions.js` file.

This file is also imported in the `cypress/e2e/accesibility.cy.js` spec.

## Acknowledgments

- [Cypress.io](https://www.cypress.io/)
- [axe-core](https://github.com/dequelabs/axe-core)
