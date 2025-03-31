import {AfterViewInit, Component, ElementRef, OnDestroy, OnInit, ViewChild} from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatIconModule } from '@angular/material/icon';
import { TablesService } from 'src/app/services/table.service';
import {NewTable, Table} from '../models/table.model';
import Konva from "konva";
import Stage = Konva.Stage;
import Layer = Konva.Layer;
import {Subscription} from "rxjs";
import {NewTableComponent} from "../new-table/new-table.component";
import {MatDialog} from "@angular/material/dialog";
import { MatDatepickerInputEvent, MatDatepickerModule } from '@angular/material/datepicker';
import { MatNativeDateModule } from '@angular/material/core';
import { MatInputModule } from '@angular/material/input';
import { Reservation } from '../models/reservation.model';
import { ReservationService } from '../services/reservation.service';
import { MatBadgeModule } from '@angular/material/badge';
import { RestaurantService } from '../services/restaurant.service';
import { Restaurant } from '../models/restaurant.model';
import { MatExpansionModule } from '@angular/material/expansion';
import { MatButtonModule } from '@angular/material/button';
import { MatSelectModule } from '@angular/material/select';
import { NewReservationComponent } from '../new-reservation/new-reservation.component';

@Component({
  standalone: true,
  selector: 'app-tables',
  imports: [
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
    MatSelectModule
  ],
  templateUrl: './tables.component.html',
  styleUrl: './tables.component.css',
})
export class TablesComponent implements OnInit, OnDestroy, AfterViewInit {
  @ViewChild('floorPlanHost', { static: true }) floorPlanHost!: ElementRef<HTMLDivElement>;

  restaurantId!: string;
  restaurant: Restaurant | null = null;
  tables: Table[] = [];
  loading = false;
  reservations: Reservation[] = [];
  selectedDate: Date = new Date();
  hourlyReservations: { [hour: string]: Reservation[] } = {};
  hours: string[] = [];
  filteredHours: string[] = [];
  pastHours: string[] = [];
  futureHours: string[] = [];
  currentHour: string = '';
  showPastHours: boolean = false;
  sidebarHidden: boolean = false;

    // Konva objects
  private stage!: Stage;
  private layer!: Layer;

  private viewInitialized = false;
  private tablesInitialized = false;

  private sub!: Subscription;

  constructor(
    private route: ActivatedRoute,
    private tablesService: TablesService,
    private reservationService: ReservationService,
    private restaurantService: RestaurantService,
    public dialog: MatDialog
  ) {}

  ngOnInit() {
    // 1. Grab restaurantId from route
    this.restaurantId = this.route.snapshot.paramMap.get('restaurantId') || '';

    // 2. Load restaurant details
    this.loadRestaurant();

    // 3. Load tables for the restaurant
    this.loadTables();
    
    // 4. Load reservations for today
    this.loadReservations(this.selectedDate);
  }

  loadRestaurant() {
    // Get restaurant details to access opening hours
    this.restaurantService.getRestaurants(undefined, undefined).subscribe({
      next: (restaurants) => {
        const restaurant = restaurants.find(r => r.id === this.restaurantId);
        if (restaurant) {
          this.restaurant = restaurant;
          // If we already have hours loaded, filter them by opening hours
          if (this.hours.length > 0) {
            this.filterHoursByOpeningHours();
          }
        }
      },
      error: (err) => {
        console.error('Error fetching restaurant details:', err);
      }
    });
  }

  loadReservations(date: Date) {
    this.reservationService.getReservations(this.restaurantId, undefined, date)
      .subscribe({
        next: (data) => {
          this.reservations = data;
          // Group reservations by hour
          this.groupReservationsByHour();
          // Update table colors based on reservations
          this.updateTableColors();
        },
        error: (err) => {
          console.error('Error fetching reservations:', err);
        }
      });
  }
  
