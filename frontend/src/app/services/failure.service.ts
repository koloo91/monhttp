import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Failure} from '../models/failure.model';
import {Wrapper} from '../models/wrapper.model';
import {map} from 'rxjs/operators';
import {FailureCount} from '../models/failure-count.model';
import {FailureCountByDay} from '../models/failure-count-by-day.model';

@Injectable({
  providedIn: 'root'
})
export class FailureService {

  constructor(private http: HttpClient) {
  }

  list(serviceId: string, from: string, to: string, pageSize: number, page: number): Observable<Wrapper<Failure>> {
    return this.http.get<Wrapper<Failure>>(`/api/services/${serviceId}/failures`, {
      params:
        {
          from,
          to,
          pageSize: `${pageSize}`,
          page: `${page}`
        }
    });
  }

  count(serviceId: string, from: string, to: string): Observable<FailureCount> {
    return this.http.get<FailureCount>(`/api/services/${serviceId}/failures/count`, {params: {from, to}});
  }

  countByDay(serviceId: string, from: string, to: string): Observable<FailureCountByDay[]> {
    return this.http.get<Wrapper<FailureCountByDay>>(`/api/services/${serviceId}/failures/countByDay`, {
      params: {
        from,
        to
      }
    })
      .pipe(
        map(wrapper => wrapper.data)
      );
  }
}
