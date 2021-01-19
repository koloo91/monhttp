import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Service} from '../models/service.model';
import {Observable} from 'rxjs';
import {map} from 'rxjs/operators';
import {Wrapper} from '../models/wrapper.model';

@Injectable({
  providedIn: 'root'
})
export class ServiceService {

  constructor(private http: HttpClient) {
  }

  create(service: Service): Observable<Service> {
    return this.http.post<Service>('/api/services', service);
  }

  list(pageSize: number, page: number): Observable<Wrapper<Service>> {
    return this.http.get<Wrapper<Service>>('/api/services', {params: {pageSize: `${pageSize}`, page: `${page}`}});
  }

  get(id: string): Observable<Service> {
    return this.http.get<Service>(`/api/services/${id}`);
  }

  put(id: string, service: Service): Observable<Service> {
    console.log(service);
    return this.http.put<Service>(`/api/services/${id}`, service);
  }

  delete(id: string): Observable<void> {
    return this.http.delete<void>(`/api/services/${id}`);
  }
}
