import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Wrapper} from '../models/wrapper.model';
import {ImportResult} from '../models/upload_result.model';

@Injectable({
  providedIn: 'root'
})
export class ImportService {

  constructor(private http: HttpClient) {
  }

  importCsv(formData: FormData): Observable<Wrapper<ImportResult>> {
    return this.http.post<Wrapper<ImportResult>>(`/api/import`, formData);
  }
}
