/// <reference types="cypress" />
/// <reference types="mocha" />

describe('Navigation Bar', () => {
  beforeEach(() => {
    cy.visit('/');
  });

  it('should display all main navigation items', () => {
    cy.get('#nav ul').first().within(() => {
      cy.contains('Business Portal').should('be.visible');
      cy.contains('Restaurants').should('be.visible');
      cy.contains('Login').should('be.visible');
      cy.contains('Register').should('be.visible');
    });
  });

  it('should navigate to correct pages when clicking links', () => {
    cy.contains('Restaurants').click();
    cy.url().should('include', '/restaurants');
    
    cy.visit('/');
    cy.contains('Login').click();
    cy.url().should('include', '/login');
    
    cy.visit('/');
    cy.contains('Register').click();
    cy.url().should('include', '/register');

    cy.visit('/');
    cy.contains('Business Portal').click();
    cy.url().should('include', '/business-portal');
  });
});