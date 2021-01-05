import {Component, OnDestroy, OnInit} from '@angular/core';
import {ServiceService} from '../../services/service.service';
import {forkJoin, Observable, of, Subscription} from 'rxjs';
import {Service} from '../../models/service.model';
import {Router} from '@angular/router';
import {ErrorService} from '../../services/error.service';
import {map, switchMap, tap} from 'rxjs/operators';
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
export class ServicesComponent implements OnInit, OnDestroy {

  isLoading = false;

  displayedColumns: string[] = ['name', 'status', /*'visibility',*/ 'failures', 'actions'];

  dataSource$: Observable<ServiceWithStatusAndFailures[]>;

  subscriptions: Subscription[] = [];

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

  ngOnDestroy(): void {
    this.subscriptions.forEach(subscription => subscription.unsubscribe());
  }

  loadServices() {
    this.isLoading = true;

    this.dataSource$ = this.serviceService.list()
      .pipe(
        map(services => services.map(service => service as ServiceWithStatusAndFailures))
      );

    this.subscriptions.push(
      this.dataSource$
        .subscribe(services => {
          this.dataSource$ = forkJoin(services.map(service => {
            return this.checkService.isOnline(service.id)
              .pipe(
                map(isOnline => {
                  service.isOnline = isOnline.online;
                  return service;
                })
              );
          }))

          this.subscriptions.push(
            this.dataSource$.subscribe(services => {
              const yesterday = new Date();
              yesterday.setDate(yesterday.getDate() - 1);

              this.dataSource$ = forkJoin(services.map(service => {
                return this.failureService.count(service.id, yesterday.toISOString(), new Date().toISOString())
                  .pipe(
                    map(failureCount => {
                      service.failureCount = failureCount.count;
                      return service;
                    }),
                    tap(() => this.isLoading = false)
                  );
              }));
            }));
        }));
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

      this.serviceService.delete(service.id)
        .pipe(
          tap(() => this.loadServices())
        ).subscribe();
    })
  }

  onEditClick(id: string): void {
    this.router.navigate(['services', 'edit', id]);
  }

  onShowChartClick(id: string): void {
    this.router.navigate(['services', id]);
  }
}
