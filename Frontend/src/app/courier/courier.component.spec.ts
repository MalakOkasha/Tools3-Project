import { ComponentFixture, TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { CourierComponent } from './courier.component'; // Corrected import
import { FormsModule } from '@angular/forms';

describe('CourierComponent', () => { // Changed to CourierComponent
  let component: CourierComponent;
  let fixture: ComponentFixture<CourierComponent>;
  let httpMock: HttpTestingController;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [CourierComponent], // Corrected to CourierComponent
      imports: [HttpClientTestingModule, FormsModule],
    }).compileComponents();

    fixture = TestBed.createComponent(CourierComponent);
    component = fixture.componentInstance;
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('should create the component', () => {
    expect(component).toBeTruthy();
  });

  it('should fetch the list of orders', () => {
    const mockOrders = [
      {
        id: '1',
        userId: '123',
        courierId: null,
        storeId: 'store123',
        itemIds: ['item1', 'item2'],
        status: 'Pending',
        totalPrice: 100.0,
        pickupLocation: 'Store Location',
        dropOffLocation: 'User Location',
        createdAt: '2024-11-12T00:00:00Z',
        updatedAt: '2024-11-12T00:00:00Z',
      },
    ];

    component.listOrders();
    const req = httpMock.expectOne('http://localhost:8080/orders/list/store123');
    expect(req.request.method).toBe('GET');
    req.flush(mockOrders);

    expect(component.orderList.length).toBe(1);
    expect(component.orderList[0].userId).toBe('123');
  });

  it('should add a new order', () => {
    spyOn(component, 'listOrders');
    component.order = { user_id: '123', store_id: 'store123', itemIdsInput: 'item1,item2', totalPrice: 100 };

    const req = httpMock.expectOne('http://localhost:8080/orders/add');
    expect(req.request.method).toBe('POST');
    req.flush({ message: 'Order added successfully' });

    expect(component.listOrders).toHaveBeenCalled();
    expect(component.order.user_id).toBe('');
    expect(component.order.itemIdsInput).toBe('');
  });

  it('should delete an order', () => {
    spyOn(component, 'listOrders');
    component.deleteOrder('1');

    const req = httpMock.expectOne('http://localhost:8080/orders/delete/1');
    expect(req.request.method).toBe('DELETE');
    req.flush({ message: 'Order deleted successfully' });

    expect(component.listOrders).toHaveBeenCalled();
  });

  it('should fetch order details', () => {
    component.getOrderById('1');
    const req = httpMock.expectOne('http://localhost:8080/orders/get/1');
    expect(req.request.method).toBe('GET');
    req.flush({ id: '1', userId: '123', status: 'Pending' });

    expect(component.order.id).toBe('1');
  });
});
