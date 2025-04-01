import { Routes } from '@angular/router';
import { HomeComponent } from './home/home.component';
import { LoginComponent } from './login/login.component';
import { RegisterComponent } from './register/register.component';
import { RestaurantsComponent } from './restaurants/restaurants.component';
import { NewRestaurantComponent } from './new-restaurant/new-restaurant.component';
import { NewReservationComponent } from './new-reservation/new-reservation.component';
import { CustomerInfoComponent } from './customer-info/customer-info.component';
import { BusinessPortalComponent } from './business-portal/business-portal.component';

export const routes: Routes = [
  { path: '', component: HomeComponent },
  { path: 'login', component: LoginComponent },
  { path: 'register', component: RegisterComponent },
  { path: 'restaurants', component: RestaurantsComponent },
  { path: 'new-restaurant', component: NewRestaurantComponent },
  { path: 'new-reservation', component: NewReservationComponent },
  { path: 'customer-info', component: CustomerInfoComponent },
  { path: 'business-portal', component: BusinessPortalComponent },
  { path: '**', redirectTo: '' }
];
