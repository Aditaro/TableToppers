import { Routes } from '@angular/router';
import { NewRestaurantComponent } from './new-restaurant/new-restaurant.component';

export const routes: Routes = [
    { path: '', redirectTo: 'new-restaurant', pathMatch: 'full' },
    { path: 'new-restaurant', component: NewRestaurantComponent }
];
