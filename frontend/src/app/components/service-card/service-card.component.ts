import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {Service} from '../../models/service.model';
import {CheckService} from '../../services/check.service';
import {Observable} from 'rxjs';
import {Average} from '../../models/average.model';
import {map, tap} from 'rxjs/operators';
import {Check} from '../../models/check.model';
import {IsOnline} from '../../models/is-online.model';

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
  isOnline$: Observable<IsOnline>;

  checks: Check[] = [];
  chartData: any = [];
  chartColorScheme = {
    domain: ['#000']
  };

  cardWidth: number;

  isLoading = false;

  constructor(private checkService: CheckService) {
  }

  ngOnInit(): void {
    this.isLoading = true;

    this.average$ = this.checkService.average(this.service.id);
    this.isOnline$ = this.checkService.isOnline(this.service.id)
      .pipe(
        tap(isOnline => {
          if (isOnline.online) {
            this.chartColorScheme = {
              domain: ['#28a745']
            };
          } else {
            this.chartColorScheme = {
              domain: ['#dd3545']
            };
          }
        })
      );

    const yesterday = new Date();
    yesterday.setDate(yesterday.getDate() - 1);

    this.checkService.list(this.service.id, yesterday.toISOString(), new Date().toISOString(), 6)
      .pipe(
        map(checks => checks.reverse()),
        tap(() => this.isLoading = false)
      )
      .subscribe(checks => {
        this.checks = checks;
        this.chartData = [{
          name: 'Latency in ms',
          series: checks.map(check => {
            return {name: new Date(check.createdAt).toLocaleTimeString(), value: check.latencyInMs};
          })
        }]

      }, console.log)
  }

  onCardClicked() {
    this.cardClicked.emit(this.service.id);
  }
}
