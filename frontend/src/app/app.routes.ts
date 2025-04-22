import { Routes } from '@angular/router';
import { AuthGuard } from './auth.guard';
import { RoleGuard } from './role.guard';
import { HomeComponent } from './home/home.component';
import { LoginComponent } from './login/login.component';
import { RegisterComponent } from './register/register.component';
import { RestaurantsComponent } from './restaurants/restaurants.component';
import { NewRestaurantComponent } from './new-restaurant/new-restaurant.component';
import { NewReservationComponent } from './new-reservation/new-reservation.component';
import { CustomerInfoComponent } from './customer-info/customer-info.component';
import { BusinessPortalComponent } from './business-portal/business-portal.component';
import { TablesComponent } from './tables/tables.component';
import { ReservationsComponent } from './reservations/reservations.component';
import { SubscriptionComponent } from './subscription/subscription.component';

export const routes: Routes = [
  { path: '', component: HomeComponent },
  { path: 'login', component: LoginComponent },
  { path: 'register', component: RegisterComponent },
  { path: 'subscription', component: SubscriptionComponent },
  { path: 'restaurants', component: RestaurantsComponent },
  { path: 'new-restaurant', component: NewRestaurantComponent, canActivate: [AuthGuard, RoleGuard], data: { expectedRoles: ['manager', 'admin'] } },
  // { path: 'new-reservation', component: NewReservationComponent },
  { path: 'customer-info', component: CustomerInfoComponent, canActivate: [AuthGuard] },
  { path: 'business-portal', component: BusinessPortalComponent },
  { path: 'restaurants/:restaurantId/tables', component: TablesComponent, canActivate: [AuthGuard, RoleGuard], data: { expectedRoles: ['manager', 'admin'] } },
  { path: 'restaurants/:restaurantId/reserve', component: ReservationsComponent, canActivate: [AuthGuard, RoleGuard], data: { expectedRoles: ['customer'] } },
  { path:'restaurants/:restaurantId/reserve/new', component: NewReservationComponent, canActivate: [AuthGuard, RoleGuard], data: { expectedRoles: ['customer'] } },
  { path: '**', redirectTo: '' }
];
