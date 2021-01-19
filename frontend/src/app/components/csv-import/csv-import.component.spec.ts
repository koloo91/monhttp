import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CsvImportComponent } from './csv-import.component';

describe('CsvImportComponent', () => {
  let component: CsvImportComponent;
  let fixture: ComponentFixture<CsvImportComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ CsvImportComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(CsvImportComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
