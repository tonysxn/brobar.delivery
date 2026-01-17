"use client"

import React from "react"
import {QueryClient, QueryClientProvider} from '@tanstack/react-query'
import {ThemeProvider} from "@/components/theme-provider"
import {Toaster} from "@/components/ui/sonner"
import "./globals.css"

const queryClient = new QueryClient()

export default function RootLayout({children}: { children: React.ReactNode }) {
    return (
        <html lang="en" suppressHydrationWarning>
        <head/>
        <body>
        <QueryClientProvider client={queryClient}>
            <ThemeProvider
                attribute="class"
                defaultTheme="system"
                enableSystem
                disableTransitionOnChange
            >
                {children}
            </ThemeProvider>
            <Toaster/>
        </QueryClientProvider>
        </body>
        </html>
    )
}
