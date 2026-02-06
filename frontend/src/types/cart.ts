import { Product, Variation } from "./product";

export interface CartItem {
    id: string;
    product: Product;
    selectedVariations: Record<string, Variation>;
    quantity: number;
}
