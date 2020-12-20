import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class SettingsService {

  constructor(private http: HttpClient) {
  }

  post(data: any): Observable<any> {
    return this.http.post<any>('/api/settings', data);
  }
}
