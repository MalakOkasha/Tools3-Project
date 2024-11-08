import { Component } from '@angular/core';
import {HttpClient, HttpClientModule} from '@angular/common/http';
import {Location} from '@angular/common';
import {FormsModule} from '@angular/forms';

@Component({
  selector: 'app-manage-orders',
  standalone: true,
  imports: [FormsModule, HttpClientModule],
  templateUrl: './manage-orders.component.html',
  styleUrl: './manage-orders.component.css'
})
export class ManageOrdersComponent {

  orders: any;

  constructor(private http: HttpClient,private location: Location) {
    this.http.get('http://localhost:8080/manage-orders').subscribe(
      (res: any) => {
        this.orders = res.data
      }
    );
  }
  updateStatus(id:any){
    this.http.get(`http://localhost:8080/update-status/${id}`).subscribe(
      (res: any) => {
        alert('order status updated successfully');
        this.location.go(this.location.path());
      }
    );
  }
  deleteOrder(id:any){
    this.http.get(`http://localhost:8080/delete-order/${id}`).subscribe(
      (res: any) => {
        alert('order deleted successfully');
        this.location.go(this.location.path());
      }
    );
  }
}
