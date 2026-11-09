import { Component, computed, DestroyRef, inject, signal } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import type { Equipment, EquipmentRequest } from '../../models/equipment';
import type { EquipmentType } from '../../models/equipment-type';
import type { Maintenance } from '../../models/maintenance';
import { EquipmentAPI } from '../../services/equipment-api';
import { EquipmentTypeAPI } from '../../services/equipment-type-api';
import { MaintenanceAPI } from '../../services/maintenance-api';
import { DashboardSummary } from './dashboard-summary/dashboard-summary';
import { DeleteConfirmationComponent } from './delete-confirmation/delete-confirmation';
import { EquipmentFormComponent } from './equipment-form/equipment-form';
import type { EquipmentFormState, EquipmentFormSubmission } from './equipment-form/equipment-form';
import { EquipmentListComponent } from './equipment-list/equipment-list';
import type { EquipmentListRow, EquipmentStatus } from './equipment-list/equipment-list';
import { MaintenanceFormComponent } from './maintenance-form/maintenance-form';
import type { MaintenanceFormState, MaintenanceFormSubmission } from './maintenance-form/maintenance-form';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.html',
  styleUrl: './dashboard.css',
  imports: [
    DashboardSummary,
    DeleteConfirmationComponent,
    EquipmentFormComponent,
    EquipmentListComponent,
    MaintenanceFormComponent,
  ],
})
export class Dashboard {
  private readonly equipmentApi = inject(EquipmentAPI);
  private readonly equipmentTypeApi = inject(EquipmentTypeAPI);
  private readonly maintenanceApi = inject(MaintenanceAPI);

  protected readonly equipment = signal<Equipment[]>([]);
  protected readonly equipmentTypes = signal<EquipmentType[]>([]);
  protected readonly loading = signal(true);
  protected readonly saving = signal(false);
  protected readonly error = signal('');
  protected readonly notice = signal('');
  protected readonly equipmentForm = signal<EquipmentFormState | null>(null);
  protected readonly maintenanceForm = signal<MaintenanceFormState | null>(null);
  protected readonly historyFor = signal<string | null>(null);
  protected readonly history = signal<Maintenance[]>([]);
  protected readonly historyLoading = signal(false);
  protected readonly historyError = signal('');
  protected readonly deleteCandidate = signal<Equipment | null>(null);
  protected readonly equipmentRows = computed<EquipmentListRow[]>(() => this.equipment().map((equipment) => {
    const intervalHours = this.typeInterval(equipment);
    const status = this.statusFor(equipment);
    return {
      equipment,
      intervalHours,
      meterPercent: intervalHours > 0
        ? Math.min(100, Math.max(0, equipment.hoursSinceMaintenance / intervalHours * 100))
        : 0,
      status,
    };
  }));
  protected readonly dueCount = computed(() => this.equipmentRows().filter((row) => row.status.key !== 'on-track').length);

  private readonly destroyRef = inject(DestroyRef);
  private loadRequestId = 0;
  private historyRequestId = 0;

  constructor() {
    this.load();
  }

  protected load(): void {
    const requestId = ++this.loadRequestId;
    this.loading.set(true);
    this.error.set('');
    let equipmentLoaded = false;
    let typesLoaded = false;
    let equipmentError = '';
    let typesError = '';
    const finish = (): void => {
      if (requestId !== this.loadRequestId || !equipmentLoaded || !typesLoaded) return;
      this.loading.set(false);
      this.error.set([equipmentError, typesError].filter(Boolean).join(' '));
    };

    this.equipmentApi.getAll().pipe(takeUntilDestroyed(this.destroyRef)).subscribe({
      next: (items) => {
        if (requestId !== this.loadRequestId) return;
        this.equipment.set(items);
        equipmentLoaded = true;
        finish();
      },
      error: () => {
        if (requestId !== this.loadRequestId) return;
        equipmentError = 'No se pudo cargar la maquinaria. Comprueba la conexión e inténtalo de nuevo.';
        equipmentLoaded = true;
        finish();
      },
    });
    this.equipmentTypeApi.getAll().pipe(takeUntilDestroyed(this.destroyRef)).subscribe({
      next: (items) => {
        if (requestId !== this.loadRequestId) return;
        this.equipmentTypes.set(items);
        typesLoaded = true;
        finish();
      },
      error: () => {
        if (requestId !== this.loadRequestId) return;
        typesError = 'No se pudieron cargar los tipos de maquinaria. Inténtalo de nuevo antes de añadir o editar equipos.';
        typesLoaded = true;
        finish();
      },
    });
  }

