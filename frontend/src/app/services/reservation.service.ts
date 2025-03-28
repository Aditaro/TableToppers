import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from 'src/environments/environment';

export interface ReservationCreate {
  reservationTime: string; // ISO string
  numberOfGuests: number;
  phoneNumber: string;
}

@Injectable({
  providedIn: 'root'
})
export class ReservationService {
  constructor(private http: HttpClient) {}

  createReservation(restaurantId: string, reservation: ReservationCreate): Observable<any> {
    const url = `${environment.apiBaseUrl}/restaurants/${restaurantId}/reservations`;
    return this.http.post(url, reservation);
  }
}
