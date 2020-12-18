import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';

import {AppRoutingModule} from './app-routing.module';
import {AppComponent} from './app.component';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {ServicesComponent} from './views/services/services.component';
import {CreateServiceComponent} from './views/create-service/create-service.component';
import {HttpClientModule} from '@angular/common/http';
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
import { DashboardComponent } from './views/dashboard/dashboard.component';
import { ServiceCardComponent } from './components/service-card/service-card.component';
import {MatDividerModule} from '@angular/material/divider';
import {NgApexchartsModule} from 'ng-apexcharts';
import { ServiceDetailsComponent } from './views/service-details/service-details.component';

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
    ServiceDetailsComponent
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
    MatButtonModule,
    MatIconModule,
    MatSnackBarModule,
    MatTableModule,
    MatDialogModule,
    MatDividerModule,
    NgApexchartsModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule {
}
