"use client"

import { Menu, Package, X, Home, UtensilsCrossed, Star } from "lucide-react";
import { usePathname, useRouter } from "next/navigation";
import Link from "next/link";
import Image from "next/image";
import Logo from "@/resources/images/logo.webp";
import { Button } from "@/components/ui/button";
import { useState } from "react";
import { useCart } from "@/contexts/cart-context";
import CartDrawer from "@/components/cart-drawer";

export default function Header() {
    const pathname = usePathname();
    const [isMenuOpen, setIsMenuOpen] = useState(false);
    const router = useRouter();
    const { cartTotal, setIsCartOpen } = useCart();

    const isMobileMenuPage = pathname === "/menu";

    const navigateTo = (path: string) => {
        setIsMenuOpen(false);
        router.push(path);
    };

    const handleCartClick = () => {
        setIsCartOpen(true);
    };

    return (
        <div className="sticky top-0 z-50 w-full">
            <header
                className="w-full bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 pt-2 pb-1">
                <div className="container mx-auto flex items-center justify-between px-4 h-full pb-2">
                    <div className="flex-1 flex justify-start">
                        <Link className="flex flex-row gap-2" href="/">
                            <Image
                                src={Logo}
                                alt="Logo"
                                width={60}
                                height={60}
                                className="block h-auto w-[60px]"
                                priority
                            />
                            <div className="flex flex-col align-middle justify-center">
                                <span className="font-bold">
                                    BROBAR
                                </span>
                                <span className="text-[12px] font-bold text-gray-400">
                                    Ресторан
                                </span>
                            </div>
                        </Link>
                    </div>

                    <div className="flex-1 flex justify-center">
                    </div>

                    <div className="flex-1 flex justify-end items-center space-x-6">
                        <Button
                            size="lg"
                            className="cart-button hidden md:flex items-center cursor-pointer"
                            onClick={handleCartClick}
                        >
                            <Package className="size-6" />
                            <div>
                                Кошик - <span>{cartTotal}</span> ₴
                            </div>
                        </Button>

                        {isMobileMenuPage && (
                            <div onClick={() => setIsMenuOpen(true)} className="cursor-pointer md:hidden">
                                <Menu size={40} />
                            </div>
                        )}
                    </div>
                </div>

                <div className={`container mx-auto px-4 gap-5 ${isMobileMenuPage ? "hidden md:flex" : "flex"}`}>
                    <Link
                        href="/"
                        className={`${pathname === "/" ? "border-b-2 border-brand" : "border-b-2 border-transparent"}`}
                    >
                        Головна
                    </Link>
                    <Link
                        href="/menu"
                        className={`${pathname === "/menu" ? "border-b-2 border-brand" : "border-b-2 border-transparent"}`}
                    >
                        Меню
                    </Link>
                    <Link
                        href="/reviews"
                        className={`${pathname === "/reviews" ? "border-b-2 border-brand" : "border-b-2 border-transparent"}`}
                    >
                        Залишити відгук
                    </Link>
                </div>
            </header>

            {isMobileMenuPage && (
                <div
                    className="md:hidden w-full bg-primary text-white py-2 px-4 flex justify-between items-center cursor-pointer"
                    onClick={handleCartClick}
                >
                    <div className="flex flex-row gap-2 items-center">
                        <Package className="size-[32px] text-black" />
                        <span className="font-medium text-black text-sm">Кошик</span>
                    </div>
                    <span className="font-medium text-black text-sm">{cartTotal} ₴</span>
                </div>
            )}

            {/* Cart Drawer */}
            <CartDrawer />

            {/* Mobile Menu Drawer */}
            <div className={`fixed inset-0 z-[60] flex justify-end transition-all duration-150 ${isMenuOpen ? "visible" : "invisible"}`}>
                {/* Backdrop */}
                <div
                    className={`absolute inset-0 bg-black/50 backdrop-blur-sm transition-opacity duration-150 ${isMenuOpen ? "opacity-100" : "opacity-0"}`}
                    onClick={() => setIsMenuOpen(false)}
                />

                {/* Drawer Content */}
                <div className={`relative w-[300px] h-full bg-background border-l border-border p-6 shadow-xl flex flex-col gap-6 transition-transform duration-150 will-change-transform ${isMenuOpen ? "translate-x-0" : "translate-x-full"}`}>
                    <div className="flex justify-between items-center">
                        <span className="text-xl font-bold">Меню</span>
                        <div onClick={() => setIsMenuOpen(false)} className="cursor-pointer hover:opacity-70">
                            <X size={32} />
                        </div>
                    </div>

                    <nav className="flex flex-col gap-1 flex-1">
                        <div
                            onClick={() => navigateTo("/")}
                            className={`flex items-center gap-3 text-lg py-2 rounded-md cursor-pointer`}
                        >
                            <Home size={24} />
                            <span>Головна</span>
                        </div>
                        <div
                            onClick={() => navigateTo("/menu")}
                            className={`flex items-center gap-3 text-lg py-2 rounded-md cursor-pointer`}
                        >
                            <UtensilsCrossed size={24} />
                            <span>Меню</span>
                        </div>
                        <div
                            onClick={() => navigateTo("/reviews")}
                            className={`flex items-center gap-3 text-lg py-2 rounded-md cursor-pointer`}
                        >
                            <Star size={24} />
                            <span>Відгуки</span>
                        </div>

                        <div className="mt-auto flex flex-col gap-2 pt-4 border-t border-border/50">
                            <div
                                onClick={() => navigateTo("/privacy")}
                                className="text-sm text-gray-400 hover:text-white cursor-pointer py-1"
                            >
                                Політика конфіденційності
                            </div>
                            <div
                                onClick={() => navigateTo("/contract")}
                                className="text-sm text-gray-400 hover:text-white cursor-pointer py-1"
                            >
                                Договір публічної оферти
                            </div>
                        </div>
                    </nav>
                </div>
            </div>
        </div>
    )
}
