"use client";

import Image from "next/image";
import Link from "next/link";
import MainImage from "@/resources/images/main.webp";
import { Clock2, Instagram, Map, Phone } from "lucide-react";
import { Button } from "@/components/ui/button";
import WelcomeMap from "@/components/welcome-map";
import Header from "@/components/header";
import Footer from "@/components/footer";
import { useShopStatus } from "@/hooks/use-shop-status";
import { formatWorkingHours } from "@/lib/working-hours";
import { useDeliveryZones } from "@/hooks/use-delivery-zones";

export default function Welcome() {
    const { workingHours } = useShopStatus();
    const deliveryHours = formatWorkingHours(workingHours, 'delivery');
    const pickupHours = formatWorkingHours(workingHours, 'pickup');
    const { zones } = useDeliveryZones();

    return (
        <div>
            <Header />

            <main>
                {/* Hero */}
                <div className="relative">
                    <Image
                        src={MainImage}
                        fill={true}
                        alt="Main background image"
                        className="object-cover object-center"
                        priority
                    />
                    <div className="flex flex-col items-center text-white py-24 md:py-36 relative">
                        <div className="flex flex-col items-center gap-2 drop-shadow-xl backdrop-blur-md bg-black/20 px-8 py-8 rounded-2xl">
                            <div className="flex flex-col items-center">
                                <span className="banner-text-lg font-extrabold w-fit mx-auto">
                                    Гарячі страви
                                </span>
                            </div>
                            <span className="banner-text-sm font-semibold w-fit mx-auto">Від Бро для Бро</span>
                            <span className="banner-text-md font-semibold text-brand w-fit mx-auto">Безкоштовна доставка</span>
                            <span className="banner-text-sm font-semibold text-brand mb-6 w-fit mx-auto">
                                при мінімальному замовленні
                            </span>
                            <Link href="/menu">
                                <Button size="lg" className="cursor-pointer background-dark text-white text-xl font-bold px-10 py-6 shadow-lg shadow-black/50 hover:scale-105 hover:shadow-xl hover:shadow-black/60 transition-all duration-200">
                                    Заглянути у меню 🍔
                                </Button>
                            </Link>
                        </div>
                    </div>
                </div>

                {/* Часи роботи */}
                <section className="background-brand py-12 flex flex-col items-center text-center text-dark">
                    <Clock2 size={30} className="mb-3" />
                    <div className="flex flex-col gap-4">
                        <span className="font-extrabold text-xl">Часи роботи</span>

                        <div className="flex flex-col gap-1">
                            <span className="font-bold text-lg">Доставка</span>
                            {deliveryHours.length > 0 ? (
                                deliveryHours.map((line, i) => (
                                    <span key={i} className="font-light text-sm">{line}</span>
                                ))
                            ) : (
                                <span className="font-light text-sm">Завантаження...</span>
                            )}
                        </div>

                        <div className="flex flex-col gap-1">
                            <span className="font-bold text-lg">Самовивіз</span>
                            {pickupHours.length > 0 ? (
                                pickupHours.map((line, i) => (
                                    <span key={i} className="font-light text-sm">{line}</span>
                                ))
                            ) : (
                                <span className="font-light text-sm">Завантаження...</span>
                            )}
                        </div>
                    </div>
                </section>

                {/* Доставка */}
                <section className="background-dark py-12 flex flex-col items-center text-white relative overflow-hidden px-4 md:px-8">
                    <h2 className="text-3xl md:text-5xl font-bold mb-3 text-center">Доставка і оплата</h2>
                    <p className="text-white text-lg md:text-xl font-bold mb-12 text-center">Бро, роби добро!</p>

                    <div className="max-w-7xl mx-auto grid grid-cols-1 lg:grid-cols-2 gap-12 lg:gap-24 text-left w-full">
                        {/* Left Column: Description */}
                        <div className="flex flex-col gap-6 text-[15px] md:text-[17px] leading-relaxed text-gray-300 font-normal">
                            <p>
                                <span className="text-brand font-bold">Brobar.delivery</span> здійснює швидку і якісну доставку страв. Ми розробили спеціальну box-упаковку, завдяки якій сет приїде до тебе у зручному форматі та ти зможеш організувати будь-яке свято у себе вдома або в офісі. Ми доставляємо замовлення через сервіс таксі &quot;On Taxi&quot; за наш рахунок.
                            </p>
                            <p>
                                Наші страви підійдуть для перекусу в офісі або вдома, для святкування важливої події, романтичного вечора або дружніх вечірок з компанією твоїх Бро.
                            </p>
                            <p>
                                Оплатити замовлення можна банківською картою онлайн через сервіс <span className="text-brand font-bold">Monobank</span> у вікні оформлення замовлення. Оплата готівкою тільки у барі при самовивозі.
                            </p>
                        </div>

                        {/* Right Column: Pricing & Zones */}
                        <div className="flex flex-col gap-8">
                            <div>
                                <h3 className="text-xl md:text-2xl font-normal mb-6">Вартість доставки:</h3>

                                <div className="space-y-6 text-[15px] md:text-[17px] text-gray-300">
                                    {zones.map((zone, index) => (
                                        <div key={index}>
                                            <h4 className="text-white text-lg underline decoration-1 underline-offset-4 mb-1 font-medium">{zone.name}</h4>
                                            <ul className="space-y-1">
                                                <li>• замовлення від {zone.freeOrderPrice}₴ доставляємо БЕЗКОШТОВНО</li>
                                                <li>• до {zone.freeOrderPrice}₴ вартість {zone.price}₴</li>
                                            </ul>
                                        </div>
                                    ))}
                                </div>
                            </div>

                            <div className="flex flex-col gap-1">
                                <h4 className="text-white text-lg font-bold">Час доставки</h4>
                                <p className="text-xl font-normal text-gray-300">30-120 хвилин</p>
                            </div>

                            <p className="text-lg md:text-xl font-medium mt-2">
                                Бро, дивись на <Link href="/delivery" className="underline decoration-1 underline-offset-4 cursor-pointer hover:text-brand transition-colors">карті</Link> в якій ти зоні доставки
                            </p>
                        </div>
                    </div>

                    <Link href="/delivery" className="self-center mt-12">
                        <Button size="lg" className="w-auto text-dark text-lg font-bold hover:scale-105 transition-transform background-brand px-8 cursor-pointer">
                            <Map className="size-5 mr-2" />
                            Мапа доставки
                        </Button>
                    </Link>
                </section>

                {/* Контакти */}
                <section className="background-brand py-12 flex flex-col items-center text-center text-dark">
                    <span className="font-extrabold text-4xl mb-2">Наші контакти</span>
                    <span className="font-light text-sm">м. Харків</span>
                    <span className="font-light text-sm mb-2">вул. Григорія Сковороди 64 (вхід з вул. Багалія)</span>

                    <div className="flex flex-col gap-1">
                        <Link href="tel:+380635009597">
                            <Button
                                size="lg"
                                className="self-center w-auto text-sm font-semibold background-dark text-white cursor-pointer hover:scale-105 transition-transform"
                            >
                                <Phone />
                                +38-(063)-500-95-97
                            </Button>
                        </Link>

                        <Link href="https://instagram.com/brobar_kh" target="_blank" rel="noopener noreferrer">
                            <Button
                                size="lg"
                                className="self-center w-auto text-sm font-semibold background-dark text-white cursor-pointer hover:scale-105 transition-transform"
                            >
                                <Instagram />
                                brobar_kh
                            </Button>
                        </Link>
                    </div>
                </section>

                <WelcomeMap />

            </main>
            <Footer />
        </div>
    );
}
