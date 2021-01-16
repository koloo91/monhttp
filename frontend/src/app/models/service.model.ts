export type ServiceType = 'HTTP' | 'ICMP_PING';

export interface Service {
  id?: string;
  name: string;
  type: ServiceType
  intervalInSeconds: number;
  nextCheckTime: string;
  endpoint: string;
  httpMethod: string;
  requestTimeoutInSeconds: number;
  httpHeaders: string;
  httpBody: string;
  expectedHttpResponseBody: string;
  expectedHttpStatusCode: number;
  followRedirects: boolean;
  verifySsl: boolean
  enableNotifications: boolean;
  notifyAfterNumberOfFailures: number;
  continuouslySendNotifications: boolean;
  notifiers: string[];
  createdAt?: string;
  updatedAt?: string;
}
