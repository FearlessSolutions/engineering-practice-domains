import { terminalLog } from '../support/functions'

describe('A11y Testing', () => {
    beforeEach(() => {
        cy.visit('https://example.cypress.io/ ')
        cy.injectAxe()
    })

    // The includedImpacts option can be set to the impacts that will be checked for.
    // These are categorizations that Deque has set to group severity of violations.
    // Can be one of "minor", "moderate", "serious", or "critical"
    it('Tests the whole page against critical impact level', () => {
        cy.checkA11y(null, {
            includedImpacts: ['critical']
        })
    })

    it('Tests the whole page against serious impact level', () => {
        cy.checkA11y(null, {
            includedImpacts: ['serious']
        })
    })

    // Tags can be specified to only run against certain rule sets
    // More information about tags at: https://www.deque.com/axe/core-documentation/api-documentation/#axecore-tags

    it('Tests the whole page against WCAG 2.2 AAA standard rules', () => {
        cy.checkA11y(null, {
            runOnly: {
                type: 'tag',
                values: ['wcag22aaa']
            }
        })
    })

    it('Tests specifically the navbar against WCAG 2.2 AAA standards', () => {
        cy.checkA11y('#navbar', {
            runOnly: {
                type: 'tag',
                values: ['wcag22aaa']
            }
        })
    })

    it('Tests specifically the "Kitchen Sink" banner against WCAG 2.2 AA standards', () => {
        cy.checkA11y('.banner', {
            runOnly: {
                type: 'tag',
                values: ['wcag22aaa']
            }
        })
    })

    it('Outputs errors to terminal using a special display function ', () => {
        cy.checkA11y(null, null, terminalLog)
    })
})
