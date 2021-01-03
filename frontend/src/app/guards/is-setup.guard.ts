import {Injectable} from '@angular/core';
import {CanActivate, ActivatedRouteSnapshot, RouterStateSnapshot, UrlTree, Router} from '@angular/router';
import {Observable} from 'rxjs';
import {SettingsService} from '../services/settings.service';
import {map} from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class IsSetupGuard implements CanActivate {
  constructor(private settingsService: SettingsService,
              private router: Router) {
  }

  canActivate(
    route: ActivatedRouteSnapshot,
    state: RouterStateSnapshot): Observable<boolean | UrlTree> | Promise<boolean | UrlTree> | boolean | UrlTree {
    return this.settingsService.setup()
      .pipe(
        map(setup => {
          if (setup.isSetup) {
            this.router.navigate(['']);
            return false;
          }
          return true;
        })
      );
  }

}
