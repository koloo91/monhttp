import {NgModule} from '@angular/core';
import {Routes, RouterModule} from '@angular/router';
import {ServicesComponent} from './views/services/services.component';
import {CreateServiceComponent} from './views/create-service/create-service.component';
import {BaseComponent} from './views/base/base.component';
import {EditServiceComponent} from './views/edit-service/edit-service.component';
import {DashboardComponent} from './views/dashboard/dashboard.component';
import {ServiceDetailsComponent} from './views/service-details/service-details.component';
import {SettingsComponent} from './views/settings/settings.component';
import {SetupComponent} from './views/setup/setup.component';
import {LoginComponent} from './views/login/login.component';
import {IsLoggedInGuard} from './guards/is-logged-in.guard';
import {IsSetupGuard} from './guards/is-setup.guard';

const routes: Routes = [
  {
    path: '',
    component: BaseComponent,
    children: [
      {
        path: '',
        component: DashboardComponent
      },
      {
        path: 'services/create',
        component: CreateServiceComponent
      },
      {
        path: 'services/edit/:id',
        component: EditServiceComponent
      },
      {
        path: 'services/:id',
        component: ServiceDetailsComponent
      },
      {
        path: 'services',
        component: ServicesComponent
      },
      {
        path: 'settings',
        component: SettingsComponent
      }
    ]
  },
  {
    path: 'setup',
    component: SetupComponent,
    canActivate: [IsSetupGuard]
  },
  {
    path: 'login',
    component: LoginComponent,
    canActivate: [IsLoggedInGuard]
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {
}
