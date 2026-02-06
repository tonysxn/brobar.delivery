"use client";

import { Utensils } from "lucide-react";
import { cn } from "@/lib/utils";

interface ContactInfoSectionProps {
    name: string;
    setName: (val: string) => void;
    phone: string;
    setPhone: (val: string) => void;
    email: string;
    setEmail: (val: string) => void;
    cutleryCount: number;
    setCutleryCount: (val: number) => void;
    isPhoneValid: boolean;
    isEmailValid: boolean;
}

export function ContactInfoSection({
    name, setName,
    phone, setPhone,
    email, setEmail,
    cutleryCount, setCutleryCount,
    isPhoneValid, isEmailValid
}: ContactInfoSectionProps) {
    const inputClasses = "w-full h-12 bg-white/5 border border-white/10 rounded-xl px-4 py-3 focus:outline-none focus:border-primary transition-colors text-white placeholder:text-gray-500";

    return (
        <div className="space-y-6">
            <h2 className="text-xl font-bold">Контактні дані</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="space-y-2">
                    <label className="text-sm font-medium text-gray-300 block mb-1">Ім'я <span className="text-red-400">*</span></label>
                    <input
                        type="text"
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        className={inputClasses}
                        placeholder="Ваше ім'я"
                    />
                </div>
                <div className="space-y-2">
                    <label className="text-sm font-medium text-gray-300 block mb-1">
                        Телефон <span className="text-red-400">*</span>
                        {!isPhoneValid && <span className="text-red-400 text-xs ml-2">Невірний формат</span>}
                    </label>
                    <input
                        type="tel"
                        value={phone}
                        onChange={(e) => setPhone(e.target.value)}
                        className={cn(inputClasses, !isPhoneValid && "border-red-500/50 focus:border-red-500")}
                        placeholder="+380..."
                    />
                </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="space-y-2">
                    <label className="text-sm font-medium text-gray-300 block mb-1">Кількість приборів <span className="text-red-400">*</span></label>
                    <div className="relative">
                        <Utensils className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-500" />
                        <input
                            type="number"
                            min="0"
                            value={cutleryCount}
                            onChange={(e) => setCutleryCount(parseInt(e.target.value) || 0)}
                            className={cn(inputClasses, "pl-12")}
                        />
                    </div>
                </div>
                <div className="space-y-2">
                    <label className="text-sm font-medium text-gray-300 block mb-1">
                        E-Mail (не обов'язково)
                        {!isEmailValid && <span className="text-red-400 text-xs ml-2">Невірний формат</span>}
                    </label>
                    <input
                        type="email"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        className={cn(inputClasses, !isEmailValid && "border-red-500/50 focus:border-red-500")}
                        placeholder="example@mail.com"
                    />
                </div>
            </div>
        </div>
    );
}
