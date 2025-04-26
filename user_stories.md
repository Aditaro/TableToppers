### Register a New User 

User Story: 
 
As a new customer, I want to create an account so that I can securely access the application’s features. 

Acceptance Criteria: 

I must provide all required fields (username, password, email, phone number, first name, last name) 

### Login for Existing User 

User Story: 
 
As an existing user (customer/staff/manager/admin), I want to log in so that I can securely access my account and the relevant system features. 

Acceptance Criteria: 

I must provide valid credentials (username/password/timestamp). Set up authentication on sign in, proper error handling and validation should be implemented. 
 

### Retrieve a List of Restaurants  

User Story (Customer): 
 
As a customer, I want to browse through a list of restaurants so that I can pick where I’d like to dine or reserve a table. 

User Story (Restaurant Manager/Staff): 
 
As a restaurant manager or staff, I want to retrieve only the restaurants under my management (or all if I’m an admin) so that I can oversee them. 

Acceptance Criteria: 

I can optionally filter restaurants by city or name 

### Join a Waitlist 

User Story: 
 
As a customer, I want to be able to enter a waitlist for reservations if the restaurant is fully booked, so that I have a chance to get a table in case of cancellations. 

Acceptance Criteria: 
If a restaurant is fully booked for a given time slot, I can join a waitlist. 
If I am next in line and a table opens, I have a limited time to confirm the reservation before it is offered to the next person. 
 
### Create a New Restaurant 

User Story: 
 
As a manager (or admin), I want to add a new restaurant to the system so that I can start managing it. 

Acceptance Criteria: 

I must provide valid restaurant data (name, location, etc.). 

### Get All Tables for a Restaurant 

User Story: 
 
As a restaurant staff or a manager, I want to see all the tables for a specific restaurant so that I can manage seating and table availability. 

Acceptance Criteria: 

I must provide the restaurant Id. 


### Create a New Table 

User Story: 
 
As a manager or staff, I want to add a new table to a restaurant so that I can increase seating capacity or reorganize the layout. 

Acceptance Criteria: 

I must provide the required fields (min capacity, max capacity) in the request body. 


### Create a New Reservation 

User Story (Customer): 
 
As a customer, I want to make a reservation for a specific restaurant so that I can secure a table in advance. 

Acceptance Criteria: 

I must provide reservation time, number of guests, and phone number. 


### Get All Reservations for a Restaurant 

User Story (Staff/Manager): 
 
As restaurant staff or a manager, I want to see all reservations for a restaurant to manage seat assignments and scheduling. 

Acceptance Criteria: 

I must provide the restaurant Id. 

I can optionally filter by date (query parameter). 
 

### Get a Single Reservation 

User Story (Staff/Manager/Customer who made it):  
 
As the reservation owner or restaurant staff/manager, I want to view the details of a specific reservation so that I can confirm the reservation or see the booking information. 

### Landing Page Development 

User Story:  

As a visitor, I want to see a welcoming landing page so that I can learn about the different restaurants 

Acceptance Criteria: 
Basic navigation for other sections of the site and the page should be optimized for mobile users. 

### Backend Environment Setup 
User Story:   
 
As a developer, setup the environment for frontend 
 
Acceptence Criteria: 
Setup of Node, VS Code, Angular CLI, Git setup/Sourcetree 

 
