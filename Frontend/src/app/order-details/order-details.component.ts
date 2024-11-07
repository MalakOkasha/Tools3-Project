import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-order-details',
  templateUrl: './order-details.component.html',
  styleUrls: ['./order-details.component.css']
})
export class OrderDetailsComponent implements OnInit {
  orderId: string | null = null;
  orderDetails: any;

  constructor(private route: ActivatedRoute, private http: HttpClient) {}

  ngOnInit(): void {
    this.orderId = this.route.snapshot.paramMap.get('id');
    if (this.orderId) {
      this.fetchOrderDetails(this.orderId);
    }
  }

  fetchOrderDetails(orderId: string) {
    this.http.get(`http://localhost:8080/orders/${orderId}`).subscribe(
      (data) => {
        this.orderDetails = data;
      },
      (error) => {
        console.error('Error fetching order details:', error);
      }
    );
  }

  cancelOrder() {
    if (this.orderDetails.status === 'pending') {
      this.http.delete(`http://localhost:8080/orders/${this.orderId}`).subscribe(
        () => {
          alert('Order cancelled successfully!');
          this.orderDetails.status = 'cancelled';
        },
        (error) => {
          console.error('Error cancelling order:', error);
          alert('Failed to cancel the order. Please try again.');
        }
      );
    }
  }
}
