import {Injectable} from '@angular/core';
import {MatSnackBar} from '@angular/material/snack-bar';
import {ApiError} from '../models/api-error.model';

@Injectable({
  providedIn: 'root'
})
export class ErrorService {

  constructor(private snackBar: MatSnackBar) {
  }

  setError(apiError: ApiError): void {
    this.snackBar.open(apiError.message, 'OK');
  }
}
