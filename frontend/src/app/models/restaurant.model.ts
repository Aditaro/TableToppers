export interface RestaurantCreate {
  name: string;
  location: string;
  img?: File|undefined;
  description?: string;
  phone?: string;
  openingHours?: string;
}