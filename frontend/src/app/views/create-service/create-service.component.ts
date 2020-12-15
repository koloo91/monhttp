import {Component, OnDestroy, OnInit} from '@angular/core';
import {FormControl, FormGroup, Validators} from '@angular/forms';
import {Service, ServiceType} from '../../models/service.model';
import {Subscription} from 'rxjs';
import {ServiceService} from '../../services/service.service';
import {Router} from '@angular/router';
import {ErrorService} from '../../services/error.service';
import {ApiError} from '../../models/api-error.model';
import {tap} from 'rxjs/operators';

@Component({
  selector: 'app-create-service',
  templateUrl: './create-service.component.html',
  styleUrls: ['./create-service.component.scss']
})
export class CreateServiceComponent implements OnInit, OnDestroy {

  formGroup = new FormGroup({
    name: new FormControl('', [Validators.required]),
    type: new FormControl('HTTP', [Validators.required]),
    checkIntervalInSeconds: new FormControl(30, [Validators.required, Validators.min(30), Validators.max(1800)]),

    endpoint: new FormControl('', [Validators.required]),
    requestTimeoutInSeconds: new FormControl(1, [Validators.min(1), Validators.max(180)]),
    httpMethod: new FormControl('GET'),
    httpBody: new FormControl(''),
    httpHeaders: new FormControl(''),
    expectedHttpResponseBody: new FormControl(''),
    expectedHttpStatusCode: new FormControl(200),
    followRedirects: new FormControl(true),
    verifySsl: new FormControl(true),

    enableNotifications: new FormControl(true),
    notifyAfterNumberOfFailures: new FormControl(2, [Validators.min(0), Validators.max(20)])
  });

  selectedServiceType: ServiceType = 'HTTP';

  isLoading = false;

  private serviceTypeSubscription: Subscription;

  constructor(private serviceService: ServiceService,
              private router: Router,
              private errorService: ErrorService) {
  }

  ngOnInit(): void {
    this.serviceTypeSubscription = this.serviceType.valueChanges
      .subscribe((serviceType) => this.selectedServiceType = serviceType);
  }

  ngOnDestroy(): void {
    this.serviceTypeSubscription?.unsubscribe();
  }

  get serviceType(): FormControl {
    return this.formGroup.get('type') as FormControl;
  }

  get checkIntervalInSeconds(): FormControl {
    return this.formGroup.get('checkIntervalInSeconds') as FormControl;
  }

  saveService(): void {
    this.isLoading = true;
    const service = this.formGroup.value as Service;
    this.serviceService.create(service)
      .pipe(
        tap(() => this.isLoading = false)
      )
      .subscribe(() => {
        this.router.navigate(['/services']);
      }, (error: ApiError) => {
        this.errorService.setError(error);
      });
  }
}
