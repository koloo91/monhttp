import {Component, OnInit} from '@angular/core';
import {ServiceService} from '../../services/service.service';

@Component({
  selector: 'app-services',
  templateUrl: './services.component.html',
  styleUrls: ['./services.component.scss']
})
export class ServicesComponent implements OnInit {

  constructor(private serviceService: ServiceService) {
  }

  ngOnInit(): void {
    this.serviceService.list().subscribe(console.log, console.log);
  }

}
