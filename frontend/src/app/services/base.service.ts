import {HttpErrorResponse} from '@angular/common/http';
import {throwError} from 'rxjs';
import {ApiError} from '../models/api-error.model';

export class BaseService {
  mapError(error: HttpErrorResponse) {
    if (error.error instanceof ErrorEvent) {
      const apiError: ApiError = {message: error.error.message};
      return throwError(apiError);
    } else {
      const apiError: ApiError = {message: error.error};
      return throwError(apiError);
    }
  }
}
