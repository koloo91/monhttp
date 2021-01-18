import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Wrapper} from '../models/wrapper.model';
import {UploadResult} from '../models/upload_result.model';

@Injectable({
  providedIn: 'root'
})
export class ImportService {

  constructor(private http: HttpClient) {
  }

  upload(formData: FormData): Observable<Wrapper<UploadResult>> {
    return this.http.post<Wrapper<UploadResult>>(`/api/upload`, formData);
  }
}
