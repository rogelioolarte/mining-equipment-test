import { Component, input, output } from '@angular/core';
import type { Equipment } from '../../../models/equipment';

@Component({
  selector: 'app-delete-confirmation',
  templateUrl: './delete-confirmation.html',
  styleUrls: ['../dashboard-shared.css', './delete-confirmation.css'],
})
export class DeleteConfirmationComponent {
  readonly candidate = input.required<Equipment>();
  readonly saving = input(false);
  readonly loading = input(false);
  readonly confirmed = output<void>();
  readonly cancelled = output<void>();
}
