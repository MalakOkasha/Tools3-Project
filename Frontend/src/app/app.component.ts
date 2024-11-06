

// import { Component } from '@angular/core';
// import { RouterOutlet } from '@angular/router'; // Import RouterOutlet

// @Component({
//   selector: 'app-root',
//   standalone: true,
//   imports: [RouterOutlet], // Add RouterOutlet here
//   templateUrl: './app.component.html',
//   styleUrls: ['./app.component.css']
// })
// export class AppComponent {
//   title = 'my-angular-app';
// }

// app.component.ts
import { Component } from '@angular/core';
import { RouterOutlet, RouterModule } from '@angular/router';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, RouterModule], // Import RouterModule here
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'my-angular-app';
}
