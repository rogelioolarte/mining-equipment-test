import { HttpClient } from '@angular/common/http';
import { inject, Service } from '@angular/core';
import { environment } from '../../environments/environment';
import { EquipmentType } from '../models/equipment-type';
import { Observable } from 'rxjs';

@Service()
export class EquipmentTypeAPI {
    http = inject(HttpClient)
    baseURL = environment.API_URL+`/equipment_type`

    getAll(): Observable<EquipmentType[]> {
        return this.http.get<EquipmentType[]>(this.baseURL);
    }
}
