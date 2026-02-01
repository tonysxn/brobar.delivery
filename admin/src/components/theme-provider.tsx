"use client"

import * as React from "react"
import { ThemeProvider as NextThemesProvider } from "next-themes"

export function ThemeProvider({
    children,
    ...props
}: React.ComponentProps<typeof NextThemesProvider>) {
    const [mounted, setMounted] = React.useState(false);

    React.useEffect(() => {
        try {
            // Check if storage is available
            localStorage.getItem("theme_test");
            localStorage.removeItem("theme_test");
            setMounted(true);
        } catch (e) {
            console.warn("Theme storage access denied, falling back to default theme");
        }
    }, []);

    if (!mounted) {
        return <>{children}</>;
    }

    return <NextThemesProvider {...props}>{children}</NextThemesProvider>
}