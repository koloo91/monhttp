import {Component, OnInit} from '@angular/core';
import {ServiceService} from '../../services/service.service';
import {CheckService} from '../../services/check.service';
import {ErrorService} from '../../services/error.service';
import {ActivatedRoute} from '@angular/router';
import {map, switchMap} from 'rxjs/operators';
import {BehaviorSubject, combineLatest, Observable} from 'rxjs';
import {Service} from '../../models/service.model';
import {FormControl, FormGroup} from '@angular/forms';
import {FailureService} from '../../services/failure.service';
import {Failure} from '../../models/failure.model';
import {Check} from '../../models/check.model';
import {MatSelectChange} from '@angular/material/select';
import {PageEvent} from '@angular/material/paginator';

@Component({
  selector: 'app-service-details',
  templateUrl: './service-details.component.html',
  styleUrls: ['./service-details.component.scss']
})
export class ServiceDetailsComponent implements OnInit {

  displayedColumns: string[] = ['reason', 'date'];

  dateTimeRangeFormGroup = new FormGroup({
    fromDateTime: new FormControl(new Date()),
    toDateTime: new FormControl(new Date()),
    timeInterval: new FormControl('')
  });

  service$: Observable<Service>;
  checks: Check[] = [];
  failures: Failure[] = [];

  chartData: any;

  timeIntervals = [
    {
      name: '1 minute',
      get: () => this.setTimeInterval(0, 0, 1)
    },
    {
      name: '5 minutes',
      get: () => this.setTimeInterval(0, 0, 5)
    },
    {
      name: '15 minutes',
      get: () => this.setTimeInterval(0, 0, 15)
    },
    {
      name: '30 minutes',
      get: () => this.setTimeInterval(0, 0, 30)
    },
    {
      name: '1 hour',
      get: () => this.setTimeInterval(0, 1, 0)
    },
    {
      name: '3 hours',
      get: () => this.setTimeInterval(0, 3, 0)
    },
    {
      name: '6 hours',
      get: () => this.setTimeInterval(0, 6, 0)
    },
    {
      name: '12 hours',
      get: () => this.setTimeInterval(0, 12, 0)
    },
    {
      name: '1 day',
      get: () => this.setTimeInterval(1, 0, 0)
    },
    {
      name: '7 days',
      get: () => this.setTimeInterval(7, 0, 0)
    },
    {
      name: '14 days',
      get: () => this.setTimeInterval(14, 0, 0)
    },
    {
      name: '30 days',
      get: () => this.setTimeInterval(30, 0, 0)
    }
  ];

  selectedInterval: any = this.timeIntervals[2];

  failureItemsPageSize = 10;
  failureItemsPerPage = [5, 10, 25, 50];
  failureItemsLength = 0;

  failurePaginatorEventSubject = new BehaviorSubject<PageEvent>({length: 0, pageSize: this.failureItemsPageSize, pageIndex: 0});
  failurePaginatorEvent$: Observable<PageEvent>;

  constructor(private serviceService: ServiceService,
              private checkService: CheckService,
              private failureService: FailureService,
              private errorService: ErrorService,
              private route: ActivatedRoute) {
    this.failurePaginatorEvent$ = this.failurePaginatorEventSubject.asObservable();
  }

  ngOnInit(): void {
    this.service$ = this.route.params
      .pipe(
        map(params => params['id'] as string),
        switchMap(id => this.serviceService.get(id))
      );

    const yesterday = new Date();
    yesterday.setDate(yesterday.getDate() - 1);

    this.loadChartData();
    combineLatest([this.route.params, this.dateTimeRangeFormGroup.valueChanges, this.failurePaginatorEvent$])
      .pipe(
        map(([params, formValues, pageEvent]) => [params['id'], formValues, pageEvent]),
        switchMap(([id, {
          fromDateTime,
          toDateTime
        }, pageEvent]) => this.failureService.list(id, fromDateTime.toISOString(), toDateTime.toISOString(), pageEvent.pageSize, pageEvent.pageIndex))
      )
      .subscribe(wrapper => {
        this.failures = wrapper.data;
        this.failureItemsLength = wrapper.totalCount;
      }, console.log);

    this.dateTimeRangeFormGroup.get('fromDateTime').setValue(yesterday);
  }

  loadChartData() {
    combineLatest([this.route.params, this.dateTimeRangeFormGroup.valueChanges])
      .pipe(
        map(([params, formValues]) => [params['id'] as string, formValues]),
        switchMap(([id, {
          fromDateTime,
          toDateTime
        }]) => this.checkService.list(id, fromDateTime.toISOString(), toDateTime.toISOString())),
        map(checks => checks.reverse()),
      )
      .subscribe(checks => {
        this.checks = checks;
        this.chartData = [{
          name: 'Latency in ms', series: checks.map(check => {
            return {name: new Date(check.createdAt).toLocaleString(), value: check.latencyInMs};
          })
        }]
      });
  }

  setTimeInterval(days: number, hours: number, minutes: number) {
    this.dateTimeRangeFormGroup.get('toDateTime').setValue(new Date());

    const fromDate = new Date();
    fromDate.setDate(fromDate.getDate() - days);
    fromDate.setHours(fromDate.getHours() - hours, fromDate.getMinutes() - minutes);
    this.dateTimeRangeFormGroup.get('fromDateTime').setValue(fromDate);
  }

  onSelectionChange($event: MatSelectChange) {
    this.timeIntervals.find(timeInterval => timeInterval.name === $event.value)?.get();
  }

  onFailurePageChanged(pageEvent: PageEvent) {
    this.failurePaginatorEventSubject.next(pageEvent);
  }
}
