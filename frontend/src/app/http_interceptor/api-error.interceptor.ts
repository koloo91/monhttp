import {Injectable} from '@angular/core';
import {HttpErrorResponse, HttpEvent, HttpHandler, HttpInterceptor, HttpRequest} from '@angular/common/http';
import {Observable, throwError} from 'rxjs';
import {catchError} from 'rxjs/operators';
import {ApiError} from '../models/api-error.model';
import {Router} from '@angular/router';
import {ErrorService} from '../services/error.service';

@Injectable()
export class ApiErrorInterceptor implements HttpInterceptor {
  constructor(private router: Router,
              private errorService: ErrorService) {
  }

  intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    return next.handle(req)
      .pipe(
        catchError(this.mapError),
        catchError((error: ApiError) => {
          if (error.message === 'monhttp needs to be setup') {
            this.router.navigate(['setup']);
          } else if (error.message === 'invalid credentials') {
            this.router.navigate(['login']);
          } else {
            this.errorService.setError(error);
          }

          return throwError(error);
        })
      );
  }

  mapError(error: HttpErrorResponse) {
    if (error.error instanceof ErrorEvent) {
      const apiError: ApiError = {message: error.error.message};
      return throwError(apiError);
    } else {
      const apiError = error.error as ApiError;
      return throwError(apiError);
    }
  }
}
