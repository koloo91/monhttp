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
    this.loadToken();
  }

  loadToken(): void {
    this.tokenSubject = new BehaviorSubject<string>(localStorage.getItem('userToken'));
  }

  checkLogin(): Observable<void> {
    return this.http.get<void>('/api/login')
  }

  storeUsernameAndPassword(username: string, password: string): Observable<string> {
    localStorage.setItem('userToken', btoa(`${username}:${password}`));
    this.loadToken();
    return of(btoa(`${username}:${password}`));
  }

  clearToken(): void {
    localStorage.setItem('userToken', null);
    this.loadToken();
  }

  public get tokenValue(): string {
    return this.tokenSubject.value;
  }
}
