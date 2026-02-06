"use client";

import { CreditCard, Banknote } from "lucide-react";
import { cn } from "@/lib/utils";

interface PaymentSectionProps {
    paymentMethod: "bank" | "cash";
    setPaymentMethod: (method: "bank" | "cash") => void;
    deliveryMethod: "delivery" | "pickup";
}

export function PaymentSection({ paymentMethod, setPaymentMethod, deliveryMethod }: PaymentSectionProps) {
    return (
        <section className="bg-white/5 rounded-2xl p-3 md:p-6 border border-white/10">
            <h2 className="text-xl font-bold mb-4">Оплата</h2>
            <div className="space-y-3">
                <label className={cn(
                    "flex items-center gap-4 p-4 rounded-xl border cursor-pointer transition-all",
                    paymentMethod === "bank"
                        ? "bg-white/10 border-primary text-white"
                        : "bg-white/5 border-white/10 hover:border-white/20 text-gray-300"
                )}>
                    <input
                        type="radio"
                        name="payment_method"
                        value="bank"
                        checked={paymentMethod === "bank"}
                        onChange={() => setPaymentMethod("bank")}
                        className="w-5 h-5 accent-primary bg-transparent border-white/20"
                    />
                    <CreditCard className={cn("w-6 h-6", paymentMethod === "bank" ? "text-primary" : "text-gray-400")} />
                    <span className={paymentMethod === "bank" ? "font-bold text-white" : "text-gray-300"}>
                        Безготівкова на сайті
                    </span>
                </label>

                {deliveryMethod === "pickup" && (
                    <label className={cn(
                        "flex items-center gap-4 p-4 rounded-xl border cursor-pointer transition-all",
                        paymentMethod === "cash"
                            ? "bg-primary/10 border-primary"
                            : "bg-white/5 border-white/10 hover:border-white/20"
                    )}>
                        <input
                            type="radio"
                            name="payment"
                            value="cash"
                            checked={paymentMethod === "cash"}
                            onChange={() => setPaymentMethod("cash")}
                            className="sr-only"
                        />
                        <Banknote className={cn("w-6 h-6", paymentMethod === "cash" ? "text-primary" : "text-gray-400")} />
                        <span className={paymentMethod === "cash" ? "font-bold text-white" : "text-gray-300"}>
                            Готівкою при отриманні
                        </span>
                    </label>
                )}
            </div>
        </section>
    );
}
