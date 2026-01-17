"use client"

import { Toaster } from "sonner"
import { useIsMobile } from "@/hooks/use-mobile"

export function SiteToaster() {
    const isMobile = useIsMobile()

    return (
        <Toaster
            position={isMobile ? "bottom-center" : "bottom-right"}
            theme="dark"
            className="toaster group"
            toastOptions={{
                classNames: {
                    toast: "group toast group-[.toaster]:bg-black/90 group-[.toaster]:backdrop-blur-xl group-[.toaster]:text-white group-[.toaster]:border-white/10 group-[.toaster]:shadow-2xl group-[.toaster]:rounded-xl font-sans",
                    description: "group-[.toast]:text-gray-400",
                    actionButton: "group-[.toast]:bg-primary group-[.toast]:text-black font-bold",
                    cancelButton: "group-[.toast]:bg-white/10 group-[.toast]:text-white",
                    icon: "group-[.toast]:text-primary",
                },
                style: {
                    background: 'rgba(0, 0, 0, 0.9)',
                    backdropFilter: 'blur(12px)',
                    border: '1px solid rgba(255, 255, 255, 0.1)',
                }
            }}
        />
    )
}
