import {Component, Input, OnInit} from '@angular/core';
import {Service} from '../../models/service.model';
import {CheckService} from '../../services/check.service';
import {Observable} from 'rxjs';
import {Average} from '../../models/average.model';

@Component({
  selector: 'app-service-card',
  templateUrl: './service-card.component.html',
  styleUrls: ['./service-card.component.scss']
})
export class ServiceCardComponent implements OnInit {

  @Input()
  service: Service;

  average$: Observable<Average>;

  constructor(private checkService: CheckService) {
  }

  ngOnInit(): void {
    this.average$ = this.checkService.average(this.service.id);
  }

}
