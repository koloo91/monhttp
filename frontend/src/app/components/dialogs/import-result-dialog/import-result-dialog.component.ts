import {Component, Inject, OnInit} from '@angular/core';
import {ImportResult} from '../../../models/upload_result.model';
import {MAT_DIALOG_DATA, MatDialogRef} from '@angular/material/dialog';

@Component({
  selector: 'app-import-result-dialog',
  templateUrl: './import-result-dialog.component.html',
  styleUrls: ['./import-result-dialog.component.scss']
})
export class ImportResultDialogComponent implements OnInit {

  constructor(public dialogRef: MatDialogRef<ImportResultDialogComponent>,
              @Inject(MAT_DIALOG_DATA) public data: ImportResult[]) {
  }

  ngOnInit(): void {
  }

}
