import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA, MatDialogModule } from '@angular/material/dialog'; // Import MatDialogModule
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms'; // Import ReactiveFormsModule
import { MatFormFieldModule } from '@angular/material/form-field'; // Import MatFormFieldModule
import { MatInputModule } from '@angular/material/input'; // Import MatInputModule
import { MatButtonModule } from '@angular/material/button'; // Import MatButtonModule
import { CommonModule } from '@angular/common'; // Import CommonModule

export interface AddWaitlistDialogData {
  restaurantId: string;
}

export interface AddWaitlistDialogResult {
  name: string;
  partySize: number;
  phoneNumber: string;
}

@Component({
  selector: 'app-add-waitlist-dialog',
  standalone: true, // Make component standalone
  imports: [ // Add necessary imports
    CommonModule,
    ReactiveFormsModule,
    MatDialogModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule
  ],
  templateUrl: './add-waitlist-dialog.component.html',
  styleUrls: ['./add-waitlist-dialog.component.css']
})
export class AddWaitlistDialogComponent {
  waitlistForm: FormGroup;

  constructor(
    public dialogRef: MatDialogRef<AddWaitlistDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: AddWaitlistDialogData,
    private fb: FormBuilder
  ) {
    this.waitlistForm = this.fb.group({
      name: ['', Validators.required],
      partySize: [1, [Validators.required, Validators.min(1)]],
      phoneNumber: ['']
    });
  }

  onCancel(): void {
    this.dialogRef.close();
  }

  onSubmit(): void {
    if (this.waitlistForm.valid) {
      const result: AddWaitlistDialogResult = {
        name: this.waitlistForm.value.name,
        partySize: this.waitlistForm.value.partySize,
        phoneNumber: this.waitlistForm.value.phoneNumber
      };
      this.dialogRef.close(result);
    }
  }
}