import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-demo-form',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './demo-form.component.html',
  styleUrls: ['./demo-form.component.css']
})
export class DemoFormComponent implements OnInit {
  isOpen = false;
  demoForm: FormGroup;
  showSuccessMessage = false;

  constructor(private fb: FormBuilder) {
    this.demoForm = this.fb.group({
      name: ['', Validators.required],
      email: ['', [Validators.required, Validators.email]],
      phone: ['', Validators.required],
      company: ['', Validators.required],
      message: ['']
    });
  }

  ngOnInit(): void {}

  openModal() {
    this.isOpen = true;
    this.showSuccessMessage = false;
  }

  closeModal() {
    this.isOpen = false;
    this.demoForm.reset();
    this.showSuccessMessage = false;
  }

  onSubmit() {
    if (this.demoForm.valid) {
      // Here you would typically send the form data to your backend
      console.log(this.demoForm.value);
      
      // Show success message
      this.showSuccessMessage = true;
      
      // Reset form and close modal after 3 seconds
      setTimeout(() => {
        this.closeModal();
      }, 3000);
    }
  }
} 