export interface Notifier {
  id: string;
  name: string;
  data: any;
  form: NotifierForm[];
}

export interface NotifierForm {
  type: string;
  title: string;
  formControlName: string;
  placeholder: string;
  required: boolean;
}
