import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Setup} from '../models/setup.model';

@Injectable({
  providedIn: 'root'
})
export class SettingsService {

  constructor(private http: HttpClient) {
  }

  setup(): Observable<Setup> {
    return this.http.get<Setup>('/api/setup');
  }

  post(data: any): Observable<any> {
    return this.http.post<any>('/api/settings', data);
  }
}
