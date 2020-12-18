import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Failure} from '../models/failure.model';
import {Wrapper} from '../models/wrapper.model';
import {map} from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class FailureService {

  constructor(private http: HttpClient) {
  }

  list(serviceId: string, from: string, to: string): Observable<Failure[]> {
    return this.http.get<Wrapper<Failure>>(`/api/services/${serviceId}/failures`, {params: {from, to}})
      .pipe(
        map(wrapper => wrapper.data)
      );
  }
}
