"use client";

import Header from "@/components/header";
import Footer from "@/components/footer";
import DeliveryMap from "@/components/delivery-map";
import { useDeliveryZones } from "@/hooks/use-delivery-zones";

export default function DeliveryPage() {
    const { zones } = useDeliveryZones();

    return (
        <div className="min-h-screen bg-black text-white selection:bg-primary/30">
            <Header />

            <main className="container mx-auto px-4 pt-16 pb-20">
                <h1 className="text-4xl md:text-5xl font-bold mb-12 text-center">Вартість доставки та оплата</h1>

                <div className="max-w-4xl mx-auto space-y-12">

                    {/* Zones Grid */}
                    <section>
                        <h2 className="text-2xl font-bold mb-6">Вартість та час доставки:</h2>
                        <div className="grid md:grid-cols-2 gap-8">
                            {zones.map((zone, index) => (
                                <div key={index} className="p-6 rounded-2xl bg-white/5 border border-white/10 transition-colors group" style={{ borderColor: `${zone.color}40` }}>
                                    <h3 className="text-xl font-bold mb-4" style={{ color: zone.color }}>{zone.name}</h3>
                                    <ul className="space-y-2 text-gray-300">
                                        <li className="flex items-start">
                                            <span className="mr-2">•</span>
                                            <span>замовлення від {zone.freeOrderPrice}₴ доставляємо <span className="font-bold text-white">БЕЗКОШТОВНО</span></span>
                                        </li>
                                        <li className="flex items-start">
                                            <span className="mr-2">•</span>
                                            <span>до {zone.freeOrderPrice}₴ вартість {zone.price}₴</span>
                                        </li>
                                    </ul>
                                </div>
                            ))}
                        </div>
                    </section>

                    {/* Additional Info Grid */}
                    <div className="grid md:grid-cols-2 gap-8">
                        {/* Time */}
                        <section className="space-y-4">
                            <h2 className="text-2xl font-bold border-b border-white/10 pb-2 inline-block">Час доставки</h2>
                            <p className="text-4xl font-bold text-primary">30-120 <span className="text-2xl font-normal text-gray-400">хвилин</span></p>
                        </section>

                        {/* Payment */}
                        <section className="space-y-4">
                            <h2 className="text-2xl font-bold border-b border-white/10 pb-2 inline-block">Оплата:</h2>
                            <p className="text-gray-300 leading-relaxed">
                                Оплатити замовлення можна банківською картою онлайн через сервіс Monobank у вікні оформлення замовлення.
                                Оплата готівкою тільки у барі при самовивозі.
                            </p>
                        </section>
                    </div>

                    {/* Map Section */}
                    <section className="space-y-6 pt-8">
                        <h2 className="text-3xl font-bold text-center">Перевірити адресу</h2>
                        <DeliveryMap />
                    </section>
                </div>
            </main>

            <Footer />
        </div>
    );
}
