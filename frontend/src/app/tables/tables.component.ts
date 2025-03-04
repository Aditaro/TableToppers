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

@Component({
  standalone: true,
  selector: 'app-tables',
  imports: [
    CommonModule,
    MatCardModule,
    MatToolbarModule,
    MatProgressSpinnerModule,
    MatIconModule
  ],
  templateUrl: './tables.component.html',
  styleUrl: './tables.component.css',
})
export class TablesComponent implements OnInit, OnDestroy, AfterViewInit {
  @ViewChild('floorPlanHost', { static: true }) floorPlanHost!: ElementRef<HTMLDivElement>;

  restaurantId!: string;
  tables: Table[] = [];
  loading = false;

    // Konva objects
  private stage!: Stage;
  private layer!: Layer;

  private viewInitialized = false;
  private tablesInitialized = false;

  private sub!: Subscription;

  constructor(
    private route: ActivatedRoute,
    private tablesService: TablesService,
    public dialog: MatDialog
  ) {}

  ngOnInit() {
    // 1. Grab restaurantId from route
    this.restaurantId = this.route.snapshot.paramMap.get('restaurantId') || '';

    // 2. Load tables for the restaurant
    this.loadTables();
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
