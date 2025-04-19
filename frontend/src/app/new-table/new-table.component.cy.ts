import { NewTableComponent } from './new-table.component';
import { FormBuilder, ReactiveFormsModule } from '@angular/forms';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';
import { NewTable } from '../models/table.model';

describe('NewTableComponent', () => {
  // Test data for different scenarios
  const createModeData: NewTable = {
    minCapacity: 0,
    maxCapacity: 0
  };

  const editModeData: NewTable = {
    name: 'Table 1',
    minCapacity: 2,
    maxCapacity: 6,
    isEdit: true
  };

  // Define mock objects
  let mockDialogRef;

  describe('Create mode', () => {
    beforeEach(() => {
      // Initialize mock objects with stubs
      mockDialogRef = {
        close: cy.stub().as('dialogClose')
      };
      
      // Mount component with all necessary dependencies
      cy.mount(NewTableComponent, {
        imports: [
          ReactiveFormsModule,
          MatCardModule,
          MatFormFieldModule,
          MatInputModule,
          MatButtonModule,
          NoopAnimationsModule
        ],
        providers: [
          FormBuilder,
          { provide: MatDialogRef, useValue: mockDialogRef },
          { provide: MAT_DIALOG_DATA, useValue: createModeData }
        ]
      });
    });

    it('should create', () => {
      cy.get('mat-card').should('exist');
    });

    it('should initialize form with default values in create mode', () => {
      cy.get('[formControlName="name"]').should('have.value', '');
      cy.get('[formControlName="minCapacity"]').should('have.value', '1');
      cy.get('[formControlName="maxCapacity"]').should('have.value', '4');
    });

    it('should validate required fields', () => {
      // Name is required
      cy.get('[formControlName="name"]').clear();
      cy.get('button[color="primary"]').should('be.disabled');
      cy.get('[formControlName="name"]').type('Table 1');
      
      // Min capacity is required and must be at least 1
      cy.get('[formControlName="minCapacity"]').clear();
      cy.get('button[color="primary"]').should('be.disabled');
      cy.get('[formControlName="minCapacity"]').clear().type('0');
      cy.get('button[color="primary"]').should('be.disabled');
      cy.get('[formControlName="minCapacity"]').clear().type('1');
      
      // Max capacity is required and must be at least 1
      cy.get('[formControlName="maxCapacity"]').clear();
      cy.get('button[color="primary"]').should('be.disabled');
      cy.get('[formControlName="maxCapacity"]').clear().type('0');
      cy.get('button[color="primary"]').should('be.disabled');
      cy.get('[formControlName="maxCapacity"]').clear().type('4');
      
      // Form should be valid now
      cy.get('button[color="primary"]').should('not.be.disabled');
    });

    it('should close dialog with form value when save is clicked', () => {
      // Set valid form values
      cy.get('[formControlName="name"]').type('Table 1');
      cy.get('[formControlName="minCapacity"]').clear().type('2');
      cy.get('[formControlName="maxCapacity"]').clear().type('6');

      // Submit form
      cy.get('button[color="primary"]').click();
      
      // Verify dialog was closed with correct values
      cy.get('@dialogClose').should('have.been.calledWith', {
        name: 'Table 1',
        minCapacity: 2,
        maxCapacity: 6
      });
    });

    it('should close dialog without value when cancel is clicked', () => {
      cy.get('button').contains('Cancel').click();
      cy.get('@dialogClose').should('have.been.calledWith');
    });
  });

  describe('Edit mode', () => {
    beforeEach(() => {
      // Initialize mock objects with stubs
      mockDialogRef = {
        close: cy.stub().as('dialogClose')
      };
      
      // Mount component with all necessary dependencies
      cy.mount(NewTableComponent, {
        imports: [
          ReactiveFormsModule,
          MatCardModule,
          MatFormFieldModule,
          MatInputModule,
          MatButtonModule,
          NoopAnimationsModule
        ],
        providers: [
          FormBuilder,
          { provide: MatDialogRef, useValue: mockDialogRef },
          { provide: MAT_DIALOG_DATA, useValue: editModeData }
        ]
      });
    });

    it('should initialize form with provided values in edit mode', () => {
      cy.get('[formControlName="name"]').should('have.value', 'Table 1');
      cy.get('[formControlName="minCapacity"]').should('have.value', '2');
      cy.get('[formControlName="maxCapacity"]').should('have.value', '6');
    });

    it('should show delete button in edit mode', () => {
      cy.get('button[color="warn"]').should('exist');
      cy.get('button[color="warn"]').should('contain.text', 'Delete');
    });

    it('should close dialog with delete flag when delete is clicked', () => {
      cy.get('button[color="warn"]').click();
      cy.get('@dialogClose').should('have.been.calledWith', { delete: true });
    });
  });
});