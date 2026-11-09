import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: '',
    loadComponent: () => import('./components/dashboard/dashboard').then((module) => module.Dashboard),
  },
  { path: '**', redirectTo: '' },
];
