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
    user_id: '',
    store_id: '',
    name: '',
    description: '',
    price: 0,
    stock: 0,
    category: '',
    cover_link: '', // Matches Swagger field
    images: [], // Ensure this is an empty array
  };
  selectedStoreId: string = '';
  isLoading: boolean = false;

  constructor(private http: HttpClient) {}

  // Add a new item
  addItem() {
    if (!this.item.store_id || !this.item.name || this.item.price <= 0) {
      alert('Please fill in the required fields.');
      return;
    }

    // Ensure images is an array
    if (typeof this.item.images === 'string') {
      this.item.images = this.item.images.split(',').map((url: string) => url.trim());
    }

    const headers = new HttpHeaders({ 'Content-Type': 'application/json' });
    this.isLoading = true;

    this.http.post('http://localhost:8080/api/items/add', this.item, { headers })
      .subscribe(
        (response: any) => {
          alert('Item added successfully.');
          this.clearItemForm();
          this.listItems();
        },
        (error) => {
          console.error('Error adding item:', error);
          alert('Failed to add item. Please try again.');
          this.isLoading = false;
        },
        () => {
          this.isLoading = false;
        }
      );
  }

  // List all items for the selected store
  listItems() {
    if (!this.selectedStoreId) {
      alert('Please enter a store ID.');
      return;
    }

    this.isLoading = true;
    this.http.get(`http://localhost:8080/api/items/list/${this.selectedStoreId}`)
      .subscribe(
        (response: any) => {
          this.itemList = Array.isArray(response) ? response : [];
          this.itemList.forEach((item: any) => {
            // Ensure `images` is an array and not empty
            item.images = item.images || [];
            item.cover_link = item.cover_link || ''; // Handle missing cover link
          });
        },
        (error) => {
          console.error('Error fetching items:', error);
          alert('Failed to fetch items.');
          this.isLoading = false;
        },
        () => {
          this.isLoading = false;
        }
      );
  }

  // Get item details by ID
  getItemById(id: string) {
    if (!id) {
      alert('Item ID is required.');
      return;
    }

    this.isLoading = true;
    this.http.get(`http://localhost:8080/api/items/get/${id}`)
      .subscribe(
        (response: any) => {
          this.item = response;
        },
        (error) => {
          console.error('Error fetching item:', error);
          alert('Item not found.');
          this.isLoading = false;
        },
        () => {
          this.isLoading = false;
        }
      );
  }

  // Delete an item by ID
  deleteItem(id: string) {
    if (!id) {
      alert('Item ID is required.');
      return;
    }

    const url = `http://localhost:8080/api/items/delete/${id}`;
    this.isLoading = true;

    this.http.delete(url).subscribe(
      (response: any) => {
        alert(response.message || 'Item deleted successfully.');
        this.listItems();
      },
      (error) => {
        console.error('Error deleting item:', error);
        alert('Failed to delete item. Please try again.');
        this.isLoading = false;
      },
      () => {
        this.isLoading = false;
      }
    );
  }

  // Helper method to clear the item form
  clearItemForm() {
    this.item = {
      user_id: '',
      store_id: '',
      name: '',
      description: '',
      price: 0,
      stock: 0,
      category: '',
      cover_link: '', // Matches Swagger field
      images: [], // Ensure this is an empty array
    };
  }
}
