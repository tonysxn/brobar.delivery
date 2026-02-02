"use client"

import {
    Package,
    MapPin, Phone, Clock2,
    LucideIcon,
    Plus,
    Minus,
    X,
    icons
} from "lucide-react";
import { useState, useRef, useEffect } from "react";
import Header from "@/components/header";
import Footer from "@/components/footer";
import { Button } from "@/components/ui/button";
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from "@/components/ui/accordion";
import { useCart, Product, Variation } from "@/contexts/cart-context";
import { toast } from "sonner";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Skeleton } from "@/components/ui/skeleton";
import { useShopStatus } from "@/hooks/use-shop-status";
import { formatWorkingHours } from "@/lib/working-hours";
import { ShopStatusAlert } from "@/components/shop-status-alert";
import { ProductCard } from "@/components/menu/product-card";
import { ImageWithFallback } from "@/components/ui/image-with-fallback";

// Helper to convert icon name to PascalCase (e.g., "utensils" -> "Utensils", "concierge-bell" -> "ConciergeBell")
const toPascalCase = (str: string): string => {
    return str
        .split(/[-_]/)
        .map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
        .join('');
};

// Dynamic icon getter - any lucide icon name will work
const getIconByName = (iconName: string | undefined | null): LucideIcon => {
    if (!iconName) return Package;
    const pascalName = toPascalCase(iconName);
    return (icons as Record<string, LucideIcon>)[pascalName] || Package;
};

interface Category {
    id: string;
    name: string;
    slug: string;
    icon: string;
    sort: number;
    products: Product[];
}

const GATEWAY_URL = process.env.NEXT_PUBLIC_GATEWAY_URL;

