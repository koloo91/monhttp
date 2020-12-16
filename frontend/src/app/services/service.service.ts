import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Service} from '../models/service.model';
import {Observable} from 'rxjs';
import {BaseService} from './base.service';
import {catchError, map} from 'rxjs/operators';
import {Wrapper} from '../models/wrapper.model';

@Injectable({
  providedIn: 'root'
})
export class ServiceService extends BaseService {

  constructor(private http: HttpClient) {
    super();
  }

  create(service: Service): Observable<Service> {
    return this.http.post<Service>('/api/services', service)
      .pipe(
        catchError(this.mapError)
      );
  }

  list(): Observable<Service[]> {
    return this.http.get<Wrapper<Service>>('/api/services')
      .pipe(
        map(wrapper => wrapper.data),
        catchError(this.mapError)
      )
  }

  get(id: string): Observable<Service> {
    return this.http.get<Service>(`/api/services/${id}`)
      .pipe(
        catchError(this.mapError)
      );
  }

  put(id: string, service: Service): Observable<Service> {
    console.log(service);
    return this.http.put<Service>(`/api/services/${id}`, service)
      .pipe(
        catchError(this.mapError)
      )
  }

  delete(id: string): Observable<void> {
    return this.http.delete<void>(`/api/services/${id}`)
      .pipe(
        catchError(this.mapError)
      );
  }
}
