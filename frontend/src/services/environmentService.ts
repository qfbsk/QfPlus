import {
  GetEnvironmentInventory,
  GetEnvironmentStatus,
  OpenEnvironmentDiagnostic,
} from '../../wailsjs/go/app/App';
import { model } from '../../wailsjs/go/models';

export type EnvironmentStatusReport = model.EnvironmentStatusReport;
export type EnvironmentInventory = model.EnvironmentInventory;

export const fetchEnvironmentStatus = () => GetEnvironmentStatus();

export const openEnvironmentDiagnostic = () => OpenEnvironmentDiagnostic();

export const fetchEnvironmentInventory = () => GetEnvironmentInventory();
