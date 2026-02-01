"use client"

import { useEffect, useState } from "react";
import { useCart } from "@/contexts/cart-context";
import { useShopStatus } from "@/hooks/use-shop-status";
import { X, ShoppingBag, Trash2, Plus, Minus } from "lucide-react";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import { toast } from "sonner";

const FILE_URL = process.env.NEXT_PUBLIC_FILE_URL || "http://localhost:3001";

function getImageUrl(imagePath: string) {
    if (!imagePath) return "https://placehold.co/80x80?text=No+Image";
    return `${FILE_URL}/${imagePath}`;
}

export default function CartDrawer() {
    const { cart, isCartOpen, setIsCartOpen, removeFromCart, cartTotal, cartItemCount, updateQuantity } = useCart();
    const [isFirefox, setIsFirefox] = useState(false);

    useEffect(() => {
        setIsFirefox(typeof navigator !== 'undefined' && /firefox/i.test(navigator.userAgent));
    }, []);

    useEffect(() => {
        if (isCartOpen) {
            toast.dismiss();
            document.body.style.overflow = "hidden";
        } else {
            document.body.style.overflow = "";
        }

        return () => {
            document.body.style.overflow = "";
        };
    }, [isCartOpen]);

    const backdropClass = isFirefox
        ? "bg-black/70"
        : "bg-black/60 backdrop-blur-sm";

    const drawerClass = isFirefox
        ? "bg-[#141414]"
        : "bg-black/95 backdrop-blur-xl supports-[backdrop-filter]:bg-black/50";

    return (
        <div className={`fixed inset-0 z-[60] flex justify-end transition-all duration-150 ${isCartOpen ? "visible" : "invisible"}`}>
            {/* Backdrop with blur */}
            <div
                className={`absolute inset-0 ${backdropClass} transition-opacity duration-150 will-change-opacity ${isCartOpen ? "opacity-100" : "opacity-0"}`}
                onClick={() => setIsCartOpen(false)}
            />

            {/* Drawer */}
            <div className={`relative w-full max-w-[420px] h-full ${drawerClass} border-l border-white/10 shadow-2xl flex flex-col transition-transform duration-200 ease-out will-change-transform ${isCartOpen ? "translate-x-0" : "translate-x-full"}`}>
                {/* Header */}
                <div className="flex items-center justify-between p-5 border-b border-white/10">
                    <div className="flex items-center gap-3">
                        <div className="w-10 h-10 rounded-full bg-gradient-to-br from-primary to-primary/80 flex items-center justify-center">
                            <ShoppingBag className="w-5 h-5 text-black" />
                        </div>
                        <div>
                            <h2 className="text-xl font-bold">Кошик</h2>
                            <p className="text-sm text-gray-400">{cartItemCount} {cartItemCount === 1 ? "товар" : cartItemCount < 5 ? "товари" : "товарів"}</p>
                        </div>
                    </div>
                    <button
                        onClick={() => setIsCartOpen(false)}
                        className="p-2 hover:bg-white/10 rounded-full transition-colors duration-100 cursor-pointer"
                    >
                        <X className="w-6 h-6" />
                    </button>
                </div>

                {/* Cart Items */}
                <div className="flex-1 overflow-y-auto p-4 space-y-3">
                    {cart.length === 0 ? (
                        <div className="flex flex-col items-center justify-center h-full text-center">
                            <div className="w-24 h-24 rounded-full bg-white/5 flex items-center justify-center mb-4">
                                <ShoppingBag className="w-12 h-12 text-gray-500" />
                            </div>
                            <h3 className="text-lg font-medium text-gray-300">Кошик порожній</h3>
                            <p className="text-sm text-gray-500 mt-1">Додайте товари з меню</p>
                        </div>
                    ) : (
                        cart.map((item) => (
                            <div
                                key={item.id}
                                className="bg-white/5 hover:bg-white/10 rounded-xl p-3 transition-all duration-150"
                            >
                                <div className="flex items-center gap-3">
                                    {/* Delete button - Left */}
                                    <button
                                        onClick={() => removeFromCart(item.id)}
                                        className="p-2 text-gray-400 hover:text-red-400 hover:bg-red-400/10 rounded-full transition-all cursor-pointer shrink-0"
                                    >
                                        <Trash2 className="w-5 h-5" />
                                    </button>

                                    {/* Product Image */}
                                    <div className="relative w-16 h-16 rounded-lg overflow-hidden shrink-0">
                                        <img
                                            src={getImageUrl(item.product.image)}
                                            alt={item.product.name}
                                            className="w-full h-full object-cover"
                                        />
                                        <div className="absolute inset-0 bg-gradient-to-t from-black/40 to-transparent" />
                                    </div>

                                    {/* Product Info */}
                                    <div className="flex-1 min-w-0">
                                        <h4 className="font-medium text-sm truncate">{item.product.name}</h4>

                                        {/* Variations */}
                                        {Object.values(item.selectedVariations).length > 0 && (
                                            <p className="text-xs text-primary/80 mt-0.5 truncate">
                                                {Object.values(item.selectedVariations).map(v => v.name).join(", ")}
                                            </p>
                                        )}

                                        {/* Description */}
                                        {item.product.description && (
                                            <p className="text-xs text-gray-400 mt-1 line-clamp-1">
                                                {item.product.description}
                                            </p>
                                        )}

                                        {/* Price */}
                                        <p className="text-primary font-bold mt-1">
                                            {item.product.price * item.quantity} ₴
                                        </p>
                                    </div>

                                    {/* Quantity Controls - Right */}
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
                            </div>
                        ))
                    )}
                </div>

                {/* Footer */}
                {cart.length > 0 && (
                    <div className="p-5 border-t border-white/10 space-y-4">
                        <div className="flex items-center justify-between">
                            <span className="text-gray-400">Разом:</span>
                            <span className="text-2xl font-bold">{cartTotal} ₴</span>
                        </div>
                        <CheckoutButton setIsCartOpen={setIsCartOpen} />
                    </div>
                )}
            </div>
        </div>
    );
}

function CheckoutButton({ setIsCartOpen }: { setIsCartOpen: (open: boolean) => void }) {
    const { isPaused, message } = useShopStatus();

    if (isPaused) {
        return (
            <div className="space-y-2">
                <Button disabled className="w-full h-14 text-lg font-bold bg-muted text-muted-foreground cursor-not-allowed">
                    Продажі призупинено
                </Button>
                <p className="text-xs text-red-500 text-center">{message}</p>
            </div>
        );
    }

    return (
        <Link href="/order" onClick={() => setIsCartOpen(false)} className="block w-full">
            <Button className="w-full h-14 text-lg font-bold text-black bg-primary hover:bg-primary/90 rounded-xl shadow-lg shadow-primary/25 transition-all duration-150 hover:shadow-primary/40 hover:scale-[1.02] cursor-pointer">
                Оформити замовлення
            </Button>
        </Link>
    );
}
