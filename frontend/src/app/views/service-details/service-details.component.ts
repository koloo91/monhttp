import {Component, OnInit} from '@angular/core';
import {ServiceService} from '../../services/service.service';
import {CheckService} from '../../services/check.service';
import {ErrorService} from '../../services/error.service';
import {ActivatedRoute} from '@angular/router';
import {map, switchMap} from 'rxjs/operators';
import {combineLatest, Observable} from 'rxjs';
import {Service} from '../../models/service.model';
import {FormControl, FormGroup} from '@angular/forms';

@Component({
  selector: 'app-service-details',
  templateUrl: './service-details.component.html',
  styleUrls: ['./service-details.component.scss']
})
export class ServiceDetailsComponent implements OnInit {


  dateTimeRangeFormGroup = new FormGroup({
    fromDateTime: new FormControl(new Date()),
    toDateTime: new FormControl(new Date())
  });

  service$: Observable<Service>;

  chartData: any = [];

  constructor(private serviceService: ServiceService,
              private checkService: CheckService,
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

    const now = new Date();

    // this.route.params
    //   .pipe(
    //     map(params => params['id'] as string),
    //     switchMap(id => this.checkService.list(id, yesterday.toISOString(), now.toISOString()))
    //   ).subscribe(checks => {
    //   this.chartData = [{
    //     name: 'Latency in ms', series: checks.map(check => {
    //       return {name: new Date(check.createdAt).toLocaleString(), value: check.latencyInMs};
    //     })
    //   }]
    // });
    this.loadChartData();
    this.dateTimeRangeFormGroup.get('fromDateTime').setValue(yesterday);

    this.dateTimeRangeFormGroup.valueChanges.subscribe(({fromDateTime, toDateTime}) => {
      console.log(fromDateTime);
      console.log(toDateTime);
    });
  }

  loadChartData() {
    combineLatest([this.route.params, this.dateTimeRangeFormGroup.valueChanges])
      .pipe(
        map(([params, formValues]) => [params['id'] as string, formValues]),
        switchMap(([id, {
          fromDateTime,
          toDateTime
        }]) => this.checkService.list(id, fromDateTime.toISOString(), toDateTime.toISOString())),
        map(checks => checks.reverse())
      ).subscribe(checks => {
      this.chartData = [{
        name: 'Latency in ms', series: checks.map(check => {
          return {name: new Date(check.createdAt).toLocaleString(), value: check.latencyInMs};
        })
      }]
    });
  }
}
