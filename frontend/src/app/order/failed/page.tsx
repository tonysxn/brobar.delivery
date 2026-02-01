"use client";

import { XCircle } from "lucide-react";
import Link from "next/link";
import { Button } from "@/components/ui/button";

export default function OrderFailedPage() {
    return (
        <div className="min-h-[60vh] flex flex-col items-center justify-center text-center px-4 animate-in fade-in zoom-in duration-500">
            <div className="w-24 h-24 bg-red-500/10 rounded-full flex items-center justify-center mb-6 ring-1 ring-red-500/20 shadow-[0_0_40px_-10px_rgba(239,68,68,0.3)]">
                <XCircle className="w-12 h-12 text-red-500" />
            </div>

            <h1 className="text-3xl font-bold mb-4 text-white">Помилка оплати</h1>
            <p className="text-gray-400 max-w-md mb-8 text-lg">
                На жаль, платіж не пройшов або був скасований. <br />
                Спробуйте ще раз або оберіть інший спосіб оплати.
            </p>

            <div className="flex gap-4">
                <Link href="/menu">
                    <Button variant="outline" size="lg" className="border-white/10 hover:bg-white/5 text-white h-12 px-8 text-base rounded-xl">
                        В меню
                    </Button>
                </Link>
                <Link href="/order">
                    <Button size="lg" className="bg-primary hover:bg-primary/90 text-black font-bold h-12 px-8 text-base rounded-xl">
                        Спробувати знову
                    </Button>
                </Link>
            </div>
        </div>
    );
}