  groupReservationsByHour() {
    // Reset hourly reservations
    this.hourlyReservations = {};
    this.hours = [];
    this.pastHours = [];
    this.futureHours = [];
    
    // Get current hour for highlighting
    const now = new Date();
    this.currentHour = now.getHours().toString().padStart(2, '0');
    
    // Default opening and closing hours if restaurant data is not available
    let openHour = 10; // Default opening hour (10 AM)
    let closeHour = 22; // Default closing hour (10 PM)
    
    // Parse opening hours if available
    if (this.restaurant && this.restaurant.openingHours) {
      const openingHoursMatch = this.restaurant.openingHours.match(/(\d+):(\d+)-(\d+):(\d+)/);
      if (openingHoursMatch && openingHoursMatch.length === 5) {
        openHour = parseInt(openingHoursMatch[1]);
        closeHour = parseInt(openingHoursMatch[3]);
      }
    }
    
    // Create slots only from opening hour to one hour before closing
    for (let i = openHour; i < closeHour -1; i++) {
      const hour = i.toString().padStart(2, '0');
      this.hourlyReservations[hour] = [];
      this.hours.push(hour);
      
      // Separate past and future hours
      if (parseInt(hour) < parseInt(this.currentHour)) {
        this.pastHours.push(hour);
      } else {
        this.futureHours.push(hour);
      }
    }
    
    // Set filtered hours to be the same as hours since we're already filtering
    this.filteredHours = [...this.hours];
    
    // Group reservations by hour
    this.reservations.forEach(reservation => {
      const reservationDate = new Date(reservation.reservationTime);
      const hour = reservationDate.getHours().toString().padStart(2, '0');
      
      if (!this.hourlyReservations[hour]) {
        this.hourlyReservations[hour] = [];
      }
      
      this.hourlyReservations[hour].push(reservation);
    });
  }
  
  filterHoursByOpeningHours() {
    if (!this.restaurant || !this.restaurant.openingHours) return;
    
    // Parse opening hours (assuming format like "10:00-22:00")
    const openingHoursMatch = this.restaurant.openingHours.match(/(\d+):(\d+)-(\d+):(\d+)/);
    
    if (openingHoursMatch && openingHoursMatch.length === 5) {
      const openHour = parseInt(openingHoursMatch[1]);
      const closeHour = parseInt(openingHoursMatch[3]);
      
      // Filter hours to only include those within opening hours
      this.filteredHours = this.hours.filter(hour => {
        const hourNum = parseInt(hour);
        return hourNum >= openHour && hourNum < closeHour;
      });
      
      // Update past and future hours based on filtered hours
      this.pastHours = this.pastHours.filter(hour => this.filteredHours.includes(hour));
      this.futureHours = this.futureHours.filter(hour => this.filteredHours.includes(hour));
    } else {
      // If parsing fails, use all hours
      this.filteredHours = this.hours;
    }
  }
  
  togglePastHours() {
    this.showPastHours = !this.showPastHours;
  }
  
  onDateSelected(event: MatDatepickerInputEvent<Date>) {
    if (event.value) {
      this.selectedDate = event.value;
      this.loadReservations(this.selectedDate);
    }
  }
  
  formatHour(hour: string): string {
    return `${hour}:00`;
  }
  
  isCurrentHour(hour: string): boolean {
    return hour === this.currentHour;
  }
  
  isPastHour(hour: string): boolean {
    return parseInt(hour) < parseInt(this.currentHour);
  }
  
  getHourClass(hour: string): string {
    if (this.isCurrentHour(hour)) {
      return 'current-hour';
    }
    if (this.isPastHour(hour)) {
      return 'past-hour';
    }
    return '';
  }
  
  getTableName(tableId: string): string {
    const table = this.tables.find(t => t.id === tableId);
    return table ? (table.name || `Table ${table.id}`) : 'Unknown';
  }

