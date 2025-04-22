import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';

@Component({
  selector: 'app-subscription',
  standalone: true,
  imports: [CommonModule, RouterModule],
  template: `
    <div class="subscription-page">
      <div class="container">
        <section class="header">
          <h1>Choose Your Plan</h1>
          <p>Select the perfect plan for your restaurant's needs</p>
        </section>

        <section class="pricing-plans">
          <div class="plans-grid">
            <div *ngFor="let plan of plans" class="plan-card" [class.popular]="plan.popular">
              <div class="plan-header">
                <h3>{{plan.name}}</h3>
                <div class="price">
                  <span class="amount">{{plan.price}}</span>
                  <span class="period">{{plan.period}}</span>
                </div>
              </div>
              <div class="plan-features">
                <ul>
                  <li *ngFor="let feature of plan.features">
                    <i class="fas fa-check"></i>
                    {{feature}}
                  </li>
                </ul>
              </div>
              <div class="plan-cta">
                <a [routerLink]="['/subscription']" class="cta-button">Get Started</a>
              </div>
            </div>
          </div>
        </section>

        <section class="faq">
          <h2>Frequently Asked Questions</h2>
          <div class="faq-grid">
            <div class="faq-item">
              <h3>Can I change my plan later?</h3>
              <p>Yes, you can upgrade or downgrade your plan at any time. Changes will be reflected in your next billing cycle.</p>
            </div>
            <div class="faq-item">
              <h3>Is there a free trial?</h3>
              <p>Yes, all plans come with a 14-day free trial. No credit card required.</p>
            </div>
            <div class="faq-item">
              <h3>What payment methods do you accept?</h3>
              <p>We accept all major credit cards, PayPal, and bank transfers.</p>
            </div>
            <div class="faq-item">
              <h3>Do you offer custom plans?</h3>
              <p>Yes, for large restaurants or chains, we offer custom plans. Contact our sales team for details.</p>
            </div>
          </div>
        </section>
      </div>
    </div>
  `,
  styleUrls: ['./subscription.component.css']
})
export class SubscriptionComponent {
  plans = [
    {
      name: 'Starter',
      price: '$49',
      period: 'per month',
      features: [
        'Up to 5 tables',
        'Basic table management',
        'Digital waitlist',
        'Email support',
        'Basic analytics'
      ],
      popular: false
    },
    {
      name: 'Professional',
      price: '$99',
      period: 'per month',
      features: [
        'Up to 20 tables',
        'Advanced table management',
        'Digital waitlist',
        'Priority support',
        'Advanced analytics',
        'Menu management',
        'Staff scheduling'
      ],
      popular: true
    },
    {
      name: 'Enterprise',
      price: '$199',
      period: 'per month',
      features: [
        'Unlimited tables',
        'Premium table management',
        'Digital waitlist',
        '24/7 support',
        'Custom analytics',
        'Menu management',
        'Staff scheduling',
        'API access',
        'Custom integrations'
      ],
      popular: false
    }
  ];
} 