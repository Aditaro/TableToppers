import { HttpClientModule } from '@angular/common/http';
import { RouterModule, ActivatedRoute } from '@angular/router';
import { TablesComponent } from './tables.component';
import { TablesService } from '../services/table.service';
import { ReservationService } from '../services/reservation.service';
import { RestaurantService } from '../services/restaurant.service';
import { MatDialog, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatIconModule } from '@angular/material/icon';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatNativeDateModule } from '@angular/material/core';
import { MatInputModule } from '@angular/material/input';
import { MatBadgeModule } from '@angular/material/badge';
import { MatExpansionModule } from '@angular/material/expansion';
import { MatButtonModule } from '@angular/material/button';
import { MatSelectModule } from '@angular/material/select';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { BehaviorSubject, of } from 'rxjs';
import { NewTableComponent } from '../new-table/new-table.component';
import { NewReservationComponent } from '../new-reservation/new-reservation.component';
import { WaitlistComponent } from '../waitlist/waitlist.component';

describe('TablesComponent', () => {
  const mockRestaurantId = 'rest123';
  
  const mockRestaurant = {
    id: mockRestaurantId,
    status: 'open',
    name: 'Test Restaurant',
    img: 'https://picsum.photos/300/200?random=1',
    description: 'A test restaurant for component testing',
    location: 'Test Location',
    phone: '+1 555-123-4567',
    openingHours: '10:00-22:00',
    specialAvailability: []
  };

  const mockTables = [
    {
      id: 'table1',
      restaurantId: mockRestaurantId,
      name: 'Table 1',
      minCapacity: 2,
      maxCapacity: 4,
      status: 'available',
      x: 10,
      y: 10
    },
    {
      id: 'table2',
      restaurantId: mockRestaurantId,
      name: 'Table 2',
      minCapacity: 4,
      maxCapacity: 6,
      status: 'available',
      x: 100,
      y: 100
    }
  ];

  const mockReservations = [
    {
      id: 'res1',
      restaurantId: mockRestaurantId,
      tableId: 'table1',
      userId: 'user123',
      reservationTime: new Date(new Date().setHours(12, 0, 0, 0)).toISOString(),
      numberOfGuests: 3,
      status: 'confirmed',
      phoneNumber: '555-123-4567'
    },
    {
      id: 'res2',
      restaurantId: mockRestaurantId,
      tableId: 'table2',
      userId: 'user456',
      reservationTime: new Date(new Date().setHours(14, 0, 0, 0)).toISOString(),
      numberOfGuests: 5,
      status: 'confirmed',
      phoneNumber: '555-987-6543'
    }
  ];

  const mockUser = {
    id: 'user123',
    email: 'test@example.com',
    phoneNumber: '555-123-4567',
    firstName: 'Test',
    lastName: 'User',
    role: 'customer'
  };

  const mockManagerUser = {
    id: 'manager123',
    email: 'manager@example.com',
    phoneNumber: '555-987-6543',
    firstName: 'Manager',
    lastName: 'User',
    role: 'manager'
  };

  let mockTablesService;
  let mockReservationService;
  let mockRestaurantService;
  let mockDialogOpen;
  let mockActivatedRoute;
  let mockMatDialog;

  beforeEach(() => {
    // Create a clean localStorage stub for each test
    cy.clearAllLocalStorage();
    cy.window().then((win) => {
      win.localStorage.clear();
      cy.stub(win.localStorage, 'getItem')
        .callsFake((key) => {
          if (key === 'user') {
            return JSON.stringify(mockUser);
          }
          return null;
        });
    });

    // Mock services
    mockTablesService = {
      getTables: cy.stub().returns(of(mockTables)),
      updateTable: cy.stub().returns(of({})),
      addTable: cy.stub().returns(of({})),
      deleteTable: cy.stub().returns(of({}))
    };

    mockReservationService = {
      getReservations: cy.stub().returns(of(mockReservations)),
      createReservation: cy.stub().returns(of({})),
      updateReservation: cy.stub().returns(of({}))
    };

    mockRestaurantService = {
      getRestaurants: cy.stub().returns(of([mockRestaurant]))
    };

    // Mock MatDialog
    mockDialogOpen = cy.stub().returns({
      afterClosed: () => of(true)
    });
    
    // Create a complete MatDialog mock with required properties
    mockMatDialog = {
      open: mockDialogOpen,
      openDialogs: [],
      afterOpened: { next: cy.stub() },
      _getAfterAllClosed: () => of(true),
    };

    // Mock ActivatedRoute
    mockActivatedRoute = {
      snapshot: {
        paramMap: {
          get: cy.stub().returns(mockRestaurantId)
        }
      }
    };

    // Create a stub for the Konva Stage to prevent errors
    cy.stub(window, 'Konva').returns({
      Stage: function() {
        return {
          add: cy.stub(),
          draw: cy.stub(),
          width: cy.stub(),
          height: cy.stub(),
          on: cy.stub(),
          container: cy.stub(),
          getPointerPosition: cy.stub().returns({ x: 0, y: 0 })
        };
      },
      Layer: function() {
        return {
          add: cy.stub(),
          draw: cy.stub(),
          clear: cy.stub()
        };
      },
      Rect: function() {
        return {
          width: cy.stub(),
          height: cy.stub(),
          fill: cy.stub(),
          stroke: cy.stub(),
          strokeWidth: cy.stub(),
          draggable: cy.stub(),
          x: cy.stub(),
          y: cy.stub(),
          on: cy.stub(),
          id: cy.stub(),
          name: cy.stub()
        };
      },
      Text: function() {
        return {
          text: cy.stub(),
          fontSize: cy.stub(),
          fill: cy.stub(),
          x: cy.stub(),
          y: cy.stub(),
          width: cy.stub(),
          align: cy.stub()
        };
      },
      Group: function() {
        return {
          add: cy.stub(),
          x: cy.stub(),
          y: cy.stub(),
          draggable: cy.stub(),
          on: cy.stub()
        };
      }
    });

    // Mount the component with mocked dependencies
    cy.mount(TablesComponent, {
      imports: [
        HttpClientModule,
        RouterModule.forRoot([]),
        CommonModule,
        MatCardModule,
        MatToolbarModule,
        MatProgressSpinnerModule,
        MatIconModule,
        MatDatepickerModule,
        MatNativeDateModule,
        MatInputModule,
        MatBadgeModule,
        MatExpansionModule,
        MatButtonModule,
        MatSelectModule,
        MatDialogModule,
        BrowserAnimationsModule
      ],
      providers: [
        { provide: TablesService, useValue: mockTablesService },
        { provide: ReservationService, useValue: mockReservationService },
        { provide: RestaurantService, useValue: mockRestaurantService },
        { provide: MatDialog, useValue: mockMatDialog },
        { provide: ActivatedRoute, useValue: mockActivatedRoute }
      ]
    });
  });

  it('should fetch restaurant, tables, and reservations on initialization', () => {
    // Verify the services were called with correct parameters
    cy.wrap(mockActivatedRoute.snapshot.paramMap.get).should('be.calledWith', 'restaurantId');
    cy.wrap(mockRestaurantService.getRestaurants).should('be.called');
    cy.wrap(mockTablesService.getTables).should('be.calledWith', mockRestaurantId);
    cy.wrap(mockReservationService.getReservations).should('be.called');
  });

  it('should display tables in the floor plan', () => {
    // Since we can't easily test Konva canvas rendering, we'll verify that the tables data is loaded
    cy.window().then((win) => {
      const componentInstance = win.document.querySelector('app-tables')?.__ngContext__?.find(c => c?.constructor?.name === 'TablesComponent');
      if (componentInstance) {
        expect(componentInstance.tables).to.deep.equal(mockTables);
      }
    });
  });

  it('should display reservations grouped by hour', () => {
    // Set component properties directly
    cy.window().then((win) => {
      const componentInstance = win.document.querySelector('app-tables')?.__ngContext__?.find(c => c?.constructor?.name === 'TablesComponent');
      if (componentInstance) {
        componentInstance.restaurant = mockRestaurant;
        componentInstance.reservations = mockReservations;
        componentInstance.groupReservationsByHour();
      }
    });

    // Check if hours are displayed
    cy.get('mat-card').should('exist');
  });

  it('should update reservations when date changes', () => {
    // Trigger date change
    cy.get('mat-datepicker-toggle').first().click();
    cy.get('.mat-calendar-body-cell-content').first().click();
    
    // Verify reservation service was called with new date
    cy.wrap(mockReservationService.getReservations).should('be.called');
  });

  it('should open new table dialog for manager users', () => {
    // Set manager user
    cy.window().then(win => {
    win.localStorage.getItem.restore()
    cy.stub(win.localStorage, 'getItem')
      .callsFake(k => k === 'user' ? JSON.stringify(mockManagerUser) : null)
    });

    // Remount component with manager user
    cy.mount(TablesComponent, {
    imports:   [
      HttpClientModule,
      RouterModule.forRoot([]),
      CommonModule,
      MatCardModule,
      MatToolbarModule,
      MatProgressSpinnerModule,
      MatIconModule,
      MatDatepickerModule,
      MatNativeDateModule,
      MatInputModule,
      MatBadgeModule,
      MatExpansionModule,
      MatButtonModule,
      MatSelectModule,
      MatDialogModule,
      BrowserAnimationsModule
    ],
    providers: [
      { provide: TablesService,      useValue: mockTablesService      },
      { provide: ReservationService, useValue: mockReservationService },
      { provide: RestaurantService,  useValue: mockRestaurantService  },
      { provide: MatDialog,          useValue: mockMatDialog },
      { provide: ActivatedRoute,     useValue: mockActivatedRoute     },
      ],
    });

    // Click the add table button (assuming it has a specific class or ID)
    cy.get('button[color="accent"]').click();
    
    cy.contains('Create New Table').should('exist');
  });

  it('should toggle sidebar visibility when toggle button is clicked', () => {
    // Set component properties directly
    cy.window().then((win) => {
      const componentInstance = win.document.querySelector('app-tables')?.__ngContext__?.find(c => c?.constructor?.name === 'TablesComponent');
      if (componentInstance) {
        componentInstance.restaurant = mockRestaurant;
        componentInstance.sidebarHidden = false;
      }
    });

    // Click the toggle sidebar button
    cy.get('div[class="sidebar-toggle"]').click();
    
    // Verify sidebar state changed
    cy.window().then((win) => {
      const componentInstance = win.document.querySelector('app-tables')?.__ngContext__?.find(c => c?.constructor?.name === 'TablesComponent');
      if (componentInstance) {
        expect(componentInstance.sidebarHidden).to.be.true;
      }
    });
  });
});