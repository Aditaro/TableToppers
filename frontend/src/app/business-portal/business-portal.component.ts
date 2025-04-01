import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-business-portal',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './business-portal.component.html',
  styleUrls: ['./business-portal.component.css']
})
export class BusinessPortalComponent {
  features = [
    {
      title: 'Table Management',
      description: 'Create and manage table layouts with detailed guest capacity settings',
      icon: 'fas fa-chair'
    },
    {
      title: 'Reservation System',
      description: 'Handle reservations, assign tables, and manage guest seating',
      icon: 'fas fa-calendar-check'
    },
    {
      title: 'Menu Management',
      description: 'Create and update your restaurant menu for customer viewing',
      icon: 'fas fa-utensils'
    },
    {
      title: 'Availability Control',
      description: 'Adjust reservation availability for holidays and special events',
      icon: 'fas fa-clock'
    },
    {
      title: 'Staff Dashboard',
      description: 'Real-time table status and wait time estimation',
      icon: 'fas fa-users'
    },
    {
      title: 'Digital Waitlist',
      description: 'Manage walk-in customers with digital waitlist notifications',
      icon: 'fas fa-list-ol'
    }
  ];

  stats = [
    { number: '50%', label: 'Faster Table Turnover' },
    { number: '30%', label: 'Increased Revenue' },
    { number: '25%', label: 'Reduced Wait Times' },
    { number: '90%', label: 'Customer Satisfaction' }
  ];
} 