import { Component } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';


@Component({
  selector: 'app-order',
  standalone: true,
  templateUrl: './my-orders.component.html',
  styleUrls: ['./my-orders.component.css'],
  imports: [FormsModule, CommonModule],
})
export class OrderComponent {
  selectedOrderId: string | null = null;  // Track selected order ID
  orderList: any[] = [];
  order: any = {
    id: '',
    user_id: '',
    store_id: '',
    itemIdsInput: '',
    total_Price: 0,
    pickupLocation: '',
    dropOffLocation: '',
    itemIds: [],
    status: 'Pending',
    store: {},
    user: {},
    items: []
  };
  selectedUserId: string = '';
  isLoading: boolean = false;

  constructor(private http: HttpClient) {}

  // Add a new order
  addOrder() {
    if (!this.order.user_id || !this.order.item_ids ) {
      alert('Please fill in the required fields.');
      return;
    }

    // Convert comma-separated item IDs to an array
    this.order.item_ids = this.order.item_ids.split(',').map((id: string) => id.trim());

    const headers = new HttpHeaders({ 'Content-Type': 'application/json' });

    this.isLoading = true;
    this.http.post('http://localhost:8080/orders/add', this.order, { headers })
      .subscribe(
        (response: any) => {
          alert('Order added successfully.');
          this.clearOrderForm();
          this.listOrders();
        },
        (error: any) => {
          console.error('Error adding order:', error);
          alert('Failed to add order. Please try again.');
        },
        () => {
          this.isLoading = false;
        }
      );
  }

  // List all orders for the selected user
  listOrders() {
    if (!this.selectedUserId) {
      alert('Please enter a user ID.');
      return;
    }

    this.isLoading = true;
    this.http.get(`http://localhost:8080/orders/list/user/${this.selectedUserId}`)
      .subscribe(
        (response: any) => {
          this.orderList = response;
        },
        (error: any) => {
          console.error('Error fetching orders:', error);
          alert('Failed to fetch orders.');
        },
        () => {
          this.isLoading = false;
        }
      );
  }

  // Get order details by ID
  getOrderById(id: string) {
    if (!id) {
      alert('Order ID is required.');
      return;
    }

    // Toggle the selectedOrderId to show or hide the details
    if (this.selectedOrderId === id) {
      this.selectedOrderId = null;  // If already selected, deselect the order
    } else {
      this.selectedOrderId = id;

      this.isLoading = true;
      this.http.get(`http://localhost:8080/orders/${id}`)
        .subscribe(
          (response: any) => {
            console.log('Order Details Response:', response);
            // Make sure we are correctly assigning the data to the order object
            this.order = response;  // Assuming response has the full order object
            if (!this.order.items) {
              this.order.items = [];  // Ensure items is an array
            }
          },
          (error: any) => {
            console.error('Error fetching order:', error);
            alert('Order not found.');
          },
          () => {
            this.isLoading = false;
          }
        );
    }
  }

  // Delete an order by ID
  deleteOrder(id: string) {
    if (!id) {
      alert('Order ID is required.');
      return;
    }

    console.log('Deleting order with ID:', id);

    const url = `http://localhost:8080/orders/delete/${id}`;

    this.isLoading = true;

    this.http.delete(url).subscribe(
      (response: any) => {
        console.log('Response:', response);

        if (response && response.message === 'Order deleted successfully') {
          alert('Order deleted successfully.');
          this.listOrders();
        } else {
          alert('Failed to delete order. Response message: ' + response.message);
        }
      },
      (error: any) => {
        console.error('Error:', error);
        alert('Failed to delete order. Please try again.');
      },
      () => {
        this.isLoading = false;
      }
    );
  }


  // Cancel an order by ID
  cancelOrder(id: string) {
    if (!id) {
      alert('Order ID is required.');
      return;
    }

    const confirmation = confirm('Are you sure you want to cancel this order?');
    if (!confirmation) return;

    this.isLoading = true;
    this.http.patch(`http://localhost:8080/orders/cancel/${id}`, {}).subscribe({
      next: (response: any) => {
        alert('Order canceled successfully.');
        this.listOrders();
      },
      error: (error: any) => {
        // Improved error handling
        if (error.error && typeof error.error === 'string') {
          alert(error.error);
        } else if (error.status === 0) {
          alert('Network error: Unable to reach the server. Please check your connection.');
        } else {
          alert(`Unexpected error: ${error.message || 'An error occurred.'}`);
        }
        console.error('Error canceling order:', error);
      },
      complete: () => {
        this.isLoading = false;
      },
    });
  }





  // Helper method to clear the order form
  clearOrderForm() {
    this.order = {
      id: '',
      user_id: '',
      store_id: '',
      itemIdsInput: '',
      total_Price: 0,
      pickupLocation: '',
      dropOffLocation: '',
      itemIds: [],
      status: 'Pending',
      store: {},
      user: {},
      items: []
    };
  }

  // Method to get the CSS class for order status
  getStatusClass(status: string): string {
    switch (status.toLowerCase()) {
      case 'pending':
        return 'pending';
      case 'shipped':
        return 'shipped';
      case 'delivered':
        return 'delivered';
      case 'cancelled':
        return 'cancelled';
      default:
        return '';
    }
  }
}
