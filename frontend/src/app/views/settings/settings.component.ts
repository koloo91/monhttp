import {Component, OnInit} from '@angular/core';
import {NotifierService} from '../../services/notifier.service';
import {Observable} from 'rxjs';
import {Notifier} from '../../models/notifier.model';
import {FormBuilder, FormGroup} from '@angular/forms';
import {tap} from 'rxjs/operators';

@Component({
  selector: 'app-settings',
  templateUrl: './settings.component.html',
  styleUrls: ['./settings.component.scss']
})
export class SettingsComponent implements OnInit {

  notifiers$: Observable<Notifier[]>;
  selectedNotifier: Notifier;

  notifierFormGroup: FormGroup = this.fb.group({});

  constructor(private notifierService: NotifierService,
              private fb: FormBuilder) {

  }

  ngOnInit(): void {
    this.loadNotifiers();
  }

  loadNotifiers(): void {
    this.selectedNotifier = null;
    this.notifierFormGroup = null;
    this.notifiers$ = this.notifierService.list()
      .pipe(
        tap(notifiers => this.notifierSelected(notifiers[0]))
      );
  }

  notifierSelected(notifier: Notifier) {
    this.notifierFormGroup = this.fb.group({});
    notifier.form.forEach(form => {
      this.notifierFormGroup.addControl(form.formControlName, this.fb.control(notifier.data[form.formControlName]));
    });

    this.selectedNotifier = notifier;
  }

  updateNotifier(): void {
    this.notifierService.put(this.selectedNotifier.id, this.notifierFormGroup.value)
      .subscribe(
        () => this.loadNotifiers(),
        console.log
      );
  }
}
