"use client";

import { Clock2 } from "lucide-react";
import { cn } from "@/lib/utils";

interface ShopStatusAlertProps {
    deliveryOpen: boolean;
    pickupOpen: boolean;
    isPaused?: boolean;
    className?: string;
}

export function ShopStatusAlert({ deliveryOpen, pickupOpen, isPaused, className }: ShopStatusAlertProps) {
    if (deliveryOpen && pickupOpen && !isPaused) return null;

    let message = "";
    if (isPaused) {
        message = "Вибачте, ми тимчасово не приймаємо замовлення 😔";
    } else if (!deliveryOpen && !pickupOpen) {
        message = "Сьогодні ми вже не працюємо. Ви можете замовити на завтра!";
    } else if (!deliveryOpen) {
        message = "Доставка на сьогодні недоступна. Ви можете замовити самовивіз або на завтра.";
    } else {
        message = "Самовивіз на сьогодні недоступний. Ви можете замовити доставку або на завтра.";
    }

    return (
        <div
            className={cn(
                "bg-yellow-500/10 border-b border-yellow-500/30 px-4 py-3 flex items-center justify-center gap-2",
                className
            )}
        >
            <Clock2 className="w-5 h-5 text-yellow-500 shrink-0" />
            <span className="text-yellow-500 text-sm text-center">
                {message}
            </span>
        </div>
    );
}
