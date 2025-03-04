/// <reference types="cypress" />
/// <reference types="mocha" />

describe('Navigation Bar', () => {
  beforeEach(() => {
    cy.visit('/');
  });

  it('should display all main navigation items', () => {
    cy.get('#nav ul').first().within(() => {
      cy.contains('Dropdown').should('be.visible');
      cy.contains('Left Sidebar').should('be.visible');
      cy.contains('Restaurants').should('be.visible');
      cy.contains('Login').should('be.visible');
      cy.contains('Register').should('be.visible');
    });
  });

  it('should show dropdown menu on hover', () => {
    cy.contains('Dropdown').parent().find('ul').invoke('show');
    
    cy.contains('Lorem ipsum dolor').should('be.visible');
    cy.contains('Magna phasellus').should('be.visible');
    cy.contains('Phasellus consequat').should('be.visible');
    cy.contains('Veroeros feugiat').should('be.visible');
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
  });

  it('should show nested dropdown menu on hover', () => {
    // Show main dropdown
    cy.contains('Dropdown').parent().find('ul').invoke('show');
    // Show nested dropdown
    cy.contains('Phasellus consequat').parent().find('ul').invoke('show');
    
    cy.contains('Lorem ipsum dolor').should('be.visible');
    cy.contains('Phasellus consequat').should('be.visible');
    cy.contains('Magna phasellus').should('be.visible');
    cy.contains('Etiam dolore nisl').should('be.visible');
  });

  it('should maintain mobile responsiveness', () => {
    cy.viewport('iphone-6', 'portrait');
    
    cy.get('#navToggle').should('exist');
    cy.get('#navToggle > .toggle').first().click({ force: true });
    
    cy.get('#navPanel').should('exist').within(() => {
      // Check each nav item with scrolling
      ['Dropdown', 'Left Sidebar', 'Restaurants', 'Login', 'Register'].forEach(text => {
        cy.contains(text)
          .scrollIntoView()
          .should('be.visible');
      });
    });
  });
});