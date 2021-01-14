import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Notifier} from '../models/notifier.model';
import {Wrapper} from '../models/wrapper.model';
import {map} from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class NotifierService {

  constructor(private http: HttpClient) {
  }

  list(): Observable<Notifier[]> {
    return this.http.get<Wrapper<Notifier>>('/api/notifiers')
      .pipe(
        map(wrapper => wrapper.data)
      );
  }

  put(notifierId: string, data: any): Observable<Notifier> {
    return this.http.put<Notifier>(`/api/notifiers/${notifierId}`, data);
  }

  testUpTemplate(notifierId: string, data: any): Observable<void> {
    return this.http.post<void>(`/api/notifiers/${notifierId}/test/up`, data)
  }

  testDownTemplate(notifierId: string, data: any): Observable<void> {
    return this.http.post<void>(`/api/notifiers/${notifierId}/test/down`, data)
  }
}
