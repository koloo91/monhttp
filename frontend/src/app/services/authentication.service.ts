import {Injectable} from '@angular/core';
import {BehaviorSubject, Observable, of} from 'rxjs';
import {HttpClient} from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class AuthenticationService {

  private tokenSubject: BehaviorSubject<string>;
  public token: Observable<string>;

  constructor(private http: HttpClient) {
    this.tokenSubject = new BehaviorSubject<string>(localStorage.getItem('userToken'));
    this.token = this.tokenSubject.asObservable();
  }

  checkLogin(): Observable<void> {
    return this.http.get<void>('/api/login')
  }

  storeUsernameAndPassword(username: string, password: string): Observable<string> {
    localStorage.setItem('userToken', btoa(`${username}:${password}`));
    this.tokenSubject.next(btoa(`${username}:${password}`))
    return of(btoa(`${username}:${password}`));
  }

  clearToken(): void {
    localStorage.setItem('userToken', null);
    this.tokenSubject.next(null);
  }

  public get tokenValue(): string {
    return this.tokenSubject.value;
  }
}
