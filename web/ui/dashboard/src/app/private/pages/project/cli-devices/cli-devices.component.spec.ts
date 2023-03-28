import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CliDevicesComponent } from './cli-devices.component';

describe('CliDevicesComponent', () => {
  let component: CliDevicesComponent;
  let fixture: ComponentFixture<CliDevicesComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ CliDevicesComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CliDevicesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
