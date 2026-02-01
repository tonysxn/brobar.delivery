import { useEffect, useState } from "react";

export function useLocalStorage<T>(key: string, fallbackValue: T): [T, (value: T) => void] {
    const [value, setValue] = useState<T>(() => {
        if (typeof window === "undefined") return fallbackValue;
        try {
            const stored = localStorage.getItem(key);
            return stored ? JSON.parse(stored) : fallbackValue;
        } catch (error) {
            console.error("Storage access denied:", error);
            return fallbackValue;
        }
    });

    useEffect(() => {
        try {
            localStorage.setItem(key, JSON.stringify(value));
        } catch (error) {
            console.error("Storage access denied:", error);
        }
    }, [key, value]);

    return [value, setValue];
}
