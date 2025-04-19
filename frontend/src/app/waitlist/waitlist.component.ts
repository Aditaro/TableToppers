import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { MatDialog, MatDialogModule } from '@angular/material/dialog'; // Import MatDialogModule
import { ActivatedRoute } from '@angular/router';
import { CommonModule } from '@angular/common'; // Import CommonModule
import { MatListModule } from '@angular/material/list'; // Import MatListModule
import { MatButtonModule } from '@angular/material/button'; // Import MatButtonModule
import { MatIconModule } from '@angular/material/icon'; // Import MatIconModule
import { MatExpansionModule } from '@angular/material/expansion'; // Import MatExpansionModule for collapsible panels
import { MatSelectModule } from '@angular/material/select'; // Import MatSelectModule for dropdown
import { FormsModule } from '@angular/forms'; // Import FormsModule for ngModel
import { AddWaitlistDialogComponent, AddWaitlistDialogResult } from './add-waitlist-dialog/add-waitlist-dialog.component'; // Import the dialog component
import { WaitlistService } from '../services/waitlist.service';
import { WaitlistEntry, WaitlistEntryCreate } from '../models/waitlist.model';
import { Table } from '../models/table.model';
import { TablesService } from '../services/table.service';

@Component({
  selector: 'app-waitlist',
  standalone: true, // Make component standalone
  imports: [ // Add necessary imports for standalone component
    CommonModule,
    MatListModule,
    MatButtonModule,
    MatIconModule,
    MatDialogModule,
    MatExpansionModule,
    MatSelectModule,
    FormsModule,
    AddWaitlistDialogComponent // Import the dialog component here as well if needed, or ensure it's standalone
  ],
  templateUrl: './waitlist.component.html',
  styleUrls: ['./waitlist.component.css']
})
export class WaitlistComponent implements OnInit {
  @Input() restaurantId!: string;
  @Output() tableStatusChanged = new EventEmitter<void>();
  waitlist: WaitlistEntry[] = [];
  tables: Table[] = [];
  selectedTableIds: {[entryId: string]: string} = {}; // Map to store selected table ID for each entry
  
  // Cohort definitions
  cohorts = [
    { name: '1-2 Groups', min: 0, max: 2, entries: [] as WaitlistEntry[], waitTime: 0 },
    { name: '2-4 Groups', min: 3, max: 4, entries: [] as WaitlistEntry[], waitTime: 0 },
    { name: '4-6 Groups', min: 5, max: 6, entries: [] as WaitlistEntry[], waitTime: 0 },
    { name: '>6 Groups', min: 7, max: 999, entries: [] as WaitlistEntry[], waitTime: 0 }
  ];

  constructor(
    private waitlistService: WaitlistService,
    private tablesService: TablesService,
    private dialog: MatDialog,
    private route: ActivatedRoute
  ) {}

  ngOnInit(): void {
    if (!this.restaurantId) {
      // Fallback if Input is not provided (e.g., direct routing)
      this.restaurantId = this.route.snapshot.parent?.paramMap.get('restaurantId') || '';
    }
    if (this.restaurantId) {
      this.fetchWaitlist();
      this.fetchTables();
    }
  }
  
  // Group waitlist entries into cohorts based on party size
  groupWaitlistByCohort(): void {
    // Reset cohorts
    this.cohorts.forEach(cohort => {
      cohort.entries = [];
      cohort.waitTime = 0;
    });
    
    this.waitlist.filter(entry => entry.status === 'waiting').forEach(entry => {
      const cohort = this.cohorts.find(c => 
        entry.partySize >= c.min && entry.partySize <= c.max
      );
      
      if (cohort) {
        cohort.entries.push(entry);
      }
    });
    
    // Calculate estimated wait time for the last entry in each cohort
    this.calculateWaitTimes();
  }
  
  // Calculate estimated wait times for each cohort
  calculateWaitTimes(): void {
    this.cohorts.forEach(cohort => {
      if (cohort.entries.length > 0) {
        // Simple estimation: 15 minutes per group in the cohort
        // This is a placeholder - real implementation would use more sophisticated logic
        cohort.waitTime = cohort.entries.length * 15;
      }
    });
  }

  fetchWaitlist(): void {
    this.waitlistService.getWaitlist(this.restaurantId)
      .subscribe(data => {
        this.waitlist = data;
        this.groupWaitlistByCohort(); // Group entries after fetching
        console.log('Fetched waitlist:', this.waitlist);
        console.log('Cohorts:', this.cohorts);
      }, error => {
        console.error('Error fetching waitlist:', error);
        // Handle error (e.g., show message to user)
      });
  }

  fetchTables(): void {
    this.tablesService.getTables(this.restaurantId)
      .subscribe(data => {
        this.tables = data;
      }, error => {
        console.error('Error fetching tables:', error);
      });
  }

