import { HttpClient } from '@angular/common/http';
import { inject, Service } from '@angular/core';
import { environment } from '../../environments/environment';
import { Maintenance, MaintenanceRequest } from '../models/maintenance';
import { Observable } from 'rxjs';

@Service()
export class MaintenanceAPI {
    http = inject(HttpClient)
    baseURL = environment.API_URL+`/maintenance`

    create(request: MaintenanceRequest): Observable<Maintenance> {
        return this.http.post<Maintenance>(this.baseURL, request)
    }

}
