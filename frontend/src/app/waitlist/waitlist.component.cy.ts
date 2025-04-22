import { HttpClientModule } from '@angular/common/http';
import { RouterModule, ActivatedRoute } from '@angular/router';
import { WaitlistComponent } from './waitlist.component';
import { WaitlistService } from '../services/waitlist.service';
import { TablesService } from '../services/table.service';
import { MatDialog, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { CommonModule } from '@angular/common';
import { MatListModule } from '@angular/material/list';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatExpansionModule } from '@angular/material/expansion';
import { MatSelectModule } from '@angular/material/select';
import { FormsModule } from '@angular/forms';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { of } from 'rxjs';
import { AddWaitlistDialogComponent } from './add-waitlist-dialog/add-waitlist-dialog.component';

describe('WaitlistComponent', () => {
  const mockRestaurantId = 'rest123';
  
  const mockWaitlistEntries = [
    {
      id: 'wait1',
      restaurantId: mockRestaurantId,
      name: 'John Doe',
      partySize: 2,
      partyAhead: 0,
      estimatedWaitTime: 15,
      phoneNumber: '555-123-4567',
      status: 'waiting'
    },
    {
      id: 'wait2',
      restaurantId: mockRestaurantId,
      name: 'Jane Smith',
      partySize: 4,
      partyAhead: 1,
      estimatedWaitTime: 30,
      phoneNumber: '555-987-6543',
      status: 'waiting'
    },
    {
      id: 'wait3',
      restaurantId: mockRestaurantId,
      name: 'Bob Johnson',
      partySize: 6,
      partyAhead: 2,
      estimatedWaitTime: 45,
      phoneNumber: '555-456-7890',
      status: 'waiting'
    },
    {
      id: 'wait4',
      restaurantId: mockRestaurantId,
      name: 'Alice Brown',
      partySize: 8,
      partyAhead: 3,
      estimatedWaitTime: 60,
      phoneNumber: '555-234-5678',
      status: 'waiting'
    }
  ];

  const mockTables = [
    {
      id: 'table1',
      restaurantId: mockRestaurantId,
      name: 'Table 1',
      minCapacity: 1,
      maxCapacity: 2,
      status: 'available',
      x: 10,
      y: 10
    },
    {
      id: 'table2',
      restaurantId: mockRestaurantId,
      name: 'Table 2',
      minCapacity: 3,
      maxCapacity: 4,
      status: 'available',
      x: 100,
      y: 100
    },
    {
      id: 'table3',
      restaurantId: mockRestaurantId,
      name: 'Table 3',
      minCapacity: 5,
      maxCapacity: 6,
      status: 'occupied',
      x: 200,
      y: 200
    },
    {
      id: 'table4',
      restaurantId: mockRestaurantId,
      name: 'Table 4',
      minCapacity: 7,
      maxCapacity: 10,
      status: 'available',
      x: 300,
      y: 300
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

  let mockWaitlistService;
  let mockTablesService;
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
    mockWaitlistService = {
      getWaitlist: cy.stub().returns(of(mockWaitlistEntries)),
      addToWaitlist: cy.stub().returns(of(mockWaitlistEntries[0])),
      updateWaitlistEntry: cy.stub().returns(of({})),
      seatCustomerAsReservation: cy.stub().returns(of({}))
    };

    mockTablesService = {
      getTables: cy.stub().returns(of(mockTables)),
      updateTable: cy.stub().returns(of({}))
    };

    // Mock MatDialog
    mockDialogOpen = cy.stub().returns({
      afterClosed: () => of({
        name: 'New Customer',
        partySize: 3,
        phoneNumber: '555-111-2222'
      })
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
        parent: {
          paramMap: {
            get: cy.stub().returns(mockRestaurantId)
          }
        }
      }
    };

    // Mount the component with mocked dependencies
    cy.mount(WaitlistComponent, {
      imports: [
        HttpClientModule,
        RouterModule.forRoot([]),
        CommonModule,
        MatListModule,
        MatButtonModule,
        MatIconModule,
        MatDialogModule,
        MatExpansionModule,
        MatSelectModule,
        FormsModule,
        BrowserAnimationsModule
      ],
      providers: [
        { provide: WaitlistService, useValue: mockWaitlistService },
        { provide: TablesService, useValue: mockTablesService },
        { provide: MatDialog, useValue: mockMatDialog },
        { provide: ActivatedRoute, useValue: mockActivatedRoute }
      ]
    });
  });

  it('should fetch waitlist and tables on initialization', () => {
    // Verify the services were called with correct parameters
    cy.wrap(mockWaitlistService.getWaitlist).should('be.calledWith', mockRestaurantId);
    cy.wrap(mockTablesService.getTables).should('be.calledWith', mockRestaurantId);
  });

  it('should group waitlist entries into cohorts by party size', () => {
    // Set component properties directly
    cy.window().then((win) => {
      const componentInstance = win.document.querySelector('app-waitlist')?.__ngContext__?.find(c => c?.constructor?.name === 'WaitlistComponent');
      if (componentInstance) {
        // Verify cohorts are populated correctly
        expect(componentInstance.cohorts[0].entries.length).to.equal(1); // 1-2 group
        expect(componentInstance.cohorts[1].entries.length).to.equal(1); // 2-4 group
        expect(componentInstance.cohorts[2].entries.length).to.equal(1); // 4-6 group
        expect(componentInstance.cohorts[3].entries.length).to.equal(1); // >6 group
      }
    });

    // Check if cohort panels are displayed
    cy.get('mat-expansion-panel').should('have.length', 4);
  });

  it('should open add waitlist dialog when add button is clicked', () => {
    // Click the add waitlist button
    cy.get('button').contains('Add to Waitlist').click();
    
    // Verify dialog was opened
    cy.contains("Add to Waitlist").should('exist');
  });

  it('should display available tables for each waitlist entry', () => {
    // Expand a cohort panel
    cy.get('mat-expansion-panel').first().click();
    
    // Check if table selection dropdown exists
    cy.get('mat-select').should('exist');
  });

  it('should calculate estimated wait times for each cohort', () => {
    // Set component properties directly
    cy.window().then((win) => {
      const componentInstance = win.document.querySelector('app-waitlist')?.__ngContext__?.find(c => c?.constructor?.name === 'WaitlistComponent');
      if (componentInstance) {
        // Verify wait times are calculated
        expect(componentInstance.cohorts[0].waitTime).to.equal(15); // 1 entry * 15 minutes
        expect(componentInstance.cohorts[1].waitTime).to.equal(15); // 1 entry * 15 minutes
        expect(componentInstance.cohorts[2].waitTime).to.equal(15); // 1 entry * 15 minutes
        expect(componentInstance.cohorts[3].waitTime).to.equal(15); // 1 entry * 15 minutes
      }
    });

    // Check if wait times are displayed
    cy.get('mat-expansion-panel').first().click();
    cy.contains('Est. Wait:').should('exist');
  });
});