  openAddWaitlistDialog(): void {
    const dialogRef = this.dialog.open(AddWaitlistDialogComponent, {
      width: '400px',
      data: { restaurantId: this.restaurantId } // Pass necessary data
    });

    dialogRef.afterClosed().subscribe((result: AddWaitlistDialogResult | undefined) => {
      if (result) {
        // Map AddWaitlistDialogResult to WaitlistEntryCreate if needed, or adjust interfaces
        const entryData: WaitlistEntryCreate = {
            name: result.name,
            partySize: result.partySize,
            phoneNumber: result.phoneNumber,
        };
        this.addCustomerToWaitlist(entryData);
      }
    });
  }

  addCustomerToWaitlist(entryData: WaitlistEntryCreate): void {
    // Logic to check for available tables first
    const availableTable = this.findAvailableTable(entryData.partySize);

    if (availableTable) {
      // Assign table (might involve another API call or just UI update)
      console.log(`Assigning table ${availableTable.name} to ${entryData.name}`);
      // TODO: Implement table assignment logic (e.g., update table status)
    } else {
      // Add to waitlist via service
      this.waitlistService.addToWaitlist(this.restaurantId, entryData)
        .subscribe(newEntry => {
          this.waitlist.push(newEntry);
          this.groupWaitlistByCohort(); // Regroup after adding
          console.log('Added to waitlist:', newEntry);
          // TODO: Provide user feedback
        }, error => {
          console.error('Error adding to waitlist:', error);
          // Handle error
        });
    }
  }

  findAvailableTable(partySize: number): Table | undefined {
    console.log('Finding available table for party size:', partySize);
    console.log('Tables:', this.tables);
    return this.tables.find(table =>
      table.status === 'available' &&
      partySize >= table.minCapacity &&
      partySize <= table.maxCapacity
    );
  }

  // Get available tables for a specific party size
  getAvailableTablesForParty(partySize: number): Table[] {
    return this.tables.filter(table => 
      table.status === 'available' &&
      partySize >= table.minCapacity &&
      partySize <= table.maxCapacity
    );
  }

  // Seat customer with selected table
  seatCustomer(entry: WaitlistEntry): void {
    console.log(`Seating customer ${entry.name}`);
    
    // Fetch tables if not already loaded
    if (this.tables.length === 0) {
      this.fetchTables();
      return; // Wait for tables to load
    }
    
    // Get the selected table ID for this entry
    const selectedTableId = this.selectedTableIds[entry.id!];
    
    if (!selectedTableId) {
      console.error('No table selected');
      return;
    }
    
    // Find the selected table
    const selectedTable = this.tables.find(table => table.id === selectedTableId);
    
    if (selectedTable) {
      // Update local table status immediately for UI feedback
      selectedTable.status = 'occupied';
      
      // First update the table status in the database
      this.tablesService.updateTable(this.restaurantId, selectedTableId, { status: 'occupied' })
        .subscribe({
          next: () => {
            // Treat seating as a special case of reservation
            this.waitlistService.seatCustomerAsReservation(this.restaurantId, entry, selectedTable.id!)
              .subscribe({
                next: () => {
                  // Update the waitlist entry status to seated
                  this.waitlistService.updateWaitlistEntry(this.restaurantId, entry.id!, 'seated', selectedTable.id!)
                    .subscribe({
                      next: () => {
                        // Update local state
                        const index = this.waitlist.findIndex(e => e.id === entry.id);
                        if (index !== -1) {
                          this.waitlist[index].status = 'seated';
                          this.groupWaitlistByCohort(); // Regroup after status change
                        }
                        
                        console.log(`Table ${selectedTable.name} assigned to ${entry.name}`);
                        // Clear the selection
                        delete this.selectedTableIds[entry.id!];
                        
                        // Force a refresh of all tables to ensure consistent state
                        this.tablesService.getTables(this.restaurantId).subscribe(updatedTables => {
                          this.tables = updatedTables;
                          // Emit event to notify parent component to update table colors
                          this.tableStatusChanged.emit();
                        });
                      },
                      error: (err) => {
                        console.error('Error updating waitlist entry status:', err);
                        // Revert local table status on error
                        selectedTable.status = 'available';
                      }
                    });
                },
                error: (err) => {
                  console.error('Error creating reservation for waitlist entry:', err);
                  // Revert local table status on error
                  selectedTable.status = 'available';
                }
              });
          },
          error: (err) => {
            console.error('Error updating table status:', err);
            // Revert local table status on error
            selectedTable.status = 'available';
          }
        });
    } else {
      console.error('Selected table not found');
    }
  }

  cancelEntry(entry: WaitlistEntry): void {
    console.log(`Cancelling entry for ${entry.name}`);
    this.waitlistService.updateWaitlistEntry(this.restaurantId, entry.id!, 'cancelled')
      .subscribe(() => {
        // Update local state
        const index = this.waitlist.findIndex(e => e.id === entry.id);
        if (index !== -1) {
          this.waitlist[index].status = 'cancelled';
          this.groupWaitlistByCohort(); // Regroup after status change
        }
      }, error => {
        console.error('Error cancelling entry:', error);
      });
  }
}

// Remove TODO for creating AddWaitlistDialogComponent as it's created