import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AssignedOrdersToCourierComponent } from './assigned-orders-to-courier.component';

describe('AssignedOrdersToCourierComponent', () => {
  let component: AssignedOrdersToCourierComponent;
  let fixture: ComponentFixture<AssignedOrdersToCourierComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AssignedOrdersToCourierComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(AssignedOrdersToCourierComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
