import { TestBed } from '@angular/core/testing';
import { EquipmentAPI } from './equipment-api';

describe('EquipmentAPI', () => {
  let service: EquipmentAPI;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(EquipmentAPI);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
