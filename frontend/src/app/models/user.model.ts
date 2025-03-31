export interface User {
  id: string;
  email: string;
  phoneNumber: string;
  firstName: string;
  lastName: string;
  role: 'customer' | 'staff' | 'manager' | 'admin';
}
