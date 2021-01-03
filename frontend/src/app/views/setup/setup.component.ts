import {Component, OnInit} from '@angular/core';
import {FormControl, FormGroup, Validators} from '@angular/forms';
import {SettingsService} from '../../services/settings.service';
import {Router} from '@angular/router';

@Component({
  selector: 'app-setup',
  templateUrl: './setup.component.html',
  styleUrls: ['./setup.component.scss']
})
export class SetupComponent implements OnInit {

  setupFormGroup = new FormGroup({
    databaseHost: new FormControl('', [Validators.required]),
    databasePort: new FormControl(5432, [Validators.required]),
    databaseUser: new FormControl('', [Validators.required]),
    databasePassword: new FormControl('', [Validators.required]),
    databaseName: new FormControl('', [Validators.required]),
    username: new FormControl('', [Validators.required]),
    password: new FormControl('', [Validators.required]),
  });

  constructor(private settingsService: SettingsService,
              private router: Router) {
  }

  ngOnInit(): void {
  }

  saveSettings() {
    this.settingsService.post(this.setupFormGroup.value)
      .subscribe(() => this.router.navigate(['']));
  }
}
