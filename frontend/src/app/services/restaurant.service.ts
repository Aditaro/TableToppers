// import { Injectable } from '@angular/core';
// import { Observable, of } from 'rxjs';
// import { Restaurant, RestaurantCreate } from '../models/restaurant.model';
// import { HttpClient } from '@angular/common/http';

// @Injectable({
//   providedIn: 'root'
// })
// export class RestaurantService {
//   private mockRestaurants: Restaurant[] = [
//     {
//       id: '1',
//       status: 'open',
//       name: 'Sushi Samba',
//       img: 'https://picsum.photos/300/200?random=1',
//       description: 'A delightful mix of Japanese and Brazilian cuisine.',
//       location: 'Tokyo',
//       phone: '+81 123 456 789',
//       openingHours: '11:00 - 21:00',
//       specialAvailability: [
//         {
//           date: '2025-12-25',
//           reason: 'Christmas Holiday',
//           status: 'closed'
//         },
//         {
//           date: '2025-01-01',
//           reason: 'New Year Holiday',
//           status: 'closed'
//         }
//       ]
//     },
//     {
//       id: '2',
//       status: 'closed',
//       name: 'Pasta Palace',
//       img: 'https://picsum.photos/300/200?random=2',
//       description: 'Authentic Italian pasta and pizza.',
//       location: 'Rome',
//       phone: '+39 06 1234 5678',
//       openingHours: '10:00 - 23:00',
//       specialAvailability: []
//     },
//     {
//       id: '3',
//       status: 'pending',
//       name: 'Burger Bonanza',
//       img: 'https://picsum.photos/300/200?random=3',
//       description: 'Classic American burgers with a modern twist.',
//       location: 'New York',
//       phone: '+1 (555) 123-4567',
//       openingHours: '11:00 - 22:00',
//       specialAvailability: [
//         {
//           date: '2025-07-04',
//           reason: '4th of July Celebration',
//           status: 'limited'
//         }
//       ]
//     }
//   ];

//   // constructor() { }
//   constructor(private http: HttpClient, private router: Router) {}

//   getRestaurants(city?: string, name?: string): Observable<Restaurant[]> {
//     let filteredData = this.mockRestaurants;

//     if (city && city.trim().length > 0) {
//       filteredData = filteredData.filter(r =>
//         r.location.toLowerCase().includes(city.toLowerCase())
//       );
//     }

//     if (name && name.trim().length > 0) {
//       filteredData = filteredData.filter(r =>
//         r.name.toLowerCase().includes(name.toLowerCase())
//       );
//     }

//     return of(filteredData);
//   }

//   // createRestaurant(data: RestaurantCreate): Observable<any> {
//   //   // Mock: simply push to local array and return it
//   //   const newRestaurant: Restaurant = {
//   //     name: data.name,
//   //     location: data.location,
//   //     phone: data.phone || '',
//   //     openingHours: data.openingHours || '',
//   //     id: (this.mockRestaurants.length + 1).toString(),
//   //     status: 'pending',
//   //     specialAvailability: [],
//   //     img: 'https://picsum.photos/300/200?random=4',
//   //     description: data.description || ''
//   //   };
//   //   this.mockRestaurants.push(newRestaurant);
//   //   return of(data);
//   // }
//   createRestaurant(data: RestaurantCreate): Observable<any> {
//     const formData = new FormData();
//     formData.append('name', data.name);
//     formData.append('location', data.location);
//     formData.append('description', data.description || '');
//     formData.append('phone', data.phone || '');
//     formData.append('openingHours', data.openingHours || '');
    
//     if (data.img) {
//       formData.append('img', data.img);
//     }
  
//     return this.http.post('http://localhost:8080/restaurants', formData);
//   }
  
// }
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Restaurant, RestaurantCreate } from '../models/restaurant.model';
import { HttpClient, HttpHeaders } from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class RestaurantService {
  private supabaseUrl = process.env.SUPABASE_URL as string;
  private supabaseApiKey = process.env.SUPABASE_API_KEY as string;
  
  private headers = new HttpHeaders({
    'apikey': this.supabaseApiKey,
    'Authorization': `Bearer ${this.supabaseApiKey}`,
    'Content-Type': 'application/json'
  });

  constructor(private http: HttpClient) {}

  // Fetch restaurants from Supabase
  getRestaurants(city?: string, name?: string): Observable<Restaurant[]> {
    // Construct query parameters dynamically based on provided filters
    let queryParams = '';
    
    if (city) {
      queryParams += `?location=ilike.${city}`;
    }
    
    if (name) {
      // If city filter is already applied, append `&`, otherwise start with `?`
      queryParams += (queryParams ? `&` : `?`) + `name=ilike.${name}`;
    }
  
    return this.http.get<Restaurant[]>(`${this.supabaseUrl}/rest/v1/restaurants${queryParams}`, { headers: this.headers });
  }
  

  // Add a new restaurant to Supabase
  createRestaurant(data: RestaurantCreate): Observable<any> {
    const newRestaurant = {
      name: data.name,
      location: data.location,
      description: data.description || '',
      phone: data.phone || '',
      opening_hours: data.openingHours || '',
      img: data.img ? data.img.name : '' // Store image filename (actual upload is separate)
    };

    return this.http.post(`${this.supabaseUrl}/rest/v1/restaurants`, newRestaurant, { headers: this.headers });
  }

  // Upload image to Supabase Storage
  uploadImage(file: File): Observable<any> {
    const formData = new FormData();
    formData.append('file', file);
    
    return this.http.post(`${this.supabaseUrl}/storage/v1/object/sign/restaurant_images/${file.name}`, formData, { headers: this.headers });
  }
}
