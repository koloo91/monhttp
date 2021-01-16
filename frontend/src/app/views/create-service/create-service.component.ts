import {Component, OnDestroy, OnInit} from '@angular/core';
import {FormControl, FormGroup, Validators} from '@angular/forms';
import {Service, ServiceType} from '../../models/service.model';
import {Observable, Subscription} from 'rxjs';
import {ServiceService} from '../../services/service.service';
import {Router} from '@angular/router';
import {ErrorService} from '../../services/error.service';
import {ApiError} from '../../models/api-error.model';
import {map, tap} from 'rxjs/operators';
import {NotifierService} from '../../services/notifier.service';
import {Notifier} from '../../models/notifier.model';

@Component({
  selector: 'app-create-service',
  templateUrl: './create-service.component.html',
  styleUrls: ['./create-service.component.scss']
})
export class CreateServiceComponent implements OnInit, OnDestroy {

  formGroup = new FormGroup({
    name: new FormControl('', [Validators.required]),
    type: new FormControl('HTTP', [Validators.required]),
    intervalInSeconds: new FormControl(30, [Validators.required, Validators.min(30), Validators.max(1800)]),

    endpoint: new FormControl('', [Validators.required]),
    requestTimeoutInSeconds: new FormControl(10, [Validators.min(1), Validators.max(180)]),
    httpMethod: new FormControl('GET'),
    httpBody: new FormControl(''),
    httpHeaders: new FormControl(''),
    expectedHttpResponseBody: new FormControl(''),
    expectedHttpStatusCode: new FormControl(200),
    followRedirects: new FormControl(true),
    verifySsl: new FormControl(true),

    enableNotifications: new FormControl(true),
    notifyAfterNumberOfFailures: new FormControl(2, [Validators.min(0), Validators.max(20)]),
    continuouslySendNotifications: new FormControl(false),
    notifiers: new FormControl(['global'])
  });

  selectedServiceType: ServiceType = 'HTTP';

  isLoading = false;

  notifiers$: Observable<Notifier[]>;

  private serviceTypeSubscription: Subscription;

  constructor(private serviceService: ServiceService,
              private errorService: ErrorService,
              private notifierService: NotifierService,
              private router: Router) {

  }

  ngOnInit(): void {
    this.notifiers$ = this.notifierService.list()
      .pipe(
        map(notifiers => {
          notifiers.splice(0, 0, {id: 'global', name: 'Global', data: null, form: []});
          return notifiers;
        })
      );

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

    this.disableFormAllFields();

    const service = this.formGroup.value as Service;
    this.serviceService.create(service)
      .pipe(
        tap(() => this.isLoading = false),
        tap(() => this.enableFormAllFields())
      )
      .subscribe(() => {
        this.router.navigate(['/services']);
      }, (error: ApiError) => {
        this.errorService.setError(error);
      });
  }

  get enableNotifications(): boolean {
    return (this.formGroup.get('enableNotifications') as FormControl).value;
  }

  get httpMethod(): string {
    return (this.formGroup.get('httpMethod') as FormControl).value;
  }

  disableFormAllFields() {
    for (let controlKey in this.formGroup.controls) {
      this.formGroup.get(controlKey).disable();
    }
  }

  enableFormAllFields() {
    for (let controlKey in this.formGroup.controls) {
      this.formGroup.get(controlKey).disable();
    }
  }
}
