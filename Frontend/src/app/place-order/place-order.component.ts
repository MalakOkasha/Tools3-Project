import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';
import { Router } from '@angular/router';

@Component({
  selector: 'app-place-order',
  standalone: true,
  imports: [FormsModule, ReactiveFormsModule, HttpClientModule],
  templateUrl: './place-order.component.html',
  styleUrls: ['./place-order.component.css']
})
export class PlaceOrderComponent {
  orderForm: FormGroup;

  constructor(private fb: FormBuilder, private http: HttpClient, private router: Router) {
    this.orderForm = this.fb.group({
      pickup: ['', Validators.required],
      dropOff: ['', Validators.required],
      delivery: ['', Validators.required],
      packageDetails: ['', Validators.required],
      terms: [false, Validators.requiredTrue]
    });
  }

  onSubmit() {
    if (this.orderForm.valid) {
      console.log('Order data:', this.orderForm.value);
      
      // Retrieve the token from local storage
      const token = localStorage.getItem('token');
      
      // Check if the token exists
      if (!token) {
        alert("You need to log in to place an order.");
        this.router.navigate(['/login']);
        return;
      }

      // Set up headers with the token
      const headers = new HttpHeaders({
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}` // Include the token here
      });

      // Send form data to the backend with the token in headers
      this.http.post('http://localhost:8080/orders', this.orderForm.value, { headers }).subscribe(
        response => {
          console.log('Order placed successfully!', response);
          alert("Order placed successfully!");
          this.orderForm.reset(); // Reset the form on success
        },
        error => {
          console.error('Error placing order:', error);
          if (error.status === 401) {
            alert("Session expired. Please log in again.");
            this.router.navigate(['/login']);
          } else {
            alert("Failed to place order. Please try again.");
          }
        }
      );
    } else {
      console.log('Form is invalid');
      alert("Please fill out all required fields and agree to the terms.");
    }
  }
}
