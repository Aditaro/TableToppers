import { Component, OnInit } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, ValidationErrors, Validators, ValidatorFn } from '@angular/forms';
import { ReservationService } from '../services/reservation.service';
import { ActivatedRoute } from '@angular/router';
import { MatSnackBar } from '@angular/material/snack-bar';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatNativeDateModule } from '@angular/material/core';
import { MatButtonModule } from '@angular/material/button';
import { MatSelectModule } from '@angular/material/select';
import { Router } from '@angular/router';
import { Restaurant } from '../models/restaurant.model';
import { NgxMaskDirective, provideNgxMask } from 'ngx-mask';

function futureDateValidator(control: AbstractControl): ValidationErrors | null {
  if (!control.value) {
    return null;
  }
  const selectedDate = new Date(control.value);
  if (selectedDate.getTime() < new Date().getTime()) {
    return { pastDate: true };
  }
  return null;
}

@Component({
  selector: 'app-new-reservation',
  standalone: true,
  templateUrl: './new-reservation.component.html',
  styleUrl: './new-reservation.component.css',
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatFormFieldModule,
    MatInputModule,
    MatDatepickerModule,
    MatNativeDateModule,
    MatButtonModule,
    MatSelectModule,
    NgxMaskDirective
  ],
  providers: [provideNgxMask()]
})
export class NewReservationComponent implements OnInit {
  reservationForm!: FormGroup;
  restaurantId!: string;
  restaurant: Restaurant;
  timeSlots: string[] = [];
  minDate: Date = new Date();
  specialAvailabilityDates: string[] = [];

  dateFilter = (date: Date | null): boolean => {
    if (!date) return false;
    const isoDate = date.toISOString().split('T')[0];
    return !this.specialAvailabilityDates.includes(isoDate);
  };

  constructor(
    private fb: FormBuilder,
    private reservationService: ReservationService,
    private route: ActivatedRoute,
    private snackBar: MatSnackBar,
    private router: Router
  ) {
    const navigation = this.router.getCurrentNavigation();
    this.restaurant = navigation?.extras.state?.['restaurant'];
  }

  ngOnInit(): void {
    this.processRestaurantData();
    this.restaurantId = this.route.snapshot.paramMap.get('restaurantId') || '';
    console.log(this.restaurantId);

    this.reservationForm = this.fb.group({
      reservationDate: [null, Validators.required],
      reservationTime: [null, Validators.required],
      numberOfGuests: [null, [Validators.required, Validators.min(1)]],
      phoneNumber: [null, Validators.required]
    }, { validators: this.dateTimeValidator });

    this.reservationForm.get('reservationDate')?.valueChanges.subscribe(date => {
      if (date) this.generateTimeSlots(date);
    });
  }

  private processRestaurantData(): void {
    this.specialAvailabilityDates = this.restaurant.specialAvailability
      .filter(sa => sa.status === 'closed')
      .map(sa => sa.date);

  }

  private parseOpeningHours(): { start: number, startMinutes: number, end: number, endMinutes: number } {
    const timeFormat = /(\d{1,2}):(\d{2})\s*-\s*(\d{1,2}):(\d{2})/;
    const match = this.restaurant.openingHours.match(timeFormat);
    
    if (!match) return { start: 11, startMinutes: 0, end: 22, endMinutes: 0};
  
    const closingHour = parseInt(match[3]);
    const closingMinute = parseInt(match[4]);
    // Convert closing time to total minutes
    const totalClosingMinutes = closingHour * 60 + closingMinute -30;
    
    return {
      start: parseInt(match[1]),
      startMinutes: parseInt(match[2]),
      end: parseInt(match[3]),
      endMinutes: parseInt(match[4]),
    };
  }
  
  private generateTimeSlots(selectedDate: Date): void {
    const { start, startMinutes, end, endMinutes } = this.parseOpeningHours();
    const now = new Date();
    const isToday = selectedDate.toDateString() === now.toDateString();
    
    const slots = [];
    for (let hour = start; hour < end; hour++) {
      for (let minute = startMinutes; minute < 60; minute += 30) {
        if (hour === end - 1 && minute > endMinutes) break;
        if (isToday) {
          const slotTime = new Date(selectedDate);
          slotTime.setHours(hour, minute);
          if (slotTime < now) continue;
        }
        slots.push(`${hour.toString().padStart(2, '0')}:${minute.toString().padStart(2, '0')}`);
      }
    }
    this.timeSlots = slots;
  }
  
  private dateTimeValidator: ValidatorFn = (group: AbstractControl): ValidationErrors | null => {
    const date = group.get('reservationDate')?.value;
    const time = group.get('reservationTime')?.value;
    
    if (!date || !time) return null;
    
    const [hours, minutes] = time.split(':');
    const selectedDateTime = new Date(date);
    selectedDateTime.setHours(parseInt(hours), parseInt(minutes));
    
    return selectedDateTime > new Date() ? null : { pastDateTime: true };
  };
  
  onSubmit(): void {
    if (this.reservationForm.invalid) {
      return;
    }
    const reservationDate = new Date(this.reservationForm.value.reservationDate);
    const [hours, minutes] = this.reservationForm.value.reservationTime.split(':');
    reservationDate.setHours(parseInt(hours), parseInt(minutes));
  
    // Convert date to ISO string for the API call
    const reservation = {
      reservationTime: reservationDate.toISOString(),
      numberOfGuests: this.reservationForm.value.numberOfGuests,
      phoneNumber: this.reservationForm.value.phoneNumber
    };

    this.reservationService.createReservation(this.restaurantId, reservation).subscribe({
      next: () => {
        this.snackBar.open('Reservation successful!', 'Close', { duration: 3000 });
        this.reservationForm.reset();
      },
      error: (err) => {
        console.error(err);
        this.snackBar.open('Failed to make reservation', 'Close', { duration: 3000 });
      }
    });
  }
}
