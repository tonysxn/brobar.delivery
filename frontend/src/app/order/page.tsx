"use client";

import { useEffect, useState, useMemo } from "react";
import { useCart } from "@/contexts/cart-context";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import DeliveryMap, { SearchResult } from "@/components/delivery-map";
import { toast } from "sonner";
import { ArrowLeft, MapPin, Truck, Store, CreditCard, Banknote, Utensils, Calendar as CalendarIcon, Trash2, Plus, Minus, Clock } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { cn } from "@/lib/utils";
import { format } from "date-fns";
import { uk } from "date-fns/locale";
import { Calendar } from "@/components/ui/calendar";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { useShopStatus } from "@/hooks/use-shop-status";
import { useSettings } from "@/contexts/settings-context";

// Utility for formatting price
const formatPrice = (price: number) => `${price} ₴`;

const FILE_URL = process.env.NEXT_PUBLIC_FILE_URL || "http://localhost:3001";
function getImageUrl(imagePath: string) {
    if (!imagePath) return "https://placehold.co/80x80?text=No+Image";
    return `${FILE_URL}/files/${imagePath}`;
}

export default function CheckoutPage() {
    const { cart, cartTotal, removeFromCart, updateQuantity, isInitialized } = useCart();
    const router = useRouter();
    const { isPaused, workingHours, serverTime, deliveryOpen, pickupOpen } = useShopStatus();
    const { getSetting } = useSettings();
    const [doorPrice, setDoorPrice] = useState(50);

    // Fetch delivery door price
    useEffect(() => {
        const priceSetting = getSetting("delivery_door_price");
        if (priceSetting && priceSetting.value) {
            const parsed = parseFloat(priceSetting.value);
            if (!isNaN(parsed)) {
                setDoorPrice(parsed);
            }
        }
    }, [getSetting]);

    // Redirect if sales are paused
    useEffect(() => {
        if (isPaused) {
            toast.error("Вибачте, ми тимчасово не приймаємо замовлення.");
            router.push("/menu");
        }
    }, [isPaused, router]);

    // Form State
    const [deliveryMethod, setDeliveryMethod] = useState<"delivery" | "pickup">("delivery");
    const [deliveryResult, setDeliveryResult] = useState<SearchResult | null>(null);

    // Delivery Details
    const [entrance, setEntrance] = useState("");
    const [toDoor, setToDoor] = useState(false);

    // Contact Info
    const [name, setName] = useState("");
    const [phone, setPhone] = useState("");
    const [email, setEmail] = useState("");

    // Time
    const [isAsap, setIsAsap] = useState(true);
    const [date, setDate] = useState<Date | undefined>(undefined);
    const [timeVal, setTimeVal] = useState("");

    // Payment
    const [paymentMethod, setPaymentMethod] = useState<"online" | "cash">("online");

    // Other
    const [cutleryCount, setCutleryCount] = useState(1);
    const [promoCode, setPromoCode] = useState("");
    const [wishes, setWishes] = useState("");

    // Check if current method is closed for today
    const DAYS = ['sunday', 'monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday'];
    const isTodayClosed = useMemo(() => {
        if (!serverTime || !workingHours) return false;
        const dayName = DAYS[serverTime.day_number];
        const currentTime = serverTime.time.substring(0, 5);
        const schedule = deliveryMethod === "delivery"
            ? workingHours.delivery?.[dayName]
            : workingHours.pickup?.[dayName];

        if (!schedule) return true;
        if (schedule.closed) return true;
        if (currentTime >= schedule.end) return true;
        return false;
    }, [serverTime, workingHours, deliveryMethod]);

    // Get working hours for current method and today
    const todaySchedule = useMemo(() => {
        if (!serverTime || !workingHours) return null;
        const dayName = DAYS[serverTime.day_number];
        return deliveryMethod === "delivery"
            ? workingHours.delivery?.[dayName]
            : workingHours.pickup?.[dayName];
    }, [serverTime, workingHours, deliveryMethod]);

    // Minimum selectable date (today if open, tomorrow if closed)
    const minDate = useMemo(() => {
        if (!serverTime) return new Date();
        const today = new Date(serverTime.datetime);
        if (isTodayClosed) {
            const tomorrow = new Date(today);
            tomorrow.setDate(tomorrow.getDate() + 1);
            return tomorrow;
        }
        return today;
    }, [serverTime, isTodayClosed]);

    // Minimum selectable time for selected date
    const minTime = useMemo(() => {
        if (!serverTime || !date) return "";
        const today = new Date(serverTime.datetime);
        const selectedDate = new Date(date);

        // If not today, no minimum time restriction (use schedule start)
        if (selectedDate.toDateString() !== today.toDateString()) {
            const dayNum = selectedDate.getDay();
            const dayName = DAYS[dayNum];
            const schedule = deliveryMethod === "delivery"
                ? workingHours?.delivery?.[dayName]
                : workingHours?.pickup?.[dayName];
            return schedule?.start || "";
        }

        // For today, minimum is current time + buffer (30 min) or schedule start, whichever is later
        const currentTime = serverTime.time.substring(0, 5);
        const [h, m] = currentTime.split(":").map(Number);
        const bufferedTime = new Date();
        bufferedTime.setHours(h, m + 30);
        const bufferedStr = `${String(bufferedTime.getHours()).padStart(2, "0")}:${String(bufferedTime.getMinutes()).padStart(2, "0")}`;

        return bufferedStr > (todaySchedule?.start || "") ? bufferedStr : (todaySchedule?.start || "");
    }, [serverTime, date, deliveryMethod, workingHours, todaySchedule]);

    // Max time based on schedule
    const maxTime = useMemo(() => {
        if (!date || !workingHours) return "";
        const selectedDate = new Date(date);
        const dayNum = selectedDate.getDay();
        const dayName = DAYS[dayNum];
        const schedule = deliveryMethod === "delivery"
            ? workingHours.delivery?.[dayName]
            : workingHours.pickup?.[dayName];
        return schedule?.end || "";
    }, [date, deliveryMethod, workingHours]);

    // Validate currently selected time
    const isTimeValid = useMemo(() => {
        if (isAsap) return true;
        if (!date || !timeVal) return false;
        if (minTime && timeVal < minTime) return false;
        if (maxTime && timeVal > maxTime) return false;
        return true;
    }, [isAsap, date, timeVal, minTime, maxTime]);

    // Disable ASAP if today is closed
    useEffect(() => {
        if (isTodayClosed && isAsap) {
            setIsAsap(false);
            setDate(minDate);
        }
    }, [isTodayClosed, isAsap, minDate]);

    // Regex
    const phoneRegex = /^(\+380\d{9}|0\d{9})$/;
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

    // Validation
    const isPhoneValid = !phone || phoneRegex.test(phone.trim());
    const isEmailValid = !email || emailRegex.test(email.trim());

    // Calculate Totals
    const deliveryPrice = deliveryMethod === "delivery"
        ? (deliveryResult?.zone?.price || 0)
        : 0;

    // Free delivery check
    const isFreeDelivery = deliveryMethod === "delivery" &&
        deliveryResult?.zone &&
        cartTotal >= deliveryResult.zone.freeOrderPrice;

    const finalDeliveryPrice = isFreeDelivery ? 0 : deliveryPrice;
    const toDoorPrice = (deliveryMethod === "delivery" && toDoor) ? doorPrice : 0;

    const total = cartTotal + finalDeliveryPrice + toDoorPrice;

    useEffect(() => {
        if (isInitialized && cart.length === 0) {
            router.push("/");
        }
    }, [cart, router, isInitialized]);

    const isValid = useMemo(() => {
        if (!name.trim()) return false;
        if (!phone.trim() || !phoneRegex.test(phone.trim())) return false;
        if (email.trim() && !emailRegex.test(email.trim())) return false;

        if (deliveryMethod === "delivery") {
            if (!deliveryResult?.address || !deliveryResult.zone) return false;
        }

        if (!isAsap) {
            if (!date || !timeVal) return false;
            if (!isTimeValid) return false;
        }

        return true;
    }, [name, phone, email, deliveryMethod, deliveryResult, isAsap, date, timeVal, isTimeValid]);

    const [isSubmitting, setIsSubmitting] = useState(false);

    const handleSubmit = async () => {
        if (!isValid) {
            toast.error("Будь ласка, заповніть всі обов'язкові поля та виправте помилки");
            return;
        }

        setIsSubmitting(true);

        try {
            const orderData = {
                name,
                phone,
                email: email || undefined,
                delivery_type_id: deliveryMethod,
                address: deliveryMethod === "delivery" ? deliveryResult?.address : "",
                coords: deliveryMethod === "delivery" && deliveryResult?.coords
                    ? `${deliveryResult.coords.lat},${deliveryResult.coords.lng}`
                    : undefined,
                entrance: deliveryMethod === "delivery" ? entrance : undefined,
                delivery_door: deliveryMethod === "delivery" ? toDoor : false,
                time: isAsap ? "ASAP" : `${format(date!, "yyyy-MM-dd")} ${timeVal}`,
                payment_method: paymentMethod,
                cutlery: cutleryCount,
                promo_code: promoCode || undefined,
                wishes: wishes || undefined,
                items: cart.map(item => {
                    // Get first selected variation ID if any
                    const variationIds = Object.values(item.selectedVariations);
                    const firstVariation = variationIds.length > 0 ? variationIds[0] : null;

                    return {
                        product_id: item.product.id,
                        product_variation_id: firstVariation?.id || undefined,
                        quantity: item.quantity
                    };
                }),
                client_total: total
            };

            const GATEWAY_URL = process.env.NEXT_PUBLIC_GATEWAY_URL || "http://localhost:8000";
            const response = await fetch(`${GATEWAY_URL}/orders`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(orderData)
            });

            const result = await response.json();

            if (!response.ok || !result.success) {
                toast.error(result.error || "Помилка при створенні замовлення");
                return;
            }

            toast.success("Замовлення оформлено успішно!");
            // Clear cart and redirect
            // TODO: clearCart()
            router.push("/menu");
        } catch (error) {
            console.error("Order error:", error);
            toast.error("Помилка з'єднання. Спробуйте ще раз.");
        } finally {
            setIsSubmitting(false);
        }
    };

    if (cart.length === 0) return null;

    const inputClasses = "w-full h-12 bg-white/5 border border-white/10 rounded-xl px-4 py-3 focus:outline-none focus:border-primary transition-colors text-white placeholder:text-gray-500";

    return (
        <div className="min-h-screen pt-24 pb-12 px-2 md:px-8 max-w-7xl mx-auto">
            <Link href="/menu" className="inline-flex items-center text-gray-400 hover:text-white mb-8 transition-colors cursor-pointer">
                <ArrowLeft className="w-4 h-4 mr-2" />
                Назад до меню
            </Link>

            <h1 className="text-2xl md:text-4xl font-bold mb-8">Оформлення замовлення</h1>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                {/* Left Column - Forms */}
                <div className="lg:col-span-2 space-y-8">

                    {/* Cart Summary (Top) */}
                    <section className="bg-white/5 rounded-2xl p-3 md:p-6 border border-white/10">
                        <h2 className="text-xl font-bold mb-4 flex items-center gap-2">
                            <Store className="w-5 h-5 text-primary" />
                            Ваше замовлення
                        </h2>
                        <div className="space-y-4 max-h-[300px] overflow-y-auto pr-2 custom-scrollbar">
                            {cart.map((item) => (
                                <div key={item.id} className="flex gap-2 md:gap-4 items-center bg-white/5 p-2 md:p-3 rounded-xl">
                                    {/* Delete Button - Left */}
                                    <button
                                        onClick={() => removeFromCart(item.id)}
                                        className="p-2 text-gray-400 hover:text-red-400 hover:bg-red-400/10 rounded-full transition-all cursor-pointer shrink-0"
                                    >
                                        <Trash2 className="w-5 h-5" />
                                    </button>

                                    <div className="w-16 h-16 rounded-lg overflow-hidden shrink-0">
                                        <img
                                            src={getImageUrl(item.product.image)}
                                            alt={item.product.name}
                                            className="w-full h-full object-cover"
                                        />
                                    </div>
                                    <div className="flex-1 min-w-0">
                                        <h4 className="font-medium truncate">{item.product.name}</h4>
                                        {Object.values(item.selectedVariations).length > 0 && (
                                            <p className="text-xs text-gray-400 truncate">
                                                {Object.values(item.selectedVariations).map(v => v.name).join(", ")}
                                            </p>
                                        )}
                                        <div className="mt-1">
                                            <span className="font-bold text-primary">{item.product.price * item.quantity} ₴</span>
                                        </div>
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
                            ))}
                        </div>
                    </section>

                    {/* Delivery Method */}
                    <section className="bg-white/5 rounded-2xl p-3 md:p-6 border border-white/10 space-y-6">
                        <h2 className="text-xl font-bold flex items-center gap-2">
                            <Truck className="w-5 h-5 text-primary" />
                            Спосіб отримання
                        </h2>

                        <div className="grid grid-cols-2 gap-4">
                            <button
                                onClick={() => setDeliveryMethod("delivery")}
                                className={`p-4 rounded-xl border-2 flex flex-col items-center gap-2 transition-all cursor-pointer ${deliveryMethod === "delivery"
                                    ? "border-primary bg-primary/10 text-white"
                                    : "border-white/10 hover:border-white/20 text-gray-400"
                                    }`}
                            >
                                <Truck className="w-6 h-6" />
                                <span className="font-medium">Доставка</span>
                            </button>
                            <button
                                onClick={() => setDeliveryMethod("pickup")}
                                className={`p-4 rounded-xl border-2 flex flex-col items-center gap-2 transition-all cursor-pointer ${deliveryMethod === "pickup"
                                    ? "border-primary bg-primary/10 text-white"
                                    : "border-white/10 hover:border-white/20 text-gray-400"
                                    }`}
                            >
                                <Store className="w-6 h-6" />
                                <span className="font-medium">Самовивіз</span>
                            </button>
                        </div>

                        {deliveryMethod === "delivery" && (
                            <div className="space-y-6 animate-in fade-in slide-in-from-top-4">
                                <DeliveryMap onLocationSelect={setDeliveryResult} cartTotal={cartTotal}>
                                    {deliveryResult?.address && (
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

                    {/* Contact Info */}
                    <section className="bg-white/5 rounded-2xl p-3 md:p-6 border border-white/10 space-y-6">
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

                        {/* Additional Fields */}
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

                        {/* Time Selection */}
                        <div className="space-y-2">
                            <label className="text-sm font-medium text-gray-300 block mb-1">На коли</label>

                            {/* Closed today banner */}
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
                                    className={`py-3 px-4 rounded-xl border transition-all flex items-center justify-center gap-2 cursor-pointer h-12 ${isAsap && !isTodayClosed
                                        ? "bg-primary/20 border-primary text-primary font-medium"
                                        : "bg-white/5 border-white/10 text-gray-400 hover:bg-white/10"
                                        } ${isTodayClosed ? "opacity-50 cursor-not-allowed" : ""}`}
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
                                            "min-w-0", // prevent overflow
                                            "[&::-webkit-calendar-picker-indicator]:hidden [&::-webkit-calendar-picker-indicator]:appearance-none" // Hide default icon
                                        )}
                                    />
                                </div>
                            </div>

                            {/* Time validation error */}
                            {!isTimeValid && timeVal && (
                                <p className="text-red-500 text-xs mt-1">
                                    Виберіть час в межах {minTime} - {maxTime}
                                </p>
                            )}
                        </div>

                        <div className="space-y-2">
                            <label className="text-sm font-medium text-gray-300 block mb-1">Промокод</label>
                            <input
                                type="text"
                                value={promoCode}
                                onChange={(e) => setPromoCode(e.target.value)}
                                className={inputClasses}
                                placeholder="Введіть промокод"
                            />
                        </div>

                        <div className="space-y-2">
                            <label className="text-sm font-medium text-gray-300 block mb-1">Ваші побажання до замовлення</label>
                            <textarea
                                value={wishes}
                                onChange={(e) => setWishes(e.target.value)}
                                className={cn(inputClasses, "min-h-[100px] resize-none h-auto")}
                                placeholder="Наприклад: не дзвонити у двері..."
                            />
                        </div>
                    </section>
                </div>

                {/* Right Column - Summary */}
                <div className="lg:col-span-1">
                    <div className="sticky top-24 space-y-6">

                        {/* Payment Method */}
                        <section className="bg-white/5 rounded-2xl p-3 md:p-6 border border-white/10">
                            <h2 className="text-xl font-bold mb-4">Оплата</h2>
                            <div className="space-y-3">
                                <label className={`flex items-center gap-4 p-4 rounded-xl border cursor-pointer transition-all ${paymentMethod === "online"
                                    ? "bg-primary/10 border-primary"
                                    : "bg-white/5 border-white/10 hover:border-white/20"
                                    }`}>
                                    <input
                                        type="radio"
                                        name="payment"
                                        value="online"
                                        checked={paymentMethod === "online"}
                                        onChange={() => setPaymentMethod("online")}
                                        className="sr-only"
                                    />
                                    <CreditCard className={`w-6 h-6 ${paymentMethod === "online" ? "text-primary" : "text-gray-400"}`} />
                                    <span className={paymentMethod === "online" ? "font-bold text-white" : "text-gray-300"}>
                                        Безготівкова на сайті
                                    </span>
                                </label>

                                {deliveryMethod === "pickup" && (
                                    <label className={`flex items-center gap-4 p-4 rounded-xl border cursor-pointer transition-all ${paymentMethod === "cash"
                                        ? "bg-primary/10 border-primary"
                                        : "bg-white/5 border-white/10 hover:border-white/20"
                                        }`}>
                                        <input
                                            type="radio"
                                            name="payment"
                                            value="cash"
                                            checked={paymentMethod === "cash"}
                                            onChange={() => setPaymentMethod("cash")}
                                            className="sr-only"
                                        />
                                        <Banknote className={`w-6 h-6 ${paymentMethod === "cash" ? "text-primary" : "text-gray-400"}`} />
                                        <span className={paymentMethod === "cash" ? "font-bold text-white" : "text-gray-300"}>
                                            Готівкою при отриманні
                                        </span>
                                    </label>
                                )}
                            </div>
                        </section>

                        {/* Order Summary */}
                        <section className="bg-white/5 rounded-2xl p-3 md:p-6 border border-white/10 space-y-4">
                            <h2 className="text-xl font-bold">Разом</h2>

                            {/* Free Delivery Notification Removed from here */}

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
                    </div>
                </div>
            </div>
        </div>
    );
}
