# User Stories

## Register a New User  
**User Story:**  
As a new customer, I want to create an account so that I can securely access the application’s features.  

**Acceptance Criteria:**  
- I must provide all required fields (username, password, email, phone number, first name, last name).  

---

## Login for Existing User  
**User Story:**  
As an existing user (customer/staff/manager/admin), I want to log in so that I can securely access my account and the relevant system features.  

**Acceptance Criteria:**  
- I must provide valid credentials (username/password/timestamp).  
- Set up authentication on sign-in, proper error handling, and validation should be implemented.  

---

## Retrieve a List of Restaurants  
**User Story (Customer):**  
As a customer, I want to browse through a list of restaurants so that I can pick where I’d like to dine or reserve a table.  

**User Story (Restaurant Manager/Staff):**  
As a restaurant manager or staff, I want to retrieve only the restaurants under my management (or all if I’m an admin) so that I can oversee them.  

**Acceptance Criteria:**  
- I can optionally filter restaurants by city or name.  

---

## Join a Waitlist  
**User Story:**  
As a customer, I want to be able to enter a waitlist for reservations if the restaurant is fully booked, so that I have a chance to get a table in case of cancellations.  

**Acceptance Criteria:**  
- If a restaurant is fully booked for a given time slot, I can join a waitlist.  
- If I am next in line and a table opens, I have a limited time to confirm the reservation before it is offered to the next person.  

---

## Create a New Restaurant  
**User Story:**  
As a manager (or admin), I want to add a new restaurant to the system so that I can start managing it.  

**Acceptance Criteria:**  
- I must provide valid restaurant data (name, location, etc.).  

---

## Get All Tables for a Restaurant  
**User Story:**  
As restaurant staff or a manager, I want to see all the tables for a specific restaurant so that I can manage seating and table availability.  

**Acceptance Criteria:**  
- I must provide the restaurant ID.  

---

## Create a New Table  
**User Story:**  
As a manager or staff, I want to add a new table to a restaurant so that I can increase seating capacity or reorganize the layout.  

**Acceptance Criteria:**  
- I must provide the required fields (min capacity, max capacity) in the request body.  

---

## Create a New Reservation  
**User Story (Customer):**  
As a customer, I want to make a reservation for a specific restaurant so that I can secure a table in advance.  

**Acceptance Criteria:**  
- I must provide reservation time, number of guests, and phone number.  

---

## Get All Reservations for a Restaurant  
**User Story (Staff/Manager):**  
As restaurant staff or a manager, I want to see all reservations for a restaurant to manage seat assignments and scheduling.  

**Acceptance Criteria:**  
- I must provide the restaurant ID.  
- I can optionally filter by date (query parameter).  

---

## Get a Single Reservation  
**User Story (Staff/Manager/Customer who made it):**  
As the reservation owner or restaurant staff/manager, I want to view the details of a specific reservation so that I can confirm the reservation or see the booking information.  

---

## Landing Page Development  
**User Story:**  
As a visitor, I want to see a welcoming landing page so that I can learn about the different restaurants.  

**Acceptance Criteria:**  
- Basic navigation for other sections of the site and the page should be optimized for mobile users.  

---

## Backend Environment Setup  
**User Story:**  
As a developer, set up the environment for the backend and frontend and integrate them.  

**Acceptance Criteria:**  
- Setup of Node, VS Code, Angular CLI, Git setup/Sourcetree.  

---

# Issues Team Planned to Address  

The issues we aim to address with this service revolve around the inefficiencies and frustrations associated with restaurant reservations. Customers often struggle to find the right restaurant, check table availability, and secure a reservation without unnecessary delays. Our platform simplifies this process by allowing users to browse restaurant details, check availability, and make reservations seamlessly. Additionally, if a restaurant is fully booked, customers can join a waitlist, ensuring they have a fair chance to secure a table in case of cancellations.  

For restaurant owners and staff, our system provides tools to manage tables, reservations, and seating capacity efficiently, ensuring better organization and a smoother dining experience. By bridging the gap between restaurants and customers, our service fosters transparency, enhances convenience, and improves the overall dining experience for both parties.  

---

# Progress So Far  

![image](https://github.com/user-attachments/assets/198539fb-ce91-4d07-816f-4df4115f8261)  

### Register a New User  
**User Story:**  
As a new customer, I want to create an account so that I can securely access the application’s features.  

**Acceptance Criteria:**  
- I must provide all required fields (username, password, email, phone number, first name, last name).  

**Status:** *Partially complete*  

---

### Login for Existing User  
**User Story:**  
As an existing user (customer/staff/manager/admin), I want to log in so that I can securely access my account and the relevant system features.  

**Acceptance Criteria:**  
- I must provide valid credentials (username/password/timestamp).  
- Set up authentication on sign-in, proper error handling, and validation should be implemented.  

**Status:** *Partially complete*  

---

### Environment Setup  
**User Story:**  
As a developer, setup the environment for the backend and frontend and integrate them.  

**Acceptance Criteria:**  
- Setup of Node, VS Code, Angular CLI, Git setup/Sourcetree.  

**Status:** *Complete*  

---

# Frontend  

- Landing Page  
- Register  
- Login  

---

# Backend  

- Designed the schema  
- Added registration functionality  

---

## Challenges Faced  

We were unable to fully complete the user registration and login functionalities due to the complexities involved in implementing secure authentication and validation mechanisms. While we managed to set up the basic structure for user registration, additional work is needed to ensure proper validation of user inputs, secure password handling, and error management.  

Similarly, for the login feature, integrating authentication, handling incorrect credentials gracefully, and ensuring session management require further refinement. These tasks demand careful implementation to maintain security and usability, which has slowed progress. However, with the backend environment fully set up, we are now in a strong position to finalize these features efficiently.  
