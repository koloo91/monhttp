import {Component, Inject, OnInit} from '@angular/core';
import {MAT_DIALOG_DATA, MatDialogRef} from '@angular/material/dialog';
import {Service} from '../../../models/service.model';

@Component({
  selector: 'app-confirm-service-delete-dialog',
  templateUrl: './confirm-service-delete-dialog.component.html',
  styleUrls: ['./confirm-service-delete-dialog.component.scss']
})
export class ConfirmServiceDeleteDialogComponent implements OnInit {

  constructor(public dialogRef: MatDialogRef<ConfirmServiceDeleteDialogComponent>,
              @Inject(MAT_DIALOG_DATA) public service: Service) {
  }

  ngOnInit(): void {
  }

  onCancelClick(): void {
    this.dialogRef.close(false);
  }

  onOkClick(): void {
    this.dialogRef.close(true);
  }
}
