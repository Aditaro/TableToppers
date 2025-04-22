import { Component, ViewChild } from '@angular/core';
import { RouterModule } from '@angular/router';
import { CommonModule } from '@angular/common';
import { DemoFormComponent } from '../demo-form/demo-form.component';
import { ReactiveFormsModule } from '@angular/forms';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule, RouterModule, DemoFormComponent, ReactiveFormsModule],
  templateUrl: './home.component.html',
  styleUrl: './home.component.css',
})
export class HomeComponent {
  @ViewChild(DemoFormComponent) demoForm!: DemoFormComponent;

  openDemoForm() {
    this.demoForm.openModal();
  }
}
