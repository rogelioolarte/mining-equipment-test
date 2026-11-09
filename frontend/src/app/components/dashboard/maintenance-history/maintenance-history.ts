import { DatePipe } from '@angular/common';
import { Component, input, output } from '@angular/core';
import type { Maintenance } from '../../../models/maintenance';

@Component({
  selector: 'app-maintenance-history',
  templateUrl: './maintenance-history.html',
  styleUrls: ['../dashboard-shared.css', './maintenance-history.css'],
  imports: [DatePipe],
})
export class MaintenanceHistoryComponent {
  readonly equipmentCode = input.required<string>();
  readonly entries = input.required<Maintenance[]>();
  readonly loading = input(false);
  readonly error = input('');
  readonly compact = input(false);
  readonly close = output<void>();
}
