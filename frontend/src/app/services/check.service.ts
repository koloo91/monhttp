import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Check} from '../models/check.model';
import {Wrapper} from '../models/wrapper.model';
import {map} from 'rxjs/operators';
import {Average} from '../models/average.model';
import {IsOnline} from '../models/is-online.model';

@Injectable({
  providedIn: 'root'
})
export class CheckService {

  constructor(private http: HttpClient) {
  }

  list(serviceId: string, from: string, to: string, reduceByFactor: number): Observable<Check[]> {
    return this.http.get<Wrapper<Check>>(`/api/services/${serviceId}/checks`, {
      params: {
        from,
        to,
        reduceByFactor: `${reduceByFactor}`
      }
    })
      .pipe(
        map(wrapper => wrapper.data)
      )
  }

  average(serviceId: string): Observable<Average> {
    return this.http.get<Average>(`/api/services/${serviceId}/average`)
  }

  isOnline(serviceId: string): Observable<IsOnline> {
    return this.http.get<IsOnline>(`/api/services/${serviceId}/online`);
  }
}
