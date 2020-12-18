import {Component, OnInit} from '@angular/core';
import {ServiceService} from '../../services/service.service';
import {CheckService} from '../../services/check.service';
import {ErrorService} from '../../services/error.service';
import {ActivatedRoute} from '@angular/router';
import {map, switchMap} from 'rxjs/operators';
import {combineLatest, Observable} from 'rxjs';
import {Service} from '../../models/service.model';
import {FormControl, FormGroup} from '@angular/forms';
import {FailureService} from '../../services/failure.service';
import {Failure} from '../../models/failure.model';
import {Check} from '../../models/check.model';

@Component({
  selector: 'app-service-details',
  templateUrl: './service-details.component.html',
  styleUrls: ['./service-details.component.scss']
})
export class ServiceDetailsComponent implements OnInit {

  displayedColumns: string[] = ['reason', 'date'];


  dateTimeRangeFormGroup = new FormGroup({
    fromDateTime: new FormControl(new Date()),
    toDateTime: new FormControl(new Date())
  });

  service$: Observable<Service>;
  checks: Check[] = [];
  failures: Failure[] = [];

  chartData: any;

  constructor(private serviceService: ServiceService,
              private checkService: CheckService,
              private failureService: FailureService,
              private errorService: ErrorService,
              private route: ActivatedRoute) {
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
    combineLatest([this.route.params, this.dateTimeRangeFormGroup.valueChanges])
      .pipe(
        map(([params, formValues]) => [params['id'], formValues]),
        switchMap(([id, {
          fromDateTime,
          toDateTime
        }]) => this.failureService.list(id, fromDateTime.toISOString(), toDateTime.toISOString()))
      )
      .subscribe(data => this.failures = data, console.log);

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
}
