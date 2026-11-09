import { Component, effect, input, output, signal } from '@angular/core';
import { FormField, FormRoot, form, maxLength, min, required } from '@angular/forms/signals';
import type { EquipmentType } from '../../../models/equipment-type';

export interface EquipmentFormState {
  id?: string;
  code: string;
  typeId: string;
  currentHourMeter: string;
}

export interface EquipmentFormSubmission {
  code: string;
  typeId: string;
  currentHourMeter: number;
}

@Component({
  selector: 'app-equipment-form',
  templateUrl: './equipment-form.html',
  styleUrls: ['../dashboard-shared.css', './equipment-form.css'],
  imports: [FormField, FormRoot],
})
export class EquipmentFormComponent {
  readonly formState = input.required<EquipmentFormState>();
  readonly equipmentTypes = input.required<EquipmentType[]>();
  readonly saving = input(false);
  readonly loading = input(false);
  readonly submitted = output<EquipmentFormSubmission>();
  readonly cancelled = output<void>();

  private readonly model = signal({
    code: '',
    typeId: '',
    currentHourMeter: 0,
  });

  protected readonly equipmentForm = form(this.model, (schemaPath) => {
    required(schemaPath.code, { message: 'Ingresa el código del equipo.' });
    maxLength(schemaPath.code, 40, { message: 'El código no puede superar los 40 caracteres.' });
    required(schemaPath.typeId, { message: 'Selecciona un tipo de equipo.' });
    required(schemaPath.currentHourMeter, { message: 'Ingresa la lectura actual del horómetro.' });
    min(schemaPath.currentHourMeter, 0, { message: 'El horómetro no puede ser negativo.' });
  }, {
    submission: {
      action: async () => {
        const values = this.model();
        this.submitted.emit({
          code: values.code.trim(),
          typeId: values.typeId,
          currentHourMeter: values.currentHourMeter,
        });
      },
    },
  });

  constructor() {
    effect(() => {
      const state = this.formState();
      this.model.set({
        code: state.code,
        typeId: state.typeId,
        currentHourMeter: Number(state.currentHourMeter),
      });
    });
  }
}
