
export interface ReservationCreate {
    reservationTime: string; // ISO string
    numberOfGuests: number;
    phoneNumber: string;
  }

  export interface Reservation {
    restaurantId: string;
    userId: string;
    tableId: string;
    reservationTime: string; // ISO date string
    numberOfGuests: number;
    status: 'pending' | 'confirmed' | 'cancelled';
    phoneNumber: string;
}