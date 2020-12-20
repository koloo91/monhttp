import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';

import {AppRoutingModule} from './app-routing.module';
import {AppComponent} from './app.component';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {ServicesComponent} from './views/services/services.component';
import {CreateServiceComponent} from './views/create-service/create-service.component';
import {HTTP_INTERCEPTORS, HttpClientModule} from '@angular/common/http';
import {BaseComponent} from './views/base/base.component';
import {MatToolbarModule} from '@angular/material/toolbar';
import {MatCardModule} from '@angular/material/card';
import {MatInputModule} from '@angular/material/input';
import {MatFormFieldModule} from '@angular/material/form-field';
import {ReactiveFormsModule} from '@angular/forms';
import {MatOptionModule} from '@angular/material/core';
import {MatSelectModule} from '@angular/material/select';
import {MatSliderModule} from '@angular/material/slider';
import {MatSlideToggleModule} from '@angular/material/slide-toggle';
import {MatButtonModule} from '@angular/material/button';
import {MatIconModule} from '@angular/material/icon';
import {MatSnackBarModule} from '@angular/material/snack-bar';
import {MatTableModule} from '@angular/material/table';
import {ConfirmServiceDeleteDialogComponent} from './components/dialogs/confirm-service-delete-dialog/confirm-service-delete-dialog.component';
import {MatDialogModule} from '@angular/material/dialog';
import {EditServiceComponent} from './views/edit-service/edit-service.component';
import {DashboardComponent} from './views/dashboard/dashboard.component';
import {ServiceCardComponent} from './components/service-card/service-card.component';
import {MatDividerModule} from '@angular/material/divider';
import {ServiceDetailsComponent} from './views/service-details/service-details.component';
import {NgxChartsModule} from '@swimlane/ngx-charts';
import {NgxMatDatetimePickerModule, NgxMatNativeDateModule} from '@angular-material-components/datetime-picker';
import {MatDatepickerModule} from '@angular/material/datepicker';
import {SettingsComponent} from './views/settings/settings.component';
import {MatListModule} from '@angular/material/list';
import {ApiErrorInterceptor} from './http_interceptor/api-error.interceptor';
import { SetupComponent } from './views/setup/setup.component';

@NgModule({
  declarations: [
    AppComponent,
    ServicesComponent,
    CreateServiceComponent,
    BaseComponent,
    ConfirmServiceDeleteDialogComponent,
    EditServiceComponent,
    DashboardComponent,
    ServiceCardComponent,
    ServiceDetailsComponent,
    SettingsComponent,
    SetupComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    BrowserAnimationsModule,
    HttpClientModule,
    ReactiveFormsModule,
    MatToolbarModule,
    MatCardModule,
    MatInputModule,
    MatFormFieldModule,
    MatSelectModule,
    MatOptionModule,
    MatSliderModule,
    MatSlideToggleModule,
    MatListModule,
    MatButtonModule,
    MatIconModule,
    MatSnackBarModule,
    MatTableModule,
    MatDialogModule,
    MatDividerModule,
    NgxChartsModule,
    NgxMatDatetimePickerModule,
    MatDatepickerModule,
    NgxMatNativeDateModule
  ],
  providers: [
    {provide: HTTP_INTERCEPTORS, useClass: ApiErrorInterceptor, multi: true}],
  bootstrap: [AppComponent]
})
export class AppModule {
}
