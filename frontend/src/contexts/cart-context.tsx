"use client"

import React, { createContext, useContext, useState, useEffect, ReactNode } from "react";

// Types
interface Variation {
    id: string;
    name: string;
    external_id: string;
    show: boolean;
}

interface VariationGroup {
    id: string;
    name: string;
    external_id: string;
    show: boolean;
    required: boolean;
    variations: Variation[];
}

interface Product {
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

export interface CartItem {
    id: string;
    product: Product;
    selectedVariations: Record<string, Variation>;
    quantity: number;
}

interface CartContextType {
    cart: CartItem[];
    addToCart: (product: Product, selectedVariations: Record<string, Variation>, quantity?: number) => void;
    updateQuantity: (cartItemId: string, delta: number) => void;
    removeFromCart: (cartItemId: string) => void;
    clearCart: () => void;
    cartTotal: number;
    cartItemCount: number;
    isCartOpen: boolean;
    setIsCartOpen: (open: boolean) => void;
    isInitialized: boolean;
}

const CartContext = createContext<CartContextType | undefined>(undefined);

function generateCartItemId(productId: string, selectedVariations: Record<string, Variation>): string {
    const variationIds = Object.values(selectedVariations)
        .map(v => v.id)
        .sort()
        .join("-");
    return variationIds ? `${productId}-${variationIds}` : productId;
}

const CART_STORAGE_KEY = "brobar_cart";
const CART_EXPIRATION_MS = 24 * 60 * 60 * 1000; // 1 day

export function CartProvider({ children }: { children: ReactNode }) {
    const [cart, setCart] = useState<CartItem[]>([]);
    const [isCartOpen, setIsCartOpen] = useState(false);
    const [isInitialized, setIsInitialized] = useState(false);

    // Initial load from localStorage
    useEffect(() => {
        if (typeof window !== "undefined") {
            const storedCart = localStorage.getItem(CART_STORAGE_KEY);
            if (storedCart) {
                try {
                    const parsed = JSON.parse(storedCart);
                    const now = new Date().getTime();

                    if (parsed.timestamp && now - parsed.timestamp < CART_EXPIRATION_MS) {
                        setCart(parsed.cart);
                    } else {
                        localStorage.removeItem(CART_STORAGE_KEY);
                    }
                } catch (e) {
                    console.error("Failed to parse cart from localStorage", e);
                }
            }
            setIsInitialized(true);
        }
    }, []);

    // Save to localStorage on change
    useEffect(() => {
        if (isInitialized && typeof window !== "undefined") {
            const data = {
                cart,
                timestamp: new Date().getTime()
            };
            localStorage.setItem(CART_STORAGE_KEY, JSON.stringify(data));
        }
    }, [cart, isInitialized]);

    const addToCart = (product: Product, selectedVariations: Record<string, Variation>, quantity: number = 1) => {
        const cartItemId = generateCartItemId(product.id, selectedVariations);

        setCart(prevCart => {
            const existingItem = prevCart.find(item => item.id === cartItemId);
            if (existingItem) {
                return prevCart.map(item =>
                    item.id === cartItemId
                        ? { ...item, quantity: item.quantity + quantity }
                        : item
                );
            } else {
                return [...prevCart, {
                    id: cartItemId,
                    product,
                    selectedVariations,
                    quantity
                }];
            }
        });
    };

    const updateQuantity = (cartItemId: string, delta: number) => {
        setCart(prevCart => {
            return prevCart
                .map(item => {
                    if (item.id === cartItemId) {
                        const newQuantity = item.quantity + delta;
                        return newQuantity > 0 ? { ...item, quantity: newQuantity } : null;
                    }
                    return item;
                })
                .filter((item): item is CartItem => item !== null);
        });
    };

    const removeFromCart = (cartItemId: string) => {
        setCart(prevCart => prevCart.filter(item => item.id !== cartItemId));
    };

    const clearCart = () => {
        setCart([]);
    };

    const cartTotal = cart.reduce((sum, item) => sum + item.product.price * item.quantity, 0);
    const cartItemCount = cart.reduce((sum, item) => sum + item.quantity, 0);

    return (
        <CartContext.Provider value={{
            cart,
            addToCart,
            updateQuantity,
            removeFromCart,
            clearCart,
            cartTotal,
            cartItemCount,
            isCartOpen,
            setIsCartOpen,
            isInitialized
        }}>
            {children}
        </CartContext.Provider>
    );
}

export function useCart() {
    const context = useContext(CartContext);
    if (context === undefined) {
        throw new Error("useCart must be used within a CartProvider");
    }
    return context;
}

export { generateCartItemId };
export type { Product, Variation, VariationGroup };
