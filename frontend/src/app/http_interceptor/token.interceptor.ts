import {Injectable} from '@angular/core';
import {HttpEvent, HttpHandler, HttpHeaders, HttpInterceptor, HttpRequest} from '@angular/common/http';
import {Observable} from 'rxjs';
import {AuthenticationService} from '../services/authentication.service';

@Injectable()
export class TokenInterceptor implements HttpInterceptor {
  constructor(private authenticationService: AuthenticationService) {
  }

  intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    if (this.authenticationService.tokenValue) {
      let headers = new HttpHeaders()
      headers = headers.append('Authorization', `Basic ${this.authenticationService.tokenValue}`)
      const securedRequest = req.clone({
        headers
      });
      return next.handle(securedRequest);
    } else {
      return next.handle(req);
    }

  }
}
