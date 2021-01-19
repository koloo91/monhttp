import {Component, OnInit} from '@angular/core';
import {ImportService} from '../../services/import.service';
import {finalize} from 'rxjs/operators';
import {MatDialog} from '@angular/material/dialog';
import {ImportResult} from '../../models/upload_result.model';
import {ImportResultDialogComponent} from '../dialogs/import-result-dialog/import-result-dialog.component';

@Component({
  selector: 'app-csv-import',
  templateUrl: './csv-import.component.html',
  styleUrls: ['./csv-import.component.scss']
})
export class CsvImportComponent implements OnInit {

  isLoading = false;
  fileToUpload: File = null;

  constructor(private importService: ImportService,
              public dialog: MatDialog) {
  }

  ngOnInit(): void {
  }

  handleFileInput(event: Event): void {
    // @ts-ignore
    this.fileToUpload = event.target.files[0];
  }

  uploadFile() {
    this.isLoading = true;
    const formData = new FormData();
    formData.append('file', this.fileToUpload);
    this.importService.importCsv(formData)
      .pipe(
        finalize(() => this.isLoading = false)
      )
      .subscribe((wrapper) => {
        this.showImportResult(wrapper.data);
      }, console.log);
  }

  showImportResult(importResults: ImportResult[]): void {
    this.dialog.open(ImportResultDialogComponent, {
      data: importResults
    });
  }
}
