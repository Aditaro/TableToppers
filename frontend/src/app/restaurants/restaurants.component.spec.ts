import { ComponentFixture, TestBed } from '@angular/core/testing';
import { RestaurantsComponent } from './restaurants.component';
import {provideAnimations} from '@angular/platform-browser/animations';

describe('RestaurantsComponent', () => {
  let component: RestaurantsComponent;
  let fixture: ComponentFixture<RestaurantsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        RestaurantsComponent
      ],
      providers: [
        provideAnimations()
      ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(RestaurantsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should mount', () => {
    expect(component).to.be.ok
  })
});