  protected startCreate(): void {
    this.notice.set('');
    this.equipmentForm.set({ code: '', typeId: '', currentHourMeter: '' });
  }

  protected startEdit(item: Equipment): void {
    const match = this.equipmentTypes().find((type) => type.id === item.type || type.name === item.type);
    this.notice.set('');
    this.equipmentForm.set({ id: item.id, code: item.code, typeId: match?.id ?? '', currentHourMeter: String(item.currentHourMeter) });
  }

  protected cancelEquipmentForm(): void { this.equipmentForm.set(null); }

  protected saveEquipment(values: EquipmentFormSubmission): void {
    const form = this.equipmentForm();
    if (!form || this.saving() || this.loading()) return;
    const request: EquipmentRequest = {
      code: values.code,
      typeId: values.typeId,
      current_hour_meter: values.currentHourMeter,
    };
    this.saving.set(true);
    this.error.set('');
    const save = form.id
      ? this.equipmentApi.updateById(form.id, request)
      : this.equipmentApi.create(request);
    save.pipe(takeUntilDestroyed(this.destroyRef)).subscribe({
      next: (saved) => {
        this.equipment.update((items) => form.id ? items.map((item) => item.id === saved.id ? saved : item) : [saved, ...items]);
        this.equipmentForm.set(null);
        this.saving.set(false);
        this.notice.set(form.id ? `Se actualizaron los datos de ${saved.code}.` : `Se añadió ${saved.code} a la flota.`);
      },
      error: () => { this.error.set('No se pudo guardar el equipo. Tus cambios siguen en el formulario.'); this.saving.set(false); },
    });
  }

  protected askDelete(item: Equipment): void { this.deleteCandidate.set(item); }
  protected cancelDelete(): void { this.deleteCandidate.set(null); }

  protected deleteEquipment(): void {
    const item = this.deleteCandidate();
    if (!item || this.saving() || this.loading()) return;
    this.saving.set(true);
    this.error.set('');
    this.equipmentApi.deleteById(item.id).pipe(takeUntilDestroyed(this.destroyRef)).subscribe({
      next: () => {
        this.equipment.update((items) => items.filter((candidate) => candidate.id !== item.id));
        this.deleteCandidate.set(null);
        this.saving.set(false);
        this.notice.set(`Se retiró ${item.code} de la flota.`);
        if (this.historyFor() === item.id) this.closeHistory();
      },
      error: () => { this.error.set(`No se pudo retirar ${item.code}. Puede tener registros de mantenimiento asociados.`); this.saving.set(false); },
    });
  }

  protected startMaintenance(item: Equipment): void {
    this.notice.set('');
    this.maintenanceForm.set({ equipmentId: item.id, hmMaintenance: String(item.currentHourMeter), maintenanceDate: new Date().toISOString().slice(0, 10), description: '' });
  }

  protected cancelMaintenance(): void { this.maintenanceForm.set(null); }

