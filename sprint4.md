# Sprint 4 Report 

#### Waitlist Feature Implementation
Implemented GET, POST, and DELETE routes for the waitlist system:

- GET: Get all waitlist entries for a restaurant.

- POST: A new customer is added to a restaurant waitlist.

- DELETE: Remove a customer from the waitlist using the generated waitlist
entry ID.

Correct retrieval, creation, and deletion of waitlist entries verified
via Postman.

==============================================================================

Readme Setup I created a full README.md for the backend that describes:

- Prerequisites and environmental setup.

- Run the backend server configuration.

- Availability of API routes and examples.

- Backend Routes and Unit Testing. Restauration/table route testing
continues.

- Validated API functionality using Postman and unit testing.

- Main functionality was added and tested - data updates (e.g. table
status) are processed correctly.
========================================================================================

### Frontend Contributions Waitlist Component A UI component

Key features include:
- allow staff to manage waitlisted customers.  

- Enter customer name, party size and phone number in input form.

- Automatic cohort assignment is based on party size.

- Figure showing estimated wait time calculation along with cohort
information.

- You can assign a table to a waiting customer and cancel their waitlist
status.

- The UI updates when waitlist status changes.

- Customer-Facing Updates Filler text and images were replaced by New
Arrivals for newly opened restaurants or dishes.

- Updated footer quick links to redirect to appropriate customer pages
(e.g. Info tab).

- Business Portal Enhancements A GIF demonstration of table management
features is embedded.

========================================================================

### Implemented modals for:

- Schedule Demo

- Contact Sales

- Start Free Trial

- Add a tiered pricing section with increasing features for each
subscription level.

- Added FAQ section for common customer questions.

======================================================================

### Cypress Testing
Tests for:  
  
- Visibility of the business portal and demo GIF.  
  
- Triggering demo/contact modals properly.  
  
- Submission of demos is required.  
  
- Navigate to the subscription page.  
  
- Validation of all footer quick links.  
  
navbar.cy.ts  
  
Navigation Bar  
1. Should display all main navigation items 
o	Visit /  
o	Assert the following items are visible in #nav ul:  
-	Business Portal  
-	Restaurants  
-	Login  
-	Register  
2. Should navigate to correct pages when clicking links   
o	Click Restaurants → expect URL to include /restaurants  
o	Click Login → expect URL to include /login  
o	Click Register → expect URL to include /register  
o	Click Business Portal → expect URL to include /business-portal  
________________________________________
business-portal.cy.ts  
  
Business Portal Page  
1. Should display the business portal page correctly   
o	Visit /business-portal  
o	Assert .business-portal exists  
o	Assert h1 contains "Transform Your Restaurant Management"  
2. Should display and play the demo GIF ✅  
o	Visit /business-portal  
o	Assert img.demo-image:  
-	Is visible  
-	Has src="assets/css/images/table-managment-demo.gif"  
-	Has alt="Table Management Demo"  
3. Should open demo form modal when clicking "Schedule Demo"   
o	Click Schedule Demo  
o	Assert .modal-overlay and .modal-content are visible  
4. Should open demo form modal when clicking "Contact Sales"   
o	Click Contact Sales  
o	Assert .modal-overlay and .modal-content are visible  
5. Should be able to fill out and submit the demo form   
o	Click Schedule Demo  
o	Fill out:  
-	#name: "Test User"  
-	#email: "test@example.com"  
-	#phone: "1234567890"  
-	#company: "Test Company"  
-	#message: "This is a test message"  
o	Assert .submit-button is not disabled  
o	Click submit  
o	Assert .success-message is visible  
6. Should navigate to subscription page when clicking "Get Started Free"  
o	Click Get Started Free  
o	Expect URL to include /subscription  
7. Should navigate to subscription page when clicking "Start Free Trial"   
o	Click Start Free Trial  
o	Expect URL to include /subscription  
________________________________________
Home Page Navigation  
1. Should navigate to login page when clicking "More Info"   
o	Click More Info (.button.alt)  
o	Expect URL to include:  
-	/login  
-	returnUrl=%2Fcustomer-info  
2. Should navigate to login page when clicking "Reserve Now"   
o	Click Reserve Now (.button)  
o	Expect URL to include /login

