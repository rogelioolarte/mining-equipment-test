import { Component, effect, input, output, signal } from '@angular/core';
import { form, FormField, FormRoot, maxLength, min, required } from '@angular/forms/signals';

export interface MaintenanceFormState {
  equipmentId: string;
  hmMaintenance: string;
  maintenanceDate: string;
  description: string;
}

export interface MaintenanceFormSubmission {
  hmMaintenance: number;
  maintenanceDate: string;
  description?: string;
}

@Component({
  selector: 'app-maintenance-form',
  templateUrl: './maintenance-form.html',
  styleUrls: ['../dashboard-shared.css', './maintenance-form.css'],
  imports: [FormField, FormRoot]
})
export class MaintenanceFormComponent {
  readonly formState = input.required<MaintenanceFormState>();
  readonly equipmentCode = input.required<string>();
  readonly saving = input(false);
  readonly loading = input(false);
  readonly submitted = output<MaintenanceFormSubmission>();
  readonly cancelled = output<void>();

  private readonly model = signal({
    hmMaintenance: 0,
    maintenanceDate: '',
    description: '',
  });

  protected readonly maintenanceForm = form(this.model, (schemaPath) => {
    required(schemaPath.hmMaintenance, { message: 'Ingresa la lectura del horómetro.' });
    min(schemaPath.hmMaintenance, 0, { message: 'El horómetro no puede ser negativo.' });
    required(schemaPath.maintenanceDate, { message: 'Selecciona la fecha del servicio.' });
    maxLength(schemaPath.description, 500, { message: 'La descripción no puede superar los 500 caracteres.' });
  }, {
    submission: {
      action: async () => {
        const values = this.model();
        this.submitted.emit({
          hmMaintenance: values.hmMaintenance,
          maintenanceDate: values.maintenanceDate,
          description: values.description.trim() || undefined,
        });
      },
    },
  });

  constructor() {
    effect(() => {
      const state = this.formState();
      this.model.set({
        hmMaintenance: Number(state.hmMaintenance),
        maintenanceDate: state.maintenanceDate,
        description: state.description,
      });
    });
  }
}
