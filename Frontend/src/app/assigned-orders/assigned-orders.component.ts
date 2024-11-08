import { Component } from '@angular/core';
import {HttpClient, HttpClientModule} from '@angular/common/http';
import {Location} from '@angular/common';
import {FormsModule} from '@angular/forms';

@Component({
  selector: 'app-assigned-orders',
  standalone: true,
  imports: [FormsModule, HttpClientModule],
  templateUrl: './assigned-orders.component.html',
  styleUrl: './assigned-orders.component.css'
})
export class AssignedOrdersComponent {


  orders: any;

  constructor(private http: HttpClient,private location: Location) {
    this.http.get('http://localhost:8080/assigned-orders').subscribe(
      (res: any) => {
        this.orders = res.data
      }
    );
  }

  accept(id:any){
    this.http.get(`http://localhost:8080/accept-order/${id}`).subscribe(
      (res: any) => {
        alert('accept order successfully');
        this.location.go(this.location.path());
      }
    );
  }

  reject(id:any){
    this.http.get(`http://localhost:8080/reject-order/${id}`).subscribe(
      (res: any) => {
        alert('reject order successfully');
        this.location.go(this.location.path());
      }
    );
  }
}
