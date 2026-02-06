"use client";

import { useEffect, useState, useMemo } from "react";
import { useCart } from "@/contexts/cart-context";
import { toast } from "sonner";
import { ArrowLeft } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { format } from "date-fns";
import { useShopStatus } from "@/hooks/use-shop-status";
import { useSettings } from "@/contexts/settings-context";
import { SearchResult } from "@/types/delivery";
import { cn } from "@/lib/utils";

// Components
import { OrderItemsList } from "./components/order-items-list";
import { DeliveryMethodSection } from "./components/delivery-method-section";
import { ContactInfoSection } from "./components/contact-info-section";
import { TimeSelectionSection } from "./components/time-selection-section";
import { PaymentSection } from "./components/payment-section";
import { OrderSummary } from "./components/order-summary";

const DAYS = ['sunday', 'monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday'];

export default function CheckoutPage() {
    const { cart, cartTotal, removeFromCart, updateQuantity, isInitialized } = useCart();
    const router = useRouter();
    const { isPaused, workingHours, serverTime } = useShopStatus();
    const { getSetting } = useSettings();
    const [doorPrice, setDoorPrice] = useState(50);

    // Form State
    const [deliveryMethod, setDeliveryMethod] = useState<"delivery" | "pickup">("delivery");
    const [deliveryResult, setDeliveryResult] = useState<SearchResult | null>(null);
    const [entrance, setEntrance] = useState("");
    const [toDoor, setToDoor] = useState(false);
    const [name, setName] = useState("");
    const [phone, setPhone] = useState("");
    const [email, setEmail] = useState("");
    const [isAsap, setIsAsap] = useState(true);
    const [date, setDate] = useState<Date | undefined>(undefined);
    const [timeVal, setTimeVal] = useState("");
    const [paymentMethod, setPaymentMethod] = useState<"bank" | "cash">("bank");
    const [cutleryCount, setCutleryCount] = useState(1);
    const [promoCode, setPromoCode] = useState("");
    const [wishes, setWishes] = useState("");
    const [isSubmitting, setIsSubmitting] = useState(false);

    // 1. Fetch delivery door price
    useEffect(() => {
        const priceSetting = getSetting("delivery_door_price");
        if (priceSetting?.value) {
            const parsed = parseFloat(priceSetting.value);
            if (!isNaN(parsed)) setDoorPrice(parsed);
        }
    }, [getSetting]);

    // 2. Redirect if sales are paused or cart is empty
    useEffect(() => {
        if (isPaused) {
            toast.error("Вибачте, ми тимчасово не приймаємо замовлення.");
            router.push("/menu");
        }
    }, [isPaused, router]);

    useEffect(() => {
        if (isInitialized && cart.length === 0) {
            router.push("/");
        }
    }, [cart, router, isInitialized]);

    // 3. Scheduling Logic
    const isTodayClosed = useMemo(() => {
        if (!serverTime || !workingHours) return false;
        const dayName = DAYS[serverTime.day_number];
        const schedule = deliveryMethod === "delivery"
            ? workingHours.delivery?.[dayName]
            : workingHours.pickup?.[dayName];

        if (!schedule || schedule.closed) return true;
        if (serverTime.time.substring(0, 5) >= schedule.end) return true;
        return false;
    }, [serverTime, workingHours, deliveryMethod]);

    const todaySchedule = useMemo(() => {
        if (!serverTime || !workingHours) return null;
        const dayName = DAYS[serverTime.day_number];
        return deliveryMethod === "delivery"
            ? workingHours.delivery?.[dayName]
            : workingHours.pickup?.[dayName];
    }, [serverTime, workingHours, deliveryMethod]);

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

    const minTime = useMemo(() => {
        if (!serverTime || !date) return "";
        const today = new Date(serverTime.datetime);
        const selectedDate = new Date(date);

        if (selectedDate.toDateString() !== today.toDateString()) {
            const dayName = DAYS[selectedDate.getDay()];
            const schedule = deliveryMethod === "delivery"
                ? workingHours?.delivery?.[dayName]
                : workingHours?.pickup?.[dayName];
            return schedule?.start || "";
        }

        const [h, m] = serverTime.time.substring(0, 5).split(":").map(Number);
        const bufferedTime = new Date();
        bufferedTime.setHours(h, m + 30);
        const bufferedStr = `${String(bufferedTime.getHours()).padStart(2, "0")}:${String(bufferedTime.getMinutes()).padStart(2, "0")}`;

        return bufferedStr > (todaySchedule?.start || "") ? bufferedStr : (todaySchedule?.start || "");
    }, [serverTime, date, deliveryMethod, workingHours, todaySchedule]);

    const maxTime = useMemo(() => {
        if (!date || !workingHours) return "";
        const dayName = DAYS[new Date(date).getDay()];
        const schedule = deliveryMethod === "delivery"
            ? workingHours.delivery?.[dayName]
            : workingHours.pickup?.[dayName];
        return schedule?.end || "";
    }, [date, deliveryMethod, workingHours]);

    const isTimeValid = useMemo(() => {
        if (isAsap) return true;
        if (!date || !timeVal) return false;
        if (minTime && timeVal < minTime) return false;
        if (maxTime && timeVal > maxTime) return false;
        return true;
    }, [isAsap, date, timeVal, minTime, maxTime]);

    useEffect(() => {
        if (isTodayClosed && isAsap) {
            setIsAsap(false);
            setDate(minDate);
        }
    }, [isTodayClosed, isAsap, minDate]);

    // 4. Validation
    const phoneRegex = /^(\+380\d{9}|0\d{9})$/;
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    const isPhoneValid = !phone || phoneRegex.test(phone.trim());
    const isEmailValid = !email || emailRegex.test(email.trim());

    const isValid = useMemo(() => {
        if (!name.trim() || !phone.trim() || !phoneRegex.test(phone.trim())) return false;
        if (email.trim() && !emailRegex.test(email.trim())) return false;
        if (deliveryMethod === "delivery" && (!deliveryResult?.address || !deliveryResult.zone)) return false;
        if (!isAsap && (!date || !timeVal || !isTimeValid)) return false;
        return true;
    }, [name, phone, email, deliveryMethod, deliveryResult, isAsap, date, timeVal, isTimeValid]);

    // 5. Totals
    const deliveryPrice = deliveryMethod === "delivery" ? (deliveryResult?.zone?.price || 0) : 0;
    const isFreeDelivery = !!(deliveryMethod === "delivery" && deliveryResult?.zone && cartTotal >= deliveryResult.zone.freeOrderPrice);
    const finalDeliveryPrice = isFreeDelivery ? 0 : deliveryPrice;
    const toDoorPrice = (deliveryMethod === "delivery" && toDoor) ? doorPrice : 0;
    const total = cartTotal + finalDeliveryPrice + toDoorPrice;

    // 6. Handle Submit
    const handleSubmit = async () => {
        if (!isValid) {
            toast.error("Будь ласка, заповніть всі обов'язкові поля та виправте помилки");
            return;
        }

        setIsSubmitting(true);
        try {
            const orderData = {
                name, phone, email: email || undefined,
                delivery_type_id: deliveryMethod,
                address: deliveryMethod === "delivery" ? deliveryResult?.address : "",
                coords: (deliveryMethod === "delivery" && deliveryResult?.coords)
                    ? `${deliveryResult.coords.lat},${deliveryResult.coords.lng}` : undefined,
                entrance: deliveryMethod === "delivery" ? entrance : undefined,
                delivery_door: deliveryMethod === "delivery" ? toDoor : false,
                time: isAsap ? "ASAP" : `${format(date!, "yyyy-MM-dd")} ${timeVal}`,
                payment_method: paymentMethod,
                cutlery: cutleryCount,
                promo_code: promoCode || undefined,
                wishes: wishes || undefined,
                items: cart.map(item => ({
                    product_id: item.product.id,
                    product_variation_id: Object.values(item.selectedVariations)[0]?.id || undefined,
                    quantity: item.quantity
                })),
                client_total: total
            };

            const response = await fetch(`${process.env.NEXT_PUBLIC_GATEWAY_URL}/orders`, {
                method: "POST", headers: { "Content-Type": "application/json" },
                body: JSON.stringify(orderData)
            });

            const result = await response.json();
            if (!response.ok || !result.success) {
                toast.error(result.error || "Помилка при створенні замовлення");
                return;
            }

            toast.success("Замовлення оформлено успішно!");
            if (result.data?.payment_url) {
                window.location.href = result.data.payment_url;
                return;
            }
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
        <div className="min-h-screen pt-18 pb-12 px-2 md:px-8 max-w-7xl mx-auto">
            <Link href="/menu" className="inline-flex items-center text-gray-400 hover:text-white mb-8 transition-colors cursor-pointer">
                <ArrowLeft className="w-4 h-4 mr-2" />
                Назад до меню
            </Link>

            <h1 className="text-2xl md:text-4xl font-bold mb-8">Оформлення замовлення</h1>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                <div className="lg:col-span-2 space-y-8">
                    <OrderItemsList
                        items={cart}
                        cartTotal={cartTotal}
                        removeFromCart={removeFromCart}
                        updateQuantity={updateQuantity}
                    />

                    <DeliveryMethodSection
                        deliveryMethod={deliveryMethod}
                        setDeliveryMethod={setDeliveryMethod}
                        setDeliveryResult={setDeliveryResult}
                        cartTotal={cartTotal}
                        deliveryResult={deliveryResult}
                        entrance={entrance}
                        setEntrance={setEntrance}
                        toDoor={toDoor}
                        setToDoor={setToDoor}
                        doorPrice={doorPrice}
                    />

                    <section className="bg-white/5 rounded-2xl p-3 md:p-6 border border-white/10 space-y-6">
                        <ContactInfoSection
                            name={name} setName={setName}
                            phone={phone} setPhone={setPhone}
                            email={email} setEmail={setEmail}
                            cutleryCount={cutleryCount} setCutleryCount={setCutleryCount}
                            isPhoneValid={isPhoneValid} isEmailValid={isEmailValid}
                        />

                        <TimeSelectionSection
                            deliveryMethod={deliveryMethod}
                            isTodayClosed={isTodayClosed}
                            isAsap={isAsap} setIsAsap={setIsAsap}
                            date={date} setDate={setDate}
                            timeVal={timeVal} setTimeVal={setTimeVal}
                            minDate={minDate} minTime={minTime} maxTime={maxTime}
                            isTimeValid={isTimeValid}
                        />

                        <div className="space-y-6 pt-2 border-t border-white/5">
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
                        </div>
                    </section>
                </div>

                <div className="lg:col-span-1">
                    <div className="sticky top-24 space-y-6">
                        <PaymentSection
                            paymentMethod={paymentMethod}
                            setPaymentMethod={setPaymentMethod}
                            deliveryMethod={deliveryMethod}
                        />

                        <OrderSummary
                            cartTotal={cartTotal}
                            deliveryMethod={deliveryMethod}
                            isFreeDelivery={isFreeDelivery}
                            deliveryPrice={deliveryPrice}
                            toDoor={toDoor}
                            toDoorPrice={toDoorPrice}
                            total={total}
                            isValid={isValid && !isSubmitting}
                            handleSubmit={handleSubmit}
                        />
                    </div>
                </div>
            </div>
        </div>
    );
}
