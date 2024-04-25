import { terminalLog } from '../support/functions'

describe('A11y Testing', () => {

    it('Tests the whole page with the accessibility analyzer', () => {
        cy.visit('https://digital.gov')
        cy.injectAxe()
        cy.checkA11y()
    })

    it('Tests a specific component for accessibility errors and expects to pass', () => {
        cy.visit('https://digital.gov')
        cy.get('[role="banner"]')
        cy.injectAxe()
        cy.checkA11y('[role="banner"]')
    })

    it('Tests a page against WCAG2 AAA standard', () => {
        cy.visit('https://digital.gov')

        cy.injectAxe()
        cy.checkA11y(null, {
            runOnly: {
                type: 'tag',
                values: ['wcag2aaa']
            }
        })
    })

    it('Tests a page against WCAG2.2 AA standard', () => {
        cy.visit('https://digital.gov')

        cy.injectAxe()
        cy.checkA11y(null, {
            runOnly: {
                type: 'tag',
                values: ['wcag22aa']
            }
        })
    })

    it('Tests another page against WCAG2 AAA standard', () => {
        cy.visit('https://digital.gov/resources/how-test-websites-for-accessibility/')

        cy.injectAxe()
        cy.checkA11y(null, {
            runOnly: {
                type: 'tag',
                values: ['wcag2aaa']
            }
        })
    })

    // The includedImpacts option can be set to the impacts that will be checked for.
    // These are categorizations that Deque has set to group severity of violations.
    // Can be one of "minor", "moderate", "serious", or "critical"
    it('Tests the whole page and filters by impact level', () => {
        cy.visit('https://digital.gov')
        cy.injectAxe()
        cy.checkA11y(null, {
            includedImpacts: ['moderate']
        })
    })

    it('Outputs errors to terminal using a special display function ', () => {
        cy.visit('https://digital.gov/resources/how-test-websites-for-accessibility/')
        cy.injectAxe()
        cy.checkA11y(null, null, terminalLog)
    })
})
