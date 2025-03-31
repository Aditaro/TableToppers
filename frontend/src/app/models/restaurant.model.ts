// export interface RestaurantCreate {
//   name: string;
//   location: string;
//   img?: File|undefined;
//   description?: string;
//   phone?: string;
//   openingHours?: string;
// }
// export interface SpecialAvailability {
//   date: string;              // "2025-12-25"
//   reason: string;            // "Holiday", "Private Event"
//   status: 'open' | 'closed' | 'limited';
// }

// export interface Restaurant {
//   id: string;
//   status: 'pending' | 'open' | 'closed';
//   name: string;
//   img: string;
//   description: string;
//   location: string;
//   phone: string;
//   openingHours: string;
//   specialAvailability: SpecialAvailability[];
// }

export interface RestaurantCreate {
  name: string;
  location: string;
  img?: string; // Store the image filename or URL (string) instead of a File object
  description?: string;
  phone?: string;
  openingHours?: string;
}

export interface SpecialAvailability {
  date: string;              // "2025-12-25"
  reason: string;            // "Holiday", "Private Event"
  status: 'open' | 'closed' | 'limited';
}

export interface Restaurant {
  id: string;
  status: 'pending' | 'open' | 'closed';
  name: string;
  img: string;  // Store the image URL or filename here, no longer a File object
  description: string;
  location: string;
  phone: string;
  openingHours: string;
  specialAvailability: SpecialAvailability[];
}

export interface NewRestaurant {
  name: string;
  img: string;
  description: string;
  location: string;
  phone: string;
  openingHours: string;
}
