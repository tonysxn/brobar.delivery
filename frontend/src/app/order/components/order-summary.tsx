"use client";

import { Separator } from "@/components/ui/separator";
import { Button } from "@/components/ui/button";
import { formatPrice } from "@/app/order/utils";

interface OrderSummaryProps {
    cartTotal: number;
    deliveryMethod: "delivery" | "pickup";
    isFreeDelivery: boolean;
    deliveryPrice: number;
    toDoor: boolean;
    toDoorPrice: number;
    total: number;
    isValid: boolean;
    handleSubmit: () => void;
}

export function OrderSummary({
    cartTotal,
    deliveryMethod,
    isFreeDelivery,
    deliveryPrice,
    toDoor,
    toDoorPrice,
    total,
    isValid,
    handleSubmit
}: OrderSummaryProps) {
    return (
        <section className="bg-white/5 rounded-2xl p-3 md:p-6 border border-white/10 space-y-4">
            <h2 className="text-xl font-bold">Разом</h2>

            <div className="space-y-2 text-sm">
                <div className="flex justify-between text-gray-300">
                    <span>Вартість товарів:</span>
                    <span>{formatPrice(cartTotal)}</span>
                </div>

                {deliveryMethod === "delivery" && (
                    <>
                        <div className="flex justify-between text-gray-300">
                            <span>Доставка:</span>
                            {isFreeDelivery ? (
                                <span className="text-green-400">Безкоштовно</span>
                            ) : (
                                <span>{formatPrice(deliveryPrice)}</span>
                            )}
                        </div>
                        {toDoor && (
                            <div className="flex justify-between text-gray-300">
                                <span>Доставка до дверей:</span>
                                <span>{formatPrice(toDoorPrice)}</span>
                            </div>
                        )}
                    </>
                )}
            </div>

            <Separator className="bg-white/10" />

            <div className="flex justify-between items-end">
                <span className="font-bold text-lg">До сплати:</span>
                <span className="font-bold text-3xl text-primary">{formatPrice(total)}</span>
            </div>

            <Button
                onClick={handleSubmit}
                disabled={!isValid}
                className="w-full h-14 text-lg font-bold text-black bg-primary hover:bg-primary/90 rounded-xl shadow-lg shadow-primary/25 transition-all hover:scale-[1.02] active:scale-[0.98] disabled:opacity-50 disabled:pointer-events-none cursor-pointer"
            >
                ЗАМОВИТИ
            </Button>

            <p className="text-xs text-center text-gray-500 mt-4">
                Натискаючи кнопку, ви погоджуєтесь з умовами публічної оферти
            </p>
        </section>
    );
}
