import { Route } from '@angular/router';
import { RestaurantsComponent } from './restaurants/restaurants.component';
import { HomeComponent } from './home/home.component';
import { LoginComponent } from './login/login.component';
import { RegisterComponent } from './register/register.component';
import {TablesComponent} from './tables/tables.component';
import {NewReservationComponent} from './new-reservation/new-reservation.component';

export const routes: Route[] = [
    { path: 'home', component: HomeComponent },
    { path: 'restaurants', component: RestaurantsComponent },
    { path: 'login', component: LoginComponent },
    { path: 'register', component: RegisterComponent },
    { path: 'restaurants/:restaurantId/tables', component: TablesComponent },
    { path: 'restaurants/:restaurantId/reserve', component: NewReservationComponent },
    { path: '', redirectTo: '/home', pathMatch: 'full' },  // Set default route to home
    { path: '**', redirectTo: '/home' }
];