  protected saveMaintenance(values: MaintenanceFormSubmission): void {
    const current = this.maintenanceForm();
    if (!current || this.saving() || this.loading()) return;
    const request = {
      equipmentId: current.equipmentId,
      hmMaintenance: values.hmMaintenance,
      maintenanceDate: this.formatToISOWithTimezone(values.maintenanceDate),
      description: values.description,
    };
    this.saving.set(true);
    this.error.set('');
    this.maintenanceApi.create(request).pipe(takeUntilDestroyed(this.destroyRef)).subscribe({
      next: (entry) => {
        this.equipment.update((items) => items.map((item) => item.id === current.equipmentId
          ? {
              ...item,
              currentHourMeter: Math.max(item.currentHourMeter, entry.hmMaintenance),
              lastMaintenanceHourMeter: entry.hmMaintenance,
              hoursSinceMaintenance: Math.max(0, item.currentHourMeter - entry.hmMaintenance),
              maintenanceStatus: 'up-to-date',
            }
          : item));
        this.maintenanceForm.set(null);
        this.saving.set(false);
        this.notice.set('Se registró el mantenimiento y se actualizó el ciclo de servicio.');
        if (this.historyFor() === current.equipmentId) this.loadHistory(current.equipmentId);
      },
      error: () => { this.error.set('No se pudo registrar el mantenimiento. Comprueba el horómetro e inténtalo de nuevo.'); this.saving.set(false); },
    });
  }

  protected showHistory(item: Equipment): void {
    this.historyFor.set(item.id);
    this.loadHistory(item.id);
  }

  protected closeHistory(): void {
    this.historyRequestId += 1;
    this.historyFor.set(null);
    this.history.set([]);
    this.historyLoading.set(false);
    this.historyError.set('');
  }

  private loadHistory(id: string): void {
    const requestId = ++this.historyRequestId;
    this.historyLoading.set(true);
    this.history.set([]);
    this.historyError.set('');
    this.equipmentApi.getAllMaintenancesByEquipmentId(id).pipe(takeUntilDestroyed(this.destroyRef)).subscribe({
      next: (entries) => {
        if (requestId !== this.historyRequestId || this.historyFor() !== id) return;
        this.history.set(entries);
        this.historyLoading.set(false);
      },
      error: () => {
        if (requestId !== this.historyRequestId || this.historyFor() !== id) return;
        this.historyError.set('No se pudo cargar el historial de mantenimiento.');
        this.historyLoading.set(false);
      },
    });
  }

  protected typeInterval(item: Equipment): number {
    return this.equipmentTypes().find((type) => type.name === item.type || type.id === item.type)?.maintenanceHours
      ?? item.maintenanceIntervalHours;
  }

  protected equipmentCode(id: string): string { return this.equipment().find((item) => item.id === id)?.code ?? ''; }

  protected statusFor(item: Equipment): EquipmentStatus {
    const interval = this.typeInterval(item);
    const remaining = interval - item.hoursSinceMaintenance;
    if (remaining <= 0) return { key: 'overdue', label: 'Requiere mantenimiento', detail: `${Math.abs(remaining).toLocaleString()} h sobre el intervalo` };
    if (remaining <= Math.max(10, interval * 0.1)) return { key: 'due-soon', label: 'Mantenimiento próximo', detail: `${remaining.toLocaleString()} h restantes` };
    return { key: 'on-track', label: 'En intervalo', detail: `${remaining.toLocaleString()} h restantes` };
  }

  protected formatToISOWithTimezone(date: string) {
    let data = new Date(date)
    const pad = (num: number) => String(num).padStart(2, '0');

    const year = data.getFullYear();
    const month = pad(data.getMonth() + 1);
    const day = pad(data.getDate());
    const hours = pad(data.getHours());
    const minutes = pad(data.getMinutes());
    const seconds = pad(data.getSeconds());

    const tzOffset = -data.getTimezoneOffset();
    const sign = tzOffset >= 0 ? '+' : '-';
    const absOffset = Math.abs(tzOffset);
    const offsetHours = pad(Math.floor(absOffset / 60));
    const offsetMinutes = pad(absOffset % 60);

    return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}${sign}${offsetHours}:${offsetMinutes}`;
  }
  
}
