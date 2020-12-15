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
}
