import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import { RestaurantCreate } from '../models/restaurant.model';

@Injectable({
  providedIn: 'root'
})
export class RestaurantService {
  private mockRestaurants: any[] = [];

  constructor() { }

  createRestaurant(data: RestaurantCreate): Observable<any> {
    // Mock: simply push to local array and return it
    this.mockRestaurants.push(data);
    return of(data);
  }
}
