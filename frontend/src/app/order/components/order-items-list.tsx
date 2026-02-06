"use client";

import { Trash2, Plus, Minus, Store } from "lucide-react";
import { ImageWithFallback } from "@/components/ui/image-with-fallback";
import { CartItem } from "@/types/cart";
import { formatPrice } from "../utils";

interface OrderItemsListProps {
    items: CartItem[];
    cartTotal: number;
    removeFromCart: (id: string) => void;
    updateQuantity: (id: string, delta: number) => void;
}

export function OrderItemsList({ items, removeFromCart, updateQuantity }: OrderItemsListProps) {
    return (
        <section className="bg-white/5 rounded-2xl p-3 md:p-6 border border-white/10">
            <h2 className="text-xl font-bold mb-4 flex items-center gap-2">
                <Store className="w-5 h-5 text-primary" />
                Ваше замовлення
            </h2>
            <div className="space-y-4 max-h-[300px] overflow-y-auto pr-2 custom-scrollbar">
                {items.map((item) => (
                    <div key={item.id} className="flex gap-2 md:gap-4 items-center bg-white/5 p-2 md:p-3 rounded-xl">
                        <button
                            onClick={() => removeFromCart(item.id)}
                            className="p-2 text-gray-400 hover:text-red-400 hover:bg-red-400/10 rounded-full transition-all cursor-pointer shrink-0"
                        >
                            <Trash2 className="w-5 h-5" />
                        </button>

                        <div className="w-16 h-16 rounded-lg overflow-hidden shrink-0">
                            <ImageWithFallback
                                src={item.product.image}
                                alt={item.product.name}
                                className="w-full h-full object-cover"
                            />
                        </div>
                        <div className="flex-1 min-w-0">
                            <h4 className="font-medium truncate">{item.product.name}</h4>
                            {Object.values(item.selectedVariations).length > 0 && (
                                <p className="text-xs text-gray-400 truncate">
                                    {Object.values(item.selectedVariations).map(v => v.name).join(", ")}
                                </p>
                            )}
                            <div className="mt-1">
                                <span className="font-bold text-primary">{formatPrice(item.product.price * item.quantity)}</span>
                            </div>
                        </div>

                        <div className="flex flex-col items-center gap-1 bg-white/5 rounded-lg p-1 shrink-0">
                            <button
                                onClick={() => updateQuantity(item.id, 1)}
                                className="p-1 hover:bg-white/10 rounded-md transition-colors cursor-pointer text-gray-300 hover:text-white"
                            >
                                <Plus className="w-4 h-4" />
                            </button>
                            <span className="text-sm font-bold w-6 text-center">{item.quantity}</span>
                            <button
                                onClick={() => updateQuantity(item.id, -1)}
                                className="p-1 hover:bg-white/10 rounded-md transition-colors cursor-pointer text-gray-300 hover:text-white"
                                disabled={item.quantity <= 1}
                            >
                                <Minus className="w-4 h-4" />
                            </button>
                        </div>
                    </div>
                ))}
            </div>
        </section>
    );
}
