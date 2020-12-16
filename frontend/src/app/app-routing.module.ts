import {NgModule} from '@angular/core';
import {Routes, RouterModule} from '@angular/router';
import {ServicesComponent} from './views/services/services.component';
import {CreateServiceComponent} from './views/create-service/create-service.component';
import {BaseComponent} from './views/base/base.component';
import {EditServiceComponent} from './views/edit-service/edit-service.component';
import {DashboardComponent} from './views/dashboard/dashboard.component';

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
        path: 'services',
        component: ServicesComponent
      }
    ]
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {
}
