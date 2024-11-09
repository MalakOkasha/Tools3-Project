import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';
import { Router } from '@angular/router';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [FormsModule, HttpClientModule],
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent {
  loginObject: Login;
  isLoading: boolean = false;

  constructor(private http: HttpClient, private router: Router) {
    this.loginObject = new Login();
  }

  async whenLogin() {
    console.log('Login button clicked:', this.loginObject);

    // Input validation
    if (!this.loginObject.Email || !this.loginObject.Password) {
      alert('Please enter both email and password.');
      return;
    }

    // Define headers for the request
    const headers = new HttpHeaders({
      'Content-Type': 'application/json',
      'Accept': 'application/json'
    });

    // Determine the login API URL based on the selected role
    let apiUrl = '';
    switch (this.loginObject.Role) {
      case 'user':
        apiUrl = 'http://localhost:8080/users/login';
        break;
      case 'courier':
        apiUrl = 'http://localhost:8080/couriers/login';
        break;
      case 'admin':
        apiUrl = 'http://localhost:8080/admins/login';
        break;
      case 'owner':
        apiUrl = 'http://localhost:8080/owners/login';
        break;
      default:
        alert('Invalid role selected.');
        return;
    }

    this.isLoading = true;

    try {
      // Make the POST request using async/await
      const res: any = await this.http.post(apiUrl, this.loginObject, { headers }).toPromise();
      console.log('Response from login API:', res);

      if (res && res.token && res.user) {
        const token = res.token;
        const user = res.user;

        console.log('Login successful!');
        console.log('Token:', token);
        console.log('User:', user);

        alert(`Welcome back, ${user.name || 'User'}!`);
        localStorage.setItem('token', token);
        this.router.navigate(['/place-order']);
      } else {
        alert('Login failed: Invalid credentials or server response.');
        console.log('Invalid response:', res);
      }
    } catch (error: any) {
      console.error('Login error:', error);
      if (error.status === 0) {
        alert('Network error: Unable to reach the server. Please check your internet connection.');
      } else if (error.status === 401) {
        alert('Unauthorized: Incorrect email or password.');
      } else if (error.status === 500) {
        alert('Server error: Please try again later.');
      } else {
        alert(`An unexpected error occurred: ${error.message}`);
      }
    } finally {
      this.isLoading = false;
    }
  }

  // Method to select role (for the role cards UI)
  selectRole(role: string) {
    this.loginObject.Role = role;
    console.log('Selected Role:', role);
  }
}

export class Login {
  Email: string;
  Password: string;
  Role: string;

  constructor() {
    this.Email = '';
    this.Password = '';
    this.Role = 'user';
  }
}
