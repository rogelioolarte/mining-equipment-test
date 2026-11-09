import { TestBed } from '@angular/core/testing';
import { EquipmentTypeAPI } from './equipment-type-api';

describe('EquipmentTypeAPI', () => {
  let service: EquipmentTypeAPI;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(EquipmentTypeAPI);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
