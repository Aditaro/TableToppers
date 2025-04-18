import { Injectable } from '@angular/core';
import { CanActivate, ActivatedRouteSnapshot, RouterStateSnapshot, UrlTree, Router } from '@angular/router';
import { Observable } from 'rxjs';

interface User {
  id: number;
  username: string;
  email: string;
  role: 'customer' | 'manager' | 'admin';
}

@Injectable({
  providedIn: 'root'
})
export class RoleGuard implements CanActivate {

  constructor(private router: Router) {}

  canActivate(
    route: ActivatedRouteSnapshot,
    state: RouterStateSnapshot): Observable<boolean | UrlTree> | Promise<boolean | UrlTree> | boolean | UrlTree {
    const expectedRoles = route.data['expectedRoles'] as Array<string>;
    const userString = localStorage.getItem('user');

    if (!userString) {
      // Should be handled by AuthGuard first, but as a fallback
      this.router.navigate(['/login'], { queryParams: { returnUrl: state.url } });
      return false;
    }

    const user: User = JSON.parse(userString);

    if (!expectedRoles || expectedRoles.length === 0) {
      // No specific role required, allow access
      return true;
    }

    const hasRole = expectedRoles.some(role => user.role === role);

    if (hasRole) {
      // User has the required role
      return true;
    } else {
      // User does not have the required role, redirect (e.g., to home or a 'forbidden' page)
      // For simplicity, redirecting to restaurants list
      this.router.navigate(['/restaurants']);
      return false;
    }
  }
}