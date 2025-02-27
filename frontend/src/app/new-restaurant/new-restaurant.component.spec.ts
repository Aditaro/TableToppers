import { NewRestaurantComponent } from './new-restaurant.component'
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import {provideAnimations} from '@angular/platform-browser/animations';


describe('NewRestaurantComponent', () => {
  let component: NewRestaurantComponent;
  let fixture: ComponentFixture<NewRestaurantComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [NewRestaurantComponent],
      providers: [
        { provide: MatDialogRef, useValue: {} },
        { provide: MAT_DIALOG_DATA, useValue: {} },
        provideAnimations()
      ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(NewRestaurantComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should mount', () => {
    expect(component).to.be.ok
  });
});
