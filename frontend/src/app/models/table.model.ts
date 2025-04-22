export interface Table {
  id: string;
  restaurantId: string;
  name?: string;
  minCapacity: number;
  maxCapacity: number;
  status: 'available' | 'occupied' | 'reserved';
  x?: number;
  y?: number;
}

export interface NewTable {
  name?: string;
  minCapacity: number;
  maxCapacity: number;
  isEdit?: boolean;
}
