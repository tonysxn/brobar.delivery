"use client";

import { Clock, Calendar as CalendarIcon } from "lucide-react";
import { format } from "date-fns";
import { uk } from "date-fns/locale";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Calendar } from "@/components/ui/calendar";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";

interface TimeSelectionSectionProps {
    deliveryMethod: "delivery" | "pickup";
    isTodayClosed: boolean;
    isAsap: boolean;
    setIsAsap: (val: boolean) => void;
    date: Date | undefined;
    setDate: (val: Date | undefined) => void;
    timeVal: string;
    setTimeVal: (val: string) => void;
    minDate: Date;
    minTime: string;
    maxTime: string;
    isTimeValid: boolean;
}

export function TimeSelectionSection({
    deliveryMethod,
    isTodayClosed,
    isAsap, setIsAsap,
    date, setDate,
    timeVal, setTimeVal,
    minDate, minTime, maxTime,
    isTimeValid
}: TimeSelectionSectionProps) {
    const inputClasses = "w-full h-12 bg-white/5 border border-white/10 rounded-xl px-4 py-3 focus:outline-none focus:border-primary transition-colors text-white placeholder:text-gray-500";

    return (
        <div className="space-y-2">
            <label className="text-sm font-medium text-gray-300 block mb-1">На коли</label>

            {isTodayClosed && (
                <div className="bg-yellow-500/10 border border-yellow-500/30 rounded-xl p-3 mb-3 flex items-center gap-2">
                    <Clock className="w-5 h-5 text-yellow-500 shrink-0" />
                    <span className="text-yellow-500 text-sm">
                        {deliveryMethod === "delivery" ? "Доставка" : "Самовивіз"} на сьогодні вже {deliveryMethod === "delivery" ? "недоступна" : "недоступний"}. Ви можете замовити на завтра.
                    </span>
                </div>
            )}

            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <button
                    onClick={() => {
                        if (isTodayClosed) return;
                        setIsAsap(true);
                        setDate(undefined);
                        setTimeVal("");
                    }}
                    disabled={isTodayClosed}
                    className={cn(
                        "py-3 px-4 rounded-xl border transition-all flex items-center justify-center gap-2 cursor-pointer h-12",
                        isAsap && !isTodayClosed
                            ? "bg-primary/20 border-primary text-primary font-medium"
                            : "bg-white/5 border-white/10 text-gray-400 hover:bg-white/10",
                        isTodayClosed && "opacity-50 cursor-not-allowed"
                    )}
                >
                    <span>Якомога швидше</span>
                </button>

                <Popover>
                    <PopoverTrigger asChild>
                        <Button
                            variant={"outline"}
                            className={cn(
                                "w-full h-12 justify-start text-left font-normal bg-white/5 border-white/10 hover:bg-white/10 hover:text-white rounded-xl transition-all",
                                !date && "text-muted-foreground",
                                isAsap && !isTodayClosed && "opacity-50"
                            )}
                            onClick={() => setIsAsap(false)}
                        >
                            <CalendarIcon className="mr-2 h-4 w-4" />
                            {date ? format(date, "P", { locale: uk }) : <span>Дата</span>}
                        </Button>
                    </PopoverTrigger>
                    <PopoverContent className="w-auto p-0" align="start">
                        <Calendar
                            mode="single"
                            selected={date}
                            onSelect={(d: Date | undefined) => {
                                setDate(d);
                                setIsAsap(false);
                            }}
                            disabled={(d) => d < new Date(minDate.getFullYear(), minDate.getMonth(), minDate.getDate())}
                            initialFocus
                            locale={uk}
                        />
                    </PopoverContent>
                </Popover>

                <div className={cn(
                    "relative rounded-xl border transition-all bg-white/5 border-white/10 hover:bg-white/10",
                    isAsap && !isTodayClosed && "opacity-50",
                    !isTimeValid && timeVal && "border-red-500"
                )}>
                    <Clock className="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
                    <input
                        type="time"
                        value={timeVal}
                        min={minTime}
                        max={maxTime}
                        onFocus={() => setIsAsap(false)}
                        onClick={(e) => {
                            setIsAsap(false);
                            // @ts-ignore
                            if (e.target.showPicker) e.target.showPicker();
                        }}
                        onChange={(e) => {
                            setTimeVal(e.target.value);
                            setIsAsap(false);
                        }}
                        className={cn(
                            "w-full h-12 bg-transparent border-none focus:ring-0 pl-11 pr-4 text-white placeholder:text-gray-500",
                            "min-w-0",
                            "[&::-webkit-calendar-picker-indicator]:hidden [&::-webkit-calendar-picker-indicator]:appearance-none"
                        )}
                    />
                </div>
            </div>

            {!isTimeValid && timeVal && (
                <p className="text-red-500 text-xs mt-1">
                    Виберіть час в межах {minTime} - {maxTime}
                </p>
            )}
        </div>
    );
}
