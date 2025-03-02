// import { Component, Inject } from '@angular/core';
// import { CommonModule } from '@angular/common';
// import { FormsModule, ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
// import { MatCardModule } from '@angular/material/card';
// import { MatFormFieldModule } from '@angular/material/form-field';
// import { MatInputModule } from '@angular/material/input';
// import { MatButtonModule } from '@angular/material/button';
// import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
// import { RestaurantService } from '../services/restaurant.service';
// import { RestaurantCreate } from '../models/restaurant.model';

// @Component({
//   selector: 'app-new-restaurant',
//   standalone: true,

//   templateUrl: './new-restaurant.component.html',
//   styleUrl: './new-restaurant.component.css',
//   imports: [
//     CommonModule,
//     FormsModule,
//     ReactiveFormsModule,
//     MatCardModule,
//     MatFormFieldModule,
//     MatInputModule,
//     MatButtonModule
//   ]
// })
// export class NewRestaurantComponent {
//   restaurantForm: FormGroup;
//   selectedFile: File | null = null;

//   constructor(
//     private fb: FormBuilder,
//     private restaurantService: RestaurantService,
//     public dialogRef: MatDialogRef<NewRestaurantComponent>,
//     @Inject(MAT_DIALOG_DATA) public data: any
//   ) {
//     this.restaurantForm = this.fb.group({
//       name: ['', Validators.required],
//       location: ['', Validators.required],
//       description: [''],
//       phone: [''],
//       openingHours: ['']
//     });
//   }

//   onFileSelected(event: any): void {
//     const file: File = event.target.files?.[0];
//     if (file) {
//       this.selectedFile = file;
//       console.log('Selected file:', file);
//     }
//   }

//   onSubmit(): void {
//     if (this.restaurantForm.invalid) {
//       return;
//     }

//     const restaurant: RestaurantCreate = {
//       name: this.restaurantForm.value.name,
//       location: this.restaurantForm.value.location,
//       description: this.restaurantForm.value.description,
//       phone: this.restaurantForm.value.phone,
//       openingHours: this.restaurantForm.value.openingHours,
//       img: this.selectedFile || undefined
//     };

//     // Perform the POST request
//     this.restaurantService.createRestaurant(restaurant).subscribe({
//       next: (created) => {
//         console.log('Restaurant created successfully:', created);
//         alert(`Restaurant "${created?.name}" created successfully!`);
//         this.restaurantForm.reset();
//         this.selectedFile = null;
//         this.dialogRef.close(true); // Close the dialog and pass true to indicate success
//       },
//       error: (err) => {
//         console.error('Failed to create restaurant:', err);
//         alert('Error creating the restaurant.');
//       }
//     });
//   }
// }

import { Component, Inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { RestaurantService } from '../services/restaurant.service';
import { RestaurantCreate } from '../models/restaurant.model';

@Component({
  selector: 'app-new-restaurant',
  standalone: true,

  templateUrl: './new-restaurant.component.html',
  styleUrl: './new-restaurant.component.css',
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    MatCardModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule
  ]
})
export class NewRestaurantComponent {
  restaurantForm: FormGroup;
  selectedFile: File | null = null;

  constructor(
    private fb: FormBuilder,
    private restaurantService: RestaurantService,
    public dialogRef: MatDialogRef<NewRestaurantComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any
  ) {
    this.restaurantForm = this.fb.group({
      name: ['', Validators.required],
      location: ['', Validators.required],
      description: [''],
      phone: [''],
      openingHours: ['']
    });
  }

  onFileSelected(event: any): void {
    const file: File = event.target.files?.[0];
    if (file) {
      this.selectedFile = file;
      console.log('Selected file:', file);
    }
  }

  onSubmit(): void {
    if (this.restaurantForm.invalid) {
      return;
    }

    // Prepare the restaurant data
    const restaurant: RestaurantCreate = {
      name: this.restaurantForm.value.name,
      location: this.restaurantForm.value.location,
      description: this.restaurantForm.value.description,
      phone: this.restaurantForm.value.phone,
      openingHours: this.restaurantForm.value.openingHours,
      // If a file is selected, use its name as the img field (replace with the actual uploaded file name after upload)
      img: this.selectedFile ? this.selectedFile.name : undefined
    };

    // Perform the POST request
    this.restaurantService.createRestaurant(restaurant).subscribe({
      next: (created) => {
        console.log('Restaurant created successfully:', created);
        alert(`Restaurant "${created?.name}" created successfully!`);
        this.restaurantForm.reset();
        this.selectedFile = null;
        this.dialogRef.close(true); // Close the dialog and pass true to indicate success
      },
      error: (err) => {
        console.error('Failed to create restaurant:', err);
        alert('Error creating the restaurant.');
      }
    });
  }
}


