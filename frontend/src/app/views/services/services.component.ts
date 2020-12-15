import {Component, OnInit} from '@angular/core';
import {ServiceService} from '../../services/service.service';
import {Observable} from 'rxjs';
import {Service} from '../../models/service.model';
import {Router} from '@angular/router';

@Component({
  selector: 'app-services',
  templateUrl: './services.component.html',
  styleUrls: ['./services.component.scss']
})
export class ServicesComponent implements OnInit {

  displayedColumns: string[] = ['name', 'status', 'visibility', 'failures', 'actions'];

  dataSource$: Observable<Service[]>;

  constructor(private serviceService: ServiceService,
              private router: Router) {
  }

  ngOnInit(): void {
    this.dataSource$ = this.serviceService.list();
    this.serviceService.list().subscribe(console.log, console.log);
  }

  onCreateClick() {
    this.router.navigate(['services', 'create']);
  }
}
