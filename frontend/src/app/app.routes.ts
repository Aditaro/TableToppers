import { Routes } from '@angular/router';
import { NewRestaurantComponent } from './new-restaurant/new-restaurant.component';
import { RestaurantsComponent } from './restaurants/restaurants.component';

export const routes: Routes = [
    { path: '', redirectTo: 'restaurants', pathMatch: 'full' },
    { path: 'restaurants', component: RestaurantsComponent }
];
