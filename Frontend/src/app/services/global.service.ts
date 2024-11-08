import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class GlobalService {
  is_login = localStorage.getItem('token') ? true : false;
  type = localStorage.getItem('user_type');

  constructor() { }
}
