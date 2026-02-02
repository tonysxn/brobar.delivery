"use client";

import { Clock2 } from "lucide-react";
import { cn } from "@/lib/utils";

interface ShopStatusAlertProps {
    deliveryOpen: boolean;
    pickupOpen: boolean;
    className?: string;
}

export function ShopStatusAlert({ deliveryOpen, pickupOpen, className }: ShopStatusAlertProps) {
    if (deliveryOpen && pickupOpen) return null;

    return (
        <div
            className={cn(
                "bg-yellow-500/10 border-b border-yellow-500/30 px-4 py-3 flex items-center justify-center gap-2",
                className
            )}
        >
            <Clock2 className="w-5 h-5 text-yellow-500 shrink-0" />
            <span className="text-yellow-500 text-sm text-center">
                {!deliveryOpen && !pickupOpen
                    ? "Сьогодні ми вже не працюємо. Ви можете замовити на завтра!"
                    : !deliveryOpen
                        ? "Доставка на сьогодні недоступна. Ви можете замовити самовивіз або на завтра."
                        : "Самовивіз на сьогодні недоступний. Ви можете замовити доставку або на завтра."}
            </span>
        </div>
    );
}
