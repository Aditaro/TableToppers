describe('LoginComponent', () => {
  beforeEach(() => {
    cy.visit('/login'); // Adjust the URL if necessary
  });

  it('should display the login form', () => {
    cy.get('form').should('be.visible');
    cy.get('input[name="email"]').should('be.visible');
    cy.get('input[name="password"]').should('be.visible');
    cy.get('button[type="submit"]').should('be.visible');
  });

  it('should show an error message for invalid login', () => {
    cy.get('input[name="email"]').type('test@example.com');
    cy.get('input[name="password"]').type('wrongpassword');
    cy.get('button[type="submit"]').click();

    cy.on('window:alert', (text) => {
      expect(text).to.contains('Login failed');
    });
  });

  it('should navigate to the homepage on successful login', () => {
    cy.get('input[name="email"]').type('test@example.com');
    cy.get('input[name="password"]').type('correctpassword');
    cy.get('button[type="submit"]').click();

    cy.url().should('eq', Cypress.config().baseUrl + '/');
  });
});
