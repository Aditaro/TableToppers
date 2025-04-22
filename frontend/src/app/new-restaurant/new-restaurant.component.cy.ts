import { NewRestaurantComponent } from './new-restaurant.component'
import { FormBuilder, ReactiveFormsModule } from '@angular/forms'
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog'
import { RestaurantService } from '../services/restaurant.service'
import { CommonModule } from '@angular/common'
import { FormsModule } from '@angular/forms'
import { MatCardModule } from '@angular/material/card'
import { MatFormFieldModule } from '@angular/material/form-field'
import { MatInputModule } from '@angular/material/input'
import { MatButtonModule } from '@angular/material/button'
import { BrowserAnimationsModule } from '@angular/platform-browser/animations'
import { of } from 'rxjs'

describe('NewRestaurantComponent', () => {
  // Mock data for testing
  const mockRestaurant = {
    id: 'rest123',
    name: 'Test Restaurant',
    location: 'Test Location',
    description: 'Test Description',
    phone: '555-123-4567',
    openingHours: '11:00 - 21:00',
    status: 'open',
    img: 'test-image.jpg',
    specialAvailability: []
  };

  // Define mock objects without stubs initially
  let mockDialogRef;
  let mockRestaurantService;

  beforeEach(() => {
    // Initialize mock objects with stubs inside the test context
    mockDialogRef = {
      close: cy.stub().as('dialogClose')
    };

    mockRestaurantService = {
      createRestaurant: cy.stub().returns(of({ id: 'new-rest-123', name: 'Test Restaurant' })).as('createRestaurant')
    };
    
    // Mount component with all necessary dependencies
    cy.mount(NewRestaurantComponent, {
      imports: [
        CommonModule,
        FormsModule,
        ReactiveFormsModule,
        MatCardModule,
        MatFormFieldModule,
        MatInputModule,
        MatButtonModule,
        BrowserAnimationsModule
      ],
      providers: [
        FormBuilder,
        { provide: MatDialogRef, useValue: mockDialogRef },
        { provide: MAT_DIALOG_DATA, useValue: {} },
        { provide: RestaurantService, useValue: mockRestaurantService }
      ]
    });
  });

  it('should display the restaurant form with all required fields', () => {
    cy.get('form').should('exist');
    cy.get('input[formControlName="name"]').should('exist');
    cy.get('input[formControlName="location"]').should('exist');
    cy.get('textarea[formControlName="description"]').should('exist');
    cy.get('input[formControlName="phone"]').should('exist');
    cy.get('input[formControlName="openingHours"]').should('exist');
    cy.get('button[type="submit"]').should('exist');
  });

  it('should validate form fields and disable submit button when invalid', () => {
    // Initially the form should be invalid because required fields are empty
    cy.get('button[type="submit"]').should('be.disabled');
    
    // Fill in only name field
    cy.get('input[formControlName="name"]').type('Test Restaurant');
    
    // Form should still be invalid with only one required field
    cy.get('button[type="submit"]').should('be.disabled');
    
  });

  it('should submit the form successfully when all required fields are valid', () => {
    // Fill in required fields
    cy.get('input[formControlName="name"]').type('Test Restaurant');
    cy.get('input[formControlName="location"]').type('Test Location');
    
    // Fill in optional fields
    cy.get('textarea[formControlName="description"]').type('Test Description');
    cy.get('input[formControlName="phone"]').type('555-123-4567');
    cy.get('input[formControlName="openingHours"]').type('11:00 - 21:00');
    
    // Submit form
    cy.get('button[type="submit"]').click();
    
    // Verify service was called
    cy.get('@createRestaurant').should('have.been.called');
    cy.get('@dialogClose').should('have.been.called');
  });

  it('should handle file selection', () => {
    // Create a test file
    const testFile = new File(['test content'], 'test-image.jpg', { type: 'image/jpeg' });
    
    // Trigger file selection
    cy.get('input[type="file"]').selectFile({
      contents: testFile,
      fileName: 'test-image.jpg',
      mimeType: 'image/jpeg',
    });
    
    // Fill in required fields
    cy.get('input[formControlName="name"]').type('Test Restaurant');
    cy.get('input[formControlName="location"]').type('Test Location');
    
    // Submit form
    cy.get('button[type="submit"]').click();
    
    // Verify service was called with file data
    cy.get('@createRestaurant').should('have.been.called');
    cy.get('@dialogClose').should('have.been.called');
  });

  it('should handle form submission errors', () => {
    // Setup service to return error
    mockRestaurantService.createRestaurant = cy.stub()
      .returns({
        subscribe: (callbacks: any) => {
          callbacks.error(new Error('Test error'));
          return { unsubscribe: () => {} };
        }
      })
      .as('createRestaurantError');
    
    // Fill in required fields
    cy.get('input[formControlName="name"]').type('Test Restaurant');
    cy.get('input[formControlName="location"]').type('Test Location');
    
    // Submit form
    cy.get('button[type="submit"]').click();
    
    // Verify error handling (alert should be shown)
    cy.on('window:alert', (text) => {
      expect(text).to.equal('Error creating the restaurant.');
    });
    
    // Dialog should not be closed on error
    cy.get('@dialogClose').should('not.have.been.called');
  });
});