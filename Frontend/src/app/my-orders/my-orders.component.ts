import { Component, OnInit } from '@angular/core';
import { HttpClient, HttpClientModule, HttpHeaders } from '@angular/common/http';
import { Router } from '@angular/router';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-my-orders',
  standalone: true,
  imports: [HttpClientModule, CommonModule],
  templateUrl: './my-orders.component.html',
  styleUrl: './my-orders.component.css'
})
export class MyOrdersComponent implements OnInit{
  orders: any[] = []; // Array to store the orders of a user
  userId: string = '12345'; // el 7eta de lesa me7taga tetzabat 3ashan el id is dynamic according to which user is logged in

  constructor(private http: HttpClient, private router: Router) {}

  ngOnInit(): void { // used to perform initialization logic, such as fetching data or setting up state
    this.fetchOrders(); //responsible for loading the orders immediately after the component is created
  }
  viewOrderDetails(orderId: string) {
    this.router.navigate(['/orders', orderId]);
  }
  fetchOrders() { // respnsible for loading data from the api
    // Retrieve the token from local storage
    const token = localStorage.getItem('token');

    // Check if the token exists
    if (!token) {
      alert("You need to log in to view your orders.");
      this.router.navigate(['/login']);
      return;
    }

    // Set up headers with the token
    const headers = new HttpHeaders({
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    });

    // Make the GET request with the token in headers
    this.http.get(`http://localhost:8080/users/${this.userId}/orders`, { headers }).subscribe(
      (ordersList: any) => {
        this.orders = ordersList;
      },
      error => {
        console.error('Error fetching orders:', error);
        if (error.status === 401) {
          alert("Session expired. Please log in again.");
          this.router.navigate(['/login']);
        } else {
          alert("Failed to fetch orders. Please try again.");
        }
      }
    );
  }
}
