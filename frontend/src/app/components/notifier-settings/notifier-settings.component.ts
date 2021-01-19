import {Component, Input, OnInit} from '@angular/core';
import {Notifier} from '../../models/notifier.model';
import {FormBuilder, FormGroup} from '@angular/forms';
import {NotifierService} from '../../services/notifier.service';

@Component({
  selector: 'app-notifier-settings',
  templateUrl: './notifier-settings.component.html',
  styleUrls: ['./notifier-settings.component.scss']
})
export class NotifierSettingsComponent implements OnInit {

  @Input()
  set notifier(notifier: Notifier) {
    this.notifierFormGroup = this.fb.group({});
    notifier.form.forEach(form => {
      this.notifierFormGroup.addControl(form.formControlName, this.fb.control(notifier.data[form.formControlName]));
    });
    this._notifier = notifier;
  }

  _notifier: Notifier

  notifierFormGroup: FormGroup = this.fb.group({});

  constructor(private notifierService: NotifierService,
              private fb: FormBuilder) {
  }

  ngOnInit(): void {

  }

  updateNotifier(): void {
    this.notifierService.put(this._notifier.id, this.notifierFormGroup.value)
      .subscribe(
        console.log,
        console.log
      );
  }

  testUpTemplate(): void {
    this.notifierService.testUpTemplate(this._notifier.id, this.notifierFormGroup.value)
      .subscribe(
        console.log,
        console.log
      );
  }

  testDownTemplate(): void {
    this.notifierService.testDownTemplate(this._notifier.id, this.notifierFormGroup.value)
      .subscribe(
        console.log,
        console.log
      );
  }
}
