import {Component, OnInit} from '@angular/core';
import {ServiceService} from '../../services/service.service';
import {BehaviorSubject, Observable} from 'rxjs';
import {Service} from '../../models/service.model';
import {Router} from '@angular/router';
import {delay, switchMap, tap} from 'rxjs/operators';
import {PageEvent} from '@angular/material/paginator';
import {Wrapper} from '../../models/wrapper.model';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent implements OnInit {

  isLoading = false;
  wrapper$: Observable<Wrapper<Service>>;

  servicesPageSize = 12;
  servicesPerPage = [3, 12, 24, 48, 96];

  servicesPaginatorEventSubject = new BehaviorSubject<PageEvent>({
    length: 0,
    pageSize: this.servicesPageSize,
    pageIndex: 0
  });
  servicesPaginatorEvent$: Observable<PageEvent>;

  constructor(private serviceService: ServiceService,
              private router: Router) {
    this.servicesPaginatorEvent$ = this.servicesPaginatorEventSubject.asObservable();
  }

  ngOnInit(): void {
    this.wrapper$ = this.servicesPaginatorEvent$
      .pipe(
        delay(0),
        tap(() => this.isLoading = true),
        switchMap(pageEvent => this.serviceService.list(pageEvent.pageSize, pageEvent.pageIndex)),
        tap(() => this.isLoading = false)
      );
  }

  showServiceDetails(serviceId: string): void {
    this.router.navigate(['services', serviceId])
  }

  onServicePageChanged(pageEvent: PageEvent) {
    this.servicesPaginatorEventSubject.next(pageEvent);
  }
}
