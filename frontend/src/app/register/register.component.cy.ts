import { HttpClientModule } from '@angular/common/http';
import { RouterModule, Router } from '@angular/router';
import { RegisterComponent } from './register.component';
import { environment } from '../../environments/environment';

describe('RegisterComponent', () => {
  beforeEach(() => {
    // Mount the component with necessary providers
    cy.mount(RegisterComponent, {
      imports: [HttpClientModule, RouterModule.forRoot([])],
    });
  });

  it('should display the registration form', () => {
    // Check if the form elements are displayed
    cy.get('h1').should('contain.text', 'Create Account');
    cy.get('input[type="email"]').should('exist');
    cy.get('input[type="password"]').should('exist');
    cy.get('button[type="submit"]').should('contain.text', 'Create Account');
  });

  it('should validate form inputs', () => {
    // Try submitting without entering data
    cy.get('button[type="submit"]').click();
    // Form should not submit with empty fields (HTML5 validation)
    cy.url().should('not.include', '/login');

    // Enter invalid email
    cy.get('input[type="email"]').type('invalid-email');
    cy.get('input[type="password"]').type('password123');
    cy.get('button[type="submit"]').click();
    // Form should not submit with invalid email (HTML5 validation)
    cy.url().should('not.include', '/login');

    // Clear fields and enter valid data
    cy.get('input[type="email"]').clear().type('test@example.com');
    cy.get('input[type="password"]').clear().type('password123');
  });

  it('should handle successful registration', () => {
    // Create a spy for router navigation
    const routerSpy = cy.spy().as('routerSpy');
    
    // Stub the Router.navigate method
    cy.window().then((win) => {
      cy.stub(Router.prototype, 'navigate').callsFake((commands) => {
        routerSpy(commands);
        return Promise.resolve(true);
      });
    });
    
    // Stub the HTTP POST request
    cy.intercept('POST', `${environment.apiBaseUrl}/register`, {
      statusCode: 200,
      body: { message: 'Registration successful' },
    }).as('registerRequest');

    // Fill in the form
    cy.get('input[type="email"]').type('test@example.com');
    cy.get('input[type="password"]').type('password123');
    
    // Submit the form
    cy.get('button[type="submit"]').click();
    
    // Wait for the request to complete
    cy.wait('@registerRequest').its('request.body').should('deep.equal', {
      email: 'test@example.com',
      password: 'password123',
    });
  });

  it('should handle registration failure', () => {
    // Stub the HTTP POST request to simulate an error
    cy.intercept('POST', `${environment.apiBaseUrl}/register`, {
      statusCode: 400,
      body: { error: 'Registration failed' },
    }).as('registerFailure');

    // Fill in the form
    cy.get('input[type="email"]').type('existing@example.com');
    cy.get('input[type="password"]').type('password123');
    
    // Submit the form
    cy.get('button[type="submit"]').click();
    
    // Wait for the request to complete
    cy.wait('@registerFailure');
  })
});