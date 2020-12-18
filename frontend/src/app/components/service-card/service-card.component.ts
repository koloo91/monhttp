import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {Service} from '../../models/service.model';
import {CheckService} from '../../services/check.service';
import {Observable} from 'rxjs';
import {Average} from '../../models/average.model';
import {map} from 'rxjs/operators';

@Component({
  selector: 'app-service-card',
  templateUrl: './service-card.component.html',
  styleUrls: ['./service-card.component.scss']
})
export class ServiceCardComponent implements OnInit {

  @Input()
  service: Service;

  @Output()
  cardClicked: EventEmitter<string> = new EventEmitter<string>();

  average$: Observable<Average>;

  chartData: any = [];

  constructor(private checkService: CheckService) {
  }

  ngOnInit(): void {
    this.average$ = this.checkService.average(this.service.id);

    const yesterday = new Date();
    yesterday.setDate(yesterday.getDate() - 1);


    this.checkService.list(this.service.id, yesterday.toISOString(), new Date().toISOString())
      .pipe(
        map(checks => checks.reverse())
      )
      .subscribe(data => {
        this.chartData = [{
          name: 'Latency in ms',
          series: data.map(check => {
            return {name: new Date(check.createdAt).toLocaleTimeString(), value: check.latencyInMs};
          })
        }]

      }, console.log)

  }

  onCardClicked() {
    this.cardClicked.emit(this.service.id);
  }
}
