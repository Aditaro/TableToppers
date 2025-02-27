import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';

@Component({
  selector: 'app-register',
  standalone: true,
  imports: [FormsModule],
  templateUrl: `./register.component.html`,
  styleUrl: `./register.component.css`,
})
export class RegisterComponent {
  email = '';
  password = '';

  constructor(private http: HttpClient, private router: Router) {}

  onRegister() {
    console.log('Registering with:', this.email, this.password);

    // Send the registration request to the backend
    this.http
      .post('http://localhost:8080/register', {
        email: this.email,
        password: this.password,
      })
      .subscribe({
        next: (response) => {
          console.log('Registration successful:', response);
          // Handle successful registration (TODO: implement navigate to login page)
          this.router.navigate(['/login']);
        },
        error: (error) => {
          alert('Registration failed');
          console.error('Registration failed:', error);
          // Handle error (TODO: show error message to the user)
        },
        complete: () => {
          alert('Registration request complete');
          console.log('Registration request complete');
        },
      });
  }
}
