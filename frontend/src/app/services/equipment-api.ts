import { HttpClient } from '@angular/common/http';
import { inject, Service } from '@angular/core';
import { environment } from '../../environments/environment';
import { Equipment, EquipmentRequest } from '../models/equipment';
import { Observable } from 'rxjs';
import { Maintenance } from '../models/maintenance';

@Service()
export class EquipmentAPI {
    http = inject(HttpClient)
    baseURL = environment.API_URL+`/equipment`

    getAll(): Observable<Equipment[]> {
        return this.http.get<Equipment[]>(this.baseURL)
    }

    create(request: EquipmentRequest): Observable<Equipment> {
        return this.http.post<Equipment>(this.baseURL, request)
    }

    updateById(id: string, request: EquipmentRequest): Observable<Equipment> {
        return this.http.put<Equipment>(`${this.baseURL}/${id}`, request)
    }

    deleteById(id: string): Observable<void> {
        return this.http.delete<void>(`${this.baseURL}/${id}`)
    }

    getAllMaintenancesByEquipmentId(id: string): Observable<Maintenance[]> {
        return this.http.get<Maintenance[]>(`${this.baseURL}/${id}/maintenance`)
    }
}
