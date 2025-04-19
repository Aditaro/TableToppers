import { NewReservationComponent } from './new-reservation.component'
import { FormBuilder, ReactiveFormsModule } from '@angular/forms'
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog'
import { MatSnackBar } from '@angular/material/snack-bar'
import { ReservationService } from '../services/reservation.service'
import { ActivatedRoute } from '@angular/router'
import { CommonModule } from '@angular/common'
import { MatFormFieldModule } from '@angular/material/form-field'
import { MatInputModule } from '@angular/material/input'
import { MatDatepickerModule } from '@angular/material/datepicker'
import { MatNativeDateModule, provideNativeDateAdapter } from '@angular/material/core'
import { MatButtonModule } from '@angular/material/button'
import { MatSelectModule } from '@angular/material/select'
import { NgxMaskDirective, provideNgxMask } from 'ngx-mask'
import { BrowserAnimationsModule } from '@angular/platform-browser/animations'
import { of } from 'rxjs'

describe('NewReservationComponent', () => {
  // Mock data for testing
  const mockRestaurant = {
    id: 'rest123',
    name: 'Test Restaurant',
    openingHours: '11:00-22:00',
    specialAvailability: [
      { date: '2023-12-25', status: 'closed' }
    ]
  };

  const mockReservation = {
    id: 'res123',
    restaurantId: 'rest123',
    userId: 'user123',
    tableId: 'table123',
    reservationTime: new Date(new Date().setDate(new Date().getDate() + 1)).toISOString(),
    numberOfGuests: 4,
    status: 'confirmed',
    phoneNumber: '(555) 123-4567'
  };

  // Define mock objects without stubs initially
  let mockDialogRef;
  let mockReservationService;
  let mockSnackBar;

  beforeEach(() => {
    // Initialize mock objects with stubs inside the test context
    mockDialogRef = {
      close: cy.stub().as('dialogClose')
    };

    mockReservationService = {
      createReservation: cy.stub().returns(of({ id: 'new-res-123' })).as('createReservation'),
      updateReservation: cy.stub().returns(of({})).as('updateReservation')
    };

    mockSnackBar = {
      open: cy.stub().as('snackBarOpen')
    };
    
    // Mount component with all necessary dependencies
    cy.mount(NewReservationComponent, {
      imports: [
        CommonModule,
        ReactiveFormsModule,
        MatFormFieldModule,
        MatInputModule,
        MatDatepickerModule,
        MatNativeDateModule,
        MatButtonModule,
        MatSelectModule,
        BrowserAnimationsModule
      ],
      providers: [
        FormBuilder,
        provideNgxMask(),
        provideNativeDateAdapter(),
        { provide: MatDialogRef, useValue: mockDialogRef },
        { provide: MAT_DIALOG_DATA, useValue: { restaurant: mockRestaurant, restaurantId: mockRestaurant.id } },
        { provide: ReservationService, useValue: mockReservationService },
        { provide: MatSnackBar, useValue: mockSnackBar },
        { provide: ActivatedRoute, useValue: { snapshot: { paramMap: { get: () => null } } } }
      ]
    });
  });

  it('should display the reservation form with all required fields', () => {
    cy.get('form').should('exist');
    cy.get('input[formControlName="reservationDate"]').should('exist');
    cy.get('mat-select[formControlName="reservationTime"]').should('exist');
    cy.get('input[formControlName="numberOfGuests"]').should('exist');
    cy.get('input[formControlName="phoneNumber"]').should('exist');
    cy.get('button[type="submit"]').should('exist').and('contain', 'Reserve Table');
  });

  it('should validate form fields and disable submit button when invalid', () => {
    // Initially the form should be invalid
    cy.get('button[type="submit"]').should('be.disabled');

    // Fill in date but leave other fields empty
    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);
    cy.get('input[formControlName="reservationDate"]').type(tomorrow.toLocaleDateString());
    cy.get('button[type="submit"]').should('be.disabled');
  });

  it('should generate time slots when a date is selected', () => {
    // Select a date
    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);
    cy.get('input[formControlName="reservationDate"]').type(tomorrow.toLocaleDateString());
    
    // Open time dropdown
    cy.get('mat-select[formControlName="reservationTime"]').click();
    
    // Verify time slots are generated
    cy.get('mat-option').should('have.length.greaterThan', 0);
  });

  it('should submit the form successfully when all fields are valid', () => {
    // Fill in all required fields
    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);
    cy.get('input[formControlName="reservationDate"]').type(tomorrow.toLocaleDateString());
    
    // Select time
    cy.get('mat-select[formControlName="reservationTime"]').click();
    cy.get('mat-option').first().click();
    
    // Fill in guests and phone
    cy.get('input[formControlName="numberOfGuests"]').type('4');
    cy.get('input[formControlName="phoneNumber"]').type('5551234567');
    
    // Submit form
    cy.get('button[type="submit"]').should('not.be.disabled').click();
    
    // Verify service was called
    cy.get('@createReservation').should('have.been.called');
    cy.get('@dialogClose').should('have.been.called');
  });

  it('should handle edit mode correctly', () => {
    // Mount component in edit mode with the reservation data
    cy.mount(NewReservationComponent, {
      imports: [
        CommonModule,
        ReactiveFormsModule,
        MatFormFieldModule,
        MatInputModule,
        MatDatepickerModule,
        MatNativeDateModule,
        MatButtonModule,
        MatSelectModule,
        BrowserAnimationsModule
      ],
      providers: [
        FormBuilder,
        provideNgxMask(),
        provideNativeDateAdapter(),
        { provide: MatDialogRef, useValue: mockDialogRef },
        { provide: MAT_DIALOG_DATA, useValue: {
          restaurant: mockRestaurant,
          restaurantId: mockRestaurant.id,
          reservation: mockReservation
        }},
        { provide: ReservationService, useValue: mockReservationService },
        { provide: MatSnackBar, useValue: mockSnackBar },
        { provide: ActivatedRoute, useValue: { snapshot: { paramMap: { get: () => null } } } }
      ]
    });
    
    // Verify button text changes in edit mode
    cy.get('button[type="submit"]').should('contain', 'Update Reservation');
    
    // Update a field and submit
    cy.get('input[formControlName="numberOfGuests"]').clear().type('6');
    cy.get('button[type="submit"]').click();
    
    // Verify update service was called
    cy.get('@updateReservation').should('have.been.called');
  });
})