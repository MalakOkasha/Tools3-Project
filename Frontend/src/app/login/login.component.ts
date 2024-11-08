import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';
import {GlobalService} from '../services/global.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [FormsModule, HttpClientModule],
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent {
  loginObject: Login;

  constructor(private http: HttpClient,private global:GlobalService) {
    this.loginObject = new Login();
  }

  whenLogin() {
    console.log('Login button clicked:', this.loginObject);

    // Define headers for the request
    const headers = new HttpHeaders({
      'Content-Type': 'application/json',
      'Accept': 'application/json'
    });

    // Make the POST request
    this.http.post('http://localhost:8080/login', this.loginObject, { headers }).subscribe(
      (res: any) => {
        console.log('Response from login API:', res);

        if (res && res.token && res.user) {
          const token = res.token; // Extract the token
          const user = res.user;   // Extract the user data
          const user_type = res.user.type;   // Extract the user type from data
          alert("You logged in successfully :)");
          localStorage.setItem('token', token); // Save the token
          localStorage.setItem('user_type', user_type); // Save the user type
          this.global.is_login = true; // make user login is global in website
          this.global.type = user_type; // make user type is global in website
        } else {
          alert('Login failed: Invalid response from the server.');
          console.log('Invalid response:', res);
        }
      },
      (error) => {
        console.error('Login error:', error);
        if (error.status === 0) {
          alert('Network error: Unable to reach the server.'); // Handle network error
        } else {
          alert('An error occurred during login. Please try again.'); // Alert user
        }
      }
    );
  }
}

export class Login {
  Email: string;
  Password: string;

  constructor() {
    this.Email = '';
    this.Password = '';
  }
}
