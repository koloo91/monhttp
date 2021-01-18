import {Service} from './service.model';

export interface UploadResult {
  rowNumber: number;
  service: Service;
  error: string;
}
