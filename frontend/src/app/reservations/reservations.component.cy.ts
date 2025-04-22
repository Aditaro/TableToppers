import { HttpClientModule } from '@angular/common/http';
import { RouterModule, Router, ActivatedRoute } from '@angular/router';
import { ReservationsComponent } from './reservations.component';
import { ReservationService } from '../services/reservation.service';
import { MatDialog, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { BehaviorSubject, of } from 'rxjs';
import { NewReservationComponent } from '../new-reservation/new-reservation.component';

describe('ReservationsComponent', () => {
  const mockReservations = [
    {
      id: '1',
      restaurantId: 'rest123',
      userId: 'user123',
      tableId: 'table1',
      reservationTime: '2023-05-15T18:00:00Z',
      numberOfGuests: 4,
      status: 'confirmed',
      phoneNumber: '555-123-4567'
    },
    {
      id: '2',
      restaurantId: 'rest123',
      userId: 'user123',
      tableId: 'table2',
      reservationTime: '2023-05-20T19:30:00Z',
      numberOfGuests: 2,
      status: 'pending',
      phoneNumber: '555-987-6543'
    }
  ];

  const mockRestaurant = {
    id: 'rest123',
    name: 'Test Restaurant',
    location: '123 Test St',
    openingHours: '9AM-10PM',
    phone: '555-123-4567'
  };

  const mockUser = {
    id: 'user123',
    email: 'test@example.com',
    phoneNumber: '555-123-4567',
    firstName: 'Test',
    lastName: 'User',
    role: 'customer'
  };

  let mockReservationService;
  let mockDialogOpen;
  let mockRouter;
  let mockActivatedRoute;

  beforeEach(() => {
    // Mock localStorage
    cy.stub(localStorage, 'getItem').returns(JSON.stringify(mockUser));

    // Mock ReservationService
    mockReservationService = {
      getReservations: cy.stub().returns(of(mockReservations)),
      deleteReservation: cy.stub().returns(of({})),
      updateReservation: cy.stub().returns(of({}))
    };

    // Mock MatDialog
    mockDialogOpen = cy.stub().returns({
      afterClosed: () => of(true)
    });

    // Mock Router
    mockRouter = {
      getCurrentNavigation: cy.stub().returns({
        extras: {
          state: {
            restaurant: mockRestaurant
          }
        }
      })
    };

    // Mock ActivatedRoute
    mockActivatedRoute = {
      snapshot: {
        paramMap: {
          get: cy.stub().returns('rest123')
        }
      }
    };

    // Mount the component with mocked dependencies
    cy.mount(ReservationsComponent, {
      imports: [HttpClientModule, RouterModule.forRoot([]), CommonModule, MatCardModule, MatDialogModule],
      providers: [
        { provide: ReservationService, useValue: mockReservationService },
        { provide: MatDialog, useValue: { open: mockDialogOpen } },
        { provide: Router, useValue: mockRouter },
        { provide: ActivatedRoute, useValue: mockActivatedRoute }
      ]
    });
  });

  it('should display the restaurant information', () => {
    cy.get('.restaurant-info h2').should('contain.text', mockRestaurant.name);
    cy.get('.restaurant-info p').eq(0).should('contain.text', mockRestaurant.location);
    cy.get('.restaurant-info p').eq(1).should('contain.text', mockRestaurant.openingHours);
    cy.get('.restaurant-info p').eq(2).should('contain.text', mockRestaurant.phone);
  });

  it('should fetch and display reservations', () => {
    // Verify the service was called with correct parameters
    cy.wrap(mockReservationService.getReservations).should('be.calledWith', 'rest123', 'user123');
    
    // Check if reservations are displayed
    cy.get('.reservations-card').should('have.length', 2);
    cy.get('.reservation-time').first().should('exist');
    cy.get('.reservation-guests').first().should('contain.text', '4 guests');
  });

  it('should open dialog when new reservation button is clicked', () => {
    cy.get('#new-reservation-btn').click();
    
    // Verify dialog was opened with correct data
    cy.wrap(mockDialogOpen).should('be.calledWith', NewReservationComponent, {
      width: '600px',
      data: {
        restaurantId: 'rest123',
        restaurant: mockRestaurant
      }
    });
  });

  it('should delete a reservation when delete button is clicked', () => {
    // Click the delete button on the first reservation
    cy.get('.delete-btn').first().click();
    
    // Verify delete service was called with correct parameters
    cy.wrap(mockReservationService.deleteReservation).should('be.calledWith', 'rest123', mockReservations[0].id);
    
    // Verify getReservations was called again to refresh the list
    cy.wrap(mockReservationService.getReservations).should('be.calledTwice');
  });

  it('should open dialog to modify a reservation when edit button is clicked', () => {
    // Click the edit button on the first reservation
    cy.get('.edit-btn').first().click();
    
    // Verify dialog was opened with correct data including the reservation
    cy.wrap(mockDialogOpen).should('be.calledWith', NewReservationComponent, {
      width: '600px',
      data: {
        restaurantId: 'rest123',
        restaurant: mockRestaurant,
        reservation: mockReservations[0]
      }
    });
  });
});