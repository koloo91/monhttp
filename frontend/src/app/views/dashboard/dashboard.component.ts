import {Component, OnInit} from '@angular/core';
import {ServiceService} from '../../services/service.service';
import {Observable} from 'rxjs';
import {Service} from '../../models/service.model';
import {Router} from '@angular/router';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent implements OnInit {

  services$: Observable<Service[]>;

  constructor(private serviceService: ServiceService,
              private router: Router) {
  }

  ngOnInit(): void {
    this.services$ = this.serviceService.list();
  }

  showServiceDetails(serviceId: string): void {
    this.router.navigate(['services', serviceId])
  }
}
