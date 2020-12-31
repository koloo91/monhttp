import {Component, OnInit} from '@angular/core';
import {FormControl, FormGroup, Validators} from '@angular/forms';
import {AuthenticationService} from '../../services/authentication.service';
import {catchError, switchMap} from 'rxjs/operators';
import {from, throwError} from 'rxjs';
import {Router} from '@angular/router';
import {ErrorService} from '../../services/error.service';
import {ApiError} from '../../models/api-error.model';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnInit {

  loginFormGroup = new FormGroup({
    username: new FormControl('', [Validators.required]),
    password: new FormControl('', [Validators.required])
  })

  constructor(private authenticationService: AuthenticationService,
              private errorService: ErrorService,
              private router: Router) {
  }

  ngOnInit(): void {
  }

  onSubmit(): void {
    if (!this.loginFormGroup.valid) {
      console.log('Invalid form');
      return;
    }

    const {username, password} = this.loginFormGroup.value;
    this.authenticationService.storeUsernameAndPassword(username, password)
      .pipe(
        switchMap(() => this.authenticationService.checkLogin()),
        switchMap(() => from(this.router.navigate(['']))),
        catchError((error: ApiError) => {
          this.errorService.setError(error);
          this.authenticationService.clearToken();
          this.loginFormGroup.get('username').setValue(username);
          this.loginFormGroup.get('password').setValue(password);
          return throwError(error);
        })
      ).subscribe(console.log, console.log);
  }
}
