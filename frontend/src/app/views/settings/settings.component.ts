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

  constructor(private notifierService: NotifierService) {

  }

  ngOnInit(): void {
    this.loadNotifiers();
  }

  loadNotifiers(): void {
    this.selectedNotifier = null;
    this.notifiers$ = this.notifierService.list()
      .pipe(
        tap((notifiers) => {
          if (!this.selectedNotifier) {
            this.selectedNotifier = notifiers[0];
          }
        })
      );
  }

  notifierSelected(notifier: Notifier): void {
    this.selectedNotifier = notifier;
  }
}
