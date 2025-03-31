import { Component } from '@angular/core';
import { ReservationService } from '../services/reservation.service';
import { ActivatedRoute, Router } from '@angular/router';
import { Restaurant } from '../models/restaurant.model';
import { User } from '../models/user.model';
import { MatDialog } from '@angular/material/dialog';
import { NewReservationComponent } from '../new-reservation/new-reservation.component';
import { CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';

@Component({
  selector: 'app-reservations',
  standalone: true,
  templateUrl: './reservations.component.html',
  styleUrl: './reservations.component.css',
  imports: [CommonModule, MatCardModule]
})
export class ReservationsComponent {
  displayedColumns: string[] = ['reservationTime', 'numberOfGuests', 'phoneNumber'];
  reservations: any[] = [];
  restaurant: Restaurant;
  restaurantId: string = '';
  user: User;

  constructor(
    private reservationService: ReservationService,
    private route: ActivatedRoute,
    private router: Router,
    public dialog: MatDialog
  ) {
    const navigation = this.router.getCurrentNavigation();
    this.restaurant = navigation?.extras.state?.['restaurant'];
  }

  ngOnInit(): void {
    this.fetchReservation();
  }
  
  fetchReservation(): void {
    this.user = JSON.parse(localStorage.getItem('user'));
    this.restaurantId = this.route.snapshot.paramMap.get('restaurantId') || '';
    this.reservationService.getReservations(this.restaurantId, this.user.id)
      .subscribe(data => {
        this.reservations = data;
        console.log(data);
      });
  }

  openNewReservationDialog(): void {
    const dialogRef = this.dialog.open(NewReservationComponent, {
      width: '600px',
      data: {
        restaurantId: this.restaurantId,
        restaurant: this.restaurant
      }
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        this.fetchReservation(); // Refresh the list if a new restaurant was created
      }
    });
  }

  deleteReservation(reservationId: string): void {
    this.reservationService.deleteReservation(this.restaurantId, reservationId)
      .subscribe(() => {
        this.fetchReservation();
      });
  }

  modifyReservation(reservation: any): void {
    const dialogRef = this.dialog.open(NewReservationComponent, {
      width: '600px',
      data: {
        restaurantId: this.restaurantId,
        restaurant: this.restaurant,
        reservation: reservation
      }
    });
    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        this.fetchReservation(); 
      }
    })
  }

  formatDate(dateString: string): string {
    if (!dateString) return '';
    
    const date = new Date(dateString);
    
    // Format: March 28, 2025, 8:21 AM
    return date.toLocaleString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: 'numeric',
      minute: '2-digit',
      hour12: true
    });
  }
}
