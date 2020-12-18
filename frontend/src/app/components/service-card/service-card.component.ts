import {Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {Service} from '../../models/service.model';
import {CheckService} from '../../services/check.service';
import {Observable} from 'rxjs';
import {Average} from '../../models/average.model';
import {ChartComponent} from 'ng-apexcharts';

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

  @ViewChild('chart') chart: ChartComponent;
  public chartOptions: Partial<any>;

  constructor(private checkService: CheckService) {
  }

  ngOnInit(): void {
    this.average$ = this.checkService.average(this.service.id);

    const yesterday = new Date();
    yesterday.setDate(yesterday.getDate() - 1);

    this.chartOptions = {
      series: [
        {
          name: 'My-series',
          data: []
        }
      ],
      chart: {
        height: 'auto',
        type: 'area'
      },
      title: {
        text: ''
      },
      yaxis: {
        labels: {
          formatter: function (val) {
            return (val / 10000).toFixed();
          }
        }
      },
      xaxis: {
        type: 'datetime'
      }
    };

    this.checkService.list(this.service.id, yesterday.toISOString(), new Date().toISOString())
      .subscribe(data => {
        this.chartOptions.series = [
          {
            name: 'My-series',
            data: data.map(check => [new Date(check.createdAt).getTime(), check.latencyInMs])
          }
        ];

      }, console.log)

  }

  onCardClicked() {
    this.cardClicked.emit(this.service.id);
  }
}
