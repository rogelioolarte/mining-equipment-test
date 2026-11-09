import { TestBed } from '@angular/core/testing';
import { MaintenanceAPI } from './maintenance-api';

describe('MaintenanceAPI', () => {
  let service: MaintenanceAPI;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(MaintenanceAPI);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