  // Cancel a reservation
  cancelReservation(reservation: Reservation) {
    // Check if the reservation is already checked in (completed)
    if (reservation.status === 'completed') {
      alert('Cannot cancel a reservation that has already been checked in.');
      return;
    }

    if (confirm(`Are you sure you want to cancel this reservation?`)) {
      // Update the reservation status to cancelled
      this.reservationService.updateReservation(
        this.restaurantId,
        reservation.id,
        { status: 'cancelled' }
      ).subscribe({
        next: () => {
          // Update the local reservation status
          reservation.status = 'cancelled';
          // Update table colors
          this.updateTableColors();
        },
        error: (err) => {
          console.error('Error cancelling reservation:', err);
        }
      });
    }
  }

  // Get available tables for a reservation
  getAvailableTablesForReservation(reservation: Reservation): Table[] {
    // Get tables that can accommodate the party size
    return this.tables.filter(table => 
      // Include the currently assigned table (if any) or available tables
      (table.id === reservation.tableId || table.status === 'available') && 
      table.minCapacity <= reservation.numberOfGuests && 
      table.maxCapacity >= reservation.numberOfGuests
    );
  }

  // Handle table selection change
  onTableSelectionChange(reservation: Reservation, tableId: string | null) {
    // If the selected table is the same as the current one, do nothing
    if (tableId === reservation.tableId) {
      return;
    }

    // Check if the reservation is already checked in (completed)
    if (reservation.status === 'completed') {
      alert('Cannot change table for a reservation that has already been checked in.');
      return;
    }

    // Update the reservation with the new table
    this.reservationService.updateReservation(
      this.restaurantId,
      reservation.id,
      { 
        tableId: tableId,
        status: 'confirmed' 
      }
    ).subscribe({
      next: () => {
        // Update the local reservation
        reservation.tableId = tableId;
        reservation.status = 'confirmed';
        // Update table colors
        this.updateTableColors();
      },
      error: (err) => {
        console.error('Error assigning table:', err);
      }
    });
  }

