export interface Check {
  id: string;
  serviceId: string;
  latencyInMs: number;
  isFailure: boolean;
  createdAt: string;
}
