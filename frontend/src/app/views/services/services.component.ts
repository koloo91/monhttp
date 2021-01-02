import {Component, OnInit} from '@angular/core';
import {ServiceService} from '../../services/service.service';
import {forkJoin, Observable, of} from 'rxjs';
import {Service} from '../../models/service.model';
import {Router} from '@angular/router';
import {ErrorService} from '../../services/error.service';
import {map, switchMap} from 'rxjs/operators';
import {MatDialog} from '@angular/material/dialog';
import {ConfirmServiceDeleteDialogComponent} from '../../components/dialogs/confirm-service-delete-dialog/confirm-service-delete-dialog.component';
import {CheckService} from '../../services/check.service';
import {FailureService} from '../../services/failure.service';

interface ServiceWithStatusAndFailures extends Service {
  isOnline: boolean;
  failureCount: number;
}

@Component({
  selector: 'app-services',
  templateUrl: './services.component.html',
  styleUrls: ['./services.component.scss']
})
export class ServicesComponent implements OnInit {

  displayedColumns: string[] = ['name', 'status', /*'visibility',*/ 'failures', 'actions'];

  dataSource$: Observable<ServiceWithStatusAndFailures[]>;

  constructor(private serviceService: ServiceService,
              private checkService: CheckService,
              private failureService: FailureService,
              private errorService: ErrorService,
              private router: Router,
              public dialog: MatDialog) {
  }

  ngOnInit(): void {
    this.loadServices();
  }

  loadServices() {
    this.dataSource$ = this.serviceService.list()
      .pipe(
        map(services => services.map(service => service as ServiceWithStatusAndFailures)),
        switchMap(services => forkJoin([...services.map(service => {
          return this.checkService.isOnline(service.id)
            .pipe(
              map(isOnline => {
                service.isOnline = isOnline.online;
                return service;
              })
            )
        })])),
        switchMap(services => forkJoin([...services.map(service => {
          const yesterday = new Date();
          yesterday.setDate(yesterday.getDate() - 1);
          return this.failureService.count(service.id, yesterday.toISOString(), new Date().toISOString())
            .pipe(
              map(failureCount => {
                service.failureCount = failureCount.count;
                return service;
              })
            );
        })]))
      );
  }

  onCreateClick() {
    this.router.navigate(['services', 'create']);
  }

  onDeleteClick(service: Service): void {
    const dialogRef = this.dialog.open(ConfirmServiceDeleteDialogComponent, {
      data: service
    });

    dialogRef.afterClosed().subscribe(deleteService => {
      if (!deleteService) {
        return;
      }

      this.dataSource$ = this.serviceService.delete(service.id)
        .pipe(
          switchMap(() => this.serviceService.list()),
          map(services => services as ServiceWithStatusAndFailures[])
        );
    })
  }

  onEditClick(id: string): void {
    this.router.navigate(['services', 'edit', id]);
  }

  onShowChartClick(id: string): void {
    this.router.navigate(['services', id]);
  }
}
