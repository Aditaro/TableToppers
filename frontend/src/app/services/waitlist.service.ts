import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from 'src/environments/environment';
import { ReservationService } from './reservation.service';
import { ReservationCreate } from '../models/reservation.model';
import { WaitlistEntry } from '../models/waitlist.model';

@Injectable({
  providedIn: 'root'
})
export class WaitlistService {
  constructor(
    private http: HttpClient,
    private reservationService: ReservationService
  ) {}

  getWaitlist(restaurantId: string): Observable<any> {
    const url = `${environment.apiBaseUrl}/restaurants/${restaurantId}/waitlist`;
    return this.http.get(url);
  }

  addToWaitlist(restaurantId: string, entryData: { name: string; partySize: number; phoneNumber: string }): Observable<any> {
    const url = `${environment.apiBaseUrl}/restaurants/${restaurantId}/waitlist`;
    return this.http.post(url, entryData);
  }

  updateWaitlistEntry(restaurantId: string, entryId: string, status: 'waiting' | 'seated' | 'cancelled', tableId?: string): Observable<any> {
    const url = `${environment.apiBaseUrl}/restaurants/${restaurantId}/waitlist/${entryId}`;
    return this.http.put(url, { status, tableId });
  }

  /**
   * Seats a customer from the waitlist by creating a reservation and updating the waitlist entry
   * @param restaurantId The restaurant ID
   * @param entry The waitlist entry to convert to a reservation
   * @param tableId The table ID to assign
   * @returns Observable of the operation result
   */
  seatCustomerAsReservation(restaurantId: string, entry: WaitlistEntry, tableId: string): Observable<any> {
    // Create a reservation from the waitlist entry
    const reservation: ReservationCreate = {
      reservationTime: new Date().toISOString(), // Current time
      numberOfGuests: entry.partySize,
      phoneNumber: entry.phoneNumber
    };

    // Create the reservation first
    return this.reservationService.createReservation(restaurantId, reservation);
  }
}