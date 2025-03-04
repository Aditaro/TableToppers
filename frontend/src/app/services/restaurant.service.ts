import { Injectable } from '@angular/core';
import {catchError, Observable, of, throwError} from 'rxjs';
import {NewRestaurant, Restaurant, RestaurantCreate} from '../models/restaurant.model';
import {HttpClient, HttpResponse} from '@angular/common/http';
import {environment} from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class RestaurantService {
  constructor(private http: HttpClient) { }

  getRestaurants(city?: string, name?: string): Observable<Restaurant[]> {
    let url
    if (city) {
      url = `${environment.apiBaseUrl}/restaurants?city=${city}`;
    } else if (name) {
      url = `${environment.apiBaseUrl}/restaurants?name=${name}`;
    } else {
      url = `${environment.apiBaseUrl}/restaurants`;
    }
    return this.http.get<Restaurant[]>(url).pipe(
      catchError(error => {
        alert('Get restaurants failed');
        console.error('Get restaurants failed:', error);
        return of([]);
      })
    );
  }

  createRestaurant(data: RestaurantCreate): Observable<NewRestaurant|null> {
    // Mock: simply push to local array and return it
    const newRestaurant: NewRestaurant = {
      name: data.name,
      location: data.location,
      phone: data.phone || '',
      openingHours: data.openingHours || '',
      img: 'https://picsum.photos/300/200?random=4',
      description: data.description || ''
    };
    return this.http.post<NewRestaurant>(`${environment.apiBaseUrl}/restaurants`, newRestaurant).pipe(
      catchError(error => {
        return throwError(() => error)
      })
    );
  }
}
