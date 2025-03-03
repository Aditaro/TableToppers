import {Component, Inject} from '@angular/core';
import {FormBuilder, FormGroup, Validators} from '@angular/forms';
import {MAT_DIALOG_DATA, MatDialogRef} from '@angular/material/dialog';
import {NewTable} from '../models/table.model';

@Component({
  selector: 'app-new-table',
  standalone: false,
  templateUrl: './new-table.component.html',
  styleUrl: './new-table.component.css'
})
export class NewTableComponent {
  tableForm: FormGroup;

  constructor(
    public dialogRef: MatDialogRef<NewTableComponent>,
    @Inject(MAT_DIALOG_DATA) public data: NewTable,
    private fb: FormBuilder
  ) {
    this.tableForm = this.fb.group({
      name: [data.name || '', Validators.required],
      minCapacity: [data.minCapacity || 1, [Validators.required, Validators.min(1)]],
      maxCapacity: [data.maxCapacity || 4, [Validators.required, Validators.min(1)]]
    });
  }

  onCancel(): void {
    this.dialogRef.close();
  }

  onSave(): void {
    if (this.tableForm.valid) {
      this.dialogRef.close(this.tableForm.value);
    }
  }

  onDelete(): void {
    // Return an object with a "delete" flag.
    this.dialogRef.close({ delete: true });
  }
}
