describe('Business Portal Page', () => {
  beforeEach(() => {
    cy.visit('/business-portal');
  });

  it('should display the business portal page correctly', () => {
    cy.get('.business-portal').should('exist');
    cy.get('.hero h1').should('contain', 'Transform Your Restaurant Management');
  });

  it('should display and play the demo GIF', () => {
    cy.get('.demo-image img')
      .should('be.visible')
      .and('have.attr', 'src', 'assets/css/images/table-managment-demo.gif')
      .and('have.attr', 'alt', 'Table Management Demo');
  });

  it('should open demo form modal when clicking Schedule Demo', () => {
    cy.get('.secondary-button').contains('Schedule Demo').click();
    cy.get('.modal-overlay').should('be.visible');
    cy.get('.modal-content').should('be.visible');
  });

  it('should open demo form modal when clicking Contact Sales', () => {
    cy.get('.secondary-button').contains('Contact Sales').click();
    cy.get('.modal-overlay').should('be.visible');
    cy.get('.modal-content').should('be.visible');
  });

  it('should be able to fill out and submit the demo form', () => {
    cy.get('.secondary-button').contains('Schedule Demo').click();
    
    // Fill out the form
    cy.get('#name').type('Test User');
    cy.get('#email').type('test@example.com');
    cy.get('#phone').type('1234567890');
    cy.get('#company').type('Test Company');
    cy.get('#message').type('This is a test message');
    
    // Submit the form
    cy.get('.submit-button').should('not.be.disabled').click();
    
    // Check for success message
    cy.get('.success-message').should('be.visible');
  });

  it('should navigate to subscription page when clicking Get Started Free', () => {
    cy.get('.cta-button').contains('Get Started Free').click();
    cy.url().should('include', '/subscription');
  });

  it('should navigate to subscription page when clicking Start Free Trial', () => {
    cy.get('.cta-button').contains('Start Free Trial').click();
    cy.url().should('include', '/subscription');
  });
});

describe('Home Page Navigation', () => {
  beforeEach(() => {
    cy.visit('/');
  });

  it('should navigate to login page when clicking More Info', () => {
    cy.get('.button.alt').contains('More info').click();
    cy.url().should('include', '/login');
    // Verify the returnUrl parameter is set correctly
    cy.url().should('include', 'returnUrl=%2Fcustomer-info');
  });

  it('should navigate to login page when clicking Reserve Now', () => {
    cy.get('.button').contains('Reserve Now').click();
    cy.url().should('include', '/login');
  });
}); 