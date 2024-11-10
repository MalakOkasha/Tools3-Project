

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
import {UpdateOrdersComponent} from './update-orders/update-orders.component';
import { ItemComponent } from './Items/Item.component';


export const routes: Routes = [
  { path: '', redirectTo: '/item', pathMatch: 'full' },
  { path: 'login', component: LoginComponent },
  { path: 'register', component: RegisterComponent },
  { path: 'place-order', component: PlaceOrderComponent },
  { path: 'my-orders', component: MyOrdersComponent },
  { path: 'order-details/:id', component: OrderDetailsComponent },
  // new routes
  { path: 'update-order', component: UpdateOrdersComponent },
  // Items Managements 
  { path: 'item', component: ItemComponent },

];
