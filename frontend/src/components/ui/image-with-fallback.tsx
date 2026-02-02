"use client";

import { useState, useEffect } from "react";
import { getAssetUrl } from "@/lib/image-url";
import { ImageOff } from "lucide-react";
import { cn } from "@/lib/utils";

interface ImageWithFallbackProps extends Omit<React.ImgHTMLAttributes<HTMLImageElement>, 'src'> {
    src: string | undefined | null;
    fallbackClassName?: string;
    iconClassName?: string;
}

export function ImageWithFallback({
    src,
    alt,
    className,
    fallbackClassName,
    iconClassName,
    ...props
}: ImageWithFallbackProps) {
    const [error, setError] = useState(false);
    const [imageSrc, setImageSrc] = useState<string>("");

    useEffect(() => {
        if (!src) {
            setError(true);
            return;
        }

        setError(false);
        setImageSrc(getAssetUrl(src));
    }, [src]);

    if (error || !imageSrc) {
        return (
            <div
                className={cn(
                    "flex items-center justify-center bg-secondary/10 text-muted-foreground",
                    className,
                    fallbackClassName
                )}
            >
                <ImageOff className={cn("w-1/3 h-1/3 opacity-50", iconClassName)} />
            </div>
        );
    }

    return (
        // eslint-disable-next-line @next/next/no-img-element
        <img
            src={imageSrc}
            alt={alt || "Product image"}
            className={className}
            onError={() => setError(true)}
            {...props}
        />
    );
}
