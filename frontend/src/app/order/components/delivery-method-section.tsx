"use client";

import { Truck, Store } from "lucide-react";
import DeliveryMap from "@/components/delivery-map";
import { SearchResult } from "@/types/delivery";
import { formatPrice } from "../utils";
import { cn } from "@/lib/utils";

interface DeliveryMethodSectionProps {
    deliveryMethod: "delivery" | "pickup";
    setDeliveryMethod: (method: "delivery" | "pickup") => void;
    setDeliveryResult: (result: SearchResult) => void;
    cartTotal: number;
    deliveryResult: SearchResult | null;
    entrance: string;
    setEntrance: (val: string) => void;
    toDoor: boolean;
    setToDoor: (val: boolean) => void;
    doorPrice: number;
}

export function DeliveryMethodSection({
    deliveryMethod,
    setDeliveryMethod,
    setDeliveryResult,
    cartTotal,
    deliveryResult,
    entrance,
    setEntrance,
    toDoor,
    setToDoor,
    doorPrice
}: DeliveryMethodSectionProps) {
    const inputClasses = "w-full h-12 bg-white/5 border border-white/10 rounded-xl px-4 py-3 focus:outline-none focus:border-primary transition-colors text-white placeholder:text-gray-500";

    return (
        <section className="bg-white/5 rounded-2xl p-3 md:p-6 border border-white/10 space-y-6">
            <h2 className="text-xl font-bold flex items-center gap-2">
                <Truck className="w-5 h-5 text-primary" />
                Спосіб отримання
            </h2>

            <div className="grid grid-cols-2 gap-4">
                <button
                    onClick={() => setDeliveryMethod("delivery")}
                    className={cn(
                        "p-4 rounded-xl border-2 flex flex-col items-center gap-2 transition-all cursor-pointer",
                        deliveryMethod === "delivery"
                            ? "border-primary bg-primary/10 text-white"
                            : "border-white/10 hover:border-white/20 text-gray-400"
                    )}
                >
                    <Truck className="w-6 h-6" />
                    <span className="font-medium">Доставка</span>
                </button>
                <button
                    onClick={() => setDeliveryMethod("pickup")}
                    className={cn(
                        "p-4 rounded-xl border-2 flex flex-col items-center gap-2 transition-all cursor-pointer",
                        deliveryMethod === "pickup"
                            ? "border-primary bg-primary/10 text-white"
                            : "border-white/10 hover:border-white/20 text-gray-400"
                    )}
                >
                    <Store className="w-6 h-6" />
                    <span className="font-medium">Самовивіз</span>
                </button>
            </div>

            {deliveryMethod === "delivery" && (
                <div className="space-y-6 animate-in fade-in slide-in-from-top-4">
                    <DeliveryMap onLocationSelect={setDeliveryResult} cartTotal={cartTotal}>
                        {!!deliveryResult?.address && (
                            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 items-end animate-in fade-in slide-in-from-top-4">
                                <div className="space-y-2">
                                    <label className="text-sm font-medium text-gray-300 block mb-1">Під'їзд / Код</label>
                                    <input
                                        type="text"
                                        value={entrance}
                                        onChange={(e) => setEntrance(e.target.value)}
                                        className={inputClasses}
                                        placeholder="Під'їзд 1, код 123"
                                    />
                                </div>
                                <div className="space-y-2">
                                    <label className="text-sm font-medium text-transparent block mb-1 select-none">Доставка</label>
                                    <div
                                        className="flex items-center space-x-3 bg-white/5 border border-white/10 rounded-xl px-4 h-12 cursor-pointer transition-colors hover:bg-white/10"
                                        onClick={() => setToDoor(!toDoor)}
                                    >
                                        <input
                                            type="checkbox"
                                            id="toDoor"
                                            checked={toDoor}
                                            onChange={(e) => setToDoor(e.target.checked)}
                                            className="w-5 h-5 rounded border-gray-600 text-primary focus:ring-primary bg-transparent cursor-pointer accent-primary"
                                            onClick={(e) => e.stopPropagation()}
                                        />
                                        <label htmlFor="toDoor" className="flex-1 cursor-pointer text-sm font-medium select-none pointer-events-none text-gray-300">
                                            Доставка до дверей (+{formatPrice(doorPrice)})
                                        </label>
                                    </div>
                                </div>
                            </div>
                        )}
                    </DeliveryMap>
                </div>
            )}

            {deliveryMethod === "pickup" && (
                <div className="p-4 bg-primary/10 border border-primary/20 rounded-xl animate-in fade-in">
                    <p className="text-center text-primary font-medium">
                        Адреса бару: вул. Григорія Сковороди 64 (вхід з вул. Багалія)
                    </p>
                </div>
            )}
        </section>
    );
}
