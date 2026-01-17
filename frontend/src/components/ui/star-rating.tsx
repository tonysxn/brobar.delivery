"use client"

import * as React from "react"
import { Star } from "lucide-react"
import { cn } from "@/lib/utils"

interface StarRatingProps {
    value: number
    onChange: (value: number) => void
    label?: string
    max?: number
}

export function StarRating({
    value,
    onChange,
    label,
    max = 5
}: StarRatingProps) {
    const [hover, setHover] = React.useState(0)

    return (
        <div className="flex flex-col gap-2">
            {label && <span className="text-sm font-medium text-gray-400 uppercase tracking-wider">{label}</span>}
            <div className="flex gap-1">
                {Array.from({ length: max }, (_, i) => i + 1).map((star) => (
                    <button
                        key={star}
                        type="button"
                        className="cursor-pointer transition-transform hover:scale-110 active:scale-95 outline-none"
                        onClick={() => onChange(star)}
                        onMouseEnter={() => setHover(star)}
                        onMouseLeave={() => setHover(0)}
                    >
                        <Star
                            className={cn(
                                "size-8 transition-colors",
                                (hover || value) >= star
                                    ? "fill-yellow-500 text-yellow-500"
                                    : "text-gray-600 fill-transparent"
                            )}
                        />
                    </button>
                ))}
            </div>
        </div>
    )
}
