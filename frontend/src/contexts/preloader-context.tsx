"use client";

import React, { createContext, useContext, useState, useEffect } from "react";

interface PreloaderContextType {
    isLoading: boolean;
    showLoader: () => void;
    hideLoader: () => void;
}

const PreloaderContext = createContext<PreloaderContextType | undefined>(undefined);

export function PreloaderProvider({ children }: { children: React.ReactNode }) {
    const [isLoading, setIsLoading] = useState(true);

    const showLoader = () => setIsLoading(true);
    const hideLoader = () => setIsLoading(false);

    useEffect(() => {
        const handleLoad = () => {
            setIsLoading(false);
        };

        if (document.readyState === "complete") {
            handleLoad();
        } else {
            window.addEventListener("load", handleLoad);
            return () => window.removeEventListener("load", handleLoad);
        }
    }, []);

    return (
        <PreloaderContext.Provider value={{ isLoading, showLoader, hideLoader }}>
            {children}
        </PreloaderContext.Provider>
    );
}

export function usePreloader() {
    const context = useContext(PreloaderContext);
    if (context === undefined) {
        throw new Error("usePreloader must be used within a PreloaderProvider");
    }
    return context;
}
