import {Component, OnInit} from '@angular/core';
import {FormControl, FormGroup, Validators} from '@angular/forms';
import {ServiceType} from '../../models/service.model';

@Component({
  selector: 'app-create-service',
  templateUrl: './create-service.component.html',
  styleUrls: ['./create-service.component.scss']
})
export class CreateServiceComponent implements OnInit {

  formGroup = new FormGroup({
    name: new FormControl('', [Validators.required]),
    type: new FormControl('HTTP', [Validators.required]),
    checkIntervalInSeconds: new FormControl(30, [Validators.required, Validators.min(30), Validators.max(1800)]),

    endpoint: new FormControl('', [Validators.required]),
    requestTimeoutInSeconds: new FormControl(1, [Validators.min(1), Validators.max(180)]),

    enableNotifications: new FormControl(true),
    notifyAfterNumberOfFailures: new FormControl(2, [Validators.min(0), Validators.max(20)])
  });

  constructor() {
  }

  ngOnInit(): void {
  }

  getServiceType(): ServiceType {
    if (!this.formGroup['type']) {
      return 'HTTP';
    }
    console.log(this.formGroup['type'].value);
    return this.formGroup['type'].value as ServiceType;
  }
}
