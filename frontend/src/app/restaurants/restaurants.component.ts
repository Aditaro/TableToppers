// import { Component } from '@angular/core';
// import { CommonModule } from '@angular/common';
// import { FormsModule } from '@angular/forms';
// import { MatCardModule } from '@angular/material/card';
// import { MatFormFieldModule } from '@angular/material/form-field';
// import { MatInputModule } from '@angular/material/input';
// import { MatButtonModule } from '@angular/material/button';
// import { Restaurant } from 'src/app/models/restaurant.model';
// import { RestaurantService } from 'src/app/services/restaurant.service';
// import { NewRestaurantComponent } from '../new-restaurant/new-restaurant.component';
// import { MatDialog } from '@angular/material/dialog';

// @Component({
//   selector: 'app-restaurants',
//   standalone: true,

//   templateUrl: './restaurants.component.html',
//   styleUrl: './restaurants.component.css',
//     imports: [
//       CommonModule,
//       FormsModule,
//       MatCardModule,
//       MatFormFieldModule,
//       MatInputModule,
//       MatButtonModule
//     ]
// })
// export class RestaurantsComponent {
//   cityFilter = '';
//   nameFilter = '';
//   restaurants: Restaurant[] = [];

//   constructor(private restaurantService: RestaurantService, public dialog: MatDialog) { }

//   ngOnInit(): void {
//     this.fetchRestaurants();
//   }

//   fetchRestaurants(): void {
//     this.restaurantService.getRestaurants(this.cityFilter, this.nameFilter)
//       .subscribe(data => {
//         this.restaurants = data;
//       });
//   }

//   onSearch(): void {
//     this.fetchRestaurants();
//   }

//   openNewRestaurantDialog(): void {
//     const dialogRef = this.dialog.open(NewRestaurantComponent, {
//       width: '600px',
//       data: {}
//     });

//     dialogRef.afterClosed().subscribe(result => {
//       if (result) {
//         this.fetchRestaurants(); // Refresh the list if a new restaurant was created
//       }
//     });
//   }
// }

import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { Restaurant } from 'src/app/models/restaurant.model';
import { RestaurantService } from 'src/app/services/restaurant.service';
import { NewRestaurantComponent } from '../new-restaurant/new-restaurant.component';
import { MatDialog } from '@angular/material/dialog';

@Component({
  selector: 'app-restaurants',
  standalone: true,

  templateUrl: './restaurants.component.html',
  styleUrl: './restaurants.component.css',
  imports: [
    CommonModule,
    FormsModule,
    MatCardModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule
  ]
})
export class RestaurantsComponent {
  cityFilter = '';
  nameFilter = '';
  restaurants: Restaurant[] = [];

  constructor(private restaurantService: RestaurantService, public dialog: MatDialog) { }

  ngOnInit(): void {
    this.fetchRestaurants();
  }

  fetchRestaurants(): void {
    // Pass the filters to the getRestaurants method
    this.restaurantService.getRestaurants(this.cityFilter, this.nameFilter)
      .subscribe(data => {
        this.restaurants = data;
      });
  }

  onSearch(): void {
    this.fetchRestaurants();  // Re-fetch restaurants based on updated filters
  }

  openNewRestaurantDialog(): void {
    const dialogRef = this.dialog.open(NewRestaurantComponent, {
      width: '600px',
      data: {}
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        this.fetchRestaurants(); // Refresh the list if a new restaurant was created
      }
    });
  }
}

