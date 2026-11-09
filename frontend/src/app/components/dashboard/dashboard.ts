import { Component, computed, inject, signal } from '@angular/core';
import { Equipment, EquipmentRequest } from '../../models/equipment';
import { EquipmentType } from '../../models/equipment-type';
import { Maintenance } from '../../models/maintenance';
import { EquipmentAPI } from '../../services/equipment-api';
import { EquipmentTypeAPI } from '../../services/equipment-type-api';
import { MaintenanceAPI } from '../../services/maintenance-api';
import { DatePipe } from '@angular/common';

type EquipmentForm = { id?: string; code: string; typeId: string; currentHourMeter: string };
type MaintenanceForm = { equipmentId: string; hmMaintenance: string; maintenanceDate: string; description: string };

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.html',
  styleUrl: './dashboard.css',
  imports: [DatePipe]
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
  protected readonly equipmentForm = signal<EquipmentForm | null>(null);
  protected readonly maintenanceForm = signal<MaintenanceForm | null>(null);
  protected readonly historyFor = signal<string | null>(null);
  protected readonly history = signal<Maintenance[]>([]);
  protected readonly historyLoading = signal(false);
  protected readonly deleteCandidate = signal<Equipment | null>(null);
  protected readonly dueCount = computed(() => this.equipment().filter((item) => this.statusFor(item).key !== 'on-track').length);

  constructor() {
    this.load();
  }

  protected load(): void {
    this.loading.set(true);
    this.error.set('');
    let equipmentLoaded = false;
    let typesLoaded = false;
    const finish = (): void => {
      if (equipmentLoaded && typesLoaded) this.loading.set(false);
    };

    this.equipmentApi.getAll().subscribe({
      next: (items) => { this.equipment.set(items); equipmentLoaded = true; finish(); },
      error: () => { this.error.set('No se pudo cargar la maquinaria. Comprueba la conexión e inténtalo de nuevo.'); equipmentLoaded = true; finish(); },
    });
    this.equipmentTypeApi.getAll().subscribe({
      next: (items) => { this.equipmentTypes.set(items); typesLoaded = true; finish(); },
      error: () => { this.error.set('No se pudieron cargar los tipos de maquinaria. Inténtalo de nuevo antes de añadir o editar equipos.'); typesLoaded = true; finish(); },
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

  protected saveEquipment(event: SubmitEvent): void {
    event.preventDefault();
    const form = this.equipmentForm();
    const element = event.currentTarget;
    if (!form || !(element instanceof HTMLFormElement) || this.saving()) return;
    const values = new FormData(element);
    const request: EquipmentRequest = {
      code: String(values.get('code') ?? '').trim(),
      typeId: String(values.get('typeId') ?? ''),
      current_hour_meter: Number(values.get('currentHourMeter')),
    };
    this.saving.set(true);
    this.error.set('');
    const save = form.id
      ? this.equipmentApi.updateById(form.id, request)
      : this.equipmentApi.create(request);
    save.subscribe({
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
    if (!item || this.saving()) return;
    this.saving.set(true);
    this.error.set('');
    this.equipmentApi.deleteById(item.id).subscribe({
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

  protected saveMaintenance(event: SubmitEvent): void {
    event.preventDefault();
    const current = this.maintenanceForm();
    const element = event.currentTarget;
    if (!current || !(element instanceof HTMLFormElement) || this.saving()) return;
    const values = new FormData(element);
    const maintenanceHour = Number(values.get('hmMaintenance'));
    const request = {
      equipmentId: current.equipmentId,
      hmMaintenance: maintenanceHour,
      maintenanceDate: this.formatToISOWithTimezone(new Date(String(values.get('maintenanceDate') ?? ''))),
      description: String(values.get('description') ?? '').trim() || undefined,
    };
    this.saving.set(true);
    this.error.set('');
    this.maintenanceApi.create(request).subscribe({
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

  protected closeHistory(): void { this.historyFor.set(null); this.history.set([]); }

  private loadHistory(id: string): void {
    this.historyLoading.set(true);
    this.history.set([]);
    this.equipmentApi.getAllMaintenancesByEquipmentId(id).subscribe({
      next: (entries) => { this.history.set(entries); this.historyLoading.set(false); },
      error: () => { this.error.set('No se pudo cargar el historial de mantenimiento.'); this.historyLoading.set(false); },
    });
  }

  protected labelForType(typeId: string): string { return this.equipmentTypes().find((type) => type.id === typeId)?.name ?? ''; }

  protected typeInterval(item: Equipment): number {
    return this.equipmentTypes().find((type) => type.name === item.type || type.id === item.type)?.maintenanceHours
      ?? item.maintenanceIntervalHours;
  }

  protected meterPercent(item: Equipment): number {
    const interval = this.typeInterval(item);
    return interval > 0 ? Math.min(100, Math.max(0, item.hoursSinceMaintenance / interval * 100)) : 0;
  }

  protected equipmentCode(id: string): string { return this.equipment().find((item) => item.id === id)?.code ?? ''; }

  protected statusFor(item: Equipment): { key: string; label: string; detail: string } {
    const interval = this.typeInterval(item);
    const remaining = interval - item.hoursSinceMaintenance;
    if (remaining <= 0) return { key: 'overdue', label: 'Requiere mantenimiento', detail: `${Math.abs(remaining).toLocaleString()} h sobre el intervalo` };
    if (remaining <= Math.max(10, interval * 0.1)) return { key: 'due-soon', label: 'Mantenimiento próximo', detail: `${remaining.toLocaleString()} h restantes` };
    return { key: 'on-track', label: 'En intervalo', detail: `${remaining.toLocaleString()} h restantes` };
  }

  protected trackEquipment(_index: number, item: Equipment): string { return item.id; }
  protected trackMaintenance(_index: number, item: Maintenance): string { return item.id; }

  protected formatToISOWithTimezone(date: Date) {
    const pad = (num: number) => String(num).padStart(2, '0');

    const year = date.getFullYear();
    const month = pad(date.getMonth() + 1);
    const day = pad(date.getDate());
    const hours = pad(date.getHours());
    const minutes = pad(date.getMinutes());
    const seconds = pad(date.getSeconds());

    const tzOffset = -date.getTimezoneOffset();
    const sign = tzOffset >= 0 ? '+' : '-';
    const absOffset = Math.abs(tzOffset);
    const offsetHours = pad(Math.floor(absOffset / 60));
    const offsetMinutes = pad(absOffset % 60);

    return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}${sign}${offsetHours}:${offsetMinutes}`;
  }
  
}
