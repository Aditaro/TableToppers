Sprint 4 Report Backend Contributions Waitlist Feature Implementation
Implemented GET, POST, and DELETE routes for the waitlist system:

GET: Get all waitlist entries for a restaurant.

POST: A new customer is added to a restaurant waitlist.

DELETE: Remove a customer from the waitlist using the generated waitlist
entry ID.

Correct retrieval, creation, and deletion of waitlist entries verified
via Postman.

Readme Setup I created a full README.md for the backend that describes:

Prerequisites and environmental setup.

Run the backend server configuration.

Availability of API routes and examples.

Backend Routes and Unit Testing. Restauration/table route testing
continues.

Validated API functionality using Postman and unit testing.

Main functionality was added and tested - data updates (e.g. table
status) are processed correctly.

Frontend Contributions Waitlist Component A UI component was developed
to allow staff to manage waitlisted customers.

Key features include:

Enter customer name, party size and phone number in input form.

Automatic cohort assignment is based on party size.

Figure showing estimated wait time calculation along with cohort
information.

You can assign a table to a waiting customer and cancel their waitlist
status.

The UI updates when waitlist status changes.

Customer-Facing Updates Filler text and images were replaced by New
Arrivals for newly opened restaurants or dishes.

Updated footer quick links to redirect to appropriate customer pages
(e.g. Info tab).

Business Portal Enhancements A GIF demonstration of table management
features is embedded.

Implemented modals for:

Schedule Demo

Contact Sales

Start Free Trial

Add a tiered pricing section with increasing features for each
subscription level.

Added FAQ section for common customer questions.

Cypress Testing Addition and passing Cypress tests for:

Visibility of the business portal and demo GIF.

Triggering demo/contact modals properly.

Submission of demos is required.

Navigate to the subscription page.

Validation of all footer quick links.

Unit Testing Unit tests for key components were implemented and
extended:

Reservation Component: Test form visibility, field input, and submission
handling.

Restaurant Component: Input for new restaurant creation is validated.

Table Component: Testing initialization behavior, floor plan visibility,
and interaction with reservation list toggle.

Waitlist Component: Validates waitlist data fetching and cohort sorting,
form dialog display, table availability display, and estimated wait time
logic.