======================================================================

### Frontend Unit Tests:

Reservation Component: Test form visibility, field input, and submission
handling.

Restaurant Component: Input for new restaurant creation is validated.

Waitlist Component: Validates waitlist data fetching and cohort sorting,
form dialog display, table availability display, and estimated wait time
logic.

Table Component: Testing initialization behavior, floor plan visibility,
and interaction with reservation list toggle.


======================================================================

### Backend Unit Tests:  
  
Backend unit tests for waitlist GET, POST, and DELETE requests.  
TestGetWaitlist  
Mock data for waitlist entries  
{  
		{  
			ID:                "34f156e9-22d4-4507-85a2-aadd843ac251",  
			RestaurantID:      "059ffaf3-1409-4da1-b1c5-187dda0e27a5",  
			Name:              "John Doe",  
			PhoneNumber:       "1234567890",  
			PartySize:         4,  
			PartyAhead:        2,  
			EstimatedWaitTime: 15,  
			CreatedAt:         "2025-04-21T12:34:56Z",  
		},  
		{  
			ID:                "ae5c0877-995f-43e1-8724-f81d16c38ef2",  
			RestaurantID:      "70336dc8-0764-432c-882c-033c2b0eac65",  
			Name:              "Jane Smith",  
			PhoneNumber:       "0987654321",  
			PartySize:         2,  
			PartyAhead:        1,  
			EstimatedWaitTime: 10,  
			CreatedAt:         "2025-04-21T13:00:00Z",  
		},  
	}  
Expected Server Response: 200  
Expected Output:   
{  
		{  
			ID:                "34f156e9-22d4-4507-85a2-aadd843ac251",  
			RestaurantID:      "059ffaf3-1409-4da1-b1c5-187dda0e27a5",  
			Name:              "John Doe",  
			PhoneNumber:       "1234567890",  
			PartySize:         4,  
			PartyAhead:        2,  
			EstimatedWaitTime: 15,  
			CreatedAt:         "2025-04-21T12:34:56Z",  
		},  
		{  
			ID:                "ae5c0877-995f-43e1-8724-f81d16c38ef2",  
			RestaurantID:      "70336dc8-0764-432c-882c-033c2b0eac65",  
			Name:              "Jane Smith",  
			PhoneNumber:       "0987654321",  
			PartySize:         2,  
			PartyAhead:        1,  
			EstimatedWaitTime: 10,  
			CreatedAt:         "2025-04-21T13:00:00Z",  
		},  
	}  

TestCreateWaitlistEntry:  
Mock data for creating a waitlist entry  
  {  
		RestaurantID:      "059ffaf3-1409-4da1-b1c5-187dda0e27a5",  
		Name:              "John Doe",  
		PhoneNumber:       "1234567890",  
		PartySize:         4,  
		PartyAhead:        2,  
		EstimatedWaitTime: 15,  
	}  
Expected Server Response: 201  
Expected Output: "Waitlist entry created successfully"  
Check that resulting data is nonempty.  
  
TestDeleteWaitlistEntry  
Mock Delete info:  
waitlistID := "34f156e9-22d4-4507-85a2-aadd843ac251"  
restaurantID := "059ffaf3-1409-4da1-b1c5-187dda0e27a5"  
  
Expected Server Response: 200  
Expected Output: "Waitlist entry deleted successfully"  
Check data becomes empty after delete  
  
Backend unit tests for tables and restaurant updating.  
  
Backend unit tests for reservation GET, POST, PATCH, and DELETE requests.  
  
#### Full backend API can be found in openapi.yaml, all endpoints implemented except the menu section, which is in our plans for future work.  
  



