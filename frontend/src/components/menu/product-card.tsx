"use client";

import { Product } from "@/contexts/cart-context";
import { Button } from "@/components/ui/button";
import { ImageWithFallback } from "@/components/ui/image-with-fallback";
import { cn } from "@/lib/utils";

interface ProductCardProps {
    product: Product;
    onAdd: (product: Product) => void;
    className?: string;
}

export function ProductCard({ product, onAdd, className }: ProductCardProps) {
    const hasVariations = product.variation_groups && product.variation_groups.length > 0;

    return (
        <div
            className={cn(
                "transition cursor-pointer my-5 border-b border-white/10 border-dashed last:border-b-0",
                className
            )}
        >
            <div className="flex flex-col sm:flex-row gap-4 lg:gap-8 xl:gap-8 pb-3">
                {/* Text Content */}
                <div className="flex flex-col gap-1 lg:gap-2 flex-1 order-2 sm:order-1">
                    <h3 className="text-lg font-medium">{product.name}</h3>
                    {/* Using a specific color class from globals or tailwind config if available, 
              otherwise falling back to a gold-ish color often used in the project */}
                    <p className="text-primary font-medium">{product.price} ₴</p>

                    {product.description && (
                        <p className="text-[13px] text-muted-foreground">
                            {product.description}
                        </p>
                    )}

                    {product.weight && (
                        <p className="text-[13px] text-gray-400">
                            <span className="font-semibold">Вага:</span>{" "}
                            {product.weight >= 1000
                                ? `${product.weight / 1000} кг`
                                : `${product.weight} г`}
                        </p>
                    )}
                </div>

                {/* Image */}
                <div className="order-1 sm:order-2 self-center sm:self-start w-full sm:w-auto">
                    <ImageWithFallback
                        src={product.image}
                        alt={product.name}
                        className="w-full sm:w-[250px] lg:w-[200px] xl:w-[250px] aspect-video sm:aspect-[4/3] object-cover rounded bg-white/5"
                    />
                </div>
            </div>

            {/* Action Button */}
            <div className="mb-3">
                <Button
                    className="w-fit text-black cursor-pointer bg-primary hover:bg-primary/90 font-medium"
                    onClick={(e) => {
                        e.stopPropagation(); // Prevent triggering parent clicks if any
                        onAdd(product);
                    }}
                    disabled={product.sold || (product.stock !== null && product.stock <= 0)}
                >
                    {(product.sold || (product.stock !== null && product.stock <= 0))
                        ? "Закінчився"
                        : hasVariations ? "Обрати" : "До кошика"}
                </Button>
                {product.stock !== null && product.stock > 0 && (
                    <span className="ml-3 text-sm text-yellow-500 font-medium">
                        В наявності: {product.stock} шт
                    </span>
                )}
            </div>
        </div>
    );
}
