export interface WaitlistEntry {
  id?: string;
  restaurantId: string;
  name: string;
  partySize: number;
  partyAhead: number;
  estimatedWaitTime?: number;
  phoneNumber: string;
  status: 'waiting' | 'seated' | 'cancelled';
}

export interface WaitlistEntryCreate {
  name: string;
  partySize: number;
  phoneNumber: string;
}
  