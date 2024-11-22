

// app.component.ts
import { Component } from '@angular/core';
import { RouterOutlet, RouterModule } from '@angular/router';
import {GlobalService} from './services/global.service';
import { Router } from '@angular/router';



@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, RouterModule], // Import RouterModule here
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'my-angular-app';


  constructor(public global:GlobalService, private router: Router ) {

  }


  
  logout() {
    // Clear user session and redirect to login page
    console.log("Logging out...");
    // Example: Clear token and navigate to login
    localStorage.removeItem('token');
    this.router.navigate(['/login']);
  }






}
