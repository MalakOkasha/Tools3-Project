import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.css'],
  standalone: true,
  imports: [FormsModule, CommonModule] // Add CommonModule here
})
export class RegisterComponent {
  registerObject: any = {
    name: '',
    email: '',
    password: '',
    confirmPassword: '',
    phone: '',
    role: '',
    location: '',
    vehicleType: '', // For Courier
    storeId: '', // For Admin and Courier
    storeName: '', // For Owner
    storeLocation: '', // For Owner
  };

  constructor(private http: HttpClient) {}

  whenRegister() {
    // Form validation
    if (!this.registerObject.name || !this.registerObject.email || !this.registerObject.password || !this.registerObject.confirmPassword || !this.registerObject.phone || !this.registerObject.role || !this.registerObject.location) {
      alert('Please fill in all required fields.');
      return;
    }

    if (this.registerObject.password !== this.registerObject.confirmPassword) {
      alert('Passwords do not match!');
      return;
    }

    const { role } = this.registerObject;
    let apiUrl = '';
    let payload: any = {};

    // Set the API URL and payload based on the role
    switch (role) {
      case 'user':
        apiUrl = 'http://localhost:8080/users/register';
        payload = {
          email: this.registerObject.email,
          location: this.registerObject.location,
          name: this.registerObject.name,
          password: this.registerObject.password,
          phone: this.registerObject.phone,
        };
        break;
      case 'courier':
        apiUrl = 'http://localhost:8080/couriers/register';
        payload = {
          email: this.registerObject.email,
          location: this.registerObject.location,
          name: this.registerObject.name,
          password: this.registerObject.password,
          phone: this.registerObject.phone,
          vehicle_type: this.registerObject.vehicleType,
          store_id: this.registerObject.storeId,
        };
        break;
      case 'admin':
        apiUrl = 'http://localhost:8080/admins/register';
        payload = {
          email: this.registerObject.email,
          location: this.registerObject.location,
          name: this.registerObject.name,
          password: this.registerObject.password,
          phone: this.registerObject.phone,
          store_id: this.registerObject.storeId,
        };
        break;
      case 'owner':
        apiUrl = 'http://localhost:8080/owners/register';
        payload = {
          email: this.registerObject.email,
          location: this.registerObject.location,
          name: this.registerObject.name,
          password: this.registerObject.password,
          phone: this.registerObject.phone,
          store_name: this.registerObject.storeName,
          store_location: this.registerObject.storeLocation,
        };
        break;
      default:
        alert('Invalid role selected');
        return;
    }

    // Send the HTTP POST request
    this.http.post(apiUrl, payload).subscribe(
      (response) => {
        console.log('Registration successful:', response);
        alert('Registration successful!');
      },
      (error) => {
        console.error('Registration error:', error);
        if (error.error && error.error.message) {
          alert(`Error: ${error.error.message}`);
        } else {
          alert('Registration failed. Please try again.');
        }
      }
    );
  }

  onRoleChange() {
    // Reset fields based on selected role
    this.registerObject.location = '';
    this.registerObject.vehicleType = '';
    this.registerObject.storeId = '';
    this.registerObject.storeName = '';
    this.registerObject.storeLocation = '';
  }
}
