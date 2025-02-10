import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [FormsModule],
  template: `
    <form (submit)="onLogin()">
      <label>Email:</label>
      <input type="email" [(ngModel)]="email" name="email" required />
      <label>Password:</label>
      <input type="password" [(ngModel)]="password" name="password" required />
      <button type="submit">Login</button>
    </form>
  `,
})
export class LoginComponent {
  email = '';
  password = '';

  constructor(private http: HttpClient, private router: Router) {}

  onLogin() {
    console.log('Logging in with:', this.email, this.password);

    // Send the login request to the backend
    this.http
      .post('http://localhost:8080/login', {
        email: this.email,
        password: this.password,
      })
      .subscribe({
        next: (response) => {
          console.log('Login successful:', response);
          // Handle successful login (TODO: implement navigate to user's homepage)
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
