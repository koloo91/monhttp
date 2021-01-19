import {Service} from './service.model';

export interface ImportResult {
  rowNumber: number;
  service: Service;
  error: string;
}
