import { Component } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-courier',
  standalone: true,
  templateUrl: './courier.component.html',
  styleUrls: ['./courier.component.css'],
  imports: [FormsModule, CommonModule],
})
export class CourierComponent {
  selectedOrderId: string | null = null;
  orderList: any[] = [];
  order: any = {
    id: '',
    user_id: '',
    store_id: '',
    itemIdsInput: '',
    totalPrice: 0,
    pickupLocation: '',
    dropOffLocation: '',
    itemIds: [],
    status: 'Pending',
    store: {},
    user: {},
    items: [],
  };
  selectedCourierId: string = '';
  isLoading: boolean = false;

  constructor(private http: HttpClient) {}


  // List all orders for the selected user
  listOrders() {
    if (!this.selectedCourierId) {
      alert('Please enter a user ID.');
      return;
    }

    this.isLoading = true;
    this.http.get(`http://localhost:8080/orders/list/courier/${this.selectedCourierId}`)
      .subscribe({
        next: (response: any) => (this.orderList = response),
        error: (error) => {
          console.error('Error fetching orders:', error);
          alert('Failed to fetch orders.');
        },
        complete: () => (this.isLoading = false),
      });
  }

  // Get order details by ID
  getOrderById(id: string) {
    if (!id) {
      alert('Order ID is required.');
      return;
    }

    this.selectedOrderId = this.selectedOrderId === id ? null : id;

    if (this.selectedOrderId) {
      this.isLoading = true;
      this.http.get(`http://localhost:8080/orders/${id}`)
        .subscribe({
          next: (response: any) => {
            this.order = response;
            this.order.items = this.order.items || [];
          },
          error: (error) => {
            console.error('Error fetching order:', error);
            alert('Order not found.');
          },
          complete: () => (this.isLoading = false),
        });
    }
  }

  // Delete an order by ID
  deleteOrder(id: string) {
    if (!id) {
      alert('Order ID is required.');
      return;
    }

    const url = `http://localhost:8080/orders/delete/${id}`;
    this.isLoading = true;

    this.http.delete(url).subscribe({
      next: (response: any) => {
        alert(response?.message || 'Order deleted successfully.');
        this.listOrders();
      },
      error: (error) => {
        console.error('Error deleting order:', error);
        alert('Failed to delete order. Please try again.');
      },
      complete: () => (this.isLoading = false),
    });
  }

  // Change the status of an order
updateOrderStatus(id: string, newStatus: string) {
  if (!id || !newStatus) {
    alert('Order ID and new status are required.');
    return;
  }

  // Check if the new status is valid and not 'delivered' or 'canceled'
  const validStatuses = ['pending', 'confirmed', 'shipped'];
  if (!validStatuses.includes(newStatus.toLowerCase())) {
    alert('Invalid status. Status cannot be changed to delivered or canceled directly.');
    return;
  }

  this.isLoading = true;
  const url = `http://localhost:8080/orders/update/${id}`;
  const headers = new HttpHeaders({ 'Content-Type': 'application/json' });
  const body = { status: newStatus };

  this.http.patch(url, body, { headers }).subscribe({
    next: (response: any) => {
      alert(response?.message || 'Order status updated successfully.');
      this.listOrders(); // Refresh the order list after updating status
    },
    error: (error) => {
      console.error('Error updating order status:', error);
      alert('Failed to update order status. Please try again.');
    },
    complete: () => (this.isLoading = false),
  });
}


  // Cancel an order by ID
  cancelOrder(id: string) {
    if (!id) {
      alert('Order ID is required.');
      return;
    }

    if (!confirm('Are you sure you want to cancel this order?')) return;

    this.isLoading = true;
    this.http.patch(`http://localhost:8080/orders/cancel/${id}`, {}).subscribe({
      next: () => {
        alert('Order canceled successfully.');
        this.listOrders();
      },
      error: (error) => {
        const errorMessage = error?.error || 'An error occurred.';
        alert(`Error canceling order: ${errorMessage}`);
        console.error('Error canceling order:', error);
      },
      complete: () => (this.isLoading = false),
    });
  }

  // Clear the order form
  clearOrderForm() {
    this.order = {
      id: '',
      user_id: '',
      store_id: '',
      itemIdsInput: '',
      totalPrice: 0,
      pickupLocation: '',
      dropOffLocation: '',
      itemIds: [],
      status: 'Pending',
      store: {},
      user: {},
      items: [],
    };
  }

  // Get CSS class based on order status
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
