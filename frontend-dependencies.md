
# 🧩 Project Dependencies Overview

This Angular frontend project is built using a modern enterprise-grade setup, including support for component libraries, internationalization, testing, code quality, and documentation tools. Below is a categorized overview of key dependencies and their purposes.

---

## 🟢 Core Framework

| Package | Purpose |
|--------|---------|
| `@angular/core`, `@angular/common`, `@angular/router`, etc. | Core Angular 19+ packages for app structure, forms, routing, and browser interaction |
| `rxjs` | Reactive programming for Angular services and components |
| `zone.js` | Required for Angular's change detection mechanism |
| `tslib` | Runtime library for TypeScript helpers |

---

## 🎨 UI & UX

| Package | Purpose |
|--------|---------|
| `@angular/material` | Material Design components |
| `@angular/cdk` | Component Dev Kit for custom behavior (e.g., overlays, drag-drop) |
| `konva` | 2D canvas rendering for custom visual elements |
| `ngx-mask` | Input masking for form fields |

---

## 🌐 Internationalization

| Package | Purpose |
|--------|---------|
| `@jsverse/transloco` | Lightweight i18n framework for Angular apps |

---

## ⚙️ Development Tooling

| Package | Purpose |
|--------|---------|
| `@angular/cli` | Angular CLI for scaffolding and build |
| `@angular-devkit/build-angular` | Webpack-based build tools for Angular apps |
| `typescript`, `ts-node` | TypeScript support and node-based execution |
| `source-map-explorer` | Visualize bundle sizes and source maps |

---

## 🧪 Testing

| Tool | Purpose |
|------|--------|
| `karma`, `jasmine-core` | Unit testing framework for Angular |
| `jest`, `jest-preset-angular` | Alternative test runner for modern JS testing |
| `cypress` | End-to-end (E2E) testing framework |
| `@storybook/test` | Component-level interaction testing |
| `jest-junit` | Outputs test results in JUnit format for CI integration |

---

## 🧹 Linting & Formatting

| Tool | Purpose |
|------|--------|
| `eslint`, `@angular-eslint/*` | ESLint-based linting configuration for Angular |
| `prettier` | Code formatting |
| `husky`, `lint-staged` | Git hook setup for automatic formatting before commit |

---

## 📦 Installation

To install all dependencies:

```bash
npm install
```

For development and production builds, use:

```bash
npm run start       # Development server
npm run build       # Production build
npm run test        # Run unit tests
```
