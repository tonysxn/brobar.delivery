export interface Variation {
    id: string;
    name: string;
    external_id: string;
    show: boolean;
}

export interface VariationGroup {
    id: string;
    name: string;
    external_id: string;
    show: boolean;
    required: boolean;
    variations: Variation[];
}

export interface Product {
    id: string;
    external_id: string;
    name: string;
    slug: string;
    description: string | null;
    price: number;
    weight: number | null;
    sort: number;
    hidden: boolean;
    alcohol: boolean;
    sold: boolean;
    image: string;
    variation_groups: VariationGroup[];
}
