import { Component } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-item',
  standalone: true,
  templateUrl: './Item.component.html',
  styleUrls: ['./Item.component.css'],
  imports: [CommonModule, FormsModule],
})
export class ItemComponent {
  itemList: any[] = [];
  item: any = {
    userId: '',
    storeId: '',
    name: '',
    description: '',
    price: 0,
    stock: 0,
    category: '',
    coverLink: '',
    images: [],
  };
  selectedStoreId: string = '';
  selectedItemId: string = '';
  isLoading: boolean = false;

  constructor(private http: HttpClient) {}

  // Add a new item
  addItem() {
    if (!this.item.storeId || !this.item.name || this.item.price <= 0) {
      alert('Please fill in the required fields.');
      return;
    }
  
    const headers = new HttpHeaders({ 'Content-Type': 'application/json' });
  
    this.isLoading = true;
    this.http.post('http://localhost:8080/items/add', this.item, { headers })
      .subscribe(
        (response: any) => {
          alert('Item added successfully.');
          this.clearItemForm();
          this.listItems();
        },
        (error) => {
          console.error('Error adding item:', error);
          alert('Failed to add item. Please try again.');
        },
        () => {
          this.isLoading = false;
        }
      );
  }
  
  
  
  

  // List all items for the selected store
  async listItems() {
    if (!this.selectedStoreId) {
      alert('Please enter a store ID.');
      return;
    }

    this.isLoading = true;
    try {
      const response = await this.http
        .get(`http://localhost:8080/items/list/${this.selectedStoreId}`)
        .toPromise();
      this.itemList = response as any[];
    } catch (error) {
      console.error('Error fetching items:', error);
      alert('Failed to fetch items.');
    } finally {
      this.isLoading = false;
    }
  }

  // Get item details by ID
  async getItemById(id: string) {
    if (!id) {
      alert('Item ID is required.');
      return;
    }

    this.isLoading = true;
    try {
      const response = await this.http
        .get(`http://localhost:8080/items/get/${id}`)
        .toPromise();
      this.item = response;
    } catch (error) {
      console.error('Error fetching item:', error);
      alert('Item not found.');
    } finally {
      this.isLoading = false;
    }
  }

  // Delete an item by ID
  async deleteItem(id: string) {
    if (!id) {
      alert('Item ID is required.');
      return;
    }

    this.isLoading = true;
    try {
      await this.http
        .delete(`http://localhost:8080/items/delete/${id}`)
        .toPromise();
      alert('Item deleted successfully.');
      this.listItems();
    } catch (error) {
      console.error('Error deleting item:', error);
      alert('Failed to delete item.');
    } finally {
      this.isLoading = false;
    }
  }

  // Helper method to clear the item form
  clearItemForm() {
    this.item = {
      userId: '',
      storeId: '',
      name: '',
      description: '',
      price: 0,
      stock: 0,
      category: '',
      coverLink: '',
      images: [],
    };
  }
}
