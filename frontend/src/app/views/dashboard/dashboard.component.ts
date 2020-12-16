import {Component, OnInit} from '@angular/core';
import {ServiceService} from '../../services/service.service';
import {Observable} from 'rxjs';
import {Service} from '../../models/service.model';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent implements OnInit {

  services$: Observable<Service[]>;

  constructor(private serviceService: ServiceService) {
  }

  ngOnInit(): void {
    this.services$ = this.serviceService.list();
  }

}
