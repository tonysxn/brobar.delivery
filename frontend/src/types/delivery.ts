export interface DeliveryZone {
    radius: number;
    innerRadius: number;
    color: string;
    price: number;
    freeOrderPrice: number;
    name: string;
}

export interface SearchResult {
    zone?: DeliveryZone;
    distance?: number;
    address?: string;
    coords?: { lat: number; lng: number };
}
