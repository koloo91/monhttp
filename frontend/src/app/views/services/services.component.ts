import {Component, OnInit} from '@angular/core';
import {ServiceService} from '../../services/service.service';
import {Observable} from 'rxjs';
import {Service} from '../../models/service.model';
import {Router} from '@angular/router';
import {ErrorService} from '../../services/error.service';
import {catchError, switchAll, switchMap} from 'rxjs/operators';
import {ApiError} from '../../models/api-error.model';
import {MatDialog} from '@angular/material/dialog';
import {ConfirmServiceDeleteDialogComponent} from '../../components/dialogs/confirm-service-delete-dialog/confirm-service-delete-dialog.component';

@Component({
  selector: 'app-services',
  templateUrl: './services.component.html',
  styleUrls: ['./services.component.scss']
})
export class ServicesComponent implements OnInit {

  displayedColumns: string[] = ['name', 'status', 'visibility', 'failures', 'actions'];

  dataSource$: Observable<Service[]>;

  constructor(private serviceService: ServiceService,
              private errorService: ErrorService,
              private router: Router,
              public dialog: MatDialog) {
  }

  ngOnInit(): void {
    this.loadServices();
  }

  loadServices() {
    this.dataSource$ = this.serviceService.list();
  }

  onCreateClick() {
    this.router.navigate(['services', 'create']);
  }

  onDeleteClick(service: Service): void {
    console.log('deleteClicked');
    const dialogRef = this.dialog.open(ConfirmServiceDeleteDialogComponent, {
      data: service
    });

    dialogRef.afterClosed().subscribe(deleteService => {
      if (!deleteService) {
        return;
      }

      this.dataSource$ = this.serviceService.delete(service.id)
        .pipe(
          switchMap(() => this.serviceService.list())
        );
    })
  }

  onEditClick(id: string) {
    console.log('editClicked');
    this.router.navigate(['services', 'edit', id]);
  }
}
