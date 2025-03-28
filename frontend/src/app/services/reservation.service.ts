import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from 'src/environments/environment';
import { Reservation, ReservationCreate } from '../models/reservation.model';

@Injectable({
  providedIn: 'root'
})
export class ReservationService {
  constructor(private http: HttpClient) {}

  createReservation(restaurantId: string, reservation: ReservationCreate): Observable<any> {
    const url = `${environment.apiBaseUrl}/restaurants/${restaurantId}/reservations`;
    return this.http.post(url, reservation);
  }

  getReservations(restaurantId: string, userId?: string, date?: Date): Observable<any> {
    let queryParams: any = {};
    if (date) {
      // Format date to YYYY-MM-DD string as required by API
      queryParams.date = date.toISOString().split('T')[0];
    }
    if (userId) {
      queryParams.userId = userId;
    }
    const url = `${environment.apiBaseUrl}/restaurants/${restaurantId}/reservations`;
    return this.http.get(url, { params: queryParams });
  }

  deleteReservation(restaurantId: string, reservationId: string): Observable<any> {
    const url = `${environment.apiBaseUrl}/restaurants/${restaurantId}/reservations/${reservationId}`;
    return this.http.delete(url);
  }

  updateReservation(restaurantId: string, reservationId: string, reservation: Partial<Reservation>): Observable<any> {
      const url = `${environment.apiBaseUrl}/restaurants/${restaurantId}/reservations/${reservationId}`;
      return this.http.put(url, reservation);
  }
}