  // Modify a reservation
  modifyReservation(reservation: Reservation) {
    // Check if the reservation is already checked in (completed)
    if (reservation.status === 'completed') {
      alert('Cannot modify a reservation that has already been checked in.');
      return;
    }
    
    // Open a dialog to edit the reservation details
    const dialogRef = this.dialog.open(NewReservationComponent, {
      width: '500px',
      data: {
        restaurant: this.restaurant,
        restaurantId: this.restaurantId,
        reservation: reservation
      }
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        // Reload reservations to get the updated data
        this.loadReservations(this.selectedDate);
      }
    });
  }

  // Check in a reservation
  checkInReservation(reservation: Reservation) {
    // Check if the reservation is already checked in
    if (reservation.status === 'completed') {
      alert('This reservation has already been checked in.');
      return;
    }
    
    // Check if the reservation is cancelled
    if (reservation.status === 'cancelled') {
      alert('Cannot check in a cancelled reservation.');
      return;
    }

    this.reservationService.updateReservation(
      this.restaurantId,
      reservation.id,
      { status: 'completed' }
    ).subscribe({
      next: () => {
        // Update the local reservation status
        reservation.status = 'completed';
        
        // Find the associated table and update its status to occupied
        const table = this.tables.find(t => t.id === reservation.tableId);
        if (table) {
          table.status = 'occupied';
          
          // Update the table status on the server
          this.tablesService.updateTable(
            this.restaurantId,
            table.id,
            { status: 'occupied' }
          ).subscribe({
            next: () => {
              // Update table colors
              this.updateTableColors();
              alert('Reservation checked in successfully!');
            },
            error: (err) => {
              console.error('Error updating table status:', err);
            }
          });
        } else {
          alert('Reservation checked in successfully!');
          this.updateTableColors();
        }
      },
      error: (err) => {
        console.error('Error checking in reservation:', err);
        alert('Error checking in reservation. Please try again.');
      }
    });
  }

  updateTableColors() {
    // Get current time
    const now = new Date();
    const currentHour = now.getHours().toString().padStart(2, '0');
    const currentHourReservations = this.hourlyReservations[currentHour] || [];
    
    // Find all tables that have reservations for the current hour
    const reservedTableIds = currentHourReservations
      .filter(r => r.status !== 'cancelled')
      .map(r => r.tableId);
    
    // Update the local table status
    this.tables.forEach(table => {
      // First check if table is occupied by a checked-in reservation
      const isOccupied = this.reservations
        .filter(r => r.status === 'completed' && r.tableId === table.id)
        .length > 0;
      
      if (isOccupied) {
        table.status = 'occupied';
      } else if (reservedTableIds.includes(table.id)) {
        // Table has a current hour reservation
        table.status = 'occupied';
      } else {
        // Check if table is reserved within the next 15 minutes
        const isReservedSoon = this.reservations
          .filter(r => r.status === 'confirmed' && r.tableId === table.id)
          .some(r => {
            const reservationTime = new Date(r.reservationTime);
            const timeDiff = reservationTime.getTime() - now.getTime();
            const minutesDiff = timeDiff / (1000 * 60);
            return minutesDiff <= 15 && minutesDiff > 0;
          });
        
        // Check if table is reserved later today
        const isReservedLater = this.reservations
          .filter(r => r.status !== 'cancelled' && r.tableId === table.id)
          .some(r => {
            const reservationTime = new Date(r.reservationTime);
            return reservationTime > now;
          });
        
        if (isReservedSoon) {
          table.status = 'reserved';
        } else if (isReservedLater) {
          table.status = 'reserved';
        } else {
          table.status = 'available';
        }
      }
    });
    
    // Redraw the tables with updated colors
    if (this.layer) {
      this.layer.destroyChildren();
      this.drawTables();
    }
  }

  ngAfterViewInit() {
    console.log('floorPlanHost in ngAfterViewInit:', this.floorPlanHost);
    this.viewInitialized = true;
    this.initTableView();
  }

  loadTables() {
    this.loading = true;
    this.sub = this.tablesService.getTables(this.restaurantId).subscribe({
      next: (data) => {
        this.tables = data || [];
        this.loading = false;
        this.tablesInitialized = true;
        console.log(this.tables);
        // Initialize the Konva stage + layer after we have the data
        this.initTableView();
      },
      error: (err) => {
        console.error('Error fetching tables:', err);
        this.loading = false;
      }
    });
  }

  initTableView() {
    if(this.viewInitialized && this.tablesInitialized) {
      this.initStage();
      this.drawTables();
    }
  }

  initStage() {
    const containerEl = this.floorPlanHost.nativeElement;
    const containerWidth = containerEl.offsetWidth;
    const containerHeight = containerEl.offsetHeight;

    // Create a Konva stage
    this.stage = new Konva.Stage({
      container: containerEl,
      width: containerWidth,
      height: containerHeight
    });

    // Create a Konva layer
    this.layer = new Konva.Layer();
    this.stage.add(this.layer);
  }

  drawTables() {
    this.tables.forEach((table) => {
      // Decide shape based on capacity or any logic you prefer
      // Let's say if maxCapacity <= 3 => circle, else rectangle
      if (table.maxCapacity <= 4) {
        this.createCircleTable(table);
      } else {
        this.createRectTable(table);
      }
    });

    // Draw the layer after adding all shapes
    this.layer.batchDraw();
  }

  // Helper function to expand a rectangle by a given margin
  private expandRect(
    rect: { x: number; y: number; width: number; height: number },
    margin: number
  ): { x: number; y: number; width: number; height: number } {
    return {
      x: rect.x - margin,
      y: rect.y - margin,
      width: rect.width + margin * 2,
      height: rect.height + margin * 2
    };
  }

  createCircleTable(table: Table) {
    const group = new Konva.Group({
      x: table.x ?? 100,
      y: table.y ?? 100,
      draggable: true
    });

    const radius = 30;
    const circle = new Konva.Circle({
      x: 0,
      y: 0,
      radius: radius,
      fill: this.getTableColor(table)
    });

    const textContent = `${table.name || ('Table ' + table.id)}\n(${table.minCapacity}-${table.maxCapacity})`;
    const text = new Konva.Text({
      x: 0,
      y: 0,
      text: textContent,
      fontSize: 12,
      fontFamily: 'Calibri',
      fill: 'white',
      align: 'center'
    });
    text.offsetX(text.width() / 2);
    text.offsetY(text.height() / 2);

    group.add(circle);
    group.add(text);

    // Add drag constraints as before...
    const boundaryMargin = 10;
    group.dragBoundFunc((pos: { x: number; y: number }) => {
      const temp = group.clone({ x: pos.x, y: pos.y });
      const potentialRect = temp.getClientRect({ skipTransform: false });
      const expandedRect = this.expandRect(potentialRect, boundaryMargin);
      let collision = false;
      const groups = this.layer.find('Group');
      groups.forEach(otherGroup => {
        if (otherGroup === group) return;
        const otherRect = otherGroup.getClientRect({ skipTransform: false });
        const expandedOtherRect = this.expandRect(otherRect, boundaryMargin);
        if (this.rectsIntersect(expandedRect, expandedOtherRect)) {
          collision = true;
        }
      });
      return collision ? { x: group.x(), y: group.y() } : pos;
    });

    group.on('dragend', (evt) => {
      this.onShapeDragEnd(table, evt);
    });

    // **Double click to open the edit dialog.**
    group.on('dblclick', () => {
      this.editTable(table, group);
    });

    this.layer.add(group);
  }

  createRectTable(table: Table) {
    const group = new Konva.Group({
      x: table.x ?? 100,
      y: table.y ?? 100,
      draggable: true
    });

    const rectWidth = 60;
    const rectHeight = 60;
    const rect = new Konva.Rect({
      x: -rectWidth / 2,
      y: -rectHeight / 2,
      width: rectWidth,
      height: rectHeight,
      fill: this.getTableColor(table),
      cornerRadius: 10
    });

    const textContent = `${table.name || ('Table ' + table.id)}\n(${table.minCapacity}-${table.maxCapacity})`;
    const text = new Konva.Text({
      x: 0,
      y: 0,
      text: textContent,
      fontSize: 12,
      fontFamily: 'Calibri',
      fill: 'white',
      align: 'center'
    });
    text.offsetX(text.width() / 2);
    text.offsetY(text.height() / 2);

    group.add(rect);
    group.add(text);

    const boundaryMargin = 10;
    group.dragBoundFunc((pos: { x: number; y: number }) => {
      const temp = group.clone({ x: pos.x, y: pos.y });
      const potentialRect = temp.getClientRect({ skipTransform: false });
      const expandedRect = this.expandRect(potentialRect, boundaryMargin);
      let collision = false;
      const groups = this.layer.find('Group');
      groups.forEach(otherGroup => {
        if (otherGroup === group) return;
        const otherRect = otherGroup.getClientRect({ skipTransform: false });
        const expandedOtherRect = this.expandRect(otherRect, boundaryMargin);
        if (this.rectsIntersect(expandedRect, expandedOtherRect)) {
          collision = true;
        }
      });
      return collision ? { x: group.x(), y: group.y() } : pos;
    });

    group.on('dragend', (evt) => {
      this.onShapeDragEnd(table, evt);
    });

    // **Double click to edit table details.**
    group.on('dblclick', () => {
      this.editTable(table, group);
    });

    this.layer.add(group);
  }


  // Helper function to check if two rectangles intersect.
  private rectsIntersect(
    r1: { x: number; y: number; width: number; height: number },
    r2: { x: number; y: number; width: number; height: number }
  ): boolean {
    return !(
      r2.x > r1.x + r1.width ||
      r2.x + r2.width < r1.x ||
      r2.y > r1.y + r1.height ||
      r2.y + r2.height < r1.y
    );
  }


  onShapeDragEnd(table: Table, evt: Konva.KonvaEventObject<DragEvent>) {
    const shape = evt.target;
    const newX = shape.x();
    const newY = shape.y();

    // Update local data
    table.x = newX;
    table.y = newY;

    // Optionally call your PUT /restaurants/:restaurantId/tables/:tableId
    this.tablesService.updateTable(this.restaurantId, table.id, {
      x: newX,
      y: newY
    }).subscribe({
      next: updated => {
        console.log('Updated table location', updated);
      },
      error: err => {
        console.error('Failed to update table location', err);
      }
    });
  }

  // Example: color table based on status
  getTableColor(table: Table): string {
    switch (table.status) {
      case 'available':
        return '#4caf50'; // green
      case 'occupied':
        return '#f44336'; // red
      case 'reserved':
        return '#ff9800'; // orange
      default:
        return '#9e9e9e'; // gray
    }
  }

  // Add a new table with default values and render it.
  addNewTable() {
    const dialogRef = this.dialog.open(NewTableComponent, {
      width: '300px',
      data: { name: '', minCapacity: 1, maxCapacity: 4 } as NewTable
    });

    dialogRef.afterClosed().subscribe((result: NewTable | undefined) => {
      if (result) {
        // Create new table payload with user-supplied details.
        const newTablePayload = {
          name: result.name,
          minCapacity: result.minCapacity,
          maxCapacity: result.maxCapacity
        };

        // Push the new table to the backend.
        this.tablesService.addTable(this.restaurantId, newTablePayload).subscribe({
          next: (createdTable: Table) => {
            // Optionally update local tables array.
            this.tables.push(createdTable);
            // Add the new table to the canvas.
            if (createdTable.maxCapacity <= 3) {
              this.createCircleTable(createdTable);
            } else {
              this.createRectTable(createdTable);
            }
          },
          error: (err) => {
            console.error('Failed to add table', err);
            alert('Failed to add table');
          }
        });
      }
    });

  }

  editTable(table: Table, group: Konva.Group) {
  // Open the dialog prepopulated with current table data.
  const dialogRef = this.dialog.open(NewTableComponent, {
    width: '300px',
    data: { name: table.name, minCapacity: table.minCapacity, maxCapacity: table.maxCapacity , isEdit: true}
  });

  dialogRef.afterClosed().subscribe((result: any) => {
    if (result) {
      if (result.delete) {
        // Delete action: call backend and remove table from canvas.
        this.tablesService.deleteTable(this.restaurantId, table.id).subscribe({
          next: () => {
            // Remove table from local list.
            this.tables = this.tables.filter(t => t.id !== table.id);
            // Remove the corresponding Konva group.
            group.destroy();
            this.layer.batchDraw();
          },
          error: err => {
            console.error('Failed to delete table', err);
            alert('Failed to delete ' + table.name);
          }
        });
      } else {
        // Save action: update table details.
        table.name = result.name;
        table.minCapacity = result.minCapacity;
        table.maxCapacity = result.maxCapacity;
        this.tablesService.updateTable(this.restaurantId, table.id, table).subscribe({
          next: () => {
            // Update the text element in the group.
            const textShape = group.findOne('Text') as Konva.Text;
            if (textShape) {
              const newText = `${table.name || ('Table ' + table.id)}\n(${table.minCapacity}-${table.maxCapacity})`;
              textShape.text(newText);
              textShape.offsetX(textShape.width() / 2);
              textShape.offsetY(textShape.height() / 2);
              this.layer.batchDraw();
            }
          },
          error: err => {
            console.error('Failed to update table', err);
            alert('Failed to update ' + table.name);
          }
        });
      }
    }
  });
}


  toggleSidebar() {
    this.sidebarHidden = !this.sidebarHidden;
    
    // Give the DOM time to update before resizing the stage
    setTimeout(() => {
      if (this.stage) {
        // Resize the stage to fit the new container size
        const containerEl = this.floorPlanHost.nativeElement;
        this.stage.width(containerEl.offsetWidth);
        this.stage.height(containerEl.offsetHeight);
        this.layer.batchDraw();
      }
    }, 300); // Match the CSS transition time (0.3s)
  }

  ngOnDestroy() {
    if (this.sub) {
      this.sub.unsubscribe();
    }
    // Konva typically cleans up if stage is removed, but we can do:
    if (this.stage) {
      this.stage.destroy();
    }
  }
}
