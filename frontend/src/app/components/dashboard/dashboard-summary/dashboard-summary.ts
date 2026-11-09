import { Component, input } from '@angular/core';

@Component({
  selector: 'app-dashboard-summary',
  templateUrl: './dashboard-summary.html',
  styleUrl: './dashboard-summary.css',
})
export class DashboardSummary {
  readonly equipmentCount = input.required<number>();
  readonly dueCount = input.required<number>();
  readonly equipmentTypeCount = input.required<number>();
}
