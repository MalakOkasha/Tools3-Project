import {Component} from '@angular/core';
import {HttpClient, HttpClientModule} from '@angular/common/http';
import {GlobalService} from '../services/global.service';
import {FormsModule} from '@angular/forms';

@Component({
  selector: 'app-update-orders',
  standalone: true,
  imports: [FormsModule, HttpClientModule],
  templateUrl: './update-orders.component.html',
  styleUrl: './update-orders.component.css'
})
export class UpdateOrdersComponent {
  updateOrderObject: updateOrder;

  constructor(private http: HttpClient) {
    this.updateOrderObject = new updateOrder()
  }

  updateOrder() {
    this.http.post('http://localhost:8080/update-order', this.updateOrderObject).subscribe(
      (res: any) => {
        alert('order updated successfully');
      }
    );
  }
}

export class updateOrder {
  order_id: any;
  status: any;

  constructor() {
    this.order_id = '';
    this.status = '';
  }
}

