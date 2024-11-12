import { Routes } from '@angular/router';
import { LoginComponent } from './login/login.component';
import { RegisterComponent } from './register/register.component';
import { OrderComponent } from './my-orders/my-orders.component';
import { ItemComponent } from './Items/Item.component';
import { CourierComponent } from './courier/courier.component';

export const routes: Routes = [
  { path: '', redirectTo: '/login', pathMatch: 'full' },
  { path: 'login', component: LoginComponent },
  { path: 'register', component: RegisterComponent },
  { path: 'my-orders', component: OrderComponent },
  { path: 'item', component: ItemComponent },  // Items Management
  { path: 'courier', component: CourierComponent },
];
