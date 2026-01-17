"use client"

import * as React from "react"
import { cn } from "@/lib/utils"

interface CustomSwitchProps {
    checked: boolean
    onCheckedChange: (checked: boolean) => void
    label: string
}

export function CustomSwitch({
    checked,
    onCheckedChange,
    label
}: CustomSwitchProps) {
    return (
        <label className="flex items-center gap-3 cursor-pointer group">
            <div className="relative">
                <input
                    type="checkbox"
                    className="sr-only"
                    checked={checked}
                    onChange={(e) => onCheckedChange(e.target.checked)}
                />
                <div className={cn(
                    "block w-12 h-7 rounded-full transition-colors",
                    checked ? "bg-yellow-500" : "bg-gray-700"
                )} />
                <div className={cn(
                    "absolute left-1 top-1 bg-white w-5 h-5 rounded-full transition-transform",
                    checked ? "translate-x-5" : "translate-x-0"
                )} />
            </div>
            <span className="text-sm font-medium text-gray-300 group-hover:text-white transition-colors">
                {label}
            </span>
        </label>
    )
}
