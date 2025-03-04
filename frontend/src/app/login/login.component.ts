import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { FormsModule } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { CommonModule } from '@angular/common';
import {UserService} from '../services/user.service';
import {environment} from '../../environments/environment';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css'],
  standalone: true,
  imports: [FormsModule, RouterLink, CommonModule]
})
export class LoginComponent {
  email = '';
  password = '';

  constructor(private http: HttpClient,
              private router: Router,
              private userService: UserService) {}

  onLogin() {
    console.log('Logging in with:', this.email, this.password);

    // Send the login request to the backend
    this.http
      .post(`${environment.apiBaseUrl}/auth/login`, {
        email: this.email,
        password: this.password,
      })
      .subscribe({
        next: (response: any) => {
          console.log('Login successful:', response);
          // Handle successful login (TODO: implement navigate to user's homepage)
          localStorage.setItem('user', JSON.stringify(response.user));
          this.router.navigate(['/']);
        },
        error: (error) => {
          alert('Login failed');
          console.error('Login failed:', error);
          // Handle error (TODO: show error message to the user)
        },
        complete: () => {
          alert('Login request complete');
          console.log('Login request complete');
        },
      });
  }
}
