
import { HttpClient } from '@angular/common/http';
import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';

@Component({
  selector: 'app-register',
  standalone: true,
  imports: [FormsModule, HttpClientModule],
  templateUrl: './register.component.html',
  styleUrl: './register.component.css'
})
export class RegisterComponent {

  registerObject: Register;

  constructor(private http: HttpClient){

    this.registerObject = new Register()
  }

  whenRegister() {
    console.log('Register button clicked:', this.registerObject);

    this.http.post('http://localhost:8080/register', this.registerObject).subscribe(
      (res: any) => {
        if (res.result) {
          alert("You registered successfully :)");
        } else {
          alert('Registration failed: ' + (res.message || 'Invalid response from the server.'));
          console.log('Server response:', res);
        }
      },
      (error) => {
        console.error('Registration error:', error);
        
        // Improved error handling
        if (error.error) {
          alert('Error: ' + (error.error.message || JSON.stringify(error.error))); // Display more specific error message
        } else if (error.status === 0) {
          alert('Network error: Unable to reach the server.');
        } else {
          alert('An error occurred during registration. Status: ' + error.status + ' - ' + error.message);
        }
      }
    );
}

  
}
export class Register{
  email: string;
  password: string;
  phone: string;
  name: string;

  constructor() {
    this.email = '';
    this.password = '';
    this.phone ='';
    this.name ='';

  }
}

