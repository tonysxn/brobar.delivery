"use client";

import { CheckCircle2 } from "lucide-react";
import Link from "next/link";
import { Button } from "@/components/ui/button";

export default function OrderSuccessPage() {
    return (
        <div className="min-h-[60vh] flex flex-col items-center justify-center text-center px-4 animate-in fade-in zoom-in duration-500">
            <div className="w-24 h-24 bg-green-500/10 rounded-full flex items-center justify-center mb-6 ring-1 ring-green-500/20 shadow-[0_0_40px_-10px_rgba(34,197,94,0.3)]">
                <CheckCircle2 className="w-12 h-12 text-green-500" />
            </div>

            <h1 className="text-3xl font-bold mb-4 text-white">Замовлення оплачено!</h1>
            <p className="text-gray-400 max-w-md mb-8 text-lg">
                Дякуємо за ваше замовлення
            </p>

            <Link href="/">
                <Button size="lg" className="cursor-pointer bg-primary hover:bg-primary/90 text-black font-bold h-12 px-8 text-base rounded-xl">
                    Повернутися на головну
                </Button>
            </Link>
        </div>
    );
}