export default function Menu() {
    const [expanded, setExpanded] = useState(false);
    const [deliveryType, setDeliveryType] = useState<"delivery" | "pickup">("delivery");
    const contentRef = useRef<HTMLDivElement>(null);
    const [height, setHeight] = useState("0px");
    const [isSticky, setIsSticky] = useState(false);

    // API data state
    const [categories, setCategories] = useState<Category[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    // Cart from context
    const { addToCart } = useCart();

    // ... modal state ...
    const [variationModalOpen, setVariationModalOpen] = useState(false);
    const [selectedProduct, setSelectedProduct] = useState<Product | null>(null);
    const [tempSelectedVariations, setTempSelectedVariations] = useState<Record<string, Variation>>({});
    const [quantity, setQuantity] = useState(1);

    // Shop status
    const { isOpen, message, workingHours, deliveryOpen, pickupOpen } = useShopStatus();
    const deliveryHours = formatWorkingHours(workingHours, 'delivery');
    const pickupHours = formatWorkingHours(workingHours, 'pickup');
    const hoursAreSame = JSON.stringify(deliveryHours) === JSON.stringify(pickupHours);
    const displayHours = hoursAreSame ? deliveryHours : null;

    // Fetch menu data from API
    useEffect(() => {
        const fetchMenu = async () => {
            try {
                setLoading(true);
                const response = await fetch(`${GATEWAY_URL}/menu`);
                const data = await response.json();

                if (data.success) {
                    setCategories(data.data);
                } else {
                    setError("Failed to load menu");
                }
            } catch (err) {
                console.error("Failed to fetch menu:", err);
                setError("Failed to load menu");
            } finally {
                setLoading(false);
            }
        };

        fetchMenu();
    }, []);

    useEffect(() => {
        const handleScroll = () => {
            const nav = document.getElementById("categories-nav");
            if (nav) {
                const rect = nav.getBoundingClientRect();
                // Disable sticky effect on desktop (lg breakpoint is 1024px)
                if (window.innerWidth >= 1024) {
                    setIsSticky(false);
                    return;
                }
                // Check sticky threshold based on screen size (matching CSS top values)
                const threshold = window.innerWidth >= 768 ? 106 : 128;
                setIsSticky(rect.top <= threshold);
            }
        };

        window.addEventListener("scroll", handleScroll);
        return () => window.removeEventListener("scroll", handleScroll);
    }, []);

    const toggleExpanded = () => {
        if (expanded) {
            setHeight("0px");
            setExpanded(false);
        } else {
            const scrollHeight = contentRef.current?.scrollHeight || 0;
            setHeight(`${scrollHeight}px`);
            setExpanded(true);
        }
    };

    const handleChange = (type: "delivery" | "pickup") => {
        setDeliveryType(type);
        setHeight("0px");
        setExpanded(false);
    };

    const deliveryLabel = deliveryType === "delivery" ? "Доставка" : "З собою";
    const DeliveryIcon = deliveryType === "delivery" ? Package : MapPin;

    const onCategoryClick = (e: React.MouseEvent<HTMLAnchorElement>, slug: string) => {
        e.preventDefault();
        const el = document.getElementById(slug);
        if (el) {
            el.scrollIntoView({ behavior: "smooth", block: "start" });
        }
    };

    const getIcon = (iconName: string): LucideIcon => {
        return getIconByName(iconName);
    };

    // Check if product has required variations
    const hasVariations = (product: Product): boolean => {
        return product.variation_groups && product.variation_groups.length > 0;
    };

    // Add to cart handler
    const handleAddToCart = (product: Product) => {
        if (hasVariations(product)) {
            // Open variation selection modal
            setSelectedProduct(product);
            setQuantity(1); // Reset quantity
            // Pre-select first variation for each group
            const initialSelections: Record<string, Variation> = {};
            product.variation_groups.forEach(group => {
                if (group.variations.length > 0) {
                    initialSelections[group.id] = group.variations[0];
                }
            });
            setTempSelectedVariations(initialSelections);
            setVariationModalOpen(true);
        } else {
            // Add directly to cart (1 item)
            addToCart(product, {}, 1);
            toast.success("Товар додано до кошика");
        }
    };

    // Confirm variation selection and add to cart
    const confirmVariationSelection = () => {
        if (selectedProduct) {
            addToCart(selectedProduct, tempSelectedVariations, quantity);
            setVariationModalOpen(false);
            toast.success("Товар додано до кошика");
        }
    };

    // Skeleton Component
    const MenuSkeleton = () => (
        <div className="flex-1 space-y-8 animate-pulse">
            {[1, 2, 3].map((categoryIndex) => (
                <div key={categoryIndex} className="space-y-4">
                    <Skeleton className="h-8 w-48" />
                    <div className="space-y-6">
                        {[1, 2, 3].map((productIndex) => (
                            <div key={productIndex} className="flex flex-col sm:flex-row gap-4 lg:gap-8 xl:gap-8 pb-3 border-b border-gray-800 border-dashed last:border-b-0">
                                <div className="flex flex-col gap-2 flex-1 order-2 sm:order-1">
                                    <Skeleton className="h-6 w-3/4" />
                                    <Skeleton className="h-4 w-20" />
                                    <Skeleton className="h-16 w-full" />
                                    <Skeleton className="h-4 w-24" />
                                </div>
                                <div className="order-1 sm:order-2 self-center sm:self-start w-full sm:w-auto">
                                    <Skeleton className="w-full sm:w-[250px] lg:w-[200px] xl:w-[250px] aspect-video sm:aspect-[4/3] rounded" />
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            ))}
        </div>
    );

    if (error) {
        return (
            <div>
                <Header />
                <main className="container mx-auto px-4 py-20 flex items-center justify-center min-h-[60vh]">
                    <div className="flex flex-col items-center gap-4">
                        <span className="text-lg text-red-500">{error}</span>
                        <Button onClick={() => window.location.reload()}>Спробувати ще раз</Button>
                    </div>
                </main>
                <Footer />
            </div>
        );
    }

    return (
        <div>
            <Header />

            {/* Shop Status Alert */}
            <ShopStatusAlert deliveryOpen={deliveryOpen} pickupOpen={pickupOpen} />

            <main className="container mx-auto px-4 pt-[25px] pb-[25px] flex flex-col lg:flex-row gap-4 lg:gap-4 xl:gap-10">
                {/* Mobile Info Accordion */}
                <div className="lg:hidden w-full mb-4">
                    <Accordion type="single" collapsible>
                        <AccordionItem value="info" className="border-b-0">
                            <AccordionTrigger
                                className="py-0 text-lg font-semibold"
                            >
                                Інформація про заклад
                            </AccordionTrigger>
                            <AccordionContent className="pt-2 pb-4">
                                <div className="flex flex-col gap-4">
                                    <div className="flex flex-row items-center gap-3">
                                        <Clock2 className="w-6 h-6 text-primary shrink-0" />
                                        <div className="flex flex-col">
                                            {hoursAreSame ? (
                                                <>
                                                    <span className="text-[15px] font-medium">Робочий час:</span>
                                                    {displayHours?.map((line, i) => (
                                                        <span key={i} className="text-[13px] text-gray-300">{line}</span>
                                                    ))}
                                                </>
                                            ) : (
                                                <>
                                                    <span className="text-[15px] font-medium">Доставка:</span>
                                                    {deliveryHours.map((line, i) => (
                                                        <span key={i} className="text-[13px] text-gray-300">{line}</span>
                                                    ))}
                                                    <span className="text-[15px] font-medium mt-2">Самовивіз:</span>
                                                    {pickupHours.map((line, i) => (
                                                        <span key={i} className="text-[13px] text-gray-300">{line}</span>
                                                    ))}
                                                </>
                                            )}
                                        </div>
                                    </div>

                                    <div className="flex flex-row items-center gap-3">
                                        <MapPin className="w-6 h-6 text-primary shrink-0" />
                                        <div className="flex flex-col">
                                            <span className="text-[15px]">Адреса:</span>
                                            <span className="text-[13px] text-gray-300">
                                                вул. Григорія Сковороди 64 (вхід з вул. Багалія), Харків, Харківська область, Україна
                                            </span>
                                        </div>
                                    </div>

                                    <div className="flex flex-row items-center gap-3">
                                        <Phone className="w-6 h-6 text-primary shrink-0" />
                                        <div className="flex flex-col">
                                            <span className="text-[15px]">Телефон:</span>
                                            <span className="text-[13px] text-gray-300">
                                                +38 (063) 500 95 97
                                            </span>
                                        </div>
                                    </div>
                                </div>
                            </AccordionContent>
                        </AccordionItem>
                    </Accordion>
                </div>

                <div className="contents lg:flex lg:flex-col gap-3 w-full lg:w-[230px] xl:w-[275px] lg:sticky lg:self-start lg:top-32 lg:z-20 pb-4 lg:pb-0 pt-2 lg:pt-0">
                    <nav
                        id="categories-nav"
                        className={`sticky top-[128px] md:top-[106px] z-30 py-2 lg:py-0 lg:static lg:bg-transparent flex flex-row lg:flex-col overflow-x-auto no-scrollbar gap-4 lg:gap-0 transition-all duration-300 ${isSticky
                            ? "bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 w-[100vw] ml-[calc(50%-50vw)] px-4 shadow-sm"
                            : "bg-transparent mx-0 px-0"
                            } lg:w-full lg:ml-0 lg:px-0`}
                    >
                        {categories.map((category) => {
                            const Icon = getIcon(category.icon);
                            return (
                                <a
                                    key={category.id}
                                    href={`#${category.slug}`}
                                    onClick={e => onCategoryClick(e, category.slug)}
                                    className={`flex items-center gap-2 mb-0 lg:mb-4 text-sm lg:text-lg cursor-pointer select-none whitespace-nowrap border lg:border-none px-3 py-2 rounded-full lg:px-0 lg:py-0 ${isSticky
                                        ? "bg-secondary/10 hover:bg-secondary/20"
                                        : "bg-secondary/10 hover:text-gray-300 lg:bg-transparent"
                                        }`}
                                >
                                    <Icon size={20} className="lg:w-[24px] lg:h-[24px]" />
                                    <span className={isSticky ? "lg:inline" : ""}>{category.name}</span>
                                </a>
                            );
                        })}
                    </nav>
                </div>

                <section className="flex-1 divide-y divide-white/10 border rounded-md p-4 mt-0">
                    {loading ? (
                        <MenuSkeleton />
                    ) : (
                        <div className="animate-in fade-in duration-1000 fill-mode-forwards divide-y divide-white/10">
                            {categories.map((category) => (
                                <div key={category.id} id={category.slug} className="scroll-mt-[182px] lg:scroll-mt-[130px] pt-6 first:pt-0">
                                    <h2 className="text-xl lg:text-2xl font-semibold mb-4">{category.name}</h2>
                                    <div className="">
                                        {category.products?.map((product) => (
                                            <ProductCard
                                                key={product.id}
                                                product={product}
                                                onAdd={(p) => handleAddToCart(p)}
                                            />
                                        ))}
                                    </div>
                                </div>
                            ))}
                        </div>
                    )}
                </section>

                <div className="hidden lg:flex flex-col gap-3 w-full lg:w-[230px] xl:w-[275px] sticky self-start top-auto lg:top-32 z-10 pb-10 lg:pb-0">
                    <div className="font-semibold text-lg lg:text-base">
                        Інформація про заклад
                    </div>

                    <div className="flex flex-row items-center gap-3">
                        <Clock2 className="w-6 h-6 text-primary shrink-0" />
                        <div className="flex flex-col">
                            {hoursAreSame ? (
                                <>
                                    <span className="text-[15px] font-medium">Робочий час:</span>
                                    {displayHours?.map((line, i) => (
                                        <span key={i} className="text-[13px] text-gray-300">{line}</span>
                                    ))}
                                </>
                            ) : (
                                <>
                                    <span className="text-[15px] font-medium">Доставка:</span>
                                    {deliveryHours.map((line, i) => (
                                        <span key={i} className="text-[13px] text-gray-300">{line}</span>
                                    ))}
                                    <span className="text-[15px] font-medium mt-2">Самовивіз:</span>
                                    {pickupHours.map((line, i) => (
                                        <span key={i} className="text-[13px] text-gray-300">{line}</span>
                                    ))}
                                </>
                            )}
                        </div>
                    </div>

                    <div className="flex flex-row items-center gap-3">
                        <MapPin className="w-6 h-6 text-primary shrink-0" />
                        <div className="flex flex-col">
                            <span className="text-[15px]">Адреса:</span>
                            <span className="text-[13px] text-gray-300">
                                вул. Григорія Сковороди 64 (вхід з вул. Багалія), Харків, Харківська область, Україна
                            </span>
                        </div>
                    </div>

                    <div className="flex flex-row items-center gap-3">
                        <Phone className="w-6 h-6 text-primary shrink-0" />
                        <div className="flex flex-col">
                            <span className="text-[15px]">Телефон:</span>
                            <span className="text-[13px] text-gray-300">
                                +38 (063) 500 95 97
                            </span>
                        </div>
                    </div>
                </div>
            </main>

            {/* Variation Selection Modal */}
            <div className={`fixed inset-0 z-[70] flex items-center justify-center transition-all duration-150 ${variationModalOpen ? "visible" : "invisible"}`}>
                {/* Backdrop with blur */}
                <div
                    className={`absolute inset-0 bg-black/70 backdrop-blur-sm transition-opacity duration-150 ${variationModalOpen ? "opacity-100" : "opacity-0"}`}
                    onClick={() => setVariationModalOpen(false)}
                />

                {/* Modal */}
                <div className={`relative bg-gradient-to-b from-[#1d1d1d] to-[#141414] border border-white/10 rounded-2xl p-6 max-w-md w-full mx-4 shadow-2xl transition duration-150 will-change-transform ${variationModalOpen ? "opacity-100 scale-100 translate-y-0" : "opacity-0 scale-95 translate-y-4"}`}>
                    {/* Close button */}
                    <button
                        onClick={() => setVariationModalOpen(false)}
                        className="absolute top-4 right-4 p-2 hover:bg-white/10 rounded-full transition-colors duration-100 cursor-pointer"
                    >
                        <X className="w-5 h-5" />
                    </button>

                    {selectedProduct && (
                        <>
                            {/* Product info */}
                            <div className="flex items-start gap-4 mb-6">
                                <div className="relative w-24 h-24 rounded-xl overflow-hidden shrink-0">
                                    <ImageWithFallback
                                        src={selectedProduct.image}
                                        alt={selectedProduct.name}
                                        className="w-full h-full object-cover"
                                    />
                                    <div className="absolute inset-0 bg-gradient-to-t from-black/40 to-transparent" />
                                </div>
                                <div>
                                    <h3 className="text-xl font-bold">{selectedProduct.name}</h3>
                                    <p className="text-2xl font-bold text-primary mt-1">{selectedProduct.price} ₴</p>
                                </div>
                            </div>

                            {/* Variation groups */}
                            <div className="space-y-5">
                                {selectedProduct.variation_groups.map((group) => (
                                    <div key={group.id}>
                                        <h4 className="text-sm font-medium text-gray-400 uppercase tracking-wider mb-3">
                                            {group.name}
                                            {group.required && <span className="text-red-400 ml-1">*</span>}
                                        </h4>
                                        <div className="flex flex-col gap-2">
                                            <RadioGroup
                                                value={tempSelectedVariations[group.id]?.id}
                                                onValueChange={(value) => {
                                                    const selectedVariation = group.variations.find(v => v.id === value);
                                                    if (selectedVariation) {
                                                        setTempSelectedVariations(prev => ({
                                                            ...prev,
                                                            [group.id]: selectedVariation
                                                        }));
                                                    }
                                                }}
                                                className="flex flex-col gap-3"
                                            >
                                                {group.variations.map((variation) => (
                                                    <div key={variation.id} className="flex items-center space-x-3">
                                                        <RadioGroupItem value={variation.id} id={variation.id} className="cursor-pointer" />
                                                        <label
                                                            htmlFor={variation.id}
                                                            className="text-base font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 cursor-pointer text-gray-200"
                                                        >
                                                            {variation.name}
                                                        </label>
                                                    </div>
                                                ))}
                                            </RadioGroup>
                                        </div>
                                    </div>
                                ))}
                            </div>

                            {/* Quantity Controls in Modal */}
                            <div className="flex items-center justify-center gap-4 my-6 py-4 border-t border-b border-white/5">
                                <button
                                    onClick={() => setQuantity(Math.max(1, quantity - 1))}
                                    className="w-8 h-8 flex items-center justify-center bg-white/10 hover:bg-white/20 rounded-full transition-colors cursor-pointer"
                                >
                                    <Minus className="w-4 h-4" />
                                </button>
                                <span className="text-xl font-bold w-8 text-center">{quantity}</span>
                                <button
                                    onClick={() => setQuantity(quantity + 1)}
                                    className="w-8 h-8 flex items-center justify-center bg-white/10 hover:bg-white/20 rounded-full transition-colors cursor-pointer"
                                >
                                    <Plus className="w-4 h-4" />
                                </button>
                            </div>

                            {/* Confirm button */}
                            <Button
                                onClick={confirmVariationSelection}
                                className="w-full h-14 text-lg font-bold text-black bg-primary hover:bg-primary/90 rounded-xl shadow-lg shadow-primary/25 transition duration-100 hover:shadow-primary/40 hover:scale-[1.02] cursor-pointer"
                            >
                                До кошику
                            </Button>
                        </>
                    )}
                </div>
            </div>

            <Footer />
        </div>
    );
}
