import { NgModule } from '@angular/core';
import { NewRestaurantComponent } from './new-restaurant/new-restaurant.component';
import { RestaurantsComponent } from './restaurants/restaurants.component';
import { TablesComponent } from './tables/tables.component';
import {HashLocationStrategy, LocationStrategy, NgIf} from '@angular/common';
import { NewTableComponent } from './new-table/new-table.component';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {MatCard, MatCardContent, MatCardTitle} from '@angular/material/card';
import {MatFormField, MatLabel} from '@angular/material/form-field';
import {MatInput} from '@angular/material/input';
import {MatButton} from '@angular/material/button';
import {MatDialogActions} from '@angular/material/dialog';
import { NewReservationComponent } from './new-reservation/new-reservation.component';
import { ReservationsComponent } from './reservations/reservations.component';

@NgModule({
  providers: [
    { provide: LocationStrategy, useClass: HashLocationStrategy }
  ],
  declarations: [
    NewTableComponent
  ],
  imports: [
    NewRestaurantComponent,
    RestaurantsComponent,
    NewReservationComponent,
    ReservationsComponent,
    TablesComponent,
    FormsModule,
    MatCard,
    MatCardContent,
    MatCardTitle,
    MatFormField,
    MatInput,
    MatLabel,
    ReactiveFormsModule,
    MatButton,
    MatDialogActions,
    NgIf
  ],
})
export class TranslocoRootModule {}
