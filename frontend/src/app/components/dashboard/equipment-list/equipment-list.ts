import { Component, computed, input, output, signal } from '@angular/core';
import type { Equipment } from '../../../models/equipment';
import type { Maintenance } from '../../../models/maintenance';
import { MaintenanceHistoryComponent } from '../maintenance-history/maintenance-history';

export interface EquipmentStatus {
  key: 'overdue' | 'due-soon' | 'on-track';
  label: string;
  detail: string;
}

export interface EquipmentListRow {
  equipment: Equipment;
  intervalHours: number;
  meterPercent: number;
  status: EquipmentStatus;
}

@Component({
  selector: 'app-equipment-list',
  templateUrl: './equipment-list.html',
  styleUrls: ['../dashboard-shared.css', './equipment-list.css'],
  host: {
    '(window:resize)': 'updateViewport()',
  },
  imports: [MaintenanceHistoryComponent],
})
export class EquipmentListComponent {
  readonly rows = input.required<EquipmentListRow[]>();
  readonly loading = input(false);
  readonly saving = input(false);
  readonly canCreate = input(false);
  readonly historyFor = input<string | null>(null);
  readonly history = input<Maintenance[]>([]);
  readonly historyLoading = input(false);
  readonly historyError = input('');

  readonly createEquipment = output<void>();
  readonly editEquipment = output<Equipment>();
  readonly registerMaintenance = output<Equipment>();
  readonly deleteEquipment = output<Equipment>();
  readonly showHistory = output<Equipment>();
  readonly closeHistory = output<void>();

  protected readonly searchTerm = signal('');
  protected readonly statusFilter = signal<EquipmentStatus['key'] | 'all'>('all');
  protected readonly typeFilter = signal('all');
  protected readonly currentPage = signal(1);
  protected readonly compactView = signal(this.matchesCompactViewport());
  protected readonly pageSize = 10;
  protected readonly typeOptions = computed(() => Array.from(
    new Set(this.rows().map((row) => row.equipment.type).filter(Boolean)),
  ).sort((left, right) => left.localeCompare(right, 'es')));
  protected readonly filteredRows = computed(() => {
    const searchTerm = this.searchTerm().trim().toLocaleLowerCase('es');
    const statusFilter = this.statusFilter();
    const typeFilter = this.typeFilter();
    return this.rows().filter((row) => {
      const matchesText = searchTerm.length === 0
        || row.equipment.code.toLocaleLowerCase('es').includes(searchTerm)
        || row.equipment.type.toLocaleLowerCase('es').includes(searchTerm);
      const matchesStatus = statusFilter === 'all' || row.status.key === statusFilter;
      const matchesType = typeFilter === 'all' || row.equipment.type === typeFilter;
      return matchesText && matchesStatus && matchesType;
    });
  });
  protected readonly pageCount = computed(() => Math.max(1, Math.ceil(this.filteredRows().length / this.pageSize)));
  protected readonly safeCurrentPage = computed(() => Math.min(this.currentPage(), this.pageCount()));
  protected readonly visibleRows = computed(() => {
    const start = (this.safeCurrentPage() - 1) * this.pageSize;
    return this.filteredRows().slice(start, start + this.pageSize);
  });

  updateViewport(): void {
    const compact = this.matchesCompactViewport();
    if (compact !== this.compactView()) this.compactView.set(compact);
  }

  private matchesCompactViewport(): boolean {
    return typeof window !== 'undefined'
      && typeof window.matchMedia === 'function'
      && window.matchMedia('(max-width: 760px)').matches;
  }

  protected updateSearch(event: Event): void {
    const input = event.currentTarget;
    if (!(input instanceof HTMLInputElement)) return;
    this.searchTerm.set(input.value);
    this.currentPage.set(1);
  }

  protected updateStatusFilter(event: Event): void {
    const select = event.currentTarget;
    if (!(select instanceof HTMLSelectElement)) return;
    const value = select.value;
    if (value !== 'all' && value !== 'overdue' && value !== 'due-soon' && value !== 'on-track') return;
    this.statusFilter.set(value);
    this.currentPage.set(1);
  }

  protected updateTypeFilter(event: Event): void {
    const select = event.currentTarget;
    if (!(select instanceof HTMLSelectElement)) return;
    this.typeFilter.set(select.value);
    this.currentPage.set(1);
  }

  protected clearFilters(): void {
    this.searchTerm.set('');
    this.statusFilter.set('all');
    this.typeFilter.set('all');
    this.currentPage.set(1);
  }

  protected previousPage(): void {
    this.currentPage.set(Math.max(1, this.safeCurrentPage() - 1));
  }

  protected nextPage(): void {
    this.currentPage.set(Math.min(this.pageCount(), this.safeCurrentPage() + 1));
  }
}
