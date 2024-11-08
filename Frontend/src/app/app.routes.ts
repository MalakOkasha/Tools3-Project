

// import { Routes } from '@angular/router';
// import { LoginComponent } from './login/login.component';
// import { RegisterComponent } from './register/register.component';

// export const routes: Routes = [
//   { path: '', redirectTo: '/login', pathMatch: 'full' },
//   { path: 'login', component: LoginComponent },
//   { path: 'register', component: RegisterComponent },
// ];
import { Routes } from '@angular/router';
import { LoginComponent } from './login/login.component';
import { RegisterComponent } from './register/register.component';
import { PlaceOrderComponent } from './place-order/place-order.component';
import { MyOrdersComponent } from './my-orders/my-orders.component';
import { OrderDetailsComponent } from './order-details/order-details.component';
import {AssignedOrdersComponent} from './assigned-orders/assigned-orders.component';
import {AssignedOrdersToCourierComponent} from './assigned-orders-to-courier/assigned-orders-to-courier.component';
import {ManageOrdersComponent} from './manage-orders/manage-orders.component';
import {UpdateOrdersComponent} from './update-orders/update-orders.component';


export const routes: Routes = [
  { path: '', redirectTo: '/login', pathMatch: 'full' },
  { path: 'login', component: LoginComponent },
  { path: 'register', component: RegisterComponent },
  { path: 'place-order', component: PlaceOrderComponent },
  { path: 'my-orders', component: MyOrdersComponent },
  { path: 'order-details/:id', component: OrderDetailsComponent },
  // new routes
  { path: 'assigned-orders', component: AssignedOrdersComponent },
  { path: 'assigned-orders-to-courier', component: AssignedOrdersToCourierComponent },
  { path: 'manage-order', component: ManageOrdersComponent },
  { path: 'update-order', component: UpdateOrdersComponent },
];
