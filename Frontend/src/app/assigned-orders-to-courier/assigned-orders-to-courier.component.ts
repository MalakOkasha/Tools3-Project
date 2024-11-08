import {Component} from '@angular/core';
import {HttpClient, HttpClientModule} from '@angular/common/http';
import {FormsModule} from '@angular/forms';
import { Location } from '@angular/common';

@Component({
  selector: 'app-assigned-orders-to-courier',
  standalone: true,
  imports: [FormsModule, HttpClientModule],
  templateUrl: './assigned-orders-to-courier.component.html',
  styleUrl: './assigned-orders-to-courier.component.css'
})
export class AssignedOrdersToCourierComponent {

  orders: any;

  constructor(private http: HttpClient,private location: Location) {
    this.http.get('http://localhost:8080/assigned-orders-to-courier').subscribe(
      (res: any) => {
        this.orders = res.data
      }
    );
  }

  reassign(id:any){
    this.http.get(`http://localhost:8080/reassigned-order-courier/${id}`).subscribe(
      (res: any) => {
        alert('reassign order courier successfully');
        this.location.go(this.location.path());
      }
    );
  }
}
