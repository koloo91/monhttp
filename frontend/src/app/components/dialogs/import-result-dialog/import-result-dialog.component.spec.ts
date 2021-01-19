import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ImportResultDialogComponent } from './import-result-dialog.component';

describe('ImportResultDialogComponent', () => {
  let component: ImportResultDialogComponent;
  let fixture: ComponentFixture<ImportResultDialogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ImportResultDialogComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ImportResultDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
